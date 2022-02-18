---
layout: "sumologic"
page_title: "SumoLogic: sumologic_content_permission"
description: |-
  Provides a way to interact with Sumologic Content
---

# sumologic_content_permission
Provides a way to configure permissions on a content to share it with a user, a role, or the entire
org. You can read more [here](https://help.sumologic.com/Manage/Content_Sharing/Share-Content).

There are three permission levels `View`, `Edit` and `Manage`. You can read more about different
levels [here](https://help.sumologic.com/Manage/Content_Sharing/Share-Content#available-permission-levels).

-> When you add a new permission to a content, all the lower level permissions are added by default.
For example, giving a user "Manage" permission on a content, implicitly gives them "Edit" and "View"
permissions on the content. Due to this behavior, when you add a higher level permission, you must
also add all the lower level permissions. For example, when you give a user "Edit" permission via
the resource, you must give them "View" permission otherwise state and configuration will be out
of sync.


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
	name = "test_role"
}

data "sumologic_user" "user" {
	email = "user@example.com"
}

// Grant user `user@example.com` "Manage" permission and role `test_role`
// "View" permission on the folder `test_permission_resource_folder`.
resource "sumologic_content_permission" "content_permission_test" {
	content_id = sumologic_content.permission_test_content.id
	notify_recipient = true
	notification_message = "You now have the permission to access this content"

	permission {
		permission_name = "View"
		source_type = "role"
		source_id = data.sumologic_role.role.id
	}

	// Note: We are explicitly adding blocks for View and Edit permissions
	// because granting Manage permissions implicitly gives user View and
	// Edit permission on the content.
	permission {
		permission_name = "View"
		source_type = "user"
		source_id = data.sumologic_user.user.id
	}
	permission {
		permission_name = "Edit"
		source_type = "user"
		source_id = data.sumologic_user.user.id
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

- `content_id` - (Required) The identifier of the content item for which you want to update
permissions.
- `notify_recipient` - (Required) Boolean value. Set it to "true" to notify the recipients by email.
- `notification_message` - (Optional) The notification message to send to the users.
- `permission` - (Required) Permission block defining permission on the content. See
[permission schema](#schema-for-permission) for details.

### Schema for `permission`
- `permission_name` - (Required) Content permission name. Valid values are `View`, `GrantView`,
`Edit`, `GrantEdit`, `Manage`, and `GrantManage`. You can read more about permission levels
[here](https://help.sumologic.com/Manage/Content_Sharing/Share-Content#available-permission-levels).
- `source_type` - (Required) Type of source for the permission. Valid values are `user`, `role`,
and `org`.
- `source_id` - (Required) An identifier that belongs to the source type chosen above. For example,
if the `sourceType` is set to `user`, `sourceId` should be identifier of the user you want to share
content with (same goes for role and org source type).

## Import
Permisions on a content item can be imported using the content identifier, e.g.:
```hcl
// import permissions for content item with identifier = 0000000008E0183E
terraform import sumologic_content_permission.dashboard_permission_import 0000000008E0183E
```
