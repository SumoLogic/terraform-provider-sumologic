package sumologic

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/url"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
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
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"view_name": {
							Type:     schema.TypeString,
							Computed: true,
						},

						"view_filter": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
			"security_data_filter": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"selection_type": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"log_analytics_filter": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"audit_data_filter": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"capabilities": {
				Type:        schema.TypeSet,
				Computed:    true,
				Description: "Set of [capabilities](https://help.sumologic.com/docs/manage/users-roles/roles/role-capabilities/) associated with this role",
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
	d.Set("name", roleV2.Name)
	d.Set("audit_data_filter", roleV2.AuditDataFilter)
	d.Set("selection_type", roleV2.SelectionType)
	d.Set("description", roleV2.Description)
	d.Set("security_data_filter", roleV2.SecurityDataFilter)
	d.Set("log_analytics_filter", roleV2.LogAnalyticsFilter)

	if err := d.Set("capabilities", schema.NewSet(schema.HashString, convertStringsToInterfaces(roleV2.Capabilities))); err != nil {
		return fmt.Errorf("error setting capabilities for datasource %s: %s", d.Id(), err)
	}

	if err := d.Set("selected_views", flattenSelectedViews(roleV2.SelectedViews)); err != nil {
		return fmt.Errorf("error setting selected views for datasource %s: %s", d.Id(), err)
	}

	return nil
}

func (s *Client) GetRoleNameV2(name string) (*RoleV2, error) {
	data, err := s.Get(fmt.Sprintf("v2/roles?name=%s", url.QueryEscape(name)))
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

func convertStringsToInterfaces(strs []string) []interface{} {
	result := make([]interface{}, len(strs))
	for i, s := range strs {
		result[i] = s
	}
	return result
}

type RoleResponseV2 struct {
	RolesV2 []RoleV2 `json:"data"`
}
