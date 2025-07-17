---
layout: "sumologic"
page_title: "SumoLogic: sumologic_apps"
description: |-
  Provides an easy way to retrieve all Sumo Logic v2 apps.
---

# sumologic_apps
Provides an easy way to retrieve all Sumo Logic v2 apps.


## Example Usage
```hcl
data "sumologic_apps" "test" {}
```

```hcl
data "sumologic_apps" "test" {
    name = "MySQL - OpenTelemetry"
	author = "Sumo Logic"
}
```


## Attributes reference

The following attributes are exported:

- `uuid` - UUID of the app.
- `name` - Name of the app.
- `description` - Description of the app.
- `latest_version` - Latest version of the app.
- `icon` - URL of the icon for the app.
- `author` - Author of the app.
- `account_types` - URL of the icon for the app
- `log_analytics_filter` - The search filter which would be applied on partitions which belong to Log Analytics product area.
- `beta` - URL of the icon for the app.
- `installs` - Number of times the app was installed.
- `appType` - Type of an app.
- `attributes` - A map of attributes for this app. Attributes allow to group apps based on different criteria.
### Values in attributes type are : 
  - `category` 
  - `use_case`
  - `collection`
- `family` - Provides a mechanism to link different apps.
- `installable` - Whether the app is installable or not as not all apps are installable.
- `show_on_marketplace` - Whether the app should show up on sumologic.com/applications webpage.


