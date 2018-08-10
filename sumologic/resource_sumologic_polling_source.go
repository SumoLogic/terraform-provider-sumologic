package sumologic

import (
	"log"
	"strconv"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/helper/validation"
)

func resourceSumologicPollingSource() *schema.Resource {
	pollingSource := resourceSumologicSource()
	pollingSource.Create = resourceSumologicPollingSourceCreate
	pollingSource.Read = resourceSumologicPollingSourceRead
	pollingSource.Update = resourceSumologicPollingSourceUpdate

	pollingSource.Schema["content_type"] = &schema.Schema{
		Type:     schema.TypeString,
		Optional: true,
		ForceNew: true,
		ValidateFunc: validation.StringInSlice([]string{"AwsS3Bucket", "AwsElbBucket", "AwsCloudFrontBucket", "AwsCloudTrailBucket", "AwsS3AuditBucket"}, false),
	}
	pollingSource.Schema["scan_interval"] = &schema.Schema{
		Type:     schema.TypeInt,
		Required: true,
		ForceNew: false,
		Default:  300000,
	}
	pollingSource.Schema["paused"] = &schema.Schema{
		Type:     schema.TypeBool,
		Required: true,
		ForceNew: false,
		Default:  false,
	}
	pollingSource.Schema["third_party_ref"] = &schema.Schema{
		Type:     schema.TypeMap,
		Required: true,
		ForceNew: true,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"resources": {
					Type:     schema.TypeList,
					Required: true,
					ForceNew: false,
					Elem: &schema.Resource{
						Schema: map[string]*schema.Schema{
							"service_type": {
								//Should have the same value as "content_type"
								Type:     schema.TypeString,
								Optional: true,
								ForceNew: true,
								ValidateFunc: validation.StringInSlice([]string{"AwsS3Bucket", "AwsElbBucket", "AwsCloudFrontBucket", "AwsCloudTrailBucket", "AwsS3AuditBucket"}, false),
							},
							"path": {
								Type:     schema.TypeList,
								Required: true,
								ForceNew: true,
								MinItems: 1,
								MaxItems: 1,
								Elem: &schema.Resource{
									Schema: map[string]*schema.Schema{
										"bucket_name": {
											Type:     schema.TypeString,
											Required: true,
											ForceNew: false,
										},
										"path_expression": {
											Type:     schema.TypeString,
											Required: true,
											ForceNew: false,
										},
									},
								},
							},
							"authentication": {
								Type:     schema.TypeList,
								Required: true,
								ForceNew: true,
								MinItems: 1,
								MaxItems: 1,
								Elem: &schema.Resource{
									Schema: map[string]*schema.Schema{
										"access_key": {
											Type:     schema.TypeString,
											Optional: true,
											ForceNew: false,
										},
										"secret_key": {
											Type:     schema.TypeString,
											Optional: true,
											ForceNew: false,
										},
									},
								},
							},
						},
					},
				},
			},
		},
	}
	return pollingSource
}

func resourceSumologicPollingSourceCreate(d *schema.ResourceData, meta interface{}) error {

	c := meta.(*Client)

	if d.Get("lookup_by_name").(bool) {
		source, err := c.GetSourceName(d.Get("collector_id").(int), d.Get("name").(string))

		if err != nil {
			return err
		}

		// Set ID of source if it exists
		if source != nil {
			d.SetId(strconv.Itoa(source.ID))
		}
	}

	// If source ID is still empty, create a new source
	if d.Id() == "" {
		sourceID, err := c.CreatePollingSource(
			d.Get("name").(string),
			d.Get("description").(string),
			d.Get("content_type").(string),
			d.Get("category").(string),
			d.Get("scan_interval").(int),
			d.Get("paused").(bool),
			d.Get("collector_id").(int),
			getAuthentication(d),
			getPathSettings(d),
		)

		if err != nil {
			return err
		}

		id := strconv.Itoa(sourceID)

		d.SetId(id)
	}

	return resourceSumologicPollingSourceRead(d, meta)
}

func resourceSumologicPollingSourceRead(d *schema.ResourceData, meta interface{}) error {

	c := meta.(*Client)

	id, _ := strconv.Atoi(d.Id())
	collectorID := d.Get("collector_id").(int)

	source, err := c.GetPollingSource(collectorID, id)

	// Source is gone, remove it from state
	if err != nil {
		log.Printf("Polling source %v: %v", id, err)
		d.SetId("")

		return nil
	}

	pollingResources := source.ThirdPartyRef.Resources
	path := getThirdyPartyPathAttributes(pollingResources)

	if err := d.Set("path", path); err != nil {
		return err
	}

	d.Set("name", source.Name)
	d.Set("description", source.Description)
	d.Set("content_type", source.ContentType)
	d.Set("category", source.Category)
	d.Set("scan_interval", source.ScanInterval)
	d.Set("paused", source.Paused)

	return nil
}

func resourceSumologicPollingSourceDelete(d *schema.ResourceData, meta interface{}) error {
	c := meta.(*Client)

	// Destroy source if `destroy` is true, otherwise ignore
	if d.Get("destroy").(bool) {
		id, _ := strconv.Atoi(d.Id())
		collectorID, _ := d.Get("collector_id").(int)

		return c.DestroySource(id, collectorID)
	}

	return nil
}

func resourceSumologicPollingSourceUpdate(d *schema.ResourceData, meta interface{}) error {
	c := meta.(*Client)

	source := resourceToPollingSource(d)

	err := c.UpdatePollingSource(source, d.Get("collector_id").(int))

	if err != nil {
		return err
	}

	return resourceSumologicPollingSourceRead(d, meta)
}

func resourceToPollingSource(d *schema.ResourceData) PollingSource {

	id, _ := strconv.Atoi(d.Id())
	source := PollingSource{}
	pollingResource := PollingResource{}

	source.ID = id
	source.Type = "Polling"
	source.Name = d.Get("name").(string)
	source.Description = d.Get("description").(string)
	source.Category = d.Get("category").(string)
	source.Paused = d.Get("paused").(bool)
	source.ScanInterval = d.Get("scan_interval").(int)
	source.ContentType = d.Get("content_type").(string)

	pollingResource.ServiceType = "AwsS3AuditBucket"
	pollingResource.Authentication = getAuthentication(d)
	pollingResource.Path = getPathSettings(d)

	source.ThirdPartyRef.Resources = append(source.ThirdPartyRef.Resources, pollingResource)

	return source
}

func getThirdyPartyPathAttributes(pollingResource []PollingResource) []map[string]interface{} {

	var s []map[string]interface{}
	for _, t := range pollingResource {
		mapping := map[string]interface{}{
			"bucket_name":     t.Path.BucketName,
			"path_expression": t.Path.PathExpression,
		}
		s = append(s, mapping)
	}

	return s
}

func getAuthentication(d *schema.ResourceData) PollingAuthentication {

	auths := d.Get("authentication").([]interface{})
	authSettings := PollingAuthentication{}

	if len(auths) > 0 {
		auth := auths[0].(map[string]interface{})
		authSettings.Type = "S3BucketAuthentication"
		authSettings.AwsID = auth["access_key"].(string)
		authSettings.AwsKey = auth["secret_key"].(string)
	}

	return authSettings
}

func getPathSettings(d *schema.ResourceData) PollingPath {
	pathSettings := PollingPath{}
	paths := d.Get("path").([]interface{})

	if len(paths) > 0 {
		path := paths[0].(map[string]interface{})
		pathSettings.Type = "S3BucketPathExpression"
		pathSettings.BucketName = path["bucket_name"].(string)
		pathSettings.PathExpression = path["path_expression"].(string)
	}

	return pathSettings
}
