---
layout: "sumologic"
page_title: "SumoLogic: sumologic_http_source"
description: |-
  Provides a Sumologic HTTP source
---

# sumologic_http_source
Provides a [Sumologic HTTP source][1].

__IMPORTANT:__ The endpoint is stored in plain-text in the state. This is a potential security issue.

## Example Usage
```hcl
resource "sumologic_http_source" "http_source" {
  name                = "HTTP"
  description         = "My description"
  category            = "my/source/category"
  collector_id        = "${sumologic_collector.collector.id}"
}

resource "sumologic_collector" "collector" {
  name        = "my-collector"
  description = "Just testing this"
}
```

## Argument reference
In addition to the common properties, the following arguments are supported:
- `message_per_request` - (Optional) When set to `true`, will create one log message per HTTP request.

## Attributes reference
The following attributes are exported:
- `id` - The internal ID of the source.
- `url` - The HTTP endpoint to use for sending data to this source.

## Import
HTTP sources can be imported using the collector and source IDs (`collector/source`), e.g.:

```hcl
terraform import sumologic_http_source.test 123/456
```

+HTTP sources can be imported using the collector name and source name (`collectorName/sourceName`), e.g.:

```hcl
terraform import sumologic_http_source.test my-test-collector/my-test-source
```

[Back to Index][0]

[0]: ../README.md
[1]: https://help.sumologic.com/Send_Data/Sources/02Sources_for_Hosted_Collectors/HTTP_Source
