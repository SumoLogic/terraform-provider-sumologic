resource "sumologic_monitor" "test_metrics_static_minDataPoints6" {
    name                = "Demo Terraform 13 - Metrics Static Monitor - Min Data Points"
    content_type        = "Monitor"
    group_notifications = true
    monitor_type        = "Metrics"

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
			min_data_points = 7
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
		min_data_points = 9
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
