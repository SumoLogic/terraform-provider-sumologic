package sumologic

import (
	"log"

	"github.com/hashicorp/terraform/helper/schema"
)

func resourceSumologicRole() *schema.Resource {
	return &schema.Resource{
		Create: resourceSumologicRoleCreate,
		Read:   resourceSumologicRoleRead,
		Delete: resourceSumologicRoleDelete,
		Update: resourceSumologicRoleUpdate,
		Exists: resourceSumologicRoleExists,
		Importer: &schema.ResourceImporter{
			State: resourceSumologicRoleImport,
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: false,
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: false,
				Default:  "",
			},
			"filter_predicate": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: false,
				Default:  "",
			},
			"users": {
				Type:     schema.TypeList,
				Optional: true,
				ForceNew: false,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"capabilities": {
				Type:     schema.TypeList,
				Optional: true,
				ForceNew: false,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"lookup_by_name": {
				Type:     schema.TypeBool,
				Optional: true,
				ForceNew: false,
				Default:  false,
			},
			"destroy": {
				Type:     schema.TypeBool,
				Optional: true,
				ForceNew: false,
				Default:  true,
			},
		},
	}
}

func resourceSumologicRoleRead(d *schema.ResourceData, meta interface{}) error {
	c := meta.(*Client)

	id := d.Id()
	role, err := c.GetRole(id)

	if err != nil {
		return err
	}

	if role == nil {
		log.Printf("[WARN] Role not found, removing from state: %v - %v", id, err)
		d.SetId("")

		return nil
	}

	d.Set("name", role.Name)
	d.Set("description", role.Description)
	d.Set("filter_predicate", role.FilterPredicate)
	d.Set("users", role.Users)
	d.Set("capabilities", role.Capabilities)

	return nil
}

func resourceSumologicRoleDelete(d *schema.ResourceData, meta interface{}) error {
	c := meta.(*Client)

	if d.Get("destroy").(bool) {
		return c.DeleteRole(d.Id())
	}

	return nil
}

func resourceSumologicRoleCreate(d *schema.ResourceData, meta interface{}) error {
	c := meta.(*Client)

	if d.Get("lookup_by_name").(bool) {
		role, err := c.GetRoleName(d.Get("name").(string))

		if err != nil {
			return err
		}

		if role != nil {
			d.SetId(role.ID)
		}
	}

	if d.Id() == "" {
		id, err := c.CreateRole(Role{
			Name: d.Get("name").(string),
		})

		if err != nil {
			return err
		}

		d.SetId(id)
	}

	return resourceSumologicRoleUpdate(d, meta)
}

func resourceSumologicRoleUpdate(d *schema.ResourceData, meta interface{}) error {

	role := resourceToRole(d)

	c := meta.(*Client)
	err := c.UpdateRole(role)

	if err != nil {
		return err
	}

	return resourceSumologicRoleRead(d, meta)
}

func resourceSumologicRoleExists(d *schema.ResourceData, meta interface{}) (bool, error) {
	c := meta.(*Client)

	_, err := c.GetRole(d.Id())

	return err == nil, nil
}

func resourceSumologicRoleImport(d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
	if err := resourceSumologicRoleRead(d, m); err != nil {
		return nil, err
	}
	return []*schema.ResourceData{d}, nil
}

func resourceToRole(d *schema.ResourceData) Role {
	rawUsers := d.Get("users").([]interface{})
	users := make([]string, len(rawUsers))
	for i, v := range rawUsers {
		users[i] = v.(string)
	}

	rawCapabilities := d.Get("capabilities").([]interface{})
	capabilitiess := make([]string, len(rawCapabilities))
	for i, v := range rawCapabilities {
		capabilitiess[i] = v.(string)
	}

	return Role{
		ID:              d.Id(),
		Name:            d.Get("name").(string),
		Description:     d.Get("description").(string),
		FilterPredicate: d.Get("filter_predicate").(string),
		Users:           users,
		Capabilities:    capabilitiess,
	}
}
