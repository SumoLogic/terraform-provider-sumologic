# Scenario C: FF OFF → create SS → FF ON → terraform plan
# Purpose: Migration path — SS created with singular notification before FF enabled
# Steps:
#   1. FF OFF: terraform apply
#   2. FF OFF: terraform plan → should show "No changes"
#   3. Turn FF ON
#   4. terraform plan → should STILL show "No changes" (no drift)

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
  name                = "TF_Scenario_C_Migration_Path"
  description         = "Scenario C: Created with FF OFF, then FF turned ON"
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
    notification {
      email_search_notification {
        include_csv_attachment = false
        include_histogram      = true
        include_query          = true
        include_result_set     = true
        subject_template       = "Alert: {{TriggerCondition}} for {{SearchName}}"
        to_list                = ["scenario_c_migration@sumologic.com"]
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
