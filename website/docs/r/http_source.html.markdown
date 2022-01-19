---
layout: "sumologic"
page_title: "SumoLogic: sumologic_http_source"
description: |-
  Provides a Sumologic HTTP source
---

# sumologic_http_source
Provides a [Sumologic HTTP source][1], [Sumologic HTTP Traces source][2], [Sumologic Kinesis Log source][3] and [Sumologic HTTP_OTLP_source][4]. To start using Traces contact your Sumo account representative to activate.

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

resource "sumologic_http_source" "http_traces_source" {
  name                = "HTTP Traces"
  description         = "My description"
  category            = "my/source/category"
  collector_id        = "${sumologic_collector.collector.id}"
  content_type        = "Zipkin"
}

resource "sumologic_http_source" "kinesisLog" {
  name = "demo-name"
  description = "demo-desc"
  category = "demo-category"
  content_type = "KinesisLog"
  collector_id = "${sumologic_collector.test.id}"
}

resource "sumologic_http_source" "http_otlp_source" {
  name = "HTTP OTLP"
  description = "My description"
  category = "my/source/category"
  content_type = "Otlp"
  collector_id = "${sumologic_collector.test.id}"
}

resource "sumologic_collector" "collector" {
  name        = "my-collector"
  description = "Just testing this"
}
```

## Argument reference

In addition to the [Common Source Properties](https://registry.terraform.io/providers/SumoLogic/sumologic/latest/docs#common-source-properties), the following arguments are supported:

- `message_per_request` - (Optional) When set to `true`, will create one log message per HTTP request.
- `content_type`        - (Optional) When configuring a HTTP Traces Source, set this property to `Zipkin`. When configuring a Kinesis Logs Source, set this property to `KinesisLog`. When configuring a HTTP OTLP Source, set this property to `Otlp`. This should only be used when creating a Traces, Kinesis Log or HTTP OTLP source.

### See also
  * [Common Source Properties](https://registry.terraform.io/providers/SumoLogic/sumologic/latest/docs#common-source-properties)

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
[3]: https://help.sumologic.com/03Send-Data/Sources/02Sources-for-Hosted-Collectors/Amazon-Web-Services/AWS_Kinesis_Firehose_for_Logs_Source
[4]: https://help.sumologic.com/03Send-Data/Sources/02Sources-for-Hosted-Collectors/OTLP_HTTP_Source
