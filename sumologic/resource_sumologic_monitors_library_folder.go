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
			},

			"modified_at": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: false,
			},

			"is_system": {
				Type:     schema.TypeBool,
				Optional: true,
				ForceNew: false,
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
			},

			"parent_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: false,
			},

			"is_mutable": {
				Type:     schema.TypeBool,
				Optional: true,
				ForceNew: false,
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
		monitor := resourceToMonitorsLibraryFolder(d)
		paramMap := make(map[string]string)
		paramMap["parentId"] = monitor.ParentID
		monitorDefinitionID, err := c.CreateMonitorsLibraryFolder(monitor, paramMap)
		if err != nil {
			return err
		}

		d.SetId(monitorDefinitionID)
	}
	return resourceSumologicMonitorsLibraryFolderRead(d, meta)
}

func resourceSumologicMonitorsLibraryFolderRead(d *schema.ResourceData, meta interface{}) error {
	c := meta.(*Client)

	monitor, err := c.MonitorsRead(d.Id())
	if err != nil {
		return err
	}

	if monitor == nil {
		log.Printf("[WARN] Monitor not found, removing from state: %v - %v", d.Id(), err)
		d.SetId("")
		return nil
	}

	d.Set("created_by", monitor.CreatedBy)
	d.Set("name", monitor.Name)
	d.Set("created_at", monitor.CreatedAt)
	d.Set("monitor_type", monitor.MonitorType)
	d.Set("description", monitor.Description)
	d.Set("modified_by", monitor.ModifiedBy)
	d.Set("is_mutable", monitor.IsMutable)
	d.Set("version", monitor.Version)
	// d.Set("type", monitor.Type)
	d.Set("parent_id", monitor.ParentID)
	d.Set("modified_at", monitor.ModifiedAt)
	d.Set("content_type", monitor.ContentType)
	d.Set("is_locked", monitor.IsLocked)
	d.Set("is_system", monitor.IsSystem)

	return nil
}

func resourceSumologicMonitorsLibraryFolderUpdate(d *schema.ResourceData, meta interface{}) error {
	c := meta.(*Client)
	monitor := resourceToMonitorsLibraryFolder(d)
	monitor.Type = "MonitorsLibraryFolderUpdate"
	// monitor.Version = monitor.Version + 1
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
