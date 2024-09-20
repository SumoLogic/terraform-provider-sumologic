---
layout: "sumologic"
page_title: "SumoLogic: sumologic_data_forwarding_destination"
description: |-
  Provides a Sumologic Data Forwarding Destination
---

# sumologic_data_forwarding_destination
Provider to manage [Sumologic Data Forwarding Destination](https://help.sumologic.com/docs/manage/data-forwarding/amazon-s3-bucket/#configure-an-s3-data-forwarding-destination)

## Example Usage
```hcl
resource "sumologic_data_forwarding_destination" "example_data_forwarding_destination" {
    destination_name = "df-destination"
    description = ""
    bucket_name = "df-bucket"
    region = "us-east-1"
    authentication_mode = "RoleBased"
    access_key_id = "accessKeyId"
    secret_access_key = "secretAccessKey"
    role_arn = "arn:aws:iam::some-valid-arn"
    encrypted = "false"
    enabled = "true"
}
```
## Argument reference

The following arguments are supported:

- `destination_name` - (Required) Name of the S3 data forwarding destination.
- `description` - (Optional) Description of the S3 data forwarding destination.
- `bucket_name` - (Required) The name of the Amazon S3 bucket.
- `region` - (Optional) The region where the S3 bucket is located.
- `authentication_mode` - (Required) AWS IAM authentication method used for access. Possible values are: 1. `AccessKey` 2. `RoleBased`
- `access_key_id` - (Optional) The AWS Access ID to access the S3 bucket.
- `secret_access_key` - (Optional) The AWS Secret Key to access the S3 bucket.
- `role_arn` - (Optional) The AWS Role ARN to access the S3 bucket.
- `encrypted` - (Optional) Enable S3 server-side encryption.
- `enabled` - (Optional) True when the data forwarding destination is enabled. Will be treated as _false_ if left blank.

The following attributes are exported:

- `id` - The internal ID of the data_forwarding_destination
