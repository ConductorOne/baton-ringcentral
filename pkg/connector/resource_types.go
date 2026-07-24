package connector

import (
	v2 "github.com/conductorone/baton-sdk/pb/c1/connector/v2"
)

// RoleResourceTypeID is the resource type ID for roles, used both to build
// roleResourceType below and to check whether role sync is enabled via
// cli.ConnectorOpts.WillSyncResourceType.
const RoleResourceTypeID = "role"

var userResourceType = &v2.ResourceType{
	Id:          "user",
	DisplayName: "User",
	Traits:      []v2.ResourceType_Trait{v2.ResourceType_TRAIT_USER},
}

var roleResourceType = &v2.ResourceType{
	Id:          RoleResourceTypeID,
	DisplayName: "Role",
	Traits:      []v2.ResourceType_Trait{v2.ResourceType_TRAIT_ROLE},
}
