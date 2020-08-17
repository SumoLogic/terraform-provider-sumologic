---
layout: "sumologic"
page_title: "SumoLogic: sumologic_aws_inventory_source"
description: |-
  Provides a Sumologic AWS Inventory source.
---

# sumologic_aws_inventory_source
Provides a Sumologic AWS Inventory source to collect AWS resource inventory data.

__IMPORTANT:__ The AWS credentials are stored in plain-text in the state. This is a potential security issue.

## Example Usage
```hcl
resource "sumologic_aws_inventory_source" "terraform_aws_inventory_source" {
  name          = "AWS Inventory"
  description   = "My description"
  category      = "aws/terraform_aws_inventory"
  content_type  = "AwsInventory"
  scan_interval = 300000
  paused        = false
  collector_id  = "${sumologic_collector.collector.id}"

  authentication {
    type = "AWSRoleBasedAuthentication"
    role_arn = "arn:aws:iam::01234567890:role/sumo-role"
  }

  path {
    type = "AwsInventoryPath"
    limit_to_regions = ["us-west-2"]
    limit_to_namespaces = ["AWS/RDS","AWS/EC2"]
  }
}

resource "sumologic_collector" "collector" {
  name        = "my-collector"
  description = "Just testing this"
}
```

## Argument reference

In addition to the common properties, the following arguments are supported:

 - `content_type` - (Required) The content-type of the collected data. This has to be `AwsInventoryPath` for AWS Inventory source.
 - `scan_interval` - (Required) Time interval in milliseconds of scans for new data. The minimum value is 1000 milliseconds. Currently this value is not respected.
 - `paused` - (Required) When set to true, the scanner is paused. To disable, set to false.
 - `authentication` - (Required) Authentication details to access AWS `Describe*` APIs.
     + `type` - (Required) Must be either `S3BucketAuthentication` or `AWSRoleBasedAuthentication`
     + `access_key` - (Required) Your AWS access key if using type `S3BucketAuthentication`
     + `secret_key` - (Required) Your AWS secret key if using type `S3BucketAuthentication`
     + `role_arn` - (Required) Your AWS role ARN if using type `AWSRoleBasedAuthentication`
 - `path` - (Required) The location to scan for new data.
     + `type` - (Required) type of polling source. This has to be `AwsInventoryPath` for AWS Inventory source.
     + `limit_to_regions` - (Optional) List of Amazon regions. 
     + `limit_to_namespaces` - (Optional) List of namespaces. By default all namespaces are selected. You can also choose a subset from
        + AWS/EC2
        + AWS/AutoScaling
        + AWS/EBS
        + AWS/ELB
        + AWS/ApplicationELB
        + AWS/NetworkELB
        + AWS/Lambda
        + AWS/RDS
        + AWS/Dynamodb
        + AWS/ECS 
        + AWS/Elasticache
        + AWS/Redshift
        + AWS/Kinesis

### See also
  * [Common Source Properties](https://github.com/terraform-providers/terraform-provider-sumologic/tree/master/website#common-source-properties)

## Attributes Reference
The following attributes are exported:

- `id` - The internal ID of the source.

## Import
AWS Inventory sources can be imported using the collector and source IDs (`collector/source`), e.g.:

```hcl
terraform import sumologic_aws_inventory_source.test 123/456
```

AWS Inventory sources can be imported using the collector name and source name (`collectorName/sourceName`), e.g.:

```hcl
terraform import sumologic_aws_inventory_source.test my-test-collector/my-test-source
```