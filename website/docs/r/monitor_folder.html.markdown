---
layout: 'sumologic'
page_title: 'SumoLogic: sumologic_monitor_folder'
description: |-
  Provides the ability to create, read, delete, and update folders for Monitors.
---

# sumologic_monitor_folder

Provides the ability to create, read, delete, and update folders for [Monitors][1].

If Fine Grain Permission (FGP) feature is enabled with Monitors Content at one's Sumo Logic account, one can also set those permission details under this monitor folder resource. 

For further details about FGP, please see this [Monitor Permission document][2]. 

## Example Monitor Folder

NOTE: Monitor folders are considered a different resource from Library content folders.

```hcl
resource "sumologic_monitor_folder" "tf_monitor_folder_1" {
  name        = "Terraform Managed Monitors"
  description = "A folder for monitors managed by terraform."
}
```

## Example Nested Monitor Folders

NOTE: 
- Monitor folders allow up to six (6) levels of sub-folders.
- `obj_permission` are added at one of the Folders to showcase how Fine Grain Permissions (FGP) are associated with two roles. 


```hcl
resource "sumologic_role" "tf_test_role_01" {
  name        = "tf_test_role_01"
  description = "Testing resource sumologic_role"
  capabilities = [
    "viewAlerts",
    "viewMonitorsV2",
    "manageMonitorsV2"
  ]
}

resource "sumologic_role" "tf_test_role_02" {
  name        = "tf_test_role_02"
  description = "Testing resource sumologic_role"
  capabilities = [
    "viewAlerts",
    "viewMonitorsV2",
    "manageMonitorsV2"
  ]
}

resource "sumologic_monitor_folder" "tf_security_team_root_folder" {
  name        = "Security Team Monitors"
  description = "Monitors used by the Security Team."
}

resource "sumologic_monitor_folder" "tf_security_team_prod_folder" {
  name        = "Production Monitors"
  description = "Monitors for the Security Team Production Environment."
  parent_id   = sumologic_monitor_folder.tf_security_team_root_folder.id
  obj_permission {
    subject_type = "role"
    subject_id = sumologic_role.tf_test_role_01.id 
    permissions = ["Create","Read","Update"] 
  }
  obj_permission {
    subject_type = "role"
    subject_id = sumologic_role.tf_test_role_02.id
    permissions = ["Create", "Read"]
  }
}

resource "sumologic_monitor_folder" "tf_security_team_stag_folder" {
  name        = "Staging Monitors"
  description = "Monitors for the Security Team Staging Environment."
  parent_id   = sumologic_monitor_folder.tf_security_team_root_folder.id
```

## Argument reference

The following arguments are supported:

- `type` - (Optional) The type of object model. Valid value:
  - `MonitorsLibraryFolder`
- `name` - (Required) The name of the monitor folder. The name must be alphanumeric.
- `description` - (Required) The description of the monitor folder.
- `parent_id` - (Optional) The identifier of the Monitor Folder that contains this Monitor Folder. Defaults to the root folder.
- `obj_permission` - (Optional) `obj_permission` construct represents a Permission Statement associated with this Folder. A set of `obj_permission` constructs can be specified under a single Folder. An `obj_permission` construct can be used to control permissions Explicitly associated with a Folder. But, it cannot be used to control permissions Inherited from a Parent / Ancestor Folder.  Default FGP would be still set to the Folder upon creation (e.g. the creating user would have full permission), even if no `obj_permission` construct is specified at a Folder and the FGP feature is enabled at the account. 
  - `subject_type` - (Required) Valid values: 
    - `role` 
    - `org` 
  - `subject_id` - (Required) A Role ID or the Org ID of the account 
  - `permissions` - (Required) A Set of Permissions. Valid Permission Values: 
    - `Create`
    - `Read`
    - `Update` 
    - `Delete` 
    - `Manage`

Additional data provided in state:

- `id` - (Computed) The identifier for this monitor folder.

## Import

Monitor folders can be imported using the monitor folder identifier, such as:

```hcl
terraform import sumologic_monitor_folder.tf_monitor_folder_1 0000000000ABC123
```

[1]: https://help.sumologic.com/?cid=10020
[2]: https://help.sumologic.com/Beta/Capabilities_and_Permissions_for_Monitors#set-permissions-for-a-monitors-folder
