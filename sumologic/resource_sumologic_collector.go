package sumologic

import (
	"log"
	"strconv"

	"github.com/hashicorp/terraform/helper/schema"
)

func resourceSumologicCollector() *schema.Resource {
	return &schema.Resource{
		Create: resourceSumologicCollectorCreate,
		Read:   resourceSumologicCollectorRead,
		Delete: resourceSumologicCollectorDelete,
		Update: resourceSumologicCollectorUpdate,

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
			"category": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: false,
				Default:  "",
			},
			"timezone": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: false,
				Default:  "Etc/UTC",
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

func resourceSumologicCollectorRead(d *schema.ResourceData, meta interface{}) error {
	c := meta.(*Client)

	id, err := strconv.Atoi(d.Id())

	if err != nil {
		return err
	}

	collector, err := c.GetCollector(id)

	if err != nil {
		log.Printf("[WARN] Collector not found, removing from state: %v - %v", id, err)
		d.SetId("")

		return nil
	}

	d.Set("name", collector.Name)
	d.Set("description", collector.Description)
	d.Set("category", collector.Category)
	d.Set("timezone", collector.TimeZone)

	return nil
}

func resourceSumologicCollectorDelete(d *schema.ResourceData, meta interface{}) error {
	c := meta.(*Client)

	if d.Get("destroy").(bool) {
		id, _ := strconv.Atoi(d.Id())
		return c.DeleteCollector(id)
	}

	return nil
}

func resourceSumologicCollectorCreate(d *schema.ResourceData, meta interface{}) error {
	c := meta.(*Client)

	if d.Get("lookup_by_name").(bool) {
		collector, err := c.GetCollectorName(d.Get("name").(string))

		if err != nil {
			return err
		}

		if collector != nil {
			d.SetId(strconv.Itoa(collector.ID))
		}
	}

	if d.Id() == "" {
		id, err := c.CreateCollector(Collector{
			CollectorType: "Hosted",
			Name:          d.Get("name").(string),
		})

		if err != nil {
			return err
		}

		d.SetId(strconv.Itoa(id))
	}

	return resourceSumologicCollectorUpdate(d, meta)
}

func resourceSumologicCollectorUpdate(d *schema.ResourceData, meta interface{}) error {

	collector := resourceToCollector(d)

	c := meta.(*Client)
	err := c.UpdateCollector(collector)

	if err != nil {
		return err
	}

	return resourceSumologicCollectorRead(d, meta)
}

func resourceToCollector(d *schema.ResourceData) Collector {
	id, _ := strconv.Atoi(d.Id())

	return Collector{
		ID:            id,
		CollectorType: "Hosted",
		Name:          d.Get("name").(string),
		Description:   d.Get("description").(string),
		Category:      d.Get("category").(string),
		TimeZone:      d.Get("timezone").(string),
	}
}
