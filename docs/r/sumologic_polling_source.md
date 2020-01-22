# sumologic_polling_source
Provides a Sumologic Polling source. This source is used to import data from  AWS S3 buckets.

__IMPORTANT:__ The AWS credentials are stored in plain-text in the state. This is a potential security issue.

## Example Usage
```hcl
locals {
  filters = [{
    name        = "Exclude Comments"
    filter_type = "Exclude"
    regexp      = "#.*"
  }]
}

resource "sumologic_polling_source" "s3_audit" {
  name          = "Amazon S3 Audit"
  description   = "My description"
  category      = "aws/s3audit"
  content_type  = "AwsS3AuditBucket"
  scan_interval = 300000
  paused        = false
  collector_id  = "${sumologic_collector.collector.id}"
  filters       = "${local.filters}"

  authentication {
    access_key = "someKey"
    secret_key = "******"
  }

  path {
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
     + `type` - (Required) Must be either `S3BucketAuthentication` or `AWSRoleBasedAuthentication`
     + `access_key` - (Required) Your AWS access key if using type `S3BucketAuthentication`
     + `secret_key` - (Required) Your AWS secret key if using type `S3BucketAuthentication`
     + `role_arn` - (Required) Your AWS role ARN if using type `AWSRoleBasedAuthentication`
 - `path` - (Required) The location to scan for new data.
     + `bucket_name` - (Required) The name of the bucket.
     + `path_expression` - (Required) The path to the data.

## Attributes reference
- `id` - The internal ID of the source.
- `url` - The HTTP endpoint to use with [SNS to notify Sumo Logic of new files](https://help.sumologic.com/03Send-Data/Sources/02Sources-for-Hosted-Collectors/Amazon-Web-Services/AWS-S3-Source#Set_up_SNS_in_AWS_(Optional)).

[Back to Index][0]

[0]: ../README.md
[1]: https://help.sumologic.com/Send_Data/Sources/03Use_JSON_to_Configure_Sources/JSON_Parameters_for_Hosted_Sources
