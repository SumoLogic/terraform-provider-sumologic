package sumologic

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/go-errors/errors"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func Provider() *schema.Provider {
	log.Printf("Sumo Logic Terraform Provider Version=%s\n", ProviderVersion)
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
			"admin_mode": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
		},
		ResourcesMap: map[string]*schema.Resource{
			"sumologic_app":                                      resourceSumologicApp(),
			"sumologic_cse_tag_schema":                           resourceSumologicCSETagSchema(),
			"sumologic_cse_context_action":                       resourceSumologicCSEContextAction(),
			"sumologic_cse_automation":                           resourceSumologicCSEAutomation(),
			"sumologic_cse_entity_normalization_configuration":   resourceSumologicCSEEntityNormalizationConfiguration(),
			"sumologic_cse_inventory_entity_group_configuration": resourceSumologicCSEInventoryEntityGroupConfiguration(),
			"sumologic_cse_entity_entity_group_configuration":    resourceSumologicCSEEntityEntityGroupConfiguration(),
			"sumologic_cse_match_list":                           resourceSumologicCSEMatchList(),
			"sumologic_cse_custom_match_list_column":             resourceSumologicCSECustomMatchListColumn(),
			"sumologic_cse_log_mapping":                          resourceSumologicCSELogMapping(),
			"sumologic_cse_rule_tuning_expression":               resourceSumologicCSERuleTuningExpression(),
			"sumologic_cse_network_block":                        resourceSumologicCSENetworkBlock(),
			"sumologic_cse_custom_entity_type":                   resourceSumologicCSECustomEntityType(),
			"sumologic_cse_custom_insight":                       resourceSumologicCSECustomInsight(),
			"sumologic_cse_entity_criticality_config":            resourceSumologicCSEEntityCriticalityConfig(),
			"sumologic_cse_insights_configuration":               resourceSumologicCSEInsightsConfiguration(),
			"sumologic_cse_insights_resolution":                  resourceSumologicCSEInsightsResolution(),
			"sumologic_cse_insights_status":                      resourceSumologicCSEInsightsStatus(),
			"sumologic_cse_aggregation_rule":                     resourceSumologicCSEAggregationRule(),
			"sumologic_cse_chain_rule":                           resourceSumologicCSEChainRule(),
			"sumologic_cse_match_rule":                           resourceSumologicCSEMatchRule(),
			"sumologic_cse_threshold_rule":                       resourceSumologicCSEThresholdRule(),
			"sumologic_cse_first_seen_rule":                      resourceSumologicCSEFirstSeenRule(),
			"sumologic_cse_outlier_rule":                         resourceSumologicCSEOutlierRule(),
			"sumologic_collector":                                resourceSumologicCollector(),
			"sumologic_installed_collector":                      resourceSumologicInstalledCollector(),
			"sumologic_http_source":                              resourceSumologicHTTPSource(),
			"sumologic_gcp_source":                               resourceSumologicGCPSource(),
			"sumologic_polling_source":                           resourceSumologicPollingSource(),
			"sumologic_s3_source":                                resourceSumologicGenericPollingSource(),
			"sumologic_s3_audit_source":                          resourceSumologicGenericPollingSource(),
			"sumologic_s3_archive_source":                        resourceSumologicGenericPollingSource(),
			"sumologic_cloudwatch_source":                        resourceSumologicGenericPollingSource(),
			"sumologic_aws_inventory_source":                     resourceSumologicGenericPollingSource(),
			"sumologic_aws_xray_source":                          resourceSumologicGenericPollingSource(),
			"sumologic_cloudtrail_source":                        resourceSumologicGenericPollingSource(),
			"sumologic_elb_source":                               resourceSumologicGenericPollingSource(),
			"sumologic_cloudfront_source":                        resourceSumologicGenericPollingSource(),
			"sumologic_gcp_metrics_source":                       resourceSumologicGenericPollingSource(),
			"sumologic_cloud_to_cloud_source":                    resourceSumologicCloudToCloudSource(),
			"sumologic_metadata_source":                          resourceSumologicMetadataSource(),
			"sumologic_cloudsyslog_source":                       resourceSumologicCloudsyslogSource(),
			"sumologic_role":                                     resourceSumologicRole(),
			"sumologic_user":                                     resourceSumologicUser(),
			"sumologic_folder":                                   resourceSumologicFolder(),
			"sumologic_content":                                  resourceSumologicContent(),
			"sumologic_scheduled_view":                           resourceSumologicScheduledView(),
			"sumologic_data_forwarding_destination":              resourceSumologicDataForwardingDestination(),
			"sumologic_data_forwarding_rule":                     resourceSumologicDataForwardingRule(),
			"sumologic_partition":                                resourceSumologicPartition(),
			"sumologic_field_extraction_rule":                    resourceSumologicFieldExtractionRule(),
			"sumologic_connection":                               resourceSumologicConnection(),
			"sumologic_monitor":                                  resourceSumologicMonitorsLibraryMonitor(),
			"sumologic_monitor_folder":                           resourceSumologicMonitorsLibraryFolder(),
			"sumologic_muting_schedule":                          resourceSumologicMutingSchedulesLibraryMutingSchedule(),
			"sumologic_slo":                                      resourceSumologicSLO(),
			"sumologic_slo_folder":                               resourceSumologicSLOLibraryFolder(),
			"sumologic_ingest_budget_v2":                         resourceSumologicIngestBudgetV2(),
			"sumologic_field":                                    resourceSumologicField(),
			"sumologic_lookup_table":                             resourceSumologicLookupTable(),
			"sumologic_subdomain":                                resourceSumologicSubdomain(),
			"sumologic_dashboard":                                resourceSumologicDashboard(),
			"sumologic_macro":                                    resourceSumologicMacro(),
			"sumologic_password_policy":                          resourceSumologicPasswordPolicy(),
			"sumologic_saml_configuration":                       resourceSumologicSamlConfiguration(),
			"sumologic_kinesis_metrics_source":                   resourceSumologicKinesisMetricsSource(),
			"sumologic_kinesis_log_source":                       resourceSumologicKinesisLogSource(),
			"sumologic_token":                                    resourceSumologicToken(),
			"sumologic_policies":                                 resourceSumologicPolicies(),
			"sumologic_hierarchy":                                resourceSumologicHierarchy(),
			"sumologic_content_permission":                       resourceSumologicPermissions(),
			"sumologic_local_file_source":                        resourceSumologicLocalFileSource(),
			"sumologic_log_search":                               resourceSumologicLogSearch(),
			"sumologic_metrics_search":                           resourceSumologicMetricsSearch(),
			"sumologic_metrics_search_v2":                        resourceSumologicMetricsSearchV2(),
			"sumologic_rum_source":                               resourceSumologicRumSource(),
			"sumologic_role_v2":                                  resourceSumologicRoleV2(),
			"sumologic_azure_event_hub_log_source":               resourceSumologicGenericPollingSource(),
			"sumologic_ot_collector":                             resourceSumologicOTCollector(),
			"sumologic_source_template":                          resourceSumologicSourceTemplate(),
			"sumologic_azure_metrics_source":                     resourceSumologicGenericPollingSource(),
			"sumologic_scan_budget":                              resourceSumologicScanBudget(),
			"sumologic_local_windows_event_log_source":           resourceSumologicLocalWindowsEventLogSource(),
		},
		DataSourcesMap: map[string]*schema.Resource{
			"sumologic_cse_log_mapping_vendor_product": dataSourceCSELogMappingVendorAndProduct(),
			"sumologic_admin_recommended_folder":       dataSourceSumologicAdminRecommendedFolder(),
			"sumologic_caller_identity":                dataSourceSumologicCallerIdentity(),
			"sumologic_collector":                      dataSourceSumologicCollector(),
			"sumologic_http_source":                    dataSourceSumologicHTTPSource(),
			"sumologic_personal_folder":                dataSourceSumologicPersonalFolder(),
			"sumologic_folder":                         dataSourceSumologicFolder(),
			"sumologic_monitor_folder":                 dataSourceSumologicMonitorFolder(),
			"sumologic_my_user_id":                     dataSourceSumologicMyUserId(),
			"sumologic_partition":                      dataSourceSumologicPartition(),
			"sumologic_partitions":                     dataSourceSumologicPartitions(),
			"sumologic_role":                           dataSourceSumologicRole(),
			"sumologic_role_v2":                        dataSourceSumologicRoleV2(),
			"sumologic_user":                           dataSourceSumologicUser(),
			"sumologic_apps":                           dataSourceSumoLogicApps(),
		},
		ConfigureFunc: providerConfigure,
	}
}

func resolveRedirectURL(accessId string, accessKey string, authJwt string) (string, error) {
	req, err := http.NewRequest(http.MethodHead, "https://api.sumologic.com/api/v1/collectors", nil)
	if err != nil {
		return "", err
	}
	if authJwt == "" {
		req.SetBasicAuth(accessId, accessKey)
	} else {
		req.Header.Add("Authorization", "Bearer "+authJwt)
	}
	client := &http.Client{CheckRedirect: func(req *http.Request, via []*http.Request) error {
		return http.ErrUseLastResponse
	}}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	location := resp.Header.Get("location")
	if location == "" {
		// location header not found implies there was no redirect needed
		// i.e, the desired environment is us1
		return "https://api.sumologic.com/api/", nil
	}
	return strings.Split(location, "v1")[0], nil
}

func providerConfigure(d *schema.ResourceData) (interface{}, error) {
	accessId := d.Get("access_id").(string)
	accessKey := d.Get("access_key").(string)
	authJwt := os.Getenv("SUMOLOGIC_AUTHJWT")
	environment := d.Get("environment").(string)
	baseUrl := d.Get("base_url").(string)
	isInAdminMode := d.Get("admin_mode").(bool)

	msg := ""
	if authJwt == "" {
		if accessId == "" || accessKey == "" {
			msg = "sumologic provider: "
		}
		if accessId == "" {
			msg = fmt.Sprintf("%s access_id should be set;", msg)
		}
		if accessKey == "" {
			msg = fmt.Sprintf("%s access_key should be set; ", msg)
		}
	}

	if environment == "" && baseUrl == "" {
		log.Printf("Attempting to resolve redirection URL from access key/id")
		url, err := resolveRedirectURL(accessId, accessKey, authJwt)
		if err != nil {
			log.Printf("[WARN] Unable to resolve redirection URL, %s", err)
			environment = "us2"
			// baseUrl will be set accordingly in NewClient constructor
			log.Printf("[WARN] environment not set, defaulting to %s", environment)
		} else {
			baseUrl = url
			log.Printf("Resolved redirection URL %s", baseUrl)
		}

	}

	if msg != "" {
		return nil, errors.New(msg)
	}

	return NewClient(
		accessId,
		accessKey,
		authJwt,
		environment,
		baseUrl,
		isInAdminMode,
	)
}
