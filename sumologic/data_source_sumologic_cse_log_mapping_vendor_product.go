package sumologic

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func dataSourceCSELogMappingVendorAndProduct() *schema.Resource {

	return &schema.Resource{
		Read: dataSourceCSELogMappingVendorAndProductRead,
		Schema: map[string]*schema.Schema{
			"guid": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"product": {
				Type:     schema.TypeString,
				Required: true,
			},
			"vendor": {
				Type:     schema.TypeString,
				Required: true,
			},
		},
	}

}

func dataSourceCSELogMappingVendorAndProductRead(d *schema.ResourceData, meta interface{}) error {
	c := meta.(*Client)

	id := d.Id()

	vendorAndProduct, err := c.GetCSELogMappingVendorsAndProducts(d.Get("product").(string), d.Get("vendor").(string))

	if err != nil {
		return err
	}

	if vendorAndProduct == nil {
		d.SetId("")
		return fmt.Errorf("Vendor product not found, removing from state: %v - %v", id, err)
	}

	d.SetId(vendorAndProduct.GUID)
	d.Set("guid", vendorAndProduct.GUID)
	d.Set("product", vendorAndProduct.Product)
	d.Set("vendor", vendorAndProduct.Vendor)

	return nil
}
