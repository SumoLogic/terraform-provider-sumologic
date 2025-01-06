package sumologic

import (
	"fmt"
	"log"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

func resourceSumologicPollingSource() *schema.Resource {
	pollingSource := resourceSumologicSource()
	pollingSource.Create = resourceSumologicPollingSourceCreate
	pollingSource.Read = resourceSumologicPollingSourceRead
	pollingSource.Update = resourceSumologicPollingSourceUpdate
	pollingSource.Importer = &schema.ResourceImporter{
		State: resourceSumologicSourceImport,
	}
	pollingSource.DeprecationMessage =
		"We are deprecating the generic sumologic polling source and in turn creating individual sources for each of the content_type currently supported."

	pollingSource.Schema["content_type"] = &schema.Schema{
		Type:         schema.TypeString,
		Required:     true,
		ForceNew:     true,
		ValidateFunc: validation.StringInSlice([]string{"AwsS3Bucket", "AwsElbBucket", "AwsCloudFrontBucket", "AwsCloudTrailBucket", "AwsS3AuditBucket", "AwsCloudWatch"}, false),
	}
	pollingSource.Schema["scan_interval"] = &schema.Schema{
		Type:     schema.TypeInt,
		Required: true,
	}
	pollingSource.Schema["paused"] = &schema.Schema{
		Type:     schema.TypeBool,
		Required: true,
	}
	pollingSource.Schema["url"] = &schema.Schema{
		Type:     schema.TypeString,
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
	pollingSource.Schema["path"] = &schema.Schema{
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
					ValidateFunc: validation.StringInSlice([]string{"S3BucketPathExpression", "CloudWatchPath"}, false),
				},
				"bucket_name": {
					Type:     schema.TypeString,
					Optional: true,
				},
				"path_expression": {
					Type:     schema.TypeString,
					Optional: true,
				},
				"limit_to_regions": {
					Type:     schema.TypeList,
					Optional: true,
					Elem: &schema.Schema{
						Type: schema.TypeString,
					},
				},
				"limit_to_namespaces": {
					Type:     schema.TypeList,
					Optional: true,
					Elem: &schema.Schema{
						Type: schema.TypeString,
					},
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
									Type: schema.TypeMap, // Accept both maps (for objects) and strings
								},
								ValidateFunc: validateTags,
							},
						},
					},
				},
			},
		},
	}

	return pollingSource
}

func validateTags(val interface{}, key string) ([]string, []error) {
	v := val.(map[string]interface{})
	var errs []error

	if tags, ok := v.([]interface{}); ok {
		for i, tag := range tags {
			switch t := tag.(type) {
			case map[string]interface{}:
				// Validate object structure
				// Validate "name" to be a string
				if _, ok := t["name"]; !ok {
					errors = append(errors, fmt.Errorf("%s[%d]: missing required field 'name'", key, i))
				} else if _, ok := t["name"].(string); !ok {
					errors = append(errors, fmt.Errorf("%s[%d]: 'name' must be a string", key, i))
				}

				// Validate "values" to be a list of strings
				if _, ok := t["values"]; !ok {
					errors = append(errors, fmt.Errorf("%s[%d]: missing required field 'values'", key, i))
				} else if values, ok := t["values"].([]interface{}); !ok {
					errors = append(errors, fmt.Errorf("%s[%d]: 'values' must be a list of strings", key, i))
				} else {
					for j, value := range values {
						if _, ok := value.(string); !ok {
							errors = append(errors, fmt.Errorf("%s[%d].values[%d]: must be a string", key, i, j))
						}
					}
				}
			case string:
				continue
			default:
				errors = append(errors, fmt.Errorf("%s[%d]: must be either a string or an object with 'name' and 'values'", key, i))
			}
		}
	}
}

func resourceSumologicPollingSourceCreate(d *schema.ResourceData, meta interface{}) error {

	c := meta.(*Client)

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

	if err := resourceSumologicSourceRead(d, source.Source); err != nil {
		return fmt.Errorf("%s", err)
	}
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
			"type":                t.Path.Type,
			"bucket_name":         t.Path.BucketName,
			"path_expression":     t.Path.PathExpression,
			"limit_to_regions":    t.Path.LimitToRegions,
			"limit_to_namespaces": t.Path.LimitToNamespaces,
			"tag_filters":         flattenTagFilters(t.Path.TagFilters),
		}
		s = append(s, mapping)
	}
	return s
}

func flattenTagFilters(v []TagFilter) []map[string]interface{} {
	var filters []map[string]interface{}
	for _, d := range v {
		filter := map[string]interface{}{
			"type":      d.Type,
			"namespace": d.Namespace,
			"tags":      d.Tags,
		}
		filters = append(filters, filter)
	}

	return filters
}

func getTagFilters(d *schema.ResourceData) []TagFilter {
	paths := d.Get("path").([]interface{})
	path := paths[0].(map[string]interface{})
	rawTagFilterConfig := path["tag_filters"].([]interface{})
	var filters []TagFilter

	for _, rawConfig := range rawTagFilterConfig {
		config := rawConfig.(map[string]interface{})
		filter := TagFilter{}
		filter.Type = config["type"].(string)
		filter.Namespace = config["namespace"].(string)

		rawTags := config["tags"].([]interface{})
		Tags := make([]string, len(rawTags))
		for i, v := range rawTags {
			Tags[i] = v.(string)
		}
		filter.Tags = Tags
		filters = append(filters, filter)
	}

	return filters
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
		switch pathType := path["type"].(string); pathType {
		case "S3BucketPathExpression":
			pathSettings.Type = "S3BucketPathExpression"
			pathSettings.BucketName = path["bucket_name"].(string)
			pathSettings.PathExpression = path["path_expression"].(string)
		case "CloudWatchPath":
			pathSettings.Type = "CloudWatchPath"
			rawLimitToRegions := path["limit_to_regions"].([]interface{})
			LimitToRegions := make([]string, len(rawLimitToRegions))
			for i, v := range rawLimitToRegions {
				LimitToRegions[i] = v.(string)
			}

			rawLimitToNamespaces := path["limit_to_namespaces"].([]interface{})
			LimitToNamespaces := make([]string, len(rawLimitToNamespaces))
			for i, v := range rawLimitToNamespaces {
				LimitToNamespaces[i] = v.(string)
			}
			pathSettings.LimitToRegions = LimitToRegions
			pathSettings.LimitToNamespaces = LimitToNamespaces
			pathSettings.TagFilters = getTagFilters(d)
		default:
			log.Printf("[ERROR] Unknown resourceType in path: %v", pathType)
		}
	}

	return pathSettings
}
