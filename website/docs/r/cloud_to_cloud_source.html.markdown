---
layout: "sumologic"
page_title: "SumoLogic: sumologic_cloud_to_cloud_source"
description: |-
  Provides a Sumologic Cloud-to-Cloud source.
---

# sumologic_cloud_to_cloud_source
Provides a [Sumologic Cloud-to-Cloud source][1].

## Supported Integrations
List of available integrations along with their corresponding `JSON` templates is present [here][2] 

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

 - `schema_ref` - (Required) Source schema details. 
     + `type` - (Required) Schema type for the Cloud-to-Cloud integration source. Available schema types can be found [here][2].
 - `config` - (Required) This is a JSON object which contains the configuration parameters for the Source. Each schema type requires different JSON parameters. Refer to `JSON Configuration` and `Config Parameters` sections in the integration page for the specific `type` you have chosen to create.


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

[1]: https://help.sumologic.com/03Send-Data/Sources/02Sources-for-Hosted-Collectors/Cloud-to-Cloud_Integration_Framework
[2]: https://help.sumologic.com/03Send-Data/Sources/02Sources-for-Hosted-Collectors/Cloud-to-Cloud_Integration_Framework#Integrations
