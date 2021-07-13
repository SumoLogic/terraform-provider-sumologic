---
layout: "sumologic"
page_title: "SumoLogic: sumologic_admin_recommended_folder"
description: |-
  Provides an easy way to retrieve the Admin Recommended Folder.
---

# sumologic_personal_folder
Provides an easy way to retrieve the Personal Folder.


## Example Usage
```hcl
data "sumologic_admin_recommended_folder" "folder" {}
```


## Attributes reference

The following attributes are exported:

- `id` - The ID of the Personal Folder.
- `name` - The name of the Personal Folder.
- `description` - The description of the Personal Folder.



