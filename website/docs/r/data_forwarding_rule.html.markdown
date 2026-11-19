---
layout: "sumologic"
page_title: "SumoLogic: sumologic_data_forwarding_rule"
description: |-
  Provides a Sumologic Data Forwarding Rule
---

# sumologic_data_forwarding_rule
Provider to manage a [Sumologic Data Forwarding Rule](https://www.sumologic.com/help/docs/manage/data-forwarding/forward-data-from-sumologic/)

A data forwarding rule sends the data of one Partition or Scheduled View to an Amazon S3 bucket, described by a [`sumologic_data_forwarding_destination`](data_forwarding_destination.html). Each Partition or Scheduled View can have at most one rule, which is why the rule is identified by its `index_id`.

## Example Usage

### Forward a Partition
```hcl
resource "sumologic_data_forwarding_destination" "example_destination" {
  destination_name = "df-destination"
  description      = "S3 bucket for forwarded data"
  bucket_name      = "df-bucket"
  s3_region        = "us-east-1"
  authentication {
    type     = "RoleBased"
    role_arn = "arn:aws:iam::123456789012:role/my-forwarding-role"
  }
  s3_server_side_encryption = false
  enabled                   = true
}

resource "sumologic_partition" "test_partition" {
  name               = "testing_rule_partitions"
  routing_expression = "_sourcecategory=abc/Terraform"
  is_compliant       = false
  retention_period   = 30
  analytics_tier     = "flex"
}

resource "sumologic_data_forwarding_rule" "example_data_forwarding_rule" {
  index_id       = sumologic_partition.test_partition.id
  destination_id = sumologic_data_forwarding_destination.example_destination.id
  enabled        = true
  file_format    = "test/{index}/{day}/{hour}/{minute}"
  payload_schema = "builtInFields"
  format         = "json"
}
```

### Forward a Scheduled View
Use the `index_id` of the Scheduled View, not its `id`. This example forwards the raw log lines as text.
```hcl
resource "sumologic_data_forwarding_destination" "example_destination" {
  destination_name = "df-destination"
  bucket_name      = "df-bucket"
  s3_region        = "us-east-1"
  authentication {
    type     = "RoleBased"
    role_arn = "arn:aws:iam::123456789012:role/my-forwarding-role"
  }
  enabled = true
}

resource "sumologic_scheduled_view" "failed_connections" {
  index_name       = "failed_connections"
  query            = "_sourceCategory=fire | count"
  start_time       = "2024-09-01T00:00:00Z"
  retention_period = 1
  lifecycle {
    prevent_destroy = true
    ignore_changes  = [index_id]
  }
}

resource "sumologic_data_forwarding_rule" "test_rule_sv" {
  index_id       = sumologic_scheduled_view.failed_connections.index_id
  destination_id = sumologic_data_forwarding_destination.example_destination.id
  enabled        = false
  file_format    = "test/{index}"
  payload_schema = "raw"
  format         = "text"
}
```

## Argument reference

The following arguments are supported:

- `index_id` - (Required) The ID of the Partition, or the `index_id` of the Scheduled View, the rule applies to. Exactly 16 characters.
- `destination_id` - (Required) The ID of the [`sumologic_data_forwarding_destination`](data_forwarding_destination.html) the data is forwarded to. Exactly 16 characters.
- `enabled` - (Optional) `true` when the data forwarding rule is enabled. Treated as `false` if left blank, so the rule is created disabled.
- `file_format` - (Optional) The path prefix to a directory in the S3 bucket, and how to format the file name. For possible values, you can use the variables `{index}` (name of the Partition or Scheduled View), `{day}` (`yyyy-MM-dd`), `{hour}` (0-23), `{minute}`, `{second}` and `{uuid}` (a random UUID), for example `SumoDataForwarding/{day}/{index}_{day}_{hour}_{minute}_{second}`. If left blank, the default `{index}_{day}_{hour}_{minute}_{second}` is used (see [Forward Data from Sumo Logic to S3 or GCS](https://www.sumologic.com/help/docs/manage/data-forwarding/forward-data-from-sumologic/)).
- `payload_schema` - (Optional) Schema for the payload. Allowed values: `builtInFields`, `allFields`, `raw`. The default is `allFields` for a Scheduled View and `builtInFields` for a Partition. In the Sumo Logic UI this is the Included Data option: `raw` is Raw, `builtInFields` is Raw + Metadata (the UI's recommended option), and `allFields` is All (Raw + Metadata + Enriched Fields). `raw` must be used together with `format = "text"`, and the other way round.
- `format` - (Optional) Format of the payload. Allowed values: `csv`, `json`, `text`. The default is `csv`. `text` must be used together with `payload_schema = "raw"`, and the other way round. Files are compressed with gzip. Data in the Infrequent Tier cannot be forwarded.

The following attributes are exported:

- `id` - The Index ID of the data_forwarding_rule, the same value as `index_id`.

See the [Logs Data Forwarding Management API](https://api.sumologic.com/docs/#tag/logsDataForwardingManagement) for the underlying operations.

## Destroy

A new or changed rule takes about five minutes to propagate, so the first five minutes of data after a change might not be forwarded. Only new data is forwarded after a rule is created. Destroying the resource deletes the data forwarding rule for the Partition or Scheduled View. The Partition, Scheduled View and destination are not changed.
