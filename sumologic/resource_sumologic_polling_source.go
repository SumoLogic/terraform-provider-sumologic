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
		Type:         schema.TypeString,
		Required:     true,
		ForceNew:     true,
		ValidateFunc: validation.StringInSlice([]string{"AwsS3Bucket", "AwsElbBucket", "AwsCloudFrontBucket", "AwsCloudTrailBucket", "AwsS3AuditBucket"}, false),
	}
	pollingSource.Schema["scan_interval"] = &schema.Schema{
		Type:     schema.TypeInt,
		Required: true,
		ForceNew: false,
	}
	pollingSource.Schema["paused"] = &schema.Schema{
		Type:     schema.TypeBool,
		Required: true,
		ForceNew: false,
	}
	pollingSource.Schema["url"] = &schema.Schema{
		Type:     schema.TypeString,
		Required: false,
		ForceNew: false,
		Computed: true,
	}
	pollingSource.Schema["authentication"] = &schema.Schema{
		Type:     schema.TypeList,
		Required: true,
		ForceNew: true,
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
					ForceNew: false,
				},
				"secret_key": {
					Type:     schema.TypeString,
					Optional: true,
					ForceNew: false,
				},
				"role_arn": {
					Type:     schema.TypeString,
					Optional: true,
					ForceNew: false,
				},
			},
		},
	}
	pollingSource.Schema["path"] = &schema.Schema{
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

		if source != nil {
			d.SetId(strconv.Itoa(source.ID))
		}
	}

	if d.Id() == "" {
		source := resourceToPollingSource(d)
		sourceID, err := c.CreatePollingSource(source, d.Get("collector_id").(int))

		if err != nil {
			return err
		}

		id := strconv.Itoa(sourceID)

		d.SetId(id)
	}

	return resourceSumologicPollingSourceRead(d, meta)
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

func resourceSumologicPollingSourceRead(d *schema.ResourceData, meta interface{}) error {
	c := meta.(*Client)

	id, _ := strconv.Atoi(d.Id())
	source, err := c.GetPollingSource(d.Get("collector_id").(int), id)

	if err != nil {
		return err
	}

	if source == nil {
		log.Printf("[WARN] Polling source not found, removing from state: %v - %v", id, err)
		d.SetId("")

		return nil
	}

	pollingResources := source.ThirdPartyRef.Resources
	path := getThirdPartyPathAttributes(pollingResources)

	if err := d.Set("path", path); err != nil {
		return err
	}

	resourceSumologicSourceRead(d, source.Source)
	d.Set("content_type", source.ContentType)
	d.Set("scan_interval", source.ScanInterval)
	d.Set("paused", source.Paused)
	d.Set("url", source.URL)

	return nil
}

func resourceToPollingSource(d *schema.ResourceData) PollingSource {
	source := resourceToSource(d)
	source.Type = "Polling"

	pollingSource := PollingSource{
		Source:       source,
		Paused:       d.Get("paused").(bool),
		ScanInterval: d.Get("scan_interval").(int),
		ContentType:  d.Get("content_type").(string),
		URL:          d.Get("url").(string),
	}

	pollingResource := PollingResource{
		ServiceType:    d.Get("content_type").(string),
		Authentication: getAuthentication(d),
		Path:           getPathSettings(d),
	}

	pollingSource.ThirdPartyRef.Resources = append(pollingSource.ThirdPartyRef.Resources, pollingResource)

	return pollingSource
}

func getThirdPartyPathAttributes(pollingResource []PollingResource) []map[string]interface{} {

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
		switch authType := auth["type"].(string); authType {
		case "S3BucketAuthentication":
			authSettings.Type = "S3BucketAuthentication"
			authSettings.AwsID = auth["access_key"].(string)
			authSettings.AwsKey = auth["secret_key"].(string)
		case "AWSRoleBasedAuthentication":
			authSettings.Type = "AWSRoleBasedAuthentication"
			authSettings.RoleARN = auth["role_arn"].(string)
		default:
			log.Printf("[ERROR] Unknown authType: %v", authType)
		}
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
