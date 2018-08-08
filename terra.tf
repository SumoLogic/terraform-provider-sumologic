provider "sumologic" {}

resource "sumologic_collector" "variablename" {
  name="collectorname"
  description="description"
  category="collectorcategory"
  timezone="Etc/UTC"
  lookup_by_name=false
  destroy=true
}

resource "sumologic_http_source" "httpsource" {
  name="httpsourcename"
  collector_id="${sumologic_collector.variablename.id}"
  category="sourcecategory"
  host_name="hostname"
  timezone="Etc/UTC"
  automatic_date_parsing=true
  multiline_processing_enabled=false
  use_autoline_matching=true
  manual_prefix_regexp=""
  force_timezone=true
  default_date_formats {
    format="YYYY"
    locator="\\[time=(.*)\\]"
  }
  //cutoff_timestamp=1234567890123
  cutoff_relative_time="-24h"
  lookup_by_name=false
  destroy=true
  message_per_request=true

  filters {
    name="excludefilter"
    filter_type="Exclude"
    regexp=".*regexp.*"
  }
}