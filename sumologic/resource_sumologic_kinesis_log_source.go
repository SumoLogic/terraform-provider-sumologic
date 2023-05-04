package sumologic

import (
	"errors"
	"fmt"
	"log"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

func resourceSumologicKinesisLogSource() *schema.Resource {
	kinesisLogSource := resourceSumologicSource()
	kinesisLogSource.Create = resourceSumologicKinesisLogSourceCreate
	kinesisLogSource.Read = resourceSumologicKinesisLogSourceRead
	kinesisLogSource.Update = resourceSumologicKinesisLogSourceUpdate
	kinesisLogSource.Importer = &schema.ResourceImporter{
		State: resourceSumologicSourceImport,
	}

	kinesisLogSource.Schema["content_type"] = &schema.Schema{
		Type:         schema.TypeString,
		Required:     true,
		ValidateFunc: validation.StringInSlice([]string{"KinesisLog"}, false),
	}

	kinesisLogSource.Schema["message_per_request"] = &schema.Schema{
		Type:     schema.TypeBool,
		Optional: true,
		Default:  false,
	}
	kinesisLogSource.Schema["url"] = &schema.Schema{
		Type:     schema.TypeString,
		Computed: true,
	}
	kinesisLogSource.Schema["authentication"] = &schema.Schema{
		Type:     schema.TypeList,
		Required: true,
		MaxItems: 1,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"type": {
					Type:         schema.TypeString,
					Required:     true,
					ValidateFunc: validation.StringInSlice([]string{"S3BucketAuthentication", "AWSRoleBasedAuthentication", "NoAuthentication"}, false),
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

	kinesisLogSource.Schema["path"] = &schema.Schema{
		Type:     schema.TypeList,
		Required: true,
		MaxItems: 1,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"type": {
					Type:         schema.TypeString,
					Required:     true,
					ValidateFunc: validation.StringInSlice([]string{"KinesisLogPath", "NoPathExpression"}, false),
				},
				"bucket_name": {
					Type:     schema.TypeString,
					Optional: true,
				},
				"path_expression": {
					Type:     schema.TypeString,
					Optional: true,
				},
				"scan_interval": {
					Type:     schema.TypeInt,
					Optional: true,
					Default:  300000,
				},
			},
		},
	}

	return kinesisLogSource
}

func resourceSumologicKinesisLogSourceCreate(d *schema.ResourceData, meta interface{}) error {
	c := meta.(*Client)

	if d.Id() == "" {
		source, err := resourceToKinesisLogSource(d)
		if err != nil {
			return err
		}

		id, err := c.CreateKinesisLogSource(source, d.Get("collector_id").(int))

		if err != nil {
			return err
		}

		d.SetId(strconv.Itoa(id))
	}

	return resourceSumologicKinesisLogSourceRead(d, meta)
}

func resourceSumologicKinesisLogSourceUpdate(d *schema.ResourceData, meta interface{}) error {
	c := meta.(*Client)

	source, err := resourceToKinesisLogSource(d)
	if err != nil {
		return err
	}

	err = c.UpdateKinesisLogSource(source, d.Get("collector_id").(int))

	if err != nil {
		return err
	}

	return resourceSumologicKinesisLogSourceRead(d, meta)
}

func resourceSumologicKinesisLogSourceRead(d *schema.ResourceData, meta interface{}) error {
	c := meta.(*Client)

	id, _ := strconv.Atoi(d.Id())
	source, err := c.GetKinesisLogSource(d.Get("collector_id").(int), id)

	if err != nil {
		return err
	}

	if source == nil {
		log.Printf("[WARN] KinesisLog source not found, removing from state: %v - %v", id, err)
		d.SetId("")

		return nil
	}

	kinesisLogResources := source.ThirdPartyRef.Resources
	path := getKinesisLogThirdPartyPathAttributes(kinesisLogResources)

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

func resourceToKinesisLogSource(d *schema.ResourceData) (KinesisLogSource, error) {
	source := resourceToSource(d)
	source.Type = "HTTP"

	kinesisLogSource := KinesisLogSource{
		Source:            source,
		MessagePerRequest: d.Get("message_per_request").(bool),
		URL:               d.Get("url").(string),
	}

	authSettings, errAuthSettings := getKinesisLogAuthentication(d)
	if errAuthSettings != nil {
		return kinesisLogSource, errAuthSettings
	}

	pathSettings, errPathSettings := getKinesisLogPathSettings(d)
	if errPathSettings != nil {
		return kinesisLogSource, errPathSettings
	}

	kinesisLogResource := KinesisLogResource{
		ServiceType:    d.Get("content_type").(string),
		Authentication: authSettings,
		Path:           pathSettings,
	}

	kinesisLogSource.ThirdPartyRef.Resources = append(kinesisLogSource.ThirdPartyRef.Resources, kinesisLogResource)

	return kinesisLogSource, nil
}

func getKinesisLogPathSettings(d *schema.ResourceData) (KinesisLogPath, error) {
	pathSettings := KinesisLogPath{}
	paths := d.Get("path").([]interface{})
	if len(paths) > 0 {
		path := paths[0].(map[string]interface{})
		switch pathType := path["type"].(string); pathType {
		case "KinesisLogPath":
			pathSettings.Type = pathType
			pathSettings.BucketName = path["bucket_name"].(string)
			pathSettings.PathExpression = path["path_expression"].(string)
			pathSettings.ScanInterval = path["scan_interval"].(int)
		case "NoPathExpression":
			pathSettings.Type = pathType
		default:
			errorMessage := fmt.Sprintf("[ERROR] Unknown resourceType in path: %v", pathType)
			log.Print(errorMessage)
			return pathSettings, errors.New(errorMessage)
		}
	} else {
		return pathSettings, errors.New(fmt.Sprintf("[ERROR] no path specification in kinesis log Soruce"))
	}
	return pathSettings, nil
}

func getKinesisLogAuthentication(d *schema.ResourceData) (PollingAuthentication, error) {
	auths := d.Get("authentication").([]interface{})
	authSettings := PollingAuthentication{}

	if len(auths) > 0 {
		auth := auths[0].(map[string]interface{})
		switch authType := auth["type"].(string); authType {
		case "S3BucketAuthentication":
			authSettings.Type = "S3BucketAuthentication"
			authSettings.AwsID = auth["access_key"].(string)
			authSettings.AwsKey = auth["secret_key"].(string)
			if auth["region"] != nil {
				authSettings.Region = auth["region"].(string)
			}
		case "AWSRoleBasedAuthentication":
			authSettings.Type = "AWSRoleBasedAuthentication"
			authSettings.RoleARN = auth["role_arn"].(string)
			if auth["region"] != nil {
				authSettings.Region = auth["region"].(string)
			}
		case "NoAuthentication":
			authSettings.Type = "NoAuthentication"
		default:
			errorMessage := fmt.Sprintf("[ERROR] Unknown authType: %v", authType)
			log.Print(errorMessage)
			return authSettings, errors.New(errorMessage)
		}
	}

	return authSettings, nil
}

func getKinesisLogThirdPartyPathAttributes(kinesisLogResource []KinesisLogResource) []map[string]interface{} {

	var s []map[string]interface{}

	for _, t := range kinesisLogResource {
		mapping := map[string]interface{}{
			"type":            t.Path.Type,
			"bucket_name":     t.Path.BucketName,
			"path_expression": t.Path.PathExpression,
			"scan_interval":   t.Path.ScanInterval,
		}
		s = append(s, mapping)
	}
	return s
}
