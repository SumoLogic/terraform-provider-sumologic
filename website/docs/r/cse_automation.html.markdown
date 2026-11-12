---
layout: "sumologic"
page_title: "SumoLogic: sumologic_cse_automation"
description: |-
  Provides a Sumologic CSE Automation
---

# sumologic_cse_automation
Provides a Sumologic CSE Automation. An automation runs a Cloud SOAR playbook on Cloud SIEM insights or entities, either automatically when an insight is created or closed, or manually. See [Automations in Cloud SIEM](https://www.sumologic.com/help/docs/cse/automation/automations-in-cloud-siem/).

The underlying API is documented in the [Cloud SIEM API reference](https://api.sumologic.com/docs/sec/) (Automations). The API types `cse_resource_type`, `execution_types` and `cse_resource_sub_types` as strings; the provider does not validate them, but Sumo Logic rejects values outside the sets listed below.

## Example Usage

### Run a playbook when an insight is created or closed
```hcl
resource "sumologic_cse_automation" "insight_automation" {
  playbook_id       = "638079aedb99cafada1e80a0"
  cse_resource_type = "INSIGHT"
  execution_types   = ["NEW_INSIGHT", "INSIGHT_CLOSED"]
  enabled           = true
}
```

### Run a playbook on demand for IP entities
```hcl
resource "sumologic_cse_automation" "entity_automation" {
  playbook_id            = "638079aedb99cafada1e80a0"
  cse_resource_type      = "ENTITY"
  cse_resource_sub_types = ["_ip"]
  execution_types        = ["ON_DEMAND"]
  enabled                = true
}
```

## Argument reference

The following arguments are supported:

- `playbook_id` - (Required) The ID of the Cloud SOAR playbook that the automation runs.
- `cse_resource_type` - (Required) The kind of object the playbook receives data for (the "Object" field in the UI). Allowed values: `INSIGHT`, `ENTITY`. Sumo Logic converts the value to upper case before validating it and rejects any other value.
- `execution_types` - (Required) When the automation runs, as a list of non-empty strings. These correspond to the "Insight Created", "Insight Closed" and "Manually Done" choices in the UI: `NEW_INSIGHT`, `INSIGHT_CLOSED`, `ON_DEMAND` (case sensitive). Sumo Logic rejects any other value. An automation without `ON_DEMAND` does not appear in the manual-run menus on insights and entities.
- `enabled` - (Required) Whether the automation is enabled (the "Status" field in the UI). There is no default; you must set it to `true` or `false`. Disabled automations do not run automatically and do not appear in the manual-run menus.
- `cse_resource_sub_types` - (Optional) A list of resource sub-types the automation applies to, the playbook runs only on entities of these types. Required, with at least one value, when `cse_resource_type` is `ENTITY`; ignored when it is `INSIGHT`. Each value must be a valid entity type: a built-in type such as `_ip` or `_mac`, or the identifier of a custom entity type.

The following attributes are exported:

- `id` - The internal ID of the automation.
- `name` - The automation name. Read-only: taken from the Cloud SOAR playbook and cannot be configured.
- `description` - The automation description. Read-only: taken from the Cloud SOAR playbook and cannot be configured.

## Import

An automation can be imported using its ID, e.g.:
```hcl
terraform import sumologic_cse_automation.insight_automation <automation_id>
```
