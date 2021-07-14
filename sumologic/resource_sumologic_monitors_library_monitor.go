package sumologic

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"log"
	"strings"
)

func resourceSumologicMonitorsLibraryMonitor() *schema.Resource {
	return &schema.Resource{
		Create: resourceSumologicMonitorsLibraryMonitorCreate,
		Read:   resourceSumologicMonitorsLibraryMonitorRead,
		Update: resourceSumologicMonitorsLibraryMonitorUpdate,
		Delete: resourceSumologicMonitorsLibraryMonitorDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{

			"version": {
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
			},

			"modified_at": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"is_system": {
				Type:     schema.TypeBool,
				Optional: true,
				Computed: true,
			},

			"content_type": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "Monitor",
			},

			"queries": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"row_id": {
							Type:     schema.TypeString,
							Required: true,
						},
						"query": {
							Type:     schema.TypeString,
							Required: true,
						},
					},
				},
			},

			"created_by": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"parent_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"is_disabled": {
				Type:     schema.TypeBool,
				Optional: true,
			},

			"is_mutable": {
				Type:     schema.TypeBool,
				Optional: true,
				Computed: true,
			},

			"triggers": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"static_condition": {
							Type:     schema.TypeList,
							MaxItems: 1,
							Optional: true,
							Elem: &schema.Resource{
								Schema: staticConditionSchema,
							},
						},
						"logs_static_condition": {
							Type:     schema.TypeList,
							MaxItems: 1,
							Optional: true,
							Elem: &schema.Resource{
								Schema: logsStaticConditionSchema,
							},
						},
						"metrics_static_condition": {
							Type:     schema.TypeList,
							MaxItems: 1,
							Optional: true,
							Elem: &schema.Resource{
								Schema: metricsStaticConditionSchema,
							},
						},
						"logs_outlier_condition": {
							Type:     schema.TypeList,
							MaxItems: 1,
							Optional: true,
							Elem: &schema.Resource{
								Schema: logsOutlierConditionSchema,
							},
						},
						"metrics_outlier_condition": {
							Type:     schema.TypeList,
							MaxItems: 1,
							Optional: true,
							Elem: &schema.Resource{
								Schema: metricsOutlierConditionSchema,
							},
						},
						"logs_missing_data_condition": {
							Type:     schema.TypeList,
							MaxItems: 1,
							Optional: true,
							Elem: &schema.Resource{
								Schema: logsMissingDataConditionSchema,
							},
						},
						"metrics_missing_data_condition": {
							Type:     schema.TypeList,
							MaxItems: 1,
							Optional: true,
							Elem: &schema.Resource{
								Schema: metricsMissingDataConditionSchema,
							},
						},
						"trigger_type": {
							Type:         schema.TypeString,
							Optional:     true,
							ValidateFunc: validation.StringInSlice([]string{"Critical", "Warning", "MissingData", "ResolvedCritical", "ResolvedWarning", "ResolvedMissingData"}, false),
						},
						"threshold": {
							Type:     schema.TypeFloat,
							Optional: true,
						},
						"threshold_type": {
							Type:         schema.TypeString,
							Optional:     true,
							ValidateFunc: validation.StringInSlice([]string{"LessThan", "LessThanOrEqual", "GreaterThan", "GreaterThanOrEqual"}, false),
						},
						"time_range": {
							Type:         schema.TypeString,
							Optional:     true,
							ValidateFunc: validation.StringInSlice([]string{"5m", "-5m", "10m", "-10m", "15m", "-15m", "30m", "-30m", "60m", "-60m", "1h", "-1h", "3h", "-3h", "6h", "-6h", "12h", "-12h", "24h", "-24h", "1d", "-1d"}, false),
						},
						"trigger_source": {
							Type:         schema.TypeString,
							Optional:     true,
							ValidateFunc: validation.StringInSlice([]string{"AllTimeSeries", "AnyTimeSeries", "AllResults"}, false),
						},
						"occurrence_type": {
							Type:         schema.TypeString,
							Optional:     true,
							ValidateFunc: validation.StringInSlice([]string{"AtLeastOnce", "Always", "ResultCount", "MissingData"}, false),
						},
						"detection_method": {
							Type:         schema.TypeString,
							Optional:     true,
							ValidateFunc: validation.StringInSlice([]string{"StaticCondition", "LogsStaticCondition", "MetricsStaticCondition", "LogsOutlierCondition", "MetricsOutlierCondition", "LogsMissingDataCondition", "MetricsMissingDataCondition"}, false),
						},
					},
				},
			},

			"notifications": {
				Type:     schema.TypeList,
				Optional: true,

				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"notification": {
							Type:     schema.TypeList,
							Required: true,
							MinItems: 1,
							MaxItems: 1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"action_type": {
										Type:       schema.TypeString,
										Optional:   true,
										Computed:   true,
										Deprecated: "The field `action_type` is deprecated and will be removed in a future release of the provider - please use `connection_type` instead.",
									},
									"connection_type": {
										Type:         schema.TypeString,
										Optional:     true,
										Computed:     true,
										ValidateFunc: validation.StringInSlice([]string{"Email", "AWSLambda", "AzureFunctions", "Datadog", "HipChat", "Jira", "NewRelic", "Opsgenie", "PagerDuty", "Slack", "MicrosoftTeams", "Webhook"}, false),
									},
									"subject": {
										Type:     schema.TypeString,
										Optional: true,
									},
									"recipients": {
										Type:     schema.TypeList,
										Optional: true,
										Elem: &schema.Schema{
											Type: schema.TypeString,
										},
									},
									"message_body": {
										Type:     schema.TypeString,
										Optional: true,
									},
									"time_zone": {
										Type:     schema.TypeString,
										Optional: true,
									},
									"connection_id": {
										Type:     schema.TypeString,
										Optional: true,
									},
									"payload_override": {
										Type:     schema.TypeString,
										Optional: true,
									},
								},
							},
						},
						"run_for_trigger_types": {
							Type:     schema.TypeList,
							Required: true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
					},
				},
			},

			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},

			"created_at": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"monitor_type": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringInSlice([]string{"Logs", "Metrics"}, false),
			},

			"is_locked": {
				Type:     schema.TypeBool,
				Optional: true,
				Computed: true,
			},
			"status": {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"group_notifications": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  true,
			},

			"type": {
				Type:         schema.TypeString,
				Optional:     true,
				Default:      "MonitorsLibraryMonitor",
				ValidateFunc: validation.StringInSlice([]string{"MonitorsLibraryMonitor", "MonitorsLibraryFolder"}, false),
			},

			"modified_by": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"post_request_map": {
				Type:     schema.TypeMap,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
		},
	}
}

func resourceSumologicMonitorsLibraryMonitorCreate(d *schema.ResourceData, meta interface{}) error {
	c := meta.(*Client)
	if d.Id() == "" {
		monitor := resourceToMonitorsLibraryMonitor(d)
		if monitor.ParentID == "" {
			rootFolder, err := c.GetMonitorsLibraryFolder("root")
			if err != nil {
				return err
			}

			monitor.ParentID = rootFolder.ID
		}
		paramMap := map[string]string{
			"parentId": monitor.ParentID,
		}
		monitorDefinitionID, err := c.CreateMonitorsLibraryMonitor(monitor, paramMap)
		if err != nil {
			return err
		}

		d.SetId(monitorDefinitionID)
	}
	return resourceSumologicMonitorsLibraryMonitorRead(d, meta)
}

func resourceSumologicMonitorsLibraryMonitorRead(d *schema.ResourceData, meta interface{}) error {
	c := meta.(*Client)

	monitor, err := c.MonitorsRead(d.Id())
	if err != nil {
		return err
	}

	if monitor == nil {
		log.Printf("[WARN] Monitor not found, removing from state: %v - %v", d.Id(), err)
		d.SetId("")
		return nil
	}

	d.Set("created_by", monitor.CreatedBy)
	d.Set("created_at", monitor.CreatedAt)
	d.Set("monitor_type", monitor.MonitorType)
	d.Set("modified_by", monitor.ModifiedBy)
	d.Set("is_mutable", monitor.IsMutable)
	d.Set("version", monitor.Version)
	d.Set("description", monitor.Description)
	d.Set("name", monitor.Name)
	d.Set("parent_id", monitor.ParentID)
	d.Set("modified_at", monitor.ModifiedAt)
	d.Set("content_type", monitor.ContentType)
	d.Set("is_locked", monitor.IsLocked)
	d.Set("is_system", monitor.IsSystem)
	d.Set("is_disabled", monitor.IsDisabled)
	d.Set("status", monitor.Status)
	d.Set("group_notifications", monitor.GroupNotifications)
	// set notifications
	notifications := make([]interface{}, len(monitor.Notifications))
	for i, n := range monitor.Notifications {
		// notification in schema should be a list of length exactly 1
		internalNotification := make(map[string]interface{})
		internalNotificationDict := n.Notification.(map[string]interface{})
		// log.Printf("monitor.Notification %v", n.Notification)
		if internalNotificationDict["connectionType"] != nil {
			internalNotification["connection_type"] = internalNotificationDict["connectionType"].(string)
		} else {
			// for backwards compatibility
			internalNotification["connection_type"] = internalNotificationDict["actionType"].(string)
			// convert from old action_type name to new connection_type name if applicable
			if internalNotification["connection_type"].(string) == "EmailAction" {
				internalNotification["connection_type"] = "Email"
			}
			if internalNotification["connection_type"].(string) == "NamedConnectionAction" {
				internalNotification["connection_type"] = "Webhook"
			}
		}
		if internalNotification["connection_type"].(string) == "Email" {
			// for backwards compatibility
			internalNotification["action_type"] = "EmailAction"
			internalNotification["subject"] = internalNotificationDict["subject"].(string)
			internalNotification["recipients"] = internalNotificationDict["recipients"].([]interface{})
			internalNotification["message_body"] = internalNotificationDict["messageBody"].(string)
			internalNotification["time_zone"] = internalNotificationDict["timeZone"].(string)
		} else {
			internalNotification["action_type"] = "NamedConnectionAction"
			internalNotification["connection_id"] = internalNotificationDict["connectionId"].(string)
			if internalNotificationDict["payloadOverride"] != nil {
				internalNotification["payload_override"] = internalNotificationDict["payloadOverride"].(string)
			}
		}

		schemaInternalNotification := []interface{}{
			internalNotification,
		}

		notifications[i] = map[string]interface{}{
			"notification":          schemaInternalNotification,
			"run_for_trigger_types": n.RunForTriggerTypes,
		}
	}
	if err := d.Set("notifications", notifications); err != nil {
		return err
	}
	// set triggers
	// NOTE: trigger blocks come in 2 forms:
	//  a. legacy version, where attributes are flattened out inside "triggers" block, and
	//  b. detection-method blocks, where each detection method gets its own sub-block.
	// Triggers of the detection method type 'StaticCondition' can exist in either form.
	// We need to make sure that when making updates, such as reading back a resource after apply,
	// we update them at their original version in order to ensure
	// that Terraform sees no state changes from apply.
	existingTriggers := []interface{}{}
	if val, ok := d.GetOk("triggers"); ok {
		existingTriggers = val.([]interface{})
	}
	triggers := make([]interface{}, len(monitor.Triggers))
	for i, t := range monitor.Triggers {
		if i < len(existingTriggers) && isLegacyTriggersBlock(existingTriggers[i].(map[string]interface{})) {
			triggers[i] = t.toLegacyTriggersBlock()
		} else {
			triggers[i] = t.toTriggersBlock()
		}
	}
	if err := d.Set("triggers", triggers); err != nil {
		return err
	}
	// set queries
	queries := make([]interface{}, len(monitor.Queries))
	for i, q := range monitor.Queries {
		queries[i] = map[string]interface{}{
			"row_id": q.RowID,
			"query":  q.Query,
		}
	}
	if err := d.Set("queries", queries); err != nil {
		return err
	}

	return nil
}

func (t *TriggerCondition) PositiveTimeRange() string {
	return strings.TrimPrefix(t.TimeRange, "-")
}

func (t *TriggerCondition) PositiveBaselineWindow() string {
	return strings.TrimPrefix(t.BaselineWindow, "-")
}

func (condition *TriggerCondition) toLegacyTriggersBlock() map[string]interface{} {
	return map[string]interface{}{
		"time_range":       condition.PositiveTimeRange(),
		"trigger_type":     condition.TriggerType,
		"threshold":        condition.Threshold,
		"threshold_type":   condition.ThresholdType,
		"occurrence_type":  condition.OccurrenceType,
		"trigger_source":   condition.TriggerSource,
		"detection_method": condition.DetectionMethod,
	}
}

func (condition *TriggerCondition) toTriggersBlock() map[string]interface{} {
	if condition == nil {
		return map[string]interface{}{}
	}
	switch condition.DetectionMethod {
	case staticConditionDetectionMethod:
		return condition.toStaticConditionTriggersBlock()
	case logsStaticConditionDetectionMethod:
		return condition.toLogsStaticConditionTriggersBlock()
	case metricsStaticConditionDetectionMethod:
		return condition.toMetricsStaticConditionTriggersBlock()
	case logsOutlierConditionDetectionMethod:
		return condition.toLogsOutlierConditionTriggersBlock()
	case metricsOutlierConditionDetectionMethod:
		return condition.toMetricsOutlierConditionTriggersBlock()
	case logsMissingDataConditionDetectionMethod:
		return condition.toLogsMissingDataConditionTriggersBlock()
	case metricsMissingDataConditionDetectionMethod:
		return condition.toMetricsMissingDataConditionTriggersBlock()
	default:
		log.Fatalln("Internal error: Bad TriggerCondition", *condition)
		return map[string]interface{}{}
	}
}

func (condition *TriggerCondition) toStaticConditionTriggersBlock() map[string]interface{} {
	return map[string]interface{}{
		staticConditionFieldName: []interface{}{map[string]interface{}{
			"time_range":      condition.PositiveTimeRange(),
			"trigger_type":    condition.TriggerType,
			"threshold":       condition.Threshold,
			"threshold_type":  condition.ThresholdType,
			"occurrence_type": condition.OccurrenceType,
			"trigger_source":  condition.TriggerSource,
			"field":           condition.Field,
		}},
	}
}

func (condition *TriggerCondition) toLogsStaticConditionTriggersBlock() map[string]interface{} {
	return map[string]interface{}{
		logsStaticConditionFieldName: []interface{}{map[string]interface{}{
			"time_range":     condition.PositiveTimeRange(),
			"trigger_type":   condition.TriggerType,
			"threshold":      condition.Threshold,
			"threshold_type": condition.ThresholdType,
			"field":          condition.Field,
		}},
	}
}

func (condition *TriggerCondition) toMetricsStaticConditionTriggersBlock() map[string]interface{} {
	return map[string]interface{}{
		metricsStaticConditionFieldName: []interface{}{map[string]interface{}{
			"time_range":      condition.PositiveTimeRange(),
			"trigger_type":    condition.TriggerType,
			"threshold":       condition.Threshold,
			"threshold_type":  condition.ThresholdType,
			"occurrence_type": condition.OccurrenceType,
		}},
	}
}

func (condition *TriggerCondition) toLogsOutlierConditionTriggersBlock() map[string]interface{} {
	return map[string]interface{}{
		logsOutlierConditionFieldName: []interface{}{map[string]interface{}{
			"trigger_type": condition.TriggerType,
			"window":       condition.Window,
			"consecutive":  condition.Consecutive,
			"direction":    condition.Direction,
			"threshold":    condition.Threshold,
			"field":        condition.Field,
		}},
	}
}

func (condition *TriggerCondition) toMetricsOutlierConditionTriggersBlock() map[string]interface{} {
	return map[string]interface{}{
		metricsOutlierConditionFieldName: []interface{}{map[string]interface{}{
			"trigger_type":    condition.TriggerType,
			"threshold":       condition.Threshold,
			"baseline_window": condition.PositiveBaselineWindow(),
			"direction":       condition.Direction,
		}},
	}
}

func (condition *TriggerCondition) toLogsMissingDataConditionTriggersBlock() map[string]interface{} {
	return map[string]interface{}{
		logsMissingDataConditionFieldName: []interface{}{map[string]interface{}{
			"trigger_type": condition.TriggerType,
			"time_range":   condition.PositiveTimeRange(),
		}},
	}
}

func (condition *TriggerCondition) toMetricsMissingDataConditionTriggersBlock() map[string]interface{} {
	return map[string]interface{}{
		metricsMissingDataConditionFieldName: []interface{}{map[string]interface{}{
			"trigger_type":   condition.TriggerType,
			"time_range":     condition.PositiveTimeRange(),
			"trigger_source": condition.TriggerSource,
		}},
	}
}

func resourceSumologicMonitorsLibraryMonitorUpdate(d *schema.ResourceData, meta interface{}) error {
	c := meta.(*Client)
	monitor := resourceToMonitorsLibraryMonitor(d)
	if d.HasChange("parentId") {
		// monitor.ParentID = d.Get("parentId").(string)
		err := c.MoveMonitorsLibraryMonitor(monitor)
		if err != nil {
			return err
		}
	}
	monitor.Type = "MonitorsLibraryMonitorUpdate"
	err := c.UpdateMonitorsLibraryMonitor(monitor)
	if err != nil {
		return err
	}
	updatedMonitor := resourceSumologicMonitorsLibraryMonitorRead(d, meta)

	return updatedMonitor
}

func resourceSumologicMonitorsLibraryMonitorDelete(d *schema.ResourceData, meta interface{}) error {
	c := meta.(*Client)
	monitor := resourceToMonitorsLibraryMonitor(d)
	err := c.DeleteMonitorsLibraryMonitor(monitor.ID)
	if err != nil {
		return err
	}
	return nil
}

func getNotifications(d *schema.ResourceData) []MonitorNotification {
	rawNotifications := d.Get("notifications").([]interface{})
	notifications := make([]MonitorNotification, len(rawNotifications))
	for i := range rawNotifications {
		notificationDict := rawNotifications[i].(map[string]interface{})
		rawNotificationAction := notificationDict["notification"].([]interface{})
		notificationActionDict := rawNotificationAction[0].(map[string]interface{})
		connectionType := ""
		actionType := ""
		if notificationActionDict["connection_type"] != nil &&
			notificationActionDict["connection_type"] != "" {
			connectionType = notificationActionDict["connection_type"].(string)
			actionType = connectionType
		} else {
			// for backwards compatibility
			actionType = notificationActionDict["action_type"].(string)
			connectionType = actionType
			// convert from old action_type name to new connection_type name if applicable
			if connectionType == "EmailAction" {
				connectionType = "Email"
			}
			if connectionType == "NamedConnectionAction" {
				connectionType = "Webhook"
			}
		}

		var n MonitorNotification
		if connectionType == "Email" {
			n.Notification = EmailNotification{
				ActionType:     "EmailAction",
				ConnectionType: connectionType,
				Subject:        notificationActionDict["subject"].(string),
				Recipients:     notificationActionDict["recipients"].([]interface{}),
				MessageBody:    notificationActionDict["message_body"].(string),
				TimeZone:       notificationActionDict["time_zone"].(string),
			}
		} else {
			n.Notification = WebhookNotificiation{
				ActionType:      "NamedConnectionAction",
				ConnectionType:  connectionType,
				ConnectionID:    notificationActionDict["connection_id"].(string),
				PayloadOverride: notificationActionDict["payload_override"].(string),
			}
		}
		n.RunForTriggerTypes = notificationDict["run_for_trigger_types"].([]interface{})
		notifications[i] = n
	}
	return notifications
}

func getTriggers(d *schema.ResourceData) []TriggerCondition {
	rawTriggers := d.Get("triggers").([]interface{})
	triggers := make([]TriggerCondition, len(rawTriggers))
	for i := range rawTriggers {
		triggerDict := rawTriggers[i].(map[string]interface{})
		triggers[i] = triggersBlockToTriggerCondition(triggerDict)
	}
	return triggers
}

func triggersBlockToTriggerCondition(triggerDict map[string]interface{}) TriggerCondition {
	if v, ok := getSingletonArrayFieldOk(triggerDict, staticConditionFieldName); ok {
		return staticConditionBlockToTriggerCondition(v)
	}
	if v, ok := getSingletonArrayFieldOk(triggerDict, logsStaticConditionFieldName); ok {
		return logsStaticConditionBlockToTriggerCondition(v)
	}
	if v, ok := getSingletonArrayFieldOk(triggerDict, metricsStaticConditionFieldName); ok {
		return metricsStaticConditionBlockToTriggerCondition(v)
	}
	if v, ok := getSingletonArrayFieldOk(triggerDict, logsOutlierConditionFieldName); ok {
		return logsOutlierConditionBlockToTriggerCondition(v)
	}
	if v, ok := getSingletonArrayFieldOk(triggerDict, metricsOutlierConditionFieldName); ok {
		return metricsOutlierConditionBlockToTriggerCondition(v)
	}
	if v, ok := getSingletonArrayFieldOk(triggerDict, logsMissingDataConditionFieldName); ok {
		return logsMissingDataConditionBlockToTriggerCondition(v)
	}
	if v, ok := getSingletonArrayFieldOk(triggerDict, metricsMissingDataConditionFieldName); ok {
		return metricsMissingDataConditionBlockToTriggerCondition(v)
	}
	// If we are here, it means this is a legacy block
	return legacyBlockToTriggerCondition(triggerDict)
}

/*
 Given a block with an inner map wrapped in a singleton array
   block {
     field {
       foo: bar
       ..
     }
   }
 returns the inner map
   ({ foo: bar, ... }, true)

 Otherwise returns (empty, false)
*/
func getSingletonArrayFieldOk(block map[string]interface{}, fieldName string) (map[string]interface{}, bool) {
	if v, ok := block[fieldName]; ok {
		if arr, ok := v.([]interface{}); ok && len(arr) == 1 {
			if mp, ok := arr[0].(map[string]interface{}); ok {
				return mp, true
			}
		}
	}
	return map[string]interface{}{}, false
}

func legacyBlockToTriggerCondition(block map[string]interface{}) TriggerCondition {
	return TriggerCondition{
		TriggerType:     block["trigger_type"].(string),
		Threshold:       block["threshold"].(float64),
		ThresholdType:   block["threshold_type"].(string),
		TimeRange:       block["time_range"].(string),
		OccurrenceType:  block["occurrence_type"].(string),
		TriggerSource:   block["trigger_source"].(string),
		DetectionMethod: block["detection_method"].(string),
	}
}

func staticConditionBlockToTriggerCondition(block map[string]interface{}) TriggerCondition {
	return TriggerCondition{
		TriggerType:     block["trigger_type"].(string),
		Threshold:       block["threshold"].(float64),
		ThresholdType:   block["threshold_type"].(string),
		TimeRange:       block["time_range"].(string),
		OccurrenceType:  block["occurrence_type"].(string),
		TriggerSource:   block["trigger_source"].(string),
		Field:           block["field"].(string),
		DetectionMethod: staticConditionDetectionMethod,
	}
}

func logsStaticConditionBlockToTriggerCondition(block map[string]interface{}) TriggerCondition {
	return TriggerCondition{
		TriggerType:     block["trigger_type"].(string),
		Threshold:       block["threshold"].(float64),
		ThresholdType:   block["threshold_type"].(string),
		Field:           block["field"].(string),
		TimeRange:       block["time_range"].(string),
		DetectionMethod: logsStaticConditionDetectionMethod,
	}
}

func metricsStaticConditionBlockToTriggerCondition(block map[string]interface{}) TriggerCondition {
	return TriggerCondition{
		TriggerType:     block["trigger_type"].(string),
		Threshold:       block["threshold"].(float64),
		ThresholdType:   block["threshold_type"].(string),
		TimeRange:       block["time_range"].(string),
		OccurrenceType:  block["occurrence_type"].(string),
		DetectionMethod: metricsStaticConditionDetectionMethod,
	}
}

func logsOutlierConditionBlockToTriggerCondition(block map[string]interface{}) TriggerCondition {
	return TriggerCondition{
		TriggerType:     block["trigger_type"].(string),
		Field:           block["field"].(string),
		Window:          block["window"].(int),
		Consecutive:     block["consecutive"].(int),
		Direction:       block["direction"].(string),
		Threshold:       block["threshold"].(float64),
		DetectionMethod: logsOutlierConditionDetectionMethod,
	}
}

func metricsOutlierConditionBlockToTriggerCondition(block map[string]interface{}) TriggerCondition {
	return TriggerCondition{
		TriggerType:     block["trigger_type"].(string),
		Threshold:       block["threshold"].(float64),
		BaselineWindow:  block["baseline_window"].(string),
		Direction:       block["direction"].(string),
		DetectionMethod: metricsOutlierConditionDetectionMethod,
	}
}

func logsMissingDataConditionBlockToTriggerCondition(block map[string]interface{}) TriggerCondition {
	return TriggerCondition{
		TriggerType:     block["trigger_type"].(string),
		TimeRange:       block["time_range"].(string),
		DetectionMethod: logsMissingDataConditionDetectionMethod,
	}
}

func metricsMissingDataConditionBlockToTriggerCondition(block map[string]interface{}) TriggerCondition {
	return TriggerCondition{
		TriggerType:     block["trigger_type"].(string),
		TimeRange:       block["time_range"].(string),
		TriggerSource:   block["trigger_source"].(string),
		DetectionMethod: metricsMissingDataConditionDetectionMethod,
	}
}

var staticConditionFieldName = "static_condition"
var logsStaticConditionFieldName = "logs_static_condition"
var metricsStaticConditionFieldName = "metrics_static_condition"
var logsOutlierConditionFieldName = "logs_outlier_condition"
var metricsOutlierConditionFieldName = "metrics_outlier_condition"
var logsMissingDataConditionFieldName = "logs_missing_data_condition"
var metricsMissingDataConditionFieldName = "metrics_missing_data_condition"
var staticConditionDetectionMethod = "StaticCondition"
var logsStaticConditionDetectionMethod = "LogsStaticCondition"
var metricsStaticConditionDetectionMethod = "MetricsStaticCondition"
var logsOutlierConditionDetectionMethod = "LogsOutlierCondition"
var metricsOutlierConditionDetectionMethod = "MetricsOutlierCondition"
var logsMissingDataConditionDetectionMethod = "LogsMissingDataCondition"
var metricsMissingDataConditionDetectionMethod = "MetricsMissingDataCondition"

func getQueries(d *schema.ResourceData) []MonitorQuery {
	rawQueries := d.Get("queries").([]interface{})
	queries := make([]MonitorQuery, len(rawQueries))
	for i := range rawQueries {
		queryDict := rawQueries[i].(map[string]interface{})
		queries[i] = MonitorQuery{
			Query: queryDict["query"].(string),
			RowID: queryDict["row_id"].(string),
		}
	}
	return queries
}

func resourceToMonitorsLibraryMonitor(d *schema.ResourceData) MonitorsLibraryMonitor {
	notifications := getNotifications(d)
	triggers := getTriggers(d)
	queries := getQueries(d)
	rawStatus := d.Get("status").([]interface{})
	status := make([]string, len(rawStatus))
	for i := range rawStatus {
		status[i] = rawStatus[i].(string)
	}

	return MonitorsLibraryMonitor{
		CreatedBy:          d.Get("created_by").(string),
		Name:               d.Get("name").(string),
		ID:                 d.Id(),
		CreatedAt:          d.Get("created_at").(string),
		MonitorType:        d.Get("monitor_type").(string),
		Description:        d.Get("description").(string),
		Queries:            queries,
		ModifiedBy:         d.Get("modified_by").(string),
		IsMutable:          d.Get("is_mutable").(bool),
		Version:            d.Get("version").(int),
		Notifications:      notifications,
		Type:               d.Get("type").(string),
		ParentID:           d.Get("parent_id").(string),
		ModifiedAt:         d.Get("modified_at").(string),
		Triggers:           triggers,
		ContentType:        d.Get("content_type").(string),
		IsLocked:           d.Get("is_locked").(bool),
		IsSystem:           d.Get("is_system").(bool),
		IsDisabled:         d.Get("is_disabled").(bool),
		Status:             status,
		GroupNotifications: d.Get("group_notifications").(bool),
	}
}

var staticConditionSchema = map[string]*schema.Schema{
	"trigger_type": {
		Type:         schema.TypeString,
		Required:     true,
		ValidateFunc: validation.StringInSlice([]string{"Critical", "Warning", "MissingData", "ResolvedCritical", "ResolvedWarning", "ResolvedMissingData"}, false),
	},
	"threshold": {
		Type:     schema.TypeFloat,
		Optional: true,
	},
	"threshold_type": {
		Type:         schema.TypeString,
		Optional:     true,
		ValidateFunc: validation.StringInSlice([]string{"LessThan", "LessThanOrEqual", "GreaterThan", "GreaterThanOrEqual"}, false),
	},
	"field": {
		Type:     schema.TypeString,
		Optional: true,
	},
	"time_range": {
		Type:         schema.TypeString,
		Required:     true,
		ValidateFunc: validation.StringInSlice([]string{"5m", "-5m", "10m", "-10m", "15m", "-15m", "30m", "-30m", "60m", "-60m", "1h", "-1h", "3h", "-3h", "6h", "-6h", "12h", "-12h", "24h", "-24h", "1d", "-1d"}, false),
	},
	"trigger_source": {
		Type:         schema.TypeString,
		Optional:     true,
		ValidateFunc: validation.StringInSlice([]string{"AllTimeSeries", "AnyTimeSeries", "AllResults"}, false),
	},
	"occurrence_type": {
		Type:         schema.TypeString,
		Required:     true,
		ValidateFunc: validation.StringInSlice([]string{"AtLeastOnce", "Always", "ResultCount", "MissingData"}, false),
	},
}

var logsStaticConditionSchema = map[string]*schema.Schema{
	"trigger_type": {
		Type:         schema.TypeString,
		Required:     true,
		ValidateFunc: validation.StringInSlice([]string{"Critical", "Warning", "ResolvedCritical", "ResolvedWarning"}, false),
	},
	"threshold": {
		Type:     schema.TypeFloat,
		Required: true,
	},
	"threshold_type": {
		Type:         schema.TypeString,
		Required:     true,
		ValidateFunc: validation.StringInSlice([]string{"LessThan", "LessThanOrEqual", "GreaterThan", "GreaterThanOrEqual"}, false),
	},
	"field": {
		Type:     schema.TypeString,
		Optional: true,
	},
	"time_range": {
		Type:         schema.TypeString,
		Required:     true,
		ValidateFunc: validation.StringInSlice([]string{"5m", "-5m", "10m", "-10m", "15m", "-15m", "30m", "-30m", "60m", "-60m", "1h", "-1h", "3h", "-3h", "6h", "-6h", "12h", "-12h", "24h", "-24h", "1d", "-1d"}, false),
	},
}

var metricsStaticConditionSchema = map[string]*schema.Schema{
	"trigger_type": {
		Type:         schema.TypeString,
		Required:     true,
		ValidateFunc: validation.StringInSlice([]string{"Critical", "Warning", "ResolvedCritical", "ResolvedWarning"}, false),
	},
	"threshold": {
		Type:     schema.TypeFloat,
		Required: true,
	},
	"threshold_type": {
		Type:         schema.TypeString,
		Required:     true,
		ValidateFunc: validation.StringInSlice([]string{"LessThan", "LessThanOrEqual", "GreaterThan", "GreaterThanOrEqual"}, false),
	},
	"time_range": {
		Type:         schema.TypeString,
		Required:     true,
		ValidateFunc: validation.StringInSlice([]string{"5m", "-5m", "10m", "-10m", "15m", "-15m", "30m", "-30m", "60m", "-60m", "1h", "-1h", "3h", "-3h", "6h", "-6h", "12h", "-12h", "24h", "-24h", "1d", "-1d"}, false),
	},
	"occurrence_type": {
		Type:         schema.TypeString,
		Required:     true,
		ValidateFunc: validation.StringInSlice([]string{"AtLeastOnce", "Always"}, false),
	},
}

var logsOutlierConditionSchema = map[string]*schema.Schema{
	"trigger_type": {
		Type:         schema.TypeString,
		Required:     true,
		ValidateFunc: validation.StringInSlice([]string{"Critical", "Warning", "ResolvedCritical", "ResolvedWarning"}, false),
	},
	"field": {
		Type:     schema.TypeString,
		Optional: true,
	},
	"window": {
		Type:         schema.TypeInt,
		Optional:     true,
		ValidateFunc: validation.IntAtLeast(1),
	},
	"consecutive": {
		Type:         schema.TypeInt,
		Optional:     true,
		ValidateFunc: validation.IntAtLeast(1),
	},
	"direction": {
		Type:         schema.TypeString,
		Optional:     true,
		ValidateFunc: validation.StringInSlice([]string{"Both", "Up", "Down"}, false),
	},
	"threshold": {
		Type:     schema.TypeFloat,
		Optional: true,
	},
}

var metricsOutlierConditionSchema = map[string]*schema.Schema{
	"trigger_type": {
		Type:         schema.TypeString,
		Required:     true,
		ValidateFunc: validation.StringInSlice([]string{"Critical", "Warning", "ResolvedCritical", "ResolvedWarning"}, false),
	},
	"threshold": {
		Type:     schema.TypeFloat,
		Optional: true,
	},
	"baseline_window": {
		Type:     schema.TypeString,
		Optional: true,
	},
	"direction": {
		Type:         schema.TypeString,
		Optional:     true,
		ValidateFunc: validation.StringInSlice([]string{"Both", "Up", "Down"}, false),
	},
}

var logsMissingDataConditionSchema = map[string]*schema.Schema{
	"trigger_type": {
		Type:         schema.TypeString,
		Required:     true,
		ValidateFunc: validation.StringInSlice([]string{"MissingData", "ResolvedMissingData"}, false),
	},
	"time_range": {
		Type:         schema.TypeString,
		Required:     true,
		ValidateFunc: validation.StringInSlice([]string{"5m", "-5m", "10m", "-10m", "15m", "-15m", "30m", "-30m", "60m", "-60m", "1h", "-1h", "3h", "-3h", "6h", "-6h", "12h", "-12h", "24h", "-24h", "1d", "-1d"}, false),
	},
}

var metricsMissingDataConditionSchema = map[string]*schema.Schema{
	"trigger_type": {
		Type:         schema.TypeString,
		Required:     true,
		ValidateFunc: validation.StringInSlice([]string{"MissingData", "ResolvedMissingData"}, false),
	},
	"time_range": {
		Type:         schema.TypeString,
		Required:     true,
		ValidateFunc: validation.StringInSlice([]string{"5m", "-5m", "10m", "-10m", "15m", "-15m", "30m", "-30m", "60m", "-60m", "1h", "-1h", "3h", "-3h", "6h", "-6h", "12h", "-12h", "24h", "-24h", "1d", "-1d"}, false),
	},
	"trigger_source": {
		Type:         schema.TypeString,
		Required:     true,
		ValidateFunc: validation.StringInSlice([]string{"AllTimeSeries", "AnyTimeSeries"}, false),
	},
}

/*
   A 'triggers' block is legacy if each one of the following fields is unset or set to an empty map
   - static_condition
   - logs_static_condition
   - metrics_static_condition
   - logs_outlier_condition
   - metrics_outlier_condition
   - logs_missing_data_condition
   - metrics_missing_data_condition
*/
func isLegacyTriggersBlock(block map[string]interface{}) bool {
	if _, ok := getSingletonArrayFieldOk(block, staticConditionFieldName); ok {
		return false
	}
	if _, ok := getSingletonArrayFieldOk(block, logsStaticConditionFieldName); ok {
		return false
	}
	if _, ok := getSingletonArrayFieldOk(block, metricsStaticConditionFieldName); ok {
		return false
	}
	if _, ok := getSingletonArrayFieldOk(block, logsOutlierConditionFieldName); ok {
		return false
	}
	if _, ok := getSingletonArrayFieldOk(block, metricsOutlierConditionFieldName); ok {
		return false
	}
	if _, ok := getSingletonArrayFieldOk(block, logsMissingDataConditionFieldName); ok {
		return false
	}
	if _, ok := getSingletonArrayFieldOk(block, metricsMissingDataConditionFieldName); ok {
		return false
	}
	return true
}
