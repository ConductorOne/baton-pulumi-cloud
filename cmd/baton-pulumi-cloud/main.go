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
	"github.com/conductorone/baton-sdk/pkg/types"
	"github.com/grpc-ecosystem/go-grpc-middleware/logging/zap/ctxzap"
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

func getConnector(ctx context.Context, pc *cfg.PulumiCloud) (types.ConnectorServer, error) {
	l := ctxzap.Extract(ctx)

	if err := cfg.ValidateConfig(pc); err != nil {
		return nil, err
	}

	c, err := client.NewClient(pc.AccessToken)
	if err != nil {
		return nil, fmt.Errorf("failed to create client: %w", err)
	}

	cb, err := connector.New(ctx, c, pc.OrgName)
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
