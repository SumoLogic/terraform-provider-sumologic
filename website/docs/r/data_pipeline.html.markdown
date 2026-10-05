---
layout: "sumologic"
page_title: "SumoLogic: sumologic_data_pipeline"
description: |-
  Provides a Sumologic Data Pipeline
---

# sumologic_data_pipeline

Provides a Sumo Logic [data pipeline][1]. A data pipeline routes messages matching a query
expression through a graph of nodes - exactly one source node, any number of router nodes,
and a single destination node.

This resource manages the pipeline shell and its routing graph only: `source`, `router`, and
`destination` nodes with their outputs. `processing_group` nodes/processors are not yet
supported - see [`sumologic_data_pipeline` data source][2] if you need to read a pipeline that
has them.

[1]: https://help.sumologic.com/docs/send-data/hosted-collectors/data-pipeline/
[2]: /docs/providers/sumologic/d/data_pipeline.html

## Example Usage

```hcl
resource "sumologic_data_pipeline" "example" {
  name             = "Production Logs Pipeline"
  route_expression = "_sourceCategory=prod/*"

  node {
    name      = "Routing Expression"
    node_type = "source"

    output {
      target = "Sumo Logic"
    }
  }

  node {
    name      = "Sumo Logic"
    node_type = "destination"
  }
}
```

### With a router node fanning out to a single destination

```hcl
resource "sumologic_data_pipeline" "example" {
  name             = "Triaged Logs Pipeline"
  route_expression = "_sourceCategory=prod/*"

  node {
    name      = "Routing Expression"
    node_type = "source"

    output {
      target = "Triage"
    }
  }

  node {
    name      = "Triage"
    node_type = "router"

    output {
      target    = "Sumo Logic"
      condition = "event_category=security"
      order     = 1
    }

    output {
      target = "Sumo Logic" # catch-all - no condition, evaluated last
    }
  }

  node {
    name      = "Sumo Logic"
    node_type = "destination"
  }
}
```

## Argument Reference

The following arguments are supported:

- `name` - (Required) The name of the data pipeline. Must be between 1 and 128 characters.
- `description` - (Optional) The description of the data pipeline. Must be between 0 and 1024
  characters.
- `pipeline_type` - (Optional, ForceNew) The type of the data pipeline. Currently only
  `route_based` is supported. Defaults to `route_based`.
- `route_expression` - (Required) The Sumo Logic query expression that determines which
  messages this pipeline applies to.
- `is_enabled` - (Optional) Whether the data pipeline is enabled. Defaults to `true`.
- `node` - (Required) The list of nodes making up the pipeline's routing graph, in the order
  they're declared. A pipeline must have exactly one `source` node and at most one
  `destination` node. See [node](#node) below.

### node

- `name` - (Required) The name of the node. The source node's name must be the literal string
  `Routing Expression`.
- `node_type` - (Required) The type of the node. One of `source`, `router`, or `destination`.
- `output` - (Optional) The list of outgoing edges from this node. See [output](#output)
  below.

### output

- `target` - (Required) The name of the node this output routes to.
- `condition` - (Optional) The Sumo Logic query expression that must match for this output to
  be taken. An output with no condition acts as a catch-all, and should be declared last.
- `order` - (Optional) The priority of this output relative to the node's other outputs; lower
  values are evaluated first.

## Attributes Reference

The following attributes are exported:

- `id` - The internal ID of the data pipeline.
- `order` - The pipeline's position in the customer's overall pipeline execution order.
- `state` - The publish state of the data pipeline (e.g. `draft`, `published`,
  `published_with_draft`).
- `version` - The current version of the data pipeline, used for optimistic locking on
  updates.
- `node.id` - The internal ID of the node. Not a stable identity across updates: the API
  assigns every node a fresh ID whenever an update touches the node list at all, even for a
  node whose ID and content are otherwise unchanged. Only stable at rest, between updates.

## Import

Data pipelines can be imported using the pipeline ID.

```hcl
terraform import sumologic_data_pipeline.example 00000000HX546
```
