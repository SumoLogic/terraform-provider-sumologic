---
layout: "sumologic"
page_title: "SumoLogic: sumologic_local_file_source"
description: |-
  Provides a Sumologic Local File Source.
---

# sumologic_local_file_source
Provides a [Sumologic Local File Source][1].

## Example Usage
```hcl
resource "sumologic_installed_collector" "installed_collector" {
  name        = "test-collector"
  category = "macos/test"
  ephemeral = true
}

resource "sumologic_local_file_source" "local" {
	name = "localfile-mac"
	description = "test"
	category = "test"
	collector_id = "${sumologic_installed_collector.installed_collector.id}"
  path_expression = "/Applications/Sumo Logic Collector/logs/*.log.*"
}
```

## Argument Reference

The following arguments are supported:

  * `name` - (Required) The name of the local file source. This is required, and has to be unique. Changing this will force recreation the source.
  * `path_expression` - (Required) A valid path expression (full path) of the file to collect. For files on Windows systems (not including Windows Events), enter the absolute path including the drive letter. Escape special characters and spaces with a backslash (). If you are collecting from Windows using CIFS/SMB, see Prerequisites for Windows Log Collection. Use a single asterisk wildcard [*] for file or folder names. Example:[var/foo/*.log]. Use two asterisks [**]to recurse within directories and subdirectories. Example: [var/*/.log].
  * `description` - (Optional) The description of the source.
  * `category` - (Optional) The default source category for the source.
  * `fields` - (Optional) Map containing [key/value pairs][2].
  * `denylist` - (Optional) Comma-separated list of valid path expressions from which logs will not be collected.
    Example: "denylist":["/var/log/**/*.bak","/var/oldlog/*.log"]
  * `encoding` - (Optional) Defines the encoding form. Default is "UTF-8". Other supported encodings are listed [here][3].

### See also
  * [Common Source Properties](https://github.com/terraform-providers/terraform-provider-sumologic/tree/master/website#common-source-properties)

## Attributes Reference
The following attributes are exported:

  * `id` - The internal ID of the local file source.

## Import
Local file sources can be imported using the collector and source IDs, e.g.:

```hcl
terraform import sumologic_local_file_source.test 123/456
```

Local file sources can also be imported using the collector name and source name, e.g.:

```hcl
terraform import sumologic_local_file_source.test my-test-collector/my-test-source
```

[1]: https://help.sumologic.com/docs/send-data/installed-collectors/sources/local-file-source/
[2]: https://help.sumologic.com/Manage/Fields
[3]: https://help.sumologic.com/docs/send-data/installed-collectors/sources/local-file-source/#supported-encoding-for-local-file-sources