---
layout: "sumologic"
page_title: "SumoLogic: sumologic_cse_network_block"
description: |-
  Provides a CSE Network Block
---

# sumologic_cse_network_block
Provides a Sumo Logic CSE Network Block.

## Example Usage
```hcl
resource "sumologic_cse_network_block" "network_block" {
  address_block         = "10.0.1.0/26"
  label     = "network block from terraform"
  internal = "true"
  suppresses_signals = "false"
}
```

## Argument reference

The following arguments are supported:

- `address_block` - (Required) The address block.
- `label` - (Required) The displayable label of the address block.
- `internal` - (Required) Internal flag.
- `suppresses_signals` - (Required) Suppresses signal flag.

The following attributes are exported:

- `id` - The internal ID of the network block.

## Import

Network Block can be imported using the field id, e.g.:
```hcl
terraform import sumologic_cse_network_block.network_block id
```