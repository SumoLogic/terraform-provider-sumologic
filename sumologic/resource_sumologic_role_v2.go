package sumologic

import (
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceSumologicRoleV2() *schema.Resource {
	return &schema.Resource{
		Create: resourceSumologicRoleV2Create,
		Read:   resourceSumologicRoleV2Read,
		Update: resourceSumologicRoleV2Update,
		Delete: resourceSumologicRoleV2Delete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{

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

			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},

			"selection_type": {
				Type:     schema.TypeString,
				Optional: true,
			},

			"name": {
				Type:     schema.TypeString,
				Required: true,
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

func resourceSumologicRoleV2Read(d *schema.ResourceData, meta interface{}) error {
	c := meta.(*Client)

	id := d.Id()
	roleV2, err := c.GetRoleV2(id)
	if err != nil {
		return err
	}

	if roleV2 == nil {
		log.Printf("[WARN] RoleV2 not found, removing from state: %v - %v", id, err)
		d.SetId("")
		return nil
	}

    if roleV2.SelectionType == "" {
        roleV2.SelectionType = "All"
    }

	d.Set("name", roleV2.Name)
	d.Set("audit_data_filter", roleV2.AuditDataFilter)
	d.Set("selection_type", roleV2.SelectionType)
	d.Set("description", roleV2.Description)
	d.Set("security_data_filter", roleV2.SecurityDataFilter)
	d.Set("log_analytics_filter", roleV2.LogAnalyticsFilter)
	if err := d.Set("capabilities", roleV2.Capabilities); err != nil {
		return fmt.Errorf("error setting capabilities for resource %s: %s", d.Id(), err)
	}
	if err := d.Set("selected_views", flattenSelectedViews(roleV2.SelectedViews)); err != nil {
		return fmt.Errorf("error setting selected views for resource %s: %s", d.Id(), err)
	}

	return nil
}

func resourceSumologicRoleV2Delete(d *schema.ResourceData, meta interface{}) error {
	c := meta.(*Client)

	return c.DeleteRoleV2(d.Id())
}

func resourceSumologicRoleV2Update(d *schema.ResourceData, meta interface{}) error {
	c := meta.(*Client)

	roleV2 := resourceToRoleV2(d)

	if roleV2.SelectionType == "" {
       roleV2.SelectionType = "All"
    }


	id := d.Id()
	retrievedRoleV2, _ := c.GetRole(id)
	roleV2.Users = retrievedRoleV2.Users

	err := c.UpdateRoleV2(roleV2)
	if err != nil {
		return err
	}

	return resourceSumologicRoleV2Read(d, meta)
}

func resourceSumologicRoleV2Create(d *schema.ResourceData, meta interface{}) error {
	c := meta.(*Client)

	if d.Id() == "" {
		roleV2 := resourceToRoleV2(d)
		id, err := c.CreateRoleV2(roleV2)
		if err != nil {
			return err
		}

		d.SetId(id)
	}

	return resourceSumologicRoleV2Read(d, meta)
}

func resourceToRoleV2(d *schema.ResourceData) RoleV2 {

	selectedViewsData := d.Get("selected_views").([]interface{})
	var selectedViews []ViewFilterDefinition
	for _, data := range selectedViewsData {
		selectedViews = append(selectedViews, resourceToViewFilterDefinition([]interface{}{data}))
	}

	capabilitiesData := d.Get("capabilities").([]interface{})
	var capabilities []string
	for _, data := range capabilitiesData {
		capabilities = append(capabilities, data.(string))
	}

	return RoleV2{
		SecurityDataFilter: d.Get("security_data_filter").(string),
		Name:               d.Get("name").(string),
		AuditDataFilter:    d.Get("audit_data_filter").(string),
		ID:                 d.Id(),
		SelectedViews:      selectedViews,
		Description:        d.Get("description").(string),
		SelectionType:      d.Get("selection_type").(string),
		LogAnalyticsFilter: d.Get("log_analytics_filter").(string),
		Capabilities:       capabilities,
		Users:              make([]string, 0),
	}
}

func resourceToViewFilterDefinition(data interface{}) ViewFilterDefinition {

	viewFilterDefinitionSlice := data.([]interface{})
	viewFilterDefinition := ViewFilterDefinition{}
	if len(viewFilterDefinitionSlice) > 0 {
		viewFilterDefinitionObj := viewFilterDefinitionSlice[0].(map[string]interface{})
		viewFilterDefinition.ViewName = viewFilterDefinitionObj["view_name"].(string)
		viewFilterDefinition.ViewFilter = viewFilterDefinitionObj["view_filter"].(string)
	}

	return viewFilterDefinition
}

func flattenSelectedViews(objs []ViewFilterDefinition) []interface{} {
	if objs == nil {
		return nil
	}
	items := []interface{}{}
	for _, item := range objs {
		i := map[string]interface{}{
			"view_name":   item.ViewName,
			"view_filter": item.ViewFilter,
		}
		items = append(items, i)
	}

	return items

}
