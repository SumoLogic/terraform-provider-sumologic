---
layout: "sumologic"
page_title: "SumoLogic: sumologic_data_forwarding_destination"
description: |-
  Provides a Sumologic Data Forwarding Destination
---

# sumologic_data_forwarding_rule
Provider to manage [Sumologic Data Forwarding Rule](https://help.sumologic.com/docs/manage/data-forwarding/amazon-s3-bucket/#forward-datato-s3)

## Example Usage

For Partitions
```hcl
resource "sumologic_partition" "test_partition" {
  name               = "testing_rule_partitions"
  routing_expression = "_sourcecategory=abc/Terraform"
  is_compliant       = false
  retention_period   = 30
  analytics_tier     = "flex"
}

resource "sumologic_data_forwarding_rule" "example_data_forwarding_rule" {
  index_id = sumologic_partition.test_partition.id
  destination_id = "00000000000732AA"
  enabled = true
  file_format = "test/{index}/{day}/{hour}/{minute}"
  payload_schema = "builtInFields"
  format = "json"
}
```
For Scheduled Views
```hcl
resource "sumologic_scheduled_view" "failed_connections" {
  index_name = "failed_connections"
  query = "_sourceCategory=fire | count"
  start_time = "2024-09-01T00:00:00Z"
  retention_period = 1
  lifecycle {
    prevent_destroy = true
    ignore_changes = [index_id]
  }
}

resource "sumologic_data_forwarding_rule" "test_rule_sv" {
  index_id       = sumologic_scheduled_view.failed_connections.index_id
  destination_id = sumologic_data_forwarding_destination.test_destination.id
  enabled        = false
  file_format    = "test/{index}"
  payload_schema = "raw"
  format         = "text"
}
```
## Argument reference

The following arguments are supported:

- `index_id` - (Required) The *id* of the Partition or *index_id* of the Scheduled View the rule applies to.
- `destination_id` - (Required) The data forwarding destination id.
- `enabled` - (Optional) True when the data forwarding rule is enabled. Will be treated as _false_ if left blank.
- `file_format` - (Optional) Specify the path prefix to a directory in the S3 bucket and how to format the file name. For possible values, kindly refer the point 6 in the [documentation](https://help.sumologic.com/docs/manage/data-forwarding/amazon-s3-bucket/#forward-datato-s3).
- `payload_schema` - (Optional) Schema for the payload. Default value of the payload schema is _allFields_ for scheduled view, and _builtInFields_ for partition.
  _raw_ payloadSchema should be used in conjunction with _text_ format and vice versa.
- `format` - (Optional) Format of the payload. Default format will be _csv_. 
  _text_ format should be used in conjunction with _raw_ payloadSchema and vice versa.

The following attributes are exported:

- `id` - The Index ID of the data_forwarding_rule