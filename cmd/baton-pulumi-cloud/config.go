package main

import (
	"github.com/conductorone/baton-sdk/pkg/field"
	"github.com/spf13/viper"
)

var (
	accessTokenField = field.StringField(
		"access-token",
		field.WithRequired(true),
		field.WithDescription("The access token for the Pulumi Cloud organization"),
	)
	orgNameField = field.StringField(
		"org-name",
		field.WithRequired(true),
		field.WithDescription("The name of the Pulumi Cloud organization"),
	)
	ConfigurationFields = []field.SchemaField{
		accessTokenField,
		orgNameField,
	}

	// FieldRelationships defines relationships between the fields listed in
	// ConfigurationFields that can be automatically validated. For example, a
	// username and password can be required together, or an access token can be
	// marked as mutually exclusive from the username password pair.
	FieldRelationships = []field.SchemaFieldRelationship{}
)

// ValidateConfig is run after the configuration is loaded, and should return an
// error if it isn't valid. Implementing this function is optional, it only
// needs to perform extra validations that cannot be encoded with configuration
// parameters.
func ValidateConfig(v *viper.Viper) error {
	return nil
}
