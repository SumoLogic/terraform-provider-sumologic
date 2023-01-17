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
```hcl

  resource "sumologic_kinesis_log_source" "kinesis_log_access_key" {
    name         = "Kinesis Log"
    description  = "Description for Kinesis Log Source"
    category     = "prod/kinesis/log"
    content_type = "KinesisLog"
    collector_id = "${sumologic_collector.collector.id}"
    authentication {
      type       = "S3BucketAuthentication"
      access_key = "someKey"
      secret_key = "******"
    }

    path {
      type            = "KinesisLogPath"
      bucket_name     = "testBucket"
      path_expression = "http-endpoint-failed/*"
      scan_interval   = 30000
    }
  }

  resource "sumologic_kinesis_log_source" "kinesis_log_role_arn" {
    name         = "Kinesis Log"
    description  = "Description for Kinesis Log Source"
    category     = "prod/kinesis/log"
    content_type = "KinesisLog"
    collector_id = "${sumologic_collector.collector.id}"

    authentication {
      type     = "AWSRoleBasedAuthentication"
      role_arn = "arn:aws:iam::604066827510:role/cw-role-SumoRole-4AOLS73TGKYI"
    }

    path {
      type            = "KinesisLogPath"
      bucket_name     = "testBucket"
      path_expression = "http-endpoint-failed/*"
      scan_interval   = 30000
    }
  }

  resource "sumologic_collector" "collector" {
    name        = "my-collector"
    description = "Just testing this"
  }

```

## Argument reference

In addition to the common properties, the following arguments are supported:

 - `content_type` - (Required) The content-type of the collected data. Details can be found in the [Sumologic documentation for hosted sources][1].
 - `authentication` - (Optional) Authentication details for connecting to the S3 bucket.
     + `type` - (Required) Must be either `S3BucketAuthentication` or `AWSRoleBasedAuthentication` or `NoAuthentication`
     + `access_key` - (Required) Your AWS access key if using type `S3BucketAuthentication`
     + `secret_key` - (Required) Your AWS secret key if using type `S3BucketAuthentication`
     + `role_arn` - (Required) Your AWS role ARN if using type `AWSRoleBasedAuthentication`
 - `path` - (Optional) The location of S3 bucket for failed Kinesis log data.
     + `type` - (Required) Must be either `KinesisLogPath` or `NoPathExpression`
     + `bucket_name` - (Optional) The name of the bucket. This is needed if using type `KinesisLogPath`. 
     + `path_expression` - (Optional) The path to the data. This is needed if using type `KinesisLogPath`. For Kinesis log source, it must include `http-endpoint-failed/`.
     + `scan_interval` - (Optional) The Time interval in milliseconds of scans for new data. The default is 300000 and the minimum value is 1000 milliseconds.

### See also
   * [Common Source Properties](https://registry.terraform.io/providers/SumoLogic/sumologic/latest/docs#common-source-properties)

## Attributes Reference
The following attributes are exported:

- `id` - The internal ID of the source.
- `url` - The HTTP endpoint to be used while creating Kinesis Firehose on AWS.

## Import
Kinesis Log sources can be imported using the collector and source IDs (`collector/source`), e.g.:

```hcl
terraform import sumologic_kinesis_log_source.test 123/456
```

HTTP sources can be imported using the collector name and source name (`collectorName/sourceName`), e.g.:

```hcl
terraform import sumologic_kinesis_log_source.test my-test-collector/my-test-source
```

[1]: https://help.sumologic.com/Send_Data/Sources/03Use_JSON_to_Configure_Sources/JSON_Parameters_for_Hosted_Sources
[2]: https://help.sumologic.com/03Send-Data/Sources/02Sources-for-Hosted-Collectors/Amazon-Web-Services/AWS_Kinesis_Firehose_for_Logs_Source
