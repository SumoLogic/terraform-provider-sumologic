---
layout: "sumologic"
page_title: "SumoLogic: sumologic_policies"
description: |-
  Sets the Sumologic Policies
---

# sumologic_policies
Sets the Sumologic Policies. Since each policy is global for the entire organization, please ensure that only a single
instance of this resource is defined. The behavior for defining more than one policies resource is undefined.

The following policies are supported:
- [Audit Policy][1]
- [Data Access Level Policy][2]
- [Maximum Web Session Timeout Policy][3]
- [Search Audit Policy][4]
- [Share a Dashboard Outside Organization Policy][5]
- [User Concurrent Sessions Limit Policy][6]
- [Check Data Ingestion Policy][7]

## Example Usage
```hcl
resource "sumologic_policies" "example_policies" {
  audit = false
  data_access_level = false
  max_user_session_timeout = "7d"
  search_audit = false
  share_dashboards_outside_organization = false
  user_concurrent_sessions_limit {
    enabled = false
    max_concurrent_sessions = 100
  }
  check_data_ingestion {
    enabled = true
    no_data_threshold = "24h"
  }
}
```

## Argument reference

The following arguments are supported:

- `audit` - (Required) Whether the [Audit Policy][1] is enabled.
- `data_access_level` - (Required) Whether the [Data Access Level Policy][2] is enabled.
- `max_user_session_timeout` - (Required) The [maximum web session timeout][3] users are able to configure within their user preferences.
- `search_audit` - (Required) Whether the [Search Audit Policy][4] is enabled.
- `share_dashboards_outside_organization` - (Required) Whether the [Share a Dashboard Outside Organization Policy][5] is enabled.
- `user_concurrent_sessions_limit` - (Block List, Max: 1, Required) The [User Concurrent Sessions Limit Policy][6]. See [user_concurrent_sessions_limit schema](#user_concurrent_sessions_limit) for details.
- `check_data_ingestion` - (Block List, Max: 1, Required) The [Check Data Ingestion Policy][7]. See [check_data_ingestion schema](#check_data_ingestion) for details.

### Schema for `user_concurrent_sessions_limit`
- `enabled` - (Required) Whether the [User Concurrent Sessions Limit Policy][6] is enabled.
- `max_concurrent_sessions` - (Optional) Maximum number of concurrent sessions a user may have. Defaults to `100`.

### Schema for `check_data_ingestion`
- `enabled` - (Required) Whether the [Check Data Ingestion Policy][7] is enabled.
- `no_data_threshold` - (Optional) Duration threshold after which an OT collector not sending data will trigger an alert. Valid values are: `1h`, `2h`, `4h`, `8h`, `16h`, or `24h`. Defaults to `24h`.

## Import
Policies can be imported using the id `org-policies`.

~> **NOTE:** Only `org-policies` id should be used when importing policies. Using any other id may have unintended consequences.

```hcl
terraform import sumologic_policies.example_policies org-policies
```

[1]: https://help.sumologic.com/Manage/Security/Audit-Index
[2]: https://help.sumologic.com/Manage/Security/Data_Access_Level_for_Shared_Dashboards
[3]: https://help.sumologic.com/Manage/Security/Set_a_Maximum_Web_Session_Timeout
[4]: https://help.sumologic.com/Manage/Security/Search_Audit_Index
[5]: https://help.sumologic.com/Visualizations-and-Alerts/Dashboards/Share_Dashboards/Share_a_Dashboard_Outside_Your_Org
[6]: https://help.sumologic.com/Manage/Security/Set_a_Limit_for_User_Concurrent_Sessions
[7]: https://help.sumologic.com/docs/ADD-URL-TO-CHECK-DATA-INGESTION-POLICY
