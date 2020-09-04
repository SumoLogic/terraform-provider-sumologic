package sumologic

import (
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceSumologicMonitorsLibraryFolder() *schema.Resource {
	return &schema.Resource{
		Create: resourceSumologicMonitorsLibraryFolderCreate,
		Read:   resourceSumologicMonitorsLibraryFolderRead,
		Update: resourceSumologicMonitorsLibraryFolderUpdate,
		Delete: resourceSumologicMonitorsLibraryFolderDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{

			"version": {
				Type:     schema.TypeInt,
				Optional: true,
				ForceNew: false,
				Computed: true,
			},

			"modified_at": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: false,
				Computed: true,
			},

			"is_system": {
				Type:     schema.TypeBool,
				Optional: true,
				ForceNew: false,
				Computed: true,
			},

			"content_type": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: false,
				Default:  "Folder",
			},

			"created_by": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: false,
				Computed: true,
			},

			"parent_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: false,
				Computed: true,
			},

			"is_mutable": {
				Type:     schema.TypeBool,
				Optional: true,
				ForceNew: false,
				Computed: true,
			},

			"description": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: false,
			},

			"created_at": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: false,
				Computed: true,
			},

			"is_locked": {
				Type:     schema.TypeBool,
				Optional: true,
				ForceNew: false,
			},

			"type": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: false,
				Default:  "MonitorsLibraryFolder",
			},

			"modified_by": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: false,
				Computed: true,
			},

			"name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: false,
			},
			"post_request_map": {
				Type:     schema.TypeMap,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
		},
	}
}

func resourceSumologicMonitorsLibraryFolderCreate(d *schema.ResourceData, meta interface{}) error {
	c := meta.(*Client)
	if d.Id() == "" {
		folder := resourceToMonitorsLibraryFolder(d)
		paramMap := make(map[string]string)
		if folder.ParentID == "" {
			rootFolder, err := c.GetMonitorsLibraryFolder("root")
			if err != nil {
				return err
			}

			folder.ParentID = rootFolder.ID
		}
		paramMap["parentId"] = folder.ParentID
		monitorDefinitionID, err := c.CreateMonitorsLibraryFolder(folder, paramMap)
		if err != nil {
			return err
		}

		d.SetId(monitorDefinitionID)
	}
	return resourceSumologicMonitorsLibraryFolderRead(d, meta)
}

func resourceSumologicMonitorsLibraryFolderRead(d *schema.ResourceData, meta interface{}) error {
	c := meta.(*Client)

	folder, err := c.GetMonitorsLibraryFolder(d.Id())
	if err != nil {
		return err
	}

	if folder == nil {
		log.Printf("[WARN] Monitor Folder not found, removing from state: %v - %v", d.Id(), err)
		d.SetId("")
		return nil
	}

	d.Set("created_by", folder.CreatedBy)
	d.Set("created_at", folder.CreatedAt)
	d.Set("modified_by", folder.ModifiedBy)
	d.Set("is_mutable", folder.IsMutable)
	d.Set("version", folder.Version)
	d.Set("name", folder.Name)
	d.Set("description", folder.Description)
	d.Set("parent_id", folder.ParentID)
	d.Set("modified_at", folder.ModifiedAt)
	d.Set("content_type", folder.ContentType)
	d.Set("is_locked", folder.IsLocked)
	d.Set("is_system", folder.IsSystem)

	return nil
}

func resourceSumologicMonitorsLibraryFolderUpdate(d *schema.ResourceData, meta interface{}) error {
	c := meta.(*Client)
	monitor := resourceToMonitorsLibraryFolder(d)
	monitor.Type = "MonitorsLibraryFolderUpdate"
	err := c.UpdateMonitorsLibraryFolder(monitor)
	if err != nil {
		return err
	}
	return resourceSumologicMonitorsLibraryFolderRead(d, meta)
}

func resourceSumologicMonitorsLibraryFolderDelete(d *schema.ResourceData, meta interface{}) error {
	c := meta.(*Client)
	monitor := resourceToMonitorsLibraryFolder(d)
	err := c.DeleteMonitorsLibraryFolder(monitor.ID)
	if err != nil {
		return err
	}
	return nil
}

func resourceToMonitorsLibraryFolder(d *schema.ResourceData) MonitorsLibraryFolder {
	return MonitorsLibraryFolder{
		CreatedBy:   d.Get("created_by").(string),
		Name:        d.Get("name").(string),
		ID:          d.Id(),
		CreatedAt:   d.Get("created_at").(string),
		Description: d.Get("description").(string),
		ModifiedBy:  d.Get("modified_by").(string),
		IsMutable:   d.Get("is_mutable").(bool),
		Version:     d.Get("version").(int),
		Type:        d.Get("type").(string),
		ParentID:    d.Get("parent_id").(string),
		ModifiedAt:  d.Get("modified_at").(string),
		ContentType: d.Get("content_type").(string),
		IsLocked:    d.Get("is_locked").(bool),
		IsSystem:    d.Get("is_system").(bool),
	}
}
