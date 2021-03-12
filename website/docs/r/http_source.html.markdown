---
layout: "sumologic"
page_title: "SumoLogic: sumologic_http_source"
description: |-
  Provides a Sumologic HTTP source
---

# sumologic_http_source
Provides a [Sumologic HTTP source][1] or [Sumologic HTTP Traces source][2].

__IMPORTANT:__ The endpoint is stored in plain-text in the state. This is a potential security issue.

## Example Usage
```hcl
resource "sumologic_http_source" "http_source" {
  name                = "HTTP"
  description         = "My description"
  category            = "my/source/category"
  collector_id        = "${sumologic_collector.collector.id}"
  filters       {
      name = "Test Exclude Debug"
      filter_type = "Exclude"
      regexp = ".*DEBUG.*"
  }
}

resource "sumologic_http_traces_source" "http_traces_source" {
  name                = "HTTP Traces"
  description         = "My description"
  category            = "my/source/category"
  collector_id        = "${sumologic_collector.collector.id}"
  content_type        = "Zipkin"
}

resource "sumologic_collector" "collector" {
  name        = "my-collector"
  description = "Just testing this"
}
```

## Argument reference

In addition to the common properties, the following arguments are supported:

- `message_per_request` - (Optional) When set to `true`, will create one log message per HTTP request.
- `content_type`        - (Optional) When configuring a HTTP Traces Source, set this property to `Zipkin`. This should only be used when creating a Traces source.

### See also
  * [Common Source Properties](https://github.com/SumoLogic/terraform-provider-sumologic/tree/master/website#common-source-properties)

## Attributes Reference
The following attributes are exported:

- `id` - The internal ID of the source.
- `url` - The HTTP endpoint to use for sending data to this source.

## Import
HTTP sources can be imported using the collector and source IDs (`collector/source`), e.g.:

```hcl
terraform import sumologic_http_source.test 123/456
```

HTTP sources can be imported using the collector name and source name (`collectorName/sourceName`), e.g.:

```hcl
terraform import sumologic_http_source.test my-test-collector/my-test-source
```

[1]: https://help.sumologic.com/Send_Data/Sources/02Sources_for_Hosted_Collectors/HTTP_Source
[2]: https://help.sumologic.com/Traces/HTTP_Traces_Source
