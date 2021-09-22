package sumologic

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"hash/fnv"
	"log"
	"strconv"
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
	d.SetId(hash("cse-insights-configuration"))
	return resourceSumologicCSEInsightsConfigurationUpdate(d, meta)
}
func hash(s string) string {
	h := fnv.New32a()
	h.Write([]byte(s))
	return strconv.Itoa(int(h.Sum32()))
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

	return CSEInsightsConfiguration{

		LookbackDays: d.Get("lookback_days").(float64),
		Threshold:    d.Get("threshold").(float64),
	}, nil
}
