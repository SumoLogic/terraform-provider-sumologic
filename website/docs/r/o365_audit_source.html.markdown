---
layout: "sumologic"
page_title: "SumoLogic: sumologic_o365_audit_source"
description: |-
  Provides a Sumologic O365 Audit source
---

# sumologic_o365_audit_source
Provides a [Sumologic O365 Audit source][1] for collecting Office 365 audit logs and security events.

## Example Usage

### Basic O365 Exchange Audit Source
```hcl
resource "sumologic_collector" "o365_collector" {
  name        = "O365 Collector"
  description = "Collector for O365 audit logs"
}

resource "sumologic_o365_audit_source" "exchange_audit" {
  name         = "o365-exchange-audit"
  description  = "O365 Exchange Audit Logs"
  category     = "o365/exchange"
  collector_id = sumologic_collector.o365_collector.id

  third_party_ref {
    resources {
      service_type = "O365AuditNotification"

      path {
        type     = "O365NotificationPath"
        workload = "Audit.Exchange"
        region   = "Commercial"
      }

      authentication {
        type          = "O365AppRegistrationAuthentication"
        tenant_id     = var.o365_tenant_id
        client_id     = var.o365_client_id
        client_secret = var.o365_client_secret
      }
    }
  }
}
```

### O365 SharePoint Audit Source
```hcl
resource "sumologic_o365_audit_source" "sharepoint_audit" {
  name         = "o365-sharepoint-audit"
  description  = "O365 SharePoint Audit Logs"
  category     = "o365/sharepoint"
  collector_id = sumologic_collector.o365_collector.id

  third_party_ref {
    resources {
      service_type = "O365AuditNotification"

      path {
        type     = "O365NotificationPath"
        workload = "Audit.SharePoint"
        region   = "Commercial"
      }

      authentication {
        type          = "O365AppRegistrationAuthentication"
        tenant_id     = var.o365_tenant_id
        client_id     = var.o365_client_id
        client_secret = var.o365_client_secret
      }
    }
  }
}
```

### O365 Azure Active Directory Audit Source for GCC High
```hcl
resource "sumologic_o365_audit_source" "azuread_audit" {
  name         = "o365-azuread-audit"
  description  = "O365 Azure AD Audit Logs"
  category     = "o365/azuread"
  collector_id = sumologic_collector.o365_collector.id

  third_party_ref {
    resources {
      service_type = "O365AuditNotification"

      path {
        type     = "O365NotificationPath"
        workload = "Audit.AzureActiveDirectory"
        region   = "GCC High"
      }

      authentication {
        type          = "O365AppRegistrationAuthentication"
        tenant_id     = var.o365_tenant_id
        client_id     = var.o365_client_id
        client_secret = var.o365_client_secret
      }
    }
  }
}
```

### O365 DLP Source
```hcl
resource "sumologic_o365_audit_source" "dlp_logs" {
  name         = "o365-dlp-logs"
  description  = "O365 Data Loss Prevention Logs"
  category     = "o365/dlp"
  collector_id = sumologic_collector.o365_collector.id

  third_party_ref {
    resources {
      service_type = "O365AuditNotification"

      path {
        type     = "O365NotificationPath"
        workload = "DLP.All"
        region   = "Commercial"
      }

      authentication {
        type          = "O365AppRegistrationAuthentication"
        tenant_id     = var.o365_tenant_id
        client_id     = var.o365_client_id
        client_secret = var.o365_client_secret
      }
    }
  }
}
```

## Argument Reference

In addition to the [Common Source Properties](https://registry.terraform.io/providers/SumoLogic/sumologic/latest/docs#common-source-properties), the following arguments are supported:

- `third_party_ref` - (Required) Configuration block for O365 third-party reference.
  - `resources` - (Required) List of resource configurations. Currently, only one resource is supported.
    - `service_type` - (Required) The service type. Must be `O365AuditNotification`.
    - `path` - (Required) Configuration block for the O365 notification path.
      - `type` - (Required) The path type. Must be `O365NotificationPath`.
      - `workload` - (Required) The Office 365 workload to collect audit logs from. Valid values are:
        - `Audit.Exchange` - Exchange audit logs
        - `Audit.AzureActiveDirectory` - Azure Active Directory audit logs
        - `Audit.SharePoint` - SharePoint audit logs
        - `Audit.General` - General audit logs
        - `DLP.All` - Data Loss Prevention logs
      - `region` - (Required) The Office 365 deployment region. Valid values are:
        - `Commercial` - Commercial cloud (default)
        - `GCC` - Government Community Cloud
        - `GCC High` - Government Community Cloud High
    - `authentication` - (Required) Configuration block for O365 app registration authentication.
      - `type` - (Required) The authentication type. Must be `O365AppRegistrationAuthentication`.
      - `tenant_id` - (Required) The Azure AD tenant ID (directory ID).
      - `client_id` - (Required) The Azure AD application (client) ID.
      - `client_secret` - (Required) The Azure AD client secret value. This is marked as sensitive.

### See also
  * [Common Source Properties](https://registry.terraform.io/providers/SumoLogic/sumologic/latest/docs#common-source-properties)

## Attributes Reference

The following attributes are exported:

- `id` - The internal ID of the source.
- `url` - The HTTP endpoint to use for receiving O365 audit notifications.

## Prerequisites

Before creating an O365 Audit source, you need to:

1. **Register an Azure AD Application** in your Office 365 tenant
2. **Configure API Permissions** for the Office 365 Management APIs:
   - `ActivityFeed.Read` - Read activity data for your organization
   - `ActivityFeed.ReadDlp` - Read DLP policy events including detected sensitive data
3. **Create a Client Secret** for the application
4. **Grant Admin Consent** for the permissions

For detailed setup instructions, see the [Office 365 Management Activity API documentation][2].

## Import

O365 Audit sources can be imported using the collector and source IDs (`collector/source`), e.g.:

```hcl
terraform import sumologic_o365_audit_source.exchange_audit 123/456
```

O365 Audit sources can be imported using the collector name and source name (`collectorName/sourceName`), e.g.:

```hcl
terraform import sumologic_o365_audit_source.exchange_audit my-o365-collector/my-exchange-source
```

## Notes

- Each workload type requires a separate source configuration.
- Make sure your Azure AD application has been granted admin consent before creating the source.
- The source will automatically subscribe to the specified Office 365 audit log content type.

[1]: https://help.sumologic.com/docs/send-data/hosted-collectors/cloud-to-cloud-integration-framework/microsoft-office-audit-source/
[2]: https://docs.microsoft.com/en-us/office/office-365-management-api/office-365-management-activity-api-reference
