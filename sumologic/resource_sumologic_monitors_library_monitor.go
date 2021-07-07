package sumologic

import (
	"log"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
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
						"field": {
							Type:     schema.TypeString,
							Optional: true,
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
	triggers := make([]interface{}, len(monitor.Triggers))
	for i, t := range monitor.Triggers {
		triggers[i] = map[string]interface{}{
			"trigger_type":     t.TriggerType,
			"threshold":        t.Threshold,
			"threshold_type":   t.ThresholdType,
			"field":            t.Field,
			"time_range":       strings.TrimPrefix(t.TimeRange, "-"),
			"occurrence_type":  t.OccurrenceType,
			"trigger_source":   t.TriggerSource,
			"detection_method": t.DetectionMethod,
			"window":           t.Window,
			"consecutive":      t.Consecutive,
			"direction":        t.Direction,
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
		triggers[i] = TriggerCondition{
			TriggerType:     triggerDict["trigger_type"].(string),
			Threshold:       triggerDict["threshold"].(float64),
			ThresholdType:   triggerDict["threshold_type"].(string),
			Field:           triggerDict["field"].(string),
			TimeRange:       triggerDict["time_range"].(string),
			OccurrenceType:  triggerDict["occurrence_type"].(string),
			TriggerSource:   triggerDict["trigger_source"].(string),
			DetectionMethod: triggerDict["detection_method"].(string),
			Window:          triggerDict["window"].(int),
			Consecutive:     triggerDict["consecutive"].(int),
			Direction:       triggerDict["direction"].(string),
		}
	}
	return triggers
}

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
