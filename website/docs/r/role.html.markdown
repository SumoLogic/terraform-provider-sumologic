---
layout: "sumologic"
page_title: "SumoLogic: sumologic_role"
description: |-
  Provides a Sumologic Role
---

# sumologic_role
Provides a [Sumologic Role][1].

## Example Usage
```hcl
resource "sumologic_role" "example_role" {
  name        = "TestRole123"
  description = "Testing resource sumologic_role"
  filter_predicate = "_sourceCategory=Test"
  capabilities = [
    "manageCollectors"
  ]
}
```

## Argument reference

The following arguments are supported:

- `name` - (Required) The name of the role.
- `description` - (Optional) The description of the role.
- `filter_predicate` - (Optional) A search filter to restrict access to specific logs.
- `capabilities` - (Optional) List of capabilities associated with this role. Valid Values are:  
  - Data Management:  
    - viewCollectors
    - manageCollectors
    - manageBudgets
    - manageDataVolumeFeed
    - viewFieldExtraction
    - manageFieldExtractionRules
    - manageS3DataForwarding
    - manageContent
    - dataVolumeIndex
    - manageConnections
    - viewScheduledViews
    - manageScheduledViews
    - viewPartitions
    - managePartitions
    - viewFields
    - manageFields
    - viewAccountOverview
    - manageTokens
  - Entity management:
    - manageEntityTypeConfig
  - Metrics:
    - metricsTransformation
    - metricsExtraction
    - metricsRules
  - Security:
    - managePasswordPolicy
    - ipAllowlisting
    - createAccessKeys
    - manageAccessKeys
    - manageSupportAccountAccess
    - manageAuditDataFeed
    - manageSaml
    - shareDashboardOutsideOrg
    - manageOrgSettings
    - changeDataAccessLevel
  - Dashboards:
    - shareDashboardWorld
    - shareDashboardAllowlist
  - UserManagement:
    - manageUsersAndRoles
  - Observability:
    - searchAuditIndex
    - auditEventIndex
  - Cloud SIEM Enterprise:
    - viewCse
  - Alerting:
    - viewMonitorsV2
    - manageMonitorsV2
    - viewAlerts


 
  





The following attributes are exported:

- `id` - The internal ID of the role.

## Import
Roles can be imported using the role id, e.g.:

```hcl
terraform import sumologic_role.role 1234567890
```

[1]: https://help.sumologic.com/Manage/Users-and-Roles/Manage-Roles
