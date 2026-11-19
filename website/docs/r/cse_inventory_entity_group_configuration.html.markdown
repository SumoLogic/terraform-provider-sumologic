---
layout: "sumologic"
page_title: "SumoLogic: sumologic_cse_inventory_entity_group_configuration"
description: |-
  Provides a CSE Inventory Entity Group Configuration
---

# sumologic_cse_inventory_entity_group_configuration
Provides a Cloud SIEM Inventory Entity Group Configuration. It groups entities from an inventory source (for example, Active Directory) and can apply a criticality and tags to them. See [Create an Entity Group](https://www.sumologic.com/help/docs/cse/records-signals-entities-insights/create-an-entity-group/). To group entities by name or IP address range instead, use `sumologic_cse_entity_entity_group_configuration`.

The underlying API is documented in the [Cloud SIEM API reference](https://api.sumologic.com/docs/sec/) (Entity Group Configurations). In the API, `inventory_type` is an enum of `computer` and `user`; the other inventory fields, including `criticality`, `inventory_source` and `tag_schema`, are strings. The API lists available criticalities (`/entity-criticality-configs`) and tag schemas (`/tag-schemas`, each identified by a `key`).

## Example Usage

### Basic
```hcl
resource "sumologic_cse_inventory_entity_group_configuration" "inventory_entity_group_configuration" {
  name             = "Inventory entity group configuration"
  description      = "Inventory entity group description"
  inventory_type   = "user"
  inventory_source = "Active Directory"
  inventory_key    = "groups"
  inventory_value  = "http_servers"
  criticality      = "HIGH"
  suppressed       = false
  tags             = ["tag"]
}
```

## Argument reference

The following arguments are supported:

- `name` - (Required) The entity group configuration name. Must not be empty.
- `inventory_type` - (Required) The type of inventory entity. Allowed values: `user`, `computer`. Any other value is rejected at plan time.
- `inventory_source` - (Required) The inventory source. The inventory source that supplies the data, as shown in the "Source" list in the UI (for example `Active Directory`). Case sensitive, spelled exactly as listed. Built-in values: `Sentinel One`, `AWS`, `Azure AD`, `Carbon Black`, `Cylance`, `Okta`, `Qualys`, `Active Directory`, `Google Workspace`, `Sailpoint`, `Tenable`, `Armis`, `CrowdStrike`, `Rapid7`, `Cisco Vulnerability Management`. Sources that have ingested inventory data, including custom inventory sources, are also accepted, so this list is not exhaustive. See [Inventory Sources and Data](https://www.sumologic.com/help/docs/cse/administration/inventory-sources-and-data/). The provider does not validate this value; Sumo Logic rejects unknown sources.
- `inventory_key` - (Optional) The inventory attribute used to match entities (the "Inventory Key" field in the UI). Use `groups` to match membership of a group in the inventory system; second-level unnormalized attributes such as `fields.distinguishedname` are also accepted in the UI. Examples: `groups`, `normalizedHostname`, `normalizedComputerName`. Must not be empty if set.
- `inventory_value` - (Optional) The value of `inventory_key` to match. The UI accepts regular expressions here (for example `.*OU\=ADFS.*`); if `inventory_key` is `groups`, this is the name of the group in the inventory system. Must not be empty if set.
- `description` - (Optional) The entity group configuration description.
- `criticality` - (Optional) The criticality applied to matching entities, selected in the UI from the criticalities defined in Cloud SIEM (see `sumologic_cse_entity_criticality_config` and [Entity Criticality](https://www.sumologic.com/help/docs/cse/records-signals-entities-insights/entity-criticality/)). Must not be empty if set. Use the criticality's name, not its ID. Setting `null`, `none` or `default` clears the criticality. The provider does not validate the value.
- `suppressed` - (Optional) Whether signals on entities in the group are suppressed. Defaults to `false`.
- `tags` - (Optional) A list of tags applied to matching entities. Each tag must not be empty.
- `dynamic_tags` - (Optional) Whether dynamic tags are enabled for the configuration. Defaults to `false`. Use together with `tag_schema`.
- `tag_schema` - (Optional) The `key` of a custom tag schema (not its label) used to generate dynamic tags ("Tag Schema" in the UI). The schema must be a freeform tag schema. Ignored if `group` is set. When dynamic tags are on, the UI sets the value to `*`, and `tags` are added in addition to the generated tags. Must not be empty if set.
- `group` - (Optional, Deprecated) The inventory group. Deprecated: use `inventory_key` and `inventory_value` instead; `group` will be removed in a future release of the provider. When set, Sumo Logic treats it as `inventory_key = "groups"` with `inventory_value` set to this value, ignores `tag_schema` and sets `dynamic_tags` to `false`. Must not be empty if set.

The following attributes are exported:

- `id` - The internal ID of the entity group configuration.

## Import

An inventory entity group configuration can be imported using its ID, e.g.:
```hcl
terraform import sumologic_cse_inventory_entity_group_configuration.inventory_entity_group_configuration <configuration_id>
```
