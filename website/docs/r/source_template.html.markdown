---
layout: "sumologic"
page_title: "SumoLogic: sumologic_source_template"
description: |-
  Provides a Sumologic Source Template
---

# sumologic_source_template
Provides a Sumologic Source Template

## Example Usage
```hcl
resource "sumologic_source_template" "example_source_template" {
    schema_ref {
        type = "Mac"
    }
    selector {
        tags =  [
             [
                {
                key = "new1"
                values =  ["Abc","Abc2"]
                }
             ]
        ]
        names = ["TestCollector1"]
    }
    is_enabled = true
    input_json = jsonencode({
        "name": "hostmetrics_test_source_template_test",
        "description": "Host metric source" ,
        "receivers": {
            "hostmetrics": {
                "receiverType": "hostmetrics",
                "collection_interval": "5m",
                "cpu_scraper_enabled": true,
                "disk_scraper_enabled": true,
                "load_scraper_enabled": true,
                "filesystem_scraper_enabled": true,
                "memory_scraper_enabled": true,
                "network_scraper_enabled": true,
                "processes_scraper_enabled": true,
                "paging_scraper_enabled": true
            }
        },
        "processors": {
            "resource": {
                "processorType": "resource",
                "user_attributes": [
                    {
                        "key": "_sourceCategory",
                        "value": "otel/host"
                    }
                ],
                "default_attributes": [
                    {
                        "key": "sumo.datasource",
                        "value": "apache"
                    },
                    {
                        "key": "host.name",
                        "value": "host1"
                    },
                    {
                        "key": "host.id",
                        "value": "hostid"
                    },
                    {
                        "key": "log.file.path",
                        "value": "filePath"
                    }
                ]
            }
        }
    })
}
```
## Argument reference

The following arguments are supported:

- `schema_ref` - (Required) Schema reference for source template.
- `input_json` - (Required) This is a JSON object which contains the configuration parameters for the source template.
- `selector` - (Optional) Conditions to select OT Agent.
- `is_enabled` - (Optional) Indicates whether the source template is enabled or disabled.

The following attributes are exported:

- `id` - The internal ID of the source_template.

### Schema for `schema_ref`
- `type` - (Required) Type of schema for the source template.
- `version` - (Optional) Version of schema used for the source template. Takes the latest version, if this field is omitted.

### Schema for `input_json`
- `name` - (Required) Name of the source template.
- `receivers` - (Required) Receiver information of the source template.
- `description` - (Optional) Description of the source template.
- `processors` - (Optional) Processors for the source template.
- `property name*` - (Optional) Any other property as needed.

### Schema for `selector`
- `tags` - (Optional) Tags filter for OT agents. It is an Array of Array of OT tag objects. Objects within same array are evaluated by 
AND logic whereas separate arrays are evaluated using OR logic.
- `names` - (Optional) Names of OT collectors to select particular agents.

### Schema for `OT tag`
- `key` - (Required) Key of the needed OT tag.
- `values` - (Required) Values of the given OT tag to be filtered.

## Import
Source Templates can be imported using the ST id, e.g.:

```hcl
terraform import sumologic_source_template.test 0000000000000004
```