package sumologic

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceSumologicCsoarPlaybook() *schema.Resource {
	return &schema.Resource{
		Create: resourceSumologicCsoarPlaybookCreate,
		Read:   resourceSumologicCsoarPlaybookRead,
		Update: resourceSumologicCsoarPlaybookUpdate,
		Delete: resourceSumologicCsoarPlaybookDelete,
		Importer: &schema.ResourceImporter{
			State: resourceSumologicCsoarPlaybookImport,
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"updated_name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"tags": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"is_deleted": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			"draft": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			"is_published": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			"last_updated": {
				Type:     schema.TypeInt,
				Computed: true, // This should be set by the API, not by user
			},
			"created_by": {
				Type:     schema.TypeInt,
				Computed: true, // This should be set by the API, not by user
			},
			"updated_by": {
				Type:     schema.TypeInt,
				Computed: true, // This should be set by the API, not by user
			},
			"nested": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			"type": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"is_enabled": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  true,
			},
			"links": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeMap,
					Elem: &schema.Schema{
						Type: schema.TypeString,
					},
				},
			},
			"nodes": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeMap,
					Elem: &schema.Schema{
						Type: schema.TypeString,
					},
				},
			},
		},
	}
}

func resourceSumologicCsoarPlaybookDelete(d *schema.ResourceData, meta interface{}) error {
	c := meta.(*Client)
	name := d.Get("name").(string)
	return c.DeletePlaybook(name)
}

func resourceSumologicCsoarPlaybookUpdate(d *schema.ResourceData, meta interface{}) error {
	c := meta.(*Client)

	playbook := Playbook{
		Description: d.Get("description").(string),
		Name:        d.Get("name").(string),
		UpdatedName: d.Get("updated_name").(string),
		Tags:        d.Get("tags").(string),
		IsDeleted:   d.Get("is_deleted").(bool),
		Draft:       d.Get("draft").(bool),
		IsPublished: d.Get("is_published").(bool),
		LastUpdated: int64(d.Get("last_updated").(int)),
		CreatedBy:   int64(d.Get("created_by").(int)),
		UpdatedBy:   int64(d.Get("updated_by").(int)),
		Nested:      d.Get("nested").(bool),
		Type:        d.Get("type").(string),
		IsEnabled:   d.Get("is_enabled").(bool),
	}

	if v, ok := d.Get("links").([]interface{}); ok {
		links := make([]map[string]interface{}, len(v))
		for i, link := range v {
			if linkMap, ok := link.(map[string]interface{}); ok {
				links[i] = linkMap
			}
		}
		playbook.Links = links
	}

	if v, ok := d.Get("nodes").([]interface{}); ok {
		nodes := make([]map[string]interface{}, len(v))
		for i, node := range v {
			if nodeMap, ok := node.(map[string]interface{}); ok {
				nodes[i] = nodeMap
			}
		}
		playbook.Nodes = nodes
	}

	return c.UpdatePlaybook(playbook)
}

func resourceSumologicCsoarPlaybookCreate(d *schema.ResourceData, meta interface{}) error {
	return fmt.Errorf("playbooks cannot be created via Terraform. Please create the playbook in the Sumo Logic UI and then import it using 'terraform import'")
}

func resourceSumologicCsoarPlaybookRead(d *schema.ResourceData, meta interface{}) error {
	if d.Id() == "" {
		return fmt.Errorf("resource ID is empty")
	}

	return nil
}

func resourceSumologicCsoarPlaybookImport(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	playbookName := d.Id()

	if playbookName == "" {
		return nil, fmt.Errorf("import ID (playbook name) cannot be empty")
	}

	d.SetId(playbookName)
	d.Set("name", playbookName)

	return []*schema.ResourceData{d}, nil
}
