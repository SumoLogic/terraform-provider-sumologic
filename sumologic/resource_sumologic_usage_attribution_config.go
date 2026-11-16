package sumologic

import (
	"context"
	"fmt"
	"regexp"
	"strings"

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
		CustomizeDiff: func(ctx context.Context, d *schema.ResourceDiff, meta interface{}) error {
			// Sumo Logic keeps a repeated name once, so the list it returns would differ from the configuration on
			// every plan. Reject it where the user can fix it.
			if duplicates := repeatedNames(d.Get("custom_tags").([]interface{})); len(duplicates) > 0 {
				return fmt.Errorf("custom_tags must not repeat a name, repeated: %s", strings.Join(duplicates, ", "))
			}
			return nil
		},

		Schema: map[string]*schema.Schema{
			"enabled": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     true,
				Description: "Turns usage attribution on (true) or off (false). Defaults to true. Turning it on starts attribution for data ingested from then on. Turning it off stops attribution and keeps the custom tags and the dashboard.",
			},
			"custom_tags": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "The names of the custom tags to attribute volume to, in addition to the standard tags. Each name must start with a letter or underscore and contain only letters, digits and underscores, must not repeat, and must not be a standard tag. The list replaces the custom tags currently configured, so an omitted or empty list removes them all. Custom tags are only changed while enabled is true: while enabled is false Sumo Logic keeps the stored tags and ignores this list. A tag applies only to data ingested after it is added.",
				Elem: &schema.Schema{
					Type: schema.TypeString,
					ValidateFunc: validation.StringMatch(
						regexp.MustCompile(usageAttributionTagNamePattern),
						"must start with a letter or underscore and contain only letters, digits and underscores",
					),
				},
			},
			"standard_tags": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "The names of the standard tags, which are always available and cannot be changed.",
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
			"breakdown_by_options": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "Everything the usage attribution dashboard can be broken down by: the standard tags, the custom tags and any columns of an ownership lookup.",
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
			"dashboard_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The ID of the Usage Attribution dashboard. Empty until the dashboard has been created, which happens in the Sumo Logic UI.",
			},
			"enabled_at": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "When usage attribution was first turned on, in ISO 8601 format (UTC). Empty while it is off.",
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

	update := resourceToUsageAttributionConfigUpdate(d)
	config, err := c.UpdateUsageAttributionConfig(update)
	if err != nil {
		return err
	}

	// The update has been applied by now. Keep the state in step with Sumo Logic first, so that the error below does
	// not leave the state behind it.
	if err := setUsageAttributionConfigResource(d, config); err != nil {
		return err
	}
	return checkCustomTagsKept(update.Enabled, update.CustomTags, config)
}

// checkCustomTagsKept fails when Sumo Logic kept fewer custom tags than were sent. It drops a custom tag that is also
// a standard tag. The standard tags can change, so the answer comes from Sumo Logic and no list is kept here. It is
// only meaningful when usage attribution is on: turning it off does not change the stored custom tags.
func checkCustomTagsKept(enabled bool, sent []UsageAttributionTag, config *UsageAttributionConfig) error {
	if !enabled {
		return nil
	}
	kept := make(map[string]bool, len(config.CustomTags))
	for _, tag := range config.CustomTags {
		kept[tag.Name] = true
	}
	var ignored []string
	for _, tag := range sent {
		if !kept[tag.Name] {
			ignored = append(ignored, tag.Name)
		}
	}
	if len(ignored) == 0 {
		return nil
	}
	return fmt.Errorf(
		"Sumo Logic ignored these custom_tags because they are standard tags: %s. Remove them from custom_tags. The standard tags are: %s",
		strings.Join(ignored, ", "),
		strings.Join(usageAttributionTagNames(config.StandardTags), ", "),
	)
}

// repeatedNames returns the names that appear more than once, in the order they first repeat.
func repeatedNames(names []interface{}) []string {
	seen := make(map[string]bool, len(names))
	var repeated []string
	for _, name := range names {
		n, _ := name.(string)
		if seen[n] {
			repeated = append(repeated, n)
		}
		seen[n] = true
	}
	return repeated
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
