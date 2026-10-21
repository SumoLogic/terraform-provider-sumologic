package sumologic

import (
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func resourceSumologicParserFolder() *schema.Resource {
	return &schema.Resource{
		Create: resourceSumologicParserFolderCreate,
		Read:   resourceSumologicParserFolderRead,
		Update: resourceSumologicParserFolderUpdate,
		Delete: resourceSumologicParserFolderDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringLenBetween(1, 255),
			},
			"description": {
				Type:         schema.TypeString,
				Optional:     true,
				Default:      "",
				ValidateFunc: validation.StringLenBetween(0, 4096),
			},
			"parent_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"version": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"content_type": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"created_at": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"created_by": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"modified_at": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"modified_by": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"is_locked": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"is_system": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"is_mutable": {
				Type:     schema.TypeBool,
				Computed: true,
			},
		},
	}
}

func resourceToParserFolder(d *schema.ResourceData) ParsersLibraryFolder {
	return ParsersLibraryFolder{
		ID:          d.Id(),
		Name:        d.Get("name").(string),
		Description: d.Get("description").(string),
		ParentID:    d.Get("parent_id").(string),
		Type:        parsersLibraryFolderType,
	}
}

func resourceSumologicParserFolderCreate(d *schema.ResourceData, meta interface{}) error {
	c := meta.(*Client)

	if d.Id() == "" {
		folder := resourceToParserFolder(d)

		if folder.ParentID == "" {
			rootID, err := c.GetParsersLibraryRootID()
			if err != nil {
				return err
			}
			folder.ParentID = rootID
		}

		id, err := c.CreateParsersLibraryFolder(folder)
		if err != nil {
			return err
		}
		d.SetId(id)
	}

	return resourceSumologicParserFolderRead(d, meta)
}

func resourceSumologicParserFolderRead(d *schema.ResourceData, meta interface{}) error {
	c := meta.(*Client)

	folder, err := c.GetParsersLibraryFolder(d.Id())
	if err != nil {
		return err
	}
	if folder == nil {
		log.Printf("[WARN] Parser folder not found, removing from state: %v", d.Id())
		d.SetId("")
		return nil
	}

	values := map[string]interface{}{
		"name":         folder.Name,
		"description":  folder.Description,
		"parent_id":    folder.ParentID,
		"version":      int(folder.Version),
		"content_type": folder.ContentType,
		"created_at":   folder.CreatedAt,
		"created_by":   folder.CreatedBy,
		"modified_at":  folder.ModifiedAt,
		"modified_by":  folder.ModifiedBy,
		"is_locked":    folder.IsLocked,
		"is_system":    folder.IsSystem,
		"is_mutable":   folder.IsMutable,
	}
	for k, v := range values {
		if err := d.Set(k, v); err != nil {
			return err
		}
	}
	return nil
}

func resourceSumologicParserFolderUpdate(d *schema.ResourceData, meta interface{}) error {
	c := meta.(*Client)

	// Move first (if requested), since the update must carry the latest version.
	if d.HasChange("parent_id") {
		if newParent := d.Get("parent_id").(string); newParent != "" {
			if err := c.MoveParsersLibraryFolder(d.Id(), newParent); err != nil {
				return err
			}
		}
	}

	current, err := c.GetParsersLibraryFolder(d.Id())
	if err != nil {
		return err
	}
	if current == nil {
		d.SetId("")
		return nil
	}

	folder := resourceToParserFolder(d)
	folder.Version = current.Version

	if err := c.UpdateParsersLibraryFolder(folder); err != nil {
		return err
	}

	return resourceSumologicParserFolderRead(d, meta)
}

func resourceSumologicParserFolderDelete(d *schema.ResourceData, meta interface{}) error {
	c := meta.(*Client)
	return c.DeleteParsersLibraryFolder(d.Id())
}
