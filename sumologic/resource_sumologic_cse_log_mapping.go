package sumologic

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"log"
)

func resourceSumologicCSELogMapping() *schema.Resource {
	return &schema.Resource{
		Create: resourceSumologicCSELogMappingCreate,
		Read:   resourceSumologicCSELogMappingRead,
		Delete: resourceSumologicCSELogMappingDelete,
		Update: resourceSumologicCSELogMappingUpdate,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"parent_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"product_guid": {
				Type:     schema.TypeString,
				Required: true,
			},
			"record_type": {
				Type:     schema.TypeString,
				Required: true,
			},
			"enabled": {
				Type:     schema.TypeBool,
				Required: true,
			},
			"relates_entities": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"skipped_values": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"fields": {
				Type:     schema.TypeList,
				Required: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{

						"name": {
							Type:     schema.TypeString,
							Required: true,
						},
						"value": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"value_type": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"skipped_values": {
							Type:     schema.TypeList,
							Optional: true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
						"default_value": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"format": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"case_insensitive": {
							Type:     schema.TypeBool,
							Optional: true,
						},
						"alternate_values": {
							Type:     schema.TypeList,
							Optional: true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
						"time_zone": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"split_delimiter": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"split_index": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"field_join": {
							Type:     schema.TypeList,
							Optional: true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
						"join_delimiter": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"format_parameters": {
							Type:     schema.TypeList,
							Optional: true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
						"lookup": {
							Type:     schema.TypeList,
							Optional: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{

									"key": {
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
					},
				},
			},
			"structured_inputs": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: getLogMappingStructuredInputFieldSchema(),
				},
			},
			"unstructured_fields": {
				Type:     schema.TypeList,
				MaxItems: 1,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{

						"pattern_names": {
							Type:     schema.TypeList,
							Required: true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
					},
				},
			},
		},
	}
}

func getLogMappingStructuredInputFieldSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"event_id_pattern": {
			Type:     schema.TypeString,
			Required: true,
		},
		"log_format": {
			Type:     schema.TypeString,
			Required: true,
		},
		"product": {
			Type:     schema.TypeString,
			Required: true,
		},
		"vendor": {
			Type:     schema.TypeString,
			Required: true,
		},
	}
}

func resourceSumologicCSELogMappingRead(d *schema.ResourceData, meta interface{}) error {
	c := meta.(*Client)

	id := d.Id()

	CSELogMapping, err := c.GetCSELogMapping(id)
	if err != nil {
		log.Printf("[WARN] CSE Insights Status not found when looking by id: %s, err: %v", id, err)
	}

	if CSELogMapping == nil {
		log.Printf("[WARN] CSE Insights Status not found, removing from state: %v - %v", id, err)
		d.SetId("")
		return nil
	}

	setLogMapping(d, CSELogMapping)

	return nil
}

func resourceSumologicCSELogMappingDelete(d *schema.ResourceData, meta interface{}) error {
	c := meta.(*Client)

	return c.DeleteCSELogMapping(d.Id())

}

func resourceSumologicCSELogMappingCreate(d *schema.ResourceData, meta interface{}) error {
	c := meta.(*Client)

	if d.Id() == "" {

		CSELogMapping := resourceToCSELogMapping(d)

		id, err := c.CreateCSELogMapping(CSELogMapping)
		if err != nil {
			return err
		}

		d.SetId(id)
	}

	return resourceSumologicCSELogMappingUpdate(d, meta)
}

func resourceSumologicCSELogMappingUpdate(d *schema.ResourceData, meta interface{}) error {
	c := meta.(*Client)

	CSELogMapping := resourceToCSELogMapping(d)
	CSELogMapping.ID = d.Id()
	err := c.UpdateCSELogMapping(CSELogMapping)
	if err != nil {
		return err
	}

	return resourceSumologicCSELogMappingRead(d, meta)
}

func resourceStringArrayToStringArray(values []interface{}) []string {
	strings := make([]string, len(values))
	for i, s := range values {
		strings[i] = s.(string)
	}
	return strings
}

func resourceToCSELogMapping(d *schema.ResourceData) CSELogMapping {

	skippedValuesData := d.Get("skipped_values").([]interface{})
	skippedValues := make([]string, len(skippedValuesData))

	for i, v := range skippedValuesData {
		skippedValues[i] = v.(string)
	}

	fieldsData := d.Get("fields").([]interface{})
	var fields []CSELogMappingField
	for _, data := range fieldsData {
		fields = append(fields, resourceToCSELogMappingField([]interface{}{data}))
	}

	structuredInputsData := d.Get("structured_inputs").([]interface{})
	var structuredInputs []CSELogMappingStructuredInputField
	for _, data := range structuredInputsData {
		structuredInputs = append(structuredInputs, resourceToCSELogMappingStructuredInputField([]interface{}{data}))
	}

	unstructuredFields := resourceToCSELogMappingUnstructuredInputField(d.Get("unstructured_fields").([]interface{}))

	enabled := d.Get("enabled").(bool)

	CSELogMapping := CSELogMapping{
		Name:            d.Get("name").(string),
		ParentId:        d.Get("parent_id").(string),
		ProductGuid:     d.Get("product_guid").(string),
		RecordType:      d.Get("record_type").(string),
		RelatesEntities: d.Get("relates_entities").(bool),
		SkippedValues:   skippedValues,
		Fields:          fields,
		Enabled:         &enabled,
	}

	if len(structuredInputs) > 0 {
		CSELogMapping.StructuredInputs = structuredInputs
	}
	if len(unstructuredFields.PatternNames) > 0 {
		CSELogMapping.UnstructuredFields = &unstructuredFields
	} else {
		CSELogMapping.UnstructuredFields = nil
	}

	return CSELogMapping
}

func resourceToCSELogMappingStructuredInputField(data interface{}) CSELogMappingStructuredInputField {
	structuredInputFieldSlice := data.([]interface{})
	structuredInputField := CSELogMappingStructuredInputField{}
	if len(structuredInputFieldSlice) > 0 {
		structuredInputFieldObj := structuredInputFieldSlice[0].(map[string]interface{})
		structuredInputField.EventIdPattern = structuredInputFieldObj["event_id_pattern"].(string)
		structuredInputField.LogFormat = structuredInputFieldObj["log_format"].(string)
		structuredInputField.Product = structuredInputFieldObj["product"].(string)
		structuredInputField.Vendor = structuredInputFieldObj["vendor"].(string)
	}
	return structuredInputField
}

func resourceToCSELogMappingUnstructuredInputField(data interface{}) CSELogMappingUnstructuredFields {
	unstructuredInputFieldSlice := data.([]interface{})
	unstructuredInputField := CSELogMappingUnstructuredFields{}
	if len(unstructuredInputFieldSlice) > 0 {
		unstructuredInputFieldObj := unstructuredInputFieldSlice[0].(map[string]interface{})
		unstructuredInputField.PatternNames = resourceStringArrayToStringArray(unstructuredInputFieldObj["pattern_names"].([]interface{}))
	}
	return unstructuredInputField
}

func resourceToCSELogMappingLookUp(data interface{}) CSELogMappingLookUp {
	lookUpSlice := data.([]interface{})
	lookup := CSELogMappingLookUp{}
	if len(lookUpSlice) > 0 {
		lookUpObj := lookUpSlice[0].(map[string]interface{})
		lookup.Key = lookUpObj["key"].(string)
		lookup.Value = lookUpObj["value"].(string)
	}
	return lookup
}

func resourceToCSELogMappingField(data interface{}) CSELogMappingField {
	fieldsSlice := data.([]interface{})
	field := CSELogMappingField{}
	if len(fieldsSlice) > 0 {
		fieldObj := fieldsSlice[0].(map[string]interface{})
		field.Name = fieldObj["name"].(string)
		field.Value = fieldObj["value"].(string)
		field.ValueType = fieldObj["value_type"].(string)
		field.SkippedValues = resourceStringArrayToStringArray(fieldObj["skipped_values"].([]interface{}))
		field.DefaultValue = fieldObj["default_value"].(string)
		field.Format = fieldObj["format"].(string)
		field.CaseInsensitive = fieldObj["case_insensitive"].(bool)
		field.AlternateValues = resourceStringArrayToStringArray(fieldObj["alternate_values"].([]interface{}))
		field.TimeZone = fieldObj["time_zone"].(string)
		field.SplitDelimiter = fieldObj["split_delimiter"].(string)
		field.SplitIndex = fieldObj["split_index"].(string)
		field.FieldJoin = resourceStringArrayToStringArray(fieldObj["field_join"].([]interface{}))
		field.JoinDelimiter = fieldObj["join_delimiter"].(string)
		field.FormatParameters = resourceStringArrayToStringArray(fieldObj["format_parameters"].([]interface{}))

		lookUpData := fieldObj["lookup"].([]interface{})
		var lookup []CSELogMappingLookUp
		for _, data := range lookUpData {
			lookup = append(lookup, resourceToCSELogMappingLookUp([]interface{}{data}))
		}

		field.LookUp = lookup

	}
	return field
}

func setLogMapping(d *schema.ResourceData, CSELogMapping *CSELogMapping) {
	d.Set("name", CSELogMapping.Name)
	if CSELogMapping.ParentId != "" {
		d.Set("parent_id", CSELogMapping.ParentId)
	}
	d.Set("product_guid", CSELogMapping.ProductGuid)
	d.Set("record_type", CSELogMapping.RecordType)
	enabled := *(CSELogMapping.Enabled)
	d.Set("enabled", enabled)
	d.Set("relates_entities", CSELogMapping.RelatesEntities)
	d.Set("skipped_values", CSELogMapping.SkippedValues)
	setFields(d, CSELogMapping.Fields)
	setStructuredInputs(d, CSELogMapping.StructuredInputs)
	if CSELogMapping.UnstructuredFields != nil {
		setUnstructuredFields(d, *(CSELogMapping.UnstructuredFields))
	}
}

func setUnstructuredFields(d *schema.ResourceData, unstructuredFields CSELogMappingUnstructuredFields) {
	resourceUnstructuredFieldsMap := make(map[string]interface{})
	resourceUnstructuredFieldsMap["patter_names"] = unstructuredFields.PatternNames

	resourceUnstructuredFields := make([]map[string]interface{}, 1)
	resourceUnstructuredFields[0] = resourceUnstructuredFieldsMap

	d.Set("unstructured_fields", resourceUnstructuredFields)

}

func setStructuredInputs(d *schema.ResourceData, structuredInputs []CSELogMappingStructuredInputField) {

	var f []map[string]interface{}

	for _, s := range structuredInputs {
		mapping := map[string]interface{}{
			"event_id_pattern": s.EventIdPattern,
			"log_format":       s.LogFormat,
			"product":          s.Product,
			"vendor":           s.Vendor,
		}
		f = append(f, mapping)
	}

	d.Set("structured_inputs", f)

}

func setFields(d *schema.ResourceData, fields []CSELogMappingField) {

	var f []map[string]interface{}

	for _, t := range fields {
		mapping := map[string]interface{}{
			"name":              t.Name,
			"value":             t.Value,
			"value_type":        t.Value,
			"skipped_values":    t.Value,
			"default_value":     t.Value,
			"format":            t.Value,
			"case_insensitive":  t.Value,
			"alternate_values":  t.Value,
			"time_zone":         t.Value,
			"split_delimiter":   t.Value,
			"split_index":       t.Value,
			"field_join":        t.Value,
			"join_delimiter":    t.Value,
			"format_parameters": t.Value,
			"lookup":            getLookUpResource(t.LookUp),
		}
		f = append(f, mapping)
	}

	d.Set("fields", f)

}

func getLookUpResource(lookUp []CSELogMappingLookUp) []map[string]interface{} {
	var s []map[string]interface{}

	for _, l := range lookUp {
		mapping := map[string]interface{}{
			"key":   l.Key,
			"value": l.Value,
		}
		s = append(s, mapping)
	}

	return s
}
