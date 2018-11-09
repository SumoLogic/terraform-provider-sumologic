package sumologic

import (
	"log"
	"os"

	"github.com/hashicorp/terraform/helper/mutexkv"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/terraform"
)

func Provider() terraform.ResourceProvider {
	defaultEnvironment := os.Getenv("SUMOLOGIC_ENVIRONMENT")
	if defaultEnvironment == "" {
		defaultEnvironment = "us2"
	}
	log.Printf("[DEBUG] sumo default environment: %s", defaultEnvironment)

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
				Default:  defaultEnvironment,
			},
		},
		ResourcesMap: map[string]*schema.Resource{
			"sumologic_collector":          resourceSumologicCollector(),
			"sumologic_http_source":        resourceSumologicHTTPSource(),
			"sumologic_polling_source":     resourceSumologicPollingSource(),
			"sumologic_cloudsyslog_source": resourceSumologicCloudsyslogSource(),
		},
		DataSourcesMap: map[string]*schema.Resource{
			"sumologic_caller_identity": dataSourceSumologicCallerIdentity(),
			"sumologic_collector":       dataSourceSumologicCollector(),
		},
		ConfigureFunc: providerConfigure,
	}
}

var SumoMutexKV = mutexkv.NewMutexKV()

func providerConfigure(d *schema.ResourceData) (interface{}, error) {
	return NewClient(
		d.Get("access_id").(string),
		d.Get("access_key").(string),
		d.Get("environment").(string),
	)
}
