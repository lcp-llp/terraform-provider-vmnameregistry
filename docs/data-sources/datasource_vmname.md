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
```

## Arguments

- `vm_name` (Required) – The VM name to look up.

## Attributes Exported

- `id` – The VM name (same as `vm_name`).
