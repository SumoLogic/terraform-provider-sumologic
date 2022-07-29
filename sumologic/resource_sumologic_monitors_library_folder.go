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
				Default:  "MonitorsLibraryFolder",
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

			"obj_permission": GetCmfFgpObjPermSetSchema(),
		},
	}
}

const fgpTargetType = "monitors"

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

		permStmts, convErr := ResourceToCmfFgpPermStmts(d, monitorDefinitionID)
		if convErr != nil {
			return convErr
		}
		_, fgpErr := c.SetCmfFgp(fgpTargetType, CmfFgpRequest{
			PermissionStatements: permStmts,
		})
		if fgpErr != nil {
			return fgpErr
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

	fgpResponse, fgpGetErr := c.GetCmfFgp(fgpTargetType, folder.ID)
	if fgpGetErr != nil {
		// if FGP endpoint is not enabled (not implemented), we should suppress this error
		suppressedErrorCode := HasErrorCode(fgpGetErr.Error(), []string{"not_implemented_yet", "api_not_enabled"})
		if suppressedErrorCode == "" {
			return fgpGetErr
		} else {
			log.Printf("[WARN] FGP Feature has not been enabled yet. Suppressing \"%s\" error under GetCmfFgp operation.", suppressedErrorCode)
		}
	} else {
		CmfFgpPermStmtsSetToResource(d, fgpResponse.PermissionStatements)
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
	monitorFolder := resourceToMonitorsLibraryFolder(d)
	monitorFolder.Type = "MonitorsLibraryFolderUpdate"
	err := c.UpdateMonitorsLibraryFolder(monitorFolder)
	if err != nil {
		return err
	}

	// converting Reource FGP to Struct
	permStmts, convErr := ResourceToCmfFgpPermStmts(d, monitorFolder.ID)
	if convErr != nil {
		return convErr
	}

	// reading FGP from Backend to reconcile
	fgpGetResponse, fgpGetErr := c.GetCmfFgp(fgpTargetType, monitorFolder.ID)
	if fgpGetErr != nil {
		// if FGP endpoint is not enabled (not implemented) and FGP feature is not used,
		// we should suppress this error
		suppressedErrorCode := HasErrorCode(fgpGetErr.Error(), []string{"not_implemented_yet", "api_not_enabled"})
		if suppressedErrorCode == "" && len(permStmts) == 0 {
			return fgpGetErr
		} else {
			log.Printf("[WARN] FGP Feature has not been enabled yet. Suppressing \"%s\" error under GetCmfFgp operation.", suppressedErrorCode)
		}
	}

	if len(permStmts) > 0 || fgpGetResponse != nil {
		_, fgpSetErr := c.SetCmfFgp(fgpTargetType, CmfFgpRequest{
			PermissionStatements: ReconcileFgpPermStmtsWithEmptyPerms(
				permStmts, fgpGetResponse.PermissionStatements,
			),
		})
		if fgpSetErr != nil {
			return fgpSetErr
		}
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
