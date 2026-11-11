---
layout: 'sumologic'
page_title: 'SumoLogic: sumologic_slo_folder'
description: |-
  Provides the ability to create, read, delete, and update folders for SLOs ([Beta feature](https://help.sumologic.com/Beta/SLO_Reliability_Management#beta-notices)).
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
  description = "folder for SLO created for testing"
  parent_id   = "0000000000000001"
}

resource "sumologic_slo" "slo_tf_test" {
  name        = "slo-tf-test1"
  description = "example SLO created for testing"
  parent_id   = sumologic_slo_folder.tf_slo_folder.id
  signal_type = "Latency"
  service     = "auth"
  application = "login"
  compliance {
    compliance_type = "Rolling"
    size            = "7d"
    target          = 99
    timezone        = "Asia/Kolkata"
  }
  indicator {
    window_based_evaluation {
      op          = "LessThan"
      query_type  = "Metrics"
      aggregation = "Avg"
      size        = "1m"
      threshold   = 200
      queries {
        query_group_type = "Threshold"
        query_group {
          row_id        = "A"
          query         = "metric=request_time_p90  service=auth api=login"
          use_row_count = false
        }
      }
    }
  }
}
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
- `type` - (Optional) The type of object model. Defaults to `SlosLibraryFolder`.
- `content_type` - (Optional) The type of the content object. Defaults to `Folder`.

Additional data provided in state:

- `id` - (Computed) The identifier for this SLO folder.
- `version` - (Computed) The version of the SLO folder, set by Sumo Logic.
- `is_mutable` - (Computed) Whether the SLO folder can be modified, as reported by Sumo Logic.
- `is_system` - (Computed) Whether the SLO folder is a system-provided one, as reported by Sumo Logic.
- `created_at` - (Computed) Time the SLO folder was created, as reported by Sumo Logic.
- `created_by` - (Computed) ID of the user who created the SLO folder.
- `modified_at` - (Computed) Time the SLO folder was last modified.
- `modified_by` - (Computed) ID of the user who last modified the SLO folder.

## Import

SLO folders can be imported using the SLO folder identifier, such as:

``` shell
terraform import sumologic_slo_folder.tf_slo_folder_1 0000000000ABC123
```

