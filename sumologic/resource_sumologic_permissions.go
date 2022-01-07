package sumologic

import (
	"errors"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceSumologicPermissions() *schema.Resource {
	return &schema.Resource{
		Create: resourceSumologicPermissionsCreate, // ?
		Read:   resourceSumologicPermissionsRead,
		Delete: resourceSumologicPermissionsDelete,
		Update: resourceSumologicPermissionsUpdate,

		Schema: map[string]*schema.Schema{
			"content_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"notify_recipient": {
				Type:     schema.TypeBool,
				Required: true,
			},
			"notification_message": {
				Type:     schema.TypeString,
				Required: true,
			},
			"permission": {
				Type:     schema.TypeList,
				Required: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"permission_name": {
							Type:     schema.TypeString,
							Required: true,
						},
						"source_type": {
							Type:     schema.TypeString,
							Required: true,
						},
						"source_id": {
							Type:     schema.TypeString,
							Required: true,
						},
					},
				},
			},
		},
	}
}

func resourceSumologicPermissionsCreate(d *schema.ResourceData, meta interface{}) error {
	c := meta.(*Client)

	if d.Id() == "" {
		err := c.UpdatePermissions(PermissionsRequest{
			PermissionAssignmentype: resourceToPermissionsArray(d.Get("permission").([]interface{}), d.Get("content_id").(string)),
			NotifyRecipients:        d.Get("notify_recipient").(bool),
			NotificationMessage:     d.Get("notification_message").(string),
		}, d.Get("content_id").(string))

		if err != nil {
			return err
		}
		d.SetId(d.Get("content_id").(string))
	}

	return resourceSumologicCSEAggregationRuleRead(d, meta)
}

func resourceSumologicPermissionsRead(d *schema.ResourceData, meta interface{}) error {
	c := meta.(*Client)

	var permissionsResponse *PermissionsResponse
	id := d.Id()

	permissionsResponse, err := c.GetPermissions(id)
	if err != nil {
		log.Printf("[WARN] Permissions not found when looking by id: %s, err: %v", id, err)
	}

	if permissionsResponse == nil {
		log.Printf("[WARN] CSE Aggregation Rule not found, removing from state: %v - %v", id, err)
		d.SetId("")
		return nil
	}

	d.Set("permission", permissionsArrayToResource(permissionsResponse.ExplicitPermissions))

	return nil
}

func resourceSumologicPermissionsDelete(d *schema.ResourceData, meta interface{}) error {
	c := meta.(*Client)
	permissionRequest, err := resourceToPermissionRequest(d)
	if err != nil {
		return err
	}
	return c.DeletePermissions(permissionRequest, d.Get("content_id").(string))
}

func resourceSumologicPermissionsUpdate(d *schema.ResourceData, meta interface{}) error {
	c := meta.(*Client)
	permissionRequest, err := resourceToPermissionRequest(d)
	if err != nil {
		return err
	}
	if err = c.UpdatePermissions(permissionRequest, d.Id()); err != nil {
		return err
	}
	return resourceSumologicPermissionsRead(d, meta)
}

func resourceToPermissionsArray(resourcePermissions []interface{}, contentId string) []Permission {
	result := make([]Permission, len(resourcePermissions))

	for i, resourcePermission := range resourcePermissions {
		PermissionMap := resourcePermission.(map[string]interface{})
		result[i] = Permission{
			PermissionName: PermissionMap["permission_name"].(string),
			SourceType:     PermissionMap["source_type"].(string),
			SourceId:       PermissionMap["source_id"].(string),
			ContentId:      contentId,
		}
	}

	return result
}

func permissionsArrayToResource(permissions []Permission) []map[string]interface{} {
	result := make([]map[string]interface{}, len(permissions))

	for i, permission := range permissions {
		result[i] = map[string]interface{}{
			"permission_name": permission.PermissionName,
			"source_type":     permission.SourceType,
			"source_id":       permission.SourceId,
		}
	}

	return result
}

func resourceToPermissionRequest(d *schema.ResourceData) (PermissionsRequest, error) {
	id := d.Id()
	if id == "" {
		return PermissionsRequest{}, errors.New("premission resource id not specified")
	}

	return PermissionsRequest{
		PermissionAssignmentype: resourceToPermissionsArray(d.Get("permission").([]interface{}), id),
		NotifyRecipients:        d.Get("notify_recipient").(bool),
		NotificationMessage:     d.Get("notification_message").(string),
	}, nil
}
