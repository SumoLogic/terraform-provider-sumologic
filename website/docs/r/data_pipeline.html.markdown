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

This resource manages the pipeline's full node graph: `source`, `router`, and `destination`
nodes with their outputs, and `processing_group` nodes with their ordered list of processors.

[1]: https://help.sumologic.com/docs/send-data/hosted-collectors/data-pipeline/

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

### With a processing_group node

```hcl
resource "sumologic_data_pipeline" "example" {
  name             = "Normalized Logs Pipeline"
  route_expression = "_sourceCategory=prod/*"

  node {
    name      = "Routing Expression"
    node_type = "source"

    output {
      target = "Log Normalization"
    }
  }

  node {
    name              = "Log Normalization"
    node_type         = "processing_group"
    filter_expression = "_sourceCategory=prod/app/*"

    processor {
      name           = "Extract Client IP"
      processor_type = "REGEX_PARSE"
      config = jsonencode({
        fieldToParseFrom = "_raw"
        regexPattern     = "(?<clientip>[\\d.]+)"
      })
    }

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
- `node_type` - (Required) The type of the node. One of `source`, `router`, `destination`, or
  `processing_group`.
- `filter_expression` - (Optional) For `processing_group` nodes, a Sumo Logic query expression
  that scopes which messages this group processes. Not applicable to other node types.
- `is_enabled` - (Optional) Whether this node is active. Defaults to `true`. Only meaningful for
  `processing_group` nodes.
- `output` - (Optional) The list of outgoing edges from this node. See [output](#output)
  below.
- `processor` - (Optional) For `processing_group` nodes, the ordered list of processors applied
  to matching messages. See [processor](#processor) below.

### output

- `target` - (Required) The name of the node this output routes to.
- `condition` - (Optional) The Sumo Logic query expression that must match for this output to
  be taken. An output with no condition acts as a catch-all, and should be declared last.
- `order` - (Optional) The priority of this output relative to the node's other outputs; lower
  values are evaluated first.

### processor

- `name` - (Required) The display name of this processor.
- `processor_type` - (Required) The processor type identifier (e.g. `REGEX_PARSE`).
  This provider does not validate processor types or their `config` contents against a fixed
  list - Sumo Logic's processor catalog is still growing, so any type the API accepts is
  accepted here too; invalid combinations surface as an error from Sumo Logic at publish time.
- `is_enabled` - (Optional) Whether this processor is active. Defaults to `true`.
- `config` - (Required) A JSON string holding this processor's type-specific configuration.
  The shape of this JSON is entirely determined by `processor_type` and is not validated by
  this provider beyond being well-formed JSON.

A processor's execution order is its position in the `processor` list - there is no separate
`order` argument.

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
- `node.processor.id` - The internal ID of the processor. Carries the same caveat as `node.id`:
  not a stable identity across updates, only stable at rest.

## Import

Data pipelines can be imported using the pipeline ID.

```hcl
terraform import sumologic_data_pipeline.example 00000000HX546
```
