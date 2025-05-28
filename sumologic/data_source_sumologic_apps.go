package sumologic

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net/url"
	"sort"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceSumoLogicApps() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceSumoLogicAppsRead,
		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"author": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"apps": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"uuid":           {Type: schema.TypeString, Computed: true},
						"name":           {Type: schema.TypeString, Computed: true},
						"description":    {Type: schema.TypeString, Computed: true},
						"latest_version": {Type: schema.TypeString, Computed: true},
						"icon":           {Type: schema.TypeString, Computed: true},
						"author":         {Type: schema.TypeString, Computed: true},
						"account_types": {
							Type:     schema.TypeList,
							Computed: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
						"beta":     {Type: schema.TypeBool, Computed: true},
						"installs": {Type: schema.TypeInt, Computed: true},
						"app_type": {Type: schema.TypeString, Computed: true},
						"attributes": {
							Type:     schema.TypeList,
							MaxItems: 1,
							Optional: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"category": {
										Type:     schema.TypeList,
										Optional: true,
										Elem: &schema.Schema{
											Type: schema.TypeString,
										},
									},
									"use_case": {
										Type:     schema.TypeList,
										Optional: true,
										Elem: &schema.Schema{
											Type: schema.TypeString,
										},
									},
									"collection": {
										Type:     schema.TypeList,
										Optional: true,
										Elem: &schema.Schema{
											Type: schema.TypeString,
										},
									},
								},
							},
						},
						"family":              {Type: schema.TypeString, Computed: true},
						"installable":         {Type: schema.TypeBool, Computed: true},
						"show_on_marketplace": {Type: schema.TypeBool, Computed: true},
					},
				},
			},
		},
	}
}

func dataSourceSumoLogicAppsRead(d *schema.ResourceData, meta interface{}) error {
	c := meta.(*Client)

	// Read apps from the API
	id, apps, err := c.getApps(d.Get("name").(string), d.Get("author").(string))
	if err != nil {
		return err
	}

	if err := d.Set("apps", flattenApps(apps)); err != nil {
		return err
	}

	d.SetId(id)

	return nil
}

func (s *Client) getApps(name string, author string) (string, []App, error) {
	// Construct the base URL
	baseURL := "v2/apps"

	// Create url.Values to hold the query parameters
	params := url.Values{}
	if name != "" {
		params.Add("name", name)
	}
	if author != "" {
		params.Add("author", author)
	}

	// Construct the full URL string
	fullURL := baseURL
	if len(params) > 0 {
		fullURL += "?" + params.Encode()
	}

	data, err := s.Get(fullURL)
	if err != nil {
		return "", nil, err
	}

	apps := AppsResponse{}
	err = json.Unmarshal(data, &apps)

	if err != nil {
		return "", nil, err
	}

	// Generate a unique ID for this data source
	id := generateDataSourceId(name, author, apps.Apps)

	return id, apps.Apps, nil
}

func generateDataSourceId(name string, author string, apps []App) string {
	// Start with the filter parameters
	idParts := []string{
		fmt.Sprintf("name:%s", name),
		fmt.Sprintf("author:%s", author),
	}

	// Add a sorted list of app UUIDs
	var uuids []string
	for _, app := range apps {
		uuids = append(uuids, app.UUID)
	}
	sort.Strings(uuids)
	idParts = append(idParts, fmt.Sprintf("apps:%s", strings.Join(uuids, ",")))

	// Join all parts and create a hash
	idString := strings.Join(idParts, "|")
	hash := sha256.Sum256([]byte(idString))
	return hex.EncodeToString(hash[:])
}

func flattenApps(apps []App) []interface{} {
	var flattenedApps []interface{}
	for _, app := range apps {

		internalAttributes := make(map[string]interface{})
		internalAttributes["category"] = app.Attributes.Category
		internalAttributes["use_case"] = app.Attributes.UseCase
		internalAttributes["collection"] = app.Attributes.Collection
		attributes := []interface{}{
			internalAttributes,
		}

		flattenedApp := map[string]interface{}{
			"uuid":                app.UUID,
			"name":                app.Name,
			"description":         app.Description,
			"latest_version":      app.LatestVersion,
			"icon":                app.Icon,
			"author":              app.Author,
			"account_types":       app.AccountTypes,
			"beta":                app.Beta,
			"installs":            app.Installs,
			"app_type":            app.AppType,
			"attributes":          attributes,
			"family":              app.Family,
			"installable":         app.Installable,
			"show_on_marketplace": app.ShowOnMarketplace,
		}
		flattenedApps = append(flattenedApps, flattenedApp)
	}
	return flattenedApps
}

type AppsResponse struct {
	Apps []App `json:"apps"`
}

type App struct {
	UUID          string   `json:"uuid"`
	Name          string   `json:"name"`
	Description   string   `json:"description"`
	LatestVersion string   `json:"latestVersion"`
	Icon          string   `json:"icon"`
	Author        string   `json:"author"`
	AccountTypes  []string `json:"accountTypes"`
	Beta          bool     `json:"beta"`
	Installs      int      `json:"installs"`
	AppType       string   `json:"appType"`
	Attributes    struct {
		Category   []string `json:"category"`
		UseCase    []string `json:"useCase"`
		Collection []string `json:"collection"`
	} `json:"attributes"`
	Family            string `json:"family"`
	Installable       bool   `json:"installable"`
	ShowOnMarketplace bool   `json:"showOnMarketplace"`
}
