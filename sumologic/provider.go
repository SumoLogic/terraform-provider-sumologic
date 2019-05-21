package sumologic

import (
	"fmt"
	"github.com/go-errors/errors"
	"log"
	"os"

	"github.com/hashicorp/terraform/helper/mutexkv"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/terraform"
)

const DefaultEnvironment = "us2"

func Provider() terraform.ResourceProvider {
	defaultEnvironment := os.Getenv("SUMOLOGIC_ENVIRONMENT")
	if defaultEnvironment == "" {
		defaultEnvironment = DefaultEnvironment
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
			"sumologic_role":               resourceSumologicRole(),
			"sumologic_user":               resourceSumologicUser(),
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
	accessId := d.Get("access_id").(string)
	accessKey := d.Get("access_key").(string)
	environment := d.Get("environment").(string)

	msg := ""
	if accessId == "" {
		msg = "sumologic provider: access_id should be set;"
	}
	if accessKey == "" {
		msg = fmt.Sprintf("%s access_key should be set; ", msg)
	}
	if msg != "" {
		if environment == DefaultEnvironment {
			msg = fmt.Sprintf("%s make sure environment is set or that the default (%s) is appropriate", msg, DefaultEnvironment)
		}
		return nil, errors.New(msg)
	}

	return NewClient(
		accessId,
		accessKey,
		environment,
	)
}
