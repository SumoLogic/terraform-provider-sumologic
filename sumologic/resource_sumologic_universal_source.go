package sumologic

import (
	"encoding/json"
	"fmt"
	"log"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/structure"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

func resourceSumologicUniversalSource() *schema.Resource {
	universalSource := resourceSumologicSource()
	universalSource.Create = resourceSumologicUniversalSourceCreate
	universalSource.Read = resourceSumologicUniversalSourceRead
	universalSource.Update = resourceSumologicUniversalSourceUpdate
	universalSource.Importer = &schema.ResourceImporter{
		State: resourceSumologicSourceImport,
	}

	universalSource.Schema["config"] = &schema.Schema{
		Type:             schema.TypeString,
		ValidateFunc:     validation.StringIsJSON,
		Required:         true,
		DiffSuppressFunc: structure.SuppressJsonDiff,
	}

	universalSource.Schema["schema_ref"] = &schema.Schema{
		Type:     schema.TypeList,
		Required: true,
		MinItems: 1,
		MaxItems: 1,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"type": {
					Type:         schema.TypeString,
					Required:     true,
					ValidateFunc: validation.StringInSlice([]string{"Okta", "CrowdStrike", "Netskope", "Carbon Black Defense", "Azure Event Hubs"}, false),
				},
				"version": {
					Type:     schema.TypeString,
					Optional: true,
				},
			},
		},
	}

	return universalSource
}

func resourceSumologicUniversalSourceCreate(d *schema.ResourceData, meta interface{}) error {
	c := meta.(*Client)

	if d.Id() == "" {
		source := resourceToUniversalSource(d)
		// log.Printf("Source %s", source.Source)
		log.Printf("SchemaRef %s", source.SchemaRef)
		log.Printf("Config: %s", source.Config)

		id, err := c.CreateUniversalSource(*source, d.Get("collector_id").(int))

		if err != nil {
			return err
		}

		d.SetId(strconv.Itoa(id))
	}

	return resourceSumologicUniversalSourceRead(d, meta)
}

func resourceSumologicUniversalSourceUpdate(d *schema.ResourceData, meta interface{}) error {
	c := meta.(*Client)

	source := resourceToUniversalSource(d)

	err := c.UpdateUniversalSource(*source, d.Get("collector_id").(int))

	if err != nil {
		return err
	}

	return resourceSumologicUniversalSourceRead(d, meta)
}

func resourceToUniversalSource(d *schema.ResourceData) *UniversalSource {
	source := resourceToSource(d)
	source.Type = "Universal"

	var universalSource UniversalSource
	var jsonRawConf json.RawMessage

	conf := []byte(d.Get("config").(string))

	err := json.Unmarshal(conf, &jsonRawConf)
	if err != nil {
		panic(err)
	}

	// _ = json.Unmarshal([]byte(d.Get("config").(string)), &universalSource)

	universalSource.Source = source
	universalSource.Config = jsonRawConf
	universalSource.SchemaRef = getSourceSchemaRef(d)

	return &universalSource
}

func resourceSumologicUniversalSourceRead(d *schema.ResourceData, meta interface{}) error {
	c := meta.(*Client)

	id, _ := strconv.Atoi(d.Id())
	source, err := c.GetUniversalSource(d.Get("collector_id").(int), id)

	if err != nil {
		return err
	}

	if source == nil {
		log.Printf("[WARN] Universal source not found, removing from state: %v - %v", id, err)
		d.SetId("")

		return nil
	}

	if err := resourceSumologicSourceRead(d, source.Source); err != nil {
		return fmt.Errorf("%s", err)
	}
	d.Set("config", source.Config)
	d.Set("schema_ref", source.SchemaRef)

	return nil
}
func getSourceSchemaRef(d *schema.ResourceData) SchemaReference {
	sourceSchemas := d.Get("schema_ref").([]interface{})
	schemaR := SchemaReference{}

	if len(sourceSchemas) > 0 {
		sourceSchema := sourceSchemas[0].(map[string]interface{})
		schemaR.Type = sourceSchema["type"].(string)
		if len(sourceSchema["version"].(string)) > 0 {
			schemaR.Version = sourceSchema["version"].(string)
		}
	}

	return schemaR
}

// func mapToConfig(in map[string]interface{}) []Headers {
// 	headers := []Headers{}
// 	for k, v := range in {
// 		headers = append(headers, Headers{Name: k, Value: v.(string)})
// 	}

// 	return headers
// }

// func ConfigToMap(in []Headers) map[string]interface{} {
// 	headerMap := map[string]interface{}{}
// 	for _, header := range in {
// 		headerMap[header.Name] = header.Value
// 	}

// 	return headerMap
// }
