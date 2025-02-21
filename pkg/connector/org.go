package connector

import (
	"context"
	"fmt"
	"strings"

	"github.com/conductorone/baton-pulumi-cloud/pkg/client"
	v2 "github.com/conductorone/baton-sdk/pb/c1/connector/v2"
	"github.com/conductorone/baton-sdk/pkg/annotations"
	"github.com/conductorone/baton-sdk/pkg/connectorbuilder"
	"github.com/conductorone/baton-sdk/pkg/pagination"
	batonEntitlement "github.com/conductorone/baton-sdk/pkg/types/entitlement"
	batonGrant "github.com/conductorone/baton-sdk/pkg/types/grant"
	batonResource "github.com/conductorone/baton-sdk/pkg/types/resource"
)

const (
	entitlementSlugAdmin  = "admin"
	entitlementSlugMember = "member"
	roleAdmin             = "admin"
	roleMember            = "member"
)

type orgBuilder struct {
	resourceType *v2.ResourceType
	client       *client.Client
	orgName      string
}

var _ connectorbuilder.ResourceSyncer = &orgBuilder{}
var _ connectorbuilder.ResourceProvisionerV2 = &orgBuilder{}

func orgResource(orgName string) (*v2.Resource, error) {
	return batonResource.NewResource(
		orgName,
		orgResourceType,
		orgName,
	)
}

func (o *orgBuilder) ResourceType(ctx context.Context) *v2.ResourceType {
	return orgResourceType
}

// List returns a single organization resource.
func (o *orgBuilder) List(ctx context.Context, parentResourceID *v2.ResourceId, _ *pagination.Token) ([]*v2.Resource, string, annotations.Annotations, error) {
	resource, err := orgResource(o.orgName)
	if err != nil {
		return nil, "", nil, fmt.Errorf("failed to create org resource: %w", err)
	}

	return []*v2.Resource{resource}, "", nil, nil
}

// formatResourceID creates a structured resource ID
func formatResourceID(parts ...string) string {
	return strings.Join(parts, ":")
}

// Entitlements returns the entitlements available for the organization.
func (o *orgBuilder) Entitlements(_ context.Context, resource *v2.Resource, _ *pagination.Token) ([]*v2.Entitlement, string, annotations.Annotations, error) {
	memberEnt := batonEntitlement.NewAssignmentEntitlement(
		resource,
		entitlementSlugMember,
		batonEntitlement.WithGrantableTo(userResourceType),
		batonEntitlement.WithDescription("Member of the Pulumi organization"),
		batonEntitlement.WithDisplayName("Member"),
	)

	adminEnt := batonEntitlement.NewPermissionEntitlement(
		resource,
		entitlementSlugAdmin,
		batonEntitlement.WithGrantableTo(userResourceType),
		batonEntitlement.WithDescription("Administrator of the Pulumi organization"),
		batonEntitlement.WithDisplayName("Administrator"),
	)

	return []*v2.Entitlement{memberEnt, adminEnt}, "", nil, nil
}

// Grants returns the granted entitlements for users in the organization.
func (o *orgBuilder) Grants(ctx context.Context, resource *v2.Resource, pToken *pagination.Token) ([]*v2.Grant, string, annotations.Annotations, error) {
	var rv []*v2.Grant
	var annotations annotations.Annotations

	var token string
	if pToken != nil {
		token = pToken.Token
	}

	// Get organization members
	resp, err := o.client.ListUsers(ctx, o.orgName, token)
	if err != nil {
		return nil, "", annotations, fmt.Errorf("failed to list org members: %w", err)
	}

	for _, member := range resp.Members {
		// Create grant for the appropriate role (admin or member)
		entSlug := entitlementSlugMember
		if member.Role == roleAdmin {
			entSlug = entitlementSlugAdmin
		}

		principalId := &v2.ResourceId{
			ResourceType: userResourceType.Id,
			Resource:     member.User.GithubLogin,
		}

		g := batonGrant.NewGrant(
			resource,
			entSlug,
			principalId,
		)
		g.Id = formatResourceID("org", o.orgName, "grant", member.User.GithubLogin, entSlug)
		g.Principal.DisplayName = member.User.Name

		rv = append(rv, g)
	}

	return rv, resp.ContinuationToken, annotations, nil
}

// Grant implements the entitlement grant operation
func (o *orgBuilder) Grant(ctx context.Context, principal *v2.Resource, entitlement *v2.Entitlement) ([]*v2.Grant, annotations.Annotations, error) {
	if principal == nil || principal.Id == nil {
		return nil, nil, fmt.Errorf("principal is nil or has nil id")
	}
	if entitlement == nil {
		return nil, nil, fmt.Errorf("entitlement is nil")
	}

	// Only users can be granted org roles
	if principal.Id.ResourceType != userResourceType.Id {
		return nil, nil, fmt.Errorf("cannot grant org role to non-user resource type: %s", principal.Id.ResourceType)
	}

	// Match against the known entitlement IDs we create
	adminEntID := formatResourceID(orgResourceType.Id, o.orgName, entitlementSlugAdmin)
	memberEntID := formatResourceID(orgResourceType.Id, o.orgName, entitlementSlugMember)

	var role string
	switch entitlement.Id {
	case adminEntID:
		role = roleAdmin
	case memberEntID:
		role = roleMember
	default:
		return nil, nil, fmt.Errorf("unknown entitlement ID: %s", entitlement.Id)
	}

	// Update the user's role in the organization
	return nil, nil, o.client.UpdateUserRole(ctx, o.orgName, principal.Id.Resource, role)
}

// Revoke implements the entitlement revoke operation
func (o *orgBuilder) Revoke(ctx context.Context, grant *v2.Grant) (annotations.Annotations, error) {
	if grant == nil {
		return nil, fmt.Errorf("grant is nil")
	}
	if grant.Principal == nil || grant.Principal.Id == nil {
		return nil, fmt.Errorf("grant principal is nil or has nil id")
	}
	if grant.Entitlement == nil {
		return nil, fmt.Errorf("grant entitlement is nil")
	}

	// Only users can have org roles revoked
	if grant.Principal.Id.ResourceType != userResourceType.Id {
		return nil, fmt.Errorf("cannot revoke org role from non-user resource type: %s", grant.Principal.Id.ResourceType)
	}

	username := grant.Principal.Id.Resource

	// Match against the known grant ID patterns we create
	adminGrantID := formatResourceID("org", o.orgName, "grant", username, entitlementSlugAdmin)
	memberGrantID := formatResourceID("org", o.orgName, "grant", username, entitlementSlugMember)

	switch grant.Id {
	case adminGrantID:
		// When admin is revoked, downgrade to member
		return nil, o.client.UpdateUserRole(ctx, o.orgName, username, roleMember)
	case memberGrantID:
		// When member is revoked, remove from org
		return nil, o.client.RemoveUser(ctx, o.orgName, username)
	default:
		return nil, fmt.Errorf("unknown grant ID: %s", grant.Id)
	}
}

func newOrgBuilder(client *client.Client, orgName string) *orgBuilder {
	return &orgBuilder{
		resourceType: orgResourceType,
		client:       client,
		orgName:      orgName,
	}
}
