# Scenario B: FF OFF + New Terraform (plural notifications)
# Purpose: What happens when new terraform sends `notifications` but BE FF is OFF
# Expected: Create may succeed but drift on subsequent plan

terraform {
  required_providers {
    sumologic = {
      source  = "sumologic/sumologic"
      version = ">= 0.0.1"
    }
  }
}

provider "sumologic" {}

resource "sumologic_log_search" "test" {
  name                = "TF_Scenario_B_PluralNotif_FF_OFF"
  description         = "Scenario B: FF OFF + new terraform using notifications (plural)"
  query_string        = "_sourceCategory=* | count by _sourceHost"
  parsing_mode        = "Manual"
  parent_id           = "000000000000271F"
  run_by_receipt_time = false

  time_range {
    begin_bounded_time_range {
      from {
        relative_time_range {
          relative_time = "-15m"
        }
      }
    }
  }

  schedule {
    cron_expression   = "0 0/15 * * * ? *"
    mute_error_emails = false
    notifications {
      email_search_notification {
        include_csv_attachment = false
        include_histogram      = true
        include_query          = true
        include_result_set     = true
        subject_template       = "Alert 1: {{SearchName}}"
        to_list                = ["scenario_b_notif1@sumologic.com"]
      }
    }
    notifications {
      email_search_notification {
        include_csv_attachment = false
        include_histogram      = false
        include_query          = false
        include_result_set     = true
        subject_template       = "Alert 2: {{SearchName}}"
        to_list                = ["scenario_b_notif2@sumologic.com"]
      }
    }
    parseable_time_range {
      begin_bounded_time_range {
        from {
          relative_time_range {
            relative_time = "-15m"
          }
        }
      }
    }
    schedule_type = "15Minutes"
    time_zone     = "Asia/Kolkata"
  }
}

output "log_search_id" {
  value = sumologic_log_search.test.id
}
