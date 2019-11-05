package sumologic

import (
	"github.com/hashicorp/terraform/helper/schema"
	"time"
)

func resourceSumologicFolder() *schema.Resource {
	return &schema.Resource{
		Create: resourceSumologicFolderCreate,
		Read:   resourceSumologicFolderRead,
		Update: resourceSumologicFolderUpdate,
		Delete: resourceSumologicFolderDelete,

		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "",
			},
			"parent_id": {
				Type:     schema.TypeString,
				Required: true,
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
			"item_type": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"permissions": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"children": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:     schema.TypeString,
							Required: true,
						},
						"description": {
							Type:     schema.TypeString,
							Required: true,
						},
						"parent_id": {
							Type:     schema.TypeString,
							Required: true,
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
						"item_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"permissions": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
					},
				},
			},
		},
	}
}

func resourceSumologicFolderCreate(d *schema.ResourceData, meta interface{}) error {
	c := meta.(*Client)

	folder_name := d.Get("name").(string)
	folder_description := d.Get("description").(string)
	folder_parent := d.Get("parent_id").(string)
	folder := FolderCreate{}
	folder.Name = folder_name
	folder.Description = folder_description
	folder.ParentId = folder_parent
	folderResponse, err := c.CreateFolder(folder)
	if err != nil {
		return err
	}
	d.SetId(folderResponse.ID)
	resourceSumologicFolderSetAttr(d, folderResponse)
	return nil
}

func resourceSumologicFolderRead(d *schema.ResourceData, meta interface{}) error {
	c := meta.(*Client)

	folder, err := c.GetFolder(d.Id())
	if err != nil {
		return err
	}
	resourceSumologicFolderSetAttr(d, folder)
	return nil
}

func resourceSumologicFolderSetAttr(d *schema.ResourceData, folder Folder) {
	d.Set("name", folder.Name)
	d.Set("description", folder.Description)
	d.Set("parent_id", folder.ParentId)
	d.Set("created_at", folder.CreatedAt)
	d.Set("created_by", folder.CreatedBy)
	d.Set("modified_at", folder.ModifiedAt)
	d.Set("modified_by", folder.ModifiedBy)
	d.Set("item_type", folder.ItemType)
	d.Set("permissions", folder.Permissions)
	children := make([]map[string]interface{}, len(folder.Children), len(folder.Children))
	for idx, folder_child := range folder.Children {
		child := make(map[string]interface{})
		child["name"] = folder_child.Name
		child["description"] = folder_child.Description
		child["parent_id"] = folder_child.ParentId
		child["item_type"] = folder_child.ItemType
		child["permissions"] = folder_child.Permissions
		child["created_at"] = folder_child.CreatedAt
		child["created_by"] = folder_child.CreatedBy
		child["modified_at"] = folder_child.ModifiedAt
		child["modified_by"] = folder_child.ModifiedBy
		child["id"] = folder_child.ID
		children[idx] = child
	}
	d.Set("children", children)
}

func resourceSumologicFolderUpdate(d *schema.ResourceData, meta interface{}) error {
	c := meta.(*Client)

	folder_name := d.Get("name").(string)
	folder_description := d.Get("description").(string)
	folder := FolderUpdate{}
	folder.Name = folder_name
	folder.Description = folder_description
	folderResponse, err := c.UpdateFolder(d.Id(), folder)
	if err != nil {
		return err
	}
	resourceSumologicFolderSetAttr(d, folderResponse)
	return nil
}

func resourceSumologicFolderDelete(d *schema.ResourceData, meta interface{}) error {
	c := meta.(*Client)
	job_id, err := c.StartDeleteFolder(d.Id())
	if err != nil {
		return err
	}
	job_finished := false
	for !job_finished {
		time.Sleep(time.Second)
		status, err := c.DeleteFolderStatus(d.Id(), job_id)
		if err != nil {
			return nil
		}
		if status == "Success" {
			job_finished = true
		}
	}
	d.SetId("")
	return nil
}
