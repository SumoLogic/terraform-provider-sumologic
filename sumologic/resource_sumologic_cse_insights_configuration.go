package sumologic

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"log"
)

func resourceSumologicCSEInsightsConfiguration() *schema.Resource {
	return &schema.Resource{
		Create: resourceSumologicCSEInsightsConfigurationCreate,
		Read:   resourceSumologicCSEInsightsConfigurationRead,
		Delete: resourceSumologicCSEInsightsConfigurationDelete,
		Update: resourceSumologicCSEInsightsConfigurationUpdate,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"lookback_days": {
				Type:     schema.TypeFloat,
				Optional: true,
				ForceNew: false,
			},
			"threshold": {
				Type:     schema.TypeFloat,
				Optional: true,
				ForceNew: false,
			},
			"global_signal_suppression_window": {
				Type:     schema.TypeFloat,
				Optional: true,
				ForceNew: false,
			},
		},
	}
}

func resourceSumologicCSEInsightsConfigurationRead(d *schema.ResourceData, meta interface{}) error {
	c := meta.(*Client)

	var CSEInsightsConfiguration *CSEInsightsConfiguration

	CSEInsightsConfiguration, err := c.GetCSEInsightsConfiguration()
	if err != nil {
		log.Printf("[WARN] CSE Insights Configuration not found, err: %v", err)

	}

	if CSEInsightsConfiguration == nil {
		log.Printf("[WARN] CSE Insights Configuration not found, removing from state: %v", err)
		d.SetId("")
		return nil
	}

	d.Set("lookback_days", CSEInsightsConfiguration.LookbackDays)
	d.Set("threshold", CSEInsightsConfiguration.Threshold)
	d.Set("global_signal_suppression_window", CSEInsightsConfiguration.GlobalSignalSuppressionWindow)

	return nil
}

func resourceSumologicCSEInsightsConfigurationDelete(d *schema.ResourceData, meta interface{}) error {
	CSEInsightsConfiguration := CSEInsightsConfiguration{}

	c := meta.(*Client)
	err := c.UpdateCSEInsightsConfiguration(CSEInsightsConfiguration)
	if err != nil {
		return err
	}

	return resourceSumologicCSEInsightsConfigurationRead(d, meta)
}

func resourceSumologicCSEInsightsConfigurationCreate(d *schema.ResourceData, meta interface{}) error {
	//we are not really creating new object in backend, using constant id for terraform resource
	d.SetId("cse-insights-configuration")
	return resourceSumologicCSEInsightsConfigurationUpdate(d, meta)
}

func resourceSumologicCSEInsightsConfigurationUpdate(d *schema.ResourceData, meta interface{}) error {
	CSEInsightsConfiguration, err := resourceToCSEInsightsConfiguration(d)
	if err != nil {
		return err
	}

	c := meta.(*Client)
	if err = c.UpdateCSEInsightsConfiguration(CSEInsightsConfiguration); err != nil {
		return err
	}

	return resourceSumologicCSEInsightsConfigurationRead(d, meta)
}

func resourceToCSEInsightsConfiguration(d *schema.ResourceData) (CSEInsightsConfiguration, error) {
	id := d.Id()
	if id == "" {
		return CSEInsightsConfiguration{}, nil
	}

	lookbackDays := d.Get("lookback_days").(float64)
	threshold := d.Get("threshold").(float64)
	globalSignalSuppressionWindow := d.Get("global_signal_suppression_window").(float64)

	return CSEInsightsConfiguration{
		LookbackDays:                  &lookbackDays,
		Threshold:                     &threshold,
		GlobalSignalSuppressionWindow: &globalSignalSuppressionWindow,
	}, nil
}
