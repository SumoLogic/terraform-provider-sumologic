---
layout: "sumologic"
page_title: "SumoLogic: sumologic_s3_audit_source"
description: |-
  Provides a AWS S3 Audit Source.
---

# sumologic_polling_source
Provides a [AWS S3 Audit Source][2].

__IMPORTANT:__ The AWS credentials are stored in plain-text in the state. This is a potential security issue.

## Example Usage
```hcl

resource "sumologic_s3_audit_source" "terraform_s3_audit" {
  name          = "Amazon S3 Audit"
  description   = "My description"
  category      = "aws/s3audit"
  content_type  = "AwsS3AuditBucket"
  scan_interval = 300000
  paused        = false
  collector_id  = "${sumologic_collector.collector.id}"

  authentication {
    type = "S3BucketAuthentication"
    access_key = "someKey"
    secret_key = "******"
  }

  path {
    type = "S3BucketPathExpression"
    bucket_name     = "Bucket1"
    path_expression = "*"
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
 - `scan_interval` - (Required) Time interval in milliseconds of scans for new data. The default is 300000 and the minimum value is 1000 milliseconds.
 - `paused` - (Required) When set to true, the scanner is paused. To disable, set to false.
 - `authentication` - (Required) Authentication details for connecting to the S3 bucket.
     + `type` - (Required) Must be either `S3BucketAuthentication` or `AWSRoleBasedAuthentication`.
     + `access_key` - (Required) Your AWS access key if using type `S3BucketAuthentication`.
     + `secret_key` - (Required) Your AWS secret key if using type `S3BucketAuthentication`.
     + `role_arn` - (Required) Your AWS role ARN if using type `AWSRoleBasedAuthentication`.This is not supported for AWS China regions.
     + `region` - (Optional) Your AWS Bucket region.
 - `path` - (Required) The location to scan for new data.
     + `type` - (Required) type of polling source. This has to be `S3BucketPathExpression` for `S3 Audit source`.
     + `bucket_name` - (Required) The name of the bucket. 
     + `path_expression` - (Required) The path to the data.
     + `sns_topic_or_subscription_arn` - (Optional) This is a computed field for SNS topic/subscription ARN. After creating the source via terraform, you can setup [SNS subscription][3] on Sumologic UI and run the terraform apply with the `-refresh-only` flag to update the state with the actual state of the infrastructure. If SNS has been successfully configured and has received a subscription confirmation request `isSuccess` will be true.

### See also
   * [Common Source Properties](https://registry.terraform.io/providers/SumoLogic/sumologic/latest/docs#common-source-properties)

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

[1]: https://help.sumologic.com/Send_Data/Sources/03Use_JSON_to_Configure_Sources/JSON_Parameters_for_Hosted_Sources
[2]: https://help.sumologic.com/03Send-Data/Sources/02Sources-for-Hosted-Collectors/Amazon-Web-Services/Amazon-S3-Audit-Source
[3]: https://help.sumologic.com/03Send-Data/Sources/02Sources-for-Hosted-Collectors/Amazon-Web-Services/AWS_Sources#set-up-sns-in-aws-highly-recommended