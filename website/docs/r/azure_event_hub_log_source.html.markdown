---
layout: "sumologic"
page_title: "SumoLogic: sumologic_azure_event_hub_log_source"
description: |-
  Provides a Sumologic Azure Event Hub Log Source.
---

# sumologic_azure_event_hub_log_source
Provides a [Sumologic Azure Event Hub Log Source][2].

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

 - `content_type` - (Required) The content-type of the collected data. Details can be found in the [Sumologic documentation for hosted sources][1].
 - `authentication` - (Required) Authentication details for connecting to the S3 bucket.
     + `type` - (Required) Must be `AzureEventHubAuthentication`.
     + `shared_access_policy_name` - (Required) Your shared access policy key name.
     + `shared_access_policy_key` - (Required) Your shared access policy key.
 - `path` - (Required) The location to scan for new data.
     + `type` - (Required) Must be `AzureEventHubPath`.
     + `namespace` - (Required) The namespace of the event hub. 
     + `event_hub_name` - (Required) The name of the event hub.
     + `consumer_group` - (Required) The consumer group of the event hub.
     + `region` - (Optional) The value can be either "Commercial" for azure or "US Gov" for azure gov. Default is for azure. 

[1]: https://help.sumologic.com/Send_Data/Sources/03Use_JSON_to_Configure_Sources/JSON_Parameters_for_Hosted_Sources
[2]: https://help.sumologic.com/03Send-Data/Sources/02Sources-for-Hosted-Collectors/XXX
