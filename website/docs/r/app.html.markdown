---
layout: "app"
page_title: "SumoLogic: sumologic_app"
description: |-
  Installs/Uninstalls/Upgrades an app (v2) from our app catalog
---

# sumologic_app
Provides a Sumologic_App.

## Example Usage
```hcl

resource "sumologic_app" "example_app" {
	uuid = "ceb7fac5-1127-4a04-a5b8-2e49190be3d5"
	version = "1.0.1"
	parameters = {
	    "k1": "v1",
	    "k2": "v2"
	}
}
```

## Argument reference

The following arguments are supported:

- `uuid` - UUID of the app to install/uninstall/upgrade.
- `version` - Version of the app to install. You can either specify a specific version of the app or use latest to install the latest version of the app.
- `parameters` - (Optional) Map of additional parameters for the app installation.

