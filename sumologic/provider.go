package sumologic

import (
	"fmt"
	"os"

	"github.com/go-errors/errors"

	"github.com/hashicorp/terraform/helper/mutexkv"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/terraform"
)

func Provider() terraform.ResourceProvider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			"access_id": {
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("SUMOLOGIC_ACCESSID", nil),
			},
			"access_key": {
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("SUMOLOGIC_ACCESSKEY", nil),
			},
			"environment": {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("SUMOLOGIC_ENVIRONMENT", nil),
			},
			"base_url": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  os.Getenv("SUMOLOGIC_BASE_URL"),
			},
		},
		ResourcesMap: map[string]*schema.Resource{
			"sumologic_collector":                          resourceSumologicCollector(),
			"sumologic_http_source":                        resourceSumologicHTTPSource(),
			"sumologic_polling_source":                     resourceSumologicPollingSource(),
			"sumologic_cloudsyslog_source":                 resourceSumologicCloudsyslogSource(),
			"sumologic_role":                               resourceSumologicRole(),
			"sumologic_user":                               resourceSumologicUser(),
			"sumologic_ingest_budget":                      resourceSumologicIngestBudget(),
			"sumologic_collector_ingest_budget_assignment": resourceSumologicCollectorIngestBudgetAssignment(),
			"sumologic_folder":                             resourceSumologicFolder(),
			"sumologic_scheduled_view":                     resourceSumologicScheduledView(),
			"sumologic_partition":                          resourceSumologicPartition(),
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
	baseUrl := d.Get("base_url").(string)

	msg := ""
	if accessId == "" {
		msg = "sumologic provider: access_id should be set;"
	}
	if accessKey == "" {
		msg = fmt.Sprintf("%s access_key should be set; ", msg)
	}
	if msg != "" {
		if environment == "" {
			msg = fmt.Sprintf("%s make sure environment is set", msg)
		}
		return nil, errors.New(msg)
	}

	return NewClient(
		accessId,
		accessKey,
		environment,
		baseUrl,
	)
}
