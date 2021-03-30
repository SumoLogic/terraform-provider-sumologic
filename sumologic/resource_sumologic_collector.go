package sumologic

import (
	"fmt"
	"log"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceSumologicCollector() *schema.Resource {
	return &schema.Resource{
		Create: resourceSumologicCollectorCreate,
		Read:   resourceSumologicCollectorRead,
		Delete: resourceSumologicCollectorDelete,
		Update: resourceSumologicCollectorUpdate,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: false,
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: false,
			},
			"category": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: false,
			},
			"timezone": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: false,
				Default:  "Etc/UTC",
			},
			"fields": {
				Type: schema.TypeMap,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Optional: true,
				ForceNew: false,
			},
		},
	}
}

func resourceSumologicCollectorRead(d *schema.ResourceData, meta interface{}) error {
	c := meta.(*Client)

	var collector *Collector
	id, err := strconv.Atoi(d.Id())
	if err != nil {
		collector, err = c.GetCollectorName(d.Id())
		if err != nil {
			log.Printf("[WARN] Collector not found when looking by name: %s, err: %v", d.Id(), err)
		} else if collector == nil {
			log.Printf("[WARN] Got a nil Collector when looking by name: %s", d.Id())
		} else {
			d.SetId(strconv.FormatInt(collector.ID, 10))
		}
	} else {
		collector, err = c.GetCollector(id)
		if err != nil {
			log.Printf("[WARN] Collector not found when looking by id: %d, err: %v", id, err)
		}
	}

	if collector == nil {
		log.Printf("[WARN] Collector not found, removing from state: %v - %v", id, err)
		d.SetId("")
		return nil
	}

	d.Set("name", collector.Name)
	d.Set("description", collector.Description)
	d.Set("category", collector.Category)
	d.Set("timezone", collector.TimeZone)
	if err := d.Set("fields", collector.Fields); err != nil {
		return fmt.Errorf("error setting fields for resource %s: %s", d.Id(), err)
	}

	return nil
}

func resourceSumologicCollectorDelete(d *schema.ResourceData, meta interface{}) error {
	c := meta.(*Client)

	id, _ := strconv.Atoi(d.Id())
	return c.DeleteCollector(id)

}

func resourceSumologicCollectorCreate(d *schema.ResourceData, meta interface{}) error {
	c := meta.(*Client)

	if d.Id() == "" {
		id, err := c.CreateCollector(Collector{
			CollectorType: "Hosted",
			Name:          d.Get("name").(string),
		})

		if err != nil {
			return err
		}

		d.SetId(strconv.FormatInt(id, 10))
	}

	return resourceSumologicCollectorUpdate(d, meta)
}

func resourceSumologicCollectorUpdate(d *schema.ResourceData, meta interface{}) error {
	collector, err := resourceToCollector(d)
	if err != nil {
		return err
	}

	c := meta.(*Client)
	if err = c.UpdateCollector(collector); err != nil {
		return err
	}

	return resourceSumologicCollectorRead(d, meta)
}

func resourceToCollector(d *schema.ResourceData) (Collector, error) {
	id, err := strconv.ParseInt(d.Id(), 10, 64)
	if err != nil {
		return Collector{}, err
	}

	return Collector{
		ID:            id,
		CollectorType: "Hosted",
		Name:          d.Get("name").(string),
		Description:   d.Get("description").(string),
		Category:      d.Get("category").(string),
		TimeZone:      d.Get("timezone").(string),
		Fields:        d.Get("fields").(map[string]interface{}),
	}, nil
}
