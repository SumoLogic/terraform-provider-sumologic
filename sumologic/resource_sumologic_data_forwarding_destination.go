package sumologic

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
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

			"destination_name": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringLenBetween(1, 255),
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"bucket_name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"authentication": {
				Type:     schema.TypeList,
				Required: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"type": {
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: validation.StringInSlice([]string{"RoleBased", "AccessKey"}, false),
						},
						"role_arn": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"access_key": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"secret_key": {
							Type:     schema.TypeString,
							Optional: true,
						},
					},
				},
			},
			"enabled": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"s3_region": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"s3_server_side_encryption": {
				Type:     schema.TypeBool,
				Optional: true,
			},
		},
	}

}

func resourceSumologicDataForwardingDestinationCreate(d *schema.ResourceData, meta interface{}) error {
	c := meta.(*Client)

	if d.Id() == "" {
		dataForwardingDestination := resourceToDataForwardingDestination(d)
		createdDataForwardingDestination, err := c.CreateDataForwardingDestination(dataForwardingDestination)

		if err != nil {
			return err
		}

		d.SetId(createdDataForwardingDestination.ID)

	}

	return resourceSumologicDataForwardingDestinationUpdate(d, meta)
}

func resourceSumologicDataForwardingDestinationUpdate(d *schema.ResourceData, meta interface{}) error {

	dataForwardingDestination := resourceToDataForwardingDestination(d)

	c := meta.(*Client)
	err := c.UpdateDataForwardingDestination(dataForwardingDestination)

	if err != nil {
		return err
	}

	return resourceSumologicDataForwardingDestinationRead(d, meta)
}

func resourceSumologicDataForwardingDestinationRead(d *schema.ResourceData, meta interface{}) error {
	c := meta.(*Client)

	id := d.Id()
	dataForwardingDestination, err := c.getDataForwardingDestination(id)

	if err != nil {
		return err
	}

	if dataForwardingDestination == nil {
		log.Printf("[WARN] Data Forwarding destination not found, removing from state: %v - %v", id, err)
		d.SetId("")

		return nil
	}

	d.Set("destination_name", dataForwardingDestination.DestinationName)
	d.Set("description", dataForwardingDestination.Description)
	d.Set("bucket_name", dataForwardingDestination.BucketName)
	d.Set("S3_region", dataForwardingDestination.S3Region)
	d.Set("enabled", dataForwardingDestination.Enabled)
	d.Set("S3_server_side_encryption", dataForwardingDestination.S3ServerSideEncryption)

	return nil
}

func resourceSumologicDataForwardingDestinationDelete(d *schema.ResourceData, meta interface{}) error {
	c := meta.(*Client)
	id := d.Id()

	return c.DeleteDataForwardingDestination(id)
}

func resourceToDataForwardingDestination(d *schema.ResourceData) DataForwardingDestination {

	authentication := extractAuthenticationDetails(d.Get("authentication").([]interface{}))

	return DataForwardingDestination{
		ID:                     d.Id(),
		DestinationName:        d.Get("destination_name").(string),
		Description:            d.Get("description").(string),
		BucketName:             d.Get("bucket_name").(string),
		AccessMethod:           authentication["type"].(string),
		AccessKey:              authentication["access_key"].(string),
		SecretKey:              authentication["secret_key"].(string),
		RoleArn:                authentication["role_arn"].(string),
		Enabled:                d.Get("enabled").(bool),
		S3Region:               d.Get("s3_region").(string),
		S3ServerSideEncryption: d.Get("s3_server_side_encryption").(bool),
	}
}

func extractAuthenticationDetails(authenticationList []interface{}) map[string]interface{} {
	if len(authenticationList) > 0 {
		authMap := authenticationList[0].(map[string]interface{})
		return authMap
	}
	return map[string]interface{}{
		"type":       nil,
		"role_arn":   nil,
		"access_key": nil,
		"secret_key": nil,
	}
}
