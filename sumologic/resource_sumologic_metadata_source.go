package sumologic

import (
	"fmt"
	"log"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

func resourceSumologicMetadataSource() *schema.Resource {
	pollingMetadataSource := resourceSumologicSource()
	pollingMetadataSource.Create = resourceSumologicMetadataSourceCreate
	pollingMetadataSource.Read = resourceSumologicMetadataSourceRead
	pollingMetadataSource.Update = resourceSumologicMetadataSourceUpdate
	pollingMetadataSource.Importer = &schema.ResourceImporter{
		State: resourceSumologicSourceImport,
	}

	pollingMetadataSource.Schema["content_type"] = &schema.Schema{
		Type:         schema.TypeString,
		Required:     true,
		ForceNew:     true,
		ValidateFunc: validation.StringInSlice([]string{"AwsMetadata"}, false),
	}
	pollingMetadataSource.Schema["scan_interval"] = &schema.Schema{
		Type:     schema.TypeInt,
		Required: true,
	}
	pollingMetadataSource.Schema["paused"] = &schema.Schema{
		Type:     schema.TypeBool,
		Required: true,
	}
	pollingMetadataSource.Schema["url"] = &schema.Schema{
		Type:     schema.TypeString,
		Computed: true,
	}
	pollingMetadataSource.Schema["authentication"] = &schema.Schema{
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
	pollingMetadataSource.Schema["path"] = &schema.Schema{
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
					ValidateFunc: validation.StringInSlice([]string{"AwsMetadataPath"}, false),
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
						Type:         schema.TypeString,
						ValidateFunc: validation.StringInSlice([]string{"AWS/EC2"}, false),
					},
				},

				"tag_filters": {
					Type:     schema.TypeList,
					Optional: true,
					Elem: &schema.Schema{
						Type: schema.TypeString,
					},
				},
			},
		},
	}

	return pollingMetadataSource
}

func resourceSumologicMetadataSourceCreate(d *schema.ResourceData, meta interface{}) error {

	c := meta.(*Client)

	if d.Id() == "" {
		source := resourceToMetadataSource(d)
		sourceID, err := c.CreateMetadataSource(source, d.Get("collector_id").(int))

		if err != nil {
			return err
		}

		id := strconv.Itoa(sourceID)

		d.SetId(id)
	}

	return resourceSumologicMetadataSourceRead(d, meta)
}

func resourceSumologicMetadataSourceUpdate(d *schema.ResourceData, meta interface{}) error {
	c := meta.(*Client)

	source := resourceToMetadataSource(d)

	err := c.UpdateMetadataSource(source, d.Get("collector_id").(int))

	if err != nil {
		return err
	}

	return resourceSumologicMetadataSourceRead(d, meta)
}

func resourceSumologicMetadataSourceRead(d *schema.ResourceData, meta interface{}) error {
	c := meta.(*Client)

	id, _ := strconv.Atoi(d.Id())
	source, err := c.GetMetadataSource(d.Get("collector_id").(int), id)

	if err != nil {
		return err
	}

	if source == nil {
		log.Printf("[WARN] Polling source not found, removing from state: %v - %v", id, err)
		d.SetId("")

		return nil
	}

	pollingResources := source.ThirdPartyRef.Resources
	path := getMetadataThirdPartyPathAttributes(pollingResources)

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

func resourceToMetadataSource(d *schema.ResourceData) MetadataSource {
	source := resourceToSource(d)
	source.Type = "Polling"

	pollingMetadataSource := MetadataSource{
		Source:       source,
		Paused:       d.Get("paused").(bool),
		ScanInterval: d.Get("scan_interval").(int),
		ContentType:  d.Get("content_type").(string),
		URL:          d.Get("url").(string),
	}

	pollingResource := MetadataResource{
		ServiceType:    d.Get("content_type").(string),
		Authentication: getMetadataAuthentication(d),
		Path:           getMetadataPathSettings(d),
	}

	pollingMetadataSource.ThirdPartyRef.Resources = append(pollingMetadataSource.ThirdPartyRef.Resources, pollingResource)

	return pollingMetadataSource
}

func getMetadataThirdPartyPathAttributes(pollingResource []MetadataResource) []map[string]interface{} {

	var s []map[string]interface{}

	for _, t := range pollingResource {
		mapping := map[string]interface{}{
			"type":                t.Path.Type,
			"limit_to_regions":    t.Path.LimitToRegions,
			"limit_to_namespaces": t.Path.LimitToNamespaces,
			"tag_filters":         t.Path.TagFilters,
		}
		s = append(s, mapping)
	}
	return s
}

func getMetadataAuthentication(d *schema.ResourceData) MetadataAuthentication {
	auths := d.Get("authentication").([]interface{})
	authSettings := MetadataAuthentication{}

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

func getMetadataPathSettings(d *schema.ResourceData) MetadataPath {
	pathSettings := MetadataPath{}
	paths := d.Get("path").([]interface{})

	if len(paths) > 0 {
		path := paths[0].(map[string]interface{})
		switch pathType := path["type"].(string); pathType {
		case "AwsMetadataPath":
			pathSettings.Type = "AwsMetadataPath"
			rawLimitToRegions := path["limit_to_regions"].([]interface{})
			LimitToRegions := make([]string, len(rawLimitToRegions))
			for i, v := range rawLimitToRegions {
				LimitToRegions[i] = v.(string)
			}

			rawtagFilters := path["tag_filters"].([]interface{})
			TagFilters := make([]string, len(rawtagFilters))
			for i, v := range rawtagFilters {
				TagFilters[i] = v.(string)
			}
			pathSettings.LimitToRegions = LimitToRegions
			pathSettings.LimitToNamespaces = []string{"AWS/EC2"}
			pathSettings.TagFilters = TagFilters
		default:
			log.Printf("[ERROR] Unknown resourceType in path: %v", pathType)
		}
	}

	return pathSettings
}
