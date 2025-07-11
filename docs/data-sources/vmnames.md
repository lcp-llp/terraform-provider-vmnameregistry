# Data Source: vmnameregistry_vmnames

Use this data source to retrieve all VM names for a given environment, optionally filtered by location. It also exports the status for each VM name.

## Example Usage

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

## Arguments

- `environment` (Required) – The environment to filter VM names (e.g., dev, prod, preprod, devtest).
- `location` (Optional) – The Azure region/location to further filter VM names (e.g., uksouth).
- `business_unit` (Optional) – The business unit to filter VM names.

## Attributes Exported

- `vm_names` – List of VM names returned by the API.
- `statuses` – Map of VM name to status (e.g., `{"lcpdevuks-0001": "Deployed"}`).
- `business_units` – Map of VM name to business unit (e.g., `{"lcpdevuks-0001": "engineering"}`).
- `id` – The data source ID (internal, not usually needed).
