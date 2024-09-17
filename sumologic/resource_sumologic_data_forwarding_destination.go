package sumologic

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"log"
)

func resourceSumologicDataForwardingDestination() *schema.Resource {
	return &schema.Resource{
		Create: resourceSumologicDataForwardingDestinationCreate,
		Read:   resourceSumologicDataForwardingDestinationRead,
		Update: resourceSumologicDataForwardingDestinationUpdate,
		Delete: resourceSumologicDataForwardingDestinationDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{

			"access_key_id": {
				Type:     schema.TypeString,
				Optional: true,
			},

			"encrypted": {
				Type:     schema.TypeBool,
				Optional: true,
			},

			"region": {
				Type:     schema.TypeString,
				Optional: true,
			},

			"destination_name": {
				Type:     schema.TypeString,
				Optional: true,
			},

			"enabled": {
				Type:     schema.TypeBool,
				Optional: true,
			},

			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},

			"authentication_mode": {
				Type:     schema.TypeString,
				Required: true,
			},

			"role_arn": {
				Type:     schema.TypeString,
				Optional: true,
			},

			"bucket_name": {
				Type:     schema.TypeString,
				Required: true,
			},

			"secret_access_key": {
				Type:     schema.TypeString,
				Optional: true,
			},
		},
	}
}

func resourceSumologicDataForwardingDestinationCreate(d *schema.ResourceData, meta interface{}) error {
	c := meta.(*Client)

	if d.Id() == "" {
		dataForwardingDestination := resourceToDataForwardingDestination(d)
		id, err := c.CreateDataForwardingDestination(dataForwardingDestination)
		if err != nil {
			return err
		}

		d.SetId(id)
	}

	return resourceSumologicDataForwardingDestinationRead(d, meta)
}

func resourceSumologicDataForwardingDestinationRead(d *schema.ResourceData, meta interface{}) error {
	c := meta.(*Client)

	id := d.Id()
	dataForwardingDestination, err := c.GetDataForwardingDestination(id)
	if err != nil {
		return err
	}

	if dataForwardingDestination == nil {
		log.Printf("[WARN] DataForwardingDestination not found, removing from state: %v - %v", id, err)
		d.SetId("")
		return nil
	}

	d.Set("destination_name", dataForwardingDestination.DestinationName)
	d.Set("role_arn", dataForwardingDestination.RoleArn)
	d.Set("authentication_mode", dataForwardingDestination.AuthenticationMode)
	d.Set("region", dataForwardingDestination.Region)
	d.Set("secret_access_key", dataForwardingDestination.SecretAccessKey)
	d.Set("bucket_name", dataForwardingDestination.BucketName)
	d.Set("enabled", dataForwardingDestination.Enabled)
	d.Set("description", dataForwardingDestination.Description)
	d.Set("access_key_id", dataForwardingDestination.AccessKeyId)
	d.Set("encrypted", dataForwardingDestination.Encrypted)

	return nil
}

func resourceSumologicDataForwardingDestinationDelete(d *schema.ResourceData, meta interface{}) error {
	c := meta.(*Client)

	return c.DeleteDataForwardingDestination(d.Id())
}

func resourceSumologicDataForwardingDestinationUpdate(d *schema.ResourceData, meta interface{}) error {
	c := meta.(*Client)

	dataForwardingDestination := resourceToDataForwardingDestination(d)
	err := c.UpdateDataForwardingDestination(dataForwardingDestination)
	if err != nil {
		return err
	}

	return resourceSumologicDataForwardingDestinationRead(d, meta)
}

func resourceToDataForwardingDestination(d *schema.ResourceData) DataForwardingDestination {

	return DataForwardingDestination{
		BucketName:         d.Get("bucket_name").(string),
		DestinationName:    d.Get("destination_name").(string),
		ID:                 d.Id(),
		RoleArn:            d.Get("role_arn").(string),
		Region:             d.Get("region").(string),
		Description:        d.Get("description").(string),
		Enabled:            d.Get("enabled").(bool),
		AccessKeyId:        d.Get("access_key_id").(string),
		Encrypted:          d.Get("encrypted").(bool),
		AuthenticationMode: d.Get("authentication_mode").(string),
		SecretAccessKey:    d.Get("secret_access_key").(string),
	}
}
