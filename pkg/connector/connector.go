package connector

import (
	"context"
	"fmt"
	"io"

	"github.com/conductorone/baton-pulumi-cloud/pkg/client"
	v2 "github.com/conductorone/baton-sdk/pb/c1/connector/v2"
	"github.com/conductorone/baton-sdk/pkg/annotations"
	"github.com/conductorone/baton-sdk/pkg/connectorbuilder"
)

// Connector implements the Pulumi connector
type Connector struct {
	client  *client.Client
	orgName string
}

// ResourceSyncers returns a ResourceSyncer for each resource type that should be synced
func (c *Connector) ResourceSyncers(ctx context.Context) []connectorbuilder.ResourceSyncer {
	return []connectorbuilder.ResourceSyncer{
		newOrgBuilder(c.client, c.orgName),
		newUserBuilder(c.client, c.orgName),
		newTeamBuilder(c.client, c.orgName),
	}
}

// Asset returns asset data for the connector
func (c *Connector) Asset(ctx context.Context, asset *v2.AssetRef) (string, io.ReadCloser, error) {
	return "", nil, nil
}

// Metadata returns metadata about the connector
func (c *Connector) Metadata(ctx context.Context) (*v2.ConnectorMetadata, error) {
	return &v2.ConnectorMetadata{
		DisplayName: "Pulumi Cloud",
	}, nil
}

// Validate ensures the connector is properly configured
func (c *Connector) Validate(ctx context.Context) (annotations.Annotations, error) {
	// Test the connection by trying to list users
	_, err := c.client.ListUsers(ctx, c.orgName, "")
	if err != nil {
		return nil, err
	}
	return nil, nil
}

// New returns a new instance of the connector
func New(ctx context.Context, client *client.Client, orgName string) (*Connector, error) {
	if client == nil {
		return nil, fmt.Errorf("pulumi client not provided")
	}

	if orgName == "" {
		return nil, fmt.Errorf("organization name not provided")
	}

	return &Connector{
		client:  client,
		orgName: orgName,
	}, nil
}
