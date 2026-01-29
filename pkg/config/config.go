package config

import (
	"github.com/conductorone/baton-sdk/pkg/field"
)

var (
	clientIDField = field.StringField(
		"ringcentral-client-id",
		field.WithDisplayName("RingCentral Client ID"),
		field.WithRequired(true),
		field.WithDescription("Client ID of the Baton App for RingCentral"),
		field.WithIsSecret(true),
	)

	clientSecretField = field.StringField(
		"ringcentral-client-secret",
		field.WithDisplayName("RingCentral Client Secret"),
		field.WithRequired(true),
		field.WithDescription("Client Secret of the Baton App for RingCentral"),
		field.WithIsSecret(true),
	)

	jwtField = field.StringField(
		"ringcentral-jwt",
		field.WithDisplayName("RingCentral JWT"),
		field.WithRequired(true),
		field.WithDescription("JWT of the admin user on RingCentral platform"),
		field.WithIsSecret(true),
	)

	configurationFields = []field.SchemaField{
		clientIDField,
		clientSecretField,
		jwtField,
	}
)

//go:generate go run ./gen
var Config = field.NewConfiguration(
	configurationFields,
	field.WithConnectorDisplayName("RingCentral"),
)
