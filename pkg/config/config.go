package config

//go:generate go run ./gen

import (
	"github.com/conductorone/baton-sdk/pkg/field"
)

var Config = field.NewConfiguration([]field.SchemaField{
	field.StringField(
		"access-token",
		field.WithRequired(true),
		field.WithDescription("The access token for the Pulumi Cloud organization"),
	),
	field.StringField(
		"org-name",
		field.WithRequired(true),
		field.WithDescription("The name of the Pulumi Cloud organization"),
	),
})

func ValidateConfig(c *PulumiCloud) error {
	return nil
}
