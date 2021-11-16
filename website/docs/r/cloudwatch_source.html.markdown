---
layout: "sumologic"
page_title: "SumoLogic: sumologic_cloudwatch_source"
description: |-
  Provides a Sumologic CloudWatch source.
---

# sumologic_cloudwatch_source
Provides a [Sumologic CloudWatch source][2].

__IMPORTANT:__ The AWS credentials are stored in plain-text in the state. This is a potential security issue.

## Example Usage
```hcl
locals {
  filters = [{
    name        = "Exclude Comments"
    filter_type = "Exclude"
    regexp      = "#.*"
  }] 
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

resource "sumologic_cloudwatch_source" "terraform_cloudwatch_source" {
  name          = "CloudWatch Metrics"
  description   = "My description"
  category      = "aws/terraform_cw"
  content_type  = "AwsCloudWatch"
  scan_interval = 300000
  paused        = false
  collector_id  = "${sumologic_collector.collector.id}"

  authentication {
    type = "AWSRoleBasedAuthentication"
    role_arn = "arn:aws:iam::01234567890:role/sumo-role"
  }

  path {
    type = "CloudWatchPath"
    limit_to_regions = ["us-west-2"]
    limit_to_namespaces = ["AWS/Route53","AWS/S3","customNamespace"]
  
    dynamic "tag_filters" {
    for_each = local.tagfilters
    content {
      type = tag_filters.value.type
      namespace = tag_filters.value.namespace
      tags = tag_filters.value.tags
    }
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
 - `scan_interval` - (Required) Time interval in milliseconds of scans for new data. The default is 300000 and the minimum value is 1000 milliseconds.
 - `paused` - (Required) When set to true, the scanner is paused. To disable, set to false.
 - `authentication` - (Required) Authentication details for connecting to the S3 bucket.
     + `type` - (Required) Must be either `S3BucketAuthentication` or `AWSRoleBasedAuthentication`
     + `access_key` - (Required) Your AWS access key if using type `S3BucketAuthentication`
     + `secret_key` - (Required) Your AWS secret key if using type `S3BucketAuthentication`
     + `role_arn` - (Required) Your AWS role ARN if using type `AWSRoleBasedAuthentication`. This is not supported for AWS China regions.
     + `region` - (Optional) Your AWS Bucket region.
 - `path` - (Required) The location to scan for new data.
     + `type` - (Required) type of polling source. This has to be `CloudWatchPath` for CloudWatch source.
     + `limit_to_regions` - (Optional) List of Amazon regions. 
     + `limit_to_namespaces` - (Optional) List of namespaces. By default all namespaces are selected. Details can be found [here](https://help.sumologic.com/03Send-Data/Sources/02Sources-for-Hosted-Collectors/Amazon-Web-Services/Amazon-CloudWatch-Source-for-Metrics#aws%C2%A0tag-filtering-namespace-support). You can also  specify custom namespace.
     + `tag_filters` - (Optional) Tag filters allow you to filter the CloudWatch metrics you collect by the AWS tags you have assigned to your AWS resources. You can define tag filters for each supported namespace. If you do not define any tag filters, all metrics will be collected for the regions and namespaces you configured for the source above. More info on tag filters can be found [here](https://help.sumologic.com/03Send-Data/Sources/02Sources-for-Hosted-Collectors/Amazon-Web-Services/Amazon-CloudWatch-Source-for-Metrics#about-aws-tag-filtering)
          + `type` - This value has to be set to `TagFilters`
          + `namespace` - Namespace for which you want to define the tag filters. Use  value as `All` to apply the tag filter for all namespaces.
          + `tags` - List of key-value pairs of tag filters. Eg: `["k3=v3"]`

### See also
   * [Common Source Properties](https://registry.terraform.io/providers/SumoLogic/sumologic/latest/docs#common-source-properties)

## Attributes Reference
The following attributes are exported:

- `id` - The internal ID of the source.
- `url` - The HTTP endpoint to use with [SNS to notify Sumo Logic of new files](https://help.sumologic.com/03Send-Data/Sources/02Sources-for-Hosted-Collectors/Amazon-Web-Services/AWS-S3-Source#Set_up_SNS_in_AWS_(Optional)).

## Import
CloudWatch sources can be imported using the collector and source IDs (`collector/source`), e.g.:

```hcl
terraform import sumologic_cloudwatch_source.test 123/456
```

CloudWatch sources can be imported using the collector name and source name (`collectorName/sourceName`), e.g.:

```hcl
terraform import sumologic_cloudwatch_source.test my-test-collector/my-test-source
```

[1]: https://help.sumologic.com/Send_Data/Sources/03Use_JSON_to_Configure_Sources/JSON_Parameters_for_Hosted_Sources
[2]: https://help.sumologic.com/03Send-Data/Sources/02Sources-for-Hosted-Collectors/Amazon-Web-Services/Amazon-CloudWatch-Source-for-Metrics
