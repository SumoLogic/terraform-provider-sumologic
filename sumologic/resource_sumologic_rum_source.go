package sumologic

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

func resourceSumologicRumSource() *schema.Resource {
	rumSource := resourceSumologicSource()
	rumSource.Create = resourceSumologicRumSourceCreate
	rumSource.Read = resourceSumologicRumSourceRead
	rumSource.Update = resourceSumologicRumSourceUpdate
	rumSource.Importer = &schema.ResourceImporter{
		State: resourceSumologicSourceImport,
	}

	rumSource.Schema["content_type"] = &schema.Schema{
		Type:         schema.TypeString,
		Required:     true,
		ValidateFunc: validation.StringInSlice([]string{"Rum"}, false),
	}

	rumSource.Schema["path"] = &schema.Schema{
		Type:     schema.TypeList,
		Optional: true,
		MaxItems: 1,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"application_name": {
					Type:     schema.TypeString,
					Optional: true,
				},
				"service_name": {
					Type:     schema.TypeString,
					Optional: false,
				},
				"deployment_environment": {
					Type:     schema.TypeString,
					Optional: true,
				},
				"sampling_rate": {
					Type:     schema.TypeFloat,
					Optional: true,
				},
				"ignore_urls": {
					Type: schema.TypeList,
					Elem: &schema.Schema{
						Type: schema.TypeString,
					},
					Optional: true,
				},
				"custom_tags": {
					Type: schema.TypeMap,
					Elem: &schema.Schema{
						Type: schema.TypeString,
					},
					Optional: true,
				},
				"propagate_trace_headers_cors_urls": {
					Type: schema.TypeList,
					Elem: &schema.Schema{
						Type: schema.TypeString,
					},
					Optional: true,
				},
				"selected_country": {
					Type:     schema.TypeString,
					Optional: true,
				},
			},
		},
	}

	return rumSource
}

func resourceSumologicRumSourceCreate(d *schema.ResourceData, meta interface{}) error {
	return nil
}

func resourceSumologicRumSourceUpdate(d *schema.ResourceData, meta interface{}) error {
	return nil
}

func resourceSumologicRumSourceRead(d *schema.ResourceData, meta interface{}) error {
	return nil
}
