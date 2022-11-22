resource "sumologic_monitor" "test_metrics_static_minDataPoints" {
    name                = "Demo Terraform 2 - Metrics Static Monitor - Min Data Points"
    content_type        = "Monitor"
    group_notifications = true
    monitor_type        = "Metrics"
evaluation_delay	= "1m"
    status              = [
        "Normal"
    ]
    queries {
        query  = "_sourceCategory=monitor-manager error"
        row_id = "A"
    }
     trigger_conditions {
         metrics_static_condition {
             critical {
                time_range = "-15m"
		occurrence_type = "Always"
                 alert {
                     threshold      = 100
                     threshold_type = "GreaterThanOrEqual"
                     min_data_points = 5
                 }
                 resolution {
                     threshold      = 80
                     threshold_type = "LessThan"
                 }
             }
             warning {
                 time_range = "-5m"
		 occurrence_type = "Always"
                 alert {
                     threshold      = 50
                     threshold_type = "GreaterThanOrEqual"
		     min_data_points = 2
                 }
                 resolution {
                     threshold      = 20
                     threshold_type = "LessThan"
                     occurrence_type = "AtLeastOnce"
                 }
             }
         }
     }
    notifications {
        run_for_trigger_types = [
            "Critical",
            "ResolvedCritical",
        ]
        notification {
            connection_type = "Email"
            recipients      = [
                "ndhar@sumologic.com"
            ]
            subject         = "Monitor Alert: {{TriggerType}} on {{SearchName}}"
            time_zone       = "PDT"
        }
    }
}
