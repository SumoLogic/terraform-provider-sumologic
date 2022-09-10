---
layout: "sumologic"
page_title: "SumoLogic: sumologic_cse_entity_normalization_configuration"
description: |-
    Provides the Sumologic CSE Entity Normalization Configuration for the whole organization. There can be only one configuration per organization.
---

# sumologic_cse_entity_normalization_configuration
Provides the Sumologic CSE Entity Normalization Configuration for the whole organization. There can be only one configuration per organization.

## Example Usage
```hcl
resource "sumologic_cse_entity_normalization_configuration" "entity_normalization_configuration" {
	windows_normalization_enabled = true
	fqdn_normalization_enabled = true
	aws_normalization_enabled = true
	default_normalized_domain = "domain.com"
	normalize_hostnames = true
	normalize_usernames = true
	domain_mappings{
		normalized_domain = "normalized.domain"
		raw_domain = "raw.domain"
	}
}
```

## Argument reference

The following arguments are supported:

- `windows_normalization_enabled` - (Required) Normalize active directory domains username and hostname formats.
- `fqdn_normalization_enabled` - (Required) Normalize names in the form user@somedomain.net or hostname.somedomain.net
- `aws_normalization_enabled` - (Required) Normalize AWS ARN and Usernames.
- `default_normalized_domain` - (Optional) When normalization is configured, at least one domain must be configured and a "Normalized Default Domain" must be provided.
- `domain_mappings` - (Required) Secondary domains.
    + `normalized_domain` - (Required) The normalized domain.
    + `raw_domain` - (Required) The raw domain to be normalized.
- `normalize_hostnames` - (Required) If hostname normalization is enabled.
- `normalize_usernames` - (Required) If username normalization is enabled.

- The following attributes are exported:

- `ID` - The internal ID of the entity normalization configuration.

## Import

Entity Normalization Configuration can be imported using the id `cse-entity-normalization-configuration`:

~> **NOTE:** Only `cse-entity-normalization-configuration` id should be used when importing the entity normalization configuration. Using any other id may have unintended consequences.

```hcl
terraform import sumologic_cse_entity_normalization_configuration.entity_normalization_configuration cse-entity-normalization-configuration
```