# Terraform Provider: vmnameregistry

A custom Terraform provider for managing and querying VM names in the LCP VM Name Registry.

## Features
- Create, update, and delete VM names for a given environment, location, and business unit
- Look up a single VM name
- List all VM names for an environment, optionally filtered by location and/or business unit
- Export VM name statuses and business units

## Requirements
- Terraform >= 1.0
- Go >= 1.18 (for building the provider)

## Provider Configuration
```hcl
provider "vmnameregistry" {
  url = "https://your-api-url/api/ManageVmName" # Base URL to the registry API (no query params)
}
```

## Resources
### `vmnameregistry_vmname`
Manages a VM name in the registry.

```hcl
resource "vmnameregistry_vmname" "example" {
  environment   = "dev"
  location      = "uksouth"
  business_unit = "engineering"
  status        = "Deployed" # Optional, defaults to "Deployed"
}

output "vm_name" {
  value = vmnameregistry_vmname.example.vm_name
}
```

## Data Sources
### `vmnameregistry_vmname`
Look up a single VM name by value.

```hcl
data "vmnameregistry_vmname" "lookup" {
  vm_name = "lcpdevuks-0001"
}

output "vm_name" {
  value = data.vmnameregistry_vmname.lookup.vm_name
}
```

### `vmnameregistry_vmnames`
List all VM names for an environment, optionally filtered by location and/or business unit, and export their statuses and business units.

```hcl
data "vmnameregistry_vmnames" "all" {
  environment   = "prod"
  location      = "uksouth"     # optional
  business_unit = "engineering" # optional
}

output "all_vm_names" {
  value = data.vmnameregistry_vmnames.all.vm_names
}

output "all_vm_statuses" {
  value = data.vmnameregistry_vmnames.all.statuses
}

output "all_vm_business_units" {
  value = data.vmnameregistry_vmnames.all.business_units
}
```

## Development
- Clone the repo
- Run `go mod tidy`
- Run `go build`
- Run `go test`

## License
MIT