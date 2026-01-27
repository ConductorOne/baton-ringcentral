package config

//go:generate go run ./gen

import (
	"github.com/conductorone/baton-sdk/pkg/field"
)

var (
	RCClientIDField = field.StringField(
		"ringcentral-client-id",
		field.WithRequired(true),
		field.WithDescription("Client ID of the Baton App for RingCentral"),
	)

	RCClientSecretField = field.StringField(
		"ringcentral-client-secret",
		field.WithRequired(true),
		field.WithDescription("Client Secret of the Baton App for RingCentral"),
	)

	RCJWTField = field.StringField(
		"ringcentral-jwt",
		field.WithRequired(true),
		field.WithDescription("JWT of the admin user on RingCentral platform"),
	)

	// ConfigurationFields defines the external configuration required for the
	// connector to run.
	ConfigurationFields = []field.SchemaField{
		RCClientIDField,
		RCClientSecretField,
		RCJWTField,
	}

	// FieldRelationships defines relationships between the fields.
	FieldRelationships = []field.SchemaFieldRelationship{}

	// Config is the configuration schema for the connector.
	Config = field.Configuration{
		Fields:      ConfigurationFields,
		Constraints: FieldRelationships,
	}
)
