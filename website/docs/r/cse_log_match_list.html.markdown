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
    description = "Match list item description"
    value = "action"
    expiration = "2022-02-27T04:00:00"
  }
}
```

## Argument reference

The following arguments are supported:

- `default_ttl` - (Required) The match list time to live. Specified in seconds.
- `description` - (Required) Match list description.
- `name` - (Required) Match list name.
- `target_column` - (Required) Target column. (possible values: Hostname, FileHash, Url, SrcIp, DstIp, Domain, Username, Ip, Asn, Isp, Org, SrcAsn, SrcIsp, SrcOrg, DstAsn, DstIsp, DstOrg or any custom column.)
- `items` - (Optional) List of match list items. See [match_list_item schema](#schema-for-match_list_item) for details.

### Schema for `match_list_item`
- `description` - (Required) Match list item description.
- `value` - (Optional) Match list item value.
- `expiration` - (Optional) Match list item expiration. (Format: YYYY-MM-DDTHH:mm:ss)

The following attributes are exported:

- `id` - The internal ID of the match list.

## Import

Mach List can be imported using the field id, e.g.:
```hcl
terraform import sumologic_cse_match_list.match_list id
```

