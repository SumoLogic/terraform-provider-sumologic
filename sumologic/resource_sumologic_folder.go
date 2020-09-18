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
		//		Importer: &schema.ResourceImporter{
		//			State: resourceSumologicFolderImport,
		//		},

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
	log.Println("====Begin Folder Read====")

	c := meta.(*Client)
	//retrieve the folder Id from the state
	id := d.Id()
	log.Printf("Folder Id from schema: %s", id)

	folder, err := c.GetFolder(id)

	//Error retrieving folder
	if err != nil {
		return err
	}

	//ensure the Folder is populated
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

	log.Println("====End Folder Read====")
	return nil
}

func resourceSumologicFolderDelete(d *schema.ResourceData, meta interface{}) error {
	log.Println("====Begin Folder Delete====")
	log.Printf("Deleting Folder Id: %s", d.Id())
	c := meta.(*Client)
	log.Println("====End Folder Delete====")
	return c.DeleteFolder(d.Id(), d.Timeout(schema.TimeoutDelete))
}

func resourceSumologicFolderCreate(d *schema.ResourceData, meta interface{}) error {
	log.Println("====Begin Folder Create====")
	c := meta.(*Client)

	//If there is no id in the state, then we need to create the object
	if d.Id() == "" {

		//Load all the data we have from the schema into a Folder Struct
		folder := resourceToFolder(d)
		log.Println("Newly populated folder values:")
		log.Printf("ParentId: %s", folder.ParentId)
		log.Printf("Name: %s", folder.Name)
		log.Printf("Description: %s", folder.Description)

		//Call create folder with our newly populated struct
		id, err := c.CreateFolder(folder)

		//Error during CreateFolder
		if err != nil {
			return err
		}

		log.Println("Saving Id to state...")
		d.SetId(id)
		log.Printf("FolderId: %s", id)
	}

	log.Println("====End Folder Create====")

	//After creating an object, we read it again to make sure the state is properly saved
	return resourceSumologicFolderRead(d, meta)
}

func resourceSumologicFolderUpdate(d *schema.ResourceData, meta interface{}) error {
	log.Println("====Begin Folder Update====")

	c := meta.(*Client)

	//Load all data from the schema into a Folder Struct
	folder := resourceToFolder(d)

	log.Printf("Parent Id: %s", folder.ParentId)
	log.Printf("Folder Id: %s", folder.ID)
	log.Printf("Name: %s", folder.Name)
	log.Printf("Description: %s", folder.Description)

	//Update the folder and return any errors
	return c.UpdateFolder(folder)

}

func resourceToFolder(d *schema.ResourceData) Folder {
	log.Println("Loading data from schema to Folder struct...")

	var folder Folder
	folder.ID = d.Id()
	folder.ParentId = d.Get("parent_id").(string)
	folder.Name = d.Get("name").(string)
	folder.Description = d.Get("description").(string)

	return folder
}
