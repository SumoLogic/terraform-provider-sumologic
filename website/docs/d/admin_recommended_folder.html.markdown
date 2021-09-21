---
layout: "sumologic"
page_title: "SumoLogic: sumologic_admin_recommended_folder"
description: |-
  Provides an easy way to retrieve the Admin Recommended Folder.
---

# sumologic_admin_recommended_folder
Provides an easy way to retrieve the Admin Recommended Folder.


## Example Usage
```hcl
data "sumologic_admin_recommended_folder" "folder" {}
```


## Attributes reference

The following attributes are exported:

- `id` - The ID of the Admin Recommended Folder.
- `name` - The name of the Admin Recommended Folder.
- `description` - The description of the Admin Recommended Folder.



