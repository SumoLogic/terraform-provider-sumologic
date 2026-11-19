---
layout: "sumologic"
page_title: "SumoLogic: sumologic_azure_event_hub_log_source"
description: |-
  Provides a Sumologic Azure Event Hub Log Source.
---

# sumologic_azure_event_hub_log_source
Provides a [Sumologic Azure Event Hub Log Source](https://help.sumologic.com/docs/send-data/collect-from-other-data-sources/azure-monitoring/ms-azure-event-hubs-source/).

__IMPORTANT:__ The Azure Event Hub credentials are stored in plain-text in the state. This is a potential security issue.

## Example Usage

The example declares the hosted collector it attaches the source to, and takes the Event Hub credentials from variables.

### Collect Azure Event Hub logs
```hcl
resource "sumologic_collector" "collector" {
  name        = "my-collector"
  description = "Hosted collector for Azure Event Hub logs"
}

variable "event_hub_policy_name" {
  type        = string
  description = "Shared access policy name, for example RootManageSharedAccessKey"
}

variable "event_hub_policy_key" {
  type      = string
  sensitive = true
}

resource "sumologic_azure_event_hub_log_source" "terraform_azure_event_hub_log_source" {
  name         = "Azure Event Hub Log Source"
  description  = "My description"
  category     = "azure/eventhub"
  content_type = "AzureEventHubLog"
  collector_id = sumologic_collector.collector.id

  authentication {
    type                      = "AzureEventHubAuthentication"
    shared_access_policy_name = var.event_hub_policy_name
    shared_access_policy_key  = var.event_hub_policy_key
  }

  path {
    type           = "AzureEventHubPath"
    namespace      = "my-event-hub-namespace"
    event_hub_name = "my-hub"
    consumer_group = "$Default"
    region         = "Commercial"
  }
}
```

## Argument reference

In addition to the [Common Source Properties](https://registry.terraform.io/providers/SumoLogic/sumologic/latest/docs#common-source-properties), the following arguments are supported:

 - `content_type` - (Required, forces a new resource) Must be `AzureEventHubLog`.
 - `authentication` - (Required, forces a new resource; exactly one block) Authentication details for connecting to Azure Event Hub.
     + `type` - (Required) Must be `AzureEventHubAuthentication`.
     + `shared_access_policy_name` - (Required) Your shared access policy name.
     + `shared_access_policy_key` - (Required) Your shared access policy key (stored in state in plain text).
 - `path` - (Required, forces a new resource; exactly one block) The location to scan for new data.
     + `type` - (Required) Must be `AzureEventHubPath`.
     + `namespace` - (Required) The namespace of the event hub. 
     + `event_hub_name` - (Required) The name of the event hub.
     + `consumer_group` - (Required) The consumer group of the event hub. If you use a custom consumer group, it must already exist for the event hub instance.
     + `region` - (Optional) The region of the event hub. The value can be either `Commercial` for Azure, or `US Gov` for Azure Government. Defaults to `Commercial`.

 - `scan_interval` - Not applicable. The provider ignores any value set for it on Azure Event Hub sources; events are read from the event hub continuously.

### See also
  * [Common Source Properties](https://registry.terraform.io/providers/SumoLogic/sumologic/latest/docs#common-source-properties)
  * [Azure Event Hubs Source](https://www.sumologic.com/help/docs/send-data/collect-from-other-data-sources/azure-monitoring/ms-azure-event-hubs-source/)

## Attributes Reference
The following attributes are exported:

- `id` - The internal ID of the source.
- `url` - The URL of the source.

## Import
Azure Event Hub Log sources can be imported using the collector and source IDs (`collector/source`), e.g.:

```hcl
terraform import sumologic_azure_event_hub_log_source.test 123/456
```

Azure Event Hub Log sources can be imported using the collector name and source name (`collectorName/sourceName`), e.g.:

```hcl
terraform import sumologic_azure_event_hub_log_source.test my-test-collector/my-test-source
```
