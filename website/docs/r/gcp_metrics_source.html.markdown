---
layout: "sumologic"
page_title: "SumoLogic: sumologic_gcp_metrics_source"
description: |-
  Provides a Sumologic GCP Metrics Source.
---

# sumologic_gcp_metrics_source
Provides a `Sumologic GCP Metrics Source`

__IMPORTANT:__ The Service Account parameters (including private key) are stored in plain-text in the state. This is a potential security issue.

## Example Usage
```hcl

resource "sumologic_gcp_metrics_source" "terraform_gcp_metrics_source" {
  name          = "GCP Metrics Source"
  description   = "Description for GCP Metrics Source"
  category      = "gcp/metrics"
  content_type  = "GcpMetrics"
  scan_interval = 300000
  paused        = false
  collector_id  = "${sumologic_collector.collector.id}"

  authentication {
    type = "service_account"
    project_id = "service_account_project_id"
    private_key_id = "service_account_private_key_id"
    private_key = <<EOPK
service_account_private_key
EOPK
    client_email = "service_account_client_email"
    client_id = "service_account_client_id"
    auth_uri = "service_account_auth_uri"
    token_uri = "service_account_token_uri"
    auth_provider_x509_cert_url = "service_account_auth_provider_x509_cert_url"
    client_x509_cert_url = "service_account_client_x509_cert_url"
  }

  path {
    type = "GcpMetricsPath"
    limit_to_regions = ["us-east1", "us-central1", "asia-south1"]
    limit_to_services = ["Compute Engine", "Firebase", "App Engine"]
    custom_services = {
        service_name = "mysql"
        prefixes = ["cloudsql.googleapis.com/database/mysql/","cloudsql.googleapis.com/database/memory/","cloudsql.googleapis.com/database/cpu","cloudsql.googleapis.com/database/disk"]
    }
    custom_services = {
        service_name = "compute_instance_and_guests"
        prefixes = ["compute.googleapis.com/instance/","compute.googleapis.com/guest/"]
    }
  }

  lifecycle {
    ignore_changes = [authentication[0].private_key]
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
 - `scan_interval` - (Required) Time interval in milliseconds of scans for new data. The default is 300000 and the minimum value is 1000 milliseconds.
 - `paused` - (Required) When set to true, the scanner is paused. To disable, set to false.
 - `authentication` - (Required) Authentication details for connecting to the  GCP Monitoring using service_account credentials.
     + `type` - (Required) Must be `service_account`.
     + `project_id` - (Required) As per the service_account.json downloaded from GCP
     + `private_key_id` - (Required) As per the service_account.json downloaded from GCP
     + `private_key` - (Required) As per the service_account.json downloaded from GCP
     + `client_email` - (Required) As per the service_account.json downloaded from GCP
     + `client_id` - (Required) As per the service_account.json downloaded from GCP
     + `auth_uri` - (Required) As per the service_account.json downloaded from GCP
     + `token_uri` - (Required) As per the service_account.json downloaded from GCP
     + `auth_provider_x509_cert_url` - (Required) As per the service_account.json downloaded from GCP
     + `client_x509_cert_url` - (Required) As per the service_account.json downloaded from GCP

 - `path` - (Required) Details about what data to ingest
     + `type` - (Required) Type of polling source. This has to be `GcpMetricsPath`.
     + `limit_to_regions` - (Optional) List of regions for which metrics would be collected (Empty to collect from all regions)
     + `limit_to_services` - (Required) List of services from which metrics would be collected
     + `custom_services` - (Optional) Sumoloigc provides list of services that can be used in limit_to_services for which metrics would be collected. Custom Services allow you to define your own service w.r.t. metric collection. You can provide list of metric prefixes that should be collected as part of the custom service. This provides fine-grain control w.r.t. what all metrics are ingested by sumologic.
          + `service_name` - Name of the custom service you want to define.
          + `prefixes` - List of metric type prefixes. Eg: `["compute.googleapis.com/instance/","compute.googleapis.com/guest/"]`
 - `lifecycle` - (Required) describe fields that should be ignored by terraform while doing comparision i.e. Sumologic backend will return `private_key` as `********` 

### See also
  * [Common Source Properties](https://registry.terraform.io/providers/SumoLogic/sumologic/latest/docs#common-source-properties)

## Attributes Reference
The following attributes are exported:

- `id` - The internal ID of the source.

## Import
GCP Metrics sources can be imported using the collector and source IDs (`collector/source`), e.g.:

```hcl
terraform import sumologic_gcp_metrics_source.test 123/456
```

GCP Metrics sources can be imported using the collector name and source name (`collectorName/sourceName`), e.g.:

```hcl
terraform import sumologic_gcp_metrics_source.test my-test-collector/my-test-source
```

[1]: https://help.sumologic.com/Send_Data/Sources/03Use_JSON_to_Configure_Sources/JSON_Parameters_for_Hosted_Sources
