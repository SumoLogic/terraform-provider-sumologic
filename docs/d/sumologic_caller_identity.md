# sumologic_caller_identity
Provides an easy way to retrive Sumologic auth details.


## Example Usage
```hcl
data "sumologic_caller_identity" "current" {}
```


## Attributes reference
The following attributes are exported:
- `access_id` - Sumologic access ID.
- `access_key` - Sumologic access key.
- `environment` - API endpoint environment.

[Back to Index][0]

[0]: ../README.md

