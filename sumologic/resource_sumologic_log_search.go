package sumologic

import (
	"log"
    "regexp"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func resourceSumologicLogSearch() *schema.Resource {

	return &schema.Resource{
		Create: resourceSumologicLogSearchCreate,
		Read:   resourceSumologicLogSearchRead,
		Update: resourceSumologicLogSearchUpdate,
		Delete: resourceSumologicLogSearchDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
				ValidateFunc: validation.All(
					validation.StringLenBetween(1, 255),
				),
			},
			"description": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringLenBetween(0, 255),
				Default:      "",
			},
			"parent_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"query_string": {
				Type:     schema.TypeString,
				Required: true,
			},
			"run_by_receipt_time": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
            "interval_time_type": {
                Type:     schema.TypeString,
                Optional: true,
                ValidateFunc: validation.StringMatch(
                    regexp.MustCompile("^(messageTime|receiptTime|searchableTime)$"),
                    "should be either 'messageTime' or 'receiptTime' or 'searchableTime'",
                ),
            },
			"time_range": {
				Type:     schema.TypeList,
				Required: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: GetTimeRangeSchema(),
				},
			},
			"parsing_mode": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"AutoParse", "Manual"}, false),
				Default:      "Manual",
			},
			"query_parameter": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:     schema.TypeString,
							Required: true,
						},
						"description": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"value": {
							Type:     schema.TypeString,
							Required: true,
						},
						"data_type": {
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: validation.StringInSlice([]string{"NUMBER", "STRING", "ANY", "KEYWORD"}, false),
						},
					},
				},
			},
			"schedule": {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"cron_expression": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"parseable_time_range": {
							Type:     schema.TypeList,
							Required: true,
							MaxItems: 1,
							Elem: &schema.Resource{
								Schema: GetTimeRangeSchema(),
							},
						},
						"time_zone": {
							Type:     schema.TypeString,
							Required: true,
						},
						"threshold": {
							Type:     schema.TypeList,
							Optional: true,
							MaxItems: 1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"threshold_type": {
										Type:         schema.TypeString,
										Required:     true,
										ValidateFunc: validation.StringInSlice([]string{"message", "group"}, false),
									},
									"operator": {
										Type:         schema.TypeString,
										Required:     true,
										ValidateFunc: validation.StringInSlice([]string{"eq", "gt", "ge", "lt", "le"}, false),
									},
									"count": {
										Type:     schema.TypeInt,
										Required: true,
									},
								},
							},
						},
						"parameter": {
							Type:     schema.TypeList,
							Optional: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"name": {
										Type:     schema.TypeString,
										Required: true,
									},
									"value": {
										Type:     schema.TypeString,
										Required: true,
									},
								},
							},
						},
						"mute_error_emails": {
							Type:     schema.TypeBool,
							Optional: true,
						},
						"notification": {
							Type:     schema.TypeList,
							Required: true,
							MaxItems: 1,
							Elem: &schema.Resource{
								Schema: getSearchNotificationSchema(),
							},
						},
						"schedule_type": {
							Type:     schema.TypeString,
							Required: true,
							ValidateFunc: validation.StringInSlice([]string{"RealTime", "15Minutes", "1Hour", "2Hours", "4Hours",
								"6Hours", "8Hours", "12Hours", "1Day", "1Week", "Custom"}, false),
						},
					},
				},
			},
		},
	}
}

func getSearchNotificationSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"alert_search_notification": {
			Type:     schema.TypeList,
			Optional: true,
			MaxItems: 1,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"source_id": {
						Type:     schema.TypeString,
						Required: true,
					},
				},
			},
		},
		"cse_signal_notification": {
			Type:     schema.TypeList,
			Optional: true,
			MaxItems: 1,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"record_type": {
						Type:     schema.TypeString,
						Required: true,
					},
				},
			},
		},
		"email_search_notification": {
			Type:     schema.TypeList,
			Optional: true,
			MaxItems: 1,
			Elem: &schema.Resource{
				Schema: getEmailSearchNotificationSchema(),
			},
		},
		"save_to_view_notification": {
			Type:     schema.TypeList,
			Optional: true,
			MaxItems: 1,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"view_name": {
						Type:     schema.TypeString,
						Required: true,
					},
				},
			},
		},
		"save_to_lookup_notification": {
			Type:     schema.TypeList,
			Optional: true,
			MaxItems: 1,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"lookup_file_path": {
						Type:     schema.TypeString,
						Required: true,
					},
					"is_lookup_merge_operation": {
						Type:     schema.TypeBool,
						Required: true,
					},
				},
			},
		},
		"service_now_search_notification": {
			Type:     schema.TypeList,
			Optional: true,
			MaxItems: 1,
			Elem: &schema.Resource{
				Schema: getServiceNowSearchNotificationSchema(),
			},
		},
		"webhook_search_notification": {
			Type:     schema.TypeList,
			Optional: true,
			MaxItems: 1,
			Elem: &schema.Resource{
				Schema: getWebhookSearchNotificationSchema(),
			},
		},
	}
}

func getEmailSearchNotificationSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"subject_template": {
			Type:     schema.TypeString,
			Optional: true,
		},
		"include_query": {
			Type:     schema.TypeBool,
			Optional: true,
			Default:  true,
		},
		"include_result_set": {
			Type:     schema.TypeBool,
			Optional: true,
			Default:  true,
		},
		"include_histogram": {
			Type:     schema.TypeBool,
			Optional: true,
			Default:  true,
		},
		"include_csv_attachment": {
			Type:     schema.TypeBool,
			Optional: true,
			Default:  false,
		},
		"to_list": {
			Type:     schema.TypeList,
			Required: true,
			Elem: &schema.Schema{
				Type: schema.TypeString,
			},
		},
	}
}

func getServiceNowSearchNotificationSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"external_id": {
			Type:     schema.TypeString,
			Required: true,
		},
		"fields": {
			Type:     schema.TypeList,
			Optional: true,
			MaxItems: 1,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"event_type": {
						Type:     schema.TypeString,
						Optional: true,
					},
					"severity": {
						Type:         schema.TypeInt,
						Optional:     true,
						ValidateFunc: validation.IntInSlice([]int{0, 1, 2, 3, 4}),
					},
					"resource": {
						Type:     schema.TypeString,
						Optional: true,
					},
					"node": {
						Type:     schema.TypeString,
						Optional: true,
					},
				},
			},
		},
	}
}

func getWebhookSearchNotificationSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"webhook_id": {
			Type:     schema.TypeString,
			Required: true,
		},
		"payload": {
			Type:     schema.TypeString,
			Optional: true,
		},
		"itemize_alerts": {
			Type:     schema.TypeBool,
			Optional: true,
			Default:  false,
		},
		"max_itemized_alerts": {
			Type:         schema.TypeInt,
			Optional:     true,
			ValidateFunc: validation.IntBetween(1, 100),
		},
	}
}

func resourceSumologicLogSearchCreate(d *schema.ResourceData, meta interface{}) error {
	c := meta.(*Client)

	if d.Id() == "" {
		logSearch := resourceToLogSearch(d)
		log.Println("=====================================================================")
		log.Printf("creating log search - %+v", logSearch)
		log.Printf("log search schedule - %+v", logSearch.Schedule)
		log.Println("=====================================================================")
		id, err := c.CreateLogSearch(logSearch)
		if err != nil {
			return err
		}

		d.SetId(id)
	}

	return resourceSumologicLogSearchRead(d, meta)
}

func resourceSumologicLogSearchRead(d *schema.ResourceData, meta interface{}) error {
	c := meta.(*Client)

	id := d.Id()
	logSearch, err := c.GetLogSearch(id)
	if err != nil {
		return err
	}

	if logSearch == nil {
		log.Printf("[WARN] LogSearch not found, removing from state: %v - %v", id, err)
		d.SetId("")
		return nil
	}
	log.Println("=====================================================================")
	log.Printf("read log search - %+v", logSearch)
	log.Printf("log search schedule - %+v", logSearch.Schedule)
	log.Println("=====================================================================")

	return setLogSearch(d, logSearch)
}

func resourceSumologicLogSearchUpdate(d *schema.ResourceData, meta interface{}) error {
	c := meta.(*Client)

	logSearch := resourceToLogSearch(d)
	log.Println("=====================================================================")
	log.Printf("updating log search - %+v", logSearch)
	log.Printf("log search schedule - %+v", logSearch.Schedule)
	log.Println("=====================================================================")
	err := c.UpdateLogSearch(logSearch)
	if err != nil {
		return err
	}

	return resourceSumologicLogSearchRead(d, meta)
}

func resourceSumologicLogSearchDelete(d *schema.ResourceData, meta interface{}) error {
	c := meta.(*Client)
	return c.DeleteLogSearch(d.Id())
}

func setLogSearch(d *schema.ResourceData, logSearch *LogSearch) error {
	if err := d.Set("name", logSearch.Name); err != nil {
		return err
	}
	if err := d.Set("description", logSearch.Description); err != nil {
		return err
	}
	if err := d.Set("parent_id", logSearch.ParentId); err != nil {
		return err
	}
	if err := d.Set("query_string", logSearch.QueryString); err != nil {
		return err
	}
	if err := d.Set("run_by_receipt_time", logSearch.RunByReceiptTime); err != nil {
		return err
	}
    if err := d.Set("interval_time_type", logSearch.IntervalTimeType); err != nil {
    	return err
    }
	if err := d.Set("parsing_mode", logSearch.ParsingMode); err != nil {
		return err
	}

	queryParameters := make([]map[string]interface{}, len(logSearch.QueryParameters))
	for i, parameter := range logSearch.QueryParameters {
		queryParameters[i] = getTerraformLogSearchQueryParameter(parameter)
	}
	if err := d.Set("query_parameter", queryParameters); err != nil {
		return err
	}

	timeRange := GetTerraformTimeRange(logSearch.TimeRange.(map[string]interface{}))
	if err := d.Set("time_range", timeRange); err != nil {
		return err
	}

	if logSearch.Schedule != nil {
		searchSchedule := getTerraformLogSearchSchedule(logSearch.Schedule)
		if err := d.Set("schedule", searchSchedule); err != nil {
			return err
		}
	}

	return nil
}

func getTerraformLogSearchQueryParameter(parameter LogSearchQueryParameter) map[string]interface{} {
	tfSearchQueryParameter := map[string]interface{}{}

	tfSearchQueryParameter["name"] = parameter.Name
	tfSearchQueryParameter["description"] = parameter.Description
	tfSearchQueryParameter["data_type"] = parameter.DataType
	tfSearchQueryParameter["value"] = parameter.Value
	return tfSearchQueryParameter
}

func getTerraformLogSearchSchedule(schedule *LogSearchSchedule) []map[string]interface{} {
	tfSearchSchedule := []map[string]interface{}{}
	tfSearchSchedule = append(tfSearchSchedule, make(map[string]interface{}))

	if schedule == nil {
		return tfSearchSchedule
	}

	tfSearchSchedule[0]["cron_expression"] = schedule.CronExpression
	tfSearchSchedule[0]["time_zone"] = schedule.TimeZone
	tfSearchSchedule[0]["mute_error_emails"] = schedule.MuteErrorEmails
	tfSearchSchedule[0]["schedule_type"] = schedule.ScheduleType

	tfSearchSchedule[0]["parseable_time_range"] =
		GetTerraformTimeRange(schedule.ParseableTimeRange.(map[string]interface{}))

	tfSearchSchedule[0]["notification"] =
		getTerraformLogSearchNotification(schedule.Notification.(map[string]interface{}))

	if schedule.Threshold != nil {
		tfSearchSchedule[0]["threshold"] = getTerraformLogSearchNotificationThreshold(schedule.Threshold)
	}

	parameters := make([]map[string]interface{}, len(schedule.Parameters))
	for i, parameter := range schedule.Parameters {
		parameters[i] = getTerraformLogSearchScheduleParameter(parameter)
	}
	tfSearchSchedule[0]["parameter"] = parameters

	return tfSearchSchedule
}

func getTerraformLogSearchNotification(notification map[string]interface{}) TerraformObject {
	tfNotification := MakeTerraformObject()

	if notification["taskType"] == "AlertSearchNotificationSyncDefinition" {
		alertNotification := MakeTerraformObject()
		alertNotification[0]["source_id"] = notification["sourceId"]
		tfNotification[0]["alert_search_notification"] = alertNotification
	} else if notification["taskType"] == "CseSignalNotificationSyncDefinition" {
		cseSignalNotification := MakeTerraformObject()
		cseSignalNotification[0]["record_type"] = notification["recordType"]
		tfNotification[0]["cse_signal_notification"] = cseSignalNotification
	} else if notification["taskType"] == "EmailSearchNotificationSyncDefinition" {
		emailNotification := MakeTerraformObject()
		emailNotification[0]["include_csv_attachment"] = notification["includeCsvAttachment"]
		emailNotification[0]["include_histogram"] = notification["includeHistogram"]
		emailNotification[0]["include_query"] = notification["includeQuery"]
		emailNotification[0]["include_result_set"] = notification["includeResultSet"]
		emailNotification[0]["subject_template"] = notification["subjectTemplate"]
		emailNotification[0]["to_list"] = notification["toList"]

		tfNotification[0]["email_search_notification"] = emailNotification
	} else if notification["taskType"] == "SaveToViewNotificationSyncDefinition" {
		saveToViewNotification := MakeTerraformObject()
		saveToViewNotification[0]["view_name"] = notification["viewName"]
		tfNotification[0]["save_to_view_notification"] = saveToViewNotification
	} else if notification["taskType"] == "SaveToLookupNotificationSyncDefinition" {
		saveToLookupNotification := MakeTerraformObject()
		saveToLookupNotification[0]["lookup_file_path"] = notification["lookupFilePath"]
		saveToLookupNotification[0]["is_lookup_merge_operation"] = notification["isLookupMergeOperation"]
		tfNotification[0]["save_to_lookup_notification"] = saveToLookupNotification
	} else if notification["taskType"] == "ServiceNowSearchNotificationSyncDefinition" {
		serviceNowNotification := MakeTerraformObject()
		serviceNowNotification[0]["external_id"] = notification["externalId"]
		if serviceNowFields := notification["fields"]; serviceNowFields != nil {
			fields := serviceNowFields.(map[string]interface{})
			tfServiceNowFields := MakeTerraformObject()
			tfServiceNowFields[0]["event_type"] = fields["eventType"]
			tfServiceNowFields[0]["severity"] = fields["severity"]
			tfServiceNowFields[0]["resource"] = fields["resource"]
			tfServiceNowFields[0]["node"] = fields["node"]
			serviceNowNotification[0]["fields"] = tfServiceNowFields
		}
		tfNotification[0]["service_now_search_notification"] = serviceNowNotification
	} else if notification["taskType"] == "WebhookSearchNotificationSyncDefinition" {
		webhookNotification := MakeTerraformObject()
		webhookNotification[0]["webhook_id"] = notification["webhookId"]
		webhookNotification[0]["payload"] = notification["payload"]
		webhookNotification[0]["itemize_alerts"] = notification["itemizeAlerts"]
		webhookNotification[0]["max_itemized_alerts"] = notification["maxItemizedAlerts"]
		tfNotification[0]["webhook_search_notification"] = webhookNotification
	}

	return tfNotification
}

func getTerraformLogSearchNotificationThreshold(threshold *SearchNotificationThreshold) TerraformObject {
	tfNotificationThreshold := MakeTerraformObject()
	if threshold == nil {
		return tfNotificationThreshold
	}

	tfNotificationThreshold[0]["threshold_type"] = threshold.ThresholdType
	tfNotificationThreshold[0]["operator"] = threshold.Operator
	tfNotificationThreshold[0]["count"] = threshold.Count

	return tfNotificationThreshold
}

func getTerraformLogSearchScheduleParameter(parameter ScheduleSearchParameter) map[string]interface{} {
	tfSearchParameter := map[string]interface{}{}

	tfSearchParameter["name"] = parameter.Name
	tfSearchParameter["value"] = parameter.Value
	return tfSearchParameter
}

func resourceToLogSearch(d *schema.ResourceData) LogSearch {

	var queryParameters []LogSearchQueryParameter
	if val, ok := d.GetOk("query_parameter"); ok {
		queryParametersData := val.([]interface{})
		for _, data := range queryParametersData {
			queryParameters = append(queryParameters, resourceToLogSearchQueryParameter([]interface{}{data}))
		}
	}

	var timeRange interface{}
	if val, ok := d.GetOk("time_range"); ok {
		tfTimeRange := val.([]interface{})[0]
		timeRange = GetTimeRange(tfTimeRange.(map[string]interface{}))
	}

	var schedule *LogSearchSchedule
	if tfSchedule, ok := d.GetOk("schedule"); ok {
		schedule = resourceToLogSearchSchedule(tfSchedule)
	}

	return LogSearch{
		ID:               d.Id(),
		Name:             d.Get("name").(string),
		Description:      d.Get("description").(string),
		ParentId:         d.Get("parent_id").(string),
		QueryString:      d.Get("query_string").(string),
		RunByReceiptTime: d.Get("run_by_receipt_time").(bool),
		IntervalTimeType: d.Get("interval_time_type").(string),
		TimeRange:        timeRange,
		ParsingMode:      d.Get("parsing_mode").(string),
		QueryParameters:  queryParameters,
		Schedule:         schedule,
	}
}

func resourceToLogSearchQueryParameter(data interface{}) LogSearchQueryParameter {
	queryParameter := LogSearchQueryParameter{}

	queryParameterSlice := data.([]interface{})
	if len(queryParameterSlice) > 0 {
		queryParameterObj := queryParameterSlice[0].(map[string]interface{})

		queryParameter.Name = queryParameterObj["name"].(string)
		queryParameter.DataType = queryParameterObj["data_type"].(string)
		queryParameter.Value = queryParameterObj["value"].(string)

		if description, found := queryParameterObj["description"]; found {
			queryParameter.Description = description.(string)
		}
	}

	return queryParameter
}

func resourceToLogSearchSchedule(data interface{}) *LogSearchSchedule {

	scheduleSlice := data.([]interface{})
	schedule := LogSearchSchedule{}
	if len(scheduleSlice) > 0 {
		scheduleObj := scheduleSlice[0].(map[string]interface{})
		schedule.Threshold = resourceToSearchNotificationThreshold(scheduleObj["threshold"])

		parametersData := scheduleObj["parameter"].([]interface{})
		parameters := make([]ScheduleSearchParameter, len(parametersData))
		for i, v := range parametersData {
			parameters[i] = resourceToScheduleSearchParameter(v)
		}
		schedule.Parameters = parameters

		tfTimeRange := scheduleObj["parseable_time_range"].([]interface{})[0]
		schedule.ParseableTimeRange = GetTimeRange(tfTimeRange.(map[string]interface{}))

		schedule.TimeZone = scheduleObj["time_zone"].(string)
		schedule.CronExpression = scheduleObj["cron_expression"].(string)
		schedule.MuteErrorEmails = scheduleObj["mute_error_emails"].(bool)
		schedule.Notification = resourceToScheduleSearchNotification(scheduleObj["notification"])
		schedule.ScheduleType = scheduleObj["schedule_type"].(string)
	}

	return &schedule
}

func resourceToSearchNotificationThreshold(data interface{}) *SearchNotificationThreshold {

	thresholdObj := data.([]interface{})
	if len(thresholdObj) == 0 {
		return nil
	}

	tfThreshold := thresholdObj[0].(map[string]interface{})

	threshold := SearchNotificationThreshold{}
	threshold.ThresholdType = tfThreshold["threshold_type"].(string)
	threshold.Operator = tfThreshold["operator"].(string)
	threshold.Count = tfThreshold["count"].(int)

	return &threshold
}

func resourceToScheduleSearchParameter(data interface{}) ScheduleSearchParameter {

	tfSearchParameter := data.(map[string]interface{})
	return ScheduleSearchParameter{
		Name:  tfSearchParameter["name"].(string),
		Value: tfSearchParameter["value"].(string),
	}
}

func resourceToScheduleSearchNotification(data interface{}) interface{} {

	notificationSlice := data.([]interface{})
	if len(notificationSlice) > 0 && notificationSlice[0] != nil {
		notificationObj := notificationSlice[0].(map[string]interface{})

		if val := notificationObj["alert_search_notification"].([]interface{}); len(val) == 1 {
			if tfAlertNotification, ok := val[0].(map[string]interface{}); ok {
				return AlertSearchNotification{
					TaskType: "AlertSearchNotificationSyncDefinition",
					SourceId: tfAlertNotification["source_id"].(string),
				}
			}
		} else if val := notificationObj["cse_signal_notification"].([]interface{}); len(val) == 1 {
			if tfCseNotification, ok := val[0].(map[string]interface{}); ok {
				return CseSignalNotification{
					TaskType:   "CseSignalNotificationSyncDefinition",
					RecordType: tfCseNotification["record_type"].(string),
				}
			}
		} else if val := notificationObj["email_search_notification"].([]interface{}); len(val) == 1 {
			if tfEmailSearchNotification, ok := val[0].(map[string]interface{}); ok {
				return getEmailSearchNotification(tfEmailSearchNotification)
			}
		} else if val := notificationObj["save_to_view_notification"].([]interface{}); len(val) == 1 {
			if tfSaveToViewNotification, ok := val[0].(map[string]interface{}); ok {
				return SaveToViewNotification{
					TaskType: "SaveToViewNotificationSyncDefinition",
					ViewName: tfSaveToViewNotification["view_name"].(string),
				}
			}
		} else if val := notificationObj["save_to_lookup_notification"].([]interface{}); len(val) == 1 {
			if tfSaveToLookupNotification, ok := val[0].(map[string]interface{}); ok {
				return SaveToLookupNotification{
					TaskType:               "SaveToLookupNotificationSyncDefinition",
					LookupFilePath:         tfSaveToLookupNotification["lookup_file_path"].(string),
					IsLookupMergeOperation: tfSaveToLookupNotification["is_lookup_merge_operation"].(bool),
				}
			}
		} else if val := notificationObj["service_now_search_notification"].([]interface{}); len(val) == 1 {
			if tfServiceNowSearchNotification, ok := val[0].(map[string]interface{}); ok {
				return getServiceNowSearchNotification(tfServiceNowSearchNotification)
			}
		} else if val := notificationObj["webhook_search_notification"].([]interface{}); len(val) == 1 {
			if tfWebhookSearchNotification, ok := val[0].(map[string]interface{}); ok {
				return getWebhookSearchNotification(tfWebhookSearchNotification)
			}
		}
	}

	return nil
}

func getEmailSearchNotification(tfEmailSearchNotification map[string]interface{}) interface{} {

	tfToList := tfEmailSearchNotification["to_list"].([]interface{})
	toList := make([]string, len(tfToList))
	for i, v := range tfToList {
		toList[i] = v.(string)
	}

	return EmailSearchNotification{
		TaskType:             "EmailSearchNotificationSyncDefinition",
		ToList:               toList,
		SubjectTemplate:      tfEmailSearchNotification["subject_template"].(string),
		IncludeQuery:         tfEmailSearchNotification["include_query"].(bool),
		IncludeResultSet:     tfEmailSearchNotification["include_result_set"].(bool),
		IncludeHistogram:     tfEmailSearchNotification["include_histogram"].(bool),
		IncludeCsvAttachment: tfEmailSearchNotification["include_csv_attachment"].(bool),
	}
}

func getServiceNowSearchNotification(tfServiceNowSearchNotification map[string]interface{}) interface{} {

	fields := ServiceNowFields{}
	if val := tfServiceNowSearchNotification["fields"].([]interface{}); len(val) == 1 {
		if tfServiceNowFields, ok := val[0].(map[string]interface{}); ok {
			fields = ServiceNowFields{
				EventType: tfServiceNowFields["event_type"].(string),
				Severity:  tfServiceNowFields["severity"].(int),
				Resource:  tfServiceNowFields["resource"].(string),
				Node:      tfServiceNowFields["node"].(string),
			}
		}
	}

	return ServiceNowSearchNotification{
		TaskType:   "ServiceNowSearchNotificationSyncDefinition",
		ExternalId: tfServiceNowSearchNotification["external_id"].(string),
		Fields:     fields,
	}
}

func getWebhookSearchNotification(tfWebhookSearchNotification map[string]interface{}) interface{} {
	var payload *string
	if p, ok := tfWebhookSearchNotification["payload"].(string); ok && p != "" {
		payload = &p
	} else {
		payload = nil
	}
	return WebhookSearchNotification{
		TaskType:          "WebhookSearchNotificationSyncDefinition",
		WebhookId:         tfWebhookSearchNotification["webhook_id"].(string),
		Payload:           payload,
		ItemizeAlerts:     tfWebhookSearchNotification["itemize_alerts"].(bool),
		MaxItemizedAlerts: tfWebhookSearchNotification["max_itemized_alerts"].(int),
	}
}
