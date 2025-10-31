package sumologic

import (
	"encoding/json"
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
				Type:        schema.TypeString,
				Optional:    true,
				Description: "JSON string representation of playbook links",
			},
			"nodes": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "JSON string representation of playbook nodes",
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
		Name: d.Get("name").(string),
	}

	playbook.Description = d.Get("description").(string)
	playbook.UpdatedName = d.Get("updated_name").(string)
	playbook.Tags = d.Get("tags").(string)
	playbook.Type = d.Get("type").(string)
	playbook.IsDeleted = d.Get("is_deleted").(bool)
	playbook.Draft = d.Get("draft").(bool)
	playbook.IsPublished = d.Get("is_published").(bool)
	playbook.Nested = d.Get("nested").(bool)
	playbook.IsEnabled = d.Get("is_enabled").(bool)

	// Handle links as JSON string
	if linksJSON, ok := d.Get("links").(string); ok && linksJSON != "" {
		var links []map[string]interface{}
		if err := json.Unmarshal([]byte(linksJSON), &links); err != nil {
			return fmt.Errorf("error parsing links JSON: %v", err)
		}
		playbook.Links = links
	}

	// Handle nodes as JSON string
	if nodesJSON, ok := d.Get("nodes").(string); ok && nodesJSON != "" {
		var nodes []map[string]interface{}
		if err := json.Unmarshal([]byte(nodesJSON), &nodes); err != nil {
			return fmt.Errorf("error parsing nodes JSON: %v", err)
		}
		playbook.Nodes = nodes
	}

	err := c.UpdatePlaybook(playbook)
	if err != nil {
		return err
	}

	// Always set state after successful update
	d.Set("description", playbook.Description)
	d.Set("tags", playbook.Tags)
	d.Set("updated_name", playbook.UpdatedName)
	d.Set("type", playbook.Type)
	d.Set("is_deleted", playbook.IsDeleted)
	d.Set("draft", playbook.Draft)
	d.Set("is_published", playbook.IsPublished)
	d.Set("nested", playbook.Nested)
	d.Set("is_enabled", playbook.IsEnabled)

	linksJSON := d.Get("links").(string)
	d.Set("links", linksJSON)

	nodesJSON := d.Get("nodes").(string)
	d.Set("nodes", nodesJSON)

	return nil
}

func resourceSumologicCsoarPlaybookCreate(d *schema.ResourceData, meta interface{}) error {
	return fmt.Errorf("playbooks cannot be created via Terraform. Please create the playbook in the CSOAR UI, export it as JSON, and then import it using 'terraform import'")
}

func resourceSumologicCsoarPlaybookRead(d *schema.ResourceData, meta interface{}) error {
	if d.Id() == "" {
		return fmt.Errorf("resource ID is empty")
	}

	// For import-only resources, ensure the name matches the ID
	d.Set("name", d.Id())

	// For import-only resources, preserve existing state values or set reasonable defaults
	// This prevents Terraform from thinking the resource doesn't exist
	if d.Get("description") == nil {
		d.Set("description", "")
	}
	if d.Get("tags") == nil {
		d.Set("tags", "")
	}
	if d.Get("is_deleted") == nil {
		d.Set("is_deleted", false)
	}
	if d.Get("draft") == nil {
		d.Set("draft", false)
	}
	if d.Get("is_published") == nil {
		d.Set("is_published", true)
	}
	if d.Get("nested") == nil {
		d.Set("nested", false)
	}
	if d.Get("type") == nil {
		d.Set("type", "General")
	}
	if d.Get("is_enabled") == nil {
		d.Set("is_enabled", true)
	}
	if d.Get("nodes") == nil {
		d.Set("nodes", "[]")
	}
	if d.Get("links") == nil {
		d.Set("links", "[]")
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

	// Set default values for all schema fields to prevent Terraform from thinking this is a new resource
	d.Set("description", "")
	d.Set("tags", "")
	d.Set("is_deleted", false)
	d.Set("draft", false)
	d.Set("is_published", true)
	d.Set("nested", false)
	d.Set("type", "General")
	d.Set("is_enabled", true)
	d.Set("nodes", "[]")
	d.Set("links", "[]")

	return []*schema.ResourceData{d}, nil
}
