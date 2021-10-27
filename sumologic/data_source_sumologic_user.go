package sumologic

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/url"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func dataSourceSumologicUser() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceSumologicUserRead,

		Schema: map[string]*schema.Schema{
			"id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"email": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"first_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"last_name": {
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

func dataSourceSumologicUserRead(d *schema.ResourceData, meta interface{}) error {
	c := meta.(*Client)

	var user *User
	var err error
	if uid, ok := d.GetOk("id"); ok {
		id := uid.(string)
		user, err = c.GetUser(id)
		if err != nil {
			return fmt.Errorf("user with id %v not found: %v", id, err)
		}
	} else {
		if uemail, ok := d.GetOk("email"); ok {
			email := uemail.(string)
			user, err = c.GetUserName(email)
			if err != nil {
				return fmt.Errorf("user with email address %s not found: %v", email, err)
			}
			if user == nil {
				return fmt.Errorf("user with email address %s not found", email)
			}
		} else {
			return errors.New("please specify either id or name")
		}
	}

	d.SetId(user.ID)
	d.Set("email", user.Email)
	d.Set("first_name", user.FirstName)
	d.Set("last_name", user.LastName)
	d.Set("is_active", user.IsActive)
	if err := d.Set("role_ids", user.RoleIds); err != nil {
		return fmt.Errorf("error setting role ids for datasource %s: %s", d.Id(), err)
	}

	log.Printf("[DEBUG] data_source_sumologic_role: retrieved %v", user)
	return nil
}

func (s *Client) GetUserName(email string) (*User, error) {
	// where to check the correct api endpoint?
	data, _, err := s.Get(fmt.Sprintf("v1/users?email=%s", url.QueryEscape(email)))
	if err != nil {
		return nil, err
	}

	if data == nil {
		return nil, fmt.Errorf("user with email address '%s' does not exist", email)
	}

	var response UserResponse
	err = json.Unmarshal(data, &response)
	if err != nil {
		return nil, err
	}

	return &response.User, nil
}

// Assume email is unique?
type UserResponse struct {
	User User `json:"user"`
}
