package sumologic

import (
	"github.com/hashicorp/terraform/helper/schema"
	"log"
)

func resourceSumologicPermission() *schema.Resource {
	return &schema.Resource{
		Create: resourceSumologicPermissionCreate,
		Read:   resourceSumologicPermissionRead,
		Delete: resourceSumologicPermissionDelete,
		Update: resourceSumologicPermissionUpdate,

		Schema: map[string]*schema.Schema{
			"content_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: false,
			},
			"permission": {
				Type:     schema.TypeSet,
				Required: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"source_id": {
							Type:     schema.TypeString,
							Required: true,
						},
						"source_type": {
							Type:     schema.TypeString,
							Required: true,
						},
						"permission_name": {
							Type:     schema.TypeString,
							Required: true,
						},
					},
				},
			},
		},
	}
}

func resourceSumologicPermissionRead(d *schema.ResourceData, meta interface{}) error {

	log.Println("====Begin Permission Read====")

	c := meta.(*Client)
	//retrieve the permission Id from the state
	id := d.Id()
	log.Printf("Permission Id from schema: %s", id)

	contentPermission, err := c.GetPermission(id)

	//Error retrieving permission
	if err != nil {
		return err
	}

	//ensure the Permission is populated
	if len(contentPermission.Permissions) == 0 {
		log.Printf("Permissions not found, removing from state: %v - %v", id, err)
		d.SetId("")
		return nil
	}

	// Write the newly read permission into the schema
	d.Set("permission", contentPermission.Permissions)

	log.Println("====End Permission Read====")
	return nil
}

func resourceSumologicPermissionDelete(d *schema.ResourceData, meta interface{}) error {
	log.Println("====Begin Permission Delete====")
	c := meta.(*Client)

	deletedPermissions := resourceToPermission(d)
	log.Printf("Deleting Permissions: %s", deletedPermissions.Permissions)

	return c.DeletePermissions(deletedPermissions)
}

func resourceSumologicPermissionCreate(d *schema.ResourceData, meta interface{}) error {
	log.Println("====Begin Permission Create====")
	c := meta.(*Client)

	//If there is no id in the state, then we need to create the object
	if d.Id() == "" {

		//Load all the data we have from the schema into a Permission Struct
		contentPermission := resourceToPermission(d)

		//Call create permission with our newly populated struct
		err := c.UpdatePermission(contentPermission)

		if err != nil {
			return err
		}

		log.Println("Saving Id to state...")
		d.SetId(contentPermission.ID)
		log.Printf("PermissionId: %s", contentPermission.ID)
	}

	log.Println("====End Permission Create====")

	//After creating an object, we read it again to make sure the state is properly saved
	return resourceSumologicPermissionRead(d, meta)
}

func resourceSumologicPermissionUpdate(d *schema.ResourceData, meta interface{}) error {
	log.Println("====Begin Permission Update====")
	c := meta.(*Client)

	//Load all data from the schema into a Permission Struct
	update := resourceToPermission(d)

	//Update the permission and return any errors
	return c.UpdatePermission(update)
}

func resourceToPermission(d *schema.ResourceData) ContentPermissions {
	log.Println("===Begin population of Permission Struct===")

	var contentPermissions ContentPermissions

	//grab all the permissions objects defined in the resource block
	permissionSet := d.Get("permission").(*schema.Set)

	//set the contentId to the ID for the struct as well
	contentPermissions.ID = d.Get("content_id").(string)

	//loop across all permission objects
	for _, m := range permissionSet.List() {
		var permission Permission

		//convert to a map
		permMap := m.(map[string]interface{})

		//populate the data
		permission.Name = permMap["permission_name"].(string)
		permission.SourceType = permMap["source_type"].(string)
		permission.SourceId = permMap["source_id"].(string)
		permission.ContentId = d.Get("content_id").(string)

		log.Printf("Name: %v", permission.Name)
		log.Printf("SourceType: %v", permission.SourceType)
		log.Printf("SourceId: %v", permission.SourceId)
		log.Printf("ContentId: %v", permission.ContentId)

		//load each permission object into the the final struct
		permissionSlice := contentPermissions.Permissions[:]
		contentPermissions.Permissions = append(permissionSlice, permission)

	}

	log.Printf("Permissions: %v", contentPermissions.Permissions)
	log.Println("===End population of Permission Struct===")

	return contentPermissions
}
