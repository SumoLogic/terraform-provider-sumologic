package sumologic

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

func resourceSumologicS3DataForwardingDestination() *schema.Resource {
	return &schema.Resource{
		Create: resourceSumologicS3DataForwardingDestinationCreate,
		Read:   resourceSumologicS3DataForwardingDestinationRead,
		Update: resourceSumologicS3DataForwardingDestinationUpdate,
		Delete: resourceSumologicS3DataForwardingDestinationDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"authentication_mode": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringInSlice([]string{"AccessKey", "RoleBased"}, false),
			},
			"access_key_id": {
				Type:      schema.TypeString,
				Sensitive: true,
				Optional:  true,
			},
			"secret_access_key": {
				Type:      schema.TypeString,
				Sensitive: true,
				Optional:  true,
			},
			"role_arn": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"region": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"encrypted": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"enabled": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"bucket_name": {
				Type:     schema.TypeString,
				Required: true,
			},
		},
	}
}

func resourceSumologicS3DataForwardingDestinationCreate(d *schema.ResourceData, m interface{}) error {
	c := m.(*Client)

	if d.Id() == "" {
		dfd := resourceToS3DataForwardingDestination(d)
		createdDfd, err := c.CreateS3DataForwardingDestination(dfd)

		if err != nil {
			return err
		}

		d.SetId(createdDfd.ID)
	}

	return resourceSumologicS3DataForwardingDestinationUpdate(d, m)
}

func resourceSumologicS3DataForwardingDestinationRead(d *schema.ResourceData, m interface{}) error {
	c := m.(*Client)
	dfd, err := c.GetS3DataForwardingDestination(d.Id())

	if err != nil {
		return err
	}

	d.Set("name", dfd.Name)
	d.Set("description", dfd.Description)
	d.Set("authentication_mode", dfd.AuthenticationMode)
	d.Set("access_key_id", dfd.AccessKeyID)
	d.Set("secret_access_key", dfd.SecretAccessKey)
	d.Set("role_arn", dfd.RoleARN)
	d.Set("region", dfd.Region)
	d.Set("encrypted", dfd.Encrypted)
	d.Set("enabled", dfd.Enabled)
	d.Set("bucket_name", dfd.BucketName)

	return nil
}

func resourceSumologicS3DataForwardingDestinationUpdate(d *schema.ResourceData, m interface{}) error {
	c := m.(*Client)
	dfd := resourceToS3DataForwardingDestination(d)
	err := c.UpdateS3DataForwardingDestination(dfd)

	if err != nil {
		return err
	}

	return resourceSumologicS3DataForwardingDestinationRead(d, m)
}

func resourceSumologicS3DataForwardingDestinationDelete(d *schema.ResourceData, m interface{}) error {
	c := m.(*Client)
	return c.DeleteS3DataForwardingDestination(d.Id())
}

func resourceToS3DataForwardingDestination(d *schema.ResourceData) S3DataForwardingDestination {
	return S3DataForwardingDestination{
		Name:               d.Get("name").(string),
		Description:        d.Get("description").(string),
		AuthenticationMode: d.Get("authentication_mode").(string),
		AccessKeyID:        d.Get("access_key_id").(string),
		SecretAccessKey:    d.Get("secret_access_key").(string),
		RoleARN:            d.Get("role_arn").(string),
		Region:             d.Get("region").(string),
		Encrypted:          d.Get("encrypted").(bool),
		Enabled:            d.Get("enabled").(bool),
		BucketName:         d.Get("bucket_name").(string),
	}
}
