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

In addition to the common properties, the following arguments are supported:

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

[1]: https://help.sumologic.com/Send_Data/Sources/03Use_JSON_to_Configure_Sources/JSON_Parameters_for_Hosted_Sources
