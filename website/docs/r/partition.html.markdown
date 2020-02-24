---
layout: "sumologic"
page_title: "SumoLogic: sumologic_partition"
description: |-
  Provides a Sumologic Partition
---

# sumologic_partition
Provides a [Sumologic Partition][1].

## Example Usage
```hcl
resource "sumologic_partition" "examplePartition" {
    name = "terraform_examplePartition"
    routing_expression = "_sourcecategory=*/Terraform"
    analytics_tier = "enhanced"
    is_compliant = false
}
```

## Argument reference

The following arguments are supported:

- `name` - (Required) The name of the partition.
- `routing_expression` - (Required) The query that defines the data to be included in the partition.
- `analytics_tier` - (Required) The Cloud Flex analytics tier for your data; only relevant if your account has basic analytics enabled. Possible values are: `enhanced`, `basic`, `cold`
- `retention_period` - (Optional) The number of days to retain data in the partition, or -1 to use the default value for your account. Only relevant if your account has variable retention enabled.
- `is_compliant` - (Required) Whether the partition is compliant or not. Mark a partition as compliant if it contains data used for compliance or audit purpose. Retention for a compliant partition can only be increased and cannot be reduced after the partition is marked compliant. A partition once marked compliant, cannot be marked non-compliant later.
- `data_forwarding_id` - (Optional) An optional ID of a data forwarding configuration to be used by the partition.

## Attributes reference

The following attributes are exported:

- `id` - Unique identifier for the partition.

[Back to Index][0]

[0]: ../README.md
[1]: https://help.sumologic.com/Manage/Partitions
