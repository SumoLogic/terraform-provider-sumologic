package sumologic

import (
	"log"
	"sort"

	"github.com/hashicorp/terraform/helper/schema"
)

func resourceSumologicUser() *schema.Resource {
	return &schema.Resource{
		Create: resourceSumologicUserCreate,
		Read:   resourceSumologicUserRead,
		Delete: resourceSumologicUserDelete,
		Update: resourceSumologicUserUpdate,
		Exists: resourceSumologicUserExists,
		Importer: &schema.ResourceImporter{
			State: resourceSumologicUserImport,
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

func resourceSumologicUserRead(d *schema.ResourceData, meta interface{}) error {
	c := meta.(*Client)

	id := d.Id()
	user, err := c.GetUser(id)

	if err != nil {
		return err
	}

	if user == nil {
		log.Printf("[WARN] User not found, removing from state: %v - %v", id, err)
		d.SetId("")

		return nil
	}

	d.Set("name", user.Name)
	d.Set("description", user.Description)
	d.Set("filter_predicate", user.FilterPredicate)
	sort.Strings(user.Users)
	d.Set("users", user.Users)
	sort.Strings(user.Capabilities)
	d.Set("capabilities", user.Capabilities)

	return nil
}

func resourceSumologicUserDelete(d *schema.ResourceData, meta interface{}) error {
	c := meta.(*Client)

	if d.Get("destroy").(bool) {
		return c.DeleteUser(d.Id())
	}

	return nil
}

func resourceSumologicUserCreate(d *schema.ResourceData, meta interface{}) error {
	c := meta.(*Client)

	if d.Get("lookup_by_name").(bool) {
		user, err := c.GetUserName(d.Get("name").(string))

		if err != nil {
			return err
		}

		if user != nil {
			d.SetId(user.ID)
		}
	}

	if d.Id() == "" {
		id, err := c.CreateUser(User{
			Name: d.Get("name").(string),
		})

		if err != nil {
			return err
		}

		d.SetId(id)
	}

	return resourceSumologicUserUpdate(d, meta)
}

func resourceSumologicUserUpdate(d *schema.ResourceData, meta interface{}) error {

	user := resourceToUser(d)

	c := meta.(*Client)
	err := c.UpdateUser(user)

	if err != nil {
		return err
	}

	return resourceSumologicUserRead(d, meta)
}

func resourceSumologicUserExists(d *schema.ResourceData, meta interface{}) (bool, error) {
	c := meta.(*Client)

	_, err := c.GetUser(d.Id())

	return err == nil, nil
}

func resourceSumologicUserImport(d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
	if err := resourceSumologicUserRead(d, m); err != nil {
		return nil, err
	}
	return []*schema.ResourceData{d}, nil
}

func resourceToUser(d *schema.ResourceData) User {
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

	return User{
		ID:              d.Id(),
		Name:            d.Get("name").(string),
		Description:     d.Get("description").(string),
		FilterPredicate: d.Get("filter_predicate").(string),
		Users:           users,
		Capabilities:    capabilitiess,
	}
}
