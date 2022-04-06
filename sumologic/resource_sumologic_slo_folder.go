package sumologic

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"log"
)

const SLOBaseApiUrl = "v1/slos"

func resourceSumologicSLOLibraryFolder() *schema.Resource {
	return &schema.Resource{
		Create: resourceSumologicSLOLibraryFolderCreate,
		Read:   resourceSumologicSLOLibraryFolderRead,
		Update: resourceSumologicSLOLibraryFolderUpdate,
		Delete: resourceSumologicSLOLibraryFolderDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{

			"version": {
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
			},

			"modified_at": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"is_system": {
				Type:     schema.TypeBool,
				Optional: true,
				Computed: true,
			},

			"content_type": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "Folder",
			},

			"created_by": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"parent_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"is_mutable": {
				Type:     schema.TypeBool,
				Optional: true,
				Computed: true,
			},

			"description": {
				Type:     schema.TypeString,
				Required: true,
			},

			"created_at": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"is_locked": {
				Type:     schema.TypeBool,
				Optional: true,
			},

			"type": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "SlosLibraryFolder",
			},

			"modified_by": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"name": {
				Type:     schema.TypeString,
				Required: true,
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

func resourceSumologicSLOLibraryFolderCreate(d *schema.ResourceData, meta interface{}) error {
	c := meta.(*Client)
	if d.Id() == "" {
		folder := resourceToSLOLibraryFolder(d)
		paramMap := make(map[string]string)
		if folder.ParentID == "" {
			rootFolder, err := c.GetSLOLibraryFolder("root")
			if err != nil {
				return err
			}

			folder.ParentID = rootFolder.ID
		}
		paramMap["parentId"] = folder.ParentID
		SLOFolderID, err := c.CreateSLOLibraryFolder(folder, paramMap)
		if err != nil {
			return err
		}

		d.SetId(SLOFolderID)
	}
	return resourceSumologicSLOLibraryFolderRead(d, meta)
}

func resourceSumologicSLOLibraryFolderRead(d *schema.ResourceData, meta interface{}) error {
	c := meta.(*Client)

	folder, err := c.GetSLOLibraryFolder(d.Id())
	if err != nil {
		return err
	}

	if folder == nil {
		log.Printf("[WARN] SLO Folder not found, removing from state: %v - %v", d.Id(), err)
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

func resourceSumologicSLOLibraryFolderUpdate(d *schema.ResourceData, meta interface{}) error {
	c := meta.(*Client)
	sloFolder := resourceToSLOLibraryFolder(d)
	sloFolder.Type = "SlosLibraryFolderUpdate"
	err := c.UpdateSLOLibraryFolder(sloFolder)
	if err != nil {
		return err
	}
	return resourceSumologicSLOLibraryFolderRead(d, meta)
}

func resourceSumologicSLOLibraryFolderDelete(d *schema.ResourceData, meta interface{}) error {
	c := meta.(*Client)
	sloFolder := resourceToSLOLibraryFolder(d)
	err := c.DeleteSLOLibraryFolder(sloFolder.ID)
	if err != nil {
		return err
	}
	return nil
}

func resourceToSLOLibraryFolder(d *schema.ResourceData) SLOLibraryFolder {
	return SLOLibraryFolder{
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
