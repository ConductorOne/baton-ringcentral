package connector

import (
	"context"
	"testing"

	v2 "github.com/conductorone/baton-sdk/pb/c1/connector/v2"
	"github.com/conductorone/baton-sdk/pkg/cli"
	rs "github.com/conductorone/baton-sdk/pkg/types/resource"
	"github.com/stretchr/testify/assert"
)

// TestUserBuilder_Grants_RoleSyncFilteredOut verifies that Grants short-circuits
// before touching b.client when role sync is disabled. b.client is left nil on
// purpose: if the guard were missing or wrong, dereferencing it would panic,
// so a clean (nil, nil, nil) return with no panic is proof the early return
// fired.
func TestUserBuilder_Grants_RoleSyncFilteredOut(t *testing.T) {
	b := &userBuilder{
		resourceType: userResourceType,
		client:       nil,
		syncRoles:    false,
	}

	userResource := &v2.Resource{
		Id: &v2.ResourceId{
			ResourceType: userResourceType.Id,
			Resource:     "user-1",
		},
	}

	grants, results, err := b.Grants(context.Background(), userResource, rs.SyncOpAttrs{})
	assert.NoError(t, err)
	assert.Nil(t, grants)
	assert.Nil(t, results)
}

// TestConnectorOpts_WillSyncResourceType_Role tests the gating logic that
// decides whether the user builder should emit role grants: it should follow
// cli.ConnectorOpts.WillSyncResourceType("role") exactly, since that's what
// NewLambdaConnector feeds into Connector.syncRoles / userBuilder.syncRoles.
func TestConnectorOpts_WillSyncResourceType_Role(t *testing.T) {
	tests := []struct {
		name                string
		syncResourceTypeIDs []string
		want                bool
	}{
		{
			name:                "no filter set syncs everything including role",
			syncResourceTypeIDs: nil,
			want:                true,
		},
		{
			name:                "role explicitly included",
			syncResourceTypeIDs: []string{"role"},
			want:                true,
		},
		{
			name:                "role excluded by explicit filter",
			syncResourceTypeIDs: []string{"user"},
			want:                false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			opts := &cli.ConnectorOpts{
				SyncResourceTypeIDs: tt.syncResourceTypeIDs,
			}
			got := opts.WillSyncResourceType(RoleResourceTypeID)
			assert.Equal(t, tt.want, got)
		})
	}
}
