---
layout: "sumologic"
page_title: "SumoLogic: sumologic_policies"
description: |-
  Sets the Sumologic Policies
---

# sumologic_policies
Sets the Sumologic Policies. Since each policy is global for the entire organization, please ensure that only a single
instance of this resource is defined. The behavior for defining more than one policies resource is undefined.

The six policies marked Required below are always managed by this resource. The OAuth Clients, Unused Access Keys and
Access Keys Expiration policies are optional: leave them out to keep their current value.

The following policies are supported:
- [Audit Policy][1]
- [Data Access Level Policy][2]
- [Maximum Web Session Timeout Policy][3]
- [Search Audit Policy][4]
- [Share a Dashboard Outside Organization Policy][5]
- [User Concurrent Sessions Limit Policy][6]
- OAuth Clients Policy (Client ID Metadata Documents, CIMD)
- Unused Access Keys Deactivation Policy
- Access Keys Expiration Policy

## Example Usage

### Defaults
Every required argument at the value Sumo Logic uses when the resource is destroyed.
```hcl
resource "sumologic_policies" "defaults" {
  audit                                 = false
  data_access_level                     = false
  max_user_session_timeout              = "7d"
  search_audit                          = false
  share_dashboards_outside_organization = false

  user_concurrent_sessions_limit {
    enabled = false
  }
}
```

### Audit and compliance
Turn on the Audit and Search Audit indexes and cap web sessions at 1 hour.
```hcl
resource "sumologic_policies" "compliance" {
  audit                                 = true
  search_audit                          = true
  data_access_level                     = true
  max_user_session_timeout              = "1h"
  share_dashboards_outside_organization = false

  user_concurrent_sessions_limit {
    enabled                 = true
    max_concurrent_sessions = 3
  }
}
```

### Access key hygiene
Deactivate access keys unused for 60 days and expire all keys after 90 days.
```hcl
resource "sumologic_policies" "access_keys" {
  audit                                 = true
  data_access_level                     = false
  max_user_session_timeout              = "7d"
  search_audit                          = true
  share_dashboards_outside_organization = false

  unused_access_keys_disable_after_in_days = 60
  access_keys_lifetime_in_days             = "90"

  user_concurrent_sessions_limit {
    enabled = false
  }
}
```

### OAuth clients
Allow only OAuth clients that an administrator created manually on the OAuth Clients page.
```hcl
resource "sumologic_policies" "oauth" {
  audit                                 = true
  data_access_level                     = false
  max_user_session_timeout              = "7d"
  search_audit                          = true
  share_dashboards_outside_organization = false

  oauth_cimd_policy = "enabled-pre-registered-only"

  user_concurrent_sessions_limit {
    enabled = false
  }
}
```

## Argument reference

The following arguments are supported:

- `audit` - (Required) Whether the [Audit Policy][1] is enabled. When `true`, audit records for your organization are written to the Audit Index.
- `data_access_level` - (Required) Whether the [Data Access Level Policy][2] is enabled. When `true`, newly created dashboards default to the viewer's role access filter. When `false`, they default to the sharer's role access filter and may show data that viewers' roles do not allow them to view.
- `max_user_session_timeout` - (Required) The [maximum web session timeout][3] users are able to configure within their user preferences. One of `5m`, `15m`, `30m`, `1h`, `2h`, `6h`, `12h`, `1d`, `2d`, `3d`, `5d` or `7d`. Users whose preference is higher than this value are lowered to it; users with a lower preference are not changed.
- `search_audit` - (Required) Whether the [Search Audit Policy][4] is enabled. When `true`, search records for your organization are written to the Search Audit Index.
- `share_dashboards_outside_organization` - (Required) Whether the [Share a Dashboard Outside Organization Policy][5] is enabled. Users also need the matching capability on their role. Setting it to `false` disables every dashboard that is currently shared outside the organization.
- `user_concurrent_sessions_limit` - (Block List, Max: 1, Required) The [User Concurrent Sessions Limit Policy][6]. See [user_concurrent_sessions_limit schema](#schema-for-user_concurrent_sessions_limit) for details.
- `oauth_cimd_policy` - (Optional) The OAuth Clients policy, shown as **OAuth Clients** in the UI. Controls whether OAuth clients can authenticate with Client ID Metadata Documents (CIMD); it does not manage individual OAuth clients. One of `disabled`, `enabled` or `enabled-pre-registered-only`. With `disabled`, CIMD authentication is off and no new CIMD clients can be created. With `enabled`, CIMD authentication is on and new CIMD clients are created automatically during authentication. With `enabled-pre-registered-only`, CIMD authentication is on but CIMD clients can only be created manually on the OAuth Clients page. If not set, the current value in Sumo Logic is kept.
- `unused_access_keys_disable_after_in_days` - (Optional) Number of days after which an unused access key is automatically deactivated, from `0` to `365`. `0` means never. If not set, the current value in Sumo Logic is kept.
- `access_keys_lifetime_in_days` - (Optional) Number of days after which an access key expires unless rotated; an expired key can no longer be used. One of `0`, `30`, `45`, `60`, `90`, `180` or `365`, given as a string. `0` means never. If not set, the current value in Sumo Logic is kept.

~> **NOTE:** Removing an optional argument (`oauth_cimd_policy`, `unused_access_keys_disable_after_in_days`, `access_keys_lifetime_in_days`) from your configuration does not revert it; the value in Sumo Logic is kept. Set it explicitly to change it.

~> **NOTE:** Writing `access_keys_lifetime_in_days` restarts the rotation clock of every access key in the organization, even if the value is unchanged. With a value above `0`, every key gets a full new lifetime from now (including keys close to expiring). With `0`, the rotation history of every key is cleared, so a later non-zero value counts from that moment. The provider only sends this argument on create or when you change it, so unrelated updates do not trigger this.

~> **NOTE:** On `terraform destroy`, the required policies are reset to the values in the [Defaults](#defaults) example: `audit`, `data_access_level`, `search_audit` and `share_dashboards_outside_organization` to `false`, `max_user_session_timeout` to `7d`, and the concurrent sessions limit disabled at `100`. The optional policies above are not reset, because their defaults vary per organization.

### Schema for `user_concurrent_sessions_limit`
- `enabled` - (Required) Whether the [User Concurrent Sessions Limit Policy][6] is enabled. When `false`, users may have an unlimited number of sessions.
- `max_concurrent_sessions` - (Optional) Maximum number of concurrent sessions a user may have, from `1` to `100`. Defaults to `100`. When a user exceeds the limit, their oldest session is logged out.

## Attributes Reference

In addition to the arguments listed above, the following attribute is exported:

- `id` - Always `org-policies`.

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
