package sumologic

import (
	"encoding/json"
	"fmt"
	"net/url"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceSumologicFolder() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceSumologicFolderRead,
		Schema: map[string]*schema.Schema{
			"path": {
				Type:     schema.TypeString,
				Required: true,
			},
			"id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"name": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func dataSourceSumologicFolderRead(d *schema.ResourceData, meta interface{}) error {
	c := meta.(*Client)

	folder, err := c.GetFolderByPath(d.Get("path").(string))

	if err != nil {
		return err
	}

	d.SetId(folder.ID)
	d.Set("name", folder.Name)

	return nil
}

func (s *Client) GetFolderByPath(path string) (*Folder, error) {
	data, _, err := s.Get(fmt.Sprintf("v2/content/path?path=%s", url.QueryEscape(path)))
	if err != nil {
		return nil, err
	}

	if data == nil {
		return nil, fmt.Errorf("folder with path '%s' does not exist", path)
	}

	var folder Folder
	err = json.Unmarshal(data, &folder)
	if err != nil {
		return nil, err
	}
	if len(folder.ID) == 0 {
		return nil, fmt.Errorf("folder with path '%s' does not exist", path)
	}

	return &folder, nil
}
