---
layout: 'sumologic'
page_title: 'SumoLogic: sumologic_csoar_playbook'
description: |-
  Provides the ability to update, and delete CSOAR playbooks.
---

# sumologic_csoar_playbook

Provides the ability to update, and delete [CSOAR Playbooks][1].

**Note**: CSOAR playbooks cannot be created via Terraform. They must first be created in the CSOAR UI, exported as JSON, and then imported into Terraform using `terraform import`.

## Example Usage

### Complete Playbook Configuration

```hcl
resource "sumologic_csoar_playbook" "complex_playbook" {
  name         = "Complex Security Playbook"
  description  = "A comprehensive incident response playbook"
  tags         = "security,advanced,automation"
  is_deleted   = false
  draft        = false
  is_published = true
  nested       = false
  type         = "General"
  is_enabled   = true
  
  nodes = [
    {
      id = "start-node-1"
      type = "devs.Start"
      position = {
        x = 0
        y = 0
      }
      attrs = {
        ".title" = {
          text = "START"
        }
      }
      stepType = "start"
      elementType = "START"
    },
    {
      id = "action-node-1"
      type = "devs.Action"
      position = {
        x = 200
        y = 0
      }
      attrs = {
        ".title" = {
          text = "Send Alert"
        }
        ".description" = {
          text = "Send notification to security team"
        }
      }
      stepType = "action"
      actionType = "Notification"
    },
    {
      id = "end-node-1"
      type = "devs.End"
      position = {
        x = 400
        y = 0
      }
      attrs = {
        ".title" = {
          text = "END"
        }
      }
      stepType = "end"
      elementType = "END"
    }
  ]
  
  links = [
    {
      id = "link-1"
      type = "link"
      source = {
        id = "start-node-1"
        port = "out"
      }
      target = {
        id = "action-node-1"
        port = "in"
      }
    },
    {
      id = "link-2"
      type = "link"
      source = {
        id = "action-node-1"
        port = "out"
      }
      target = {
        id = "end-node-1"
        port = "in"
      }
    }
  ]
}
```

### Playbook with Dynamic Updates Using JSON File

```hcl
locals {
  playbooks = jsondecode(file("${path.module}/playbooks.json"))
}

resource "sumologic_csoar_playbook" "playbooks" {
  for_each = local.playbooks
  
  name         = each.value.name
  description  = each.value.description
  updated_name = try(each.value.updated_name, null)
  tags         = each.value.tags
  is_deleted   = each.value.is_deleted
  draft        = each.value.draft
  is_published = each.value.is_published
  nested       = each.value.nested
  type         = each.value.type
  is_enabled   = each.value.is_enabled
  nodes        = jsonencode(each.value.nodes)
  links        = jsonencode(each.value.links)
}
```

### Playbook with Name Update

```hcl
resource "sumologic_csoar_playbook" "renamed_playbook" {
  name         = "Original Playbook Name"
  updated_name = "New Playbook Name"
  description  = "Playbook with updated name"
  tags         = "renamed,updated"
  is_enabled   = true
}
```

## Argument Reference

The following arguments are supported:

- `name` - (Required) The current name of the playbook. This is used to identify the playbook for import and updates.
- `description` - (Optional) The description of the playbook. Supports HTML markup.
- `updated_name` - (Optional) When specified, this will rename the playbook to the new name.
- `tags` - (Optional) Comma-separated tags for the playbook.
- `is_deleted` - (Optional) Whether the playbook is deleted.
- `draft` - (Optional) Whether the playbook is in draft state.
- `is_published` - (Optional) Whether the playbook is published.
- `nested` - (Optional) Whether the playbook is nested.
- `type` - (Optional) The type of playbook.
- `is_enabled` - (Optional) Whether the playbook is enabled.
- `nodes` - (Optional) JSON string representation of playbook nodes.
- `links` - (Optional) JSON string representation of playbook links that connect the nodes together to form the workflow.

Additional data provided in state:

- `id` - (Computed) The ID for this playbook (same as name).
- `last_updated` - (Computed) Timestamp of when the playbook was last updated.
- `created_by` - (Computed) ID of the user who created the playbook.
- `updated_by` - (Computed) ID of the user who last updated the playbook.

## Import

CSOAR playbooks must be imported before they can be managed by Terraform. Import using the playbook name:

```bash
terraform import sumologic_csoar_playbook.example "My Playbook Name"
```

## Important Notes

1. **Import-Only Resource**: Playbooks cannot be created through Terraform. They must first be created in the CSOAR UI.

2. **Name-Based Identification**: Playbooks are identified by their name.

3. **JSON Encoding**: The `nodes` and `links` fields require proper JSON encoding. Use `jsonencode()` function when defining these in HCL.

4. **State Management**: Since playbooks cannot be read via API, Terraform manages state based on the configuration provided.

5. **Updates Only**: This resource is designed for updating existing playbook metadata, not for comprehensive playbook management.

## Example JSON File Structure

For use with dynamic configuration, here's an example `playbooks.json` structure:

```json
{
  "security_incident": {
    "name": "Security Incident Response",
    "description": "Automated security incident response playbook",
    "tags": "security,incident,automation",
    "is_deleted": false,
    "draft": false,
    "is_published": true,
    "nested": false,
    "type": "General",
    "is_enabled": true,
    "nodes": [
      {
        "id": "start",
        "type": "devs.Start",
        "stepType": "start"
      }
    ],
    "links": []
  }
}
```

[1]: https://www.sumologic.com/help/docs/platform-services/automation-service/playbooks/