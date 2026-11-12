---
layout: "sumologic"
page_title: "SumoLogic: sumologic_usage_attribution_config"
description: |-
  Turns usage attribution on or off and sets its custom tags.
---

# sumologic_usage_attribution_config
Turns usage attribution on or off for your organization and sets the custom tags that ingest volume can be attributed to.
Usage attribution (also called Granular Usage Visibility) breaks down ingest volume by tags such as `sourceCategory` or
`collector`, and by custom tags you add, on the Usage Attribution page.

There is exactly one usage attribution configuration per organization. Define only one `sumologic_usage_attribution_config`
resource; the behavior of several is undefined.

The user or access key the provider runs as needs the **Manage Usage Management** capability. Reading the
configuration, including the data source, needs **View Usage Management**.

## Example Usage

### Turn usage attribution on
```hcl
resource "sumologic_usage_attribution_config" "attribution" {
  enabled = true
}
```

### With custom tags
```hcl
resource "sumologic_usage_attribution_config" "attribution" {
  enabled     = true
  custom_tags = ["department", "cost_center"]
}
```

### Use the available breakdowns elsewhere
```hcl
resource "sumologic_usage_attribution_config" "attribution" {
  custom_tags = ["department"]
}

output "usage_attribution_breakdowns" {
  value = sumologic_usage_attribution_config.attribution.breakdown_by_options
}
```

## Argument reference

The following arguments are supported:

- `enabled` - (Optional) Turns usage attribution on (`true`) or off (`false`). Defaults to `true`. Turning it on starts
  attribution for data ingested from then on; data ingested earlier is not attributed. Turning it off stops attribution
  and keeps the custom tags and the dashboard.
- `custom_tags` - (Optional) The names of the custom tags to attribute volume to, in addition to the standard tags. Each
  name must start with a letter or underscore and contain only letters, digits and underscores
  (`^[A-Za-z_][A-Za-z0-9_]*$`). A name that is already a standard tag is ignored and a name given twice is kept once. The
  list replaces the custom tags currently configured, so an omitted or empty list removes them all. A tag added here
  applies only to data ingested after it is added.

## Attributes reference

The following attributes are exported:

- `id` - The fixed value `usageAttributionConfig`.
- `standard_tags` - The names of the standard tags, which are always available and cannot be changed.
- `breakdown_by_options` - Everything the usage attribution dashboard can be broken down by: the standard tags, the custom
  tags and any columns of an ownership lookup.
- `dashboard_id` - The ID of the Usage Attribution dashboard. Empty until the dashboard has been created, which happens
  in the Sumo Logic UI.
- `enabled_at` - When usage attribution was first turned on, in ISO 8601 format (UTC). Empty while it is off.

## Destroy behavior

The configuration cannot be deleted. `terraform destroy` sets `enabled` to `false`. It does not change the custom tags
or the dashboard.

## Import

The configuration is imported with its fixed ID:

```shell
terraform import sumologic_usage_attribution_config.attribution usageAttributionConfig
```
