package main

import (
	cfg "github.com/conductorone/baton-pulumi-cloud/pkg/config"
	"github.com/conductorone/baton-sdk/pkg/config"
)

func main() {
	config.Generate("pulumi-cloud", cfg.Config)
}
