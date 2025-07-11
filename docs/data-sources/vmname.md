# Data Source: vmnameregistry_vmname

Use this data source to look up a single VM name by its value.

## Example Usage

```hcl
data "vmnameregistry_vmname" "lookup" {
  vm_name = "lcpdevuks-0001"
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
```

## Arguments

- `vm_name` (Required) – The VM name to look up.

## Attributes Exported

- `vm_name` – The VM name (same as input).
- `status` – The status of the VM name (e.g., "Deployed", "Reserved", "Available").
- `business_unit` – The business unit associated with the VM name.
- `location` – The location/region extracted from the VM name pattern.
- `id` – The VM name (same as `vm_name`).
