---
layout: "sumologic"
page_title: "SumoLogic: sumologic_admin_recommended_folder"
description: |-
  Provides an easy way to retrieve the Admin Recommended Folder.
---

# sumologic_admin_recommended_folder
Provides an easy way to retrieve the Admin Recommended Folder.

In order to use the Admin Recommended Folder, you should configure the provider to run in admin mode.
Please refer to the [Example Usage](#example-usage) section below for more details. 

## Example Usage
```hcl
# Configure the Sumo Logic Provider in Admin Mode
provider "sumologic" {
  ...
  admin_mode  = true
  alias       = "admin"
}

# Look up the Admin Recommended Folder
data "sumologic_admin_recommended_folder" "folder" {}

# Create a folder underneath the Admin Recommended Folder (which requires Admin Mode)
resource "sumologic_folder" "test" {
  provider    = sumologic.admin
  name        = "test"
  description = "A test folder"
  parent_id   = data.sumologic_admin_recommended_folder.folder.id
}
```


## Attributes reference

The following attributes are exported:

- `id` - The ID of the Admin Recommended Folder.
- `name` - The name of the Admin Recommended Folder.
- `description` - The description of the Admin Recommended Folder.



