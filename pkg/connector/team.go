package connector

import (
	"context"
	"fmt"

	"github.com/conductorone/baton-pulumi-cloud/pkg/client"
	v2 "github.com/conductorone/baton-sdk/pb/c1/connector/v2"
	"github.com/conductorone/baton-sdk/pkg/annotations"
	"github.com/conductorone/baton-sdk/pkg/connectorbuilder"
	"github.com/conductorone/baton-sdk/pkg/pagination"
	batonEntitlement "github.com/conductorone/baton-sdk/pkg/types/entitlement"
	batonGrant "github.com/conductorone/baton-sdk/pkg/types/grant"
	batonResource "github.com/conductorone/baton-sdk/pkg/types/resource"
)

type teamBuilder struct {
	resourceType *v2.ResourceType
	client       *client.Client
	orgName      string
}

var _ connectorbuilder.ResourceSyncer = &teamBuilder{}
var _ connectorbuilder.ResourceProvisionerV2 = &teamBuilder{}

func teamResource(team client.Team, parentResourceId *v2.ResourceId) (*v2.Resource, error) {
	profile := map[string]interface{}{
		"name":         team.Name,
		"display_name": team.DisplayName,
		"description":  team.Description,
	}

	return batonResource.NewGroupResource(
		team.DisplayName,
		teamResourceType,
		team.Name,
		[]batonResource.GroupTraitOption{
			batonResource.WithGroupProfile(profile),
		},
		batonResource.WithParentResourceID(parentResourceId),
	)
}

func (o *teamBuilder) ResourceType(ctx context.Context) *v2.ResourceType {
	return teamResourceType
}

// List returns all the teams from Pulumi as resource objects.
func (o *teamBuilder) List(ctx context.Context, parentResourceID *v2.ResourceId, _ *pagination.Token) ([]*v2.Resource, string, annotations.Annotations, error) {
	teams := []*v2.Resource{}
	var annotations annotations.Annotations

	// Create the parent org resource ID if not provided
	if parentResourceID == nil {
		parentResourceID = &v2.ResourceId{
			ResourceType: orgResourceType.Id,
			Resource:     o.orgName,
		}
	}

	resp, err := o.client.ListTeams(ctx, o.orgName)
	if err != nil {
		return nil, "", annotations, fmt.Errorf("failed to list teams: %w", err)
	}

	for _, team := range resp {
		teamResource, err := teamResource(team, parentResourceID)
		if err != nil {
			return nil, "", annotations, err
		}
		teams = append(teams, teamResource)
	}

	// ListTeams doesn't support pagination
	return teams, "", annotations, nil
}

// Entitlements returns the entitlements available for a team.
func (o *teamBuilder) Entitlements(_ context.Context, resource *v2.Resource, _ *pagination.Token) ([]*v2.Entitlement, string, annotations.Annotations, error) {
	memberEnt := batonEntitlement.NewAssignmentEntitlement(
		resource,
		"member",
		batonEntitlement.WithGrantableTo(userResourceType),
		batonEntitlement.WithDescription("Member of the team"),
		batonEntitlement.WithDisplayName("Member"),
	)

	return []*v2.Entitlement{memberEnt}, "", nil, nil
}

// Grants returns the granted entitlements for users in the team.
func (o *teamBuilder) Grants(ctx context.Context, resource *v2.Resource, _ *pagination.Token) ([]*v2.Grant, string, annotations.Annotations, error) {
	var rv []*v2.Grant
	var annotations annotations.Annotations

	// Get team details including members
	team, err := o.client.GetTeam(ctx, o.orgName, resource.Id.Resource)
	if err != nil {
		return nil, "", annotations, fmt.Errorf("failed to get team: %w", err)
	}

	for _, member := range team.Members {
		grant := batonGrant.NewGrant(
			resource,
			"member",
			&v2.ResourceId{
				ResourceType: userResourceType.Id,
				Resource:     member.GithubLogin,
			},
		)

		rv = append(rv, grant)
	}

	return rv, "", annotations, nil
}

// Grant implements the entitlement grant operation
func (o *teamBuilder) Grant(ctx context.Context, principal *v2.Resource, entitlement *v2.Entitlement) ([]*v2.Grant, annotations.Annotations, error) {
	if principal == nil || principal.Id == nil {
		return nil, nil, fmt.Errorf("principal is nil or has nil id")
	}
	if entitlement == nil || entitlement.Resource == nil || entitlement.Resource.Id == nil {
		return nil, nil, fmt.Errorf("entitlement is nil or has nil resource")
	}
	if principal.Id.ResourceType != userResourceType.Id {
		return nil, nil, fmt.Errorf("cannot grant team membership to non-user resource type: %s", principal.Id.ResourceType)
	}

	err := o.client.UpdateTeamMembership(ctx, o.orgName, entitlement.Resource.Id.Resource, principal.Id.Resource, "add")
	if err != nil {
		return nil, nil, fmt.Errorf("failed to add team member: %w", err)
	}

	return nil, nil, nil
}

// Revoke implements the entitlement revoke operation
func (o *teamBuilder) Revoke(ctx context.Context, grant *v2.Grant) (annotations.Annotations, error) {
	if grant == nil || grant.Principal == nil || grant.Principal.Id == nil {
		return nil, fmt.Errorf("grant is nil or has nil principal")
	}
	if grant.Entitlement == nil || grant.Entitlement.Resource == nil || grant.Entitlement.Resource.Id == nil {
		return nil, fmt.Errorf("grant has nil entitlement or resource")
	}

	err := o.client.UpdateTeamMembership(ctx, o.orgName, grant.Entitlement.Resource.Id.Resource, grant.Principal.Id.Resource, "remove")
	if err != nil {
		return nil, fmt.Errorf("failed to remove team member: %w", err)
	}

	return nil, nil
}

func newTeamBuilder(client *client.Client, orgName string) *teamBuilder {
	return &teamBuilder{
		resourceType: teamResourceType,
		client:       client,
		orgName:      orgName,
	}
}
