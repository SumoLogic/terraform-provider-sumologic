---
layout: "sumologic"
page_title: "SumoLogic: sumologic_data_pipeline_test"
description: |-
  Runs sample data through a data pipeline definition and returns how each message is transformed at every node and processor.
---

# sumologic_data_pipeline_test

Runs sample data through a data pipeline definition and returns the per-message, per-node and
per-processor transformation results. Use it to check how data flows through a pipeline before
creating or changing it with the `sumologic_data_pipeline` resource.

The pipeline is defined inline with the same arguments as the `sumologic_data_pipeline` resource.
Nothing is created or changed in Sumo Logic. The pipeline configuration is validated before it is
executed, and the route expression is not evaluated against the sample data.

## Example Usage

```hcl
data "sumologic_data_pipeline_test" "this" {
  name             = "nginx-pipeline"
  route_expression = "_sourceCategory=prod/nginx"
  sample_data      = "192.168.1.1 GET /api/health 200\n192.168.1.2 GET /api/users 500"

  metadata_fields = {
    _sourceCategory = "prod/nginx"
  }

  node {
    name      = "Routing Expression"
    node_type = "source"
    output {
      target = "Parse Nginx"
    }
  }

  node {
    name      = "Parse Nginx"
    node_type = "processing_group"
    processor {
      name           = "extract-client-ip"
      processor_type = "REGEX_PARSE"
      config = jsonencode({
        fieldToParseFrom = "_raw"
        regexPattern     = "(?<clientip>[\\d.]+) "
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

output "parsed" {
  value = [for m in data.sumologic_data_pipeline_test.this.message : m.output]
}
```

## Argument Reference

The following arguments are supported:

- `name` - (Required) The name of the pipeline definition being tested.
- `description` - (Optional) The description of the pipeline definition.
- `pipeline_type` - (Optional) The type of the pipeline. Currently only `route_based` is supported.
  Defaults to `route_based`.
- `route_expression` - (Required) The Sumo Logic query expression that determines which messages the
  pipeline applies to.
- `is_enabled` - (Optional) Whether the pipeline definition is enabled. Defaults to `true`.
- `node` - (Required) The list of nodes making up the pipeline's routing/processing graph. It accepts
  the same `node`, `output` and `processor` arguments as the `sumologic_data_pipeline` resource.
- `sample_data` - (Required) The raw sample data to run through the pipeline.
- `message_separator` - (Optional) How `sample_data` is split into individual messages. One of
  `\n` (newline, the default), `\r\n`, `---`, `\|`, `(?<=\})(?=\{)` (JSON boundary),
  `(?=\d{4}-\d{2}-\d{2}[T ])` (ISO timestamp), `(?=\w{3}\s+\d+\s+\d{2}:\d{2}:\d{2})`
  (syslog timestamp), or `(empty)` to treat the whole input as a single message. These are
  literal tokens made of a backslash and a letter, not escape sequences, so each backslash must be
  doubled in HCL, for example `message_separator = "\\n"`. An actual newline character, an empty
  string, and any other value are rejected.
- `metadata_fields` - (Optional) A map of metadata fields, such as `_sourceCategory`, attached to
  each message before pipeline processing. These simulate metadata that would normally be extracted
  in production.

## Attributes Reference

In addition to the arguments above, the following attributes are exported:

- `id` - A hash of the request, stable for identical inputs.
- `message` - One entry per message in `sample_data`, in input order. Each message exports:
  - `message_index` - The zero-based index of the message in the input.
  - `input` - The original message after splitting.
  - `output` - The final message after all pipeline processing.
  - `truncated` - Whether the output was truncated for exceeding the size limit.
  - `node_result` - The `processing_group` and `router` nodes the message visited, in order. Each exports:
    - `node_name` - The name of the node.
    - `node_type` - The type of the node, `processing_group` or `router`.
    - `matched` - For processing groups, whether the filter expression matched the message. Always
      `true` for routers.
    - `input` - The message before the node ran.
    - `output` - The message after the node ran.
    - `selected_route_condition` - For routers, the condition of the route taken. Empty for the
      default route.
    - `selected_route_target` - For routers, the name of the node the message was routed to.
    - `processor_result` - For processing groups, the result of each processor. Each exports:
      - `processor_name` - The name of the processor.
      - `processor_type` - The type of the processor.
      - `input` - The message before the processor ran.
      - `output` - The message after the processor ran.
      - `message_dropped` - Whether the processor dropped the message.
  - `edge_result` - The connections the message traversed, in order. Each exports:
    - `source_node` - The name of the node the edge leaves.
    - `target_node` - The name of the node the edge enters.
    - `condition_expression` - The routing condition on the edge, if any.
    - `edge_order` - The priority order of the edge.
