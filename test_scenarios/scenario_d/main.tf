# Scenario D: FF ON + New Terraform + Multiple Notifications (Happy Path)
# Purpose: Full end-to-end — create SS with multiple notifications when FF is ON
# Expected: Create succeeds, no drift, import works

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
  name                = "TF_Scenario_D_MultiNotif_HappyPath"
  description         = "Scenario D: FF ON + new terraform with multiple notifications"
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
        subject_template       = "Alert 1: {{TriggerCondition}} for {{SearchName}}"
        to_list                = ["scenario_d_team_a@sumologic.com"]
      }
    }
    notifications {
      email_search_notification {
        include_csv_attachment = true
        include_histogram      = false
        include_query          = true
        include_result_set     = true
        subject_template       = "Alert 2: {{TriggerCondition}} for {{SearchName}}"
        to_list                = ["scenario_d_team_b@sumologic.com"]
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
    threshold {
      count          = 10
      operator       = "gt"
      threshold_type = "group"
    }
    time_zone = "Asia/Kolkata"
  }
}

output "log_search_id" {
  value = sumologic_log_search.test.id
}
