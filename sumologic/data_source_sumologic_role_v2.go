package sumologic

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/url"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func dataSourceSumologicRoleV2() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceSumologicRoleV2Read,

		Schema: map[string]*schema.Schema{
			"id": {
				Type:     schema.TypeString,
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
			"selected_views": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "List of views which with specific view level filters in accordance to the selectionType chosen.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{

						"view_name": {
							Type:     schema.TypeString,
							Required: true,
						},

						"view_filter": {
							Type:     schema.TypeString,
							Optional: true,
						},
					},
				},
			},
			"security_data_filter": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"selection_type": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"log_analytics_filter": {
				Type:     schema.TypeString,
				Optional: true,
			},

			"audit_data_filter": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"capabilities": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "List of [capabilities](https://help.sumologic.com/docs/manage/users-roles/roles/role-capabilities/) associated with this role",
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
		},
	}
}

func dataSourceSumologicRoleV2Read(d *schema.ResourceData, meta interface{}) error {
	c := meta.(*Client)

	var roleV2 *RoleV2
	var err error
	if rid, ok := d.GetOk("id"); ok {
		id := rid.(string)
		roleV2, err = c.GetRoleV2(id)
		if err != nil {
			return fmt.Errorf("role with id %v not found: %v", id, err)
		}
		if roleV2 == nil {
			return fmt.Errorf("role with id %v not found", id)
		}
	} else {
		if rname, ok := d.GetOk("name"); ok {
			name := rname.(string)
			roleV2, err = c.GetRoleNameV2(name)
			if err != nil {
				return fmt.Errorf("role with name %s not found: %v", name, err)
			}
			if roleV2 == nil {
				return fmt.Errorf("role with name %s not found", name)
			}
		} else {
			return errors.New("please specify either id or name")
		}
	}

	d.SetId(roleV2.ID)
	d.Set("selected_views", roleV2.SelectedViews)
	d.Set("name", roleV2.Name)
	d.Set("audit_data_filter", roleV2.AuditDataFilter)
	d.Set("selection_type", roleV2.SelectionType)
	d.Set("capabilities", roleV2.Capabilities)
	d.Set("description", roleV2.Description)
	d.Set("security_data_filter", roleV2.SecurityDataFilter)
	d.Set("log_analytics_filter", roleV2.LogAnalyticsFilter)

	log.Printf("[DEBUG] data_source_sumologic_role: retrieved %v", roleV2)
	return nil
}

func (s *Client) GetRoleNameV2(name string) (*RoleV2, error) {
	data, _, err := s.Get(fmt.Sprintf("v2/roles?name=%s", url.QueryEscape(name)))
	if err != nil {
		return nil, err
	}

	if data == nil {
		return nil, fmt.Errorf("role with name '%s' does not exist", name)
	}

	var response RoleResponseV2
	err = json.Unmarshal(data, &response)
	if err != nil {
		return nil, err
	}
	if len(response.RolesV2) == 0 {
		return nil, fmt.Errorf("role with name '%s' does not exist", name)
	}

	return &response.RolesV2[0], nil
}

type RoleResponseV2 struct {
	RolesV2 []RoleV2 `json:"data"`
}
