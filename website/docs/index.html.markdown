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

# Setup authentication variables. See "Authentication" section for more details.
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

# Configure the Sumo Logic Provider in Admin Mode
provider "sumologic" {
    access_id   = "${var.sumologic_access_id}"
    access_key  = "${var.sumologic_access_key}"
    environment = "us2"
    admin_mode  = true
    alias       = "admin"
}

# Look up the Admin Recommended Folder
data "sumologic_admin_recommended_folder" "folder" {}

# Create a folder underneath the Admin Recommended Folder (which requires Admin Mode)
resource "sumologic_folder" "test" {
    provider    = sumologic.admin
    name        = "test"
    description = "A test folder"
    parent_id   = data.sumologic_admin_recommended_folder.folder.id
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

## Common Source Properties

The following properties are common to ALL sources and can be used to configure each source.

- `collector_id` - (Required) The ID of the collector to attach this source to.
- `name` - (Required) The name of the source. This is required, and has to be unique in the scope of the collector. Changing this will force recreation the source.
- `description` - (Optional) Description of the source.
- `category` - (Optional) The source category this source logs to.
- `host_name` - (Optional) The source host this source logs to.
- `timezone` - (Optional) The timezone assigned to the source. The value follows the [tzdata][2] naming convention.
- `automatic_date_parsing` - (Optional) Determines if timestamp information is parsed or not. Type true to enable automatic parsing of dates (the default setting); type false to disable. If disabled, no timestamp information is parsed at all.
- `multiline_processing_enabled` - (Optional) Type true to enable; type false to disable. The default setting is true. Consider setting to false to avoid unnecessary processing if you are collecting single message per line files (for example, Linux system.log). If you're working with multiline messages (for example, log4J or exception stack traces), keep this setting enabled.
- `use_autoline_matching` - (Optional) Type true to enable if you'd like message boundaries to be inferred automatically; type false to prevent message boundaries from being automatically inferred (equivalent to the Infer Boundaries option in the UI). The default setting is true.
- `manual_prefix_regexp` - (Optional) When using useAutolineMatching=false, type a regular expression that matches the first line of the message to manually create the boundary. Note that any special characters in the regex, such as backslashes or double quotes, must be escaped.
- `force_timezone` - (Optional) Type true to force the source to use a specific time zone, otherwise type false to use the time zone found in the logs. The default setting is false.
- `default_date_formats` - (Optional) Define the format for the timestamps present in your log messages. You can specify a locator regex to identify where timestamps appear in log lines. Requires 'automatic_date_parsing' set to True. 
  + `format` - (Required) The timestamp format supplied as a Java SimpleDateFormat, or "epoch" if the timestamp is in epoch format.
  + `locator` - (Optional) Regular expression to locate the timestamp within the messages.  

  Usage:
  ```hcl
     default_date_formats {
       format = "MM-dd-yyyy HH:mm:ss"
       locator = "timestamp:(.*)\\s"
     }
  ```
- `filters` - (Optional) If you'd like to add a filter to the source.
  + `filter_type` - (Required) The type of filter to apply. (Exclude, Include, Mask, or Hash)
  + `name` - (Required) The Name for the filter. 
  + `regexp` - (Required) Regular expression to match within the messages. When used with Incude/Exclude the expression must match the entire message. When used with Mask/Hash rules the expression must contain an unnamed capture group to hash/mask. 
  + `mask` - (Optional) When applying a Mask rule, replaces the detected expression with this string.

  Usage:
  ```hcl
     filters {
       filter_type = "Include"
       name = "Sample Include"
       regexp = ".*\\d{16}.*"
     }
     filters {
       filter_type = "Mask"
       name = "Sample Mask"
       regexp = "(\\d{16})"
       mask = "MaskedID"
     }
  ```  
- `cutoff_timestamp` - (Optional) Only collect data more recent than this timestamp, specified as milliseconds since epoch (13 digit). This maps to the `Collection should begin` field on the UI. Example: using `1663786159000` will set the cutoff timestamp to `Wednesday, September 21, 2022 6:49:19 PM GMT`
- `cutoff_relative_time` - (Optional) Can be specified instead of cutoffTimestamp to provide a relative offset with respect to the current time.This maps to the `Collection should begin` field on the UI. Example: use -1h, -1d, or -1w to collect data that's less than one hour, one day, or one week old, respectively.
- `fields` - (Optional) Map containing key/value pairs.
 
   Usage:
   ```hcl
     fields = {
       environment = "production"
       service = "apache"
     }
   ```
## Configuring SNS Subscription
This is supported in the following resources.
 - `sumologic_cloudfront_source`
 - `sumologic_cloudtrail_source`
 - `sumologic_elb_source`
 - `sumologic_s3_audit_source`
 - `sumologic_s3_source`

Steps to configure SNS subscription and sync the state in terrform:
     - Step 1: Create the source via terraform.
     - Step 2: Setup [SNS subscription][3] outside of terraform on Sumologic UI.
     - Step 3: Run `terraform plan -refresh-only` to review the changes and verify the state with the SNS subscription information. Make sure only `sns_topic_or_subscription_arn` is updated. If SNS has been successfully configured and has received a subscription confirmation request `isSuccess` parameter will be true.
     - Step 4: Apply the changes with `terraform apply -refresh-only`.

[2]: https://en.wikipedia.org/wiki/Tz_database
[3]: https://help.sumologic.com/03Send-Data/Sources/02Sources-for-Hosted-Collectors/Amazon-Web-Services/AWS_Sources#set-up-sns-in-aws-highly-recommended

