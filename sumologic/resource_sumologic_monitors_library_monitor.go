package sumologic

import (
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
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
							Type:     schema.TypeString,
							Required: true,
						},
						"threshold": {
							Type:     schema.TypeFloat,
							Optional: true,
						},
						"threshold_type": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"time_range": {
							Type:     schema.TypeString,
							Required: true,
						},
						"trigger_source": {
							Type:     schema.TypeString,
							Required: true,
						},
						"occurrence_type": {
							Type:     schema.TypeString,
							Required: true,
						},
						"detection_method": {
							Type:     schema.TypeString,
							Optional: true,
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
										Type:     schema.TypeString,
										Required: true,
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
				Required: true,
			},

			"created_at": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"monitor_type": {
				Type:     schema.TypeString,
				Required: true,
			},

			"is_locked": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"status": {
				Type:     schema.TypeList,
				Optional: true,
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
				Type:     schema.TypeString,
				Optional: true,
				Default:  "MonitorsLibraryMonitor",
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
		paramMap := make(map[string]string)
		if monitor.ParentID == "" {
			rootFolder, err := c.GetMonitorsLibraryFolder("root")
			if err != nil {
				return err
			}

			monitor.ParentID = rootFolder.ID
		}
		paramMap["parentId"] = monitor.ParentID
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

	return nil
}

func resourceSumologicMonitorsLibraryMonitorUpdate(d *schema.ResourceData, meta interface{}) error {
	c := meta.(*Client)
	monitor := resourceToMonitorsLibraryMonitor(d)
	monitor.Type = "MonitorsLibraryMonitorUpdate"
	err := c.UpdateMonitorsLibraryMonitor(monitor)
	if err != nil {
		return err
	}
	return resourceSumologicMonitorsLibraryMonitorRead(d, meta)
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

func resourceToMonitorsLibraryMonitor(d *schema.ResourceData) MonitorsLibraryMonitor {
	// handle notifications
	rawNotifications := d.Get("notifications").([]interface{})
	notifications := make([]MonitorNotification, len(rawNotifications))
	for i := range rawNotifications {
		notificationDict := rawNotifications[i].(map[string]interface{})
		n := MonitorNotification{}
		rawNotificationAction := notificationDict["notification"].([]interface{})
		notificationActionDict := rawNotificationAction[0].(map[string]interface{})
		if notificationActionDict["action_type"].(string) == "EmailAction" {
			notificationAction := EmailNotification{}
			notificationAction.ActionType = notificationActionDict["action_type"].(string)
			notificationAction.Subject = notificationActionDict["subject"].(string)
			notificationAction.Recipients = notificationActionDict["recipients"].([]interface{})
			notificationAction.MessageBody = notificationActionDict["message_body"].(string)
			notificationAction.TimeZone = notificationActionDict["time_zone"].(string)
			n.Notification = notificationAction
		}
		if notificationActionDict["action_type"].(string) == "NamedConnectionAction" {
			notificationAction := WebhookNotificiation{}
			notificationAction.ActionType = notificationActionDict["action_type"].(string)
			notificationAction.ConnectionID = notificationActionDict["connection_id"].(string)
			notificationAction.PayloadOverride = notificationActionDict["payload_override"].(string)
			n.Notification = notificationAction
		}
		n.RunForTriggerTypes = notificationDict["run_for_trigger_types"].([]interface{})
		notifications[i] = n
	}
	// handle triggers
	rawTriggers := d.Get("triggers").([]interface{})
	triggers := make([]TriggerCondition, len(rawTriggers))
	for i := range rawTriggers {
		triggerDict := rawTriggers[i].(map[string]interface{})
		t := TriggerCondition{}
		t.TriggerType = triggerDict["trigger_type"].(string)
		t.Threshold = triggerDict["threshold"].(float64)
		t.ThresholdType = triggerDict["threshold_type"].(string)
		t.TimeRange = triggerDict["time_range"].(string)
		t.OccurrenceType = triggerDict["occurrence_type"].(string)
		t.TriggerSource = triggerDict["trigger_source"].(string)
		t.DetectionMethod = triggerDict["detection_method"].(string)
		triggers[i] = t
	}
	// handle queries
	rawQueries := d.Get("queries").([]interface{})
	queries := make([]MonitorQuery, len(rawQueries))
	for i := range rawQueries {
		queryDict := rawQueries[i].(map[string]interface{})
		q := MonitorQuery{}
		q.Query = queryDict["query"].(string)
		q.RowID = queryDict["row_id"].(string)
		queries[i] = q
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
		Status:             d.Get("status").(string),
		GroupNotifications: d.Get("group_notifications").(bool),
	}
}
