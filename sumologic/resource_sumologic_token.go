package sumologic

import (
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

func resourceSumologicToken() *schema.Resource {
	return &schema.Resource{
		Create: resourceSumologicTokenCreate,
		Read:   resourceSumologicTokenRead,
		Update: resourceSumologicTokenUpdate,
		Delete: resourceSumologicTokenDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"version": {
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
			},

			"status": {
				Type:     schema.TypeString,
				Required: true,
			},

			"name": {
				Type:     schema.TypeString,
				Required: true,
			},

			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},

			"type": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringInSlice([]string{"CollectorRegistration", "CollectorRegistrationTokenResponse"}, false),
			},
		},
	}
}

func resourceSumologicTokenCreate(d *schema.ResourceData, meta interface{}) error {
	c := meta.(*Client)

	if d.Id() == "" {
		token := resourceToToken(d)
		id, err := c.CreateToken(token)
		if err != nil {
			return err
		}

		d.SetId(id)
	}

	return resourceSumologicTokenRead(d, meta)
}

func resourceSumologicTokenRead(d *schema.ResourceData, meta interface{}) error {
	c := meta.(*Client)

	id := d.Id()
	token, err := c.GetToken(id)
	if err != nil {
		return err
	}

	if token == nil {
		log.Printf("[WARN] Token not found, removing from state: %v - %v", id, err)
		d.SetId("")
		return nil
	}

	d.Set("name", token.Name)
	d.Set("status", token.Status)
	d.Set("description", token.Description)
	d.Set("version", token.Version)

	return nil
}

func resourceSumologicTokenDelete(d *schema.ResourceData, meta interface{}) error {
	c := meta.(*Client)

	return c.DeleteToken(d.Id())
}

func resourceSumologicTokenUpdate(d *schema.ResourceData, meta interface{}) error {
	c := meta.(*Client)

	token := resourceToToken(d)
	err := c.UpdateToken(token)
	if err != nil {
		return err
	}

	return resourceSumologicTokenRead(d, meta)
}

func resourceToToken(d *schema.ResourceData) Token {

	return Token{
		Name:        d.Get("name").(string),
		ID:          d.Id(),
		Description: d.Get("description").(string),
		Version:     d.Get("version").(int),
		Type:        d.Get("type").(string),
		Status:      d.Get("status").(string),
	}
}
