## 2.20.0 (Unreleased)
BUG FIXES:
* Fix typo on cse_match_list documentation (GH-460)

## 2.19.2 (November 4, 2022)
ENHANCEMENTS:
* Suppress diffs for equivalent values of some time attributes. This should reduce output of `terraform plan` that didn't disappear after running `terraform apply`. (GH-442)
* Add better validation and documentation for some time attributes (GH-443)

## 2.19.1 (October 6, 2022)
FEATURES:
* Add new optional `resolution_window` field to resource/sumologic_monitor (GH-418)

BUG FIXES:
* CSE rules hard failing if passing tags with empty strings. (GH-445)
* Return error when unable to read collectors. (GH-446)

## 2.19.0 (September 20, 2022)
FEATURES:
* **New Resource:** sumologic_cse_entity_normalization_configuration (GH-430)

ENHANCEMENTS:
* Updated maxdepth level for hierarchy resource (GH-433)

## 2.18.2 (September 1, 2022)
BUG FIXES:
* Fix bug for validation for monitor name and description regex (GH-428)

## 2.18.1 (August 31, 2022)

BUG FIXES:
* Fix compliance period validation for SLOs (GH-424)
* Adding validations for name, description and payload_override in monitors (GH-420)

## 2.18.0 (August 8, 2022)

FEATURES:
* **New Resource:** sumologic_installed_collector (GH-182)

ENHANCEMENTS:
* Upgraded dependencies (GH-393)

## 2.17.0 (July 28, 2022)
FEATURES:
* **New Resource:** sumologic_cse_entity_entity_group_configuration (GH-376)
* **New Resource:** sumologic_cse_inventory_entity_group_configuration (GH-376)
* Add new optional `notification_group_fields` field to resource/sumologic_monitor (GH-403)
* Add new optional `obj_permission` set to resource/sumologic_monitor for Fine Grain Permission (FGP) support (GH-397)
* Add use_versioned_api parameter for s3 source (GH-401)

BUG FIXES:
* Default to NIL for optional timezome field in SumoLogic source (GH-392)
* Allow Monitor move between Monitor folders (GH-405)

## 2.16.2 (June 12, 2022)

BUG FIXES:
* Monitor Folder provider now handles more error codes: "api_not_enabled", in addition to: "not_implemented_yet" (GH-389)

## 2.16.1 (June 6, 2022)

BUG FIXES:
* Allow locator field in DefaultDateFormat to be empty (GH-384)

## 2.16.0 (May 20, 2022)

FEATURES:
* Add new optional `obj_permission` set to resource/sumologic_monitor_folder for Fine Grain Permission (FGP) support (GH-373)

BUG FIXES:
* Fix bug in cse match list items creation (was timing out due to StateChangeConf on an infinite loop) (GH-377)

## 2.15.0 (May 13, 2022)

FEATURES:
* **New Data Source:** `sumologic_folder` (GH-374)
* **New Resource:** `sumologic_slo` (GH-362)
* **New Resource:** `sumologic_slo_folder` (GH-362)
* Add support for slo based monitors (GH-363)
* Add new optional `alert_name` field to resource/sumologic_monitor (GH-359)


BUG FIXES:
* Add CRITICAL as a valid value for cse_custom_insight severity field (GH-367)
* Fix bug preventing to create more than 100 cse match list items within a cs_match_list (GH-368)

## 2.14.0 (March 30, 2022)

FEATURES:
* **New Resource:** sumologic_cse_match_list (GH-353)

ENHANCEMENTS:
* Add support for SumoCloudSOAR webhook connection (GH-352)

## 2.13.0 (February 24, 2022)

FEATURES:
* **New Resource:** sumologic_content_permission (GH-340)

ENHANCEMENTS:
* Add support for importing folder resource (GH-345)
* Allow AtLeastOnce resolution conditions for Metrics monitors  (GH-346)

## 2.12.0 (February 7, 2022)

FEATURES:
* **New Resource:** Gcp Metrics Source `sumologic_gcp_metrics_source` (GH-329, 332)

ENHANCEMENTS:
* Add support for OTLP in HTTP source resource (GH-335)
* Add backoff on http 429s (GH-338)
* Add `domain` field to the dashboard resource (GH-341)

BUG FIXES:
* Fix to allow more than one topology_label for Dashboard resource (GH-336)
* sumologic_cse_log_mapping split_index as int (GH-333)

## 2.11.5 (December 14, 2021)

BUG FIXES:

* Set admin mode in PostRawPayload method (GH-322)

## 2.11.4 (November 19, 2021)

BUG FIXES:

* Fix unexpected end of JSON input error in folder and dashboard resources (GH-319)

## 2.11.3 (November 17, 2021)

BUG FIXES:

* Fix provider crash when user / role data source is declared with a non-existent identifier (GH-316)

## 2.11.2 (November 11, 2021)

ENHANCEMENTS:

* Add support for SNS subscription in polling sources (GH-311)

## 2.11.1 (November 8, 2021)

FEATURES:

* **New Datasource:** sumologic_user (GH-299)

BUG FIXES:

* Fix occurrence_type for metrics resolution conditions (GH-297)
* Relaxed validation for monitor time range (GH-306)

## 2.11.0 (October 19, 2021)

FEATURES:

* **New Resource:** sumologic_cse_rule_tuning_expression (GH-281)
* **New Resource:** sumologic_cse_entity_criticality_config (GH-275)
* **New Resource:** sumologic_cse_custom_entity_type (GH-275)
* **New Resource:** sumologic_cse_insights_resolution (GH-274)
* **New Resource:** sumologic_cse_insights_status (GH-274)
* **New Resource:** sumologic_cse_insights_configuration (GH-274)
* **New Resource:** sumologic_cse_log_mapping (GH-284)
* **New Datasource:** sumologic_cse_log_mapping_vendor_product (GH-284)
* **New Resource:** sumologic_cse_aggregation_rule (GH-290)
* **New Resource:** sumologic_cse_chain_rule (GH-290)
* **New Resource:** sumologic_cse_match_rule (GH-287)
* **New Resource:** sumologic_cse_threshold_rule (GH-287)
* **New Resource:** sumologic_cse_custom_insight (GH-289)

BUG FIXES:

* Fix hierarchy without a filter not being accepted

## 2.10.0 (September 22, 2021)

* Add a provider option `admin_mode`

FEATURES:

* **New Resource:** sumologic_hierarchy (GH-260)
* **New Resource:** sumologic_cse_network_block (GH-271)

POTENTIALLY BREAKING CHANGES:

* resource/sumologic_policies: Changed all policies to be required. Configurations might need to be updated in
  case some policies were not specified previously. (GH-279)

DEPRECATIONS:

* resource/sumologic_monitor: Deprecated `triggers` in favor of `trigger_conditions` (GH-267)

## 2.9.10 (August 24, 2021)

FEATURES:

* **New Resource:** sumologic_policies (GH-248)
* Add a new optional field `playbook` to resource/sumologic_monitor.
* Add a new optional field `evaluation_delay` to resource/sumologic_monitor.

## 2.9.9 (August 12, 2021)

BUG FIXES:

* resource/sumologic_monitor: Removed deprecation warning for `triggers`.
* seperated docs for sumologic_monitor_folder from docs for sumologic_monitor.
* resource/sumologic_monitor: Fixed docs for `trigger_conditions`.

FEATURES:

* Adding "entityId" as part of SAML API response object.

## 2.9.8 (July 30, 2021)

FEATURES:

* Add support for ServiceNow Incident and Event webhook connection (GH-250)
* Add support for new detection methods to sumologic_monitor (GH-239)

DEPRECATIONS:

* resource/sumologic_monitor: Deprecated `triggers` in favor of `trigger_conditions` (GH-239)

BUG FIXES:

* datasource/sumologic_http_source: fix int64 conversion for `collector_id` (GH-251)

## 2.9.7 (July 22, 2021)

ENHANCEMENTS:

* Upgrade GoLang to support arm_64 (GH-241)

## 2.9.6 (July 9, 2021)

BUF FIXES:

* Allow negative terse values for monitor threshold (GH-230)

## 2.9.5 (July 8, 2021)

ENHANCEMENTS:

* Add validation for monitor resource (GH-223)

BUG FIXES:

* Set error message on failure for content resource (GH-224)

## 2.9.4 (June 24, 2021)

FEATURES:

* **New Resource:** sumologic_token (GH-203)
* **New Datasource:** sumologic_admin_recommended_folder (GH-215)

ENHANCEMENTS:

* Remove requirement of placeholder values for `path` and `authentication` for `sumologic_gcp_source` resource (GH-205)
* Add assertion consumer url to terraform saml configuration (GH-200)

## 2.9.3 (April 26, 2021)

FEATURES:

* Add support Microsoft Teams as Connection Type (GH-186)

## 2.9.2 (April 13, 2021)

FEATURES:

* Kinesis Metrics Source (GH-176)

BUG FIXES:

* Handle optional time range in Panel (GH-175)

## 2.9.1 (March 19, 2021)

BUG FIXES:

* Role Data Source now supports names with spaces (GH-172)

ENHANCEMENTS:

* Trace Source Docs (GH-170)

## 2.9.0 (March 9, 2021)

FEATURES:

* Password Policy Resource (GH-161)
* SAML Resource (GH-163)

## 2.8.0 (February 19, 2021)

FEATURES:

* Dashboard (New) Native Terraform Support (GH-146)

## 2.7.1 (Febuary 11, 2021)

FEATURES:

* Import Scheduled Views (GH-152)
* Partitions now use `continuous` as default value for `analytics_tier` (GH-154)

## 2.7.0 (Febuary 5, 2021)

BUG FIXES:

* is_active is now required in sumologic_user, doc fixes
* monitor doc fixes and improvements

## 2.6.3 (January 15, 2021)

FEATURES:

* Roles Data Source (GH-136)

## 2.6.2 (December 23, 2020)

ENHANCEMENTS:

* Partitions now use isActive field to handle decommissioned partitions (GH-131)

## 2.6.1 (December 14, 2020)

ENHANCEMENTS:

* partitions now support import (GH-126)
* adds a new GCP source (GH-120)

## 2.6.0 (November 23, 2020)

FEATURES:

* **New Resource:** sumologic_subdomain (GH-114)

## 2.5.0 (November 13, 2020)

FEATURES:

* **New Resource:** sumologic_lookup_table (GH-87)

## 2.4.1 (November 5, 2020)

BUG FIXES:

* Updated docs for sumologic_cloud_to_cloud_source (GH-111)
* Monitors moving between folders (GH-107)

## 2.4.0 (November 3, 2020)

FEATURES:

* **New Resource:** sumologic_cloud_to_cloud_source (GH-93)

## 2.3.3 (October 16, 2020)

BUG FIXES:

* resource/sumologic_collector: Log when collector is nil. (GH-99)

## 2.3.2 (October 15, 2020)

DEPRECATIONS:

* resource/sumologic_monitor: Deprecated `action_type` in notifications in favor of `connection_type`. (GH-94)

DOCS:

* Improved docs for sumologic_monitor resources with webhook connection example

## 2.3.1 (October 15, 2020)

ENHANCEMENTS:

* Added `in` and `fed` to `environment` section (GH-96)

BUG FIXES:

* Fix a bug where `parse_expression` would present a diff without any changes (GH-95)
* Check for erros when getting collector from client (GF-92)

## 2.3.0 (September 29, 2020)

FEATURES:

* **New Resource:** sumologic_field (GH-82)

## 2.2.2 (September 11, 2020)

FEATURES:

* **New Resource:** sumologic_ingest_budget_v2 (GH-78)

## 2.2.1 (September 8, 2020)

FEATURES:

* **New Resource:** sumologic_monitor (GH-74)
* **New Resource:** sumologic_monitor_folder (GH-74)

## 2.2.0 (September 1, 2020)
DEPRECATIONS:

* resource/sumologic_polling_source: Deprecated in faovur of having individul sources for each content type currently supported. (GH-64)

FEATURES:

* **New Resource:** sumologic_aws_inventory_source (GH-69)
* **New Resource:** sumologic_aws_xray_source (GH-68)
* **New Resource:** sumologic_s3_source (GH-64)
* **New Resource:** sumologic_s3_audit_source (GH-64)
* **New Resource:** sumologic_cloudwatch_source (GH-64)
* **New Resource:** sumologic_cloudtrail_source (GH-64)
* **New Resource:** sumologic_elb_source (GH-64)
* **New Resource:** sumologic_cloudfront_source (GH-64)
* **New Resource:** sumologic_metadata_source (GH-61)

## 2.1.3 (July 30, 2020)
ENHANCEMENTS:
* Allow updates to content resources so that dashboard links do not exprie. This creates a known bug - do not update the name of a resource.

## 2.1.2 (July 24, 2020)
ENHANCEMENTS:
* Now parrt of the Terraform Registry - compatible with Terraform 0.13

## 2.1.1 (July 17, 2020)

DOCS:
* Add docs for common source properties ([#50](https://github.com/SumoLogic/terraform-provider-sumologic/pull/50))

BUG FIXES:
* Fix bug for detecting polling source path type changes on polling source ([#44](https://github.com/SumoLogic/terraform-provider-sumologic/pull/44))

## 2.1.0 (June 22, 2020)

ENHANCEMENTS:
* Add support for cloud watch metrics in the sumologic polling source ([#24](https://github.com/SumoLogic/terraform-provider-sumologic/pull/24))

DOCS:
* Fixed usage example for sumologic_content resource.

DEPRECATIONS:
* resource/sumologic_collector: Remove deprecated attributes `lookup_by_name` and `destroy` ([#32](https://github.com/SumoLogic/terraform-provider-sumologic/pull/32))
* resource/sumologic_sources: Remove deprecated attributes `lookup_by_name` and `destroy` ([#32](https://github.com/SumoLogic/terraform-provider-sumologic/pull/32))


## 2.0.3 (June 02, 2020)

BUG FIXES:

* Fix URL redirection to prod when no base url or deployment is provided ([#20](https://github.com/terraform-providers/terraform-provider-github/issues/20))

## 2.0.2 (May 29, 2020)

ENHANCEMENTS:

* Check for errors when setting aggregate values on read ([#8](https://github.com/terraform-providers/terraform-provider-github/issues/8))

BUG FIXES:

* Fixes updates to content items not being recognized ([#11](https://github.com/terraform-providers/terraform-provider-github/issues/11))

## 2.0.1 (April 30, 2020)

FEATURES:

* **New Resource:** `sumologic_connection` ([#3](https://github.com/terraform-providers/terraform-provider-github/issues/3))

## 2.0.0 (April 09, 2020)

FEATURES:

* **New Resource:** `sumologic_partition` ([#79](https://github.com/terraform-providers/terraform-provider-github/issues/79))
* **New Resource:** `sumologic_field_extraction_rule` ([#83](https://github.com/terraform-providers/terraform-provider-github/issues/83))
* **New Resource:** `sumologic_content` ([#85](https://github.com/terraform-providers/terraform-provider-github/issues/85))

ENHANCEMENTS:

* Travis running acceptance tests ([#92](https://github.com/terraform-providers/terraform-provider-github/issues/92))

BUG FIXES:

* resource/sumologic_partition: Fixes decomissioning of partitions ([#86](https://github.com/terraform-providers/terraform-provider-github/issues/86))

DEPRECATIONS:

* resource/sumologic_collector_ingest_budget_assignment: Deprecated in favor of assigning ingest budgets through the _budget field attribute of collectors ([#135](https://github.com/terraform-providers/terraform-provider-github/issues/135))
* resource/sumologic_collector: Deprecated `lookup_by_name` and `destroy` attributes ([#136](https://github.com/terraform-providers/terraform-provider-github/issues/136))
* resource/sumologic_sources: Deprecated `lookup_by_name` and `destroy` attributes ([#137](https://github.com/terraform-providers/terraform-provider-github/issues/137))
