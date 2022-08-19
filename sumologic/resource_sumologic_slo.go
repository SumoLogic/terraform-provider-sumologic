package sumologic

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"log"
)

const fieldNameWindowBasedEvaluation = `window_based_evaluation`
const fieldNameRequestBasedEvaluation = `request_based_evaluation`
const sloContentTypeString = "Slo"

func resourceSumologicSLO() *schema.Resource {

	queryGroupElemSchema := &schema.Resource{
		Schema: map[string]*schema.Schema{
			"row_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"query": {
				Type:     schema.TypeString,
				Required: true,
			},
			"use_row_count": {
				Type:     schema.TypeBool,
				Required: true,
			},
			"field": {
				Type:     schema.TypeString,
				Optional: true,
			},
		},
	}

	queriesSchema := &schema.Schema{
		Type:     schema.TypeList,
		Required: true,
		MinItems: 1,
		MaxItems: 6,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"query_group_type": {
					Type:     schema.TypeString,
					Required: true,
					ValidateFunc: validation.StringInSlice([]string{
						"Successful", "Unsuccessful", "Total", "Threshold",
					}, false),
				},
				"query_group": {
					Type:     schema.TypeList,
					MinItems: 1,
					MaxItems: 6,
					Required: true,
					Elem:     queryGroupElemSchema,
				},
			},
		},
	}

	requestBasedIndicatorSchema := &schema.Resource{
		Schema: map[string]*schema.Schema{
			"queries": queriesSchema,
			"query_type": {
				Type:     schema.TypeString,
				Required: true,
				ValidateFunc: validation.StringInSlice([]string{
					"Logs", "Metrics",
				}, false),
			},
			"threshold": {
				Type:     schema.TypeFloat,
				Optional: true,
			},
			"op": {
				Type:     schema.TypeString,
				Optional: true,
				ValidateFunc: validation.StringInSlice([]string{
					"LessThan", "GreaterThan", "LessThanOrEqual", "GreaterThanOrEqual",
				}, false),
			},
		},
	}

	windowBasedIndicatorSchema := &schema.Resource{
		Schema: map[string]*schema.Schema{
			"queries": queriesSchema,
			"query_type": {
				Type:     schema.TypeString,
				Required: true,
				ValidateFunc: validation.StringInSlice([]string{
					"Logs", "Metrics",
				}, false),
			},
			"threshold": {
				Type:     schema.TypeFloat,
				Required: true,
			},
			"op": {
				Type:     schema.TypeString,
				Required: true,
				ValidateFunc: validation.StringInSlice([]string{
					"LessThan", "GreaterThan", "LessThanOrEqual", "GreaterThanOrEqual",
				}, false),
			},
			"aggregation": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"size": {
				Type:     schema.TypeString,
				Required: true,
			},
		},
	}

	return &schema.Resource{
		Create: resourceSumologicSLOCreate,
		Read:   resourceSLORead,
		Update: resourceSumologicSLOUpdate,
		Delete: resourceSumologicSLODelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{

			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"version": {
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
			},
			"created_at": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"created_by": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"modified_at": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"modified_by": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"parent_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"is_system": {
				Type:     schema.TypeBool,
				Optional: true,
				Computed: true,
			},
			"signal_type": {
				Type:     schema.TypeString,
				Required: true,
				ValidateFunc: validation.StringInSlice([]string{
					"Latency", "Error", "Throughput", "Availability", "Other",
				}, false),
			},
			"compliance": {
				Type:     schema.TypeList,
				Required: true,
				//MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"compliance_type": {
							Type:     schema.TypeString,
							Required: true,
							ValidateFunc: validation.StringInSlice([]string{
								"Rolling",
								"Calendar",
							}, false),
						},
						"target": {
							Type:         schema.TypeFloat,
							Required:     true,
							ValidateFunc: validation.FloatBetween(0, 100),
						},
						"timezone": {
							Type:     schema.TypeString,
							Required: true,
						},
						"size": {
							Type:     schema.TypeString,
							Required: true,
							ValidateFunc: validation.StringInSlice([]string{
								"Week", "Month", "Quarter",
								"1d", "2d", "3d", "4d", "5d", "6d", "7d", "8d", "9d", "10d", "11d", "12d", "13d", "14d",
							}, false),
						},
						"start_from": {
							Type:     schema.TypeString,
							Optional: true,
						},
					},
				},
			},
			"indicator": {
				Type:     schema.TypeList,
				MaxItems: 1,
				MinItems: 1,
				Required: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						fieldNameWindowBasedEvaluation: {
							Type:     schema.TypeList,
							MaxItems: 1,
							Optional: true,
							Elem:     windowBasedIndicatorSchema,
						},
						fieldNameRequestBasedEvaluation: {
							Type:     schema.TypeList,
							MaxItems: 1,
							Optional: true,
							Elem:     requestBasedIndicatorSchema,
						},
					},
				},
			},
			"is_mutable": {
				Type:     schema.TypeBool,
				Optional: true,
				Computed: true,
			},
			"is_locked": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"service": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"application": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
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

func resourceSumologicSLOCreate(d *schema.ResourceData, meta interface{}) error {
	c := meta.(*Client)
	if d.Id() == "" {
		slo, err := resourceToSLO(d)

		if err != nil {
			return err
		}

		slo.Type = "SlosLibrarySlo"
		if slo.ParentID == "" {
			rootFolder, err := c.GetSLOLibraryFolder("root")
			if err != nil {
				return err
			}

			slo.ParentID = rootFolder.ID
		}
		paramMap := map[string]string{
			"parentId": slo.ParentID,
		}
		sloDefinitionID, err := c.CreateSLO(*slo, paramMap)
		if err != nil {
			return err
		}

		d.SetId(sloDefinitionID)
	}

	return resourceSLORead(d, meta)
}

func resourceSLORead(d *schema.ResourceData, meta interface{}) error {
	c := meta.(*Client)

	slo, err := c.SLORead(d.Id())
	if err != nil {
		return err
	}

	if slo == nil {
		log.Printf("[WARN] SLO not found, removing from state: %v - %v", d.Id(), err)
		d.SetId("")
		return nil
	}

	d.Set("name", slo.Name)
	d.Set("description", slo.Description)
	d.Set("version", slo.Version)
	d.Set("created_at", slo.CreatedAt)
	d.Set("created_by", slo.CreatedBy)
	d.Set("modified_at", slo.ModifiedAt)
	d.Set("modified_by", slo.ModifiedBy)
	d.Set("parent_id", slo.ParentID)
	d.Set("content_type", sloContentTypeString)
	d.Set("is_mutable", slo.IsMutable)
	d.Set("is_locked", slo.IsLocked)
	d.Set("is_system", slo.IsSystem)
	d.Set("service", slo.Service)
	d.Set("application", slo.Application)

	flatCompliance, err := flattenSLOCompliance(slo.Compliance)
	if err != nil {
		return err
	}

	if err := d.Set("compliance", flatCompliance); err != nil {
		return fmt.Errorf("error setting compliance fields for resource %s: %s", d.Id(), err)
	}

	flatIndicator, err := flattenSLOIndicator(slo.Indicator)

	if err := d.Set("indicator", []interface{}{flatIndicator}); err != nil {
		return fmt.Errorf("error setting indicator field for resource %s: %s", d.Id(), err)
	}

	return nil
}

func flattenSLOCompliance(c SLOCompliance) ([]interface{}, error) {

	switch c.ComplianceType {
	case "Rolling":
		return []interface{}{flattenRollingCompliance(c)}, nil
	case "Calendar":
		return []interface{}{flattenCalendarCompliance(c)}, nil
	default:
		return nil, fmt.Errorf("unhandled compliance type found : %s", c.ComplianceType)
	}
}

func flattenRollingCompliance(c SLOCompliance) map[string]interface{} {
	comp := map[string]interface{}{}

	comp["compliance_type"] = c.ComplianceType
	comp["timezone"] = c.Timezone
	comp["timezone"] = c.Timezone
	comp["target"] = c.Target
	comp["size"] = c.Size

	return comp
}

func flattenCalendarCompliance(c SLOCompliance) map[string]interface{} {
	comp := map[string]interface{}{}

	comp["compliance_type"] = c.ComplianceType
	comp["timezone"] = c.Timezone
	comp["target"] = c.Target
	comp["size"] = c.WindowType
	comp["start_from"] = c.StartFrom

	return comp
}

func flattenSLOIndicator(ind SLOIndicator) (map[string]interface{}, error) {

	switch ind.EvaluationType {
	case "Request":
		return flattenRequestIndicator(ind)
	case "Window":
		return flattenWindowIndicator(ind)
	default:
		return nil, fmt.Errorf("unhandled indicator type found : %s", ind.EvaluationType)
	}
}

func flattenRequestIndicator(ind SLOIndicator) (map[string]interface{}, error) {
	reqIndicator := map[string]interface{}{}
	reqIndicator["query_type"] = ind.QueryType
	queries, err := flattenSLOQueries(ind.Queries)

	if err != nil {
		return nil, err
	}

	reqIndicator["queries"] = queries
	reqIndicator["threshold"] = ind.Threshold
	reqIndicator["op"] = ind.Op

	return map[string]interface{}{
		fieldNameRequestBasedEvaluation: []interface{}{reqIndicator},
		//fieldNameWindowBasedEvaluation:  []interface{}{},
	}, nil
}

func flattenWindowIndicator(ind SLOIndicator) (map[string]interface{}, error) {
	windowIndicator := map[string]interface{}{}
	windowIndicator["query_type"] = ind.QueryType
	queries, err := flattenSLOQueries(ind.Queries)

	if err != nil {
		return nil, err
	}

	windowIndicator["queries"] = queries
	windowIndicator["aggregation"] = ind.Aggregation
	windowIndicator["threshold"] = ind.Threshold
	windowIndicator["op"] = ind.Op
	windowIndicator["size"] = ind.Size

	return map[string]interface{}{
		//fieldNameRequestBasedEvaluation: []interface{}{},
		fieldNameWindowBasedEvaluation: []interface{}{windowIndicator},
	}, nil
}

func flattenSLOQueries(queries []SLIQueryGroup) ([]interface{}, error) {
	var queryList []interface{}

	for _, query := range queries {
		queryMap := map[string]interface{}{}
		queryMap["query_group_type"] = query.QueryGroupType
		queries, err := flattenSLOQueryGroup(query.QueryGroup)

		if err != nil {
			return nil, err
		}

		queryMap["query_group"] = queries

		queryList = append(queryList, queryMap)
	}

	return queryList, nil
}

func flattenSLOQueryGroup(queries []SLIQuery) ([]interface{}, error) {
	var queryList []interface{}

	for _, query := range queries {
		queryMap := map[string]interface{}{}
		queryMap["row_id"] = query.RowId
		queryMap["query"] = query.Query
		queryMap["use_row_count"] = query.UseRowCount
		queryMap["field"] = query.Field

		queryList = append(queryList, queryMap)
	}

	return queryList, nil
}

func resourceToSLO(d *schema.ResourceData) (*SLOLibrarySLO, error) {
	compliance := getSLOCompliance(d)
	indicator, err := getSLOIndicator(d)

	if err != nil {
		return nil, err
	}

	slo := SLOLibrarySLO{
		ID:          d.Id(),
		Name:        d.Get("name").(string),
		Description: d.Get("description").(string),
		Version:     d.Get("version").(int),
		CreatedAt:   d.Get("created_at").(string),
		CreatedBy:   d.Get("created_by").(string),
		ContentType: sloContentTypeString,
		ModifiedAt:  d.Get("modified_at").(string),
		ModifiedBy:  d.Get("modified_by").(string),
		ParentID:    d.Get("parent_id").(string),
		IsSystem:    d.Get("is_system").(bool),
		IsMutable:   d.Get("is_mutable").(bool),
		IsLocked:    d.Get("is_locked").(bool),
		SignalType:  d.Get("signal_type").(string),
		Compliance:  *compliance,
		Indicator:   *indicator,
		Service:     d.Get("service").(string),
		Application: d.Get("application").(string),
	}

	err = verifySLOObject(slo)
	if err != nil {
		return nil, err
	}

	return &slo, nil
}

func getSLOCompliance(d *schema.ResourceData) *SLOCompliance {
	complianceDict := d.Get("compliance").([]interface{})[0].(map[string]interface{})

	complianceType := complianceDict["compliance_type"].(string)

	startFrom := ""
	windowType := ""
	size := complianceDict["size"].(string)

	if complianceType == "Calendar" {
		// field windowType needs to be specified instead of `size` for calendar compliance
		windowType = size
		size = ""

		if complianceDict["start_from"] != nil {
			startFrom = complianceDict["start_from"].(string)
		}
	}

	return &SLOCompliance{
		ComplianceType: complianceType,
		Target:         complianceDict["target"].(float64),
		Timezone:       complianceDict["timezone"].(string),
		Size:           size,
		WindowType:     windowType,
		StartFrom:      startFrom,
	}
}

func getSLOIndicator(d *schema.ResourceData) (*SLOIndicator, error) {
	indicatorWrapperDict := d.Get("indicator").([]interface{})[0].(map[string]interface{})

	if indicatorWrapperDict[fieldNameRequestBasedEvaluation] != nil &&
		len(indicatorWrapperDict[fieldNameRequestBasedEvaluation].([]interface{})) > 0 {
		indicatorDict := indicatorWrapperDict[fieldNameRequestBasedEvaluation].([]interface{})[0].(map[string]interface{})
		queriesRaw := indicatorDict["queries"].([]interface{})

		threshold := float64(0)

		if indicatorDict["threshold"] != nil {
			threshold = indicatorDict["threshold"].(float64)
		}

		op := ""
		if indicatorDict["op"] != nil {
			op = indicatorDict["op"].(string)
		}

		return &SLOIndicator{
			EvaluationType: "Request",
			QueryType:      indicatorDict["query_type"].(string),
			Queries:        GetSLOIndicatorQueries(queriesRaw),
			Threshold:      threshold,
			Op:             op,
		}, nil
	}

	if indicatorWrapperDict[fieldNameWindowBasedEvaluation] != nil &&
		len(indicatorWrapperDict[fieldNameWindowBasedEvaluation].([]interface{})) > 0 {
		indicatorDict := indicatorWrapperDict[fieldNameWindowBasedEvaluation].([]interface{})[0].(map[string]interface{})
		queriesRaw := indicatorDict["queries"].([]interface{})

		aggregation := ""
		if indicatorDict["aggregation"] != nil {
			aggregation = indicatorDict["aggregation"].(string)
		}

		return &SLOIndicator{
			EvaluationType: "Window",
			QueryType:      indicatorDict["query_type"].(string),
			Queries:        GetSLOIndicatorQueries(queriesRaw),
			Threshold:      indicatorDict["threshold"].(float64),
			Op:             indicatorDict["op"].(string),
			Aggregation:    aggregation,
			Size:           indicatorDict["size"].(string),
		}, nil
	}

	return nil, fmt.Errorf("can't find indicator in resource, valid types are '%s' and '%s'", fieldNameRequestBasedEvaluation, fieldNameWindowBasedEvaluation)
}

func GetSLOIndicatorQueries(queriesRaw []interface{}) []SLIQueryGroup {

	queries := make([]SLIQueryGroup, len(queriesRaw))

	for i := range queries {
		qDict := queriesRaw[i].(map[string]interface{})

		queries[i].QueryGroupType = qDict["query_group_type"].(string)

		qGroupRaw := qDict["query_group"].([]interface{})
		qGroups := make([]SLIQuery, len(qGroupRaw))

		for j := range qGroups {
			qRaw := qGroupRaw[j].(map[string]interface{})
			field := ""

			if qRaw["field"] != nil {
				field = qRaw["field"].(string)
			}

			qGroup := SLIQuery{
				RowId:       qRaw["row_id"].(string),
				Query:       qRaw["query"].(string),
				Field:       field,
				UseRowCount: qRaw["use_row_count"].(bool),
			}

			qGroups[j] = qGroup
		}

		queries[i].QueryGroup = qGroups
	}

	return queries
}

func resourceSumologicSLOUpdate(d *schema.ResourceData, meta interface{}) error {
	c := meta.(*Client)
	slo, err := resourceToSLO(d)
	if err != nil {
		return err
	}

	slo.Type = "SlosLibrarySloUpdate"
	if d.HasChange("parent_id") {
		err := c.MoveSLOLibraryToFolder(slo.ID, slo.ParentID)
		if err != nil {
			return err
		}
	}
	sloUpdated, err := c.SLORead(d.Id())

	if err != nil {
		return err
	}
	slo.Version = sloUpdated.Version
	slo.ModifiedAt = sloUpdated.ModifiedAt

	err = c.UpdateSLO(*slo)
	if err != nil {
		return err
	}
	return resourceSLORead(d, meta)
}

func resourceSumologicSLODelete(d *schema.ResourceData, meta interface{}) error {

	c := meta.(*Client)
	slo, err := resourceToSLO(d)

	if err != nil {
		return err
	}

	err = c.DeleteSLO(slo.ID)
	if err != nil {
		return err
	}

	return nil
}

func verifySLOObject(slo SLOLibrarySLO) error {

	if slo.Name == "" {
		return fmt.Errorf("name is required")
	}

	for _, q := range slo.Indicator.Queries {
		for _, qg := range q.QueryGroup {
			if qg.Field != "" && qg.UseRowCount {
				return fmt.Errorf("'field' for the query can not be specfied when 'use_row_count' is true")
			}
		}
	}
	return nil
}
