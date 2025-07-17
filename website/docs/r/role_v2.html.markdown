---
layout: "sumologic"
page_title: "SumoLogic: sumologic_role_v2"
description: |-
  Provides a Sumologic Role v2
---

# sumologic_role_v2 ( beta )
Provider to manage Sumologic Role v2

## Example Usage
```hcl
resource "sumologic_role_v2" "example_role_v2" {
    selected_views {
      view_name = "view1"
    }
    selected_views {
      view_name = "view2"
    }
    name = "DataAdmin"
    audit_data_filter = "info"
    selection_type = "Allow"
    capabilities = ["manageContent","manageDataVolumeFeed","manageFieldExtractionRules","manageS3DataForwarding"]
    description = "Manage data of the org."
    security_data_filter = "error"
    log_analytics_filter = "!_sourceCategory=collector"
}
```
## Argument reference

The following arguments are supported:

- `selected_views` - (Optional) List of views with specific view level filters in accordance to the selectionType chosen.
- `name` - (Required) Name of the role.
- `capabilities` - (Optional) List of [capabilities](https://help.sumologic.com/docs/manage/users-roles/roles/role-capabilities/) associated with this role. 
- `description` - (Optional) Description of the role.
- `security_data_filter` - (Optional) A search filter which would be applied on partitions which belong to Security Data product area.
- `log_analytics_filter` - (Optional) A search filter which would be applied on partitions which belong to Log Analytics product area.
- `audit_data_filter` - (Optional) A search filter which would be applied on partitions which belong to Audit Data product area. Help Doc : (https://help.sumologic.com/docs/manage/security/audit-index/).
- `selection_type` - (Optional) Describes the Permission Construct for the list of views in "selectedViews" parameter. 
### Valid Values are : 
  - `All` selectionType would allow access to all views in the org.
  - `Allow` selectionType would allow access to specific views mentioned in "selectedViews" parameter.
  - `Deny` selectionType would deny access to specific views mentioned in "selectedViews" parameter.

The following attributes are exported:

- `id` - The internal ID of the role_v2



[Back to Index][0]

[0]: ../README.md
