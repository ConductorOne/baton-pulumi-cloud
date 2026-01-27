package main

import (
	"context"
	"fmt"
	"os"

	"github.com/conductorone/baton-pulumi-cloud/pkg/client"
	cfg "github.com/conductorone/baton-pulumi-cloud/pkg/config"
	"github.com/conductorone/baton-pulumi-cloud/pkg/connector"
	"github.com/conductorone/baton-sdk/pkg/config"
	"github.com/conductorone/baton-sdk/pkg/connectorbuilder"
	"github.com/conductorone/baton-sdk/pkg/connectorrunner"
	"github.com/conductorone/baton-sdk/pkg/types"
	"github.com/grpc-ecosystem/go-grpc-middleware/logging/zap/ctxzap"
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

var version = "dev"

func main() {
	ctx := context.Background()

	_, cmd, err := config.DefineConfiguration(
		ctx,
		"baton-pulumi-cloud",
		getConnector,
		cfg.Config,
		connectorrunner.WithDefaultCapabilitiesConnectorBuilder(&connector.Connector{}),
	)
	if err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(1)
	}

	cmd.Version = version

	err = cmd.Execute()
	if err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(1)
	}
}

func getConnector(ctx context.Context, v *viper.Viper) (types.ConnectorServer, error) {
	l := ctxzap.Extract(ctx)

	token := v.GetString(cfg.AccessTokenField.FieldName)
	orgName := v.GetString(cfg.OrgNameField.FieldName)

	c, err := client.NewClient(token)
	if err != nil {
		return nil, fmt.Errorf("failed to create client: %w", err)
	}

	cb, err := connector.New(ctx, c, orgName)
	if err != nil {
		l.Error("error creating connector", zap.Error(err))
		return nil, err
	}
	connector, err := connectorbuilder.NewConnector(ctx, cb)
	if err != nil {
		return nil, err
	}
	return connector, nil
}
