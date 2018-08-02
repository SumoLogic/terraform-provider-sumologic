# sumologic_polling_source
Provides a Sumologic Polling source. This source is used to import data from  AWS S3 buckets.

__IMPORTANT:__ The AWS credentials are stored in plain-text in the state. This is a potential security issue.

## Example Usage
```hcl
resource "sumologic_polling_source" "s3_audit" {
  name          = "Amazon S3 Audit"
  description   = "My description"
  category      = "aws/s3audit"
  content_type  = "AwsS3AuditBucket"
  scan_interval = 1
  paused        = false
  collector_id  = "${sumologic_collector.collector.id}"

  authentication {
    access_key = "AKIAIOSFODNN7EXAMPLE"
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
 - `name` - (Required) The name of the source. This is required, and has to be unique in the scope of the collector. Changing this will force recreation the source.
 - `description` - (Optional) Description of the source.
 - `collector_id` - (Required) The ID of the collector to attach this source to.
 - `category` - (Required) The source category this source logs to.
 - `content_type` - (Required) The content-type of the collected data. Details can be found in the [Sumologic documentation for hosted sources][2].
 - `scan_interval` - (Required) Time interval of scans for new data.
 - `paused` - (Required) When set to true, the scanner is paused. To disable, set to false.
 - `authentication` - (Required) Authentication details for connecting to the S3 bucket.
     + `access_key` - (Required) Your AWS access key
     + `secret_key` - (Required) Your AWS secret key
 - `path` - (Required) The location to scan for new data.
     + `bucket_name` - (Required) The name of the bucket.
     + `path_expression` - (Required) The path to the data.

## Attributes reference
- `id` - The internal ID of the source.

[Back to Index][0]

[0]: ../README.md
[1]: https://help.sumologic.com/Send_Data/Sources/03Use_JSON_to_Configure_Sources/JSON_Parameters_for_Hosted_Sources
