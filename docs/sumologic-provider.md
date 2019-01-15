# Sumologic Provider
This provider is used to manage collectors and sources supported by Sumologic. The provider needs to be configured with the proper credentials before it can be used.

## Example Usage
```hcl
# Configure the Sumologic Provider
provider "sumologic" {
    access_id   = "${var.sumologic_access_id}"
    access_key  = "${var.sumologic_access_key}"
    environment = "us2"
}

# Create a collector
resource "sumologic_collector" "collector" {
    name = "MyCollector"
}

# Create a HTTP source
resource "sumologic_http_source" "http_source" {
    name         = "http-source"
    category     = "my/source/category"
    collector_id = "${sumologic_collector.collector.id}"
}

# Create HTTP source / SNS auto-confirm
resource "sumologic_http_source_sns_autoconfirm" "sns_confirm" {
  category             = "aws/config/${var.project_name}/${var.env}"
  confirmation_timeout = 3

  triggers {
    http_source_url = "${sumologic_http_source.http_source.url}"
  }
}
```

## Authentication
The Sumo Logic Provider offers a flexible means of providing credentials for authentication. The following methods are supported and explained below:

 - Static credentials
 - Environment variables

### Static credentials
Static credentials can be provided by adding an `access_id` and `access_key` in-line in the Sumo Logic provider block:

Usage:
```hcl
provider "sumologic" {
    environment = "us2"
    access_id   = "your-access-id"
    access_key  = "your-access-key"
}
```

### Environment variables
You can provide your credentials via the `SUMOLOGIC_ACCESSID` and `SUMOLOGIC_ACCESSKEY` environment variables, representing your Sumo Logic Access ID and Sumo Logic Access Key, respectively.

Usage:
```hcl
provider "sumologic" { }
```

```bash
$ export SUMOLOGIC_ACCESSID="your-access-id"
$ export SUMOLOGIC_ACCESSKEY="your-access-id"
$ terraform plan
```

## Argument Reference
- `access_id` - (Optional) This is the Sumo Logic Access ID. It must be provided, but it can also be source from the SUMOLOGIC_ACCESSID environment variable.
- `access_key` - (Optional) This is the Sumo Logic Access Key. It must be provided, but it can also be sourced from the SUMOLOGIC_ACCESSKEY variable.
- `environment` - (Optional) This is the API endpoint to use. Default is `us2`. See the [Sumo Logic documentation][1] for details on which environment you should use.

[Back to Index][0]

[0]: README.md
[1]: https://help.sumologic.com/APIs/General_API_Information/Sumo_Logic_Endpoints_and_Firewall_Security
