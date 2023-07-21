---
layout: "sumologic"
page_title: "SumoLogic: sumologic_cse_match_list"
description: |-
  Provides a Sumologic CSE Match List
---

# match_list
Provides a Sumologic CSE Match List.

## Example Usage
```hcl
resource "sumologic_cse_match_list" "match_list" {
  default_ttl = 10800
  description = "Match list description"
  name = "Match list name"
  target_column = "SrcIp"
  items {
    description = "IP address"
    value = "192.168.0.1"
    expiration = "2022-02-27T04:00:00"
  }
}
```

## Argument reference

The following arguments are supported:

- `default_ttl` - (Optional) The default time to live for match list items added through the UI. Specified in seconds.
- `description` - (Required) Match list description.
- `name` - (Required) Match list name.
- `target_column` - (Required) Target column. (possible values: Hostname, FileHash, Url, SrcIp, DstIp, Domain, Username, Ip, Asn, Isp, Org, SrcAsn, SrcIsp, SrcOrg, DstAsn, DstIsp, DstOrg or any custom column.)
- `items` - (Optional) List of match list items. See [match_list_item schema](#schema-for-match_list_item) for details.

**Note:** When managing CSE match list items outside of terraform, omit the `items` argument and add `items` to the [ignore_changes](https://developer.hashicorp.com/terraform/language/meta-arguments/lifecycle#ignore_changes) list in a lifecycle block on the `sumologic_cse_match_list` resource. As match list items are added or removed outside of terraform, terraform will ignore these changes, protecting match list items from accidental deletion.

For example:

```hcl
resource "sumologic_cse_match_list" "match_list" {
  default_ttl = 10800
  description = "Match list description"
  name = "Match list name"
  target_column = "SrcIp"
  
  lifecycle {
    # protects match list items added outside terraform from accidental deletion
    ignore_changes = [items]
  }
}
```

### Schema for `match_list_item`
- `description` - (Required) Match list item description.
- `value` - (Optional) Match list item value.
- `expiration` - (Optional) Match list item expiration. (Format: YYYY-MM-DDTHH:mm:ss)

The following attributes are exported:

- `id` - The internal ID of the match list.

## Import

Match List can be imported using the field id, e.g.:
```hcl
terraform import sumologic_cse_match_list.match_list id
```

