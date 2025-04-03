package main

import (
	"github.com/conductorone/baton-sdk/pkg/field"
	"github.com/spf13/viper"
)

const (
	ringCentralClientID     = "ringcentral-client-id"
	ringCentralClientSecret = "ringcentral-client-secret"
	ringCentralJWT          = "ringcentral-jwt"
)

var (
	rcClientIDField = field.StringField(
		ringCentralClientID,
		field.WithRequired(true),
		field.WithDescription("Client ID of the Baton App for RingCentral"),
	)

	rcClientSecretField = field.StringField(
		ringCentralClientSecret,
		field.WithRequired(true),
		field.WithDescription("Client Secret of the Baton App for RingCentral"),
	)

	rcJWTField = field.StringField(
		ringCentralJWT,
		field.WithRequired(true),
		field.WithDescription("JWT of the admin user on RingCentral platform"),
	)

	// ConfigurationFields defines the external configuration required for the
	// connector to run. Note: these fields can be marked as optional or
	// required.
	ConfigurationFields = []field.SchemaField{
		rcClientIDField,
		rcClientSecretField,
		rcJWTField,
	}
)

// ValidateConfig is run after the configuration is loaded, and should return an
// error if it isn't valid. Implementing this function is optional, it only
// needs to perform extra validations that cannot be encoded with configuration
// parameters.
func ValidateConfig(_ *viper.Viper) error {
	return nil
}
