package connector

import (
	"context"
	"testing"

	v2 "github.com/conductorone/baton-sdk/pb/c1/connector/v2"
	"github.com/conductorone/baton-sdk/pkg/annotations"
	"github.com/conductorone/baton-sdk/pkg/cli"
	"github.com/stretchr/testify/assert"
)

// TestUserBuilder_ResourceType_SkipAnnotations verifies that ResourceType
// annotates the returned (cloned) resource type based on skipRoleGrants:
//   - skipRoleGrants=false (role sync enabled): SkipEntitlements only, since
//     Entitlements() always returns nil regardless of role-filter state.
//   - skipRoleGrants=true (role sync filtered out): SkipEntitlementsAndGrants,
//     since nothing meaningful would be emitted at all.
//
// It also asserts the shared package-level userResourceType var is never
// mutated by ResourceType, since other code paths read it directly.
func TestUserBuilder_ResourceType_SkipAnnotations(t *testing.T) {
	t.Run("role sync enabled: SkipEntitlements only", func(t *testing.T) {
		b := &userBuilder{resourceType: userResourceType, skipRoleGrants: false}

		rt := b.ResourceType(context.Background())

		annos := annotations.Annotations(rt.GetAnnotations())
		assert.True(t, annos.Contains(&v2.SkipEntitlements{}))
		assert.False(t, annos.Contains(&v2.SkipEntitlementsAndGrants{}))
	})

	t.Run("role sync filtered out: SkipEntitlementsAndGrants", func(t *testing.T) {
		b := &userBuilder{resourceType: userResourceType, skipRoleGrants: true}

		rt := b.ResourceType(context.Background())

		annos := annotations.Annotations(rt.GetAnnotations())
		assert.True(t, annos.Contains(&v2.SkipEntitlementsAndGrants{}))
	})

	assert.Len(t, userResourceType.Annotations, 0, "ResourceType must not mutate the shared package-level userResourceType var")
}

// TestConnectorOpts_WillSyncResourceType_Role tests the gating logic that
// decides whether the user builder should emit role grants: it should follow
// cli.ConnectorOpts.WillSyncResourceType("role") exactly, since that's what
// NewLambdaConnector feeds into Connector.skipRoleGrants / userBuilder.skipRoleGrants
// (inverted: skipRoleGrants = !WillSyncResourceType(RoleResourceTypeID)).
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
