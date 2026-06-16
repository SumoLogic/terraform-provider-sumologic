---
layout: "sumologic"
page_title: "SumoLogic: sumologic_data_mask_rule"
description: |-
  Provides a Sumologic Data Mask Rule
---

# sumologic_data_mask_rule
Provides a [Sumologic Data Mask Rule][1].

Data mask rules allow you to define patterns for masking sensitive data (PII) in your log messages. Rules can be scoped to a single org, specific child orgs, or all child orgs for multi-org management.

## Example Usage

### Basic rule scoped to current org
```hcl
resource "sumologic_data_mask_rule" "email_mask" {
  name        = "mask_email_addresses"
  pattern     = "[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\\.[a-zA-Z]{2,}"
  pii_type    = "email"
  replacement = "[EMAIL REDACTED]"
  scope       = "org"
  enabled     = true
  description = "Masks email addresses in log data"
}
```

### Rule applied to specific child orgs
```hcl
resource "sumologic_data_mask_rule" "ssn_mask_child_orgs" {
  name                 = "mask_ssn_child_orgs"
  pattern              = "\\d{3}-\\d{2}-\\d{4}"
  pii_type             = "ssn"
  replacement          = "[SSN REDACTED]"
  scope                = "child_org"
  scope_target_org_ids = ["000000000ABC1234", "000000000DEF5678"]
  enabled              = true
  description          = "Masks SSNs in selected child orgs"
}
```

### Rule applied to all child orgs
```hcl
resource "sumologic_data_mask_rule" "phone_mask_all" {
  name        = "mask_phone_all_orgs"
  pattern     = "\\b\\d{3}[-.\\s]?\\d{3}[-.\\s]?\\d{4}\\b"
  pii_type    = "phone"
  replacement = "[PHONE REDACTED]"
  scope       = "all_orgs"
  enabled     = true
  description = "Masks phone numbers across all child orgs"
}
```

## Argument Reference

The following arguments are supported:

- `name` - (Required) The name of the data mask rule. Must be between 1 and 255 characters.
- `pattern` - (Required) A valid regular expression pattern to match sensitive data.
- `pii_type` - (Required) The type of PII the rule targets. Supported values: `phone`, `email`, `ip`, `ssn`, `credit_card`, `custom`.
- `replacement` - (Required) The replacement string to use when masking matched data. Must be between 1 and 512 characters.
- `scope` - (Required) The scope of the rule. Supported values: `org`, `child_org`, `all_orgs`.
- `scope_target_org_ids` - (Optional) A list of child org IDs to apply the rule to. Required when `scope` is `child_org`.
- `enabled` - (Optional) Whether the rule is enabled. Defaults to `true`.
- `description` - (Optional) A description of the rule. Maximum 1024 characters.

## Attributes Reference

The following attributes are exported:

- `id` - The unique identifier of the data mask rule.
- `is_active` - Whether the rule is currently active.

## Import

Data mask rules can be imported using the rule ID.

```hcl
terraform import sumologic_data_mask_rule.example 000000000ABC1234
```

[1]: https://help.sumologic.com/docs/manage/data-masking/
