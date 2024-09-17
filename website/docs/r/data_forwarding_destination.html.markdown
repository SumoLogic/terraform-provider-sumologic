---
layout: "sumologic"
page_title: "SumoLogic: sumologic_data_forwarding_destination"
description: |-
  Provides a Sumologic Data ForwardingDestination
---

# sumologic_dataForwardingDestination
Provider to manage Sumologic Data ForwardingDestinations

## Example Usage
```hcl
resource "sumologic_data_forwarding_destination" "example_data_forwarding_destination" {
    destination_name = "df-destination"
    role_arn = "roleArn"
    authentication_mode = "RoleBased"
    region = "us-east-1"
    secret_access_key = "secretAccessKey"
    bucket_name = "df-bucket"
    enabled = "true"
    description = ""
    access_key_id = "accessKeyId"
    encrypted = ""
}
```
## Argument reference

The following arguments are supported:

- `destination_name` - (Optional) Name of the S3 data forwarding destination.
- `role_arn` - (Optional) The AWS Role ARN to access the S3 bucket.
- `authentication_mode` - (Required) AWS IAM authentication method used for access. Possible values are: 1. `AccessKey` 2. `RoleBased`
- `region` - (Optional) The region where the S3 bucket is located.
- `secret_access_key` - (Optional) The AWS Secret Key to access the S3 bucket.
- `bucket_name` - (Required) The name of the Amazon S3 bucket.
- `enabled` - (Optional) True if the destination is Active.
- `description` - (Optional) Description of the S3 data forwarding destination.
- `access_key_id` - (Optional) The AWS Access ID to access the S3 bucket.
- `encrypted` - (Optional) Enable S3 server-side encryption.

The following attributes are exported:

- `id` - The internal ID of the data_forwarding_destination



[Back to Index][0]

[0]: ../README.md
