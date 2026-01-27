package config

//go:generate go run ./gen

import (
	"github.com/conductorone/baton-sdk/pkg/field"
)

var Config = field.NewConfiguration([]field.SchemaField{
	field.StringField(
		"ringcentral-client-id",
		field.WithRequired(true),
		field.WithDescription("Client ID of the Baton App for RingCentral"),
	),
	field.StringField(
		"ringcentral-client-secret",
		field.WithRequired(true),
		field.WithDescription("Client Secret of the Baton App for RingCentral"),
	),
	field.StringField(
		"ringcentral-jwt",
		field.WithRequired(true),
		field.WithDescription("JWT of the admin user on RingCentral platform"),
	),
})

func ValidateConfig(c *Ringcentral) error {
	return nil
}
