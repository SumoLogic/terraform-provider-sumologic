---
layout: "sumologic"
page_title: "SumoLogic: sumologic_usage_attribution_config"
description: |-
  Provides the current usage attribution configuration.
---

# sumologic_usage_attribution_config
Provides the current usage attribution configuration of your organization without changing it. Use it to read the
custom tags or the dashboard ID when the configuration is managed outside Terraform, or by a
[`sumologic_usage_attribution_config` resource](../r/usage_attribution_config.html) in another workspace.

The user or access key the provider runs as needs the **View Usage Management** capability.

## Example Usage

### Read the current configuration
```hcl
data "sumologic_usage_attribution_config" "current" {
}

output "usage_attribution_enabled" {
  value = data.sumologic_usage_attribution_config.current.enabled
}

output "usage_attribution_breakdowns" {
  value = data.sumologic_usage_attribution_config.current.breakdown_by_options
}
```

### Fail a plan when usage attribution is off
```hcl
data "sumologic_usage_attribution_config" "current" {
}

resource "terraform_data" "require_usage_attribution" {
  lifecycle {
    precondition {
      condition     = data.sumologic_usage_attribution_config.current.enabled
      error_message = "Usage attribution must be turned on."
    }
  }
}
```

## Attributes reference

The data source takes no arguments. The following attributes are exported:

- `id` - The fixed value `usageAttributionConfig`.
- `enabled` - Whether usage attribution is on.
- `custom_tags` - The names of the custom tags, in the order Sumo Logic returns them.
- `standard_tags` - The names of the standard tags.
- `breakdown_by_options` - Everything the usage attribution dashboard can be broken down by: the standard tags, the
  custom tags and any columns of an ownership lookup.
- `dashboard_id` - The ID of the Usage Attribution dashboard. Empty until the dashboard has been created.
- `enabled_at` - When usage attribution was first turned on, in ISO 8601 format (UTC). Empty while it is off.
