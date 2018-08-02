package sumologic

import (
	"os"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/terraform"
)

func Provider() terraform.ResourceProvider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			"access_id": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  os.Getenv("SUMOLOGIC_ACCESSID"),
			},
			"access_key": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  os.Getenv("SUMOLOGIC_ACCESSKEY"),
			},
			"environment": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "us2",
			},
		},
		ResourcesMap: map[string]*schema.Resource{
			"sumologic_collector":                   resourceSumologicCollector(),
			"sumologic_http_source":                 resourceSumologicHTTPSource(),
			"sumologic_polling_source":              resourceSumologicPollingSource(),
			"sumologic_cloudsyslog_source":          resourceSumologicCloudsyslogSource(),
			"sumologic_http_source_sns_autoconfirm": resourceSumologicHTTPSourceSNSAutoConfirm(),
		},
		DataSourcesMap: map[string]*schema.Resource{
			"sumologic_caller_identity": dataSourceSumologicCallerIdentity(),
		},
		ConfigureFunc: providerConfigure,
	}
}

func providerConfigure(d *schema.ResourceData) (interface{}, error) {
	return NewClient(
		d.Get("access_id").(string),
		d.Get("access_key").(string),
		d.Get("environment").(string),
	)
}
