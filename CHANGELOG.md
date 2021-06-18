## 2.9.4 (Unreleased)

Features:

* **New Resource:** sumologic_token (GH-203)

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
