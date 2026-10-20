package sumologic

import (
	"encoding/json"
	"fmt"
	"net/url"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

const parsersLibraryFolderContentType = "Folder"
const parsersSearchPageSize = 100

func dataSourceSumologicParserFolder() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceSumologicParserFolderRead,
		Schema: map[string]*schema.Schema{
			"id": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ExactlyOneOf: []string{"id", "name"},
			},
			"name": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ExactlyOneOf: []string{"id", "name"},
			},
			// Only used to disambiguate when looking up by name.
			"parent_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"description": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"is_locked": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"is_system": {
				Type:     schema.TypeBool,
				Computed: true,
			},
		},
	}
}

func dataSourceSumologicParserFolderRead(d *schema.ResourceData, meta interface{}) error {
	c := meta.(*Client)

	var folder *ParsersLibraryFolder
	var err error

	if id, ok := d.GetOk("id"); ok {
		folder, err = c.GetParsersLibraryFolder(id.(string))
		if err != nil {
			return err
		}
		if folder == nil {
			return fmt.Errorf("parser folder with id '%s' does not exist", id.(string))
		}
	} else {
		name := d.Get("name").(string)
		parentID := d.Get("parent_id").(string)

		folder, err = c.GetParsersLibraryFolderByName(name, parentID)
		if err != nil {
			return err
		}
	}

	d.SetId(folder.ID)
	if err := d.Set("name", folder.Name); err != nil {
		return err
	}
	if err := d.Set("parent_id", folder.ParentID); err != nil {
		return err
	}
	if err := d.Set("description", folder.Description); err != nil {
		return err
	}
	if err := d.Set("is_locked", folder.IsLocked); err != nil {
		return err
	}
	return d.Set("is_system", folder.IsSystem)
}

// ---------- ENDPOINTS ----------

// GetParsersLibraryFolder reads a folder by id. Returns nil if it does not exist.
// An error is returned if the id refers to an item that is not a folder.
func (s *Client) GetParsersLibraryFolder(id string) (*ParsersLibraryFolder, error) {
	data, err := s.Get(fmt.Sprintf("%s/%s", parsersLibraryBaseApiUrl, id))
	if err != nil {
		return nil, err
	}
	if data == nil {
		return nil, nil
	}

	var folder ParsersLibraryFolder
	if err := json.Unmarshal(data, &folder); err != nil {
		return nil, err
	}
	if folder.ContentType != parsersLibraryFolderContentType {
		return nil, fmt.Errorf("parsers library item '%s' is not a folder (contentType: %s)", id, folder.ContentType)
	}
	return &folder, nil
}

// GetParsersLibraryFolderByName finds a folder by exact name. If parentID is not empty, only
// folders directly under that parent are considered. It is an error if zero or more than one
// folder matches.
func (s *Client) GetParsersLibraryFolderByName(name, parentID string) (*ParsersLibraryFolder, error) {
	query := url.QueryEscape("type:folder " + name)

	var matches []ParsersLibraryFolder
	for offset := 0; ; offset += parsersSearchPageSize {
		data, err := s.Get(fmt.Sprintf("%s/search?query=%s&limit=%d&offset=%d",
			parsersLibraryBaseApiUrl, query, parsersSearchPageSize, offset))
		if err != nil {
			return nil, err
		}
		if data == nil {
			break
		}

		var page []ParsersLibraryItemWithPath
		if err := json.Unmarshal(data, &page); err != nil {
			return nil, err
		}

		for _, result := range page {
			item := result.Item
			if item.ContentType != parsersLibraryFolderContentType || item.Name != name {
				continue
			}
			if parentID != "" && item.ParentID != parentID {
				continue
			}
			matches = append(matches, item)
		}

		if len(page) < parsersSearchPageSize {
			break
		}
	}

	switch len(matches) {
	case 0:
		return nil, fmt.Errorf("parser folder with name '%s' does not exist", name)
	case 1:
		return &matches[0], nil
	default:
		return nil, fmt.Errorf("found %d parser folders named '%s'; set parent_id or look up by id", len(matches), name)
	}
}

// ---------- TYPES ----------

type ParsersLibraryFolder struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	ParentID    string `json:"parentId"`
	ContentType string `json:"contentType"`
	IsLocked    bool   `json:"isLocked"`
	IsSystem    bool   `json:"isSystem"`
}

type ParsersLibraryItemWithPath struct {
	Item ParsersLibraryFolder `json:"item"`
	Path string               `json:"path"`
}
