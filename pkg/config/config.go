package config

import (
	"github.com/conductorone/baton-sdk/pkg/field"
)

var (
	clientIDField = field.StringField(
		"ringcentral-client-id",
		field.WithDisplayName("RingCentral Client ID"),
		field.WithDescription("Client ID of the Baton App for RingCentral"),
		field.WithRequired(true),
	)

	clientSecretField = field.StringField(
		"ringcentral-client-secret",
		field.WithDisplayName("RingCentral Client Secret"),
		field.WithDescription("Client Secret of the Baton App for RingCentral"),
		field.WithIsSecret(true),
		field.WithRequired(true),
	)

	jwtField = field.StringField(
		"ringcentral-jwt",
		field.WithDisplayName("RingCentral JWT"),
		field.WithDescription("JWT of the admin user on RingCentral platform"),
		field.WithIsSecret(true),
		field.WithRequired(true),
	)

	// ConfigurationFields defines the external configuration required for the
	// connector to run.
	ConfigurationFields = []field.SchemaField{
		clientIDField,
		clientSecretField,
		jwtField,
	}

	// FieldRelationships defines relationships between the fields listed in
	// ConfigurationFields that can be automatically validated.
	FieldRelationships = []field.SchemaFieldRelationship{}
)

//go:generate go run ./gen
var Config = field.NewConfiguration(
	ConfigurationFields,
	field.WithConstraints(FieldRelationships...),
	field.WithConnectorDisplayName("RingCentral"),
	field.WithHelpUrl("/docs/baton/ringcentral"),
	field.WithIconUrl("/static/app-icons/ringcentral.svg"),
)
