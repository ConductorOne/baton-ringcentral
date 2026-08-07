package connector

import (
	"context"
	"io"

	"github.com/conductorone/baton-ringcentral/pkg/client"
	cfg "github.com/conductorone/baton-ringcentral/pkg/config"

	v2 "github.com/conductorone/baton-sdk/pb/c1/connector/v2"
	"github.com/conductorone/baton-sdk/pkg/annotations"
	"github.com/conductorone/baton-sdk/pkg/cli"
	"github.com/conductorone/baton-sdk/pkg/connectorbuilder"
)

type Connector struct {
	client *client.RingCentralClient
	// skipRoleResourceType reports whether role is excluded from the sync
	// filter. Named for the skip condition so the zero value is safe: main.go
	// registers a zero-value Connector{} as the capabilities factory.
	skipRoleResourceType bool
}

// ResourceSyncers returns a ResourceSyncerV2 for each resource type that should be synced from the upstream service.
func (d *Connector) ResourceSyncers(_ context.Context) []connectorbuilder.ResourceSyncerV2 {
	return []connectorbuilder.ResourceSyncerV2{
		newUserBuilder(d.client, d.skipRoleResourceType),
		newRoleBuilder(d.client),
	}
}

// Asset takes an input AssetRef and attempts to fetch it using the connector's authenticated http client
// It streams a response, always starting with a metadata object, following by chunked payloads for the asset.
func (d *Connector) Asset(_ context.Context, _ *v2.AssetRef) (string, io.ReadCloser, error) {
	return "", nil, nil
}

// Metadata returns metadata about the connector.
func (d *Connector) Metadata(_ context.Context) (*v2.ConnectorMetadata, error) {
	return &v2.ConnectorMetadata{
		DisplayName: "Baton RingCentral Connector",
		Description: "Connector to sync users and permissions data from RingCentral. It allows the grant and revoke of user roles.",
	}, nil
}

// Validate is called to ensure that the connector is properly configured. It should exercise any API credentials
// to be sure that they are valid.
func (d *Connector) Validate(_ context.Context) (annotations.Annotations, error) {
	return nil, nil
}

// NewConnector returns a new instance of the connector.
func NewConnector(ctx context.Context, ac *cfg.Ringcentral, skipRoleResourceType bool) (connectorbuilder.ConnectorBuilderV2, []connectorbuilder.Opt, error) {
	c, err := client.New(
		ctx,
		client.WithClientID(ac.RingcentralClientId),
		client.WithClientSecret(ac.RingcentralClientSecret),
		client.WithJWT(ac.RingcentralJwt),
	)
	if err != nil {
		return nil, nil, err
	}

	return &Connector{
		client:               c,
		skipRoleResourceType: skipRoleResourceType,
	}, nil, nil
}

// NewLambdaConnector returns a new instance of the connector for Lambda deployments.
func NewLambdaConnector(ctx context.Context, ac *cfg.Ringcentral, opts *cli.ConnectorOpts) (connectorbuilder.ConnectorBuilderV2, []connectorbuilder.Opt, error) {
	// nil opts means no filter, so nothing is skipped.
	return NewConnector(ctx, ac, opts != nil && !opts.WillSyncResourceType(RoleResourceTypeID))
}
