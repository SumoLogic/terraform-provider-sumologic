---
layout: "sumologic"
page_title: "SumoLogic: sumologic_gcp_source"
description: |-
  Provides a Sumo Logic Google Cloud Platform Source.
---

# sumologic_gcp_source
Provides a [Sumo Logic Google Cloud Platform Source][2].

***Note:*** Google no longer requires a pub/sub domain to be [verified][3]. You no longer have to set up domain verification with your GCP Source endpoint.

## Example Usage
```hcl

resource "sumologic_gcp_source" "terraform_gcp_source" {
  name          = "GCP Source"
  description   = "My description"
  category      = "gcp"
  collector_id  = "${sumologic_collector.collector.id}"
}

resource "sumologic_collector" "collector" {
  name        = "my-collector"
  description = "Just testing this"
}
```

## Argument reference
  * [Common Source Properties](https://github.com/SumoLogic/terraform-provider-sumologic/tree/master/website#common-source-properties)

## Attributes Reference
The following attributes are exported:

- `id` - The internal ID of the source.
- `url` - The HTTP endpoint to use for sending data to this source.

## Import
Sumo Logic Google Cloud Platform sources can be imported using the collector and source IDs (`collector/source`), e.g.:

```hcl
terraform import sumologic_gcp_source.test 100000001/100000001
```

Sumo Logic Google Cloud Platform sources can be imported using the collector name and source name (`collectorName/sourceName`), e.g.:

```hcl
terraform import sumologic_gcp_source.test my-test-collector/my-test-source
```

[1]: https://help.sumologic.com/Send_Data/Sources/03Use_JSON_to_Configure_Sources/JSON_Parameters_for_Hosted_Sources
[2]: https://help.sumologic.com/03Send-Data/Sources/02Sources-for-Hosted-Collectors/Google-Cloud-Platform-Source
[3]: https://cloud.google.com/pubsub/docs/push
