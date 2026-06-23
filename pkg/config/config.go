package config

import "github.com/conductorone/baton-sdk/pkg/field"

const (
	RingCentralClientID     = "ringcentral-client-id"
	RingCentralClientSecret = "ringcentral-client-secret"
	RingCentralJWT          = "ringcentral-jwt"
)

var (
	rcClientIDField = field.StringField(
		RingCentralClientID,
		field.WithRequired(true),
		field.WithDisplayName("Client ID"),
		field.WithPlaceholder("Enter your RingCentral Client ID"),
		field.WithDescription("Client ID of the Baton App for RingCentral"),
	)

	rcClientSecretField = field.StringField(
		RingCentralClientSecret,
		field.WithRequired(true),
		field.WithDisplayName("Client Secret"),
		field.WithPlaceholder("Enter your RingCentral Client Secret"),
		field.WithDescription("Client Secret of the Baton App for RingCentral"),
		field.WithIsSecret(true),
	)

	rcJWTField = field.StringField(
		RingCentralJWT,
		field.WithRequired(true),
		field.WithDisplayName("JWT"),
		field.WithPlaceholder("Enter your RingCentral admin JWT"),
		field.WithDescription("JWT of the admin user on RingCentral platform"),
		field.WithIsSecret(true),
	)
)

//go:generate go run ./gen
var Config = field.NewConfiguration(
	[]field.SchemaField{rcClientIDField, rcClientSecretField, rcJWTField},
	field.WithConnectorDisplayName("RingCentral"),
	field.WithIconUrl("/static/app-icons/ringcentral.svg"),
	field.WithHelpUrl("/docs/baton/ringcentral"),
)
