---
layout: 'sumologic'
page_title: 'SumoLogic: sumologic_monitor_folder'
description: |-
  Provides the ability to create, read, delete, and update folders for Monitors.
---

# sumologic_monitor

Provides the ability to create, read, delete, and update folders for [Monitors][1].

## Example Monitor Folder

NOTE: Monitor folders are considered a different resource from Library content folders.

```hcl
resource "sumologic_monitor_folder" "tf_monitor_folder_1" {
  name        = "test terraform folder"
  description = "a folder for monitors"
}
```

## Argument reference

The following arguments are supported:

- `type` - (Optional) The type of object model. Valid value:
  - `MonitorsLibraryFolder`
- `name` - (Required) The name of the monitor folder. The name must be alphanumeric.
- `description` - (Required) The description of the monitor folder.
- `parent_id` - (Optional) The ID of the Monitor Folder that contains this Monitor Folder. Defaults to the root folder.

## Import

Monitor folders can be imported using the monitor ID, such as:

```hcl
terraform import sumologic_monitor_folder.test 1234567890
```

[1]: https://help.sumologic.com/?cid=10020
