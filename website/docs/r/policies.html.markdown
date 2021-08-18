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
}
```

## Argument reference

The following arguments are supported:

- `audit` - (Optional) Whether the [Audit Policy][1] is enabled. Defaults to `false`.
- `data_access_level` - (Optional) Whether the [Data Access Level Policy][2] is enabled. Defaults to `false`.
- `max_user_session_timeout` - (Optional) The [maximum web session timeout][3] users are able to configure within their user preferences. Defaults to `7d`.
- `search_audit` - (Optional) Whether the [Search Audit Policy][4] is enabled. Defaults to `false`.
- `share_dashboards_outside_organization` - (Optional) Whether the [Share a Dashboard Outside Organization Policy][5] is enabled. Defaults to `false`.
- `user_concurrent_sessions_limit` - (Block List, Max: 1, Optional) The [User Concurrent Sessions Limit Policy][6]. See [user_concurrent_sessions_limit schema](#user_concurrent_sessions_limit) for details.

### Schema for `user_concurrent_sessions_limit`
- `enabled` - (Optional) Whether the [User Concurrent Sessions Limit Policy][6] is enabled. Defaults to `false`.
- `max_concurrent_sessions` - (Optional) Maximum number of concurrent sessions a user may have. Defaults to `100`.

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
