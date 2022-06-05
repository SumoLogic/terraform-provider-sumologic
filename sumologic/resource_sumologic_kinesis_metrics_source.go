package sumologic

import (
	"errors"
	"fmt"
	"log"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

func resourceSumologicKinesisMetricsSource() *schema.Resource {
	kinesisMetricsSource := resourceSumologicSource()
	kinesisMetricsSource.Create = resourceSumologicKinesisMetricsSourceCreate
	kinesisMetricsSource.Read = resourceSumologicKinesisMetricsSourceRead
	kinesisMetricsSource.Update = resourceSumologicKinesisMetricsSourceUpdate
	kinesisMetricsSource.Importer = &schema.ResourceImporter{
		State: resourceSumologicSourceImport,
	}

	kinesisMetricsSource.Schema["content_type"] = &schema.Schema{
		Type:         schema.TypeString,
		Required:     true,
		ValidateFunc: validation.StringInSlice([]string{"KinesisMetric"}, false),
	}

	kinesisMetricsSource.Schema["message_per_request"] = &schema.Schema{
		Type:     schema.TypeBool,
		Optional: true,
		Default:  false,
	}
	kinesisMetricsSource.Schema["url"] = &schema.Schema{
		Type:     schema.TypeString,
		Computed: true,
	}
	kinesisMetricsSource.Schema["authentication"] = &schema.Schema{
		Type:     schema.TypeList,
		Required: true,
		MinItems: 1,
		MaxItems: 1,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"type": {
					Type:         schema.TypeString,
					Required:     true,
					ValidateFunc: validation.StringInSlice([]string{"S3BucketAuthentication", "AWSRoleBasedAuthentication"}, false),
				},
				"access_key": {
					Type:     schema.TypeString,
					Optional: true,
				},
				"secret_key": {
					Type:     schema.TypeString,
					Optional: true,
				},
				"role_arn": {
					Type:     schema.TypeString,
					Optional: true,
				},
			},
		},
	}

	kinesisMetricsSource.Schema["path"] = &schema.Schema{
		Type:     schema.TypeList,
		Required: true,
		MinItems: 1,
		MaxItems: 1,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"type": {
					Type:         schema.TypeString,
					Required:     true,
					ValidateFunc: validation.StringInSlice([]string{"KinesisMetricPath"}, false),
				},

				"tag_filters": {
					Type:     schema.TypeList,
					Optional: true,
					Elem: &schema.Resource{
						Schema: map[string]*schema.Schema{
							"type": {
								Type:     schema.TypeString,
								Optional: true,
							},
							"namespace": {
								Type:     schema.TypeString,
								Optional: true,
							},
							"tags": {
								Type:     schema.TypeList,
								Optional: true,
								Elem: &schema.Schema{
									Type: schema.TypeString,
								},
							},
						},
					},
				},
			},
		},
	}

	return kinesisMetricsSource
}

func resourceSumologicKinesisMetricsSourceCreate(d *schema.ResourceData, meta interface{}) error {
	c := meta.(*Client)

	if d.Id() == "" {
		source, err := resourceToKinesisMetricsSource(d)
		if err != nil {
			return err
		}

		id, err := c.CreateKinesisMetricsSource(source, d.Get("collector_id").(int))

		if err != nil {
			return err
		}

		d.SetId(strconv.Itoa(id))
	}

	return resourceSumologicKinesisMetricsSourceRead(d, meta)
}

func resourceSumologicKinesisMetricsSourceUpdate(d *schema.ResourceData, meta interface{}) error {
	c := meta.(*Client)

	source, err := resourceToKinesisMetricsSource(d)
	if err != nil {
		return err
	}

	err = c.UpdateKinesisMetricsSource(source, d.Get("collector_id").(int))

	if err != nil {
		return err
	}

	return resourceSumologicKinesisMetricsSourceRead(d, meta)
}

func resourceSumologicKinesisMetricsSourceRead(d *schema.ResourceData, meta interface{}) error {
	c := meta.(*Client)

	id, _ := strconv.Atoi(d.Id())
	source, err := c.GetKinesisMetricsSource(d.Get("collector_id").(int), id)

	if err != nil {
		return err
	}

	if source == nil {
		log.Printf("[WARN] KinesisMetric source not found, removing from state: %v - %v", id, err)
		d.SetId("")

		return nil
	}

	kinesisMetricsResources := source.ThirdPartyRef.Resources
	path := getKinesisMetricsThirdPartyPathAttributes(kinesisMetricsResources)

	if err := d.Set("path", path); err != nil {
		return err
	}

	if err := resourceSumologicSourceRead(d, source.Source); err != nil {
		return fmt.Errorf("%s", err)
	}
	d.Set("content_type", source.ContentType)
	d.Set("message_per_request", source.MessagePerRequest)
	d.Set("url", source.URL)

	return nil
}

func resourceToKinesisMetricsSource(d *schema.ResourceData) (KinesisMetricsSource, error) {
	source := resourceToSource(d)
	source.Type = "HTTP"

	kinesisMetricsSource := KinesisMetricsSource{
		Source:            source,
		MessagePerRequest: d.Get("message_per_request").(bool),
		URL:               d.Get("url").(string),
	}

	authSettings, errAuthSettings := getPollingAuthentication(d)
	if errAuthSettings != nil {
		return kinesisMetricsSource, errAuthSettings
	}

	pathSettings, errPathSettings := getKinesisMetricsPathSettings(d)
	if errPathSettings != nil {
		return kinesisMetricsSource, errPathSettings
	}

	kinesisMetricsResource := PollingResource{
		ServiceType:    d.Get("content_type").(string),
		Authentication: authSettings,
		Path:           pathSettings,
	}

	kinesisMetricsSource.ThirdPartyRef.Resources = append(kinesisMetricsSource.ThirdPartyRef.Resources, kinesisMetricsResource)

	return kinesisMetricsSource, nil
}

func getKinesisMetricsPathSettings(d *schema.ResourceData) (PollingPath, error) {
	pathSettings := PollingPath{}
	paths := d.Get("path").([]interface{})
	if len(paths) > 0 {
		path := paths[0].(map[string]interface{})
		pathType := path["type"].(string)
		pathSettings.Type = pathType
		pathSettings.TagFilters = getPollingTagFilters(d)
	} else {
		return pathSettings, errors.New(fmt.Sprintf("[ERROR] no path specification in kinesis metric Soruce"))
	}
	return pathSettings, nil
}

func getKinesisMetricsThirdPartyPathAttributes(pollingResource []PollingResource) []map[string]interface{} {

	var s []map[string]interface{}

	for _, t := range pollingResource {
		mapping := map[string]interface{}{
			"type":        t.Path.Type,
			"tag_filters": flattenPollingTagFilters(t.Path.TagFilters),
		}
		s = append(s, mapping)
	}
	return s
}
