---
layout: "sumologic"
page_title: "SumoLogic: sumologic_cse_tag_schema"
description: |-
  Provides a Sumologic CSE Tag Schema
---

# sumologic_cse_tag_schema
Provides a Sumologic CSE Tag Schema.

## Example Usage
```hcl
resource "sumologic_cse_tag_schema" "tag_schema" {
	key = "location"
	label = "label"
	content_types = ["entity"]
	free_form = "true"	    
	value_options {
    	value = "option value"
    	label = "option label"
		link = "http://foo.bar.com"
    }
}

```

## Argument reference

The following arguments are supported:

- `key` - (Required) Tag Schema key.
- `label` - (Required) Tag Schema label.
- `content_types` - (Required) Applicable content types. Valid values: "customInsight", "entity", "rule", "threatIntelligence".
- `free_form` - (Required) Whether the tag schema accepts free form custom values.
- `value_options` - (At least one need to be added) 
  + `value` - (Required) Value option value.
  + `label` - (Required) Value option label.
  + `link` - (Optional) Value option link.



The following attributes are exported:

- `id` - The internal ID of the Tag Schema.

## Import

Tag Schema can be imported using the field id, e.g.:
```hcl
terraform import sumologic_cse_tag_schema.tag_schema id
```
