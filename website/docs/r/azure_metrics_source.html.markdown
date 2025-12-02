---
layout: "sumologic"
page_title: "SumoLogic: sumologic_azure_metrics_source"
description: |-
  Provides a Sumologic Azure Metrics Source.
---

# sumologic_azure_metrics_source
Provides a [Sumologic Azure Metrics Source](https://help.sumologic.com/docs/send-data/hosted-collectors/microsoft-source/azure-metrics-source/)

## Example Usage
```hcl
locals {
  tagfilters = [{
    "type" = "AzureTagFilters"
    "namespace" = "Microsoft.ClassicStorage/storageAccounts"
    "tags" {
      "name" = "test-name-1"
      "values" = ["value1", "value2"]
    }
  },{
    "type" = "AzureTagFilters"
    "namespace" = "Microsoft.ClassicStorage/storageAccounts"
    "tags" {
      "name" = "test-name-2"
      "values" = ["value3"]
    }
  }]
}

resource "sumologic_azure_metrics_source" "terraform_azure_metrics_source" {
  name = "Azure Metrics Source"
  description = "My description"
  category = "azure/metrics"
  content_type = "AzureMetrics"
  scan_interval = 300000
  paused = false
  collector_id = "${sumologic_collector.collector.id}"

  authentication {
    type = "AzureClientSecretAuthentication"
    tenant_id = "azure_tenant_id"
    client_id = "azure_client_id"
    client_secret = "azure_client_secret"
  }

  path {
    type = "AzureMetricsPath"
    environment = "Azure"
    limit_to_regions = ["eastus2", "westeurope"]
    limit_to_namespaces = ["Microsoft.ClassicStorage/storageAccounts"]
    dynamic "azure_tag_filters" {
      for_each = local.tagfilters
        content {
          type = azure_tag_filters.value.type
          namespace = azure_tag_filters.value.namespace
          tags {
            name = azure_tag_filters.value.tags.name
            values = azure_tag_filters.value.tags.values
          }
        }
    }
  }
}

resource "sumologic_collector" "collector" {
  name        = "my-collector"
  description = "Just testing this"
}
```

## Argument reference
In addition to the [Common Source Properties](https://registry.terraform.io/providers/SumoLogic/sumologic/latest/docs#common-source-properties), the following arguments are supported:
- `content_type` - (Required) Must be `AzureMetrics`.
- `scan_interval` - (Optional) Time interval in milliseconds of scans for new data. The default is 300000 and the minimum value is 1000 milliseconds.
- `paused` - (Optional) When set to true, the scanner is paused. To disable, set to false. The default is false.
- `authentication` - (Required) Authentication details for connecting to ingest metrics from Azure.
     + `type` - (Required) Must be `AzureClientSecretAuthentication`.
     + `tenant_id` - (Required) Your tenant id collected from [Azure platform](https://help.sumologic.com/docs/send-data/hosted-collectors/microsoft-source/azure-metrics-source/#vendor-configuration).
     + `client_id` - (Required) Your client id collected from [Azure platform](https://help.sumologic.com/docs/send-data/hosted-collectors/microsoft-source/azure-metrics-source/#vendor-configuration).
     + `client_secret` - (Required) Your client secret collected from [Azure platform](https://help.sumologic.com/docs/send-data/hosted-collectors/microsoft-source/azure-metrics-source/#vendor-configuration).
 - `path` - (Required) The location to scan for new data.
     + `type` - (Required) Must be `AzureMetricsPath`.
     + `environment` - (Required) The Azure environment from which to collect metrics.
     + `limit_to_regions` - (Optional) The list of regions from which to collect metrics. By default all regions are selected.
     + `limit_to_namespaces` - (Optional) The list of namespaces from which to collect metrics. By default all namespaces are selected.
     + `azure_tag_filters` - (Optional) Tag filters allow you to filter the Azure metrics by the tags you have assigned to your Azure resources. You can define tag filters for each supported namespace. If you do not define any tag filters, all metrics will be collected for the namespaces configured above.
          + `type` - (Required) Must be `AzureTagFilters`
          + `namespace` - Namespace for which you want to define the tag filters.
          + `tags` - List of key and value pairs of tag filters.
              + `name`: The name of the tag.
              + `values`: The list of accepted values for the tag name.
