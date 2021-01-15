package sumologic

import (
	"errors"
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func dataSourceSumologicRole() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceSumologicRoleRead,

		Schema: map[string]*schema.Schema{
			"id": {
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
			},
			"name": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"description": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"filter_predicate": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"capabilities": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
		},
	}
}

func dataSourceSumologicRoleRead(d *schema.ResourceData, meta interface{}) error {
	c := meta.(*Client)

	var role *Role
	var err error
	if rid, ok := d.GetOk("id"); ok {
		id := rid.(string)
		role, err = c.GetRole(id)
		if err != nil {
			return fmt.Errorf("role with id %v not found: %v", id, err)
		}
	} else {
		if rname, ok := d.GetOk("name"); ok {
			name := rname.(string)
			role, err = c.GetRoleName(name)
			if err != nil {
				return fmt.Errorf("role with name %s not found: %v", name, err)
			}
			if role == nil {
				return fmt.Errorf("role with name %s not found", name)
			}
		} else {
			return errors.New("please specify either id or name")
		}
	}

	d.SetId(role.ID)
	d.Set("name", role.Name)
	d.Set("description", role.Description)
	d.Set("filter_predicate", role.FilterPredicate)
	if err := d.Set("capabilities", role.Capabilities); err != nil {
		return fmt.Errorf("error setting capabilities for datasource %s: %s", d.Id(), err)
	}

	log.Printf("[DEBUG] data_source_sumologic_role: retrieved %v", role)
	return nil
}
