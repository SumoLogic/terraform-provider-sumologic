package sumologic

import (
	"encoding/json"
	"log"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/structure"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

func resourceSumologicCloudToCloudSource() *schema.Resource {
	return &schema.Resource{
		Create: resourceSumologicCloudToCloudSourceCreate,
		Read:   resourceSumologicCloudToCloudSourceRead,
		Update: resourceSumologicCloudToCloudSourceUpdate,
		Delete: resourceSumologicCloudToCloudSourceDelete,
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

func resourceSumologicCloudToCloudSourceCreate(d *schema.ResourceData, meta interface{}) error {
	c := meta.(*Client)

	if d.Id() == "" {
		source := resourceToCloudToCloudSource(d)
		log.Printf("SchemaRef %s", source.SchemaRef)
		log.Printf("Config: %s", source.Config)

		id, err := c.CreateCloudToCloudSource(*source, d.Get("collector_id").(int))

		if err != nil {
			return err
		}

		d.SetId(strconv.Itoa(id))
	}

	return resourceSumologicCloudToCloudSourceRead(d, meta)
}

func resourceSumologicCloudToCloudSourceUpdate(d *schema.ResourceData, meta interface{}) error {
	c := meta.(*Client)

	source := resourceToCloudToCloudSource(d)

	err := c.UpdateCloudToCloudSource(*source, d.Get("collector_id").(int))

	if err != nil {
		return err
	}

	return resourceSumologicCloudToCloudSourceRead(d, meta)
}

func resourceSumologicCloudToCloudSourceDelete(d *schema.ResourceData, meta interface{}) error {
	c := meta.(*Client)

	id, _ := strconv.Atoi(d.Id())
	collectorID, _ := d.Get("collector_id").(int)

	return c.DestroySource(id, collectorID)

}

func resourceToCloudToCloudSource(d *schema.ResourceData) *CloudToCloudSource {
	id, _ := strconv.Atoi(d.Id())
	var cloudToCloudSource CloudToCloudSource
	var jsonRawConf json.RawMessage

	cloudToCloudSource.Type = "Universal"

	conf := []byte(d.Get("config").(string))

	err := json.Unmarshal(conf, &jsonRawConf)
	if err != nil {
		log.Println("Unable to unmarshal the Json configuration")
		return nil
	}

	cloudToCloudSource.ID = id
	cloudToCloudSource.Config = jsonRawConf
	cloudToCloudSource.SchemaRef = getSourceSchemaRef(d)

	return &cloudToCloudSource
}

func resourceSumologicCloudToCloudSourceRead(d *schema.ResourceData, meta interface{}) error {
	c := meta.(*Client)

	id, _ := strconv.Atoi(d.Id())
	source, err := c.GetCloudToCloudSource(d.Get("collector_id").(int), id)

	if err != nil {
		return err
	}

	if source == nil {
		log.Printf("[WARN] Cloud-to-Cloud source not found, removing from state: %v - %v", id, err)
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
	}

	return schemaR
}
