package connector

import (
	"context"
	"github.com/conductorone/baton-ringcentral/pkg/client"
	v2 "github.com/conductorone/baton-sdk/pb/c1/connector/v2"
	"github.com/conductorone/baton-sdk/pkg/annotations"
	"github.com/conductorone/baton-sdk/pkg/pagination"
	"github.com/conductorone/baton-sdk/pkg/types/entitlement"
	//"github.com/conductorone/baton-sdk/pkg/types/grant"
	rs "github.com/conductorone/baton-sdk/pkg/types/resource"
)

type roleBuilder struct {
	client       *client.RingCentralClient
	resourceType *v2.ResourceType
}

func (b *roleBuilder) ResourceType(_ context.Context) *v2.ResourceType {
	return roleResourceType
}

func (b *roleBuilder) List(ctx context.Context, _ *v2.ResourceId, pToken *pagination.Token) ([]*v2.Resource, string, annotations.Annotations, error) {
	var roleResources []*v2.Resource
	_, pageToken, err := getToken(pToken, userResourceType)
	if err != nil {
		return nil, "", nil, err
	}

	roles, nextPageToken, err := b.client.ListAllAvailableRoles(ctx, client.PageOptions{
		Page:    pageToken,
		PerPage: pToken.Size,
	})
	if err != nil {
		return nil, "", nil, err
	}

	for _, role := range roles {
		roleResource, err := parseIntoRoleResource(role)
		if err != nil {
			return nil, "", nil, err
		}

		roleResources = append(roleResources, roleResource)
	}

	//err = bag.Next(nextPageToken)
	//if err != nil {
	//	return nil, "", nil, err
	//}

	return roleResources, nextPageToken, nil, nil
}

func (b *roleBuilder) Entitlements(ctx context.Context, resource *v2.Resource, pToken *pagination.Token) ([]*v2.Entitlement, string, annotations.Annotations, error) {
	var roleEntitlements []*v2.Entitlement
	permissionName := "assigned"

	assigmentOptions := []entitlement.EntitlementOption{
		entitlement.WithGrantableTo(userResourceType),
		entitlement.WithDescription(resource.Description),
		entitlement.WithDisplayName(resource.DisplayName),
	}

	roleEntitlements = append(roleEntitlements, entitlement.NewPermissionEntitlement(resource, permissionName, assigmentOptions...))

	return roleEntitlements, "", nil, nil
}

func (b *roleBuilder) Grants(ctx context.Context, resource *v2.Resource, pToken *pagination.Token) ([]*v2.Grant, string, annotations.Annotations, error) {

	return nil, "", nil, nil
}

func parseIntoRoleResource(role client.Role) (*v2.Resource, error) {
	profile := map[string]interface{}{
		"role_id":      role.Id,
		"description":  role.Description,
		"display_name": role.DisplayName,
		"scope":        role.Scope,
		"hidden":       role.Hidden,
		"custom":       role.Custom,
	}

	roleTraits := []rs.RoleTraitOption{
		rs.WithRoleProfile(profile),
	}

	ret, err := rs.NewRoleResource(
		role.DisplayName,
		roleResourceType,
		role.Id,
		roleTraits,
	)

	if err != nil {
		return nil, err
	}

	return ret, nil
}

func newRoleBuilder(c *client.RingCentralClient) *roleBuilder {
	return &roleBuilder{
		resourceType: roleResourceType,
		client:       c,
	}
}
