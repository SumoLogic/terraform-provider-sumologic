---
layout: "sumologic"
page_title: "SumoLogic: sumologic_folder"
description: |-
  Provides an easy way to retrieve a folder.
---

# sumologic_folder
Provides an easy way to retrieve a folder.

You must specify the absolute path of the folder to retrieve. The content library has "Library"
folder at the root level. For items in "Personal" folder, the base path is "/Library/Users/user@sumologic.com"
where "user@sumologic.com" is the email address of the user. For example, if a user with email address
`wile@acme.com` has `Rockets` folder inside Personal folder, the path of Rockets folder will be
`/Library/Users/wile@acme.com/Rockets`.

For items in "Admin Recommended" folder, the base path is "/Library/Admin Recommended". For example,
given a folder `Acme` in Admin Recommended folder, the path will be `/Library/Admin Recommended/Acme`.


## Example Usage
```hcl
provider "sumologic" {
  environment = "us2"
  access_id = "..."
  access_key = "..."
}

# Provider with admin mode set to true
provider "sumologic" {
  ...
  admin_mode  = true
  alias       = "admin"
}

# Look up folder named "Rockets" under Personal folder of user "wile@acme.com"
data "sumologic_folder" "rockets" {
  path = "/Library/Users/wile@acme.com/Rockets"
}

# Look up folder named "Acme" under "Admin Recommended" folder (must use provider with
admin mode set to true)
data "sumologic_folder" "acme" {
  provider    = sumologic.admin
  path = "/Library/Admin Recommended/Acme"
}
```


## Attributes reference

The following attributes are exported:

- `id` - The ID of the folder.
- `name` - The name of the folder.



