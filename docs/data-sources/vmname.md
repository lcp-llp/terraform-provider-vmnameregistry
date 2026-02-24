# Data Source: vmnameregistry_vmname

Use this data source to look up a single VM name by its value.

## Example Usage

```hcl
data "vmnameregistry_vmname" "lookup" {
  environment = "dev"
  vm_name     = "devuks-0001"
}

output "vm_name" {
  value = data.vmnameregistry_vmname.lookup.vm_name
}

output "vm_status" {
  value = data.vmnameregistry_vmname.lookup.status
}

output "vm_business_unit" {
  value = data.vmnameregistry_vmname.lookup.business_unit
}

output "vm_location" {
  value = data.vmnameregistry_vmname.lookup.location
}

output "vm_cloud_platform" {
  value = data.vmnameregistry_vmname.lookup.cloud_platform
}
```

## Arguments

- `environment` (Required) – The environment for the VM name (e.g., dev, prod, preprod, devtest).
- `vm_name` (Required) – The VM name to look up.

## Attributes Exported

- `vm_name` – The VM name (same as input).
- `status` – The status of the VM name (e.g., `Deployed`, `Reserved`, `Available`, `Decommissioned`).
- `business_unit` – The business unit associated with the VM name.
- `cloud_platform` – The cloud platform of the VM name (e.g., `aws`, `azure`, `gcp`, `onprem`).
- `location` – The location/region of the VM name.
- `id` – The VM name (same as `vm_name`).
