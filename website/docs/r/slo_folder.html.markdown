---
layout: 'sumologic'
page_title: 'SumoLogic: sumologic_slo_folder'
description: |-
  Provides the ability to create, read, delete, and update folders for SLOs.
---

# sumologic_slo_folder

Provides the ability to create, read, delete, and update folders for SLO's.

## Example SLO Folder

NOTE: SLO folders are considered a different resource from Library content and monitor folders.

```hcl
resource "sumologic_slo_folder" "tf_slo_folder" {
  name        = "Terraform Managed SLO's"
  description = "A folder for SLO's managed by terraform."
}
```

## Example Creating SLO Folders and then SLO's in them

```hcl

resource "sumologic_slo_folder" "tf_slo_folder" {
  name        = "slo-tf-folder"
  description = "Root folder for SLO created for testing"
  parent_id   = "0000000000000001"
}

resource "sumologic_slo" "slo_tf_test" {
  name        = "slo-tf-test1"
  description = "example SLO created for testing"
  parent_id   = sumologic_slo_folder.slo_tf_test_folder.id
 ...
```

## Example Nested SLO Folders

```hcl
resource "sumologic_slo_folder" "tf_payments_team_root_folder" {
  name        = "Security Team SLOs"
  description = "SLO's payments services."
}

resource "sumologic_slo_folder" "tf_payments_team_prod_folder" {
  name        = "Production SLOs"
  description = "SLOs for the Payments service on Production Environment."
  parent_id   = sumologic_slo_folder.tf_payments_team_root_folder.id
}

resource "sumologic_slo_folder" "tf_payments_team_stag_folder" {
  name        = "Staging SLOs"
  description = "SLOs for the payments service on Staging Environment."
  parent_id   = sumologic_slo_folder.tf_payments_team_root_folder.id
}
```

## Argument reference

The following arguments are supported:

- `name` - (Required) The name of the SLO folder. The name must be alphanumeric.
- `description` - (Required) The description of the SLO folder.
- `parent_id` - (Optional) The identifier of the SLO Folder that contains this SLO Folder. Defaults to the root folder.

Additional data provided in state:

- `id` - (Computed) The identifier for this SLO folder.

## Import

SLO folders can be imported using the SLO folder identifier, such as:

``` shell
terraform import sumologic_slo_folder.tf_slo_folder_1 0000000000ABC123
```

