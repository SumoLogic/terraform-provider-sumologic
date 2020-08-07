---
layout: "sumologic"
page_title: "SumoLogic: sumologic_metadata_source"
description: |-
  Provides a Sumologic Metadata (Tag) source. This source allows you to collect tags from EC2 instances running on AWS.
---

# sumologic_polling_source
Provides a Sumologic Metadata (Tag) source. This source allows you to collect tags from EC2 instances running on AWS.

__IMPORTANT:__ The AWS credentials are stored in plain-text in the state. This is a potential security issue.

## Example Usage
```hcl
resource "sumologic_metadata_source" "terraform_metadata" {
  name          = "Metadata source"
  description   = "My description"
  category      = "aws/metadata"
  content_type  = "AwsMetadata"
  scan_interval = 300000
  paused        = false
  collector_id  = "${sumologic_collector.collector.id}"

  authentication {
    type = "AWSRoleBasedAuthentication"
    role_arn = "arn:aws:iam::604066827510:role/cw-role-SumoRole-4AOLS73TGKYI"
  }

  path {
    type = "AwsMetadataPath"
    limit_to_regions = ["us-west-2"]
    limit_to_namespaces = ["AWS/EC2"]
    tag_filters = ["Deploy*,", "!DeployStatus,", "Cluster"]
}

resource "sumologic_collector" "collector" {
  name        = "my-collector"
  description = "Just testing this"
}
```

## Argument reference

In addition to the common properties, the following arguments are supported:

 - `content_type` - (Required) The content-type of the collected data. For Metadata source this is `AwsMetadata`. Details can be found in the [Sumologic documentation for hosted sources][1].
 - `scan_interval` - (Required) Time interval in milliseconds of scans for new data. The default is 300000 and the minimum value is 1000 milliseconds.
 - `paused` - (Required) When set to true, the scanner is paused. To disable, set to false.
 - `authentication` - (Required) Authentication details for AWS access.
     + `type` - (Required) Must be either `S3BucketAuthentication` or `AWSRoleBasedAuthentication`
     + `access_key` - (Required) Your AWS access key if using type `S3BucketAuthentication`
     + `secret_key` - (Required) Your AWS secret key if using type `S3BucketAuthentication`
     + `role_arn` - (Required) Your AWS role ARN if using type `AWSRoleBasedAuthentication`
 - `path` - (Required) The location to scan for new data.
     + `type` - (Required) type of polling source. Only allowed value is `AwsMetadataPath`.
     + `limit_to_regions` - (Optional) List of Amazon regions.
     + `limit_to_namespaces` - List of namespaces. For `AwsMetadataPath` the only valid namespace is `AWS/EC2`. 
     + `tag_filters` - (Optional) Leave this field blank to collect all tags configured for the EC2 instance. To collect a subset of tags, follow the instructions in [Define EC2 tag filters][2]

### See also
  * [Sumologic > Sources > Sources for Hosted Collectors > AWS > AWS Metadata (Tag) Source][3]
  * [Common Source Properties][4]

## Attributes Reference
The following attributes are exported:

- `id` - The internal ID of the source.
- `url` - The HTTP endpoint to use with [SNS to notify Sumo Logic of new files](https://help.sumologic.com/03Send-Data/Sources/02Sources-for-Hosted-Collectors/Amazon-Web-Services/AWS-S3-Source#Set_up_SNS_in_AWS_(Optional)).

## Import
Metadata sources can be imported using the collector and source IDs (`collector/source`), e.g.:

```hcl
terraform import sumologic_metadata_source.test 123/456
```

Metadata sources can be imported using the collector name and source name (`collectorName/sourceName`), e.g.:

```hcl
terraform import sumologic_metadata_source.test my-test-collector/my-test-source
```

[1]: https://help.sumologic.com/Send_Data/Sources/03Use_JSON_to_Configure_Sources/JSON_Parameters_for_Hosted_Sources
[2]:https://help.sumologic.com/03Send-Data/Sources/02Sources-for-Hosted-Collectors/Amazon-Web-Services/AWS-Metadata-(Tag)-Source#Define_EC2_tag_filters
[3]:https://help.sumologic.com/03Send-Data/Sources/02Sources-for-Hosted-Collectors/Amazon-Web-Services/AWS-Metadata-(Tag)-Source
[4]:https://github.com/terraform-providers/terraform-provider-sumologic/tree/master/website#common-source-properties