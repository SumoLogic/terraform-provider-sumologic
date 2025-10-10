---
layout: "sumologic"
page_title: "SumoLogic: sumologic_macro"
description: |-
  Provides a Sumologic Macro
---

# sumologic_macro (beta)
Provides a [Sumologic Macro (Beta)][1].

The feature is in beta, will not function if not enabled for your org. Please reach out to sumologic support engineer to have the feature enabled for your org.

## Example Usage
Single macro
```hcl

resource "sumologic_macro" "ip_macro" {
  name = "ip_macro"
  definition = "_sourceCategory=yourcategory | where ip = {{ip_address}} | timeslice 5m | count by _timeslice"
  argument {
    name = "ip_address"
    type = "String"
  }
  argument_validation {
    eval_expression = "isValidIP(ip_address)"
    error_message = "The ip you provided is invalid"
  }
}
```

Macros with dependencies
```hcl

resource "sumologic_macro" "ip_macro" {
  depends_on = [sumologic_macro.ip_macro_nested]
  name = "ip_macro"
  definition = "_sourceCategory=yourcategory | `ip_macro_nested`"
}

resource "sumologic_macro" "ip_macro_nested" {
  name = "ip_macro_nested"
  definition = "_sourceCategory=yourcategory | count"
}
```

## Argument reference

The following arguments are supported:

- `name` - (Required) Name of the macro.
- `description` - (Optional) Description of the macro.
- `definition` - (Required) The definition of your macro
- `enabled` - (Optional) Whether the macro will be enabled. Default true.
- `argument` - (Block List, Optional) A list of arguments for the macro. They must match the arguments in the definition. See [argument schema](#schema-for-argument) for details.
- `argumentValidations` - (Block List, Optional) A list validations for the arguments in the macro. See [argumentValidation schema](#schema-for-argumentvalidation)
for details.

## Attributes reference
In addition to all arguments above, the following attributes are exported:

- `id` - The ID of the macro.

### Schema for `argument`
- `name` - (Required) Name of the argument.
- `type` - (Required) Type of the argument. Must be String, Any, Number or Keyword

### Schema for `argumentValidation`
- `evalExpression` - (Required) The expression to validate a macro argument.
- `errorMessage` - (Required) Error message to show when the argument validation failed.


[1]: https://help.sumologic.com/docs/search/search-query-language/search-operators/macro/


