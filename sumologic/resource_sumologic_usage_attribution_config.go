package sumologic

import (
	"regexp"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

const usageAttributionConfigId = "usageAttributionConfig"

// Tag names are used as field names in the dashboard queries, so Sumo Logic only accepts valid field names.
const usageAttributionTagNamePattern = "^[A-Za-z_][A-Za-z0-9_]*$"

func resourceSumologicUsageAttributionConfig() *schema.Resource {
	return &schema.Resource{
		Create: resourceSumologicUsageAttributionConfigCreate,
		Read:   resourceSumologicUsageAttributionConfigRead,
		Update: resourceSumologicUsageAttributionConfigUpdate,
		Delete: resourceSumologicUsageAttributionConfigDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"enabled": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  true,
			},
			"custom_tags": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
					ValidateFunc: validation.StringMatch(
						regexp.MustCompile(usageAttributionTagNamePattern),
						"must start with a letter or underscore and contain only letters, digits and underscores",
					),
				},
			},
			"standard_tags": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"breakdown_by_options": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"dashboard_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"enabled_at": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceSumologicUsageAttributionConfigRead(d *schema.ResourceData, meta interface{}) error {
	c := meta.(*Client)

	config, err := c.GetUsageAttributionConfig()
	if err != nil {
		return err
	}

	return setUsageAttributionConfigResource(d, config)
}

func resourceSumologicUsageAttributionConfigCreate(d *schema.ResourceData, meta interface{}) error {
	// The configuration of an organization always exists and can only be updated, so create applies the given values.
	if err := resourceSumologicUsageAttributionConfigUpdate(d, meta); err != nil {
		return err
	}

	d.SetId(usageAttributionConfigId)
	return nil
}

func resourceSumologicUsageAttributionConfigUpdate(d *schema.ResourceData, meta interface{}) error {
	c := meta.(*Client)

	config, err := c.UpdateUsageAttributionConfig(resourceToUsageAttributionConfigUpdate(d))
	if err != nil {
		return err
	}

	return setUsageAttributionConfigResource(d, config)
}

func resourceSumologicUsageAttributionConfigDelete(d *schema.ResourceData, meta interface{}) error {
	c := meta.(*Client)

	// The configuration cannot be deleted. Destroying the resource turns usage attribution off. Turning it off does not
	// change the stored custom tags or the dashboard.
	_, err := c.UpdateUsageAttributionConfig(UsageAttributionConfigUpdate{
		Enabled:    false,
		CustomTags: usageAttributionTags(d.Get("custom_tags").([]interface{})),
	})
	return err
}

func resourceToUsageAttributionConfigUpdate(d *schema.ResourceData) UsageAttributionConfigUpdate {
	return UsageAttributionConfigUpdate{
		Enabled:    d.Get("enabled").(bool),
		CustomTags: usageAttributionTags(d.Get("custom_tags").([]interface{})),
	}
}

// usageAttributionTags always returns a non-nil slice so that an empty list is sent as [] and not as null.
func usageAttributionTags(names []interface{}) []UsageAttributionTag {
	tags := make([]UsageAttributionTag, 0, len(names))
	for _, name := range names {
		tags = append(tags, UsageAttributionTag{Name: name.(string), Enabled: true})
	}
	return tags
}

func usageAttributionTagNames(tags []UsageAttributionTag) []string {
	names := make([]string, 0, len(tags))
	for _, tag := range tags {
		names = append(names, tag.Name)
	}
	return names
}

func setUsageAttributionConfigResource(d *schema.ResourceData, config *UsageAttributionConfig) error {
	if err := d.Set("enabled", config.Enabled); err != nil {
		return err
	}
	if err := d.Set("custom_tags", usageAttributionTagNames(config.CustomTags)); err != nil {
		return err
	}
	if err := d.Set("standard_tags", usageAttributionTagNames(config.StandardTags)); err != nil {
		return err
	}
	if err := d.Set("breakdown_by_options", config.BreakdownByOptions); err != nil {
		return err
	}
	if err := d.Set("dashboard_id", config.DashboardId); err != nil {
		return err
	}
	return d.Set("enabled_at", config.EnabledAt)
}
