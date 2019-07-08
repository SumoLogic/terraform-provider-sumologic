package sumologic

import (
	"log"
	"sort"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/helper/validation"
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
			"first_name": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     false,
				ValidateFunc: validation.NoZeroValues,
			},
			"last_name": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     false,
				ValidateFunc: validation.NoZeroValues,
			},
			"email": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     false,
				ValidateFunc: validation.NoZeroValues,
			},
			"active": {
				Type:     schema.TypeBool,
				Optional: true,
				ForceNew: false,
				Default:  true,
			},
			"role_ids": {
				Type:     schema.TypeList,
				Required: true,
				ForceNew: false,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
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

	d.Set("firstName", user.FirstName)
	d.Set("lastName", user.LastName)
	d.Set("email", user.Email)
	sort.Strings(user.RoleIds)
	d.Set("roleIds", user.RoleIds)

	return nil
}

func resourceSumologicUserDelete(d *schema.ResourceData, meta interface{}) error {
	c := meta.(*Client)
	return c.DeleteUser(d.Id())
}

func resourceSumologicUserCreate(d *schema.ResourceData, meta interface{}) error {
	c := meta.(*Client)

	if d.Id() == "" {
		user := resourceToUser(d)
		id, err := c.CreateUser(user)

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

	user, err := c.GetUser(d.Id())
	if err != nil {
		return false, err
	}

	return user != nil, nil
}

func resourceSumologicUserImport(d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
	if err := resourceSumologicUserRead(d, m); err != nil {
		return nil, err
	}
	return []*schema.ResourceData{d}, nil
}

func resourceToUser(d *schema.ResourceData) User {
	rawRoles := d.Get("role_ids").([]interface{})
	roleIds := make([]string, len(rawRoles))
	for i, v := range rawRoles {
		roleIds[i] = v.(string)
	}

	return User{
		ID:        d.Id(),
		FirstName: d.Get("first_name").(string),
		LastName:  d.Get("last_name").(string),
		Email:     d.Get("email").(string),
		Active:    d.Get("active").(bool),
		RoleIds:   roleIds,
	}
}
