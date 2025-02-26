package connector

import (
	"context"
	"github.com/conductorone/baton-ringcentral/pkg/client"
	v2 "github.com/conductorone/baton-sdk/pb/c1/connector/v2"
	"github.com/conductorone/baton-sdk/pkg/annotations"
	"github.com/conductorone/baton-sdk/pkg/pagination"
	"github.com/conductorone/baton-sdk/pkg/types/grant"
	rs "github.com/conductorone/baton-sdk/pkg/types/resource"
)

type userBuilder struct {
	resourceType *v2.ResourceType
	client       *client.RingCentralClient
}

func (b *userBuilder) ResourceType(_ context.Context) *v2.ResourceType {
	return userResourceType
}

// List returns all the users from the database as resource objects.
// Users include a UserTrait because they are the 'shape' of a standard user.
func (b *userBuilder) List(ctx context.Context, _ *v2.ResourceId, pToken *pagination.Token) ([]*v2.Resource, string, annotations.Annotations, error) {
	var userResources []*v2.Resource

	bag, pageToken, err := getToken(pToken, userResourceType)
	if err != nil {
		return nil, "", nil, err
	}

	users, nextPageToken, err := b.client.ListAllUsers(ctx, client.PageOptions{
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

	for _, user := range users {
		userResource, err := parseIntoUserResource(user)
		if err != nil {
			return nil, "", nil, err
		}

		userResources = append(userResources, userResource)
	}

	nextPageToken, err = bag.Marshal()
	if err != nil {
		return nil, "", nil, err
	}

	return userResources, nextPageToken, nil, nil
}

// Entitlements always returns an empty slice for users.
func (b *userBuilder) Entitlements(_ context.Context, _ *v2.Resource, _ *pagination.Token) ([]*v2.Entitlement, string, annotations.Annotations, error) {
	return nil, "", nil, nil
}

// Grants always returns an empty slice for users since they don't have any entitlements.
// In this case, Grants will create the Role Grants, since the Roles assigned are an internal data of each user that should be requested using the User ID.
func (b *userBuilder) Grants(ctx context.Context, userResource *v2.Resource, _ *pagination.Token) ([]*v2.Grant, string, annotations.Annotations, error) {
	var roleGrants []*v2.Grant

	userRoles, err := b.client.GetUserAssignedRoles(ctx, userResource)
	if err != nil {
		return nil, "", nil, err
	}

	for _, userRole := range userRoles {
		roleResource := &v2.Resource{
			Id: &v2.ResourceId{
				ResourceType: roleResourceType.Id,
				Resource:     userRole.Id,
			},
		}
		roleGrants = append(roleGrants, grant.NewGrant(roleResource, rolePermissionName, userResource))
	}

	return roleGrants, "", nil, nil
}

// parseIntoUserResource - This function parses an Extension (users from RingCentral) into a User Resource.
func parseIntoUserResource(extension client.Extension) (*v2.Resource, error) {
	var userStatus = v2.UserTrait_Status_STATUS_ENABLED

	profile := map[string]interface{}{
		"user_id":    extension.ID,
		"email":      extension.ContactInfo.Email,
		"first_name": extension.ContactInfo.FirstName,
		"last_name":  extension.ContactInfo.LastName,
		"status":     extension.Status,
	}

	userTraits := []rs.UserTraitOption{
		rs.WithUserProfile(profile),
		rs.WithStatus(userStatus),
		rs.WithUserLogin(extension.ContactInfo.Email),
		rs.WithEmail(extension.ContactInfo.Email, true),
	}

	displayName := extension.Name
	if displayName == "" {
		displayName = extension.ContactInfo.Email
	}

	ret, err := rs.NewUserResource(
		displayName,
		userResourceType,
		extension.ID,
		userTraits,
	)
	if err != nil {
		return nil, err
	}

	return ret, nil
}

func newUserBuilder(c *client.RingCentralClient) *userBuilder {
	return &userBuilder{
		resourceType: userResourceType,
		client:       c,
	}
}
