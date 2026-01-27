package config

//go:generate go run ./gen

import (
	"github.com/conductorone/baton-sdk/pkg/field"
)

var (
	AccessTokenField = field.StringField(
		"access-token",
		field.WithRequired(true),
		field.WithDescription("The access token for the Pulumi Cloud organization"),
	)
	OrgNameField = field.StringField(
		"org-name",
		field.WithRequired(true),
		field.WithDescription("The name of the Pulumi Cloud organization"),
	)

	// ConfigurationFields defines the external configuration required for the
	// connector to run.
	ConfigurationFields = []field.SchemaField{AccessTokenField, OrgNameField}

	// FieldRelationships defines relationships between the fields.
	FieldRelationships = []field.SchemaFieldRelationship{}

	// Config is the configuration schema for the connector.
	Config = field.Configuration{
		Fields:      ConfigurationFields,
		Constraints: FieldRelationships,
	}
)
