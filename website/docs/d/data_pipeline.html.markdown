---
layout: "sumologic"
page_title: "SumoLogic: sumologic_data_pipeline"
description: |-
  Provides a way to retrieve details of a Sumo Logic data pipeline, including its full routing/processing node graph.
---

# sumologic_data_pipeline

Provides a way to retrieve details of a Sumo Logic [data pipeline][1], including its full
node graph (source, router, destination, and processing group nodes, with their processors
and outputs).

[1]: https://help.sumologic.com/docs/send-data/hosted-collectors/data-pipeline/

## Example Usage

```hcl
data "sumologic_data_pipeline" "this" {
  name = "Production Logs Pipeline"
}
```

```hcl
data "sumologic_data_pipeline" "that" {
  id = "00000000HX546"
}
```

A data pipeline can be looked up by either `id` or `name`. Exactly one of those attributes
must be specified.

## Argument Reference

- `id` - (Optional) The id of the data pipeline. Exactly one of `id` or `name` is required.
- `name` - (Optional) The name of the data pipeline. Exactly one of `id` or `name` is required.

## Attributes Reference

The following attributes are exported:

- `id` - The internal ID of the data pipeline.
- `name` - The name of the data pipeline.
- `description` - The description of the data pipeline.
- `pipeline_type` - The type of the data pipeline. Currently only `route_based` is supported.
- `route_expression` - The Sumo Logic query expression that determines which messages this
  pipeline applies to.
- `is_enabled` - Whether the data pipeline is enabled.
- `order` - The pipeline's position in the customer's overall pipeline execution order.
- `state` - The publish state of the data pipeline (e.g. `draft`, `published`,
  `published_with_draft`).
- `version` - The current version of the data pipeline, used for optimistic locking on
  updates.
- `customer_id` - The id of the customer that owns the data pipeline.
- `created_at` - The timestamp of when the data pipeline was created.
- `modified_at` - The timestamp of when the data pipeline was last modified.
- `created_by_user_id` - The id of the user who created the data pipeline.
- `modified_by_user_id` - The id of the user who last modified the data pipeline.
- `node` - The list of nodes making up the pipeline's routing/processing graph, in the order
  returned by the API. Each node exports:
  - `id` - The internal ID of the node.
  - `name` - The name of the node.
  - `node_type` - The type of the node: `source`, `router`, `destination`, or
    `processing_group`.
  - `filter_expression` - For `processing_group` nodes, the expression that further filters
    which messages entering the node get processed.
  - `is_enabled` - Whether the node is enabled.
  - `processor` - For `processing_group` nodes, the ordered list of processors run against
    messages passing through the node. Each processor exports:
    - `id` - The internal ID of the processor.
    - `name` - The name of the processor.
    - `processor_type` - The type of the processor.
    - `order` - The processor's position within the processing group.
    - `is_enabled` - Whether the processor is enabled.
    - `config` - The processor's type-specific configuration, as a raw JSON string.
  - `output` - The list of outgoing edges from this node. Each output exports:
    - `target` - The name of the node this output routes to.
    - `condition` - The Sumo Logic query expression that must match for this output to be
      taken. An output with no condition acts as a catch-all.
    - `order` - The priority of this output relative to the node's other outputs; lower
      values are evaluated first.
