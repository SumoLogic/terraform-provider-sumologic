# sumologic_cloudsyslog_source

Provides a [Sumo Logic Cloud Syslog source][1].

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

Only the common properties are valid.

## Attributes reference

The following attributes are exported:
- `id` - The internal ID of the source.
- `token` - The token to use for sending data to this source.

[Back to Index][0]

[0]: ../README.md
[1]: https://help.sumologic.com/Send_Data/Sources/02Sources_for_Hosted_Collectors/Cloud_Syslog_Source
