package sumologic

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceSumologicDataPipeline() *schema.Resource {
	return &schema.Resource{
		Read:   dataSourceSumologicDataPipelineRead,
		Schema: dataSourceDataPipelineSchema(),
	}
}

func dataSourceDataPipelineSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"id": {
			Type:         schema.TypeString,
			Optional:     true,
			Computed:     true,
			ExactlyOneOf: []string{"id", "name"},
		},
		"name": {
			Type:         schema.TypeString,
			Optional:     true,
			Computed:     true,
			ExactlyOneOf: []string{"id", "name"},
		},
		"description": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"pipeline_type": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"route_expression": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"is_enabled": {
			Type:     schema.TypeBool,
			Computed: true,
		},
		"order": {
			Type:     schema.TypeInt,
			Computed: true,
		},
		"state": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"version": {
			Type:     schema.TypeInt,
			Computed: true,
		},
		"customer_id": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"created_at": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"modified_at": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"created_by_user_id": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"modified_by_user_id": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"node": {
			Type:     schema.TypeList,
			Computed: true,
			Elem: &schema.Resource{
				Schema: dataSourceDataPipelineNodeSchema(),
			},
		},
	}
}

func dataSourceDataPipelineNodeSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"id": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"name": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"node_type": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"filter_expression": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"is_enabled": {
			Type:     schema.TypeBool,
			Computed: true,
		},
		"processor": {
			Type:     schema.TypeList,
			Computed: true,
			Elem: &schema.Resource{
				Schema: dataSourceDataPipelineProcessorSchema(),
			},
		},
		"output": {
			Type:     schema.TypeList,
			Computed: true,
			Elem: &schema.Resource{
				Schema: dataSourceDataPipelineOutputSchema(),
			},
		},
	}
}

func dataSourceDataPipelineProcessorSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"id": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"name": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"processor_type": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"order": {
			Type:     schema.TypeInt,
			Computed: true,
		},
		"is_enabled": {
			Type:     schema.TypeBool,
			Computed: true,
		},
		"config": {
			Type:     schema.TypeString,
			Computed: true,
		},
	}
}

func dataSourceDataPipelineOutputSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"target": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"condition": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"order": {
			Type:     schema.TypeInt,
			Computed: true,
		},
	}
}

func dataSourceSumologicDataPipelineRead(d *schema.ResourceData, meta interface{}) error {
	c := meta.(*Client)

	var pipeline *DataPipeline
	var err error

	if id, ok := d.GetOk("id"); ok {
		pipeline, err = c.GetDataPipeline(id.(string))
		if err != nil {
			return fmt.Errorf("error retrieving data pipeline with id %v: %v", id, err)
		}
		if pipeline == nil {
			return fmt.Errorf("data pipeline with id %v not found", id)
		}
	} else {
		name := d.Get("name").(string)
		pipeline, err = c.FindDataPipelineByName(name)
		if err != nil {
			return fmt.Errorf("error retrieving data pipeline named %v: %v", name, err)
		}
		if pipeline == nil {
			return fmt.Errorf("data pipeline named %v not found", name)
		}
	}

	d.SetId(pipeline.ID)

	if err := setDataPipelineFields(d, pipeline); err != nil {
		return err
	}

	return nil
}

func setDataPipelineFields(d *schema.ResourceData, pipeline *DataPipeline) error {
	d.Set("name", pipeline.Name)
	d.Set("description", pipeline.Description)
	d.Set("pipeline_type", pipeline.PipelineType)
	d.Set("route_expression", pipeline.RouteExpression)
	d.Set("is_enabled", pipeline.IsEnabled)
	d.Set("order", pipeline.Order)
	d.Set("state", pipeline.State)
	d.Set("version", pipeline.Version)
	d.Set("customer_id", pipeline.CustomerID)
	d.Set("created_at", pipeline.CreatedAt)
	d.Set("modified_at", pipeline.ModifiedAt)
	d.Set("created_by_user_id", pipeline.CreatedByUserID)
	d.Set("modified_by_user_id", pipeline.ModifiedByUserID)

	if err := d.Set("node", flattenDataPipelineNodes(pipeline.Nodes)); err != nil {
		return fmt.Errorf("error setting node for data pipeline %s: %v", pipeline.ID, err)
	}

	return nil
}

func flattenDataPipelineNodes(nodes []DataPipelineNode) []interface{} {
	if nodes == nil {
		return []interface{}{}
	}

	flattened := make([]interface{}, len(nodes))
	for i, node := range nodes {
		isEnabled := false
		if node.IsEnabled != nil {
			isEnabled = *node.IsEnabled
		}

		flattened[i] = map[string]interface{}{
			"id":                node.ID,
			"name":              node.Name,
			"node_type":         node.NodeType,
			"filter_expression": node.FilterExpression,
			"is_enabled":        isEnabled,
			"processor":         flattenDataPipelineProcessors(node.Processors),
			"output":            flattenDataPipelineOutputs(node.Outputs),
		}
	}

	return flattened
}

func flattenDataPipelineProcessors(processors []DataPipelineProcessor) []interface{} {
	if processors == nil {
		return []interface{}{}
	}

	flattened := make([]interface{}, len(processors))
	for i, processor := range processors {
		config := ""
		if len(processor.Config) > 0 {
			config = string(processor.Config)
		}

		flattened[i] = map[string]interface{}{
			"id":             processor.ID,
			"name":           processor.Name,
			"processor_type": processor.ProcessorType,
			"order":          processor.Order,
			"is_enabled":     processor.IsEnabled,
			"config":         config,
		}
	}

	return flattened
}

func flattenDataPipelineOutputs(outputs []DataPipelineOutput) []interface{} {
	if outputs == nil {
		return []interface{}{}
	}

	flattened := make([]interface{}, len(outputs))
	for i, output := range outputs {
		order := 0
		if output.Order != nil {
			order = *output.Order
		}

		flattened[i] = map[string]interface{}{
			"target":    output.Target,
			"condition": output.Condition,
			"order":     order,
		}
	}

	return flattened
}
