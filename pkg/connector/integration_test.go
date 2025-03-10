package connector

import (
	"context"
	"fmt"
	"os"
	"testing"

	"github.com/conductorone/baton-ringcentral/pkg/client"
	v2 "github.com/conductorone/baton-sdk/pb/c1/connector/v2"
	"github.com/conductorone/baton-sdk/pkg/pagination"
	"github.com/stretchr/testify/assert"
)

var (
	ctx              = context.Background()
	rcClientID       = os.Getenv("RINGCENTRAL_CLIENT_ID")
	rcClientSecret   = os.Getenv("RINGCENTRAL_CLIENT_SECRET")
	rcJWT            = os.Getenv("RINGCENTRAL_JWT")
	parentResourceID = &v2.ResourceId{}

	message    string
	rolesCache []*v2.Resource
)

// TestUserBuilder_List tests the List function for User Resources.
func TestUserBuilder_List(t *testing.T) {
	if rcClientID == "" {
		t.Fatal("rcClientID env variable is required")
	}
	if rcClientSecret == "" {
		t.Fatal("rcClientSecret env variable is required")
	}
	if rcJWT == "" {
		t.Fatal("rcJWT env variable is required")
	}

	c, err := client.New(
		ctx,
		client.WithClientID(rcClientID),
		client.WithClientSecret(rcClientSecret),
		client.WithJWT(rcJWT),
	)
	if err != nil {
		message = fmt.Sprintf("error creating the client: %v", err)
		t.Fatal(message)
	}

	b := newUserBuilder(c)

	var users []*v2.Resource
	paginationToken := &pagination.Token{
		Size: 2, Token: "",
	}
	for {
		userResources, nextPageToken, _, err := b.List(ctx, parentResourceID, paginationToken)
		if err != nil {
			message = fmt.Sprintf("error listing users: %v", err)
			t.Fatal(message)
		}
		users = append(users, userResources...)
		if nextPageToken == "" {
			break
		}
		paginationToken.Token = nextPageToken
	}

	assert.NotNil(t, users)
}

// TestRoleBuilder_List tests the List function for Role Resources.
func TestRoleBuilder_List(t *testing.T) {
	if rcClientID == "" {
		t.Fatal("rcClientID env variable is required")
	}
	if rcClientSecret == "" {
		t.Fatal("rcClientSecret env variable is required")
	}
	if rcJWT == "" {
		t.Fatal("rcJWT env variable is required")
	}

	c, err := client.New(
		ctx,
		client.WithClientID(rcClientID),
		client.WithClientSecret(rcClientSecret),
		client.WithJWT(rcJWT),
	)
	if err != nil {
		message = fmt.Sprintf("error creating the client: %v", err)
		t.Fatal(message)
	}

	b := newRoleBuilder(c)

	var roles []*v2.Resource
	paginationToken := &pagination.Token{
		Size: 5, Token: "",
	}
	for {
		roleResources, nextPageToken, _, err := b.List(ctx, parentResourceID, paginationToken)
		if err != nil {
			message = fmt.Sprintf("error listing roles: %v", err)
			t.Fatal(message)
		}
		roles = append(roles, roleResources...)
		if nextPageToken == "" {
			break
		}
		paginationToken.Token = nextPageToken
	}

	assert.NotNil(t, roles)
	rolesCache = append(rolesCache, roles...)
}

// TestRoleBuilder_Entitlements tests the Entitlements function for Role Resources using the previously requested Roles on the TestRoleBuilder_List test.
func TestRoleBuilder_Entitlements(t *testing.T) {
	if len(rolesCache) == 0 {
		t.Skip("role entitlements test will be skipped since the rolesCache is empty")
	}

	var entitlements []*v2.Entitlement
	c, _ := client.New(ctx)
	b := newRoleBuilder(c)

	for _, role := range rolesCache {
		entitlementResource, _, _, err := b.Entitlements(ctx, role, nil)
		if err != nil {
			message = fmt.Sprintf("error creating entitlement: %v", err)
			t.Fatal(message)
		}

		entitlements = append(entitlements, entitlementResource...)
	}

	assert.NotNil(t, entitlements)
}
