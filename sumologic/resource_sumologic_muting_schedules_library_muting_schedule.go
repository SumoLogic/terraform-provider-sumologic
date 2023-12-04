package sumologic

import (
	"fmt"
	"log"
	"regexp"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

func resourceSumologicMutingSchedulesLibraryMutingSchedule() *schema.Resource {
	return &schema.Resource{
		Create: resourceSumologicMutingSchedulesLibraryMutingScheduleCreate,
		Read:   resourceSumologicMutingSchedulesLibraryMutingScheduleRead,
		Update: resourceSumologicMutingSchedulesLibraryMutingScheduleUpdate,
		Delete: resourceSumologicMutingSchedulesLibraryMutingScheduleDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: getMutingScheduleSchema(),
	}
}

func getMutingScheduleBaseSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{

		"name": {
			Type:     schema.TypeString,
			Required: true,
			ValidateFunc: validation.All(
				validation.StringDoesNotContainAny("/"),
				validation.StringLenBetween(1, 255),
				validation.StringMatch(regexp.MustCompile(`(?s)^[^\ ].*[^\ ]$`),
					"name must not contain leading or trailing spaces"),
			),
		},

		"description": {
			Type:     schema.TypeString,
			Optional: true,
			ValidateFunc: validation.All(
				validation.StringLenBetween(0, 4096),
				validation.StringMatch(regexp.MustCompile(`(?s)^[^\ ].*[^\ ]$`),
					"description must not contain leading or trailing spaces"),
			),
		},

		"parent_id": {
			Type:     schema.TypeString,
			Optional: true,
			Computed: true,
		},

		"monitor": {
			Type:     schema.TypeList,
			Optional: true,
			MaxItems: 1,
			Elem: &schema.Resource{
				Schema: getMonitorScopeSchema(),
			},
			AtLeastOneOf: monitorAtleastOneKey,
		},

		"schedule": {
			Type:     schema.TypeList,
			Required: true,
			MaxItems: 1,
			Elem: &schema.Resource{
				Schema: getScheduleDefinitionSchemma(),
			},
		},
		"notification_groups": {
			Type:     schema.TypeList,
			Optional: true,
			MaxItems: 10,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"group_key": {
						Type:         schema.TypeString,
						ValidateFunc: validation.StringLenBetween(1, 128),
						Required:     true,
					},
					"group_values": {
						Type:     schema.TypeList,
						Required: true,
						Elem: &schema.Schema{
							Type:         schema.TypeString,
							ValidateFunc: validation.StringLenBetween(1, 256),
						},
					},
				},
			},
		},

		"version": {
			Type:     schema.TypeInt,
			Optional: true,
			Computed: true,
		},
	}
}

func getMonitorScopeSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"ids": {
			Type:     schema.TypeList,
			Optional: true,
			Elem: &schema.Schema{
				Type:         schema.TypeString,
				ValidateFunc: validation.StringIsNotEmpty,
			},
		},
		"all": {
			Type:     schema.TypeBool,
			Optional: true,
		},
	}
}

func getScheduleDefinitionSchemma() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"timezone": {
			Type:     schema.TypeString,
			Required: true,
			ForceNew: false,
		},
		"start_date": {
			Type:     schema.TypeString,
			ForceNew: false,
			Required: true,
			ValidateFunc: validation.All(
				validation.StringMatch(regexp.MustCompile(`^\d{4}-\d{2}-\d{2}$`),
					"start date in format of yyyy-mm-dd"),
				StartDateIsAfterYesterday(),
			),
		},
		"start_time": {
			Type:     schema.TypeString,
			Required: true,
			ForceNew: false,
			ValidateFunc: validation.All(
				validation.StringMatch(regexp.MustCompile(`^(?:[01]\d|2[0-3]):[0-5]\d$`),
					"start time in format of 00:00"),
			),
		},
		"duration": {
			Type:         schema.TypeInt,
			Required:     true,
			ForceNew:     false,
			ValidateFunc: validation.IntAtLeast(15),
		},
		"rrule": {
			Type:     schema.TypeString,
			Optional: true,
		},
	}
}

func getMutingScheduleSchema() map[string]*schema.Schema {
	tfSchema := getMutingScheduleBaseSchema()

	additionalAttributes := map[string]*schema.Schema{

		"type": {
			Type:         schema.TypeString,
			Optional:     true,
			Default:      "MutingSchedulesLibraryMutingSchedule",
			ValidateFunc: validation.StringInSlice([]string{"MutingSchedulesLibraryMutingSchedule", "MutingSchedulesLibraryFolder"}, false),
		},

		"content_type": {
			Type:     schema.TypeString,
			Optional: true,
			Default:  "MutingSchedule",
		},

		"is_system": {
			Type:     schema.TypeBool,
			Optional: true,
			Computed: true,
		},

		"is_mutable": {
			Type:     schema.TypeBool,
			Optional: true,
			Computed: true,
		},

		"created_by": {
			Type:     schema.TypeString,
			Optional: true,
			Computed: true,
		},

		"created_at": {
			Type:     schema.TypeString,
			Optional: true,
			Computed: true,
		},

		"modified_by": {
			Type:     schema.TypeString,
			Optional: true,
			Computed: true,
		},

		"modified_at": {
			Type:     schema.TypeString,
			Optional: true,
			Computed: true,
		},
	}

	for k, v := range additionalAttributes {
		tfSchema[k] = v
	}

	return tfSchema
}

func resourceSumologicMutingSchedulesLibraryMutingScheduleCreate(d *schema.ResourceData, meta interface{}) error {
	c := meta.(*Client)

	if d.Id() == "" {
		mutingSchedule := resourceToMutingSchedulesLibraryMutingSchedule(d)
		if mutingSchedule.ParentID == "" {
			rootFolder, err := c.GetMutingSchedulesLibraryFolder("root")
			if err != nil {
				return err
			}

			mutingSchedule.ParentID = rootFolder.ID
		}
		paramMap := map[string]string{
			"parentId": mutingSchedule.ParentID,
		}
		mutingScheduleDefinitionID, err := c.CreateMutingSchedulesLibraryMutingSchedule(mutingSchedule, paramMap)
		if err != nil {
			return err
		}
		d.SetId(mutingScheduleDefinitionID)
	}
	return resourceSumologicMutingSchedulesLibraryMutingScheduleRead(d, meta)
}

func resourceSumologicMutingSchedulesLibraryMutingScheduleRead(d *schema.ResourceData, meta interface{}) error {
	c := meta.(*Client)

	mutingSchedule, err := c.MutingSchedulesRead(d.Id())
	if err != nil {
		return err
	}

	if mutingSchedule == nil {
		log.Printf("[WARN] MutingSchedule not found, removing from state: %v - %v", d.Id(), err)
		d.SetId("")
		return nil
	}

	monitorScope := make([]interface{}, 1)
	monitorScope[0] = map[string]interface{}{
		"all": mutingSchedule.Monitor.All,
		"ids": mutingSchedule.Monitor.Ids,
	}
	schedule := make([]interface{}, 1)
	schedule[0] = map[string]interface{}{
		"timezone":   mutingSchedule.Schedule.TimeZone,
		"start_date": mutingSchedule.Schedule.StartDate,
		"start_time": mutingSchedule.Schedule.StartTime,
		"duration":   mutingSchedule.Schedule.Duration,
		"rrule":      mutingSchedule.Schedule.RRule,
	}

	d.Set("created_by", mutingSchedule.CreatedBy)
	d.Set("created_at", mutingSchedule.CreatedAt)
	d.Set("modified_by", mutingSchedule.ModifiedBy)
	d.Set("is_mutable", mutingSchedule.IsMutable)
	d.Set("version", mutingSchedule.Version)
	d.Set("description", mutingSchedule.Description)
	d.Set("name", mutingSchedule.Name)
	d.Set("parent_id", mutingSchedule.ParentID)
	d.Set("modified_at", mutingSchedule.ModifiedAt)
	d.Set("content_type", mutingSchedule.ContentType)
	d.Set("is_system", mutingSchedule.IsSystem)
	d.Set("monitor", monitorScope)
	d.Set("schedule", schedule)
	d.Set("notification_groups", notificationGroupArrayToResource(mutingSchedule.NotificationGroups))

	return nil
}

func resourceSumologicMutingSchedulesLibraryMutingScheduleUpdate(d *schema.ResourceData, meta interface{}) error {
	c := meta.(*Client)
	mutingSchedule := resourceToMutingSchedulesLibraryMutingSchedule(d)

	mutingSchedule.Type = "MutingSchedulesLibraryMutingScheduleUpdate"
	err := c.UpdateMutingSchedulesLibraryMutingSchedule(mutingSchedule)
	if err != nil {
		return err
	}
	updatedMutingSchedule := resourceSumologicMutingSchedulesLibraryMutingScheduleRead(d, meta)
	return updatedMutingSchedule
}

func resourceSumologicMutingSchedulesLibraryMutingScheduleDelete(d *schema.ResourceData, meta interface{}) error {
	c := meta.(*Client)
	mutingSchedule := resourceToMutingSchedulesLibraryMutingSchedule(d)
	err := c.DeleteMutingSchedulesLibraryMutingSchedule(mutingSchedule.ID)
	if err != nil {
		return err
	}
	return nil
}

func getMonitorScope(d *schema.ResourceData) *MonitorScope {
	monitorMap := d.Get("monitor").([]interface{})
	if len(monitorMap) == 0 {
		return nil
	} else {
		monitorScopeDict := monitorMap[0].(map[string]interface{})
		monitorScope := MonitorScope{
			Ids: fieldsToStringArray(monitorScopeDict["ids"].([]interface{})),
			All: monitorScopeDict["all"].(bool),
		}
		return &monitorScope
	}
}

func getScheduleDefinition(d *schema.ResourceData) ScheduleDefinition {
	scheduleDefinitionMap := d.Get("schedule").([]interface{})
	scheduleDefinitionDict := scheduleDefinitionMap[0].(map[string]interface{})
	scheduleDefinition := ScheduleDefinition{
		TimeZone:  scheduleDefinitionDict["timezone"].(string),
		StartDate: scheduleDefinitionDict["start_date"].(string),
		StartTime: scheduleDefinitionDict["start_time"].(string),
		Duration:  scheduleDefinitionDict["duration"].(int),
		RRule:     scheduleDefinitionDict["rrule"].(string),
	}
	return scheduleDefinition
}

func notificationGroupArrayToResource(notificationGroups []NotificationGroupDefinition) []map[string]interface{} {
	result := make([]map[string]interface{}, len(notificationGroups))

	for i, notificationGroup := range notificationGroups {
		result[i] = map[string]interface{}{
			"group_key":    notificationGroup.GroupKey,
			"group_values": notificationGroup.GroupValues,
		}
	}
	return result
}

func getNotificationGroupArray(resourceNotificationGroups []interface{}) []NotificationGroupDefinition {
	result := make([]NotificationGroupDefinition, len(resourceNotificationGroups))

	for i, resourceNotificationGroup := range resourceNotificationGroups {
		resourceNotificationGroupMap := resourceNotificationGroup.(map[string]interface{})
		result[i] = NotificationGroupDefinition{
			GroupKey:    resourceNotificationGroupMap["group_key"].(string),
			GroupValues: resourceToStringArray(resourceNotificationGroupMap["group_values"].([]interface{})),
		}
	}
	return result
}

func StartDateIsAfterYesterday() schema.SchemaValidateFunc {
	return func(i interface{}, k string) (warnings []string, errors []error) {
		v, ok := i.(string)
		if !ok {
			return warnings, []error{fmt.Errorf("expected type of %q to be string", k)}
		}

		date, err := time.Parse("2006-01-02", v)

		if err != nil {
			return warnings, []error{fmt.Errorf("expected %q to be an valid start date yyyy-mm-dd : got %v", k, v)}
		}

		yesterday := time.Now().AddDate(0, 0, -1).Truncate(24 * time.Hour)

		if date.Before(yesterday) {
			return warnings, []error{fmt.Errorf("expected %q to be an valid start date : got %v", k, v)}
		}
		return warnings, errors
	}
}

var monitorAtleastOneKey = []string{
	"monitor.0.ids",
	"monitor.0.all",
}

func resourceToMutingSchedulesLibraryMutingSchedule(d *schema.ResourceData) MutingSchedulesLibraryMutingSchedule {
	monitorScope := getMonitorScope(d)
	scheduleDefinition := getScheduleDefinition(d)

	return MutingSchedulesLibraryMutingSchedule{
		CreatedBy:          d.Get("created_by").(string),
		Name:               d.Get("name").(string),
		ID:                 d.Id(),
		CreatedAt:          d.Get("created_at").(string),
		Description:        d.Get("description").(string),
		ModifiedBy:         d.Get("modified_by").(string),
		IsMutable:          d.Get("is_mutable").(bool),
		Version:            d.Get("version").(int),
		Type:               d.Get("type").(string),
		ParentID:           d.Get("parent_id").(string),
		ModifiedAt:         d.Get("modified_at").(string),
		ContentType:        d.Get("content_type").(string),
		IsSystem:           d.Get("is_system").(bool),
		Schedule:           scheduleDefinition,
		Monitor:            monitorScope,
		NotificationGroups: getNotificationGroupArray(d.Get("notification_groups").([]interface{})),
	}
}
