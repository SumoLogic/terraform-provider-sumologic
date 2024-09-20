---
layout: "sumologic"
page_title: "SumoLogic: sumologic_data_forwarding_destination"
description: |-
  Provides a Sumologic Data Forwarding Destination
---

# sumologic_data_forwarding_rule
Provider to manage [Sumologic Data Forwarding Rule](https://help.sumologic.com/docs/manage/data-forwarding/amazon-s3-bucket/#forward-datato-s3)

## Example Usage
```hcl
resource "sumologic_data_forwarding_rule" "example_data_forwarding_rule" {
    index_id = "00000000024C6155"
    destination_id = "00000000000732AA"
    enabled = "true"
    file_format = "test/{index}/{day}/{hour}/{minute}"
    payload_schema = "builtInFields"
    format = "json"
}
```
## Argument reference

The following arguments are supported:

- `index_id` - (Required) The _id_ of the Partition or Scheduled View the rule applies to.
- `destination_id` - (Required) The data forwarding destination id.
- `enabled` - (Optional) True when the data forwarding rule is enabled. Will be treated as _false_ if left blank.
- `file_format` - (Optional) Specify the path prefix to a directory in the S3 bucket and how to format the file name.
- `payload_schema` - (Optional) Schema for the payload. Default value of the payload schema is _allFields_ for scheduled view, and _builtInFields_ for partition.
  _raw_ payloadSchema should be used in conjunction with _text_ format and vice versa.
- `format` - (Optional) Format of the payload. Default format will be _csv_. 
  _text_ format should be used in conjunction with _raw_ payloadSchema and vice versa.
