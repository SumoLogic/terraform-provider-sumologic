package sumologic

import (
	"errors"
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

func resourceSumologicGenericPollingSource() *schema.Resource {
	pollingSource := resourceSumologicSource()
	pollingSource.Create = resourceSumologicGenericPollingSourceCreate
	pollingSource.Read = resourceSumologicGenericPollingSourceRead
	pollingSource.Update = resourceSumologicGenericPollingSourceUpdate
	pollingSource.Importer = &schema.ResourceImporter{
		State: resourceSumologicSourceImport,
	}

	pollingSource.Schema["content_type"] = &schema.Schema{
		Type:     schema.TypeString,
		Required: true,
		ForceNew: true,
		ValidateFunc: validation.StringInSlice([]string{"AwsS3Bucket", "AwsElbBucket", "AwsCloudFrontBucket",
			"AwsCloudTrailBucket", "AwsS3AuditBucket", "AwsCloudWatch", "AwsInventory", "AwsXRay", "GcpMetrics", "AwsS3ArchiveBucket", "AzureEventHubLog", "AzureMetrics"}, false),
	}
	pollingSource.Schema["scan_interval"] = &schema.Schema{
		Type:     schema.TypeInt,
		Optional: true,
		Default:  300000,
		DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
			contentType := d.Get("content_type").(string)
			return contentType == "AzureEventHubLog"
		},
	}
	pollingSource.Schema["paused"] = &schema.Schema{
		Type:     schema.TypeBool,
		Optional: true,
		Default:  false,
	}
	pollingSource.Schema["url"] = &schema.Schema{
		Type:     schema.TypeString,
		Computed: true,
	}
	pollingSource.Schema["authentication"] = &schema.Schema{
		Type:     schema.TypeList,
		Required: true,
		ForceNew: true,
		MinItems: 1,
		MaxItems: 1,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"type": {
					Type:         schema.TypeString,
					Required:     true,
					ValidateFunc: validation.StringInSlice([]string{"S3BucketAuthentication", "AWSRoleBasedAuthentication", "service_account", "AzureEventHubAuthentication", "AzureClientSecretAuthentication"}, false),
				},
				"access_key": {
					Type:     schema.TypeString,
					Optional: true,
				},
				"secret_key": {
					Type:     schema.TypeString,
					Optional: true,
				},
				"role_arn": {
					Type:     schema.TypeString,
					Optional: true,
				},
				"region": {
					Type:     schema.TypeString,
					Optional: true,
				},
				"project_id": {
					Type:     schema.TypeString,
					Optional: true,
				},
				"private_key_id": {
					Type:     schema.TypeString,
					Optional: true,
				},
				"private_key": {
					Type:     schema.TypeString,
					Optional: true,
				},
				"client_email": {
					Type:     schema.TypeString,
					Optional: true,
				},
				"client_id": {
					Type:     schema.TypeString,
					Optional: true,
				},
				"auth_uri": {
					Type:     schema.TypeString,
					Optional: true,
				},
				"token_uri": {
					Type:     schema.TypeString,
					Optional: true,
				},
				"auth_provider_x509_cert_url": {
					Type:     schema.TypeString,
					Optional: true,
				},
				"client_x509_cert_url": {
					Type:     schema.TypeString,
					Optional: true,
				},
				"shared_access_policy_name": {
					Type:     schema.TypeString,
					Optional: true,
				},
				"shared_access_policy_key": {
					Type:     schema.TypeString,
					Optional: true,
				},
				"tenant_id": {
					Type:     schema.TypeString,
					Optional: true,
				},
				"client_secret": {
					Type:     schema.TypeString,
					Optional: true,
				},
			},
		},
	}
	pollingSource.Schema["path"] = &schema.Schema{
		Type:     schema.TypeList,
		Required: true,
		ForceNew: true,
		MinItems: 1,
		MaxItems: 1,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"type": {
					Type:     schema.TypeString,
					Required: true,
					ValidateFunc: validation.StringInSlice([]string{"S3BucketPathExpression", "CloudWatchPath",
						"AwsInventoryPath", "AwsXRayPath", "GcpMetricsPath", "AzureEventHubPath", "AzureMetricsPath"}, false),
				},
				"bucket_name": {
					Type:     schema.TypeString,
					Optional: true,
				},
				"use_versioned_api": {
					Type:     schema.TypeBool,
					Optional: true,
					Default:  nil,
				},
				"path_expression": {
					Type:     schema.TypeString,
					Optional: true,
				},
				"limit_to_regions": {
					Type:     schema.TypeList,
					Optional: true,
					Elem: &schema.Schema{
						Type: schema.TypeString,
					},
				},
				"limit_to_namespaces": {
					Type:     schema.TypeList,
					Optional: true,
					Elem: &schema.Schema{
						Type: schema.TypeString,
					},
				},
				"limit_to_services": {
					Type:     schema.TypeList,
					Optional: true,
					Elem: &schema.Schema{
						Type: schema.TypeString,
					},
				},
				"custom_services": {
					Type:     schema.TypeList,
					Optional: true,
					Elem: &schema.Resource{
						Schema: map[string]*schema.Schema{
							"service_name": {
								Type:     schema.TypeString,
								Optional: true,
							},
							"prefixes": {
								Type:     schema.TypeList,
								Optional: true,
								Elem: &schema.Schema{
									Type: schema.TypeString,
								},
							},
						},
					},
				},
				"tag_filters": {
					Type:     schema.TypeList,
					Optional: true,
					Elem: &schema.Resource{
						Schema: map[string]*schema.Schema{
							"type": {
								Type:     schema.TypeString,
								Optional: true,
							},
							"namespace": {
								Type:     schema.TypeString,
								Optional: true,
							},
							"tags": {
								Type:     schema.TypeList,
								Optional: true,
								Elem: &schema.Schema{
									Type: schema.TypeString,
								},
							},
						},
					},
				},
				"azure_tag_filters": {
					Type:     schema.TypeList,
					Optional: true,
					Elem: &schema.Resource{
						Schema: map[string]*schema.Schema{
							"type": {
								Type:     schema.TypeString,
								Required: true,
							},
							"namespace": {
								Type:     schema.TypeString,
								Optional: true,
							},
							"tags": {
								Type:     schema.TypeList,
								Optional: true,
								Elem: &schema.Resource{
									Schema: map[string]*schema.Schema{
										"name": {
											Type:     schema.TypeString,
											Required: true,
										},
										"values": {
											Type:     schema.TypeList,
											Optional: true,
											Elem: &schema.Schema{
												Type: schema.TypeString,
											},
										},
									},
								},
							},
						},
					},
				},
				"sns_topic_or_subscription_arn": {
					Type:     schema.TypeList,
					Computed: true,
					Elem: &schema.Resource{
						Schema: map[string]*schema.Schema{
							"is_success": {
								Type:     schema.TypeBool,
								Computed: true,
							},
							"arn": {
								Type:     schema.TypeString,
								Computed: true,
							},
						},
					},
				},
				"namespace": {
					Type:     schema.TypeString,
					Optional: true,
				},
				"event_hub_name": {
					Type:     schema.TypeString,
					Optional: true,
				},
				"consumer_group": {
					Type:     schema.TypeString,
					Optional: true,
				},
				"region": {
					Type:     schema.TypeString,
					Optional: true,
				},
				"environment": {
					Type:     schema.TypeString,
					Optional: true,
				},
			},
		},
	}
	return pollingSource
}

func resourceSumologicGenericPollingSourceCreate(d *schema.ResourceData, meta interface{}) error {
	c := meta.(*Client)

	if d.Id() == "" {
		source, err := resourceToGenericPollingSource(d)
		if err != nil {
			return err
		}

		sourceID, err := c.CreatePollingSource(source, d.Get("collector_id").(int))
		if err != nil {
			return err
		}

		d.SetId(strconv.Itoa(sourceID))
	}

	return resourceSumologicGenericPollingSourceRead(d, meta)
}

func resourceSumologicGenericPollingSourceUpdate(d *schema.ResourceData, meta interface{}) error {
	c := meta.(*Client)

	source, err := resourceToGenericPollingSource(d)
	if err != nil {
		return err
	}

	err = c.UpdatePollingSource(source, d.Get("collector_id").(int))
	if err != nil {
		return err
	}

	return resourceSumologicGenericPollingSourceRead(d, meta)
}

func resourceSumologicGenericPollingSourceRead(d *schema.ResourceData, meta interface{}) error {
	c := meta.(*Client)

	id, _ := strconv.Atoi(d.Id())
	source, err := c.GetPollingSource(d.Get("collector_id").(int), id)

	if err != nil {
		return err
	}

	if source == nil {
		log.Printf("[WARN] Polling source not found, removing from state: %v - %v", id, err)
		d.SetId("")

		return nil
	}

	pollingResources := source.ThirdPartyRef.Resources
	path := getPollingThirdPartyPathAttributes(pollingResources)

	if err := d.Set("path", path); err != nil {
		return err
	}

	authSettings := getPollingThirdPartyAuthenticationAttributes(pollingResources)
	if err := d.Set("authentication", authSettings); err != nil {
		return err
	}

	if err := resourceSumologicSourceRead(d, source.Source); err != nil {
		return fmt.Errorf("%s", err)
	}
	d.Set("content_type", source.ContentType)
	d.Set("scan_interval", source.ScanInterval)
	d.Set("paused", source.Paused)
	d.Set("url", source.URL)

	return nil
}

func resourceToGenericPollingSource(d *schema.ResourceData) (PollingSource, error) {
	source := resourceToSource(d)
	source.Type = "Polling"

	pollingSource := PollingSource{
		Source:       source,
		Paused:       d.Get("paused").(bool),
		ScanInterval: d.Get("scan_interval").(int),
		ContentType:  d.Get("content_type").(string),
		URL:          d.Get("url").(string),
	}

	if source.ContentType == "AzureEventHubLog" {
		pollingSource.ScanInterval = -1
	}

	authSettings, errAuthSettings := getPollingAuthentication(d)
	if errAuthSettings != nil {
		return pollingSource, errAuthSettings
	}

	pathSettings, errPathSettings := getPollingPathSettings(d)
	if errPathSettings != nil {
		return pollingSource, errPathSettings
	}

	pollingResource := PollingResource{
		ServiceType:    d.Get("content_type").(string),
		Authentication: authSettings,
		Path:           pathSettings,
	}

	pollingSource.ThirdPartyRef.Resources = append(pollingSource.ThirdPartyRef.Resources, pollingResource)

	return pollingSource, nil
}

func getPollingThirdPartyPathAttributes(pollingResource []PollingResource) []map[string]interface{} {

	var s []map[string]interface{}

	for _, t := range pollingResource {
		mapping := map[string]interface{}{
			"type":                          t.Path.Type,
			"bucket_name":                   t.Path.BucketName,
			"path_expression":               t.Path.PathExpression,
			"limit_to_regions":              t.Path.LimitToRegions,
			"limit_to_namespaces":           t.Path.LimitToNamespaces,
			"limit_to_services":             t.Path.LimitToServices,
			"use_versioned_api":             t.Path.UseVersionedApi,
			"custom_services":               flattenCustomServices(t.Path.CustomServices),
			"tag_filters":                   flattenPollingTagFilters(t.Path.TagFilters),
			"azure_tag_filters":             flattenPollingAzureTagFilters(t.Path.TagFilters),
			"sns_topic_or_subscription_arn": flattenPollingSnsTopicOrSubscriptionArn(t.Path.SnsTopicOrSubscriptionArn),
			"namespace":                     t.Path.Namespace,
			"event_hub_name":                t.Path.EventHubName,
			"consumer_group":                t.Path.ConsumerGroup,
			"region":                        t.Path.Region,
		}

		s = append(s, mapping)
	}
	return s
}

func getPollingThirdPartyAuthenticationAttributes(pollingResource []PollingResource) []map[string]interface{} {

	var s []map[string]interface{}

	for _, t := range pollingResource {
		mapping := map[string]interface{}{
			"type":                        t.Authentication.Type,
			"access_key":                  t.Authentication.AwsID,
			"secret_key":                  t.Authentication.AwsKey,
			"role_arn":                    t.Authentication.RoleARN,
			"region":                      t.Authentication.Region,
			"project_id":                  t.Authentication.ProjectId,
			"private_key_id":              t.Authentication.PrivateKeyId,
			"private_key":                 t.Authentication.PrivateKey,
			"client_email":                t.Authentication.ClientEmail,
			"client_id":                   t.Authentication.ClientId,
			"auth_uri":                    t.Authentication.AuthUrl,
			"token_uri":                   t.Authentication.TokenUrl,
			"auth_provider_x509_cert_url": t.Authentication.AuthProviderX509CertUrl,
			"client_x509_cert_url":        t.Authentication.ClientX509CertUrl,
			"shared_access_policy_name":   t.Authentication.SharedAccessPolicyName,
			"shared_access_policy_key":    t.Authentication.SharedAccessPolicyKey,
			"tenant_id":                   t.Authentication.TenantId,
			"client_secret":               t.Authentication.ClientSecret,
		}
		s = append(s, mapping)
	}
	return s
}

func flattenCustomServices(v []string) []map[string]interface{} {
	var custom_services []map[string]interface{}
	for _, d := range v {
		custom_service_name_and_prefixes := strings.Split(d, "=")
		custom_service_name := custom_service_name_and_prefixes[0]
		custom_service_prefixes := strings.Split(custom_service_name_and_prefixes[1], ";")
		custom_service := map[string]interface{}{
			"service_name": custom_service_name,
			"prefixes":     custom_service_prefixes,
		}
		custom_services = append(custom_services, custom_service)
	}
	return custom_services
}

func getCustomServices(path map[string]interface{}) []string {
	var customServices []string
	rawCustomServicesConfig := path["custom_services"].([]interface{})
	for _, rawCustomServiceConfig := range rawCustomServicesConfig {
		customServiceConfig := rawCustomServiceConfig.(map[string]interface{})
		customServiceName := customServiceConfig["service_name"].(string)
		customServicePrefixesInterface := customServiceConfig["prefixes"].([]interface{})
		var customServicePrefixes []string
		for _, v := range customServicePrefixesInterface {
			customServicePrefixes = append(customServicePrefixes, v.(string))
		}
		customServices = append(customServices,
			fmt.Sprintf("%s=%s", customServiceName, strings.Join(customServicePrefixes[:], ";")))
	}
	return customServices
}

func flattenPollingTagFilters(v []interface{}) []map[string]interface{} {
	var filters []map[string]interface{}
	if len(v) > 0 {
		filter := v[0].(map[string]interface{})
		switch filterType, ok := filter["type"].(string); {
		case !ok || filterType == "AzureTagFilters": // do nothing
		default:
			for _, d := range v {
				rawFilter := d.(map[string]interface{})
				filter := map[string]interface{}{
					"type":      rawFilter["type"].(string),
					"namespace": rawFilter["namespace"].(string),
					"tags":      rawFilter["tags"].([]interface{}),
				}
				filters = append(filters, filter)
			}
		}
	}
	return filters
}

func flattenPollingAzureTagFilters(v []interface{}) []map[string]interface{} {
	var filters []map[string]interface{}
	if len(v) > 0 {
		filter := v[0].(map[string]interface{})
		switch filterType, ok := filter["type"].(string); {
		case !ok || filterType == "AzureTagFilters":
			for _, d := range v {
				rawFilter := d.(map[string]interface{})
				filter := map[string]interface{}{
					"type":      "AzureTagFilters",
					"namespace": rawFilter["namespace"].(string),
					"tags":      flattenAzureTagKeyValuePair(rawFilter["tags"].([]interface{})),
				}
				filters = append(filters, filter)
			}
		default: // do nothing
		}
	}
	return filters
}

func flattenAzureTagKeyValuePair(v []interface{}) []map[string]interface{} {
	var tags []map[string]interface{}
	for _, d := range v {
		rawTag := d.(map[string]interface{})
		tag := map[string]interface{}{
			"name":   rawTag["name"].(string),
			"values": rawTag["values"].([]interface{}),
		}
		tags = append(tags, tag)
	}
	return tags
}

func getPollingTagFilters(d *schema.ResourceData) []interface{} {
	paths := d.Get("path").([]interface{})
	path := paths[0].(map[string]interface{})
	rawTagFilterConfig := path["tag_filters"].([]interface{})
	var filters []interface{}

	for _, rawConfig := range rawTagFilterConfig {
		config := rawConfig.(map[string]interface{})
		filter := TagFilter{}
		filter.Type = config["type"].(string)
		filter.Namespace = config["namespace"].(string)

		rawTags := config["tags"].([]interface{})
		Tags := make([]string, len(rawTags))
		for i, v := range rawTags {
			Tags[i] = v.(string)
		}
		filter.Tags = Tags
		filters = append(filters, filter)
	}
	return filters
}

func getPollingAzureTagFilters(d *schema.ResourceData) []interface{} {
	paths := d.Get("path").([]interface{})
	path := paths[0].(map[string]interface{})
	rawTagFilterConfig := path["azure_tag_filters"].([]interface{})
	var filters []interface{}

	for _, rawConfig := range rawTagFilterConfig {
		config := rawConfig.(map[string]interface{})
		filter := AzureTagFilter{}
		filter.Type = config["type"].(string)
		filter.Namespace = config["namespace"].(string)

		rawTags := config["tags"].([]interface{})
		Tags := make([]AzureTagKeyValuePair, len(rawTags))
		for i, rawTagPair := range rawTags {
			tagPair := rawTagPair.(map[string]interface{})
			pair := AzureTagKeyValuePair{}
			pair.Name = tagPair["name"].(string)
			// get the values from the azure key value pair
			rawValues := tagPair["values"].([]interface{})
			valueList := make([]string, len(rawValues))
			for j, v := range rawValues {
				valueList[j] = v.(string)
			}
			pair.Values = valueList
			Tags[i] = pair
		}
		filter.Tags = Tags
		filters = append(filters, filter)
	}
	return filters
}

func flattenPollingSnsTopicOrSubscriptionArn(v PollingSnsTopicOrSubscriptionArn) []map[string]interface{} {
	var snsTopicOrSubscriptionArn []map[string]interface{}
	snsTopic := map[string]interface{}{
		"is_success": v.IsSuccess,
		"arn":        v.Arn,
	}
	snsTopicOrSubscriptionArn = append(snsTopicOrSubscriptionArn, snsTopic)
	return snsTopicOrSubscriptionArn
}

func getPollingSnsTopicOrSubscriptionArn(d *schema.ResourceData) PollingSnsTopicOrSubscriptionArn {
	paths := d.Get("path").([]interface{})
	path := paths[0].(map[string]interface{})
	snsConfig := path["sns_topic_or_subscription_arn"].([]interface{})
	snsTopicOrSubscriptionArn := PollingSnsTopicOrSubscriptionArn{}

	if len(snsConfig) > 0 {
		for _, rawConfig := range snsConfig {
			config := rawConfig.(map[string]interface{})
			snsTopicOrSubscriptionArn.IsSuccess = config["is_success"].(bool)
			snsTopicOrSubscriptionArn.Arn = config["arn"].(string)
		}
	}
	return snsTopicOrSubscriptionArn
}

func addGcpServiceAccountDetailsToAuth(authSettings *PollingAuthentication, auth map[string]interface{}) error {
	authSettings.Type = "service_account"
	authSettings.ProjectId = auth["project_id"].(string)
	authSettings.PrivateKeyId = auth["private_key_id"].(string)
	authSettings.PrivateKey = auth["private_key"].(string)
	authSettings.ClientEmail = auth["client_email"].(string)
	authSettings.ClientId = auth["client_id"].(string)
	authSettings.AuthUrl = auth["auth_uri"].(string)
	authSettings.TokenUrl = auth["token_uri"].(string)
	authSettings.AuthProviderX509CertUrl = auth["auth_provider_x509_cert_url"].(string)
	authSettings.ClientX509CertUrl = auth["client_x509_cert_url"].(string)

	errTxt := ""
	if len(strings.Trim(authSettings.ProjectId, " \t")) == 0 {
		errTxt = errTxt + "\nproject_id is mandatory while using service_account authentication"
	}
	if len(authSettings.ClientEmail) == 0 {
		errTxt = errTxt + "\nclient_email is mandatory while using service_account authentication"
	}
	if len(authSettings.PrivateKey) == 0 {
		errTxt = errTxt + "\nprivate_key is mandatory while using service_account authentication"
	}

	if len(errTxt) == 0 {
		return nil
	} else {
		return errors.New(errTxt)
	}
}

func getPollingAuthentication(d *schema.ResourceData) (PollingAuthentication, error) {
	auths := d.Get("authentication").([]interface{})
	authSettings := PollingAuthentication{}

	if len(auths) > 0 {
		auth := auths[0].(map[string]interface{})
		switch authType := auth["type"].(string); authType {
		case "S3BucketAuthentication":
			if d.Get("content_type").(string) == "AwsInventory" {
				return authSettings, errors.New(
					fmt.Sprintf("[ERROR] Unsupported authType: %v for AwsInventory source", authType))
			}
			authSettings.Type = "S3BucketAuthentication"
			authSettings.AwsID = auth["access_key"].(string)
			authSettings.AwsKey = auth["secret_key"].(string)
			if auth["region"] != nil {
				authSettings.Region = auth["region"].(string)
			}
		case "AWSRoleBasedAuthentication":
			authSettings.Type = "AWSRoleBasedAuthentication"
			authSettings.RoleARN = auth["role_arn"].(string)
			if auth["region"] != nil {
				authSettings.Region = auth["region"].(string)
			}
		case "service_account":
			err := addGcpServiceAccountDetailsToAuth(&authSettings, auth)
			if err != nil {
				return authSettings, err
			}
		case "AzureEventHubAuthentication":
			authSettings.Type = "AzureEventHubAuthentication"
			authSettings.SharedAccessPolicyName = auth["shared_access_policy_name"].(string)
			authSettings.SharedAccessPolicyKey = auth["shared_access_policy_key"].(string)
		case "AzureClientSecretAuthentication":
			authSettings.Type = "AzureClientSecretAuthentication"
			authSettings.TenantId = auth["tenant_id"].(string)
			authSettings.AzureClientId = auth["client_id"].(string)
			authSettings.ClientSecret = auth["client_secret"].(string)
		default:
			errorMessage := fmt.Sprintf("[ERROR] Unknown authType: %v", authType)
			log.Print(errorMessage)
			return authSettings, errors.New(errorMessage)
		}
	}

	return authSettings, nil
}

func getLimitToRegions(path map[string]interface{}) []string {
	rawLimitToRegions := path["limit_to_regions"].([]interface{})
	limitToRegions := make([]string, len(rawLimitToRegions))
	for i, v := range rawLimitToRegions {
		limitToRegions[i] = v.(string)
	}
	return limitToRegions
}

func getLimitToServices(path map[string]interface{}) []string {
	rawLimitToServices := path["limit_to_services"].([]interface{})
	limitToServices := make([]string, len(rawLimitToServices))
	for i, v := range rawLimitToServices {
		limitToServices[i] = v.(string)
	}
	return limitToServices
}

func addGcpMetricsPathSettings(pathSettings *PollingPath, path map[string]interface{}) {
	pathSettings.LimitToRegions = getLimitToRegions(path)
	pathSettings.LimitToServices = getLimitToServices(path)
	pathSettings.CustomServices = getCustomServices(path)
}

func getPollingPathSettings(d *schema.ResourceData) (PollingPath, error) {
	pathSettings := PollingPath{}
	paths := d.Get("path").([]interface{})

	if len(paths) > 0 {
		path := paths[0].(map[string]interface{})
		switch pathType := path["type"].(string); pathType {
		case "S3BucketPathExpression":
			pathSettings.Type = "S3BucketPathExpression"
			pathSettings.BucketName = path["bucket_name"].(string)
			pathSettings.PathExpression = path["path_expression"].(string)

			if isFieldSet(d, "path.0.use_versioned_api") {
				val := path["use_versioned_api"].(bool)
				pathSettings.UseVersionedApi = &val
			} else {
				pathSettings.UseVersionedApi = nil
			}

			pathSettings.SnsTopicOrSubscriptionArn = getPollingSnsTopicOrSubscriptionArn(d)
		case "CloudWatchPath", "AwsInventoryPath":
			pathSettings.Type = pathType
			rawLimitToRegions := path["limit_to_regions"].([]interface{})
			LimitToRegions := make([]string, 0, len(rawLimitToRegions))
			for _, v := range rawLimitToRegions {
				if v != nil {
					LimitToRegions = append(LimitToRegions, v.(string))
				}
			}

			rawLimitToNamespaces := path["limit_to_namespaces"].([]interface{})
			LimitToNamespaces := make([]string, 0, len(rawLimitToNamespaces))
			for _, v := range rawLimitToNamespaces {
				if v != nil {
					LimitToNamespaces = append(LimitToNamespaces, v.(string))
				}
			}
			pathSettings.LimitToRegions = LimitToRegions
			pathSettings.LimitToNamespaces = LimitToNamespaces
			if pathType == "CloudWatchPath" {
				pathSettings.TagFilters = getPollingTagFilters(d)
			}
		case "AwsXRayPath":
			pathSettings.Type = "AwsXRayPath"
			rawLimitToRegions := path["limit_to_regions"].([]interface{})
			LimitToRegions := make([]string, 0, len(rawLimitToRegions))
			for _, v := range rawLimitToRegions {
				if v != nil {
					LimitToRegions = append(LimitToRegions, v.(string))
				}
			}
			pathSettings.LimitToRegions = LimitToRegions
		case "GcpMetricsPath":
			pathSettings.Type = pathType
			addGcpMetricsPathSettings(&pathSettings, path)
		case "AzureEventHubPath":
			pathSettings.Type = "AzureEventHubPath"
			pathSettings.Namespace = path["namespace"].(string)
			pathSettings.EventHubName = path["event_hub_name"].(string)
			pathSettings.ConsumerGroup = path["consumer_group"].(string)
			if path["region"] != nil {
				pathSettings.Region = path["region"].(string)
			}
		case "AzureMetricsPath":
			pathSettings.Type = "AzureMetricsPath"
			pathSettings.Environment = path["environment"].(string)
			rawLimitToNamespaces := path["limit_to_namespaces"].([]interface{})
			LimitToNamespaces := make([]string, 0, len(rawLimitToNamespaces))
			for _, v := range rawLimitToNamespaces {
				if v != nil {
					LimitToNamespaces = append(LimitToNamespaces, v.(string))
				}
			}
			pathSettings.LimitToNamespaces = LimitToNamespaces
			pathSettings.TagFilters = getPollingAzureTagFilters(d)
		default:
			errorMessage := fmt.Sprintf("[ERROR] Unknown resourceType in path: %v", pathType)
			log.Print(errorMessage)
			return pathSettings, errors.New(errorMessage)
		}
	}

	return pathSettings, nil
}

func isFieldSet(d *schema.ResourceData, key string) bool {
	_, ok := d.GetOkExists(key)
	return ok
}
