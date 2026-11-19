---
layout: "sumologic"
page_title: "SumoLogic: sumologic_s3_audit_source"
description: |-
  Provides a AWS S3 Audit Source.
---

# sumologic_s3_audit_source
Provides a [AWS S3 Audit Source][2].

__IMPORTANT:__ The AWS credentials are stored in plain-text in the state. This is a potential security issue.

## Example Usage

### Amazon S3 Audit source with an access key
```hcl
resource "sumologic_s3_audit_source" "terraform_s3_audit" {
  name          = "Amazon S3 Audit"
  description   = "My description"
  category      = "aws/s3audit"
  content_type  = "AwsS3AuditBucket"
  scan_interval = 300000
  paused        = false
  collector_id  = sumologic_collector.collector.id

  authentication {
    type = "S3BucketAuthentication"
    access_key = var.aws_access_key
    secret_key = var.aws_secret_key
  }

  path {
    type = "S3BucketPathExpression"
    bucket_name     = "my-s3-access-logs-bucket"
    path_expression = "*"
  }
}

resource "sumologic_collector" "collector" {
  name        = "my-collector"
  description = "Just testing this"
}

variable "aws_access_key" {
  type = string
}

variable "aws_secret_key" {
  type      = string
  sensitive = true
}
```

### Amazon S3 Audit source with an IAM role
```hcl
resource "sumologic_collector" "collector" {
  name        = "my-collector"
  description = "Just testing this"
}

resource "sumologic_s3_audit_source" "terraform_s3_audit_role" {
  name         = "Amazon S3 Audit (role)"
  description  = "My description"
  category     = "aws/s3audit"
  content_type = "AwsS3AuditBucket"
  collector_id = sumologic_collector.collector.id

  authentication {
    type     = "AWSRoleBasedAuthentication"
    role_arn = "arn:aws:iam::123456789012:role/sumo-s3-audit-role"
  }

  path {
    type            = "S3BucketPathExpression"
    bucket_name     = "my-s3-access-logs-bucket"
    path_expression = "logs/*"
  }
}
```

## Argument reference

In addition to the [Common Source Properties](https://registry.terraform.io/providers/SumoLogic/sumologic/latest/docs#common-source-properties), the following arguments are supported:

 - `content_type` - (Required, forces a new resource) The content-type of the collected data. Must be `AwsS3AuditBucket`.
 - `scan_interval` - (Optional) Time interval in milliseconds of scans for new data. Defaults to `300000`.
 - `paused` - (Optional) When set to `true`, the scanner is paused. Defaults to `false`.
 - `authentication` - (Required, forces a new resource; exactly one block) Authentication details for connecting to the S3 bucket.
     + `type` - (Required) Must be either `S3BucketAuthentication` or `AWSRoleBasedAuthentication`.
     + `access_key` - (Optional) Required when `type` is `S3BucketAuthentication`. Your AWS access key if using type `S3BucketAuthentication`.
     + `secret_key` - (Optional) Required when `type` is `S3BucketAuthentication`. Your AWS secret key if using type `S3BucketAuthentication`.
     + `role_arn` - (Optional) Required when `type` is `AWSRoleBasedAuthentication`. Your AWS role ARN if using type `AWSRoleBasedAuthentication`.This is not supported for AWS China regions.
     + `region` - (Optional) Your AWS Bucket region.
 - `path` - (Required, forces a new resource; exactly one block) The location to scan for new data.
     + `type` - (Required) type of polling source. This has to be `S3BucketPathExpression` for `S3 Audit source`.
     + `bucket_name` - (Required) The name of the bucket. 
     + `path_expression` - (Required) The path to the data.
     + `sns_topic_or_subscription_arn` - (Computed, read-only) SNS topic/subscription ARN, set by Sumo Logic when SNS notification is configured; list of objects with `is_success` and `arn`.

### See also
   * [Common Source Properties](https://registry.terraform.io/providers/SumoLogic/sumologic/latest/docs#common-source-properties)
   * [Configuring SNS Subscription](https://registry.terraform.io/providers/SumoLogic/sumologic/latest/docs#configuring-sns-subscription)

## Attributes Reference
The following attributes are exported:

- `id` - The internal ID of the source.
- `url` - The HTTP endpoint to use with [SNS to notify Sumo Logic of new files](https://help.sumologic.com/03Send-Data/Sources/02Sources-for-Hosted-Collectors/Amazon-Web-Services/AWS-S3-Source#Set_up_SNS_in_AWS_(Optional)).

## Import
S3 Audit sources can be imported using the collector and source IDs (`collector/source`), e.g.:

```hcl
terraform import sumologic_s3_audit_source.test 123/456
```

S3 Audit sources can be imported using the collector name and source name (`collectorName/sourceName`), e.g.:

```hcl
terraform import sumologic_s3_audit_source.test my-test-collector/my-test-source
```

[2]: https://www.sumologic.com/help/docs/send-data/hosted-collectors/amazon-aws/amazon-s3-audit-source/