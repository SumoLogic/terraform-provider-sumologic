---
layout: "sumologic"
page_title: "SumoLogic: sumologic_cse_log_mapping_vendor_product"
description: |-
  Provides a way to retrieve Sumo Logic CSE log mapping vendor product guid  to be used by another terraform stack.
---

# sumologic_cse_log_mapping_vendor_product

Provides a way to retrieve Sumo Logic CSE log mapping vendor product guid to be used by another terraform stack.
Refer to [Products with Log Mappings](https://help.sumologic.com/Cloud_SIEM_Enterprise/Ingestion_Guides/00Products_with_Log_Mappings) for a list of available product and vendor names.

## Example Usage
```hcl
data "sumologic_cse_log_mapping_vendor_product" "web_gateway" {
  product = "Web Gateway"
  vendor = "McAfee"
}
```

A Log mapping vendor product can be looked up by providing values of `product` and `vendor`
Both `product` and `vendor` values are mandatory. If not provided an error will be generated.

## Attributes reference

The following attributes are exported:

- `guid` - The internal GUID of the log mapping vendor product.
- `product` - The name of the product.
- `vendor` - The name of the vendor.


