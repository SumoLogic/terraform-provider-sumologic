package sumologic

import (
	"fmt"
	"time"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/pborman/uuid"
)

func resourceSumologicHTTPSourceSNSAutoConfirm() *schema.Resource {
	return &schema.Resource{
		Create: resourceSumologicHTTPSourceSNSAutoConfirmCreate,
		Read:   resourceSumologicHTTPSourceSNSAutoConfirmRead,
		Delete: resourceSumologicHTTPSourceSNSAutoConfirmDelete,

		Schema: map[string]*schema.Schema{
			"category": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"confirmation_timeout": {
				Type:     schema.TypeInt,
				Optional: true,
				ForceNew: true,
				Default:  3,
			},
			"triggers": &schema.Schema{
				Type:     schema.TypeMap,
				Optional: true,
				ForceNew: true,
			},
		},
	}
}

func resourceSumologicHTTPSourceSNSAutoConfirmCreate(d *schema.ResourceData, meta interface{}) error {
	now := time.Now()
	confirmationTimeoutInMinutes := d.Get("confirmation_timeout").(int)

	err := resource.Retry(time.Duration(confirmationTimeoutInMinutes)*time.Minute, func() *resource.RetryError {
		c := meta.(*Client)

		if err := c.createHTTPSourceSNSAutoConfirm(d.Get("category").(string), now); err != nil {
			return resource.RetryableError(fmt.Errorf("Error fetching SubscriptionConfirmation message"))
		}

		return nil

	})

	if err != nil {
		return err
	}
	d.SetId(uuid.NewRandom().String())
	return nil
}

func resourceSumologicHTTPSourceSNSAutoConfirmRead(d *schema.ResourceData, meta interface{}) error {
	return nil
}

func resourceSumologicHTTPSourceSNSAutoConfirmDelete(d *schema.ResourceData, meta interface{}) error {
	d.SetId("")
	return nil
}
