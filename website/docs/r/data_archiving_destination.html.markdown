---
layout: "sumologic"
page_title: "SumoLogic: sumologic_data_archiving_destination"
description: |-
  Provides a Sumologic Data Archiving Destination
---

# sumologic_data_archiving_destination
Provides a [Sumologic Data Archiving Destination][1].

A data archiving destination describes where archived log data is written. The set of
supported arguments depends on the `destination_type` you choose.

## Example Usage

### S3 destination with role-based authentication
```hcl
resource "sumologic_data_archiving_destination" "s3_role_based" {
  destination_name = "archive-s3-role-based"

  destination_config {
    destination_type = "S3"
    bucket_name      = "my-archive-bucket"
    region           = "us-east-1"
    encrypted        = true
    enabled          = true

    auth_config {
      authentication_mode = "RoleBased"
      role_arn            = "arn:aws:iam::123456789012:role/my-archiving-role"
    }
  }
}
```

### S3 destination with access key authentication
```hcl
resource "sumologic_data_archiving_destination" "s3_access_key" {
  destination_name = "archive-s3-access-key"

  destination_config {
    destination_type = "S3"
    bucket_name      = "my-archive-bucket"
    region           = "us-east-1"
    encrypted        = true
    enabled          = true

    auth_config {
      authentication_mode = "AccessKey"
      access_key_id       = "your access key id"
      access_key_secret   = "your access key secret"
    }
  }
}
```

### Syslog destination
```hcl
resource "sumologic_data_archiving_destination" "syslog" {
  destination_name = "archive-syslog"

  destination_config {
    destination_type = "Syslog"
    protocol         = "udp"
    host             = "10.20.30.40"
    port             = 514
  }
}
```

### Hitachi destination
```hcl
resource "sumologic_data_archiving_destination" "hitachi" {
  destination_name = "archive-hitachi"

  destination_config {
    destination_type = "Hitachi"
    url              = "https://hitachi.example.com"
    object_id        = "archive_path/logname_{day}_{hour}_{minute}_{second}_{uuid}.log"
    username         = "your username"
    password         = "your password"
  }
}
```

### Generic REST API destination
```hcl
resource "sumologic_data_archiving_destination" "rest_api" {
  destination_name = "archive-rest-api"

  destination_config {
    destination_type = "RestAPI"
    url              = "https://example.com/receiver/v1/http/your_token"
    object_id        = "archive_test"
    username         = "your username"
    password         = "your password"
  }
}
```

## Argument Reference

The following arguments are supported:

- `destination_name` - (Required) Name of the data archiving destination. Must be between 1 and 128 characters.
- `destination_config` - (Required) Configuration of the destination. Only one block is allowed. See [destination_config](#destination_config) below.

### destination_config

- `destination_type` - (Required, ForceNew) The type of destination. Possible values are `S3`, `Syslog`, `Hitachi`, and `RestAPI`. Cannot be changed after creation.
- `description` - (Optional) Description of the destination.
- `bucket_name` - (Optional, ForceNew) The name of the Amazon S3 bucket. Required when `destination_type` is `S3`. Cannot be changed after creation.
- `region` - (Optional) The region where the S3 bucket is located.
- `encrypted` - (Optional) Enable server-side encryption. Must be set explicitly when `destination_type` is `S3`.
- `enabled` - (Optional) Whether the destination is enabled. Must be set explicitly when `destination_type` is `S3`.
- `auth_config` - (Optional) Authentication settings for the destination. Only one block is allowed. Required when `destination_type` is `S3`. See [auth_config](#auth_config) below.
- `protocol` - (Optional) Transport protocol. Possible values are `tcp` and `udp`. Required when `destination_type` is `Syslog`.
- `host` - (Optional) Hostname or IP address of the destination. Required when `destination_type` is `Syslog`.
- `port` - (Optional) Port of the destination. Must be between 1 and 65535. Required when `destination_type` is `Syslog`.
- `token` - (Optional, Sensitive) Token used to authenticate with a `Syslog` destination.
- `url` - (Optional) URL of the destination. Required when `destination_type` is `Hitachi` or `RestAPI`.
- `object_id` - (Optional) Object identifier at the destination. Required when `destination_type` is `Hitachi`.
- `username` - (Optional) Username used to authenticate with the destination. Required when `destination_type` is `Hitachi`. For `RestAPI` it is optional on creation but required for any later update, so set it up front if the destination will ever change.
- `password` - (Optional, Sensitive) Password used to authenticate with the destination. Required when `destination_type` is `Hitachi`.

### auth_config

- `authentication_mode` - (Required) AWS IAM authentication method used for access. Possible values are `AccessKey` and `RoleBased`.
- `access_key_id` - (Optional) The AWS Access Key ID used to access the S3 bucket. Used when `authentication_mode` is `AccessKey`.
- `access_key_secret` - (Optional, Sensitive) The AWS Secret Access Key used to access the S3 bucket. Used when `authentication_mode` is `AccessKey`.
- `role_arn` - (Optional) The AWS Role ARN used to access the S3 bucket. Used when `authentication_mode` is `RoleBased`.

## Attributes Reference

The following attributes are exported:

- `id` - The unique identifier of the data archiving destination.
- `created_at` - When the destination was created.
- `created_by` - The identifier of the user who created the destination.
- `modified_at` - When the destination was last modified.
- `modified_by` - The identifier of the user who last modified the destination.
- `destination_config.invalidated_by_system` - True when the system has invalidated the destination, for example because its credentials no longer work.

Sensitive values such as `token`, `password`, and `access_key_secret` are returned masked by the
API, so they are read back from state rather than from the API response.

## Import

Data archiving destinations can be imported using the destination ID.

```hcl
terraform import sumologic_data_archiving_destination.example 000000000ABC1234
```

[1]: https://help.sumologic.com/docs/manage/data-archiving/
