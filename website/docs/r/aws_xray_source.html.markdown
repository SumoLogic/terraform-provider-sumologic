---
layout: "sumologic"
page_title: "SumoLogic: sumologic_aws_xray_source"
description: |-
  Provides a Sumologic AWS XRay source.
---

# sumologic_aws_xray_source
Provides a Sumologic AWS XRay source to collect metrics derived from XRay traces.

__IMPORTANT:__ The AWS credentials are stored in plain-text in the state. This is a potential security issue.

## Example Usage
```hcl
resource "sumologic_aws_xray_source" "terraform_aws_xray_source" {
  name          = "AWS XRay Metrics"
  description   = "My description"
  category      = "aws/terraform_xray"
  content_type  = "AwsXRay"
  scan_interval = 300000
  paused        = false
  collector_id  = "${sumologic_collector.collector.id}"

  authentication {
    type = "AWSRoleBasedAuthentication"
    role_arn = "arn:aws:iam::01234567890:role/sumo-role"
  }

  path {
    type = "AwsXRayPath"
    limit_to_regions = ["us-west-2"]
  }
}

resource "sumologic_collector" "collector" {
  name        = "my-collector"
  description = "Just testing this"
}
```

## Argument reference

In addition to the common properties, the following arguments are supported:

 - `content_type` - (Required) The content-type of the collected data. This has to be `AwsXRay` for AWS XRay source.
 - `scan_interval` - (Required) Time interval in milliseconds of scans for new data. The minimum value is 1000 milliseconds. Currently this value is not respected, and collection happens at a default interval of 1 minute.
 - `paused` - (Required) When set to true, the scanner is paused. To disable, set to false.
 - `authentication` - (Required) Authentication details for making `xray:Get*` calls.
     + `type` - (Required) Must be either `S3BucketAuthentication` or `AWSRoleBasedAuthentication`
     + `access_key` - (Required) Your AWS access key if using type `S3BucketAuthentication`
     + `secret_key` - (Required) Your AWS secret key if using type `S3BucketAuthentication`
     + `role_arn` - (Required) Your AWS role ARN if using type `AWSRoleBasedAuthentication`
 - `path` - (Required) The location to scan for new data.
     + `type` - (Required) type of polling source. This has to be `AwsXRayPath` for AWS XRay source.
     + `limit_to_regions` - (Optional) List of Amazon regions. 

### See also
  * [Common Source Properties](https://github.com/SumoLogic/terraform-provider-sumologic/tree/master/website#common-source-properties)

## Attributes Reference
The following attributes are exported:

- `id` - The internal ID of the source.

## Import
AWS XRay sources can be imported using the collector and source IDs (`collector/source`), e.g.:

```hcl
terraform import sumologic_aws_xray_source.test 123/456
```

AWS XRay sources can be imported using the collector name and source name (`collectorName/sourceName`), e.g.:

```hcl
terraform import sumologic_aws_xray_source.test my-test-collector/my-test-source
```