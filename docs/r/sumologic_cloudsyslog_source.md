# sumologic_cloudsyslog_source

Provides a [Sumologic Cloud Syslog source][1].

__IMPORTANT:__ The token is stored in plain-text in the state. This is a potential security issue.

## Example Usage
```hcl
resource "sumologic_cloudsyslog_source" "cloudsyslog_source" {
  name         = "CLOUDSYSLOG"
  description  = "My description"
  category     = "my/source/category"
  collector_id = "${sumologic_collector.collector.id}"
}

resource "sumologic_collector" "collector" {
  name        = "my-collector"
  description = "Just testing this"
}
```

## Argument reference

The following arguments are supported:
- `name` - (Required) The name of the source. This is required, and has to be unique in the scope of the collector. Changing this will force recreation the source.
- `description` - (Optional) Description of the source.
- `collector_id` - (Required) The ID of the collector to attach this source to.
- `category` - (Optional) The source category this source logs to.

## Attributes reference

The following attributes are exported:
- `id` - The internal ID of the source.
- `token` - The token to use for sending data to this source.

[Back to Index][0]

[0]: ../README.md
[1]: https://help.sumologic.com/Send_Data/Sources/02Sources_for_Hosted_Collectors/Cloud_Syslog_Source
