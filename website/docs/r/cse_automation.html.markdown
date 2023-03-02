---
layout: "sumologic"
page_title: "SumoLogic: sumologic_cse_automation"
description: |-
  Provides a Sumologic CSE Automation
---

# sumologic_cse_automation
Provides a Sumologic CSE Automation.

## Example Usage
```hcl
resource "sumologic_cse_automation" "insight_automation" {
  playbook_id = "638079aedb99cafada1e80a0"
  cse_resource_type = "INSIGHT"
  execution_types = ["NEW_INSIGHT","INSIGHT_CLOSED"]
}

resource "sumologic_cse_automation" "entity_automation" {
  playbook_id = "638079aedb99cafada1e80a0"
  cse_resource_type = "ENTITY"
  cse_resource_sub_types = ["_ip"]
  execution_types = ["ON_DEMAND"]
}
```

## Argument reference

The following arguments are supported:

- `plabook_id` - (Required) CSOAR playbook Id.
- `name` - (Computed) Automation name.
- `description` - (Computed) Automation description.
- `cse_resource_type` - (Required) CSE Resource type for automation. Valid values: "INSIGHT", "ENTITY".
- `execution_types` - (Required) Automation execution type. Valid values: "NEW_INSIGHT", "INSIGHT_CLOSED", "ON_DEMAND".
- `cse_resource_sub_types` - (Optional) CSE Resource sub-type when cse_resource_type is specified as "ENTITY". Examples: "_ip", "_mac".

The following attributes are exported:

- `id` - The internal ID of the Automation.

## Import

Automation can be imported using the field id, e.g.:
```hcl
terraform import sumologic_cse_automation.automation id
```
