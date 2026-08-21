---
layout: "sumologic"
page_title: "SumoLogic: sumologic_lookup_table"
description: |-
  Provides a Sumologic Lookup Table
---

# sumologic_lookup_table
Provides a [Sumologic Lookup Table][1].

## Example Usage
```hcl
data "sumologic_personal_folder" "personalFolder" {}

resource "sumologic_lookup_table" "lookupTable" {
    name = "Sample Lookup Table"
    fields {
      field_name = "FieldName1"
      field_type = "boolean"
    }
    fields {
      field_name = "FieldName2"
      field_type = "string"
    }
    ttl               = 100
    primary_keys      = ["FieldName1"]
    parent_folder_id  = "${data.sumologic_personal_folder.personalFolder.id}"
    size_limit_action = "DeleteOldData"
    description       = "some description"
    content           = file("${path.module}/lookup.csv")
}
```

## Argument reference

The following arguments are supported:

- `name` - (Required) The name of the lookup table.
- `parent_folder_id` - (Required) The parent-folder-path identifier of the lookup table in the Library.
- `description` - (Required) The description of the lookup table.
- `fields` - (Required) The list of fields in the lookup table.
  - `fieldName` - (Required) The name of the field.
  - `fieldType` - (Required) The data type of the field. Supported types: boolean, int, long, double, string
- `primaryKeys` - (Required) The names of the fields that make up the primary key for the lookup table. These will be a subset of the fields that the table will contain.
- `ttl` - (Optional) A time to live for each entry in the lookup table (in minutes). 365 days is the maximum time to live for each entry that you can specify. Setting it to 0 means that the records will not expire automatically.
- `sizeLimitAction` - (Optional) The action that needs to be taken when the size limit is reached for the table. The possible values can be StopIncomingMessages or DeleteOldData. DeleteOldData will start deleting old data once size limit is reached whereas StopIncomingMessages will discard all the updates made to the lookup table once size limit is reached.
- `content` - (Optional) CSV content to populate the lookup table with, e.g. `file("${path.module}/lookup.csv")`. Notes:
  - The first row must be a header with column names matching `fields.field_name` (case-insensitive).
  - Every field in `fields` must have a matching column, and every field listed in `primary_keys` must be present.
  - Uploading replaces all existing rows in the table (it is not merged with existing data).
  - Removing this argument truncates the table rather than leaving it unmanaged.
  - This value cannot be read back from the Sumo Logic API. Importing an existing table will show a diff on the first plan if `content` is set locally, and the state stores a hash of the content rather than the content itself.
  - CSV structure (headers, column count, quoting) is validated at plan time. Cell values are not validated against `field_type` locally - that is enforced by the API at apply time.
  - The API's upload limit is 100MB per CSV; values over 10MB produce a plan-time warning recommending an ingest pipeline or scheduled search instead.
  - If the upload completes with some rows rejected (e.g. duplicate primary keys), `terraform apply` still succeeds but shows a warning with the rejected-row detail.

## Attributes reference

The following attributes are exported:

- `id` - Unique identifier for the partition.

## Import
Lookup Tables can be imported using the id, e.g.:

```hcl
terraform import sumologic_lookup_table.test 1234567890
```

[1]: https://help.sumologic.com/05Search/Lookup_Tables
