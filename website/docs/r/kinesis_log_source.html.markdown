---
layout: "sumologic"
page_title: "SumoLogic: sumologic_kinesis_log_source"
description: |-
  Provides a Sumologic Kinesis Log source. This source is used to integrate with Log Stream via Kinesis Firehose from AWS.
---

# sumologic_kinesis_log_source

Provides a [Sumologic Kinesis Log source][2]. This source is used to ingest log via Kinesis Firehose from AWS.

__IMPORTANT:__ The AWS credentials are stored in plain-text in the state. This is a potential security issue.

## Example Usage

Each example declares the hosted collector it attaches the source to, so it can be copied on its own.

### Kinesis Firehose logs, with failed deliveries read from S3 using an access key
```hcl
resource "sumologic_collector" "collector" {
  name        = "my-collector"
  description = "Hosted collector for Kinesis logs"
}

variable "aws_access_key" {
  type = string
}

variable "aws_secret_key" {
  type      = string
  sensitive = true
}

resource "sumologic_kinesis_log_source" "kinesis_log_access_key" {
  name         = "Kinesis Log"
  description  = "Description for Kinesis Log Source"
  category     = "prod/kinesis/log"
  content_type = "KinesisLog"
  collector_id = sumologic_collector.collector.id

  authentication {
    type       = "S3BucketAuthentication"
    access_key = var.aws_access_key
    secret_key = var.aws_secret_key
    region     = "us-east-1"
  }

  path {
    type            = "KinesisLogPath"
    bucket_name     = "my-firehose-backup-bucket"
    path_expression = "http-endpoint-failed/*"
    scan_interval   = 30000
  }
}
```

### Kinesis Firehose logs, with failed deliveries read from S3 using an IAM role
```hcl
resource "sumologic_collector" "collector" {
  name        = "my-collector"
  description = "Hosted collector for Kinesis logs"
}

resource "sumologic_kinesis_log_source" "kinesis_log_role_arn" {
  name         = "Kinesis Log"
  description  = "Description for Kinesis Log Source"
  category     = "prod/kinesis/log"
  content_type = "KinesisLog"
  collector_id = sumologic_collector.collector.id

  authentication {
    type     = "AWSRoleBasedAuthentication"
    role_arn = "arn:aws:iam::123456789012:role/sumo-kinesis-role"
  }

  path {
    type            = "KinesisLogPath"
    bucket_name     = "my-firehose-backup-bucket"
    path_expression = "http-endpoint-failed/*"
    scan_interval   = 30000
  }
}
```

### Kinesis Firehose logs only (no S3 backup)
```hcl
resource "sumologic_collector" "collector" {
  name        = "my-collector"
  description = "Hosted collector for Kinesis logs"
}

resource "sumologic_kinesis_log_source" "kinesis_log_basic" {
  name         = "Kinesis Log"
  category     = "prod/kinesis/log"
  content_type = "KinesisLog"
  collector_id = sumologic_collector.collector.id
}
```

## Argument reference

In addition to the [Common Source Properties](https://registry.terraform.io/providers/SumoLogic/sumologic/latest/docs#common-source-properties), the following arguments are supported:

 - `content_type` - (Required) The content-type of the collected data. Must be `KinesisLog`.
 - `message_per_request` - (Optional) Set to `true` if each HTTP request from Firehose carries a single message. Defaults to `false`.
 - `authentication` - (Optional, at most one block) Authentication details for connecting to the S3 bucket that holds failed deliveries. Defaults to `type = "NoAuthentication"`.
     + `type` - (Optional) One of `S3BucketAuthentication`, `AWSRoleBasedAuthentication` or `NoAuthentication`. Defaults to `NoAuthentication`.
     + `access_key` - (Optional) Required when `type` is `S3BucketAuthentication`. Your AWS access key if using type `S3BucketAuthentication`
     + `secret_key` - (Optional) Required when `type` is `S3BucketAuthentication`. Your AWS secret key if using type `S3BucketAuthentication`
     + `role_arn` - (Optional) Required when `type` is `AWSRoleBasedAuthentication`. Your AWS role ARN if using type `AWSRoleBasedAuthentication`
     + `region` - (Optional) Your AWS bucket region.
 - `path` - (Optional, at most one block) The location of the S3 bucket for failed Kinesis log data. Defaults to `type = "NoPathExpression"`.
     + `type` - (Optional) `KinesisLogPath` or `NoPathExpression`. Defaults to `NoPathExpression`.
     + `bucket_name` - (Optional) The name of the bucket. This is needed if using type `KinesisLogPath`. 
     + `path_expression` - (Optional) The path to the data. This is needed if using type `KinesisLogPath`. For Kinesis log source, it must include `http-endpoint-failed/`, for example `http-endpoint-failed/*`. Do not use the same bucket and prefix for different Kinesis sources.
     + `scan_interval` - (Optional) Time interval in milliseconds of scans of the S3 path for new data. No default is set in the provider.

### See also
   * [Common Source Properties](https://registry.terraform.io/providers/SumoLogic/sumologic/latest/docs#common-source-properties)
   * [AWS Kinesis Firehose for Logs Source](https://www.sumologic.com/help/docs/send-data/hosted-collectors/amazon-aws/aws-kinesis-firehose-logs-source/)

## Attributes Reference
The following attributes are exported:

- `id` - The internal ID of the source.
- `url` - The HTTP endpoint to be used while creating Kinesis Firehose on AWS.

## Import
Kinesis Log sources can be imported using the collector and source IDs (`collector/source`), e.g.:

```hcl
terraform import sumologic_kinesis_log_source.test 123/456
```

Kinesis Log sources can be imported using the collector name and source name (`collectorName/sourceName`), e.g.:

```hcl
terraform import sumologic_kinesis_log_source.test my-test-collector/my-test-source
```

[2]: https://www.sumologic.com/help/docs/send-data/hosted-collectors/amazon-aws/aws-kinesis-firehose-logs-source/
