---
layout: "sumologic"
page_title: "SumoLogic: sumologic_partition"
description: |-
  Provides a way to retrieve all Sumo Logic partitions.
---

# sumologic_partition

Provides a way to retrieve all [Sumologic Partitions][1].

## Example Usage

```hcl
data "sumologic_partitions" "partitions" {}
```

## Attributes reference

The following attributes are exported:

- `name` - The name of the partition.
- `routing_expression` - The query that defines the data to be included in the partition.
- `analytics_tier` - The Data Tier where the data in the partition will reside. Possible values are: `continuous`, `frequent`, `infrequent`
- `retention_period` - The number of days to retain data in the partition.
- `is_compliant` - Whether the partition is used for compliance or audit purposes.
- `is_included_in_default_search` - Whether the partition is included in the default search scope.
- `total_bytes` - The size of the data in the partition in bytes.
- `is_active` - Whether the partition is currently active or decommissioned.
- `index_type` - The type of partition index. Possible values are: `DefaultIndex`, `AuditIndex`or `Partition`
- `data_forwarding_id` - The ID of the data forwarding configuration to be used by the partition.

[1]: https://help.sumologic.com/docs/manage/partitions/data-tiers/
