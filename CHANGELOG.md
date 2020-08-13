## 2.2.0 (Unreleased)

## 2.1.3 (July 30, 2020)
ENHANCEMENTS:
* Allow updates to content resources so that dashboard links do not exprie. This creates a known bug - do not update the name of a resource.

## 2.1.2 (July 24, 2020)
ENHANCEMENTS:
* Now parrt of the Terraform Registry - compatible with Terraform 0.13

## 2.1.1 (July 17, 2020)

DOCS:
* Add docs for common source properties ([#50](https://github.com/terraform-providers/terraform-provider-sumologic/pull/50))

BUG FIXES:
* Fix bug for detecting polling source path type changes on polling source ([#44](https://github.com/terraform-providers/terraform-provider-sumologic/pull/44))

## 2.1.0 (June 22, 2020)

ENHANCEMENTS:
* Add support for cloud watch metrics in the sumologic polling source ([#24](https://github.com/terraform-providers/terraform-provider-sumologic/pull/24))

DOCS:
* Fixed usage example for sumologic_content resource.

DEPRECATIONS:
* resource/sumologic_collector: Remove deprecated attributes `lookup_by_name` and `destroy` ([#32](https://github.com/terraform-providers/terraform-provider-sumologic/pull/32))
* resource/sumologic_sources: Remove deprecated attributes `lookup_by_name` and `destroy` ([#32](https://github.com/terraform-providers/terraform-provider-sumologic/pull/32))


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
