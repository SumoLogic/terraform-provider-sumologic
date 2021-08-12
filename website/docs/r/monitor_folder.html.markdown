---
layout: 'sumologic'
page_title: 'SumoLogic: sumologic_monitor_folder'
description: |-
  Provides the ability to create, read, delete, and update folders for Monitors.
---

# sumologic_monitor_folder

Provides the ability to create, read, delete, and update folders for [Monitors][1].

## Example Monitor Folder

NOTE: Monitor folders are considered a different resource from Library content folders.

```hcl
resource "sumologic_monitor_folder" "tf_monitor_folder_1" {
  name        = "Terraform Managed Monitors"
  description = "A folder for monitors managed by terraform."
}
```

## Example Nested Monitor Folders

NOTE: Monitor folders allow up to six (6) levels of sub-folders.

```hcl
resource "sumologic_monitor_folder" "tf_security_team_root_folder" {
  name        = "Security Team Monitors"
  description = "Monitors used by the Security Team."
}

resource "sumologic_monitor_folder" "tf_security_team_prod_folder" {
  name        = "Production Monitors"
  description = "Monitors for the Security Team Production Environment."
  parent_id   = sumologic_monitor_folder.tf_security_team_root_folder.id
}

resource "sumologic_monitor_folder" "tf_security_team_stag_folder" {
  name        = "Staging Monitors"
  description = "Monitors for the Security Team Staging Environment."
  parent_id   = sumologic_monitor_folder.tf_security_team_root_folder.id
```

## Argument reference

The following arguments are supported:

- `type` - (Optional) The type of object model. Valid value:
  - `MonitorsLibraryFolder`
- `name` - (Required) The name of the monitor folder. The name must be alphanumeric.
- `description` - (Required) The description of the monitor folder.
- `parent_id` - (Optional) The identifier of the Monitor Folder that contains this Monitor Folder. Defaults to the root folder.

Additional data provided in state:

- `id` - (Computed) The identifier for this monitor folder.

## Import

Monitor folders can be imported using the monitor folder identifier, such as:

```hcl
terraform import sumologic_monitor_folder.tf_monitor_folder_1 0000000000ABC123
```

[1]: https://help.sumologic.com/?cid=10020
