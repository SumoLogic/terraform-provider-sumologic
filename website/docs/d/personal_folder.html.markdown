---
layout: "sumologic"
page_title: "SumoLogic: sumologic_personal_folder"
description: |-
  Provides an easy way to retrieve the Personal Folder.
---

# sumologic_personal_folder
Provides an easy way to retrieve the Personal Folder.


## Example Usage
```hcl
data "sumologic_personal_folder" "personalFolder" {}
```


## Attributes reference

The following attributes are exported:

- `id` - The ID of the Personal Folder.
- `name` - The name of the Personal Folder.
- `description` - The description of the Personal Folder.



