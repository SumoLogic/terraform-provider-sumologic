---
layout: "sumologic"
page_title: "Provider: SumoLogic"
description: |-
  This provider is used to manage resources supported by Sumo Logic.
---

# Sumo Logic Provider
This provider is used to manage resources supported by Sumo Logic. The provider needs to be configured with the proper credentials before it can be used.

## Example Usage
```hcl
terraform {
    required_providers {
        sumologic = {
            source = "sumologic/sumologic"
            version = "" # set the Sumo Logic Terraform Provider version
        }
    }
    required_version = ">= 0.13"
}

# Setup authentication variables. See "Authentication" for more details.
variable "sumologic_access_id" {
    type = string
    description = "Sumo Logic Access ID"
}
variable "sumologic_access_key" {
    type = string
    description = "Sumo Logic Access Key"
    sensitive = true
}

# Configure the Sumo Logic Provider
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
$ export SUMOLOGIC_ACCESSKEY="your-access-key"
$ export SUMOLOGIC_ENVIRONMENT=us2
$ terraform plan
```

## Argument Reference
- `access_id` - (Required) This is the Sumo Logic Access ID. It must be provided, but it can also be source from the SUMOLOGIC_ACCESSID environment variable.
- `access_key` - (Required) This is the Sumo Logic Access Key. It must be provided, but it can also be sourced from the SUMOLOGIC_ACCESSKEY variable.
- `environment` - (Required) This is the API endpoint to use. See the [Sumo Logic documentation](https://help.sumologic.com/APIs/General_API_Information/Sumo_Logic_Endpoints_and_Firewall_Security) for details on which environment you should use. It must be provided, but it can be sourced from the SUMOLOGIC_ENVIRONMENT variable.
