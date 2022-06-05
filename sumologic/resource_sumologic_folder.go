package sumologic

import (
	"log"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceSumologicFolder() *schema.Resource {
	return &schema.Resource{
		Create: resourceSumologicFolderCreate,
		Read:   resourceSumologicFolderRead,
		Delete: resourceSumologicFolderDelete,
		Update: resourceSumologicFolderUpdate,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"parent_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"description": {
				Type:     schema.TypeString,
				Required: true,
			},
		},
		Timeouts: &schema.ResourceTimeout{
			Delete: schema.DefaultTimeout(time.Minute),
		},
	}
}

func resourceSumologicFolderRead(d *schema.ResourceData, meta interface{}) error {
	c := meta.(*Client)
	// Retrieve the folder Id from the state
	id := d.Id()
	log.Printf("[DEBUG] Folder id from schema: %s", id)

	folder, err := c.GetFolder(id)
	if err != nil {
		return err
	}

	// Ensure the Folder is populated
	if folder == nil {
		log.Printf("Folder not found, removing from state: %v - %v", id, err)
		d.SetId("")
		return nil
	}

	// Write the newly read folder into the schema
	d.Set("parent_id", folder.ParentId)
	d.Set("name", folder.Name)
	d.Set("description", folder.Description)
	d.SetId(folder.ID)

	return nil
}

func resourceSumologicFolderDelete(d *schema.ResourceData, meta interface{}) error {
	c := meta.(*Client)
	log.Printf("[DEBUG] Deleting folder: %s", d.Id())
	return c.DeleteFolder(d.Id(), d.Timeout(schema.TimeoutDelete))
}

func resourceSumologicFolderCreate(d *schema.ResourceData, meta interface{}) error {
	c := meta.(*Client)

	// If there is no id in the state, then we need to create the object
	if d.Id() == "" {
		// Load all the data we have from the schema into a Folder Struct
		folder := resourceToFolder(d)

		id, err := c.CreateFolder(folder)
		if err != nil {
			return err
		}

		d.SetId(id)
	}

	return resourceSumologicFolderRead(d, meta)
}

func resourceSumologicFolderUpdate(d *schema.ResourceData, meta interface{}) error {
	c := meta.(*Client)

	// Load all data from the schema into a Folder Struct
	folder := resourceToFolder(d)

	// Update the folder and return any errors
	return c.UpdateFolder(folder)
}

func resourceToFolder(d *schema.ResourceData) Folder {
	var folder Folder

	folder.ID = d.Id()
	folder.ParentId = d.Get("parent_id").(string)
	folder.Name = d.Get("name").(string)
	folder.Description = d.Get("description").(string)

	return folder
}
