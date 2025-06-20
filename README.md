# Terraform Provider: vmnameregistry

A custom Terraform provider for managing and querying VM names in the LCP VM Name Registry.

## Features
- Create, update, and delete VM names for a given environment and location
- Look up a single VM name
- List all VM names for an environment, optionally filtered by location
- Export VM name statuses

## Requirements
- Terraform >= 1.0
- Go >= 1.18 (for building the provider)

## Development
- Clone the repo
- Run `go mod tidy`
- Run `go build`
- Run `go test`

