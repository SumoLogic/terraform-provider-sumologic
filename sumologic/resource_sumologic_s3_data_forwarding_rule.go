package sumologic

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceSumologicS3DataForwardingRule() *schema.Resource {
	return &schema.Resource{
		Create: resourceSumologicS3DataForwardingRuleCreate,
		Read:   resourceSumologicS3DataForwardingRuleRead,
		Update: resourceSumologicS3DataForwardingRuleUpdate,
		Delete: resourceSumologicS3DataForwardingRuleDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"index_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"destination_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"enabled": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  true,
			},
			"file_format": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "{index}_{day}_{hour}_{minute}_{second}",
			},
			"format": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "csv",
			},
		},
	}
}

func resourceSumologicS3DataForwardingRuleCreate(d *schema.ResourceData, meta interface{}) error {
	c := meta.(*Client)

	if d.Id() == "" {
		dfr := resourceToS3DataForwardingRule(d)
		createdDfd, err := c.CreateS3DataForwardingRule(dfr)

		if err != nil {
			return err
		}

		d.SetId(createdDfd.IndexID)
	}

	return resourceSumologicS3DataForwardingRuleRead(d, meta)
}

func resourceSumologicS3DataForwardingRuleRead(d *schema.ResourceData, meta interface{}) error {
	c := meta.(*Client)
	id := d.Id()
	dfr, err := c.GetS3DataForwardingRule(id)

	if err != nil {
		return err
	}

	if dfr == nil {
		d.SetId("")
		return fmt.Errorf("S3DataForwardingRule for id %s not found", id)
	}

	d.SetId(dfr.IndexID)
	d.Set("index_id", dfr.IndexID)
	d.Set("destination_id", dfr.DestinationID)
	d.Set("enabled", dfr.Enabled)
	d.Set("file_format", dfr.FileFormat)
	d.Set("format", dfr.Format)

	return nil
}

func resourceSumologicS3DataForwardingRuleUpdate(d *schema.ResourceData, meta interface{}) error {
	c := meta.(*Client)
	dfr := resourceToS3DataForwardingRule(d)
	err := c.UpdateS3DataForwardingRule(dfr)

	if err != nil {
		return err
	}

	return resourceSumologicS3DataForwardingRuleRead(d, meta)
}

func resourceSumologicS3DataForwardingRuleDelete(d *schema.ResourceData, meta interface{}) error {
	c := meta.(*Client)
	return c.DeleteS3DataForwardingRule(d.Id())
}

func resourceToS3DataForwardingRule(d *schema.ResourceData) S3DataForwardingRule {
	return S3DataForwardingRule{
		IndexID:       d.Get("index_id").(string),
		DestinationID: d.Get("destination_id").(string),
		Enabled:       d.Get("enabled").(bool),
		FileFormat:    d.Get("file_format").(string),
		Format:        d.Get("format").(string),
	}
}
