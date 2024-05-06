
terraform {
  required_providers {
    sumologic = {
      source = "sumologic/sumologic"
      version = "~> 2.28.3" # set the Sumo Logic Terraform Provider version
    }
  }
  required_version = ">= 0.13"
}

provider "sumologic" {
  base_url   = "https://stag-api.sumologic.net/api/"
  access_id   = "suKiV4IMObZbEh"
  access_key  = "8kiP2y2paTe6vUMWKAqeEZPUPQfnOR7EMgrrxSOO9puS28iE5GStUFLM2SgWJ4w3"
}

#resource "sumologic_monitor" "pgupta_tf_logs_monitor_1" {
#  name         = "pgupta Terraform Logs Monitor"
#  description  = "tf logs monitor"
#  type         = "MonitorsLibraryMonitor"
#  is_disabled  = false
#  content_type = "Monitor"
#  monitor_type = "Logs"
#  evaluation_delay = "5m"
#  tags = {
#    "team" = "metrics"
#    "application" = "sumologic"
#  }
#
#  queries {
#    row_id = "A"
#    query  = "_sourceCategory=event-action info"
#  }
#
#  trigger_conditions {
#    logs_static_condition {
#      critical {
#        time_range = "15m"
#        alert {
#          threshold      = 40.0
#          threshold_type = "GreaterThan"
#        }
#        resolution {
#          threshold      = 40.0
#          threshold_type = "LessThanOrEqual"
#        }
#      }
#    }
#  }
#}

data "sumologic_folder" "pgupta" {
  path = "/Library/Users/pgupta@sumologic.com"
}

resource "sumologic_lookup_table" "pgupta_lookup" {
  name = "Sample Lookup Table"
  fields {
    field_name = "prefix"
    field_type = "string"
  }
  fields {
    field_name = "assembly"
    field_type = "string"
  }
  ttl               = 100
  primary_keys      = ["prefix", "assembly"]
  parent_folder_id  = "${data.sumologic_folder.pgupta.id}"
  size_limit_action = "DeleteOldData"
  description       = "some new description ieddss"
  csv_file_path     = "sumologic/prefix_assembly_csv.csv"

}
