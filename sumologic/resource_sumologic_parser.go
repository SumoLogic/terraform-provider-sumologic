package sumologic

import (
	"log"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func resourceSumologicParser() *schema.Resource {
	return &schema.Resource{
		Create: resourceSumologicParserCreate,
		Read:   resourceSumologicParserRead,
		Update: resourceSumologicParserUpdate,
		Delete: resourceSumologicParserDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringLenBetween(1, 255),
			},
			"description": {
				Type:         schema.TypeString,
				Optional:     true,
				Default:      "",
				ValidateFunc: validation.StringLenBetween(0, 4096),
			},
			"stanzas": {
				Type:             schema.TypeString,
				Required:         true,
				DiffSuppressFunc: suppressParserWhitespaceDiff,
			},
			"local_stanzas": {
				Type:             schema.TypeString,
				Optional:         true,
				DiffSuppressFunc: suppressParserWhitespaceDiff,
			},
			"model_path": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"sourcetype_path": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"families": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"is_partial": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			"parent_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"version": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"content_type": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"created_at": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"created_by": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"modified_at": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"modified_by": {
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
			"is_mutable": {
				Type:     schema.TypeBool,
				Computed: true,
			},
		},
	}
}

// The API may normalize leading/trailing whitespace in stanzas, which would otherwise cause a perpetual diff.
func suppressParserWhitespaceDiff(_, old, new string, _ *schema.ResourceData) bool {
	return strings.TrimSpace(old) == strings.TrimSpace(new)
}

func resourceToParser(d *schema.ResourceData) ParsersLibraryParser {
	return ParsersLibraryParser{
		ID:             d.Id(),
		Name:           d.Get("name").(string),
		Description:    d.Get("description").(string),
		Stanzas:        d.Get("stanzas").(string),
		LocalStanzas:   d.Get("local_stanzas").(string),
		ModelPath:      d.Get("model_path").(string),
		SourcetypePath: d.Get("sourcetype_path").(string),
		Families:       d.Get("families").(string),
		IsPartial:      d.Get("is_partial").(bool),
		ParentID:       d.Get("parent_id").(string),
		Type:           parsersLibraryParserType,
	}
}

func resourceSumologicParserCreate(d *schema.ResourceData, meta interface{}) error {
	c := meta.(*Client)

	parser := resourceToParser(d)

	if parser.ParentID == "" {
		rootID, err := c.GetParsersLibraryRootID()
		if err != nil {
			return err
		}
		parser.ParentID = rootID
	}

	id, err := c.CreateParsersLibraryParser(parser)
	if err != nil {
		return err
	}
	d.SetId(id)

	return resourceSumologicParserRead(d, meta)
}

func resourceSumologicParserRead(d *schema.ResourceData, meta interface{}) error {
	c := meta.(*Client)

	parser, err := c.GetParsersLibraryParser(d.Id())
	if err != nil {
		return err
	}
	if parser == nil {
		log.Printf("[WARN] Parser not found, removing from state: %v", d.Id())
		d.SetId("")
		return nil
	}

	if err := d.Set("name", parser.Name); err != nil {
		return err
	}
	if err := d.Set("description", parser.Description); err != nil {
		return err
	}
	if err := d.Set("stanzas", parser.Stanzas); err != nil {
		return err
	}
	if err := d.Set("local_stanzas", parser.LocalStanzas); err != nil {
		return err
	}
	if err := d.Set("model_path", parser.ModelPath); err != nil {
		return err
	}
	if err := d.Set("sourcetype_path", parser.SourcetypePath); err != nil {
		return err
	}
	if err := d.Set("families", parser.Families); err != nil {
		return err
	}
	if err := d.Set("is_partial", parser.IsPartial); err != nil {
		return err
	}
	if err := d.Set("parent_id", parser.ParentID); err != nil {
		return err
	}
	if err := d.Set("version", int(parser.Version)); err != nil {
		return err
	}
	if err := d.Set("content_type", parser.ContentType); err != nil {
		return err
	}
	if err := d.Set("created_at", parser.CreatedAt); err != nil {
		return err
	}
	if err := d.Set("created_by", parser.CreatedBy); err != nil {
		return err
	}
	if err := d.Set("modified_at", parser.ModifiedAt); err != nil {
		return err
	}
	if err := d.Set("modified_by", parser.ModifiedBy); err != nil {
		return err
	}
	if err := d.Set("is_locked", parser.IsLocked); err != nil {
		return err
	}
	if err := d.Set("is_system", parser.IsSystem); err != nil {
		return err
	}
	if err := d.Set("is_mutable", parser.IsMutable); err != nil {
		return err
	}
	return nil
}

func resourceSumologicParserUpdate(d *schema.ResourceData, meta interface{}) error {
	c := meta.(*Client)

	// Move first (if requested), since the update must carry the latest version.
	if d.HasChange("parent_id") {
		if newParent := d.Get("parent_id").(string); newParent != "" {
			if err := c.MoveParsersLibraryParser(d.Id(), newParent); err != nil {
				return err
			}
		}
	}

	current, err := c.GetParsersLibraryParser(d.Id())
	if err != nil {
		return err
	}
	if current == nil {
		d.SetId("")
		return nil
	}

	parser := resourceToParser(d)
	parser.Version = current.Version

	if err := c.UpdateParsersLibraryParser(parser); err != nil {
		return err
	}

	return resourceSumologicParserRead(d, meta)
}

func resourceSumologicParserDelete(d *schema.ResourceData, meta interface{}) error {
	c := meta.(*Client)
	return c.DeleteParsersLibraryParser(d.Id())
}
