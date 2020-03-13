## 2.0.0 (Unreleased)

FEATURES:

* **New Resource:** `sumologic_partition` [GH-79]
* **New Resource:** `sumologic_field_extraction_rule` [GH-83]
* **New Resource:** `sumologic_content` [GH-85]

ENHANCEMENTS:

* Travis running acceptance tests [GH-92]

BUG FIXES:

* resource/sumologic_partition: Fixes decomissioning of partitions [GH-86]

DEPRECATIONS:

* resource/sumologic_collector_ingest_budget_assignment: Deprecated in favor of assigning ingest budgets through the _budget field attribute of collectors [GH-135]
* resource/sumologic_collector: Deprecated `lookup_by_name` and `destroy` attributes [GH-136]
* resource/sumologic_sources: Deprecated `lookup_by_name` and `destroy` attributes [GH-137]