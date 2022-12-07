package sumologic

import (
	"log"
	"regexp"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

var (
	triggerConditionsAtleatOneKey = []string{
		"trigger_conditions.0.logs_static_condition",
		"trigger_conditions.0.metrics_static_condition",
		"trigger_conditions.0.logs_outlier_condition",
		"trigger_conditions.0.metrics_outlier_condition",
		"trigger_conditions.0.logs_missing_data_condition",
		"trigger_conditions.0.metrics_missing_data_condition",
		"trigger_conditions.0.slo_sli_condition",
		"trigger_conditions.0.slo_burn_rate_condition",
	}
	criticalOrWarningAtleastOneKeys = []string{
		"trigger_conditions.0.logs_static_condition.0.warning",
		"trigger_conditions.0.logs_static_condition.0.critical",
	}
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
				Type:       schema.TypeList,
				Optional:   true,
				Deprecated: "The field `triggers` is deprecated and will be removed in a future release of the provider -- please use `trigger_conditions` instead.",
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
						"time_range": {
							Type:             schema.TypeString,
							Optional:         true,
							ValidateFunc:     validation.StringMatch(regexp.MustCompile(`^-?(\d)+[smhd]$`), "Time range must be in the format '-?\\d+[smhd]'. Examples: -15m, 1d, etc."),
							DiffSuppressFunc: SuppressEquivalentTimeDiff(false),
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
							Type:     schema.TypeString,
							Optional: true,
							ValidateFunc: validation.StringInSlice([]string{"StaticCondition", "LogsStaticCondition",
								"MetricsStaticCondition", "LogsOutlierCondition", "MetricsOutlierCondition",
								"LogsMissingDataCondition", "MetricsMissingDataCondition", "SloSliCondition",
								"SloBurnRateCondition"}, false),
						},
						"resolution_window": &resolutionWindowSchema,
					},
				},
			},

			"trigger_conditions": {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"logs_static_condition": {
							Type:     schema.TypeList,
							MaxItems: 1,
							Optional: true,
							Elem: &schema.Resource{
								Schema: logsStaticTriggerConditionSchema,
							},
							AtLeastOneOf: triggerConditionsAtleatOneKey,
						},
						"metrics_static_condition": {
							Type:     schema.TypeList,
							MaxItems: 1,
							Optional: true,
							Elem: &schema.Resource{
								Schema: metricsStaticTriggerConditionSchema,
							},
						},
						"logs_outlier_condition": {
							Type:     schema.TypeList,
							MaxItems: 1,
							Optional: true,
							Elem: &schema.Resource{
								Schema: logsOutlierTriggerConditionSchema,
							},
						},
						"metrics_outlier_condition": {
							Type:     schema.TypeList,
							MaxItems: 1,
							Optional: true,
							Elem: &schema.Resource{
								Schema: metricsOutlierTriggerConditionSchema,
							},
						},
						"logs_missing_data_condition": {
							Type:     schema.TypeList,
							MaxItems: 1,
							Optional: true,
							Elem: &schema.Resource{
								Schema: logsMissingDataTriggerConditionSchema,
							},
						},
						"metrics_missing_data_condition": {
							Type:     schema.TypeList,
							MaxItems: 1,
							Optional: true,
							Elem: &schema.Resource{
								Schema: metricsMissingDataTriggerConditionSchema,
							},
						},
						sloSLIConditionFieldName: {
							Type:     schema.TypeList,
							MaxItems: 1,
							Optional: true,
							Elem: &schema.Resource{
								Schema: sloSLITriggerConditionSchema,
							},
						},
						sloBurnRateConditionFieldName: {
							Type:     schema.TypeList,
							MaxItems: 1,
							Optional: true,
							Elem: &schema.Resource{
								Schema: sloBurnRateTriggerConditionSchema,
							},
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
										ValidateFunc: validation.StringInSlice([]string{"Email", "AWSLambda", "AzureFunctions", "Datadog", "HipChat", "Jira", "NewRelic", "Opsgenie", "PagerDuty", "Slack", "MicrosoftTeams", "ServiceNow", "Webhook"}, false),
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
										Type:         schema.TypeString,
										Optional:     true,
										ValidateFunc: validation.StringIsJSON,
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
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringMatch(regexp.MustCompile(`(?s)^[^\ ].*[^\ ]$`), "description must not contain leading or trailing spaces"),
			},

			"created_at": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"monitor_type": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringInSlice([]string{"Logs", "Metrics", "Slo"}, false),
			},

			"evaluation_delay": {
				Type:             schema.TypeString,
				Optional:         true,
				Computed:         true,
				ValidateFunc:     validation.StringMatch(regexp.MustCompile(`^((\d)+[smh])+$`), "This value is not in correct format. Example: 1m30s"),
				DiffSuppressFunc: SuppressEquivalentTimeDiff(false),
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
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringMatch(regexp.MustCompile(`(?s)^[^\ ].*[^\ ]$`), "name must not contain leading or trailing spaces"),
			},

			"post_request_map": {
				Type:     schema.TypeMap,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},

			"playbook": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"slo_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"alert_name": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringLenBetween(1, 512),
			},
			"notification_group_fields": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"obj_permission": GetCmfFgpObjPermSetSchema(),
		},
	}
}

// Trigger Condition schemas
var logsStaticTriggerConditionSchema = map[string]*schema.Schema{
	"field": {
		Type:     schema.TypeString,
		Optional: true,
	},
	"critical": nestedWithAtleastOneOfKeys(true, schemaMap{
		"time_range": &timeRangeSchema,
		"alert": nested(false, schemaMap{
			"threshold":      &thresholdSchema,
			"threshold_type": &thresholdTypeSchema,
		}),
		"resolution": nested(false, schemaMap{
			"threshold":         &thresholdSchema,
			"threshold_type":    &thresholdTypeSchema,
			"resolution_window": &resolutionWindowSchema,
		}),
	}, criticalOrWarningAtleastOneKeys),
	"warning": nestedWithAtleastOneOfKeys(true, schemaMap{
		"time_range": &timeRangeSchema,
		"alert": nested(false, schemaMap{
			"threshold":      &thresholdSchema,
			"threshold_type": &thresholdTypeSchema,
		}),
		"resolution": nested(false, schemaMap{
			"threshold":         &thresholdSchema,
			"threshold_type":    &thresholdTypeSchema,
			"resolution_window": &resolutionWindowSchema,
		}),
	}, criticalOrWarningAtleastOneKeys),
}

var metricsStaticTriggerConditionSchema = map[string]*schema.Schema{
	"critical": nestedWithAtleastOneOfKeys(true, schemaMap{
		"time_range":      &timeRangeSchema,
		"occurrence_type": &occurrenceTypeSchema,
		"alert": nested(false, schemaMap{
			"threshold":      &thresholdSchema,
			"threshold_type": &thresholdTypeSchema,
		}),
		"resolution": nested(false, schemaMap{
			"threshold":       &thresholdSchema,
			"threshold_type":  &thresholdTypeSchema,
			"occurrence_type": &occurrenceTypeOptSchema,
		}),
	}, criticalOrWarningAtleastOneKeys),
	"warning": nestedWithAtleastOneOfKeys(true, schemaMap{
		"time_range":      &timeRangeSchema,
		"occurrence_type": &occurrenceTypeSchema,
		"alert": nested(false, schemaMap{
			"threshold":      &thresholdSchema,
			"threshold_type": &thresholdTypeSchema,
		}),
		"resolution": nested(false, schemaMap{
			"threshold":       &thresholdSchema,
			"threshold_type":  &thresholdTypeSchema,
			"occurrence_type": &occurrenceTypeOptSchema,
		}),
	}, criticalOrWarningAtleastOneKeys),
}

var logsOutlierTriggerConditionSchema = map[string]*schema.Schema{
	"field": {
		Type:     schema.TypeString,
		Optional: true,
	},
	"direction": {
		Type:         schema.TypeString,
		Optional:     true,
		ValidateFunc: validation.StringInSlice([]string{"Both", "Up", "Down"}, false),
	},
	"critical": nestedWithAtleastOneOfKeys(true, schemaMap{
		"window":      &windowSchema,
		"consecutive": &consecutiveSchema,
		"threshold":   &thresholdSchema,
	}, criticalOrWarningAtleastOneKeys),
	"warning": nestedWithAtleastOneOfKeys(true, schemaMap{
		"window":      &windowSchema,
		"consecutive": &consecutiveSchema,
		"threshold":   &thresholdSchema,
	}, criticalOrWarningAtleastOneKeys),
}

var metricsOutlierTriggerConditionSchema = map[string]*schema.Schema{
	"direction": {
		Type:         schema.TypeString,
		Optional:     true,
		ValidateFunc: validation.StringInSlice([]string{"Both", "Up", "Down"}, false),
	},
	"critical": nestedWithAtleastOneOfKeys(true, schemaMap{
		"baseline_window": &baselineWindowSchema,
		"threshold":       &thresholdSchema,
	}, criticalOrWarningAtleastOneKeys),
	"warning": nestedWithAtleastOneOfKeys(true, schemaMap{
		"baseline_window": &baselineWindowSchema,
		"threshold":       &thresholdSchema,
	}, criticalOrWarningAtleastOneKeys),
}

var logsMissingDataTriggerConditionSchema = map[string]*schema.Schema{
	"time_range": &timeRangeSchema,
}

var metricsMissingDataTriggerConditionSchema = map[string]*schema.Schema{
	"time_range": &timeRangeSchema,
	"trigger_source": {
		Type:         schema.TypeString,
		Required:     true,
		ValidateFunc: validation.StringInSlice([]string{"AllTimeSeries", "AnyTimeSeries"}, false),
	},
}

var sloSLITriggerConditionSchema = map[string]*schema.Schema{
	"critical": nestedWithAtleastOneOfKeys(true, schemaMap{
		"sli_threshold": {
			Type:         schema.TypeFloat,
			Required:     true,
			ValidateFunc: validation.FloatBetween(0, 100),
		},
	}, criticalOrWarningAtleastOneKeys),
	"warning": nestedWithAtleastOneOfKeys(true, schemaMap{
		"sli_threshold": {
			Type:         schema.TypeFloat,
			Required:     true,
			ValidateFunc: validation.FloatBetween(0, 100),
		},
	}, criticalOrWarningAtleastOneKeys),
}

var sloBurnRateTriggerConditionSchema = map[string]*schema.Schema{
	"critical": nestedWithAtleastOneOfKeys(true, schemaMap{
		"time_range": &timeRangeSchema,
		"burn_rate_threshold": {
			Type:         schema.TypeFloat,
			Required:     true,
			ValidateFunc: validation.FloatAtLeast(0),
		},
	}, criticalOrWarningAtleastOneKeys),
	"warning": nestedWithAtleastOneOfKeys(true, schemaMap{
		"time_range": &timeRangeSchema,
		"burn_rate_threshold": {
			Type:         schema.TypeFloat,
			Required:     true,
			ValidateFunc: validation.FloatAtLeast(0),
		},
	}, criticalOrWarningAtleastOneKeys),
}

var occurrenceTypeSchema = schema.Schema{
	Type:         schema.TypeString,
	Required:     true,
	ValidateFunc: validation.StringInSlice([]string{"AtLeastOnce", "Always"}, false),
}

var occurrenceTypeOptSchema = schema.Schema{
	Type:         schema.TypeString,
	Optional:     true,
	ValidateFunc: validation.StringInSlice([]string{"AtLeastOnce", "Always"}, false),
}

var windowSchema = schema.Schema{
	Type:         schema.TypeInt,
	Optional:     true,
	ValidateFunc: validation.IntAtLeast(1),
}

var baselineWindowSchema = schema.Schema{
	Type:     schema.TypeString,
	Optional: true,
}

var consecutiveSchema = schema.Schema{
	Type:         schema.TypeInt,
	Optional:     true,
	ValidateFunc: validation.IntAtLeast(1),
}

var timeRangeSchema = schema.Schema{
	Type:             schema.TypeString,
	Required:         true,
	ValidateFunc:     validation.StringMatch(regexp.MustCompile(`^-?(\d)+[smhd]$`), "Time range must be in the format '-?\\d+[smhd]'. Examples: -15m, 1d, etc."),
	DiffSuppressFunc: SuppressEquivalentTimeDiff(false),
}

var resolutionWindowSchema = schema.Schema{
	Type:             schema.TypeString,
	Optional:         true,
	ValidateFunc:     validation.StringMatch(regexp.MustCompile(`^(\d)+[smhd]`), "Resolution window must be in the format '\\d+[smhd]'. Examples: 0m, 15m, 1d, etc."),
	DiffSuppressFunc: SuppressEquivalentTimeDiff(false),
}

var thresholdSchema = schema.Schema{
	Type:     schema.TypeFloat,
	Optional: true,
}

var thresholdTypeSchema = schema.Schema{
	Type:         schema.TypeString,
	Optional:     true,
	ValidateFunc: validation.StringInSlice([]string{"LessThan", "LessThanOrEqual", "GreaterThan", "GreaterThanOrEqual"}, false),
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

		permStmts, convErr := ResourceToCmfFgpPermStmts(d, monitorDefinitionID)
		if convErr != nil {
			return convErr
		}
		_, fgpErr := c.SetCmfFgp(fgpTargetType, CmfFgpRequest{
			PermissionStatements: permStmts,
		})
		if fgpErr != nil {
			return fgpErr
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

	fgpResponse, fgpErr := c.GetCmfFgp(fgpTargetType, monitor.ID)
	if fgpErr != nil {
		suppressedErrorCode := HasErrorCode(fgpErr.Error(), []string{"not_implemented_yet", "api_not_enabled"})
		if suppressedErrorCode == "" {
			return fgpErr
		} else {
			log.Printf("[WARN] FGP Feature has not been enabled yet. Suppressing \"%s\" error under GetCmfFgp operation.", suppressedErrorCode)
		}
	} else {
		CmfFgpPermStmtsSetToResource(d, fgpResponse.PermissionStatements)
	}

	d.Set("created_by", monitor.CreatedBy)
	d.Set("created_at", monitor.CreatedAt)
	d.Set("monitor_type", monitor.MonitorType)
	d.Set("evaluation_delay", monitor.EvaluationDelay)
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
	d.Set("playbook", monitor.Playbook)
	d.Set("alert_name", monitor.AlertName)
	d.Set("slo_id", monitor.SloID)
	d.Set("notification_group_fields", monitor.NotificationGroupFields)

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

	// set either 'trigger_conditions' or 'triggers', but not both, based on whichever the plan uses.
	// we avoid converting between the two so as to prevent plan mismatches before and after an apply.
	var has_trigger_conditions = false
	if val, ok := d.GetOk("trigger_conditions"); ok {
		if arr, ok := val.([]interface{}); ok && len(arr) > 0 {
			has_trigger_conditions = true
			if err :=
				d.Set("trigger_conditions", toSingletonArray(jsonToTriggerConditionsBlock(monitor.Triggers))); err != nil {
				return err
			}
		}
	}
	if !has_trigger_conditions {
		triggers := make([]interface{}, len(monitor.Triggers))
		for i, t := range monitor.Triggers {
			triggers[i] = map[string]interface{}{
				"time_range":        t.PositiveTimeRange(),
				"trigger_type":      t.TriggerType,
				"threshold":         t.Threshold,
				"threshold_type":    t.ThresholdType,
				"occurrence_type":   t.OccurrenceType,
				"trigger_source":    t.TriggerSource,
				"detection_method":  t.DetectionMethod,
				"resolution_window": t.PositiveResolutionWindow(),
			}
		}
		if err := d.Set("triggers", triggers); err != nil {
			return err
		}
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

	if d.HasChange("parent_id") {
		updatedMonitor, err := c.MoveMonitorsLibraryMonitor(monitor.ID, monitor.ParentID)
		if err != nil {
			return err
		}
		monitor = *updatedMonitor
	}
	monitor.Type = "MonitorsLibraryMonitorUpdate"
	err := c.UpdateMonitorsLibraryMonitor(monitor)
	if err != nil {
		return err
	}

	// converting Resource FGP to Struct
	permStmts, convErr := ResourceToCmfFgpPermStmts(d, monitor.ID)
	if convErr != nil {
		return convErr
	}

	// reading FGP from Backend to reconcile
	fgpGetResponse, fgpGetErr := c.GetCmfFgp(fgpTargetType, monitor.ID)
	if fgpGetErr != nil {
		/*
		   |errCode         |  len  | logic                   |
		   |--------------------------------------------------|
		   |server_error    |   0   | return err at Get       |
		   |server_error    |   1   | warn; return err at Set |
		   |not_enabled     |   0   | warn                    |
		   |not_enabled     |   1   | warn; return err at Set |
		*/
		suppressedErrorCode := HasErrorCode(fgpGetErr.Error(), []string{"not_implemented_yet", "api_not_enabled"})
		if suppressedErrorCode == "" && len(permStmts) == 0 {
			return fgpGetErr
		} else {
			log.Printf("[WARN] FGP Feature has not been enabled yet. Suppressing \"%s\" error under GetCmfFgp operation.", suppressedErrorCode)
		}
	}

	if len(permStmts) > 0 || fgpGetResponse != nil {
		_, fgpSetErr := c.SetCmfFgp(fgpTargetType, CmfFgpRequest{
			PermissionStatements: ReconcileFgpPermStmtsWithEmptyPerms(
				permStmts, fgpGetResponse.PermissionStatements,
			),
		})
		if fgpSetErr != nil {
			return fgpSetErr
		}
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
	if triggerCondition, ok := singletonFromResourceData(d, "trigger_conditions"); ok {
		ret := triggerConditionsBlockToJson(triggerCondition)
		return ret
	} else {
		rawTriggers := d.Get("triggers").([]interface{})
		triggers := make([]TriggerCondition, len(rawTriggers))
		for i := range rawTriggers {
			triggerDict := rawTriggers[i].(map[string]interface{})
			triggers[i] = TriggerCondition{
				TriggerType:      triggerDict["trigger_type"].(string),
				Threshold:        triggerDict["threshold"].(float64),
				ThresholdType:    triggerDict["threshold_type"].(string),
				TimeRange:        triggerDict["time_range"].(string),
				OccurrenceType:   triggerDict["occurrence_type"].(string),
				TriggerSource:    triggerDict["trigger_source"].(string),
				DetectionMethod:  triggerDict["detection_method"].(string),
				ResolutionWindow: triggerDict["resolution_window"].(string),
			}
		}
		return triggers
	}
}

// 'trigger_conditions' block to TriggerCondition JSON structure
func triggerConditionsBlockToJson(block map[string]interface{}) []TriggerCondition {
	conditions := make([]TriggerCondition, 0)
	if sc, ok := fromSingletonArray(block, logsStaticConditionFieldName); ok {
		conditions = append(conditions, logsStaticConditionBlockToJson(sc)...)
	}
	if sc, ok := fromSingletonArray(block, metricsStaticConditionFieldName); ok {
		conditions = append(conditions, metricsStaticConditionBlockToJson(sc)...)
	}
	if sc, ok := fromSingletonArray(block, logsOutlierConditionFieldName); ok {
		conditions = append(conditions, logsOutlierConditionBlockToJson(sc)...)
	}
	if sc, ok := fromSingletonArray(block, metricsOutlierConditionFieldName); ok {
		conditions = append(conditions, metricsOutlierConditionBlockToJson(sc)...)
	}
	if sc, ok := fromSingletonArray(block, logsMissingDataConditionFieldName); ok {
		conditions = append(conditions, logsMissingDataConditionBlockToJson(sc)...)
	}
	if sc, ok := fromSingletonArray(block, metricsMissingDataConditionFieldName); ok {
		conditions = append(conditions, metricsMissingDataConditionBlockToJson(sc)...)
	}
	if sc, ok := fromSingletonArray(block, sloSLIConditionFieldName); ok {
		conditions = append(conditions, sloSLIConditionBlockToJson(sc)...)
	}
	if sc, ok := fromSingletonArray(block, sloBurnRateConditionFieldName); ok {
		conditions = append(conditions, sloBurnConditionBlockToJson(sc)...)
	}

	return conditions
}

func logsStaticConditionBlockToJson(block map[string]interface{}) []TriggerCondition {
	base := TriggerCondition{
		Field:           block["field"].(string),
		DetectionMethod: logsStaticConditionDetectionMethod,
	}
	return base.cloneReadingFromNestedBlocks(block)
}

func metricsStaticConditionBlockToJson(block map[string]interface{}) []TriggerCondition {
	base := TriggerCondition{
		DetectionMethod: metricsStaticConditionDetectionMethod,
	}
	triggerConditions := base.cloneReadingFromNestedBlocks(block)
	for i, _ := range triggerConditions {
		if (triggerConditions[i].TriggerType == "ResolvedCritical" && triggerConditions[i].OccurrenceType == "") ||
			(triggerConditions[i].TriggerType == "ResolvedWarning" && triggerConditions[i].OccurrenceType == "") {
			triggerConditions[i].OccurrenceType = "Always"
		}
	}
	return triggerConditions
}

func logsOutlierConditionBlockToJson(block map[string]interface{}) []TriggerCondition {
	base := TriggerCondition{
		Field:           block["field"].(string),
		Direction:       block["direction"].(string),
		DetectionMethod: logsOutlierConditionDetectionMethod,
	}
	// outlier blocks do not have 'alert' and 'resolution' sub-blocks. Here we generate empty blocks
	// for reading to work
	for _, triggerType := range []string{"critical", "warning"} {
		if subBlock, ok := fromSingletonArray(block, triggerType); ok {
			subBlock["alert"] = toSingletonArray(map[string]interface{}{})
			subBlock["resolution"] = toSingletonArray(map[string]interface{}{})
		}
	}
	return base.cloneReadingFromNestedBlocks(block)
}

func metricsOutlierConditionBlockToJson(block map[string]interface{}) []TriggerCondition {
	base := TriggerCondition{
		Direction:       block["direction"].(string),
		DetectionMethod: metricsOutlierConditionDetectionMethod,
	}
	// outlier blocks do not have 'alert' and 'resolution' sub-blocks. Here we generate empty blocks
	// for reading to work
	for _, triggerType := range []string{"critical", "warning"} {
		if subBlock, ok := fromSingletonArray(block, triggerType); ok {
			subBlock["alert"] = toSingletonArray(map[string]interface{}{})
			subBlock["resolution"] = toSingletonArray(map[string]interface{}{})
		}
	}
	return base.cloneReadingFromNestedBlocks(block)
}

func logsMissingDataConditionBlockToJson(block map[string]interface{}) []TriggerCondition {
	alert := TriggerCondition{
		TimeRange:       block["time_range"].(string),
		DetectionMethod: logsMissingDataConditionDetectionMethod,
		TriggerType:     "MissingData",
	}
	resolution := TriggerCondition{
		TimeRange:       block["time_range"].(string),
		DetectionMethod: logsMissingDataConditionDetectionMethod,
		TriggerType:     "ResolvedMissingData",
	}
	return []TriggerCondition{alert, resolution}
}

func metricsMissingDataConditionBlockToJson(block map[string]interface{}) []TriggerCondition {
	// The TF model for metrics missing data does not have explicit threshold blocks. We implicitly create
	// the 2 trigger types: MissingData and ResolvedMissingData.
	alert := TriggerCondition{
		TimeRange:       block["time_range"].(string),
		TriggerSource:   block["trigger_source"].(string),
		DetectionMethod: metricsMissingDataConditionDetectionMethod,
		TriggerType:     "MissingData",
	}

	resolution := TriggerCondition{
		TimeRange:       block["time_range"].(string),
		TriggerSource:   block["trigger_source"].(string),
		DetectionMethod: metricsMissingDataConditionDetectionMethod,
		TriggerType:     "ResolvedMissingData",
	}
	return []TriggerCondition{alert, resolution}
}

func sloSLIConditionBlockToJson(block map[string]interface{}) []TriggerCondition {
	base := TriggerCondition{
		DetectionMethod: sloSLIConditionDetectionMethod,
	}

	return base.sloCloneReadingFromNestedBlocks(block)
}

func sloBurnConditionBlockToJson(block map[string]interface{}) []TriggerCondition {
	base := TriggerCondition{
		DetectionMethod: sloBurnRateConditionDetectionMethod,
	}

	return base.sloCloneReadingFromNestedBlocks(block)
}

// TriggerCondition JSON model to 'trigger_conditions' block
func jsonToTriggerConditionsBlock(conditions []TriggerCondition) map[string]interface{} {
	missingDataConditions := make([]TriggerCondition, 0)
	dataConditions := make([]TriggerCondition, 0)
	for _, condition := range conditions {
		if condition.TriggerType == "MissingData" || condition.TriggerType == "ResolvedMissingData" {
			missingDataConditions = append(missingDataConditions, condition)
		} else {
			dataConditions = append(dataConditions, condition)
		}
	}
	triggerConditionsBlock := map[string]interface{}{}
	if len(dataConditions) > 0 {
		switch dataConditions[0].DetectionMethod {
		case logsStaticConditionDetectionMethod:
			triggerConditionsBlock[logsStaticConditionFieldName] = toSingletonArray(jsonToLogsStaticConditionBlock(dataConditions))
		case metricsStaticConditionDetectionMethod:
			triggerConditionsBlock[metricsStaticConditionFieldName] = toSingletonArray(jsonToMetricsStaticConditionBlock(dataConditions))
		case logsOutlierConditionDetectionMethod:
			triggerConditionsBlock[logsOutlierConditionFieldName] = toSingletonArray(jsonToLogsOutlierConditionBlock(dataConditions))
		case metricsOutlierConditionDetectionMethod:
			triggerConditionsBlock[metricsOutlierConditionFieldName] = toSingletonArray(jsonToMetricsOutlierConditionBlock(dataConditions))
		case sloSLIConditionDetectionMethod:
			triggerConditionsBlock[sloSLIConditionFieldName] = toSingletonArray(jsonToSloSliConditionBlock(dataConditions))
		case sloBurnRateConditionDetectionMethod:
			triggerConditionsBlock[sloBurnRateConditionFieldName] = toSingletonArray(jsonToSloBurnRateConditionBlock(dataConditions))
		}
	}
	if len(missingDataConditions) > 0 {
		switch missingDataConditions[0].DetectionMethod {
		case logsMissingDataConditionDetectionMethod:
			triggerConditionsBlock[logsMissingDataConditionFieldName] = toSingletonArray(jsonToLogsMissingDataConditionBlock(missingDataConditions))
		case metricsMissingDataConditionDetectionMethod:
			triggerConditionsBlock[metricsMissingDataConditionFieldName] = toSingletonArray(jsonToMetricsMissingDataConditionBlock(missingDataConditions))
		}
	}
	return triggerConditionsBlock
}

func jsonToLogsStaticConditionBlock(conditions []TriggerCondition) map[string]interface{} {
	var criticalAlrt, criticalRslv = dict{}, dict{}
	var warningAlrt, warningRslv = dict{}, dict{}
	criticalDict := dict{"alert": toSingletonArray(criticalAlrt), "resolution": toSingletonArray(criticalRslv)}
	warningDict := dict{"alert": toSingletonArray(warningAlrt), "resolution": toSingletonArray(warningRslv)}
	block := map[string]interface{}{}

	block["field"] = conditions[0].Field
	block["critical"] = toSingletonArray(criticalDict)
	block["warning"] = toSingletonArray(warningDict)

	var hasCritical, hasWarning = false, false
	for _, condition := range conditions {
		switch condition.TriggerType {
		case "Critical":
			hasCritical = true
			criticalDict["time_range"] = condition.PositiveTimeRange()
			criticalAlrt["threshold"] = condition.Threshold
			criticalAlrt["threshold_type"] = condition.ThresholdType
		case "ResolvedCritical":
			hasCritical = true
			criticalDict["time_range"] = condition.PositiveTimeRange()
			criticalRslv["threshold"] = condition.Threshold
			criticalRslv["threshold_type"] = condition.ThresholdType
			criticalRslv["resolution_window"] = condition.PositiveResolutionWindow()
		case "Warning":
			hasWarning = true
			warningDict["time_range"] = condition.PositiveTimeRange()
			warningAlrt["threshold"] = condition.Threshold
			warningAlrt["threshold_type"] = condition.ThresholdType
		case "ResolvedWarning":
			hasWarning = true
			warningDict["time_range"] = condition.PositiveTimeRange()
			warningRslv["threshold"] = condition.Threshold
			warningRslv["threshold_type"] = condition.ThresholdType
			warningRslv["resolution_window"] = condition.PositiveResolutionWindow()
		}
	}
	if !hasCritical {
		delete(block, "critical")
	}
	if !hasWarning {
		delete(block, "warning")
	}
	return block
}

func jsonToMetricsStaticConditionBlock(conditions []TriggerCondition) map[string]interface{} {
	var criticalAlrt, criticalRslv = dict{}, dict{}
	var warningAlrt, warningRslv = dict{}, dict{}
	criticalDict := dict{"alert": toSingletonArray(criticalAlrt), "resolution": toSingletonArray(criticalRslv)}
	warningDict := dict{"alert": toSingletonArray(warningAlrt), "resolution": toSingletonArray(warningRslv)}
	block := map[string]interface{}{}

	block["critical"] = toSingletonArray(criticalDict)
	block["warning"] = toSingletonArray(warningDict)

	var hasCritical, hasWarning = false, false
	for _, condition := range conditions {
		switch condition.TriggerType {
		case "Critical":
			hasCritical = true
			criticalDict["time_range"] = condition.PositiveTimeRange()
			criticalDict["occurrence_type"] = condition.OccurrenceType
			criticalAlrt["threshold"] = condition.Threshold
			criticalAlrt["threshold_type"] = condition.ThresholdType
		case "ResolvedCritical":
			hasCritical = true
			criticalDict["time_range"] = condition.PositiveTimeRange()
			criticalRslv["threshold"] = condition.Threshold
			criticalRslv["threshold_type"] = condition.ThresholdType
			if condition.OccurrenceType == "AtLeastOnce" {
				criticalRslv["occurrence_type"] = condition.OccurrenceType
			} else {
				// otherwise, the canonical translation is to leave out occurrenceType in the Resolved block
			}
		case "Warning":
			hasWarning = true
			warningDict["time_range"] = condition.PositiveTimeRange()
			warningDict["occurrence_type"] = condition.OccurrenceType
			warningAlrt["threshold"] = condition.Threshold
			warningAlrt["threshold_type"] = condition.ThresholdType
		case "ResolvedWarning":
			hasWarning = true
			warningDict["time_range"] = condition.PositiveTimeRange()
			warningRslv["threshold"] = condition.Threshold
			warningRslv["threshold_type"] = condition.ThresholdType
			if condition.OccurrenceType == "AtLeastOnce" {
				criticalRslv["occurrence_type"] = condition.OccurrenceType
			} else {
				// otherwise, the canonical translation is to leave out occurrenceType in the Resolved block
			}
		}
	}
	if !hasCritical {
		delete(block, "critical")
	}
	if !hasWarning {
		delete(block, "warning")
	}
	return block
}

func jsonToLogsOutlierConditionBlock(conditions []TriggerCondition) map[string]interface{} {
	var criticalDict, warningDict = dict{}, dict{}
	block := map[string]interface{}{}

	block["field"] = conditions[0].Field
	block["direction"] = conditions[0].Direction
	block["critical"] = toSingletonArray(criticalDict)
	block["warning"] = toSingletonArray(warningDict)

	var hasCritical, hasWarning = false, false
	for _, condition := range conditions {
		switch condition.TriggerType {
		case "Critical":
			hasCritical = true
			criticalDict["window"] = condition.Window
			criticalDict["consecutive"] = condition.Consecutive
			criticalDict["threshold"] = condition.Threshold
		case "ResolvedCritical":
			hasCritical = true
			criticalDict["window"] = condition.Window
			criticalDict["consecutive"] = condition.Consecutive
			criticalDict["threshold"] = condition.Threshold
		case "Warning":
			hasWarning = true
			warningDict["window"] = condition.Window
			warningDict["consecutive"] = condition.Consecutive
			warningDict["threshold"] = condition.Threshold
		case "ResolvedWarning":
			hasWarning = true
			warningDict["window"] = condition.Window
			warningDict["consecutive"] = condition.Consecutive
			warningDict["threshold"] = condition.Threshold
		}
	}
	if !hasCritical {
		delete(block, "critical")
	}
	if !hasWarning {
		delete(block, "warning")
	}
	return block
}

func jsonToMetricsOutlierConditionBlock(conditions []TriggerCondition) map[string]interface{} {
	var criticalDict, warningDict = dict{}, dict{}
	block := map[string]interface{}{}

	block["direction"] = conditions[0].Direction
	block["critical"] = toSingletonArray(criticalDict)
	block["warning"] = toSingletonArray(warningDict)

	var hasCritical, hasWarning = false, false
	for _, condition := range conditions {
		switch condition.TriggerType {
		case "Critical":
			hasCritical = true
			criticalDict["baseline_window"] = condition.PositiveBaselineWindow()
			criticalDict["threshold"] = condition.Threshold
		case "ResolvedCritical":
			hasCritical = true
			criticalDict["baseline_window"] = condition.PositiveBaselineWindow()
			criticalDict["threshold"] = condition.Threshold
		case "Warning":
			hasWarning = true
			warningDict["baseline_window"] = condition.PositiveBaselineWindow()
			warningDict["threshold"] = condition.Threshold
		case "ResolvedWarning":
			hasWarning = true
			warningDict["baseline_window"] = condition.PositiveBaselineWindow()
			warningDict["threshold"] = condition.Threshold
		}
	}
	if !hasCritical {
		delete(block, "critical")
	}
	if !hasWarning {
		delete(block, "warning")
	}
	return block
}

func jsonToSloSliConditionBlock(conditions []TriggerCondition) map[string]interface{} {
	var criticalAlrt, warningAlrt = dict{}, dict{}
	block := map[string]interface{}{}

	block["critical"] = toSingletonArray(criticalAlrt)
	block["warning"] = toSingletonArray(warningAlrt)

	var hasCritical, hasWarning = false, false
	for _, condition := range conditions {
		switch condition.TriggerType {
		case "Critical":
			hasCritical = true
			criticalAlrt["sli_threshold"] = condition.SLIThreshold
		case "Warning":
			hasWarning = true
			warningAlrt["sli_threshold"] = condition.SLIThreshold
		}
	}
	if !hasCritical {
		delete(block, "critical")
	}
	if !hasWarning {
		delete(block, "warning")
	}
	return block
}

func jsonToSloBurnRateConditionBlock(conditions []TriggerCondition) map[string]interface{} {

	var criticalAlrt, warningAlrt = dict{}, dict{}
	block := map[string]interface{}{}

	block["critical"] = toSingletonArray(criticalAlrt)
	block["warning"] = toSingletonArray(warningAlrt)

	var hasCritical, hasWarning = false, false
	for _, condition := range conditions {
		switch condition.TriggerType {
		case "Critical":
			hasCritical = true
			criticalAlrt["time_range"] = condition.TimeRange
			criticalAlrt["burn_rate_threshold"] = condition.BurnRateThreshold
		case "Warning":
			hasWarning = true
			warningAlrt["time_range"] = condition.TimeRange
			warningAlrt["burn_rate_threshold"] = condition.BurnRateThreshold
		}
	}
	if !hasCritical {
		delete(block, "critical")
	}
	if !hasWarning {
		delete(block, "warning")
	}
	return block
}

func jsonToLogsMissingDataConditionBlock(conditions []TriggerCondition) map[string]interface{} {
	block := map[string]interface{}{}
	firstCondition := conditions[0]
	block["time_range"] = firstCondition.PositiveTimeRange()
	return block
}

func jsonToMetricsMissingDataConditionBlock(conditions []TriggerCondition) map[string]interface{} {
	block := map[string]interface{}{}
	firstCondition := conditions[0]
	block["time_range"] = firstCondition.PositiveTimeRange()
	block["trigger_source"] = firstCondition.TriggerSource
	return block
}

const logsStaticConditionFieldName = "logs_static_condition"
const metricsStaticConditionFieldName = "metrics_static_condition"
const logsOutlierConditionFieldName = "logs_outlier_condition"
const metricsOutlierConditionFieldName = "metrics_outlier_condition"
const logsMissingDataConditionFieldName = "logs_missing_data_condition"
const metricsMissingDataConditionFieldName = "metrics_missing_data_condition"
const sloSLIConditionFieldName = "slo_sli_condition"
const sloBurnRateConditionFieldName = "slo_burn_rate_condition"

const logsStaticConditionDetectionMethod = "LogsStaticCondition"
const metricsStaticConditionDetectionMethod = "MetricsStaticCondition"
const logsOutlierConditionDetectionMethod = "LogsOutlierCondition"
const metricsOutlierConditionDetectionMethod = "MetricsOutlierCondition"
const logsMissingDataConditionDetectionMethod = "LogsMissingDataCondition"
const metricsMissingDataConditionDetectionMethod = "MetricsMissingDataCondition"
const sloSLIConditionDetectionMethod = "SloSliCondition"
const sloBurnRateConditionDetectionMethod = "SloBurnRateCondition"

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
	rawGroupFields := d.Get("notification_group_fields").([]interface{})
	notificationGroupFields := make([]string, len(rawGroupFields))
	for i := range rawGroupFields {
		notificationGroupFields[i] = rawGroupFields[i].(string)
	}

	return MonitorsLibraryMonitor{
		CreatedBy:               d.Get("created_by").(string),
		Name:                    d.Get("name").(string),
		ID:                      d.Id(),
		CreatedAt:               d.Get("created_at").(string),
		MonitorType:             d.Get("monitor_type").(string),
		Description:             d.Get("description").(string),
		EvaluationDelay:         d.Get("evaluation_delay").(string),
		Queries:                 queries,
		ModifiedBy:              d.Get("modified_by").(string),
		IsMutable:               d.Get("is_mutable").(bool),
		Version:                 d.Get("version").(int),
		Notifications:           notifications,
		Type:                    d.Get("type").(string),
		ParentID:                d.Get("parent_id").(string),
		ModifiedAt:              d.Get("modified_at").(string),
		Triggers:                triggers,
		ContentType:             d.Get("content_type").(string),
		IsLocked:                d.Get("is_locked").(bool),
		IsSystem:                d.Get("is_system").(bool),
		IsDisabled:              d.Get("is_disabled").(bool),
		Status:                  status,
		GroupNotifications:      d.Get("group_notifications").(bool),
		Playbook:                d.Get("playbook").(string),
		AlertName:               d.Get("alert_name").(string),
		SloID:                   d.Get("slo_id").(string),
		NotificationGroupFields: notificationGroupFields,
	}
}

func toResource(schm ...map[string]*schema.Schema) *schema.Resource {
	mapMerge := func(dest map[string]*schema.Schema, src map[string]*schema.Schema) {
		for k, v := range src {
			dest[k] = v
		}
	}
	allSchms := map[string]*schema.Schema{}
	for _, _schm := range schm {
		mapMerge(allSchms, _schm)
	}
	return &schema.Resource{
		Schema: allSchms,
	}
}

func nested(optional bool, sch map[string]*schema.Schema) *schema.Schema {
	if optional {
		return &schema.Schema{
			Type:     schema.TypeList,
			Optional: true,
			MaxItems: 1,
			Elem:     toResource(sch),
		}
	} else {
		return &schema.Schema{
			Type:     schema.TypeList,
			Required: true,
			MaxItems: 1,
			Elem:     toResource(sch),
		}
	}
}

func nestedWithAtleastOneOfKeys(optional bool, sch map[string]*schema.Schema, atleastOneOfKeys []string) *schema.Schema {
	if optional {
		return &schema.Schema{
			Type:         schema.TypeList,
			Optional:     true,
			MaxItems:     1,
			Elem:         toResource(sch),
			AtLeastOneOf: atleastOneOfKeys,
		}
	} else {
		return &schema.Schema{
			Type:         schema.TypeList,
			Required:     true,
			MaxItems:     1,
			Elem:         toResource(sch),
			AtLeastOneOf: atleastOneOfKeys,
		}
	}
}

func (condition *TriggerCondition) readFrom(block map[string]interface{}) {
	for k, v := range block {
		switch k {
		case "time_range":
			condition.TimeRange = v.(string)
		case "trigger_type":
			condition.TriggerType = v.(string)
		case "threshold":
			condition.Threshold = v.(float64)
		case "threshold_type":
			condition.ThresholdType = v.(string)
		case "occurrence_type":
			condition.OccurrenceType = v.(string)
		case "trigger_source":
			condition.TriggerSource = v.(string)
		case "detection_method":
			condition.DetectionMethod = v.(string)
		case "field":
			condition.Field = v.(string)
		case "window":
			condition.Window = v.(int)
		case "baseline_window":
			condition.BaselineWindow = v.(string)
		case "consecutive":
			condition.Consecutive = v.(int)
		case "direction":
			condition.Direction = v.(string)
		case "sli_threshold":
			condition.SLIThreshold = v.(float64)
		case "burn_rate_threshold":
			condition.BurnRateThreshold = v.(float64)
		case "resolution_window":
			condition.ResolutionWindow = v.(string)
		default:
		}
	}
}

// Clones the TriggerCondition with fields from nested blocks.
// Expects the following nested block structure:
// +-critical
//
//	+-alert
//	+-resolution
//
// +-warning
//
//	+-alert
//	+-resolution
//
// Adds any 'flat' fields appear at a level to the trigger condition, as per the mapping defined in [[readFrom]].
func (base TriggerCondition) cloneReadingFromNestedBlocks(block map[string]interface{}) []TriggerCondition {
	var conditions = []TriggerCondition{}
	var criticalCondition, resolvedCriticalCondition, warningCondition, resolvedWarningCondition = base, base, base, base
	criticalCondition.TriggerType = "Critical"
	resolvedCriticalCondition.TriggerType = "ResolvedCritical"
	warningCondition.TriggerType = "Warning"
	resolvedWarningCondition.TriggerType = "ResolvedWarning"

	if critical, ok := fromSingletonArray(block, "critical"); ok {
		criticalCondition.readFrom(critical)
		resolvedCriticalCondition.readFrom(critical)
		if resolvedCriticalCondition.DetectionMethod == metricsStaticConditionDetectionMethod {
			// do not inherit the top-level occurrence type into resolution blocks for MetricsStaticConditions
			// we want the caller to be able to tell whether the resolution block had set its own occurrence type
			resolvedCriticalCondition.OccurrenceType = ""
		}
		if alert, ok := fromSingletonArray(critical, "alert"); ok {
			criticalCondition.readFrom(alert)
		}
		if resolved, ok := fromSingletonArray(critical, "resolution"); ok {
			resolvedCriticalCondition.readFrom(resolved)
		}
		conditions = append(conditions, criticalCondition, resolvedCriticalCondition)
	}
	if warning, ok := fromSingletonArray(block, "warning"); ok {
		warningCondition.readFrom(warning)
		resolvedWarningCondition.readFrom(warning)
		if resolvedCriticalCondition.DetectionMethod == metricsStaticConditionDetectionMethod {
			// do not inherit the top-level occurrence type into resolution blocks for MetricsStaticConditions
			// we want the caller to be able to tell whether the resolution block had set its own occurrence type
			resolvedCriticalCondition.OccurrenceType = ""
		}
		if alert, ok := fromSingletonArray(warning, "alert"); ok {
			warningCondition.readFrom(alert)
		}
		if resolved, ok := fromSingletonArray(warning, "resolution"); ok {
			resolvedWarningCondition.readFrom(resolved)
		}
		conditions = append(conditions, warningCondition, resolvedWarningCondition)
	}
	return conditions
}

// adapted version of cloneReadingFromNestedBlocks for slo conditions
func (base TriggerCondition) sloCloneReadingFromNestedBlocks(block map[string]interface{}) []TriggerCondition {
	var conditions = []TriggerCondition{}
	var criticalCondition, warningCondition = base, base
	criticalCondition.TriggerType = "Critical"
	warningCondition.TriggerType = "Warning"

	if critical, ok := fromSingletonArray(block, "critical"); ok {
		criticalCondition.readFrom(critical)
		conditions = append(conditions, criticalCondition)
	}
	if warning, ok := fromSingletonArray(block, "warning"); ok {
		warningCondition.readFrom(warning)
		conditions = append(conditions, warningCondition)
	}
	return conditions
}

func toSingletonArray(m map[string]interface{}) []map[string]interface{} {
	return []map[string]interface{}{m}
}

func fromSingletonArray(block map[string]interface{}, field string) (map[string]interface{}, bool) {
	emptyDict := dict{}
	if iface, ok := block[field]; ok {
		switch arr := iface.(type) {
		case []map[string]interface{}:
			switch len(arr) {
			case 0:
				return emptyDict, false
			case 1:
				return arr[0], true
			default:
				log.Fatalf("Expected field '%s' to be a singleton array if present, got: %s", field, iface)
			}
		case []interface{}:
			switch len(arr) {
			case 0:
				return emptyDict, false
			case 1:
				return arr[0].(map[string]interface{}), true
			default:
				log.Fatalf("Expected field '%s' to be a singleton array if present, got: %s", field, iface)
			}
		default:
			log.Fatalf("Expected field '%s' to be a singleton array if present, got: %s", field, iface)
		}
	}
	return emptyDict, false
}

func singletonFromResourceData(block *schema.ResourceData, field string) (map[string]interface{}, bool) {
	if i, ok := block.GetOk(field); ok {
		if arr, ok := i.([]interface{}); ok && len(arr) == 1 {
			if elem, ok := arr[0].(map[string]interface{}); ok {
				return elem, true
			}
		}
	}
	return map[string]interface{}{}, false
}

func (t *TriggerCondition) PositiveTimeRange() string {
	return strings.TrimPrefix(t.TimeRange, "-")
}

func (t *TriggerCondition) PositiveBaselineWindow() string {
	return strings.TrimPrefix(t.BaselineWindow, "-")
}

func (t *TriggerCondition) PositiveResolutionWindow() string {
	return strings.TrimPrefix(t.ResolutionWindow, "-")
}

type schemaMap = map[string]*schema.Schema
type dict = map[string]interface{}
