package sumologic

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/structure"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"encoding/json"
	"fmt"
	"log"
)

func resourceSumologicSourceTemplate() *schema.Resource {
	return &schema.Resource{
		Create: resourceSumologicSourceTemplateCreate,
		Read:   resourceSumologicSourceTemplateRead,
		Update: resourceSumologicSourceTemplateUpdate,
		Delete: resourceSumologicSourceTemplateDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"created_at": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"modified_by": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"created_by": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"modified_at": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"total_collector_linked": {
				Type:     schema.TypeInt,
				Computed: true,
			},

			"schema_ref": {
				Type:        schema.TypeList,
				Required:    true,
				MaxItems:    1,
				Description: "schema reference for source template.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{

						"type": {
							Type:     schema.TypeString,
							Required: true,
						},

						"version": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},

						"latest_version": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},

			"selector": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "Agent selector conditions",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{

						"tags": {
							Type:        schema.TypeList,
							Optional:    true,
							Description: "tags filter for agents",
							Elem: &schema.Schema{
								Type: schema.TypeList,
								Elem: &schema.Resource{
									Schema: map[string]*schema.Schema{

										"key": {
											Type:     schema.TypeString,
											Required: true,
										},

										"values": {
											Type:        schema.TypeList,
											Required:    true,
											Description: "values of the given tag.",
											Elem: &schema.Schema{
												Type: schema.TypeString,
											},
										},
									},
								},
							},
						},

						"names": {
							Type:        schema.TypeList,
							Optional:    true,
							Description: "names to select custom agents",
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
					},
				},
			},

			"config": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"input_json": {
				Type:             schema.TypeString,
				Description:      "inputJson of source template",
				ValidateFunc:     validation.StringIsJSON,
				Required:         true,
				DiffSuppressFunc: structure.SuppressJsonDiff,
			},
		},
	}
}

func resourceSumologicSourceTemplateRead(d *schema.ResourceData, meta interface{}) error {
	c := meta.(*Client)

	id := d.Id()
	sourceTemplate, err := c.GetSourceTemplate(id)
	if err != nil {
		return err
	}

	if sourceTemplate == nil {
		log.Printf("[WARN] SourceTemplate not found, removing from state: %v - %v", id, err)
		d.SetId("")
		return nil
	}

	if err := d.Set("schema_ref", schemaRefToList(&sourceTemplate.SchemaRef)); err != nil {
		return fmt.Errorf("Error setting schema_ref for resource %s: %s", d.Id(), err)
	}

	if err := d.Set("selector", selectorToList(&sourceTemplate.Selector)); err != nil {
		return fmt.Errorf("Error setting selector for resource %s: %s", d.Id(), err)
	}

	if err := d.Set("input_json", string(sourceTemplate.InputJson)); err != nil {
		return fmt.Errorf("Error setting input_json for resource %s: %s", d.Id(), err)
	}

	d.Set("created_by", sourceTemplate.CreatedBy)
	d.Set("modified_by", sourceTemplate.ModifiedBy)
	d.Set("created_at", sourceTemplate.CreatedAt)
	d.Set("modified_at", sourceTemplate.ModifiedAt)
	d.Set("config", sourceTemplate.Config)
	d.Set("total_collector_linked", sourceTemplate.TotalCollectorLinked)

	return nil
}

func resourceSumologicSourceTemplateUpdate(d *schema.ResourceData, meta interface{}) error {
	c := meta.(*Client)

	sourceTemplate := resourceToSourceTemplate(d)
	err := c.UpdateSourceTemplate(sourceTemplate)
	if err != nil {
		return err
	}

	return resourceSumologicSourceTemplateRead(d, meta)
}

func resourceSumologicSourceTemplateDelete(d *schema.ResourceData, meta interface{}) error {
	c := meta.(*Client)

	return c.DeleteSourceTemplate(d.Id())
}

func resourceSumologicSourceTemplateCreate(d *schema.ResourceData, meta interface{}) error {
	c := meta.(*Client)

	if d.Id() == "" {
		sourceTemplate := resourceToSourceTemplate(d)
		id, err := c.CreateSourceTemplate(sourceTemplate)
		if err != nil {
			return err
		}
		d.SetId(id)
	}

	return resourceSumologicSourceTemplateRead(d, meta)
}

func resourceToSourceTemplate(d *schema.ResourceData) SourceTemplate {

	schemaRef := resourceToSchemaRef(d.Get("schema_ref"))
	selector := resourceToSelector(d.Get("selector"))

	var jsonRawConf json.RawMessage
	conf := []byte(d.Get("input_json").(string))

	err := json.Unmarshal(conf, &jsonRawConf)
	if err != nil {
		log.Println("Unable to unmarshal the input json configuration")
		return SourceTemplate{}
	}

	return SourceTemplate{
		Config:               d.Get("config").(string),
		ID:                   d.Id(),
		SchemaRef:            schemaRef,
		TotalCollectorLinked: d.Get("total_collector_linked").(int),
		Selector:             selector,
		InputJson:            jsonRawConf,
		CreatedBy:            d.Get("created_by").(string),
		ModifiedBy:           d.Get("modified_by").(string),
		ModifiedAt:           d.Get("modified_at").(string),
		CreatedAt:            d.Get("created_at").(string),
	}
}

func resourceToSchemaRef(data interface{}) SchemaRef {

	schemaRefSlice := data.([]interface{})
	schemaRef := SchemaRef{}
	if len(schemaRefSlice) > 0 {
		schemaRefObj := schemaRefSlice[0].(map[string]interface{})
		schemaRef.Type = schemaRefObj["type"].(string)
		schemaRef.Version = schemaRefObj["version"].(string)
		schemaRef.LatestVersion = schemaRefObj["latest_version"].(string)
	}

	return schemaRef
}

func resourceToSelector(data interface{}) Selector {

	selectorSlice := data.([]interface{})
	selector := Selector{}
	if len(selectorSlice) > 0 {
		selectorObj := selectorSlice[0].(map[string]interface{})

		tagsData := selectorObj["tags"].([]interface{})
		tags := make([][]OtTag, len(tagsData))
		if len(tagsData) > 0 {
			for i, v := range tagsData {
				tags[i] = resourceToOtTag(v)
			}
		}
		selector.Tags = tags

		namesData := selectorObj["names"].([]interface{})
		names := make([]string, len(namesData))
		for i, v := range namesData {
			names[i] = v.(string)
		}
		selector.Names = names

	}

	return selector
}

func resourceToOtTag(data interface{}) []OtTag {

	otTagSlice := data.([]interface{})
	otTag := make([]OtTag, len(otTagSlice))
	if len(otTagSlice) > 0 {
		for i, v := range otTagSlice {
			otTagObj := v.(map[string]interface{})
			otTag[i].Key = otTagObj["key"].(string)

			valuesData := otTagObj["values"].([]interface{})
			values := make([]string, len(valuesData))
			for j, val := range valuesData {
				values[j] = val.(string)
			}
			otTag[i].Values = values
		}
	}

	return otTag
}

func selectorToList(selector *Selector) []map[string]interface{} {
	selectorList := make([]map[string]interface{}, 0, 1)

	if selector == nil {
		return selectorList
	}

	selectorObj := make(map[string]interface{})
	if len(selector.Tags) > 0 {
		var flattenedTags [][]map[string]interface{}
		for _, tagList := range selector.Tags {
			var tagElem []map[string]interface{}
			for _, tag := range tagList {
				tagElem = append(tagElem, map[string]interface{}{
					"key":    tag.Key,
					"values": tag.Values,
				})
			}
			flattenedTags = append(flattenedTags, tagElem)
		}
		selectorObj["tags"] = flattenedTags
	}
	if len(selector.Names) > 0 {
		selectorObj["names"] = selector.Names
	}

	selectorList = append(selectorList, selectorObj)
	return selectorList
}

func schemaRefToList(schema_ref *SchemaRef) []map[string]string {
	schemaList := make([]map[string]string, 0, 1)
	schemaObj := make(map[string]string)

	if schema_ref != nil {
		schemaObj["type"] = schema_ref.Type
		schemaObj["version"] = schema_ref.Version
		schemaObj["latest_version"] = schema_ref.LatestVersion
		schemaList = append(schemaList, schemaObj)
	}

	return schemaList
}
