---
layout: "sumologic"
page_title: "SumoLogic: sumologic_role_v2"
description: |-
  Provides a way to retrieve Sumo Logic role details (id, names, etc) for a role managed outside of terraform.
---

# sumologic_role_v2 ( beta )

Provides a way to retrieve Sumo Logic role details (id, names, etc) for a role
managed by another terraform stack.


## Example Usage
```hcl
data "sumologic_role_v2" "this" {
  name = "MyRole"
}
```

```hcl
data "sumologic_role_v2" "that" {
  id = "1234567890"
}
```

A role can be looked up by either `id` or `name`. One of those attributes needs to be specified.

If both `id` and `name` have been specified, `id` takes precedence.

## Attributes reference

The following attributes are exported:

- `id` - The internal ID of the role. This can be used to create users having that role.
- `name` - The name of the role.
- `description` - The description of the role.
- `capabilities` - The list of capabilities associated with the role.
- `selected_views` - List of views with specific view level filters in accordance to the selectionType chosen.
- `audit_data_filter` - The search filter which would be applied on partitions which belong to Audit Data product area. Help Doc : (https://help.sumologic.com/docs/manage/security/audit-index/). Applicable with only `All` selectionType
- `security_data_filter` - The search filter which would be applied on partitions which belong to Security Data product area. Applicable with only `All` selectionType.
- `log_analytics_filter` - The search filter which would be applied on partitions which belong to Log Analytics product area. Applicable with only `All` selectionType
- `selection_type` - Describes the Permission Construct for the list of views in "selectedViews" parameter.
### Values in selection type are : 
  - `All` selectionType would allow access to all views in the org.
  - `Allow` selectionType would allow access to specific views mentioned in "selectedViews" parameter.
  - `Deny` selectionType would deny access to specific views mentioned in "selectedViews" parameter.

