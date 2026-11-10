---
layout: "sumologic"
page_title: "SumoLogic: sumologic_monitor_folder"
description: |-
  Provides an easy way to retrieve a Monitor folder by its path.
---

# sumologic_monitor_folder
Provides an easy way to retrieve a [Monitor][1] folder by its path.

You must specify the absolute path of the folder to retrieve. The Monitors library has a root folder named
`Monitor`, so every path starts with `/Monitor`. For example, a folder named `Security` inside the root
folder has the path `/Monitor/Security`, and a sub-folder named `Production` inside it has the path
`/Monitor/Security/Production`.

-> Monitor folders are a different object type from Content Library folders. To look up a Content Library folder,
use the [sumologic_folder](folder.html) data source instead.

## Example Usage
```hcl
# Look up the monitor folder "Production" inside the "Security" folder
data "sumologic_monitor_folder" "prod" {
  path = "/Monitor/Security/Production"
}

# Use the looked-up folder as the parent of a new monitor folder
resource "sumologic_monitor_folder" "prod_alerts" {
  name        = "Alerts"
  description = "Alert monitors for production."
  parent_id   = data.sumologic_monitor_folder.prod.id
}
```

## Argument reference

The following arguments are supported:

- `path` - (Required) The absolute path of the monitor folder, starting with `/Monitor`. For example,
`/Monitor/Security/Production`. The data source returns an error if no folder exists at this path.

## Attributes reference

The following attributes are exported:

- `id` - The ID of the monitor folder.
- `name` - The name of the monitor folder.
- `description` - The description of the monitor folder.

[1]: https://help.sumologic.com/docs/alerts/monitors/
