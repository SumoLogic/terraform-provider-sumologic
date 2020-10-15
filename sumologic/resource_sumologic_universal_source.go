package sumologic

import (
	"encoding/json"
	"log"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/structure"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

func resourceSumologicUniversalSource() *schema.Resource {
	return &schema.Resource{
		Create: resourceSumologicUniversalSourceCreate,
		Read:   resourceSumologicUniversalSourceRead,
		Update: resourceSumologicUniversalSourceUpdate,
		Delete: resourceSumologicUniversalSourceDelete,
		Importer: &schema.ResourceImporter{
			State: resourceSumologicSourceImport,
		},
		Schema: map[string]*schema.Schema{
			"config": {
				Type:             schema.TypeString,
				ValidateFunc:     validation.StringIsJSON,
				Required:         true,
				DiffSuppressFunc: structure.SuppressJsonDiff,
			},
			"schema_ref": {
				Type:     schema.TypeMap,
				Required: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"collector_id": {
				Type:     schema.TypeInt,
				Required: true,
				ForceNew: true,
			},
		},
	}
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

func resourceSumologicUniversalSourceDelete(d *schema.ResourceData, meta interface{}) error {
	c := meta.(*Client)

	id, _ := strconv.Atoi(d.Id())
	collectorID, _ := d.Get("collector_id").(int)

	return c.DestroySource(id, collectorID)

}

func resourceToUniversalSource(d *schema.ResourceData) *UniversalSource {
	id, _ := strconv.Atoi(d.Id())
	var universalSource UniversalSource
	var jsonRawConf json.RawMessage

	universalSource.Type = "Universal"

	conf := []byte(d.Get("config").(string))

	err := json.Unmarshal(conf, &jsonRawConf)
	if err != nil {
		log.Println("Unable to unmarshal the Json configuration")
		return nil
	}

	universalSource.ID = id
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
