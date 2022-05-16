---
layout: "sumologic"
page_title: "SumoLogic: sumologic_kinesis_metrics_source"
description: |-
  Provides a Sumologic Kinesis Metrics source. This source is used to integrate with Metrics Stream via Kinesis Firehose from AWS.
---

# sumologic_kinesis_metrics_source

Provides a Sumologic Kinesis Metrics source. This source is used to ingest data from Cloudwatch Metrics Stream via Kinesis Firehose from AWS.

__IMPORTANT:__ The AWS credentials are stored in plain-text in the state. This is a potential security issue.

## Example Usage
```hcl
locals {
  tagfilters = [{
          "type" = "TagFilters"
          "namespace" = "All"
          "tags" = ["k3=v3"]
        },{
          "type" = "TagFilters"
          "namespace" = "AWS/Route53"
          "tags" = ["k1=v1"]
        },{
          "type" = "TagFilters"
          "namespace" = "AWS/S3"
          "tags" = ["k2=v2"]
        }]
}

resource "sumologic_kinesis_metrics_source" "kinesis_metrics_access_key" {
  name          = "Kinesis Metrics"
  description   = "Description for Kinesis Metrics Source"
  category      = "prod/kinesis/metrics"
  content_type  = "KinesisMetric"
  collector_id  = "${sumologic_collector.collector.id}"
  authentication {
    type = "S3BucketAuthentication"
    access_key = "someKey"
    secret_key = "******"
  }

  path {
    type = "KinesisMetricPath"
    tag_filters {
        type = "TagFilters"
        namespace = "All"
        tags = ["k3=v3"]
    }
    tag_filters {
        type = "TagFilters"
        namespace = "AWS/Route53"
        tags = ["k1=v1"]
    }
  }
}

resource "sumologic_kinesis_metrics_source" "kinesis_metrics_role_arn" {
  name          = "Kinesis Metrics"
  description   = "Description for Kinesis Metrics Source"
  category      = "prod/kinesis/metrics"
  content_type  = "KinesisMetric"
  collector_id  = "${sumologic_collector.collector.id}"

  authentication {
    type = "AWSRoleBasedAuthentication"
    role_arn = "arn:aws:iam::604066827510:role/cw-role-SumoRole-4AOLS73TGKYI"
  }

  path {
    type = "KinesisMetricPath"
    tag_filters {
        type = "TagFilters"
        namespace = "All"
        tags = ["k3=v3"]
    }
    tag_filters {
        type = "TagFilters"
        namespace = "AWS/Route53"
        tags = ["k1=v1"]
    }
  }
}

resource "sumologic_collector" "collector" {
  name        = "my-collector"
  description = "Just testing this"
}
```

## Argument reference

In addition to the [Common Source Properties](https://registry.terraform.io/providers/SumoLogic/sumologic/latest/docs#common-source-properties), the following arguments are supported:

 - `content_type` - (Required) The content-type of the collected data. Details can be found in the [Sumologic documentation for hosted sources][1].
 - `authentication` - (Required) Authentication details for connecting to the S3 bucket.
     + `type` - (Required) Must be either `S3BucketAuthentication` or `AWSRoleBasedAuthentication`
     + `access_key` - (Required) Your AWS access key if using type `S3BucketAuthentication`
     + `secret_key` - (Required) Your AWS secret key if using type `S3BucketAuthentication`
     + `role_arn` - (Required) Your AWS role ARN if using type `AWSRoleBasedAuthentication`
 - `path` - (Required) The location to scan for new data.
     + `type` - (Required) Must be `KinesisMetricPath`
     + `tag_filters` - (Optional) Tag filters allow you to filter the CloudWatch metrics you collect by the AWS tags you have assigned to your AWS resources. You can define tag filters for each supported namespace. If you do not define any tag filters, all metrics will be collected for the regions and namespaces you configured for the source above. More info on tag filters can be found [here](https://help.sumologic.com/03Send-Data/Sources/02Sources-for-Hosted-Collectors/Amazon-Web-Services/Amazon-CloudWatch-Source-for-Metrics#about-aws-tag-filtering)
          + `type` - This value has to be set to `TagFilters`
          + `namespace` - Namespace for which you want to define the tag filters. Use  value as `All` to apply the tag filter for all namespaces.
          + `tags` - List of key-value pairs of tag filters. Eg: `["k3=v3"]`

### See also
   * [Common Source Properties](https://registry.terraform.io/providers/SumoLogic/sumologic/latest/docs#common-source-properties)

## Attributes Reference
The following attributes are exported:

- `id` - The internal ID of the source.
- `url` - The HTTP endpoint to used while creating Kinesis Firehose on AWS.

## Import
Kinesis Metrics sources can be imported using the collector and source IDs (`collector/source`), e.g.:

```hcl
terraform import sumologic_kinesis_metrics_source.test 123/456
```

HTTP sources can be imported using the collector name and source name (`collectorName/sourceName`), e.g.:

```hcl
terraform import sumologic_kinesis_metrics_source.test my-test-collector/my-test-source
```

## Full Example (Including terraform for AWS asset creation)
```hcl
terraform {
  required_providers {
    sumologic = {
      source = "sumologic/sumologic"
    }
    aws = {
      source = "hashicorp/aws"
    }
  }
}

provider "sumologic" {}
provider "aws" {}

locals {
  account_id = ""
  aws_access_key = ""
  aws_secret_key = ""

  description = "update your terraform description here"
  identifier = "SumologicMetricStream"

  region = "us-west-2"

  tagfilters = [
    { type = "TagFilters", namespace = "AWS/ApplicationELB", tags = ["Deployment=prod"] },
  ]

}

resource "sumologic_collector" "collector_for_kinesis_metrics" {
  name = "AWS Metrics via Kinesis"
}

resource "sumologic_kinesis_metrics_source" "kinesis_source" {
  name          = "CloudWatch Metrics via Kinesis"
  description   = "Description for Sumologic source"
  category      = "aws/cloudwatch"
  content_type  = "KinesisMetric"
  collector_id  = sumologic_collector.collector_for_kinesis_metrics.id

  authentication {
    type = "S3BucketAuthentication"
    access_key = local.aws_access_key
    secret_key = local.aws_secret_key
  }

  path {
    type = "KinesisMetricPath"

    dynamic "tag_filters" {
      for_each = local.tagfilters
      content {
        type      = tag_filters.value.type
        namespace = tag_filters.value.namespace
        tags      = tag_filters.value.tags
      }
    }
  }
}

// ------------------------------------ AWS Kinesis part

resource "aws_cloudwatch_metric_stream" "main" {
  name          = local.identifier
  role_arn      = aws_iam_role.metric_stream_to_firehose.arn
  firehose_arn  = aws_kinesis_firehose_delivery_stream.kinesis_stream.arn
  output_format = "opentelemetry0.7"

// Edit and uncomment below lines to add include_filter (or exclude_filter on similar lines)
//  include_filter {
//    namespace = "AWS/ApplicationELB"
//  }
//  include_filter {
//    namespace = "AWS/DynamoDB"
//  }

}

resource "aws_iam_role" "metric_stream_to_firehose" {
  name = "${local.identifier}-stream_to_firehose"

  assume_role_policy = <<EOF
{
  "Version": "2012-10-17",
  "Statement": [
    {
      "Action": "sts:AssumeRole",
      "Principal": {
        "Service": "streams.metrics.cloudwatch.amazonaws.com"
      },
      "Effect": "Allow",
      "Sid": ""
    }
  ]
}
EOF
}

resource "aws_iam_role_policy" "metric_stream_to_firehose" {
  name = "default"
  role = aws_iam_role.metric_stream_to_firehose.id

  policy = <<EOF
{
    "Version": "2012-10-17",
    "Statement": [
        {
            "Effect": "Allow",
            "Action": [
                "firehose:PutRecord",
                "firehose:PutRecordBatch"
            ],
            "Resource": "${aws_kinesis_firehose_delivery_stream.kinesis_stream.arn}"
        }
    ]
}
EOF
}

resource "aws_iam_role" "firehose_role" {
  name = "${local.identifier}_firehose"

  assume_role_policy = <<EOF
{
  "Version": "2012-10-17",
  "Statement": [
    {
      "Action": "sts:AssumeRole",
      "Principal": {
        "Service": "firehose.amazonaws.com"
      },
      "Effect": "Allow",
      "Sid": ""
    }
  ]
}
EOF
}

resource "aws_iam_role_policy" "firehose_can_log_errors_to_Cloudwatch" {
  role = aws_iam_role.firehose_role.id

  policy = <<EOF
{
    "Version": "2012-10-17",
    "Statement": [
        {
            "Effect": "Allow",
            "Action": [
                "logs:PutLogEvents"
            ],
            "Resource": [
                "arn:aws:logs:${local.region}:${local.account_id}:log-group:/aws/kinesisfirehose/${local.identifier}:*",
                "arn:aws:logs:${local.region}:${local.account_id}:log-group:/aws/kinesisfirehose/${local.identifier}:*:log-stream:*"
            ]
        }
    ]
}
EOF
}

resource "aws_iam_role_policy" "firehose_can_use_s3_bucket_for_failures" {
  role = aws_iam_role.firehose_role.id

  policy = <<EOF
{
    "Version": "2012-10-17",
    "Statement": [
        {
            "Action": [
                "s3:AbortMultipartUpload",
                "s3:GetBucketLocation",
                "s3:GetObject",
                "s3:ListBucket",
                "s3:ListBucketMultipartUploads",
                "s3:PutObject"
            ],
            "Resource": [
                "arn:aws:s3:::${aws_s3_bucket.bucket_for_Kinesis_failures.bucket}/*",
                "arn:aws:s3:::${aws_s3_bucket.bucket_for_Kinesis_failures.bucket}"
            ],
            "Effect": "Allow"
        }
    ]
}
EOF
}

resource "aws_s3_bucket" "bucket_for_Kinesis_failures" {
  bucket = "${replace(lower(local.identifier),"_", "-")}-kinesisfailures"
}
resource "aws_s3_bucket_acl" "bucket_for_Kinesis_failures" {
  bucket = aws_s3_bucket.bucket_for_Kinesis_failures.id
  acl    = "private"
}
resource "aws_kinesis_firehose_delivery_stream" "kinesis_stream" {
  name        = local.identifier
  destination = "http_endpoint"

  http_endpoint_configuration {
    name = "ToSumo"
    url = sumologic_kinesis_metrics_source.kinesis_source.url
    role_arn   = aws_iam_role.firehose_role.arn
    buffering_interval = 60
    s3_backup_mode = "FailedDataOnly"

    request_configuration {
      content_encoding = "GZIP"
    }

    cloudwatch_logging_options {
      enabled = true
      log_group_name = "/aws/kinesisfirehose/${local.identifier}"
      log_stream_name = "DestinationDelivery"
    }
  }

  s3_configuration {
    role_arn   = aws_iam_role.firehose_role.arn
    bucket_arn = aws_s3_bucket.bucket_for_Kinesis_failures.arn
  }
}


// ------------------------------------ authorizing Sumo to use our AWS accounts

resource "aws_iam_policy" "cloudwatch_ingest" {
  name        = "policy-for-cloudwatch-ingest"
  description = "Managed by Terraform"

  policy = <<EOF
{
  "Version": "2012-10-17",
  "Statement": [
      {
          "Action": [
              "cloudwatch:ListMetrics",
              "cloudwatch:GetMetricStatistics",
              "tag:GetResources"
          ],
          "Effect": "Allow",
          "Resource": "*"
      }
  ]
}
EOF
}

data "aws_iam_policy_document" "sumo_can_use_our_AWS" {
  statement {
    actions = ["sts:AssumeRole"]
    condition {
      test = "StringEquals"
      variable = "sts:ExternalId"
      values = [local.account_id]
    }
    principals {
      identifiers = ["arn:aws:iam::${local.account_id}:root"]
      type = "AWS"
    }
  }
}

resource "aws_iam_role" "cloudwatch_role" {
  name               = "role-for-cloudwatch-ingest"
  assume_role_policy = data.aws_iam_policy_document.sumo_can_use_our_AWS.json
}

resource "aws_iam_role_policy_attachment" "test-attach" {
  role       = aws_iam_role.cloudwatch_role.name
  policy_arn = aws_iam_policy.cloudwatch_ingest.arn
}

```


[1]: https://help.sumologic.com/Send_Data/Sources/03Use_JSON_to_Configure_Sources/JSON_Parameters_for_Hosted_Sources
