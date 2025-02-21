package connector

import (
	"context"
	"fmt"

	"github.com/conductorone/baton-pulumi-cloud/pkg/client"
	v2 "github.com/conductorone/baton-sdk/pb/c1/connector/v2"
	"github.com/conductorone/baton-sdk/pkg/annotations"
	"github.com/conductorone/baton-sdk/pkg/pagination"
	batonResource "github.com/conductorone/baton-sdk/pkg/types/resource"
)

type userBuilder struct {
	client  *client.Client
	orgName string
}

func userResource(user *client.User, parentResourceId *v2.ResourceId) (*v2.Resource, error) {
	if user == nil {
		return nil, fmt.Errorf("user is nil")
	}

	profile := map[string]interface{}{
		"github_login": user.User.GithubLogin,
		"name":         user.User.Name,
		"role":         user.Role,
	}

	userStatus := v2.UserTrait_Status_STATUS_ENABLED
	accountType := v2.UserTrait_ACCOUNT_TYPE_HUMAN

	userTraits := []batonResource.UserTraitOption{
		batonResource.WithUserLogin(user.User.GithubLogin),
		batonResource.WithUserProfile(profile),
		batonResource.WithStatus(userStatus),
		batonResource.WithAccountType(accountType),
	}

	name := user.User.Name
	if name == "" {
		name = user.User.GithubLogin
	}

	return batonResource.NewUserResource(
		name,
		userResourceType,
		user.User.GithubLogin,
		userTraits,
		batonResource.WithParentResourceID(parentResourceId),
	)
}

func (ub *userBuilder) ResourceType(ctx context.Context) *v2.ResourceType {
	return userResourceType
}

// List returns all the users from the database as resource objects.
func (ub *userBuilder) List(ctx context.Context, parentResourceID *v2.ResourceId, pToken *pagination.Token) ([]*v2.Resource, string, annotations.Annotations, error) {
	var token string
	if pToken != nil {
		token = pToken.Token
	}

	// Create organization resource ID as parent
	orgParentID := &v2.ResourceId{
		ResourceType: orgResourceType.Id,
		Resource:     ub.orgName,
	}

	resp, err := ub.client.ListUsers(ctx, ub.orgName, token)
	if err != nil {
		return nil, "", nil, fmt.Errorf("failed to list users: %w", err)
	}

	resources := make([]*v2.Resource, 0, len(resp.Members))
	for _, member := range resp.Members {
		resource, err := userResource(&member, orgParentID)
		if err != nil {
			return nil, "", nil, err
		}
		resources = append(resources, resource)
	}

	return resources, resp.ContinuationToken, nil, nil
}

// Entitlements returns an empty list since users don't have their own entitlements
func (ub *userBuilder) Entitlements(_ context.Context, resource *v2.Resource, _ *pagination.Token) ([]*v2.Entitlement, string, annotations.Annotations, error) {
	return nil, "", nil, nil
}

// Grants returns an empty list since role grants are handled at the org level
func (ub *userBuilder) Grants(_ context.Context, resource *v2.Resource, _ *pagination.Token) ([]*v2.Grant, string, annotations.Annotations, error) {
	return nil, "", nil, nil
}

func newUserBuilder(client *client.Client, orgName string) *userBuilder {
	return &userBuilder{
		client:  client,
		orgName: orgName,
	}
}
