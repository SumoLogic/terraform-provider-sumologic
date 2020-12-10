package sumologic

import (
	"errors"
	"fmt"
	"log"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

func resourceSumologicGCPSource() *schema.Resource {
	gcpSource := resourceSumologicSource()
	gcpSource.Create = resourceSumologicGCPSourceCreate
	gcpSource.Read = resourceSumologicGCPSourceRead
	gcpSource.Update = resourceSumologicGCPSourceUpdate
	gcpSource.Importer = &schema.ResourceImporter{
		State: resourceSumologicSourceImport,
	}

	gcpSource.Schema["content_type"] = &schema.Schema{
		Type:         schema.TypeString,
		Optional:     true,
		ForceNew:     true,
		Default:      "GoogleCloudLogs",
		ValidateFunc: validation.StringInSlice([]string{"GoogleCloudLogs"}, false),
	}
	gcpSource.Schema["message_per_request"] = &schema.Schema{
		Type:     schema.TypeBool,
		Optional: true,
		Default:  false,
	}
	gcpSource.Schema["url"] = &schema.Schema{
		Type:     schema.TypeString,
		Computed: true,
	}
	gcpSource.Schema["authentication"] = &schema.Schema{
		Type:     schema.TypeList,
		Optional: true,
		ForceNew: true,
		MinItems: 1,
		MaxItems: 1,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"type": {
					Type:         schema.TypeString,
					Optional:     true,
					Default:      "NoAuthentication",
					ValidateFunc: validation.StringInSlice([]string{"NoAuthentication"}, false),
				},
			},
		},
	}
	gcpSource.Schema["path"] = &schema.Schema{
		Type:     schema.TypeList,
		Optional: true,
		ForceNew: true,
		MinItems: 1,
		MaxItems: 1,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"type": {
					Type:         schema.TypeString,
					Optional:     true,
					Default:      "NoPathExpression",
					ValidateFunc: validation.StringInSlice([]string{"NoPathExpression"}, false),
				},
			},
		},
	}

	return gcpSource
}

func resourceSumologicGCPSourceCreate(d *schema.ResourceData, meta interface{}) error {
	c := meta.(*Client)

	if d.Id() == "" {
		source, err := resourceToGCPSource(d)
		if err != nil {
			return err
		}

		sourceID, err := c.CreateGCPSource(source, d.Get("collector_id").(int))
		if err != nil {
			return err
		}

		d.SetId(strconv.Itoa(sourceID))
	}

	return resourceSumologicGCPSourceRead(d, meta)
}

func resourceSumologicGCPSourceUpdate(d *schema.ResourceData, meta interface{}) error {
	c := meta.(*Client)

	source, err := resourceToGCPSource(d)
	if err != nil {
		return err
	}

	err = c.UpdateGCPSource(source, d.Get("collector_id").(int))
	if err != nil {
		return err
	}

	return resourceSumologicGCPSourceRead(d, meta)
}

func resourceSumologicGCPSourceRead(d *schema.ResourceData, meta interface{}) error {
	c := meta.(*Client)

	id, _ := strconv.Atoi(d.Id())
	source, err := c.GetGCPSource(d.Get("collector_id").(int), id)

	if err != nil {
		return err
	}

	if source == nil {
		log.Printf("[WARN] GCP source not found, removing from state: %v - %v", id, err)
		d.SetId("")

		return nil
	}

	GCPResources := source.ThirdPartyRef.Resources
	path := getGCPThirdPartyPathAttributes(GCPResources)

	if err := d.Set("path", path); err != nil {
		return err
	}

	if err := resourceSumologicSourceRead(d, source.Source); err != nil {
		return fmt.Errorf("%s", err)
	}
	d.Set("content_type", source.ContentType)
	d.Set("url", source.URL)

	return nil
}

func resourceToGCPSource(d *schema.ResourceData) (GCPSource, error) {
	source := resourceToSource(d)
	source.Type = "HTTP"

	gcpSource := GCPSource{
		Source:            source,
		MessagePerRequest: d.Get("message_per_request").(bool),
	}

	authSettings, errAuthSettings := getGCPAuthentication(d)
	if errAuthSettings != nil {
		return gcpSource, errAuthSettings
	}

	pathSettings, errPathSettings := getGCPPathSettings(d)
	if errPathSettings != nil {
		return gcpSource, errPathSettings
	}

	GCPResource := GCPResource{
		ServiceType:    d.Get("content_type").(string),
		Authentication: authSettings,
		Path:           pathSettings,
	}

	gcpSource.ThirdPartyRef.Resources = append(gcpSource.ThirdPartyRef.Resources, GCPResource)

	return gcpSource, nil
}

func getGCPThirdPartyPathAttributes(GCPResource []GCPResource) []map[string]interface{} {

	var s []map[string]interface{}

	for _, t := range GCPResource {
		mapping := map[string]interface{}{
			"type": t.Path.Type,
		}
		s = append(s, mapping)
	}
	return s
}

func getGCPAuthentication(d *schema.ResourceData) (GCPAuthentication, error) {
	auths := d.Get("authentication").([]interface{})
	authSettings := GCPAuthentication{}

	if len(auths) > 0 {
		auth := auths[0].(map[string]interface{})
		switch authType := auth["type"].(string); authType {
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

func getGCPPathSettings(d *schema.ResourceData) (GCPPath, error) {
	pathSettings := GCPPath{}
	paths := d.Get("path").([]interface{})

	if len(paths) > 0 {
		path := paths[0].(map[string]interface{})
		switch pathType := path["type"].(string); pathType {
		case "NoPathExpression":
			pathSettings.Type = "NoPathExpression"
		default:
			errorMessage := fmt.Sprintf("[ERROR] Unknown resourceType in path: %v", pathType)
			log.Print(errorMessage)
			return pathSettings, errors.New(errorMessage)
		}
	}

	return pathSettings, nil
}
