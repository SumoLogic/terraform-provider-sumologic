package sumologic

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"log"
)

func resourceSumologicDataForwarding() *schema.Resource {
	return &schema.Resource{

		Create: resourceSumologicDataForwardingCreate,
		Read:   resourceSumologicDataForwardingRead,
		Update: resourceSumologicDataForwardingUpdate,
		Delete: resourceSumologicDataForwardingDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{

			"destination_name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"bucket_name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"authentication": {
				Type:     schema.TypeList,
				Required: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"type": {
							Type:     schema.TypeString,
							Required: true,
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

func resourceSumologicDataForwardingCreate(d *schema.ResourceData, meta interface{}) error {
	c := meta.(*Client)

	if d.Id() == "" {
		dataForwarding := resourceToDataForwarding(d)
		createdDataForwarding, err := c.CreateDataForwarding(dataForwarding)

		if err != nil {
			return err
		}

		d.SetId(createdDataForwarding.ID)

	}

	return resourceSumologicDataForwardingUpdate(d, meta)
}

func resourceSumologicDataForwardingUpdate(d *schema.ResourceData, meta interface{}) error {

	dataForwarding := resourceToDataForwarding(d)

	c := meta.(*Client)
	err := c.UpdateDataForwarding(dataForwarding)

	if err != nil {
		return err
	}

	return resourceSumologicDataForwardingRead(d, meta)
}

func resourceSumologicDataForwardingRead(d *schema.ResourceData, meta interface{}) error {
	c := meta.(*Client)

	id := d.Id()
	dataForwarding, err := c.getDataForwarding(id)

	if err != nil {
		return err
	}

	if dataForwarding == nil {
		log.Printf("[WARN] Data Forwarding not found, removing from state: %v - %v", id, err)
		d.SetId("")

		return nil
	}

	d.Set("destination_name", dataForwarding.DestinationName)
	d.Set("description", dataForwarding.Description)
	d.Set("bucket_name", dataForwarding.BucketName)
	d.Set("S3_region", dataForwarding.S3Region)
	d.Set("S3_server_side_encryption", dataForwarding.S3ServerSideEncryption)

	return nil
}

func resourceSumologicDataForwardingDelete(d *schema.ResourceData, meta interface{}) error {
	c := meta.(*Client)
	id := d.Id()

	return c.DeleteDataForwarding(id)
}

func resourceToDataForwarding(d *schema.ResourceData) DataForwarding {

	authentication := extractAuthenticationDetails(d.Get("authentication").([]interface{}))

	return DataForwarding{
		ID:                     d.Id(),
		DestinationName:        d.Get("destination_name").(string),
		Description:            d.Get("description").(string),
		BucketName:             d.Get("bucket_name").(string),
		AccessMethod:           authentication["type"].(string),
		AccessKey:              authentication["access_key"].(string),
		SecretKey:              authentication["secret_key"].(string),
		RoleArn:                authentication["role_arn"].(string),
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
