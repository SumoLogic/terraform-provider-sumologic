package sumologic

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceSumologicServiceAccount() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceSumologicServiceAccountRead,

		Schema: map[string]*schema.Schema{
			"id": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ExactlyOneOf: []string{"id", "email"},
			},
			"email": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ExactlyOneOf: []string{"id", "email"},
			},
			"name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"is_active": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"role_ids": {
				Type:     schema.TypeSet,
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
		},
	}
}

func dataSourceSumologicServiceAccountRead(d *schema.ResourceData, meta interface{}) error {
	c := meta.(*Client)

	var account *ServiceAccount
	if id, ok := d.GetOk("id"); ok {
		found, err := c.GetServiceAccount(id.(string))
		if err != nil {
			return err
		}
		if found == nil {
			return fmt.Errorf("service account with id %v not found", id)
		}
		account = found
	} else {
		email := d.Get("email").(string)
		accounts, err := c.ListServiceAccounts()
		if err != nil {
			return err
		}
		for i := range accounts {
			if accounts[i].Email == email {
				account = &accounts[i]
				break
			}
		}
		if account == nil {
			return fmt.Errorf("service account with email %v not found", email)
		}
	}

	d.SetId(account.ID)
	d.Set("email", account.Email)
	d.Set("name", account.Name)
	d.Set("is_active", account.IsActive)
	if err := d.Set("role_ids", account.RoleIds); err != nil {
		return fmt.Errorf("error setting role ids for service account %s: %s", account.ID, err)
	}
	return nil
}
