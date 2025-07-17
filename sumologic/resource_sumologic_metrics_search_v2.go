package sumologic

import (
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceSumologicMetricsSearchV2() *schema.Resource {
	return &schema.Resource{
		Create: resourceSumologicMetricsSearchV2Create,
		Read:   resourceSumologicMetricsSearchV2Read,
		Update: resourceSumologicMetricsSearchV2Update,
		Delete: resourceSumologicMetricsSearchV2Delete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"title": {
				Type:     schema.TypeString,
				Required: true,
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"time_range": {
				Type:     schema.TypeList,
				Required: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: GetTimeRangeSchema(),
				},
			},
			"queries": {
				Type:     schema.TypeList,
				Required: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"query_key": {
							Type:     schema.TypeString,
							Required: true,
						},
						"query_string": {
							Type:     schema.TypeString,
							Required: true,
						},
						"query_type": {
							Type:     schema.TypeString,
							Required: true,
						},
						"metrics_query_mode": {
							Type:     schema.TypeString,
							Optional: true,
						},
					},
				},
			},
			"visual_settings": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"folder_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
		},
	}
}

func resourceSumologicMetricsSearchV2Create(d *schema.ResourceData, meta interface{}) error {
	c := meta.(*Client)

	if d.Id() == "" {
		metricsSearchV2 := resourceToMetricsSearchV2(d)
		log.Println("=====================================================================")
		log.Printf("creating metrics search - %+v", metricsSearchV2)
		log.Println("=====================================================================")
		id, err := c.CreateMetricsSearchV2(metricsSearchV2)
		if err != nil {
			return err
		}

		d.SetId(id)
	}

	return resourceSumologicMetricsSearchV2Read(d, meta)
}

func resourceSumologicMetricsSearchV2Read(d *schema.ResourceData, meta interface{}) error {
	c := meta.(*Client)

	id := d.Id()
	metricsSearchV2, err := c.GetMetricsSearchV2(id)
	if err != nil {
		return err
	}

	if metricsSearchV2 == nil {
		log.Printf("[WARN] MetricsSearch not found, removing from state: %v - %v", id, err)
		d.SetId("")
		return nil
	}
	log.Println("=====================================================================")
	log.Printf("read metrics search - %+v", metricsSearchV2)
	log.Println("=====================================================================")

	return setMetricsSearchV2(d, metricsSearchV2)
}

func resourceSumologicMetricsSearchV2Update(d *schema.ResourceData, meta interface{}) error {
	c := meta.(*Client)

	metricsSearchV2 := resourceToMetricsSearchV2(d)
	log.Println("=====================================================================")
	log.Printf("updating metrics search - %+v", metricsSearchV2)
	log.Println("=====================================================================")
	err := c.UpdateMetricsSearchV2(metricsSearchV2)
	if err != nil {
		return err
	}

	return resourceSumologicMetricsSearchV2Read(d, meta)
}

func resourceSumologicMetricsSearchV2Delete(d *schema.ResourceData, meta interface{}) error {
	c := meta.(*Client)
	return c.DeleteMetricsSearchV2(d.Id())
}

func resourceToMetricsSearchV2(d *schema.ResourceData) MetricsSearchV2 {
	var timeRange interface{}
	if val, ok := d.GetOk("time_range"); ok {
		tfTimeRange := val.([]interface{})[0]
		timeRange = GetTimeRange(tfTimeRange.(map[string]interface{}))
	}

	var metricsQueries []MetricsSearchQueryV2
	if val, ok := d.GetOk("queries"); ok {
		metricsQueryData := val.([]interface{})
		for _, data := range metricsQueryData {
			metricsQueries = append(metricsQueries, resourceToMetricsSearchMetricsQueryV2([]interface{}{data}))
		}
	}

	return MetricsSearchV2{
		ID:             d.Id(),
		Title:          d.Get("title").(string),
		Description:    d.Get("description").(string),
		FolderID:       d.Get("folder_id").(string),
		VisualSettings: d.Get("visual_settings").(string),
		TimeRange:      timeRange,
		Queries:        metricsQueries,
	}
}

func resourceToMetricsSearchMetricsQueryV2(data interface{}) MetricsSearchQueryV2 {
	metricsSearchQuery := MetricsSearchQueryV2{}

	metricsSearchQuerySlice := data.([]interface{})
	if len(metricsSearchQuerySlice) > 0 {
		metricsSearchQueryObj := metricsSearchQuerySlice[0].(map[string]interface{})

		metricsSearchQuery.QueryKey = metricsSearchQueryObj["query_key"].(string)
		metricsSearchQuery.QueryString = metricsSearchQueryObj["query_string"].(string)
		metricsSearchQuery.QueryType = metricsSearchQueryObj["query_type"].(string)
		metricsSearchQuery.MetricsQueryMode = metricsSearchQueryObj["metrics_query_mode"].(string)
	}

	return metricsSearchQuery
}

func setMetricsSearchV2(d *schema.ResourceData, metricsSearchV2 *MetricsSearchV2) error {
	if err := d.Set("title", metricsSearchV2.Title); err != nil {
		return err
	}
	if err := d.Set("description", metricsSearchV2.Description); err != nil {
		return err
	}
	if err := d.Set("folder_id", metricsSearchV2.FolderID); err != nil {
		return err
	}
	if err := d.Set("visual_settings", metricsSearchV2.VisualSettings); err != nil {
		return err
	}

	metricsSearchQueriesV2 := make([]map[string]interface{}, len(metricsSearchV2.Queries))
	for i, searchQuery := range metricsSearchV2.Queries {
		metricsSearchQueriesV2[i] = getTerraformMetricsSearchQueryV2(searchQuery)
	}
	if err := d.Set("queries", metricsSearchQueriesV2); err != nil {
		return err
	}

	timeRange := GetTerraformTimeRange(metricsSearchV2.TimeRange.(map[string]interface{}))
	if err := d.Set("time_range", timeRange); err != nil {
		return err
	}

	return nil
}

func getTerraformMetricsSearchQueryV2(metricsSearchQueryV2 MetricsSearchQueryV2) map[string]interface{} {
	tfMetricsSearchQueryV2 := map[string]interface{}{}

	tfMetricsSearchQueryV2["query_key"] = metricsSearchQueryV2.QueryKey
	tfMetricsSearchQueryV2["query_string"] = metricsSearchQueryV2.QueryString
	tfMetricsSearchQueryV2["query_type"] = metricsSearchQueryV2.QueryType
	tfMetricsSearchQueryV2["metrics_query_mode"] = metricsSearchQueryV2.MetricsQueryMode
	return tfMetricsSearchQueryV2
}
