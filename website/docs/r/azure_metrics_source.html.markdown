---
layout: "sumologic"
page_title: "SumoLogic: sumologic_azure_metrics_source"
description: |-
  Provides a Sumologic Azure Metrics Source.
---

# sumologic_azure_metrics_source
Provides a [Sumologic Azure Metrics Source](https://help.sumologic.com/docs/send-data/hosted-collectors/microsoft-source/azure-metrics-source/)

## Example Usage

Each example declares the hosted collector it attaches the source to, so it can be copied on its own.

### Collect all regions and namespaces
```hcl
resource "sumologic_collector" "collector" {
  name        = "my-collector"
  description = "Hosted collector for Azure metrics"
}

resource "sumologic_azure_metrics_source" "all_metrics" {
  name          = "Azure Metrics Source"
  description   = "My description"
  category      = "azure/metrics"
  content_type  = "AzureMetrics"
  scan_interval = 300000
  collector_id  = sumologic_collector.collector.id

  authentication {
    type          = "AzureClientSecretAuthentication"
    tenant_id     = var.azure_tenant_id
    client_id     = var.azure_client_id
    client_secret = var.azure_client_secret
  }

  path {
    type        = "AzureMetricsPath"
    environment = "Azure"
  }
}

variable "azure_tenant_id" {
  type = string
}

variable "azure_client_id" {
  type = string
}

variable "azure_client_secret" {
  type      = string
  sensitive = true
}
```

### Limit by region, namespace and resource tags
```hcl
resource "sumologic_collector" "collector" {
  name        = "my-collector"
  description = "Hosted collector for Azure metrics"
}

resource "sumologic_azure_metrics_source" "filtered_metrics" {
  name          = "Azure Metrics Source (filtered)"
  description   = "Storage account metrics for tagged resources"
  category      = "azure/metrics"
  content_type  = "AzureMetrics"
  scan_interval = 300000
  collector_id  = sumologic_collector.collector.id

  authentication {
    type          = "AzureClientSecretAuthentication"
    tenant_id     = var.azure_tenant_id
    client_id     = var.azure_client_id
    client_secret = var.azure_client_secret
  }

  path {
    type                = "AzureMetricsPath"
    environment         = "Azure"
    limit_to_regions    = ["eastus2", "westeurope"]
    limit_to_namespaces = ["Microsoft.ClassicStorage/storageAccounts"]

    azure_tag_filters {
      type      = "AzureTagFilters"
      namespace = "Microsoft.ClassicStorage/storageAccounts"
      tags {
        name   = "test-name-1"
        values = ["value1", "value2"]
      }
      tags {
        name   = "test-name-2"
        values = ["value3"]
      }
    }
  }
}

variable "azure_tenant_id" {
  type = string
}

variable "azure_client_id" {
  type = string
}

variable "azure_client_secret" {
  type      = string
  sensitive = true
}
```

## Argument reference
In addition to the [Common Source Properties](https://registry.terraform.io/providers/SumoLogic/sumologic/latest/docs#common-source-properties), the following arguments are supported:
- `content_type` - (Required, forces a new resource) Must be `AzureMetrics`.
- `scan_interval` - (Optional) Time interval in milliseconds of scans for new data. Defaults to `300000` (5 minutes). The provider accepts any integer; the Sumo Logic help page lists `1m` or `5m` as the supported options (see the [help page](https://www.sumologic.com/help/docs/send-data/hosted-collectors/microsoft-source/azure-metrics-source/)), so use `60000` or `300000`.
- `authentication` - (Required, forces a new resource; exactly one block) Authentication details for connecting to ingest metrics from Azure.
     + `type` - (Required) Must be `AzureClientSecretAuthentication`.
     + `tenant_id` - (Required) Your tenant id collected from [Azure platform](https://help.sumologic.com/docs/send-data/hosted-collectors/microsoft-source/azure-metrics-source/#vendor-configuration).
     + `client_id` - (Required) Your client id collected from [Azure platform](https://help.sumologic.com/docs/send-data/hosted-collectors/microsoft-source/azure-metrics-source/#vendor-configuration).
     + `client_secret` - (Required) Your client secret collected from [Azure platform](https://help.sumologic.com/docs/send-data/hosted-collectors/microsoft-source/azure-metrics-source/#vendor-configuration).
 - `path` - (Required, forces a new resource; exactly one block) The location to scan for new data.
     + `type` - (Required) Must be `AzureMetricsPath`.
     + `environment` - (Required) The Azure environment from which to collect metrics, for example `Azure`.
     + `limit_to_regions` - (Optional) The list of regions from which to collect metrics. By default all regions are selected.
     + `limit_to_namespaces` - (Optional) The list of namespaces from which to collect metrics. By default all namespaces are selected.
     + `azure_tag_filters` - (Optional) Tag filters allow you to filter the Azure metrics by the tags you have assigned to your Azure resources. You can define tag filters for each supported namespace. If you do not define any tag filters, all metrics will be collected for the namespaces configured above.
          + `type` - (Required) Must be `AzureTagFilters`
          + `namespace` - (Optional) Namespace for which you want to define the tag filters.
          + `tags` - (Optional) Repeatable block of tag name and accepted values.
              + `name` - (Required) The name of the tag.
              + `values` - (Optional) The list of accepted values for the tag name.

### See also
  * [Common Source Properties](https://registry.terraform.io/providers/SumoLogic/sumologic/latest/docs#common-source-properties)
  * [Azure Metrics Source](https://www.sumologic.com/help/docs/send-data/hosted-collectors/microsoft-source/azure-metrics-source/) (the help page recommends one source per Azure subscription)

## Attributes Reference
The following attributes are exported:

- `id` - The internal ID of the source.
- `url` - The URL of the source.

## Import
Azure Metrics sources can be imported using the collector and source IDs (`collector/source`), e.g.:

```hcl
terraform import sumologic_azure_metrics_source.test 123/456
```

Azure Metrics sources can be imported using the collector name and source name (`collectorName/sourceName`), e.g.:

```hcl
terraform import sumologic_azure_metrics_source.test my-test-collector/my-test-source
```
