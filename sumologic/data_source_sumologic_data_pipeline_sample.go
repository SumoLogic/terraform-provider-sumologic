package sumologic

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

var dataPipelineMessageSeparators = []string{
	`\n`,
	`\r\n`,
	`---`,
	`\|`,
	`(?<=\})(?=\{)`,
	`(?=\d{4}-\d{2}-\d{2}[T ])`,
	`(?=\w{3}\s+\d+\s+\d{2}:\d{2}:\d{2})`,
	`(empty)`,
}

func dataSourceSumologicDataPipelineTest() *schema.Resource {
	return &schema.Resource{
		Read:   dataSourceSumologicDataPipelineTestRead,
		Schema: dataSourceDataPipelineTestSchema(),
	}
}

func dataSourceDataPipelineTestSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"name": {
			Type:         schema.TypeString,
			Required:     true,
			ValidateFunc: validation.StringLenBetween(1, 128),
		},
		"description": {
			Type:         schema.TypeString,
			Optional:     true,
			ValidateFunc: validation.StringLenBetween(0, 1024),
		},
		"pipeline_type": {
			Type:         schema.TypeString,
			Optional:     true,
			Default:      "route_based",
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
		"node": {
			Type:     schema.TypeList,
			Required: true,
			MinItems: 1,
			Elem: &schema.Resource{
				Schema: dataPipelineNodeResourceSchema(),
			},
		},
		"sample_data": {
			Type:     schema.TypeString,
			Required: true,
		},
		"message_separator": {
			Type:         schema.TypeString,
			Optional:     true,
			ValidateFunc: validation.StringInSlice(dataPipelineMessageSeparators, false),
		},
		"metadata_fields": {
			Type:     schema.TypeMap,
			Optional: true,
			Elem:     &schema.Schema{Type: schema.TypeString},
		},
		"message": {
			Type:     schema.TypeList,
			Computed: true,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"message_index": {
						Type:     schema.TypeInt,
						Computed: true,
					},
					"input": {
						Type:     schema.TypeString,
						Computed: true,
					},
					"output": {
						Type:     schema.TypeString,
						Computed: true,
					},
					"truncated": {
						Type:     schema.TypeBool,
						Computed: true,
					},
					"node_result": {
						Type:     schema.TypeList,
						Computed: true,
						Elem: &schema.Resource{
							Schema: map[string]*schema.Schema{
								"node_name": {
									Type:     schema.TypeString,
									Computed: true,
								},
								"node_type": {
									Type:     schema.TypeString,
									Computed: true,
								},
								"matched": {
									Type:     schema.TypeBool,
									Computed: true,
								},
								"input": {
									Type:     schema.TypeString,
									Computed: true,
								},
								"output": {
									Type:     schema.TypeString,
									Computed: true,
								},
								"selected_route_condition": {
									Type:     schema.TypeString,
									Computed: true,
								},
								"selected_route_target": {
									Type:     schema.TypeString,
									Computed: true,
								},
								"processor_result": {
									Type:     schema.TypeList,
									Computed: true,
									Elem: &schema.Resource{
										Schema: map[string]*schema.Schema{
											"processor_name": {
												Type:     schema.TypeString,
												Computed: true,
											},
											"processor_type": {
												Type:     schema.TypeString,
												Computed: true,
											},
											"input": {
												Type:     schema.TypeString,
												Computed: true,
											},
											"output": {
												Type:     schema.TypeString,
												Computed: true,
											},
											"message_dropped": {
												Type:     schema.TypeBool,
												Computed: true,
											},
										},
									},
								},
							},
						},
					},
					"edge_result": {
						Type:     schema.TypeList,
						Computed: true,
						Elem: &schema.Resource{
							Schema: map[string]*schema.Schema{
								"source_node": {
									Type:     schema.TypeString,
									Computed: true,
								},
								"target_node": {
									Type:     schema.TypeString,
									Computed: true,
								},
								"condition_expression": {
									Type:     schema.TypeString,
									Computed: true,
								},
								"edge_order": {
									Type:     schema.TypeInt,
									Computed: true,
								},
							},
						},
					},
				},
			},
		},
	}
}

func dataSourceSumologicDataPipelineTestRead(d *schema.ResourceData, meta interface{}) error {
	c := meta.(*Client)

	request := expandDataPipelineSampleTestRequest(d)

	response, err := c.TestDataPipeline(request)
	if err != nil {
		return fmt.Errorf("error testing data pipeline %q: %v", request.Pipeline.Name, err)
	}

	id, err := dataPipelineSampleTestID(request)
	if err != nil {
		return err
	}
	d.SetId(id)

	if err := d.Set("message", flattenDataPipelineSampleTestMessages(response.Messages)); err != nil {
		return fmt.Errorf("error setting message for data pipeline test: %v", err)
	}

	return nil
}

func expandDataPipelineSampleTestRequest(d *schema.ResourceData) DataPipelineSampleTestRequest {
	metadataFields := map[string]string{}
	for k, v := range d.Get("metadata_fields").(map[string]interface{}) {
		metadataFields[k] = v.(string)
	}

	return DataPipelineSampleTestRequest{
		Pipeline:         expandDataPipelineRequest(d),
		SampleData:       d.Get("sample_data").(string),
		MessageSeparator: d.Get("message_separator").(string),
		MetadataFields:   metadataFields,
	}
}

func dataPipelineSampleTestID(request DataPipelineSampleTestRequest) (string, error) {
	body, err := json.Marshal(request)
	if err != nil {
		return "", err
	}
	sum := sha256.Sum256(body)
	return hex.EncodeToString(sum[:]), nil
}

func flattenDataPipelineSampleTestMessages(messages []DataPipelineSampleTestMessage) []interface{} {
	flattened := make([]interface{}, len(messages))
	for i, message := range messages {
		flattened[i] = map[string]interface{}{
			"message_index": message.MessageIndex,
			"input":         message.Input,
			"output":        message.Output,
			"truncated":     message.Truncated,
			"node_result":   flattenDataPipelineSampleTestNodes(message.NodeResults),
			"edge_result":   flattenDataPipelineSampleTestEdges(message.EdgeResults),
		}
	}
	return flattened
}

func flattenDataPipelineSampleTestNodes(nodes []DataPipelineSampleTestNode) []interface{} {
	flattened := make([]interface{}, len(nodes))
	for i, node := range nodes {
		processors := make([]interface{}, len(node.ProcessorResults))
		for j, processor := range node.ProcessorResults {
			processors[j] = map[string]interface{}{
				"processor_name":  processor.ProcessorName,
				"processor_type":  processor.ProcessorType,
				"input":           processor.Input,
				"output":          processor.Output,
				"message_dropped": processor.MessageDropped,
			}
		}

		flattened[i] = map[string]interface{}{
			"node_name":                node.NodeName,
			"node_type":                node.NodeType,
			"matched":                  node.Matched,
			"input":                    node.Input,
			"output":                   node.Output,
			"selected_route_condition": node.SelectedRouteCondition,
			"selected_route_target":    node.SelectedRouteTarget,
			"processor_result":         processors,
		}
	}
	return flattened
}

func flattenDataPipelineSampleTestEdges(edges []DataPipelineSampleTestEdge) []interface{} {
	flattened := make([]interface{}, len(edges))
	for i, edge := range edges {
		flattened[i] = map[string]interface{}{
			"source_node":          edge.SourceNode,
			"target_node":          edge.TargetNode,
			"condition_expression": edge.ConditionExpression,
			"edge_order":           edge.EdgeOrder,
		}
	}
	return flattened
}
