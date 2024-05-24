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
```hcl
  resource "sumologic_azure_event_hub_log_source" "terraform_azure_event_hub_log_source" {
    name          = "Azure Event Hub Log Source"
    description   = "My description"
    category      = "azure/eventhub"
    content_type  = "AzureEventHubLog"
    collector_id  = "${sumologic_collector.collector.id}"
    authentication {
      type = "AzureEventHubAuthentication"
      shared_access_policy_name = "%s"
      shared_access_policy_key = "%s"
    }
    path {
      type = "AzureEventHubPath"
      namespace     = "%s"
      event_hub_name = "%s"
      consumer_group = "%s"
      region = "%s"
    }
  }
  
  resource "sumologic_collector" "collector" {
    name        = "my-collector"
    description = "Just testing this"
  }
```

## Argument reference

In addition to the [Common Source Properties](https://registry.terraform.io/providers/SumoLogic/sumologic/latest/docs#common-source-properties), the following arguments are supported:

 - `content_type` - (Required) Must be `AzureEventHubLog`.
 - `authentication` - (Required) Authentication details for connecting to Azure Event Hub.
     + `type` - (Required) Must be `AzureEventHubAuthentication`.
     + `shared_access_policy_name` - (Required) Your shared access policy name.
     + `shared_access_policy_key` - (Required) Your shared access policy key.
 - `path` - (Required) The location to scan for new data.
     + `type` - (Required) Must be `AzureEventHubPath`.
     + `namespace` - (Required) The namespace of the event hub. 
     + `event_hub_name` - (Required) The name of the event hub.
     + `consumer_group` - (Required) The consumer group of the event hub.
     + `region` - (Optional) The region of the event hub. The value can be either `Commercial` for Azure, or `US Gov` for Azure Government. Defaults to `Commercial`.
