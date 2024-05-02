---
layout: "sumologic"
page_title: "SumoLogic: sumologic_role_v2"
description: |-
  Provides a Sumologic Role V2
---

# sumologic_roleV2
Provider to manage Sumologic Role V2s

## Example Usage
```hcl
resource "sumologic_role_v2" "example_role_v2" {
    selected_views = ""
    name = "DataAdmin"
    audit_data_filter = "info"
    selection_type = "Allow"
    capabilities = "["manageContent","manageDataVolumeFeed","manageFieldExtractionRules","manageS3DataForwarding"]"
    description = "Manage data of the org."
    security_data_filter = "error"
    log_analytics_filter = "!_sourceCategory=collector"
    selected_views = "[{"viewName": "view1"}, {"viewName": "view2"}]"
}
```
## Argument reference

The following arguments are supported:

- `selected_views` - (Optional) List of views which with specific view level filters in accordance to the selectionType chosen.
- `name` - (Required) Name of the role.
- `audit_data_filter` - (Optional) A search filter which would be applied on partitions which belong to Audit Data product area. Help Doc : (https://help.sumologic.com/docs/manage/security/audit-index/). Applicable with only `All` selectionType
- `selection_type` - (Optional) Describes the Permission Construct for the list of views in "selectedViews" parameter. 
### Valid Values are : 
  - `All` selectionType would allow access to all views in the org.
  - `Allow` selectionType would allow access to specific views mentioned in "selectedViews" parameter.
  - `Deny` selectionType would deny access to specific views mentioned in "selectedViews" parameter.
- `capabilities` - (Optional) List of [capabilities](https://help.sumologic.com/docs/manage/users-roles/roles/role-capabilities/) associated with this role. Valid values are
### Data Management
  - viewCollectors
  - manageCollectors
  - manageBudgets
  - manageDataVolumeFeed
  - viewFieldExtraction
  - manageFieldExtractionRules
  - manageS3DataForwarding
  - manageContent
  - manageApps
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
  - downloadSearchResults
  - manageIndexes
  - manageDataStreams
  - viewParsers
  - viewDataStreams

### Entity management
  - manageEntityTypeConfig

### Metrics
  - metricsTransformation
  - metricsExtraction
  - metricsRules

### Security
  - managePasswordPolicy
  - ipAllowlisting
  - ipWhitelisting
  - createAccessKeys
  - manageAccessKeys
  - manageSupportAccountAccess
  - manageAuditDataFeed
  - manageSaml
  - shareDashboardOutsideOrg
  - manageOrgSettings
  - changeDataAccessLevel

### Dashboards
  - shareDashboardWorld
  - shareDashboardAllowlist
  - shareDashboardWhitelist

### UserManagement
  - manageUsersAndRoles

### Observability
  - searchAuditIndex
  - auditEventIndex

### Cloud SIEM Enterprise
  - viewCse
  - cseViewAutomations
  - cseManageContextActions
  - cseViewNetworkBlocks
  - cseManageInsightTags
  - cseViewRules
  - cseViewThreatIntelligence
  - cseCommentOnInsights
  - cseViewEntityGroups
  - cseManageEntityConfiguration
  - cseManageNetworkBlocks
  - cseManageMatchLists
  - cseViewCustomInsights
  - cseManageActions
  - cseManageAutomations
  - cseManageMappings
  - cseManageThreatIntelligence
  - cseViewActions
  - cseCreateInsights
  - cseManageTagSchemas
  - cseInvokeInsights
  - cseManageCustomEntityType
  - cseViewTagSchemas
  - cseDeleteInsights
  - cseManageCustomInsights
  - cseViewFileAnalysis
  - cseManageFileAnalysis
  - cseManageEntityCriticality
  - cseViewEntityCriticality
  - cseViewEntity
  - cseManageCustomInsightStatuses
  - cseViewContextActions
  - cseViewMappings
  - cseViewCustomEntityType
  - cseManageEntityGroups
  - cseViewCustomInsightStatuses
  - cseViewEnrichments
  - cseManageInsightSignals
  - cseManageRules
  - cseManageArtifacts
  - cseViewMatchLists
  - cseManageInsightPolicy
  - cseManageEnrichments
  - cseViewEntityConfiguration
  - cseManageEntity
  - cseExecuteAutomations
  - cseManageSuppressedEntities
  - cseManageInsightStatus  
  - cseManageInsightAssignee
  - cseManageFavoriteFields
  - cseViewSuppressedEntities

### Alerting
  - viewMonitorsV2
  - manageMonitorsV2
  - viewAlerts
  - viewMutingSchedules
  - manageMutingSchedules
  - adminMonitorsV2

### SLO
  - viewSlos
  - manageSlos

### CloudSoar
  - cloudSoarPlaybooksAccess
  - cloudSoarNotificationConfigure
  - cloudSoarReportAll
  - cloudSoarIncidentTriageAccess
  - cloudSoarIncidentTaskView
  - cloudSoarIncidentChangeOwnership
  - cloudSoarIncidentNotesEdit
  - cloudSoarAPIEmailEdit
  - cloudSoarIncidentTemplatesAccess
  - cloudSoarIncidentPlaybooksManage
  - cloudSoarGeneralConfigure
  - cloudSoarEntitiesAccess
  - cloudSoarEntitiesBulkPhysicalDelete
  - cloudSoarIncidentAttachmentsAccess
  - cloudSoarAppCentralAccess
  - cloudSoarBridgeMonitoringAccess
  - viewCloudSoar
  - cloudSoarIncidentView
  - cloudSoarObservabilityAccess
  - cloudSoarAPIEmailRead
  - cloudSoarAppCentralExport
  - cloudSoarWidgetsAll
  - cloudSoarIncidentTaskReassign
  - cloudSoarIntegrationsAccess
  - cloudSoarCustomizationIncidentLabels
  - cloudSoarAutomationRulesConfigure
  - cloudSoarIncidentTaskAccessAll
  - cloudSoarAuditAndInformationConfigureAuditTrail
  - cloudSoarIncidentTriageEdit
  - cloudSoarIncidentEdit
  - cloudSoarNotificationTriage
  - cloudSoarIncidentTriageBulkPhysicalDelete
  - cloudSoarIncidentNotesAccess
  - cloudSoarAPIUse
  - cloudSoarIncidentPlaybooksEdit
  - cloudSoarDashboardAll
  - cloudSoarEntitiesManage
  - cloudSoarIncidentTemplatesConfigure
  - cloudSoarIncidentTriageAccessAll
  - cloudSoarPlaybooksConfigure
  - cloudSoarIncidentAccessAll
  - cloudSoarCustomizationLogo
  - cloudSoarIncidentTaskAccess
  - cloudSoarIncidentTriageView
  - cloudSoarIntegrationsConfigure
  - cloudSoarIncidentManageInvestigators
  - cloudSoarIncidentAccess
  - cloudSoarAuditAndInformationLicenseInformation
  - cloudSoarIncidentBulkOperations
  - cloudSoarCustomizationFields
  - cloudSoarIncidentTaskEdit
  - cloudSoarDashboardAccess
  - cloudSoarIncidentAttachmentsEdit
  - cloudSoarIncidentFoldersEdit
  - cloudSoarUserManagementGroups
  - cloudSoarIncidentPlaybooksAccess
  - cloudSoarIncidentWarRoomUse
  - cloudSoarReportAccess
  - cloudSoarAuditAndInformationAuditTrail
  - cloudSoarAutomationRulesAccess
  - cloudSoarIncidentTriageChangeOwnership
  - cloudSoarObservabilityManagement
- `description` - (Optional) Description of the role.
- `security_data_filter` - (Optional) A search filter which would be applied on partitions which belong to Security Data product area. Applicable with only `All` selectionType.
- `log_analytics_filter` - (Optional) A search filter which would be applied on partitions which belong to Log Analytics product area. Applicable with only `All` selectionType

The following attributes are exported:

- `id` - The internal ID of the role_v2



[Back to Index][0]

[0]: ../README.md
