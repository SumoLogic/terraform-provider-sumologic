package sumologic

import (
	"log"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceSumologicSubdomain() *schema.Resource {
	return &schema.Resource{
		Create: resourceSumologicSubdomainCreate,
		Read:   resourceSumologicSubdomainRead,
		Update: resourceSumologicSubdomainUpdate,
		Delete: resourceSumologicSubdomainDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"subdomain": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: false,
			},
		},
	}
}

func resourceSumologicSubdomainRead(d *schema.ResourceData, meta interface{}) error {
	c := meta.(*Client)

	id := d.Id()
	subdomain, err := c.GetSubdomain()

	if err != nil {
		return err
	}

	if subdomain == nil {
		log.Printf("[WARN] Subdomain not found, removing from state: %v - %v", id, err)
		d.SetId("")
		return nil
	}

	d.Set("subdomain", subdomain.Subdomain)

	return nil
}

func resourceSumologicSubdomainCreate(d *schema.ResourceData, meta interface{}) error {
	c := meta.(*Client)

	if d.Id() == "" {
		subdomain := resourceToSubdomain(d)
		id, err := c.CreateSubdomain(subdomain)

		if err != nil {
			if strings.Contains(err.Error(), "subdomain:already_configured") {
				updatedID, updateErr := c.UpdateSubdomain(subdomain)
				if updateErr != nil {
					if updatedID == "" {
						id = subdomain.Subdomain
					} else {
						return updateErr
					}
				}
				id = updatedID
			} else {
				return err
			}
		}

		d.SetId(id)
	}

	return resourceSumologicSubdomainRead(d, meta)
}

func resourceSumologicSubdomainDelete(d *schema.ResourceData, meta interface{}) error {
	c := meta.(*Client)

	return c.DeleteSubdomain()
}

func resourceSumologicSubdomainUpdate(d *schema.ResourceData, meta interface{}) error {
	c := meta.(*Client)

	subdomain := resourceToSubdomain(d)

	_, err := c.UpdateSubdomain(subdomain)
	if err != nil {
		return err
	}

	return resourceSumologicSubdomainRead(d, meta)
}

func resourceToSubdomain(d *schema.ResourceData) Subdomain {
	return Subdomain{
		Subdomain: d.Get("subdomain").(string),
	}
}
