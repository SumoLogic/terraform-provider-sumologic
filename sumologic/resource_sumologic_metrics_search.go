package sumologic

import (
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceSumologicMetricsSearch() *schema.Resource {
	return &schema.Resource{
		Create: resourceSumologicMetricsSearchCreate,
		Read:   resourceSumologicMetricsSearchRead,
		Update: resourceSumologicMetricsSearchUpdate,
		Delete: resourceSumologicMetricsSearchDelete,
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
			"parent_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"log_query": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"desired_quantization_in_secs": {
				Type:     schema.TypeInt,
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
			"properties": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"metrics_queries": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"row_id": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"query": {
							Type:     schema.TypeString,
							Optional: true,
						},
					},
				},
			},
		},
	}
}

func resourceSumologicMetricsSearchCreate(d *schema.ResourceData, meta interface{}) error {
	c := meta.(*Client)

	if d.Id() == "" {
		metricsSearch := resourceToMetricsSearch(d)
		log.Println("=====================================================================")
		log.Printf("creating metrics search - %+v", metricsSearch)
		log.Println("=====================================================================")
		id, err := c.CreateMetricsSearch(metricsSearch)
		if err != nil {
			return err
		}

		d.SetId(id)
	}

	return resourceSumologicMetricsSearchRead(d, meta)
}

func resourceSumologicMetricsSearchRead(d *schema.ResourceData, meta interface{}) error {
	c := meta.(*Client)

	id := d.Id()
	metricsSearch, err := c.GetMetricsSearch(id)
	if err != nil {
		return err
	}

	if metricsSearch == nil {
		log.Printf("[WARN] MetricsSearch not found, removing from state: %v - %v", id, err)
		d.SetId("")
		return nil
	}
	log.Println("=====================================================================")
	log.Printf("read metrics search - %+v", metricsSearch)
	log.Println("=====================================================================")

	return setMetricsSearch(d, metricsSearch)
}

func resourceSumologicMetricsSearchUpdate(d *schema.ResourceData, meta interface{}) error {
	c := meta.(*Client)

	metricsSearch := resourceToMetricsSearch(d)
	log.Println("=====================================================================")
	log.Printf("updating metrics search - %+v", metricsSearch)
	log.Println("=====================================================================")
	err := c.UpdateMetricsSearch(metricsSearch)
	if err != nil {
		return err
	}

	return resourceSumologicMetricsSearchRead(d, meta)
}

func resourceSumologicMetricsSearchDelete(d *schema.ResourceData, meta interface{}) error {
	c := meta.(*Client)
	return c.DeleteMetricsSearch(d.Id())
}

func setMetricsSearch(d *schema.ResourceData, metricsSearch *MetricsSearch) error {
	if err := d.Set("title", metricsSearch.Title); err != nil {
		return err
	}
	if err := d.Set("description", metricsSearch.Description); err != nil {
		return err
	}
	if err := d.Set("parent_id", metricsSearch.ParentId); err != nil {
		return err
	}
	if err := d.Set("log_query", metricsSearch.LogQuery); err != nil {
		return err
	}
	if err := d.Set("desired_quantization_in_secs", metricsSearch.DesiredQuantizationInSecs); err != nil {
		return err
	}
	if err := d.Set("properties", metricsSearch.Properties); err != nil {
		return err
	}

	metricsSearchQueries := make([]map[string]interface{}, len(metricsSearch.MetricsQueries))
	for i, searchQuery := range metricsSearch.MetricsQueries {
		metricsSearchQueries[i] = getTerraformMetricsSearchQuery(searchQuery)
	}
	if err := d.Set("metrics_queries", metricsSearchQueries); err != nil {
		return err
	}

	timeRange := GetTerraformTimeRange(metricsSearch.TimeRange.(map[string]interface{}))
	if err := d.Set("time_range", timeRange); err != nil {
		return err
	}

	return nil
}

func getTerraformMetricsSearchQuery(metricsSearchQuery MetricsSearchQuery) map[string]interface{} {
	tfMetricsSearchQuery := map[string]interface{}{}

	tfMetricsSearchQuery["row_id"] = metricsSearchQuery.RowId
	tfMetricsSearchQuery["query"] = metricsSearchQuery.Query
	return tfMetricsSearchQuery
}

func resourceToMetricsSearch(d *schema.ResourceData) MetricsSearch {
	var timeRange interface{}
	if val, ok := d.GetOk("time_range"); ok {
		tfTimeRange := val.([]interface{})[0]
		timeRange = GetTimeRange(tfTimeRange.(map[string]interface{}))
	}

	var metricsQueries []MetricsSearchQuery
	if val, ok := d.GetOk("metrics_queries"); ok {
		metricsQueryData := val.([]interface{})
		for _, data := range metricsQueryData {
			metricsQueries = append(metricsQueries, resourceToMetricsSearchMetricsQuery([]interface{}{data}))
		}
	}

	return MetricsSearch{
		ID:                        d.Id(),
		Title:                     d.Get("title").(string),
		Description:               d.Get("description").(string),
		ParentId:                  d.Get("parent_id").(string),
		LogQuery:                  d.Get("log_query").(string),
		DesiredQuantizationInSecs: d.Get("desired_quantization_in_secs").(int),
		TimeRange:                 timeRange,
		Properties:                d.Get("properties").(string),
		MetricsQueries:            metricsQueries,
	}
}

func resourceToMetricsSearchMetricsQuery(data interface{}) MetricsSearchQuery {
	metricsSearchQuery := MetricsSearchQuery{}

	metricsSearchQuerySlice := data.([]interface{})
	if len(metricsSearchQuerySlice) > 0 {
		metricsSearchQueryObj := metricsSearchQuerySlice[0].(map[string]interface{})

		metricsSearchQuery.RowId = metricsSearchQueryObj["row_id"].(string)
		metricsSearchQuery.Query = metricsSearchQueryObj["query"].(string)
	}

	return metricsSearchQuery
}
