package config

//go:generate go run ./gen

import (
	"github.com/conductorone/baton-sdk/pkg/field"
)

var (
	AccessTokenField = field.StringField(
		"access-token",
		field.WithRequired(true),
		field.WithIsSecret(true),
		field.WithDescription("The access token for the Pulumi Cloud organization"),
	)
	OrgNameField = field.StringField(
		"org-name",
		field.WithRequired(true),
		field.WithDescription("The name of the Pulumi Cloud organization"),
	)
	BaseURLField = field.StringField(
		"base-url",
		field.WithDescription("Override the Pulumi Cloud API URL (for testing)"),
		field.WithHidden(true),
		field.WithExportTarget(field.ExportTargetCLIOnly),
	)

	Config = field.NewConfiguration([]field.SchemaField{
		AccessTokenField,
		OrgNameField,
		BaseURLField,
	})
)

func ValidateConfig(c *PulumiCloud) error {
	return nil
}
