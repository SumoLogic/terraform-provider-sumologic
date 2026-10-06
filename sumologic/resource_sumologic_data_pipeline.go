package sumologic

import (
	"encoding/json"
	"fmt"
	"sort"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/structure"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func resourceSumologicDataPipeline() *schema.Resource {
	return &schema.Resource{
		Create: resourceSumologicDataPipelineCreate,
		Read:   resourceSumologicDataPipelineRead,
		Update: resourceSumologicDataPipelineUpdate,
		Delete: resourceSumologicDataPipelineDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringLenBetween(1, 128),
			},
			"description": {
				Type:         schema.TypeString,
				Optional:     true,
				Default:      "",
				ValidateFunc: validation.StringLenBetween(0, 1024),
			},
			"pipeline_type": {
				Type:         schema.TypeString,
				Optional:     true,
				Default:      "route_based",
				ForceNew:     true,
				ValidateFunc: validation.StringInSlice([]string{"route_based"}, false),
			},
			"route_expression": {
				Type:     schema.TypeString,
				Required: true,
			},
			"is_enabled": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  true,
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
			"node": {
				Type:     schema.TypeList,
				Required: true,
				MinItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"name": {
							Type:     schema.TypeString,
							Required: true,
						},
						"node_type": {
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: validation.StringInSlice([]string{"source", "router", "destination", "processing_group"}, false),
						},
						"filter_expression": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"is_enabled": {
							Type:     schema.TypeBool,
							Optional: true,
							Default:  true,
						},
						"output": {
							Type:     schema.TypeList,
							Optional: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"target": {
										Type:     schema.TypeString,
										Required: true,
									},
									"condition": {
										Type:     schema.TypeString,
										Optional: true,
									},
									"order": {
										Type:     schema.TypeInt,
										Optional: true,
									},
								},
							},
						},
						"processor": {
							Type:     schema.TypeList,
							Optional: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"id": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"name": {
										Type:     schema.TypeString,
										Required: true,
									},
									"processor_type": {
										Type:     schema.TypeString,
										Required: true,
									},
									"is_enabled": {
										Type:     schema.TypeBool,
										Optional: true,
										Default:  true,
									},
									"config": {
										Type:             schema.TypeString,
										Required:         true,
										ValidateFunc:     validation.StringIsJSON,
										DiffSuppressFunc: structure.SuppressJsonDiff,
									},
								},
							},
						},
					},
				},
			},
		},
	}
}

func resourceSumologicDataPipelineCreate(d *schema.ResourceData, meta interface{}) error {
	c := meta.(*Client)

	request := expandDataPipelineRequest(d)

	pipeline, err := c.CreateDataPipeline(request)
	if err != nil {
		return fmt.Errorf("error creating data pipeline %q: %v", request.Name, err)
	}

	d.SetId(pipeline.ID)

	if _, err := c.PublishDataPipeline(d.Id()); err != nil {
		return fmt.Errorf("error publishing data pipeline %s: %v", d.Id(), err)
	}

	if err := c.SetDataPipelineEnabled(d.Id(), request.IsEnabled); err != nil {
		return fmt.Errorf("error setting enabled state for data pipeline %s: %v", d.Id(), err)
	}

	return resourceSumologicDataPipelineRead(d, meta)
}

func resourceSumologicDataPipelineRead(d *schema.ResourceData, meta interface{}) error {
	c := meta.(*Client)

	pipeline, err := c.GetDataPipeline(d.Id())
	if err != nil {
		return fmt.Errorf("error retrieving data pipeline %s: %v", d.Id(), err)
	}
	if pipeline == nil {
		d.SetId("")
		return nil
	}

	pipeline.Nodes = reorderDataPipelineNodes(pipeline.Nodes, dataPipelineNodeNameOrder(d, pipeline))

	for i, node := range pipeline.Nodes {
		if node.NodeType != "processing_group" {
			continue
		}
		pipeline.Nodes[i].Processors = reorderDataPipelineProcessors(
			node.Processors,
			dataPipelineProcessorNameOrder(d, node),
		)
	}

	return setDataPipelineResourceFields(d, pipeline)
}

func dataPipelineProcessorNameOrder(d *schema.ResourceData, node DataPipelineNode) []string {
	raw := d.Get("node").([]interface{})
	for _, r := range raw {
		m := r.(map[string]interface{})
		if m["name"].(string) != node.Name {
			continue
		}
		processors := m["processor"].([]interface{})
		if len(processors) == 0 {
			break
		}
		names := make([]string, len(processors))
		for i, p := range processors {
			names[i] = p.(map[string]interface{})["name"].(string)
		}
		return names
	}

	sorted := append([]DataPipelineProcessor{}, node.Processors...)
	sort.SliceStable(sorted, func(i, j int) bool {
		return sorted[i].Order < sorted[j].Order
	})
	names := make([]string, len(sorted))
	for i, p := range sorted {
		names[i] = p.Name
	}
	return names
}

func reorderDataPipelineProcessors(processors []DataPipelineProcessor, referenceOrder []string) []DataPipelineProcessor {
	byName := make(map[string]DataPipelineProcessor, len(processors))
	for _, p := range processors {
		byName[p.Name] = p
	}

	ordered := make([]DataPipelineProcessor, 0, len(processors))
	seen := make(map[string]bool, len(processors))
	for _, name := range referenceOrder {
		if p, ok := byName[name]; ok && !seen[name] {
			ordered = append(ordered, p)
			seen[name] = true
		}
	}
	for _, p := range processors {
		if !seen[p.Name] {
			ordered = append(ordered, p)
			seen[p.Name] = true
		}
	}
	return ordered
}

func dataPipelineNodeNameOrder(d *schema.ResourceData, pipeline *DataPipeline) []string {
	raw := d.Get("node").([]interface{})
	if len(raw) > 0 {
		names := make([]string, len(raw))
		for i, r := range raw {
			names[i] = r.(map[string]interface{})["name"].(string)
		}
		return names
	}
	return topologicalDataPipelineNodeOrder(pipeline.Nodes)
}

func topologicalDataPipelineNodeOrder(nodes []DataPipelineNode) []string {
	byName := make(map[string]DataPipelineNode, len(nodes))
	var source string
	for _, n := range nodes {
		byName[n.Name] = n
		if n.NodeType == "source" {
			source = n.Name
		}
	}

	order := make([]string, 0, len(nodes))
	seen := make(map[string]bool, len(nodes))

	queue := []string{}
	if source != "" {
		queue = append(queue, source)
	}
	for len(queue) > 0 {
		name := queue[0]
		queue = queue[1:]
		if seen[name] {
			continue
		}
		seen[name] = true
		order = append(order, name)

		outputs := append([]DataPipelineOutput{}, byName[name].Outputs...)
		sort.SliceStable(outputs, func(i, j int) bool {
			oi, oj := outputs[i].Order, outputs[j].Order
			if oi == nil {
				return false
			}
			if oj == nil {
				return true
			}
			return *oi < *oj
		})
		for _, o := range outputs {
			queue = append(queue, o.Target)
		}
	}

	for _, n := range nodes {
		if !seen[n.Name] {
			order = append(order, n.Name)
			seen[n.Name] = true
		}
	}

	return order
}

func reorderDataPipelineNodes(nodes []DataPipelineNode, referenceOrder []string) []DataPipelineNode {
	byName := make(map[string]DataPipelineNode, len(nodes))
	for _, n := range nodes {
		byName[n.Name] = n
	}

	ordered := make([]DataPipelineNode, 0, len(nodes))
	seen := make(map[string]bool, len(nodes))
	for _, name := range referenceOrder {
		if n, ok := byName[name]; ok && !seen[name] {
			ordered = append(ordered, n)
			seen[name] = true
		}
	}
	for _, n := range nodes {
		if !seen[n.Name] {
			ordered = append(ordered, n)
			seen[n.Name] = true
		}
	}
	return ordered
}

func resourceSumologicDataPipelineUpdate(d *schema.ResourceData, meta interface{}) error {
	c := meta.(*Client)

	request := DataPipelineUpdateRequest{
		DataPipelineRequest: expandDataPipelineRequest(d),
		Version:             d.Get("version").(int),
	}

	if _, err := c.UpdateDataPipeline(d.Id(), request); err != nil {
		return fmt.Errorf("error updating data pipeline %s: %v", d.Id(), err)
	}

	if _, err := c.PublishDataPipeline(d.Id()); err != nil {
		return fmt.Errorf("error publishing data pipeline %s: %v", d.Id(), err)
	}

	if err := c.SetDataPipelineEnabled(d.Id(), request.IsEnabled); err != nil {
		return fmt.Errorf("error setting enabled state for data pipeline %s: %v", d.Id(), err)
	}

	if d.HasChanges("name", "description") {
		if _, err := c.UpdateDataPipelineMetadata(d.Id(), DataPipelineMetadataRequest{
			Name:        d.Get("name").(string),
			Description: d.Get("description").(string),
		}); err != nil {
			return fmt.Errorf("error updating metadata for data pipeline %s: %v", d.Id(), err)
		}
	}

	return resourceSumologicDataPipelineRead(d, meta)
}

func resourceSumologicDataPipelineDelete(d *schema.ResourceData, meta interface{}) error {
	c := meta.(*Client)

	if err := c.DeleteDataPipeline(d.Id()); err != nil {
		return fmt.Errorf("error deleting data pipeline %s: %v", d.Id(), err)
	}

	return nil
}

func expandDataPipelineRequest(d *schema.ResourceData) DataPipelineRequest {
	return DataPipelineRequest{
		Name:            d.Get("name").(string),
		Description:     d.Get("description").(string),
		PipelineType:    d.Get("pipeline_type").(string),
		IsEnabled:       d.Get("is_enabled").(bool),
		RouteExpression: d.Get("route_expression").(string),
		Nodes:           expandDataPipelineNodes(d.Get("node").([]interface{})),
	}
}

func expandDataPipelineNodes(raw []interface{}) []DataPipelineNode {
	nodes := make([]DataPipelineNode, len(raw))
	for i, r := range raw {
		m := r.(map[string]interface{})
		isEnabled := m["is_enabled"].(bool)
		nodes[i] = DataPipelineNode{
			ID:               m["id"].(string),
			Name:             m["name"].(string),
			NodeType:         m["node_type"].(string),
			FilterExpression: m["filter_expression"].(string),
			IsEnabled:        &isEnabled,
			Outputs:          expandDataPipelineOutputs(m["output"].([]interface{})),
			Processors:       expandDataPipelineProcessors(m["processor"].([]interface{})),
		}
	}
	return nodes
}

func expandDataPipelineProcessors(raw []interface{}) []DataPipelineProcessor {
	processors := make([]DataPipelineProcessor, len(raw))
	for i, r := range raw {
		m := r.(map[string]interface{})

		var config json.RawMessage
		json.Unmarshal([]byte(m["config"].(string)), &config)

		processors[i] = DataPipelineProcessor{
			ID:            m["id"].(string),
			Name:          m["name"].(string),
			ProcessorType: m["processor_type"].(string),
			Order:         i + 1,
			IsEnabled:     m["is_enabled"].(bool),
			Config:        config,
		}
	}
	return processors
}

func expandDataPipelineOutputs(raw []interface{}) []DataPipelineOutput {
	outputs := make([]DataPipelineOutput, len(raw))
	for i, r := range raw {
		m := r.(map[string]interface{})
		outputs[i] = DataPipelineOutput{
			Target:    m["target"].(string),
			Condition: m["condition"].(string),
		}
		if order := m["order"].(int); order != 0 {
			outputs[i].Order = &order
		}
	}
	return outputs
}

func setDataPipelineResourceFields(d *schema.ResourceData, pipeline *DataPipeline) error {
	d.Set("name", pipeline.Name)
	d.Set("description", pipeline.Description)
	d.Set("pipeline_type", pipeline.PipelineType)
	d.Set("route_expression", pipeline.RouteExpression)
	d.Set("is_enabled", pipeline.IsEnabled)
	d.Set("order", pipeline.Order)
	d.Set("state", pipeline.State)
	d.Set("version", pipeline.Version)

	if err := d.Set("node", flattenDataPipelineResourceNodes(pipeline.Nodes)); err != nil {
		return fmt.Errorf("error setting node for data pipeline %s: %v", pipeline.ID, err)
	}

	return nil
}

func flattenDataPipelineResourceNodes(nodes []DataPipelineNode) []interface{} {
	flattened := make([]interface{}, len(nodes))
	for i, node := range nodes {
		isEnabled := true
		if node.IsEnabled != nil {
			isEnabled = *node.IsEnabled
		}
		flattened[i] = map[string]interface{}{
			"id":                node.ID,
			"name":              node.Name,
			"node_type":         node.NodeType,
			"filter_expression": node.FilterExpression,
			"is_enabled":        isEnabled,
			"output":            flattenDataPipelineOutputs(node.Outputs),
			"processor":         flattenDataPipelineResourceProcessors(node.Processors),
		}
	}
	return flattened
}

func flattenDataPipelineResourceProcessors(processors []DataPipelineProcessor) []interface{} {
	flattened := make([]interface{}, len(processors))
	for i, processor := range processors {
		flattened[i] = map[string]interface{}{
			"id":             processor.ID,
			"name":           processor.Name,
			"processor_type": processor.ProcessorType,
			"is_enabled":     processor.IsEnabled,
			"config":         string(processor.Config),
		}
	}
	return flattened
}
