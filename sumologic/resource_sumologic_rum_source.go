package sumologic

import (
	"errors"
	"fmt"
	"log"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
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
		Optional:     true,
		Default:      "Rum",
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
					Required: true,
				},
				"deployment_environment": {
					Type:     schema.TypeString,
					Optional: true,
				},
				"sampling_rate": {
					Type:         schema.TypeFloat,
					Optional:     true,
					ValidateFunc: validation.FloatBetween(0, 1),
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
				"propagate_trace_header_cors_urls": {
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
	c := meta.(*Client)

	if d.Id() == "" {
		source, err := resourceToRumSource(d)
		if err != nil {
			return err
		}

		id, err := c.CreateRumSource(source, d.Get("collector_id").(int))

		if err != nil {
			return err
		}

		d.SetId(strconv.Itoa(id))
	}

	return resourceSumologicRumSourceRead(d, meta)
}

func resourceSumologicRumSourceUpdate(d *schema.ResourceData, meta interface{}) error {
	c := meta.(*Client)

	source, err := resourceToRumSource(d)
	if err != nil {
		return err
	}

	err = c.UpdateRumSource(source, d.Get("collector_id").(int))

	if err != nil {
		return err
	}

	return resourceSumologicRumSourceRead(d, meta)
}

func resourceSumologicRumSourceRead(d *schema.ResourceData, meta interface{}) error {
	c := meta.(*Client)

	id, _ := strconv.Atoi(d.Id())
	source, err := c.GetRumSource(d.Get("collector_id").(int), id)

	if err != nil {
		return err
	}

	if source == nil {
		log.Printf("[WARN] Rum source not found, removing from state: %v - %v", id, err)
		d.SetId("")

		return nil
	}

	if err := resourceSumologicSourceRead(d, source.Source); err != nil {
		return fmt.Errorf("%s", err)
	}
	d.Set("content_type", source.ContentType)

	return nil
}

func resourceToRumSource(d *schema.ResourceData) (RumSource, error) {
	source := resourceToSource(d)
	source.Type = "HTTP"

	rumSource := RumSource{
		Source: source,
	}

	path, errPath := getRumSourcePath(d)
	if errPath != nil {
		return rumSource, errPath
	}

	auth := RumAuthentication{
		Type: "NoAuthentication",
	}

	rumThirdPartyResource := RumThirdPartyResource{
		ServiceType:    "Rum",
		Path:           path,
		Authentication: auth,
	}

	rumSource.RumThirdPartyRef.Resources = append(rumSource.RumThirdPartyRef.Resources, rumThirdPartyResource)

	return rumSource, nil
}

func getRumSourcePath(d *schema.ResourceData) (RumSourcePath, error) {
	rumSourcePath := RumSourcePath{}
	paths := d.Get("path").([]interface{})

	if len(paths) > 0 {
		path := paths[0].(map[string]interface{})
		rumSourcePath.Type = "RumPath"
		rumSourcePath.ApplicationName = path["application_name"].(string)
		rumSourcePath.ServiceName = path["service_name"].(string)
		rumSourcePath.DeploymentEnvironment = path["deployment_environment"].(string)
		rumSourcePath.SamplingRate = path["sampling_rate"].(float64)

		ignoreUrls_int := path["ignore_urls"].([]interface{})
		IgnoreUrls := make([]string, len(ignoreUrls_int))
		for i, v := range ignoreUrls_int {
			IgnoreUrls[i] = v.(string)
		}
		rumSourcePath.IgnoreUrls = IgnoreUrls

		rumSourcePath.CustomTags = path["custom_tags"].(map[string]interface{})

		propagateTraceHeaderCorsUrls_int := path["ignore_urls"].([]interface{})
		PropagateTraceHeaderCorsUrls := make([]string, len(propagateTraceHeaderCorsUrls_int))
		for i, v := range propagateTraceHeaderCorsUrls_int {
			PropagateTraceHeaderCorsUrls[i] = v.(string)
		}
		rumSourcePath.PropagateTraceHeaderCorsUrls = PropagateTraceHeaderCorsUrls

		rumSourcePath.SelectedCountry = path["selected_country"].(string)

		return rumSourcePath, nil
	} else {
		return rumSourcePath, errors.New("[ERROR] Rum path not configured")
	}
}
