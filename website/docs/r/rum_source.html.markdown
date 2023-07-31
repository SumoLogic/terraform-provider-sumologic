---
layout: "sumologic"
page_title: "SumoLogic: sumologic_rum_source"
description: |-
  Provides a Sumologic Rum Source.
---

# sumologic_rum_source
Provides a Sumologic Rum Source.

## Example Usage
```hcl
resource "sumologic_collector" "collector" {
  name       = "test-collector"
  category   = "macos/test"
}

resource "sumologic_rum_source" "testRumSource" {
  name       = "rum_source_test"
  description = "Rum source created via terraform"
  category = "source/category"
  collector_id = sumologic_collector.collector.id
  path {
    application_name = "test_application"
    service_name = "test_service"
    deployment_environment = "test_environment"
    sampling_rate = 0.5
    ignore_urls = ["/^https:\\/\\/www.tracker.com\\/.*/", "/^https:\\/\\/api.mydomain.com\\/log\\/.*/"]
    custom_tags = { test_tag = "test_value" }
    propagate_trace_header_cors_urls = ["/^https:\\/\\/api.mydomain.com\\/apiv3\\/.*/", "/^https:\\/\\/www.3rdparty.com\\/.*/"]
    selected_country = "Poland"
  }
}
  ```

## Argument Reference

In addition to the common properties, the following arguments are supported:

 - `path`
     + `application_name` - (Optional) (Recommended) Add an Application Name tag of a text string to show for the app name in spans (for example, bookings-app). This groups services in the Application Service View. If left blank, services will belong to a "default" application.
     + `service_name` - (Required) Add a Service Name of a text string to show for the service name in spans (for example, "bookings-web-app").
     + `deployment_environment` - (Optional) Your production, staging, or development environment name.
     + `sampling_rate` - (Optional) Add a Probabilistic sampling rate for heavy traffic sites in a decimal value based on percentage, for example, 10% would be entered as 0.1. Supports floating values between 0.0 and 1.0, defaults to 1.0 (all data is passed).
     + `ignore_urls` - (Optional) Add a list of URLs not to collect trace data from. Supports regex. Make sure provided URLs are valid JavaScript flavor regexes. For example: "/^https:\/\/www.tracker.com\/.*/, /^https:\/\/api.mydomain.com\/log\/.*/"
     + `custom_tags` - (Optional) Defines custom tags attached to the spans. For example: "internal.version = 0.1.21"
     + `propagate_trace_header_cors_urls` - (Optional) (Recommended) Add a list of URLs or URL patterns that pass tracing context to construct traces end-to-end. Provided URLs should be valid JavaScript flavor regexes. Some examples are "/^https:\/\/api.mydomain.com\/apiv3\/.*/" and "/^https:\/\/www.3rdparty.com\/.*/".
     + `selected_country` - (Optional) Specify if you want to enrich spans with the details level up to the city - if left blank, enrichment works down to the state level.

### See also
  * [Common Source Properties](https://github.com/terraform-providers/terraform-provider-sumologic/tree/master/website#common-source-properties)
  * [Configuration of the Rum collection](https://help.sumologic.com/docs/apm/real-user-monitoring/configure-data-collection/)

## Attributes Reference
The following attributes are exported:

- `id` - The internal ID of the Rum source.
- `url` - The HTTP endpoint to be used while sending the data to the Sumo.

## Import
Rum sources can be imported using the collector and source IDs, e.g.:

```hcl
terraform import sumologic_rum_source.test 123/456
```

Rum sources can also be imported using the collector name and source name, e.g.:

```hcl
terraform import sumologic_rum_source.test my-test-collector/my-test-source
```