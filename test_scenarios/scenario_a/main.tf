# Scenario A: FF ON + Old Terraform (singular notification)
# Purpose: Verify backward compatibility — singular `notification` still works when FF is ON

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
  name                = "TF_Scenario_A_SingularNotif_FF_ON"
  description         = "Scenario A: FF ON + old terraform using singular notification block"
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
        to_list                = ["scenario_a_old_tf@sumologic.com"]
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
