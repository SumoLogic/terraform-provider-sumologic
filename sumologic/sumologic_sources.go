package sumologic

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

type Source struct {
	ID                         int                    `json:"id,omitempty"`
	Type                       string                 `json:"sourceType"`
	Name                       string                 `json:"name"`
	Description                string                 `json:"description,omitempty"`
	Category                   string                 `json:"category,omitempty"`
	HostName                   string                 `json:"hostName,omitempty"`
	TimeZone                   string                 `json:"timeZone,omitempty"`
	AutomaticDateParsing       bool                   `json:"automaticDateParsing"`
	MultilineProcessingEnabled bool                   `json:"multilineProcessingEnabled"`
	UseAutolineMatching        bool                   `json:"useAutolineMatching"`
	ManualPrefixRegexp         string                 `json:"manualPrefixRegexp,omitempty"`
	ForceTimeZone              bool                   `json:"forceTimeZone"`
	DefaultDateFormats         []DefaultDateFormat    `json:"defaultDateFormats,omitempty"`
	Filters                    []Filter               `json:"filters,omitempty"`
	CutoffTimestamp            int                    `json:"cutoffTimestamp,omitempty"`
	CutoffRelativeTime         string                 `json:"cutoffRelativeTime,omitempty"`
	Fields                     map[string]interface{} `json:"fields,omitempty"`
	Url                        string                 `json:"url,omitempty"`
	ContentType                string                 `json:"contentType,omitempty"`
}

type DefaultDateFormat struct {
	Format  string `json:"format"`
	Locator string `json:"locator"`
}

type Filter struct {
	Name       string `json:"name"`
	FilterType string `json:"filterType"`
	Regexp     string `json:"regexp"`
	Mask       string `json:"mask"`
}

type SourceList struct {
	Sources []Source `json:"sources"`
}

func resourceSumologicSource() *schema.Resource {
	return &schema.Resource{
		Delete: resourceSumologicSourceDelete,
		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "",
			},
			"category": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "",
			},
			"host_name": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "",
			},
			"timezone": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "Etc/UTC",
			},
			"automatic_date_parsing": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  true,
			},
			"multiline_processing_enabled": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  true,
			},
			"use_autoline_matching": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  true,
			},
			"manual_prefix_regexp": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  nil,
			},
			"force_timezone": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			"default_date_formats": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"format": {
							Type:     schema.TypeString,
							Required: true,
						},
						"locator": {
							Type:     schema.TypeString,
							Optional: true,
						},
					},
				},
			},
			"filters": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:     schema.TypeString,
							Required: true,
						},
						"filter_type": {
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: validation.StringInSlice([]string{"Exclude", "Include", "Hash", "Mask", "Forward"}, false),
						},
						"regexp": {
							Type:     schema.TypeString,
							Required: true,
						},
						"mask": {
							Type:     schema.TypeString,
							Optional: true,
						},
					},
				},
			},
			"cutoff_timestamp": {
				Type:     schema.TypeInt,
				Optional: true,
				Default:  0,
			},
			"cutoff_relative_time": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				Default:  nil,
			},
			"fields": {
				Type: schema.TypeMap,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Optional: true,
				Default:  nil,
			},
			"collector_id": {
				Type:     schema.TypeInt,
				Required: true,
				ForceNew: true,
			},
			"content_type": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  nil,
			},
		},
	}
}

func resourceSumologicSourceDelete(d *schema.ResourceData, meta interface{}) error {
	c := meta.(*Client)

	id, _ := strconv.Atoi(d.Id())
	collectorID, _ := d.Get("collector_id").(int)

	return c.DestroySource(id, collectorID)

}

func resourceSumologicSourceImport(d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
	ids := strings.Split(d.Id(), "/")
	c := m.(*Client)

	if len(ids) != 2 {
		return nil, fmt.Errorf("expected collector/source, got %s", d.Id())
	}

	collectorID, err := strconv.Atoi(ids[0])
	if err == nil {
		// collectorId/sourceId
		d.SetId(ids[1])
		d.Set("collector_id", collectorID)
	} else {
		// collectorName/sourceName
		collector, _ := c.GetCollectorName(ids[0])
		if collector != nil {
			source, _ := c.GetSourceName(collector.ID, ids[1])
			if source != nil {
				d.SetId(strconv.Itoa(source.ID))
				d.Set("collector_id", collector.ID)
			} else {
				return nil, fmt.Errorf("source with name '%s' does not exist", ids[1])
			}
		} else {
			return nil, fmt.Errorf("collector with name '%s' does not exist", ids[0])
		}
	}

	return []*schema.ResourceData{d}, nil
}

func resourceToSource(d *schema.ResourceData) Source {
	id, _ := strconv.Atoi(d.Id())

	source := Source{}
	source.ID = id
	source.Name = d.Get("name").(string)
	source.Description = d.Get("description").(string)
	source.Category = d.Get("category").(string)
	source.HostName = d.Get("host_name").(string)
	source.TimeZone = d.Get("timezone").(string)
	source.AutomaticDateParsing = d.Get("automatic_date_parsing").(bool)
	source.MultilineProcessingEnabled = d.Get("multiline_processing_enabled").(bool)
	source.UseAutolineMatching = d.Get("use_autoline_matching").(bool)
	source.ManualPrefixRegexp = d.Get("manual_prefix_regexp").(string)
	source.ForceTimeZone = d.Get("force_timezone").(bool)
	source.DefaultDateFormats = getDefaultDateFormats(d)
	source.Filters = getFilters(d)
	source.CutoffTimestamp = d.Get("cutoff_timestamp").(int)
	source.CutoffRelativeTime = d.Get("cutoff_relative_time").(string)
	source.Fields = d.Get("fields").(map[string]interface{})
	source.ContentType = d.Get("content_type").(string)

	return source
}

func resourceSumologicSourceRead(d *schema.ResourceData, source Source) error {
	d.Set("name", source.Name)
	d.Set("description", source.Description)
	d.Set("category", source.Category)
	d.Set("host_name", source.HostName)
	d.Set("timezone", source.TimeZone)
	d.Set("automatic_date_parsing", source.AutomaticDateParsing)
	d.Set("multiline_processing_enabled", source.MultilineProcessingEnabled)
	d.Set("use_autoline_matching", source.UseAutolineMatching)
	d.Set("manual_prefix_regexp", source.ManualPrefixRegexp)
	d.Set("force_timezone", source.ForceTimeZone)
	if err := d.Set("default_date_formats", flattenDateFormats(source.DefaultDateFormats)); err != nil {
		return fmt.Errorf("error setting default date formats for resource %s: %s", d.Id(), err)
	}
	if err := d.Set("filters", flattenFilters(source.Filters)); err != nil {
		return fmt.Errorf("error setting filters for resource %s: %s", d.Id(), err)
	}
	d.Set("cutoff_timestamp", source.CutoffTimestamp)
	d.Set("cutoff_relative_time", source.CutoffRelativeTime)
	if err := d.Set("fields", source.Fields); err != nil {
		return fmt.Errorf("error setting fields for resource %s: %s", d.Id(), err)
	}
	d.Set("content_type", source.ContentType)
	return nil
}

func flattenDateFormats(v []DefaultDateFormat) []map[string]interface{} {
	var defaultDateDormats []map[string]interface{}
	for _, d := range v {
		defaultDateFormat := map[string]interface{}{
			"format":  d.Format,
			"locator": d.Locator,
		}
		defaultDateDormats = append(defaultDateDormats, defaultDateFormat)
	}
	return defaultDateDormats
}

func flattenFilters(v []Filter) []map[string]interface{} {
	var filters []map[string]interface{}
	for _, d := range v {
		filter := map[string]interface{}{
			"name":        d.Name,
			"filter_type": d.FilterType,
			"regexp":      d.Regexp,
			"mask":        d.Mask,
		}
		filters = append(filters, filter)
	}
	return filters
}

func getDefaultDateFormats(d *schema.ResourceData) []DefaultDateFormat {

	rawDefaultDateFormatsConfig := d.Get("default_date_formats").([]interface{})
	var defaultDateDormats []DefaultDateFormat

	for _, rawConfig := range rawDefaultDateFormatsConfig {
		config := rawConfig.(map[string]interface{})
		defaultDateFormat := DefaultDateFormat{}
		defaultDateFormat.Format = config["format"].(string)
		defaultDateFormat.Locator = config["locator"].(string)
		defaultDateDormats = append(defaultDateDormats, defaultDateFormat)
	}

	return defaultDateDormats
}

func getFilters(d *schema.ResourceData) []Filter {

	rawFilterConfig := d.Get("filters").([]interface{})
	var filters []Filter

	for _, rawConfig := range rawFilterConfig {
		config := rawConfig.(map[string]interface{})
		filter := Filter{}
		filter.Name = config["name"].(string)
		filter.FilterType = config["filter_type"].(string)
		filter.Regexp = config["regexp"].(string)
		filter.Mask = config["mask"].(string)
		filters = append(filters, filter)
	}

	return filters
}

func (s *Client) DestroySource(sourceID int, collectorID int) error {

	_, err := s.Delete(fmt.Sprintf("v1/collectors/%d/sources/%d", collectorID, sourceID))

	return err
}

func (s *Client) GetSourceName(collectorID int64, sourceName string) (*Source, error) {

	data, _, err := s.Get(fmt.Sprintf("v1/collectors/%d/sources", collectorID), false)

	if err != nil {
		return nil, err
	}

	if data == nil {
		return nil, nil
	}

	var response SourceList
	err = json.Unmarshal(data, &response)

	if err != nil {
		return nil, err
	}

	for _, source := range response.Sources {
		if source.Name == sourceName {
			return &source, nil
		}
	}

	return nil, nil
}
