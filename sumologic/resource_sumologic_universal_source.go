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
		Type:     schema.TypeMap,
		Required: true,
		Elem: &schema.Schema{
			Type: schema.TypeString,
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
		log.Println("Unable to unmarshal the Json configuration")
		return nil
	}

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
	sourceSchema := d.Get("schema_ref").(map[string]interface{})
	schemaR := SchemaReference{}

	if len(sourceSchema) > 0 {
		schemaR.Type = sourceSchema["type"].(string)
		if sourceSchema["version"] != nil {
			schemaR.Version = sourceSchema["version"].(string)
		}
	}

	return schemaR
}
