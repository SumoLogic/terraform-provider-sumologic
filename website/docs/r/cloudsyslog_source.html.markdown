---
layout: "sumologic"
page_title: "SumoLogic: sumologic_cloudsyslog_source"
description: |-
  Provides a Sumo Logic Cloud Syslog source.
---

# sumologic_cloudsyslog_source

Provides a [Sumo Logic Cloud Syslog source][1].

__IMPORTANT:__ The token is stored in plain-text in the state. This is a potential security issue.

## Example Usage
```hcl
resource "sumologic_cloudsyslog_source" "cloudsyslog_source" {
  name         = "CLOUDSYSLOG"
  description  = "My description"
  category     = "my/source/category"
  collector_id = "${sumologic_collector.collector.id}"
}

resource "sumologic_collector" "collector" {
  name        = "my-collector"
  description = "Just testing this"
}
```

## Argument reference

  * Common Source Properties(https://github.com/terraform-providers/terraform-provider-sumologic/tree/master/website#common-source-properties)

## Attributes reference

The following attributes are exported:

- `id` - The internal ID of the source.
- `token` - The token to use for sending data to this source.

## Import
Cloud Syslog sources can be imported using the collector and source IDs (`collector/source`), e.g.:

```hcl
terraform import sumologic_cloudsyslog_source.test 123/456
```

HTTP sources can be imported using the collector name and source name (`collectorName/sourceName`), e.g.:

```hcl
terraform import sumologic_cloudsyslog_source.test my-test-collector/my-test-source
```

[1]: https://help.sumologic.com/Send_Data/Sources/02Sources_for_Hosted_Collectors/Cloud_Syslog_Source
