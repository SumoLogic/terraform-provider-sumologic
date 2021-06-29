package sumologic

import (
	"encoding/json"
	"log"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/structure"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

func resourceSumologicContent() *schema.Resource {
	return &schema.Resource{
		Create: resourceSumologicContentCreate,
		Read:   resourceSumologicContentRead,
		Update: resourceSumologicContentUpdate,
		Delete: resourceSumologicContentDelete,

		Schema: map[string]*schema.Schema{
			"parent_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"config": {
				Type:             schema.TypeString,
				ValidateFunc:     validation.StringIsJSON,
				Required:         true,
				DiffSuppressFunc: structure.SuppressJsonDiff,
			},
		},
		Timeouts: &schema.ResourceTimeout{
			Read:   schema.DefaultTimeout(1 * time.Minute),
			Create: schema.DefaultTimeout(10 * time.Minute),
			Update: schema.DefaultTimeout(10 * time.Minute),
			Delete: schema.DefaultTimeout(1 * time.Minute),
		},
	}
}

func resourceSumologicContentRead(d *schema.ResourceData, meta interface{}) error {
	c := meta.(*Client)
	// Retrieve the content Id from the state
	id := d.Id()
	log.Printf("[DEBUG] Looking for content with id: %s", id)

	content, err := c.GetContent(id, d.Timeout(schema.TimeoutRead))
	if err != nil {
		return err
	}
	if content == nil {
		log.Printf("[WARN] Content not found, removing from state: %v - %v", id, err)
		d.SetId("")
		return nil
	}

	log.Printf("[DEBUG] content: %s", content.Name)
	log.Printf("[DEBUG] parent of content: %s", content.ParentId)
	log.Printf("[DEBUG] content config: %s", content.Config)

	// Write the newly read content object into the schema
	d.Set("config", content.Config)
	return nil
}

func resourceSumologicContentDelete(d *schema.ResourceData, meta interface{}) error {
	c := meta.(*Client)
	log.Printf("Deleting content with id: %s", d.Id())
	return c.DeleteContent(d.Id(), d.Timeout(schema.TimeoutDelete))
}

func resourceSumologicContentCreate(d *schema.ResourceData, meta interface{}) error {
	c := meta.(*Client)

	// If there is no id in the state, then we need to create the object
	if d.Id() == "" {
		// Load all the data we have from the schema into a Content Struct
		content := resourceToContent(d)

		id, err := c.CreateOrUpdateContent(*content, d.Timeout(schema.TimeoutCreate), false)
		if err != nil {
			return err
		}

		d.SetId(id)
		log.Printf("Created content with id=%s, type=%s", id, content.Type)
	}

	return resourceSumologicContentRead(d, meta)
}

func resourceSumologicContentUpdate(d *schema.ResourceData, meta interface{}) error {
	c := meta.(*Client)

	content := resourceToContent(d)

	id, err := c.CreateOrUpdateContent(*content, d.Timeout(schema.TimeoutUpdate), true)
	if err != nil {
		return err
	}

	d.SetId(id)
	log.Printf("Updated content with id=%s, type=%s", id, content.Type)

	return resourceSumologicContentRead(d, meta)
}

func resourceToContent(d *schema.ResourceData) *Content {
	var content Content

	_ = json.Unmarshal([]byte(d.Get("config").(string)), &content)

	content.Children = []Content{}
	content.ParentId = d.Get("parent_id").(string)
	content.Config = d.Get("config").(string)
	content.ID = d.Id()

	return &content
}
