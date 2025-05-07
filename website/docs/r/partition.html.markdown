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
  analytics_tier = "continuous"
  is_compliant = false
  is_included_in_default_search = true
  lifecycle {
    prevent_destroy = true
  }
}
```

## Argument reference

The following arguments are supported:

- `name` - (Required, Forces new resource) The name of the partition.
- `routing_expression` - (Required) The query that defines the data to be included in the partition.
- `analytics_tier` - (Optional) The cloud flex analytics tier for your data; only relevant if your account has basic analytics enabled. If no value is supplied, partition will be created in continuous tier. Other possible values are : "frequent" and "infrequent". For flex partition, you can leave it empty or send flex.
This field is set during partition creation and cannot be modified later by the user. However, Terraform can detect and reconcile backend-driven changes, such as those resulting from a Flex migration.
- `retention_period` - (Optional) The number of days to retain data in the partition, or -1 to use the default value for your account. Only relevant if your account has variable retention enabled.
- `is_compliant` - (Optional) Whether the partition is compliant or not. Mark a partition as compliant if it contains data used for compliance or audit purpose. Retention for a compliant partition can only be increased and cannot be reduced after the partition is marked compliant. A partition once marked compliant, cannot be marked non-compliant later.
- `reduce_retention_period_immediately` - (Optional) This is required on update if the newly specified retention period is less than the existing retention period. In such a situation, a value of true says that data between the existing retention period and the new retention period should be deleted immediately; if false, such data will be deleted after seven days. This property is optional and ignored if the specified retentionPeriod is greater than or equal to the current retention period.
- `is_included_in_default_search` - Indicates whether the partition is included in the default search scope. When executing a query such as "error | count," certain partitions are automatically part of the search scope. However, for specific partitions, the user must explicitly mention the partition using the _index term, as in "_index=webApp error | count". This property governs the default inclusion of the partition in the search scope. Configuring this property is exclusively permitted for flex partitions.

## Attributes reference

The following attributes are exported:

- `id` - Unique identifier for the partition.
## Import
Partitions can can be imported using the id. The list of partitions and their ids can be obtained using the Sumologic [partions api][2].

```hcl
terraform import sumologic_partition.partition 1234567890
```

[1]: https://help.sumologic.com/docs/manage/partitions/data-tiers/
[2]: https://api.sumologic.com/docs/#operation/listPartitions

