package sumologic

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"log"
	"regexp"
)

const sloAggregationRegexString = `^(Avg|Min|Max|Sum|(p[5-9][0-9])(\.\d{1,3})?$)$` // TODO update it to allow lower pct values
const sloAggregationWindowRegexString = `^[0-9]{1,2}(m|h)$`                        // TODO make it exact of min 1m and max 1h

func resourceSumologicSLO() *schema.Resource {

	aggrRegex := regexp.MustCompile(sloAggregationRegexString)
	windowRegex := regexp.MustCompile(sloAggregationWindowRegexString)

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

	return &schema.Resource{
		Create: resourceSumologicSLOCreate,
		Read:   resourceSumologicSLORead,
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
								"1d", "2d", "3d", "4d", "5d", "6d", "7d", "8d", "9d", "10d", "11d", "12d", "13d", "14d",
							}, false),
						},
					},
				},
			},
			"indicator": {
				Type:     schema.TypeList,
				MaxItems: 1,
				Required: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"evaluation_type": {
							Type:     schema.TypeString,
							Required: true,
							ValidateFunc: validation.StringInSlice([]string{
								"Window",
								"Request",
							}, false),
						},
						"queries": {
							Type:     schema.TypeList,
							Required: true,
							MinItems: 1,
							MaxItems: 2,
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
						},
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
							Type:         schema.TypeString,
							Optional:     true,
							Default:      "Avg",
							ValidateFunc: validation.StringMatch(aggrRegex, `value must match : `+sloAggregationRegexString),
						},
						"size": {
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: validation.StringMatch(windowRegex, `value must match : `+sloAggregationWindowRegexString),
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
		slo := resourceToSLO(d)
		slo.Type = "SlosLibrarySlo"
		slo.ContentType = "Slo"
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
		sloDefinitionID, err := c.CreateSLO(slo, paramMap)
		if err != nil {
			return err
		}

		d.SetId(sloDefinitionID)
	}
	return resourceSumologicSLORead(d, meta)
}

func resourceSumologicSLORead(d *schema.ResourceData, meta interface{}) error {
	c := meta.(*Client)

	slo, err := c.SLORead(d.Id(), nil)
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
	d.Set("content_type", "Slo")
	d.Set("is_mutable", slo.IsMutable)
	d.Set("is_locked", slo.IsLocked)
	d.Set("is_system", slo.IsSystem)
	d.Set("service", slo.Service)
	d.Set("application", slo.Application)
	// set compliance
	//if err := d.Set("compliance", []SLOCompliance{slo.Compliance}); err != nil {
	//	return fmt.Errorf("error setting fields for resource %s: %s", d.Id(), err)
	//}
	//if err := d.Set("indicator", slo.Indicator); err != nil {
	//	return fmt.Errorf("error setting fields for resource %s: %s", d.Id(), err)
	//}

	return nil
}

func resourceToSLO(d *schema.ResourceData) SLOLibrarySLO {
	compliance := getSLOCompliance(d)
	indicator := getSLOIndicator(d)
	return SLOLibrarySLO{
		ID:          d.Id(),
		Name:        d.Get("name").(string),
		Description: d.Get("description").(string),
		Version:     d.Get("version").(int),
		CreatedAt:   d.Get("created_at").(string),
		CreatedBy:   d.Get("created_by").(string),
		ModifiedAt:  d.Get("modified_at").(string),
		ModifiedBy:  d.Get("modified_by").(string),
		ParentID:    d.Get("parent_id").(string),
		//ContentType: d.Get("content_type").(string),
		//Type:        d.Get("type").(string),
		IsSystem:    d.Get("is_system").(bool),
		IsMutable:   d.Get("is_mutable").(bool),
		IsLocked:    d.Get("is_locked").(bool),
		SignalType:  d.Get("signal_type").(string),
		Compliance:  compliance,
		Indicator:   indicator,
		Service:     d.Get("service").(string),
		Application: d.Get("application").(string),
	}
}

func getSLOCompliance(d *schema.ResourceData) SLOCompliance {
	complianceDict := d.Get("compliance").([]interface{})[0].(map[string]interface{})
	return SLOCompliance{
		ComplianceType: complianceDict["compliance_type"].(string),
		Target:         complianceDict["target"].(float64),
		Timezone:       complianceDict["timezone"].(string),
		Size:           complianceDict["size"].(string),
	}
}

func getSLOIndicator(d *schema.ResourceData) SLOIndicator {
	indicatorDict := d.Get("indicator").([]interface{})[0].(map[string]interface{})
	queriesRaw := indicatorDict["queries"].([]interface{})
	return SLOIndicator{
		EvaluationType: indicatorDict["evaluation_type"].(string),
		QueryType:      indicatorDict["query_type"].(string),
		Queries:        GetSLOIndicatorQueries(queriesRaw),
		Threshold:      indicatorDict["threshold"].(float64),
		Op:             indicatorDict["op"].(string),
		Aggregation:    indicatorDict["aggregation"].(string),
		Size:           indicatorDict["size"].(string),
	}
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
			qGroup := SLIQuery{
				RowId: qRaw["row_id"].(string),
				Query: qRaw["query"].(string),
			}
			qGroups[j] = qGroup
		}
		queries[i].QueryGroup = qGroups
	}

	return queries
}

func resourceSumologicSLOUpdate(d *schema.ResourceData, meta interface{}) error {
	c := meta.(*Client)
	slo := resourceToSLO(d)
	slo.Type = "SlosLibrarySloUpdate"
	if d.HasChange("parent_id") {
		err := c.MoveSLOLibraryToFolder(slo)
		if err != nil {
			return err
		}
	}

	err := c.UpdateSLO(slo)
	if err != nil {
		return err
	}
	return resourceSumologicSLORead(d, meta)
}

func resourceSumologicSLODelete(d *schema.ResourceData, meta interface{}) error {

	c := meta.(*Client)
	slo := resourceToSLO(d)
	err := c.DeleteSLO(slo.ID)
	if err != nil {
		return err
	}
	return nil
}
