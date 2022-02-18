package sumologic

import (
	"errors"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

func resourceSumologicPermissions() *schema.Resource {
	return &schema.Resource{
		Create: resourceSumologicPermissionsCreate,
		Read:   resourceSumologicPermissionsRead,
		Delete: resourceSumologicPermissionsDelete,
		Update: resourceSumologicPermissionsUpdate,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

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
				Optional: true,
				Default:  "",
			},
			"permission": {
				Type:     schema.TypeSet,
				Required: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"permission_name": {
							Type:     schema.TypeString,
							Required: true,
							ValidateFunc: validation.StringInSlice(
								[]string{"View", "GrantView", "Edit", "GrantEdit", "Manage", "GrantManage"}, false),
						},
						"source_type": {
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: validation.StringInSlice([]string{"user", "role", "org"}, false),
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
		id, err := c.UpdatePermissions(PermissionsRequest{
			Permissions: resourceToPermissionsArray(d.Get("permission").(*schema.Set),
				d.Get("content_id").(string)),
			NotifyRecipients:    d.Get("notify_recipient").(bool),
			NotificationMessage: d.Get("notification_message").(string),
		}, d.Get("content_id").(string))

		if err != nil {
			return err
		}
		d.SetId(id)
	}

	return resourceSumologicPermissionsRead(d, meta)
}

func resourceSumologicPermissionsRead(d *schema.ResourceData, meta interface{}) error {
	c := meta.(*Client)

	var permissionsResponse *PermissionsResponse
	id := d.Id()

	permissionsResponse, err := c.GetPermissions(id)
	if err != nil {
		log.Printf("[WARN] Error getting permissions for content(id=%s), err: %v", id, err)
		return err
	}

	if permissionsResponse == nil {
		log.Printf("[WARN] Permissions not found for content(id=%s), removing from state. err: %v",
			id, err)
		d.SetId("")
		return nil
	}

	creatorId, _ := getCreatorId(id, meta)
	if creatorId == "" {
		log.Printf("[WARN] Creator id is empty for this content %v", id)
	}

	d.Set("permission",
		permissionsArrayToResource(permissionsResponse.ExplicitPermissions, creatorId))
	log.Printf("[WARN] Content id %v", d.Get("content_id"))

	return nil
}

func resourceSumologicPermissionsDelete(d *schema.ResourceData, meta interface{}) error {
	resourceSumologicPermissionsRead(d, meta)
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
	resourceSumologicPermissionsDelete(d, meta)
	if _, err = c.UpdatePermissions(permissionRequest, d.Id()); err != nil {
		return err
	}
	return resourceSumologicPermissionsRead(d, meta)
}

func getCreatorId(contentId string, meta interface{}) (string, error) {
	c := meta.(*Client)
	path, err := c.GetContentPath(contentId)
	if err != nil {
		log.Printf("[WARN] Cannot get path for content %v - %v", contentId, err)
		return "", err
	}
	if path == "" {
		log.Printf("[WARN] Path is empty %v", contentId)
		return "", nil
	}
	creatorId, err := c.GetCreatorId(path)
	if err != nil {
		log.Printf("[WARN] Cannot get content by path %v - %v", contentId, err)
		return "", err
	}
	if creatorId == "" {
		log.Printf("[WARN] Creator ID is empty %v", contentId)
	}
	return creatorId, nil
}

func resourceToPermissionsArray(resourcePermissions *schema.Set, contentId string) []Permission {
	result := make([]Permission, resourcePermissions.Len())

	for i, resourcePermission := range resourcePermissions.List() {
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

func permissionsArrayToResource(permissions []Permission, creatorId string) []map[string]interface{} {

	result := make([]map[string]interface{}, 0)

	for _, permission := range permissions {
		if permission.SourceType == "user" && permission.SourceId == creatorId {
			continue
		}
		result = append(result, map[string]interface{}{
			"permission_name": permission.PermissionName,
			"source_type":     permission.SourceType,
			"source_id":       permission.SourceId,
		})
	}

	return result
}

func resourceToPermissionRequest(d *schema.ResourceData) (PermissionsRequest, error) {
	id := d.Id()
	if id == "" {
		return PermissionsRequest{}, errors.New("premission resource id not specified")
	}

	return PermissionsRequest{
		Permissions:         resourceToPermissionsArray(d.Get("permission").(*schema.Set), id),
		NotifyRecipients:    d.Get("notify_recipient").(bool),
		NotificationMessage: d.Get("notification_message").(string),
	}, nil
}
