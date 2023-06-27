package sumologic

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"log"
)

func resourceSumologicCSEEntityNormalizationConfiguration() *schema.Resource {
	return &schema.Resource{
		Create: resourceSumologicCSEEntityNormalizationConfigurationCreate,
		Read:   resourceSumologicCSEEntityNormalizationConfigurationRead,
		Delete: resourceSumologicCSEEntityNormalizationConfigurationDelete,
		Update: resourceSumologicCSEEntityNormalizationConfigurationUpdate,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"windows_normalization_enabled": {
				Type:     schema.TypeBool,
				Required: true,
			},
			"fqdn_normalization_enabled": {
				Type:     schema.TypeBool,
				Required: true,
			},
			"aws_normalization_enabled": {
				Type:     schema.TypeBool,
				Required: true,
			},
			"default_normalized_domain": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringIsNotEmpty,
			},
			"normalize_hostnames": {
				Type:     schema.TypeBool,
				Required: true,
			},
			"normalize_usernames": {
				Type:     schema.TypeBool,
				Required: true,
			},
			"domain_mappings": {
				Type:     schema.TypeList,
				Required: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"normalized_domain": {
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: validation.StringIsNotEmpty,
						},
						"raw_domain": {
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: validation.StringIsNotEmpty,
						},
					},
				},
			},
		},
	}
}

func resourceSumologicCSEEntityNormalizationConfigurationRead(d *schema.ResourceData, meta interface{}) error {
	c := meta.(*Client)

	var CSEEntityNormalizationConfiguration *CSEEntityNormalizationConfiguration

	CSEEntityNormalizationConfiguration, err := c.GetCSEEntityNormalizationConfiguration()
	if err != nil {
		log.Printf("[WARN] CSE Entity Normalization Configuration not found, err: %v", err)

	}

	if CSEEntityNormalizationConfiguration == nil {
		log.Printf("[WARN] CSE Entity Normalization Configuration not found, removing from state: %v", err)
		d.SetId("")
		return nil
	}

	d.Set("windows_normalization_enabled", CSEEntityNormalizationConfiguration.WindowsNormalizationEnabled)
	d.Set("fqdn_normalization_enabled", CSEEntityNormalizationConfiguration.FqdnNormalizationEnabled)
	d.Set("aws_normalization_enabled", CSEEntityNormalizationConfiguration.AwsNormalizationEnabled)
	d.Set("default_normalized_domain", CSEEntityNormalizationConfiguration.DefaultNormalizedDomain)
	d.Set("normalized_hostnames", CSEEntityNormalizationConfiguration.NormalizeHostnames)
	d.Set("normalized_usernames", CSEEntityNormalizationConfiguration.NormalizeUsernames)
	d.Set("domain_mappings", domainMappingsToResource(CSEEntityNormalizationConfiguration.DomainMappings))

	return nil
}

func domainMappingsToResource(domainMappings []DomainMapping) []map[string]interface{} {
	result := make([]map[string]interface{}, len(domainMappings))

	for i, domainMapping := range domainMappings {
		result[i] = map[string]interface{}{
			"normalized_domain": domainMapping.NormalizedDomain,
			"raw_domain":        domainMapping.RawDomain,
		}
	}

	return result
}

func resourceSumologicCSEEntityNormalizationConfigurationDelete(d *schema.ResourceData, meta interface{}) error {
	CSEEntityNormalizationConfiguration := CSEEntityNormalizationConfiguration{false, false, false, "", false, false, []DomainMapping{}}

	c := meta.(*Client)
	err := c.UpdateCSEEntityNormalizationConfiguration(CSEEntityNormalizationConfiguration)
	if err != nil {
		return err
	}

	return resourceSumologicCSEEntityNormalizationConfigurationRead(d, meta)
}

func resourceSumologicCSEEntityNormalizationConfigurationCreate(d *schema.ResourceData, meta interface{}) error {
	//we are not really creating new object in backend, using constant id for terraform resource
	d.SetId("cse-entity-normalization-configuration")
	return resourceSumologicCSEEntityNormalizationConfigurationUpdate(d, meta)
}

func resourceSumologicCSEEntityNormalizationConfigurationUpdate(d *schema.ResourceData, meta interface{}) error {
	CSEEntityNormalizationConfiguration, err := resourceToCSEEntityNormalizationConfiguration(d)
	if err != nil {
		return err
	}

	c := meta.(*Client)
	if err = c.UpdateCSEEntityNormalizationConfiguration(CSEEntityNormalizationConfiguration); err != nil {
		return err
	}

	return resourceSumologicCSEEntityNormalizationConfigurationRead(d, meta)
}

func resourceToCSEEntityNormalizationConfiguration(d *schema.ResourceData) (CSEEntityNormalizationConfiguration, error) {
	id := d.Id()
	if id == "" {
		return CSEEntityNormalizationConfiguration{}, nil
	}

	return CSEEntityNormalizationConfiguration{
		WindowsNormalizationEnabled: d.Get("windows_normalization_enabled").(bool),
		FqdnNormalizationEnabled:    d.Get("fqdn_normalization_enabled").(bool),
		AwsNormalizationEnabled:     d.Get("aws_normalization_enabled").(bool),
		DefaultNormalizedDomain:     d.Get("default_normalized_domain").(string),
		NormalizeHostnames:          d.Get("normalize_hostnames").(bool),
		NormalizeUsernames:          d.Get("normalize_usernames").(bool),
		DomainMappings:              resourceToDomainMappingArray(d.Get("domain_mappings").([]interface{})),
	}, nil
}

func resourceToDomainMappingArray(resourceDomainMappings []interface{}) []DomainMapping {
	result := make([]DomainMapping, len(resourceDomainMappings))

	for i, resourceDomainMapping := range resourceDomainMappings {
		result[i] = DomainMapping{
			NormalizedDomain: resourceDomainMapping.(map[string]interface{})["normalized_domain"].(string),
			RawDomain:        resourceDomainMapping.(map[string]interface{})["raw_domain"].(string),
		}
	}

	return result
}
