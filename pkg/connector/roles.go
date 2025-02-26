package connector

import (
	"context"
	"fmt"

	"github.com/conductorone/baton-ringcentral/pkg/client"
	v2 "github.com/conductorone/baton-sdk/pb/c1/connector/v2"
	"github.com/conductorone/baton-sdk/pkg/annotations"
	"github.com/conductorone/baton-sdk/pkg/pagination"
	"github.com/conductorone/baton-sdk/pkg/types/entitlement"
	rs "github.com/conductorone/baton-sdk/pkg/types/resource"
	"github.com/grpc-ecosystem/go-grpc-middleware/logging/zap/ctxzap"
	"go.uber.org/zap"
)

const rolePermissionName = "assigned"

type roleBuilder struct {
	client       *client.RingCentralClient
	resourceType *v2.ResourceType
}

func (b *roleBuilder) ResourceType(_ context.Context) *v2.ResourceType {
	return roleResourceType
}

func (b *roleBuilder) List(ctx context.Context, _ *v2.ResourceId, pToken *pagination.Token) ([]*v2.Resource, string, annotations.Annotations, error) {
	var roleResources []*v2.Resource

	bag, pageToken, err := getToken(pToken, userResourceType)
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

	err = bag.Next(nextPageToken)
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

	nextPageToken, err = bag.Marshal()
	if err != nil {
		return nil, "", nil, err
	}

	return roleResources, nextPageToken, nil, nil
}

func (b *roleBuilder) Entitlements(_ context.Context, resource *v2.Resource, _ *pagination.Token) ([]*v2.Entitlement, string, annotations.Annotations, error) {
	var roleEntitlements []*v2.Entitlement

	assigmentOptions := []entitlement.EntitlementOption{
		entitlement.WithGrantableTo(userResourceType),
		entitlement.WithDescription(resource.Description),
		entitlement.WithDisplayName(resource.DisplayName),
	}

	roleEntitlements = append(roleEntitlements, entitlement.NewPermissionEntitlement(resource, rolePermissionName, assigmentOptions...))

	return roleEntitlements, "", nil, nil
}

// Grants function isn't implemented here because they are build in the Grants function of the Roles.
// This was made like this since it was convenient considering the data model of the platform.
func (b *roleBuilder) Grants(_ context.Context, _ *v2.Resource, _ *pagination.Token) ([]*v2.Grant, string, annotations.Annotations, error) {
	return nil, "", nil, nil
}

func (b *roleBuilder) Grant(ctx context.Context, principal *v2.Resource, entitlement *v2.Entitlement) (annotations.Annotations, error) {
	l := ctxzap.Extract(ctx)
	if principal.Id.ResourceType != userResourceType.Id {
		l.Warn("ringcentral-connector: only users can be granted with role membership",
			zap.String("principal_id", principal.Id.Resource),
			zap.String("principal_type", principal.Id.Resource))
		return nil, fmt.Errorf("ringcentral-connector: only users can be granted with role membership")
	}

	roleID := entitlement.Resource.Id.Resource
	err := b.client.UpdateUserRoles(ctx, principal, roleID, false)
	if err != nil {
		return nil, err
	}

	return nil, nil
}

func (b *roleBuilder) Revoke(ctx context.Context, grant *v2.Grant) (annotations.Annotations, error) {
	roleID := grant.Entitlement.Resource.Id.Resource

	err := b.client.UpdateUserRoles(ctx, grant.Principal, roleID, true)
	if err != nil {
		return nil, err
	}

	return nil, nil
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
