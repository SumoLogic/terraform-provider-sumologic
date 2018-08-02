package sumologic

import (
	"log"
	"strconv"

	"github.com/hashicorp/terraform/helper/schema"
)

func resourceSumologicCloudsyslogSource() *schema.Resource {
	return &schema.Resource{
		Create: resourceSumologicCloudsyslogSourceCreate,
		Read:   resourceSumologicCloudsyslogSourceRead,
		Update: resourceSumologicCloudsyslogSourceUpdate,
		Delete: resourceSumologicCloudsyslogSourceDelete,

		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: false,
				Default:  "",
			},
			"collector_id": {
				Type:     schema.TypeInt,
				Required: true,
				ForceNew: true,
			},
			"category": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: false,
				Default:  "",
			},
			"token": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"lookup_by_name": {
				Type:     schema.TypeBool,
				Optional: true,
				ForceNew: false,
				Default:  false,
			},
			"destroy": {
				Type:     schema.TypeBool,
				Optional: true,
				ForceNew: false,
				Default:  true,
			},
		},
	}
}

func resourceSumologicCloudsyslogSourceCreate(d *schema.ResourceData, meta interface{}) error {
	c := meta.(*Client)

	if d.Get("lookup_by_name").(bool) {
		source, err := c.GetSourceName(d.Get("collector_id").(int), d.Get("name").(string))

		if err != nil {
			return err
		}

		// Set ID of source if it exists
		if source != nil {
			d.SetId(strconv.Itoa(source.ID))
		}
	}

	// If source ID is still empty, create a new source
	if d.Id() == "" {
		source := CloudsyslogSource{}
		source.Name = d.Get("name").(string)

		id, err := c.CreateCloudsyslogSource(
			d.Get("name").(string),
			d.Get("description").(string),
			d.Get("collector_id").(int),
		)

		if err != nil {
			return err
		}

		d.SetId(strconv.Itoa(id))
	}

	return resourceSumologicCloudsyslogSourceUpdate(d, meta)
}

func resourceSumologicCloudsyslogSourceUpdate(d *schema.ResourceData, meta interface{}) error {
	c := meta.(*Client)

	source := resourceToCloudsyslogSource(d)

	err := c.UpdateCloudsyslogSource(source, d.Get("collector_id").(int))

	if err != nil {
		return err
	}

	return resourceSumologicCloudsyslogSourceRead(d, meta)
}

func resourceToCloudsyslogSource(d *schema.ResourceData) CloudsyslogSource {

	id, _ := strconv.Atoi(d.Id())

	source := CloudsyslogSource{}
	source.ID = id
	source.Type = "Cloudsyslog"
	source.Name = d.Get("name").(string)
	source.Description = d.Get("description").(string)
	source.Category = d.Get("category").(string)

	return source
}

func resourceSumologicCloudsyslogSourceRead(d *schema.ResourceData, meta interface{}) error {
	c := meta.(*Client)

	id, _ := strconv.Atoi(d.Id())
	source, err := c.GetCloudsyslogSource(d.Get("collector_id").(int), id)

	// Source is gone, remove it from state
	if err != nil {
		log.Printf("Cloudsyslog source %v: %v", id, err)
		d.SetId("")

		return nil
	}

	d.Set("name", source.Name)
	d.Set("description", source.Description)
	d.Set("token", source.Token)

	return nil
}

func resourceSumologicCloudsyslogSourceDelete(d *schema.ResourceData, meta interface{}) error {
	c := meta.(*Client)

	// Destroy source if `destroy` is true, otherwise ignore
	if d.Get("destroy").(bool) {
		id, _ := strconv.Atoi(d.Id())
		collectorID, _ := d.Get("collector_id").(int)

		return c.DestroySource(id, collectorID)
	}

	return nil
}
