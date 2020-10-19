---
layout: "sumologic"
page_title: "SumoLogic: sumologic_cloud_to_cloud_source"
description: |-
  Provides a Sumologic Cloud-to-Cloud source.
---

# sumologic_cloud_to_cloud_source
Provides a [Sumologic Cloud-to-Cloud source][1].

__IMPORTANT:__ The API credentials are stored in plain-text in the state. This is a potential security issue.

## Example Usage
```hcl

resource "sumologic_cloud_to_cloud_source" "okta_source" {
 collector_id    = sumologic_collector.collector.id
 schema_ref = {
   type = "Okta"
   }
 config = jsonencode({"name":"okta source",
    "domain":"dev-xxx-admin.okta.com",
    "collectAll":true,
    "apiKey":"xxx",
    "fields":{
      "_siemForward":false
    },
    "pollingInterval": 30
    })

}
resource "sumologic_collector" "collector" {
  name        = "my-collector"
  description = "Just testing this"
}
```

## Argument reference
The following arguments are supported:

 - `config` - (Required) This is a JSON object which contains the configuration parameters for the Source.
 - `schema_ref` - (Required) Source schema details. 
     + `type` - (Required) Schema type for the Cloud-to-Cloud source.

## Attributes Reference
The following attributes are exported:

- `id` - The internal ID of the source.

## Import
Cloud-to-Cloud sources can be imported using the collector and source IDs (`collector/source`), e.g.:

```hcl
terraform import sumologic_cloud_to_cloud_source.test 100000001/100000001
```

Cloud-to-Cloud sources can be imported using the collector name and source name (`collectorName/sourceName`), e.g.:

```hcl
terraform import sumologic_cloud_to_cloud_source.test my-test-collector/my-test-source
```

[1]: https://help.sumologic.com/Beta/Cloud-to-Cloud_Integration_Framework
