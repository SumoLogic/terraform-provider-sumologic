---
layout: "sumologic"
page_title: "SumoLogic: sumologic_content_permission"
description: |-
  Provides a way to interact with Sumologic Content
---

# sumologic_content_permission
Provides a way for users to configure permissions on content.

__IMPORTANT:__  
* Based on [content permission APIs] (https://api.sumologic.com/docs/#tag/contentPermissions), when a new permission is added to a content, all lower level permissions are also added to that content. Because of this behavior, the best practice of updating this resource would be redefine the whole block instead of modifying it in place.
* While destroying this resource, all permissions related to the content get deleted except for the permissions of the content creator. 

 

## Example Usage
```hcl
data "sumologic_personal_folder" "personalFolder" {}

resource "sumologic_content" "permission_test_content" {
	parent_id = data.sumologic_personal_folder.personalFolder.id
	config = jsonencode({
		"type": "FolderSyncDefinition",
		"name": "test_permission_resource_folder",
		"description": "",
		"children": []
	})
}
	
data "sumologic_role" "role" { 
	name = "test-role"
}

data "sumologic_user" "user" { 
	email = "user@example.com"
}

resource "sumologic_content_permission" "content_permission_test" {
	content_id = sumologic_content.permission_test_content.id
	notify_recipient = true
	notification_message = "You now have the permission to access this content"
	permission {
		permission_name = "View"
		source_type = "role"
		source_id = data.sumologic_role.role.id
	}
	permission {
		permission_name = "Manage"
		source_type = "user"
		source_id = data.sumologic_user.user.id
	}
}
```

## Argument reference

The following arguments are supported:

- `content_id` - (Required) The identifier of the content item.
- `notify_recipient` - (Required) Boolean value. Set to "true" to notify the users who had a permission update.
- `notification_message` - (Required) The notification message sent to the users who had a permission update.

## Attributes reference

### Schema for `permission`
- `permission_name` - (Required) Content permission name. Valid values are: View, GrantView, Edit, GrantEdit, Manage, and GrantManage.
- `source_type` - (Required) Type of source for the permission. Valid values are: user, role, and org.
- `source_id` - (Required) An identifier that belongs to the source type chosen above. For e.g. if the sourceType is set to "user", sourceId should be identifier of a user (same goes for role and org sourceType)