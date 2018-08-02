# sumologic_http_source
Provides a [HTTP source / SNS auto-confirm][1].

__IMPORTANT:__ The endpoint is stored in plain-text in the state. This is a potential security issue.

## Example Usage
```hcl
resource "sumologic_http_source_sns_autoconfirm" "sns_confirm" {
  category             = "${sumologic_http_source.http_source.category}"
  confirmation_timeout = 3

  triggers {
    http_source_url = "${sumologic_http_source.http_source.url}"
  }
}

resource "sumologic_http_source" "http_source" {
  name                = "HTTP"
  description         = "My description"
  message_per_request = true
  category            = "my/source/category"
  collector_id        = "${sumologic_collector.collector.id}"
}
```

## Argument reference
The following arguments are supported:
- `category` - (Required) The source category of the HTTP source. This is required, should be the same as the HTTP source category. Changing this will force recreation the HTTP source / SNS auto-confirm.
- `confirmation_timeout` - (Optional) Integer indicating number of minutes to wait in retying mode for confirmation message before marking it as failure. (default is 3 minute)
- `triggers` -   A mapping of values which should trigger a rerun of resource. Values are meant to be interpolated references to variables or attributes of other resources.


[Back to Index][0]

[0]: ../README.md