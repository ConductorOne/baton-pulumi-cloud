![Baton Logo](./baton-logo.png)

# `baton-pulumi-cloud` [![Go Reference](https://pkg.go.dev/badge/github.com/conductorone/baton-pulumi-cloud.svg)](https://pkg.go.dev/github.com/conductorone/baton-pulumi-cloud) ![main ci](https://github.com/conductorone/baton-pulumi-cloud/actions/workflows/main.yaml/badge.svg)

`baton-pulumi-cloud` is a connector for Pulumi Cloud built using the [Baton SDK](https://github.com/conductorone/baton-sdk).

Check out [Baton](https://github.com/conductorone/baton) to learn more the project in general.

# Getting Started

## brew

```
brew install conductorone/baton/baton conductorone/baton/baton-pulumi-cloud
baton-pulumi-cloud
baton resources
```

## docker

```
docker run --rm -v $(pwd):/out -e BATON_DOMAIN_URL=domain_url -e BATON_API_KEY=apiKey -e BATON_USERNAME=username ghcr.io/conductorone/baton-pulumi-cloud:latest -f "/out/sync.c1z"
docker run --rm -v $(pwd):/out ghcr.io/conductorone/baton:latest -f "/out/sync.c1z" resources
```

## source

```
go install github.com/conductorone/baton/cmd/baton@main
go install github.com/conductorone/baton-pulumi-cloud/cmd/baton-pulumi-cloud@main

baton-pulumi-cloud

baton resources
```

# Data Model

`baton-pulumi-cloud` will pull down information about the following resources:

- Organizations
- Teams
- Users

# Contributing, Support and Issues

We started Baton because we were tired of taking screenshots and manually
building spreadsheets. We welcome contributions, and ideas, no matter how
small&mdash;our goal is to make identity and permissions sprawl less painful for
everyone. If you have questions, problems, or ideas: Please open a GitHub Issue!

See [CONTRIBUTING.md](https://github.com/ConductorOne/baton/blob/main/CONTRIBUTING.md) for more details.

# `baton-pulumi-cloud` Command Line Usage

```
baton-pulumi-cloud

Usage:
  baton-pulumi-cloud [flags]
  baton-pulumi-cloud [command]

Available Commands:
  capabilities       Get connector capabilities
  completion         Generate the autocompletion script for the specified shell
  help               Help about any command

Flags:
      --access-token string              required: The access token for the Pulumi Cloud organization ($BATON_ACCESS_TOKEN)
      --client-id string                 The client ID used to authenticate with ConductorOne ($BATON_CLIENT_ID)
      --client-secret string             The client secret used to authenticate with ConductorOne ($BATON_CLIENT_SECRET)
  -f, --file string                      The path to the c1z file to sync with ($BATON_FILE) (default "sync.c1z")
  -h, --help                             help for baton-pulumi-cloud
      --log-format string                The output format for logs: json, console ($BATON_LOG_FORMAT) (default "json")
      --log-level string                 The log level: debug, info, warn, error ($BATON_LOG_LEVEL) (default "info")
      --org-name string                  required: The name of the Pulumi Cloud organization ($BATON_ORG_NAME)
      --otel-collector-endpoint string   The endpoint of the OpenTelemetry collector to send observability data to ($BATON_OTEL_COLLECTOR_ENDPOINT)
  -p, --provisioning                     This must be set in order for provisioning actions to be enabled ($BATON_PROVISIONING)
      --skip-full-sync                   This must be set to skip a full sync ($BATON_SKIP_FULL_SYNC)
      --ticketing                        This must be set to enable ticketing support ($BATON_TICKETING)
  -v, --version                          version for baton-pulumi-cloud

Use "baton-pulumi-cloud [command] --help" for more information about a command.
```
