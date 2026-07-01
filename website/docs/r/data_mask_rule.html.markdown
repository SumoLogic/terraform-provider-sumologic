---
layout: "sumologic"
page_title: "SumoLogic: sumologic_data_mask_rule"
description: |-
  Provides a Sumologic Data Mask Rule
---

# sumologic_data_mask_rule
Provides a [Sumologic Data Mask Rule][1].

Data mask rules allow you to define regex patterns for masking sensitive data in your search results at query time.

## Example Usage

### Basic rule
```hcl
resource "sumologic_data_mask_rule" "email_mask" {
  name          = "mask_email_addresses"
  regex_pattern = "[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\\.[a-zA-Z]{2,}"
  mask_string   = "[EMAIL REDACTED]"
  enabled       = true
  description   = "Masks email addresses in search results"
}
```

### Rule with default mask string
```hcl
resource "sumologic_data_mask_rule" "ip_mask" {
  name          = "mask_ip_addresses"
  regex_pattern = "\\b\\d{1,3}\\.\\d{1,3}\\.\\d{1,3}\\.\\d{1,3}\\b"
  enabled       = true
  description   = "Masks IP addresses in search results"
}
```

## Argument Reference

The following arguments are supported:

- `name` - (Required, ForceNew) The name of the data mask rule. Must be between 1 and 128 characters. Cannot be changed after creation.
- `regex_pattern` - (Required) A valid regular expression pattern to match sensitive data that should be masked.
- `mask_string` - (Optional) The replacement string to use when masking matched data. Must be between 1 and 64 characters. Defaults to `##redactedPII##`.
- `enabled` - (Required) Whether the rule is active. Only enabled rules are applied to search results.
- `description` - (Optional) A description of the rule. Maximum 512 characters.

## Attributes Reference

The following attributes are exported:

- `id` - The unique identifier of the data mask rule.

## Import

Data mask rules can be imported using the rule ID.

```hcl
terraform import sumologic_data_mask_rule.example 000000000ABC1234
```

[1]: https://help.sumologic.com/docs/manage/data-masking/
