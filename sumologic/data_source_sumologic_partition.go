package sumologic

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func dataSourceSumologicPartition() *schema.Resource {
	resource := resourceSumologicPartition()
	return &schema.Resource{
		Read:   resource.Read,
		Schema: resource.Schema,
	}
}
