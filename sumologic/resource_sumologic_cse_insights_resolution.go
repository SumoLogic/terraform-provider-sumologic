package sumologic

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"log"
	"strconv"
)

func resourceSumologicCSEInsightsResolution() *schema.Resource {
	return &schema.Resource{
		Create: resourceSumologicCSEInsightsResolutionCreate,
		Read:   resourceSumologicCSEInsightsResolutionRead,
		Delete: resourceSumologicCSEInsightsResolutionDelete,
		Update: resourceSumologicCSEInsightsResolutionUpdate,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"description": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: false,
			},
			"name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: false,
			},
			"parent": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     false,
				ValidateFunc: validation.StringInSlice([]string{"Resolved", "False Positive", "No Action", "Duplicate"}, false),
			},
		},
	}
}

func resourceSumologicCSEInsightsResolutionRead(d *schema.ResourceData, meta interface{}) error {
	c := meta.(*Client)

	var CSEInsightsResolutionGet *CSEInsightsResolutionGet
	id, _ := strconv.Atoi(d.Id())

	CSEInsightsResolutionGet, err := c.GetCSEInsightsResolution(id)
	if err != nil {
		log.Printf("[WARN] CSE Insights Resolution not found when looking by id: %d, err: %v", id, err)

	}

	if CSEInsightsResolutionGet == nil {
		log.Printf("[WARN] CSE Insights Resolution not found, removing from state: %v - %v", id, err)
		d.SetId("")
		return nil
	}

	d.Set("name", CSEInsightsResolutionGet.Name)
	d.Set("description", CSEInsightsResolutionGet.Description)
	d.Set("parent", parentIdToParentName(CSEInsightsResolutionGet.Parent.ID))

	return nil
}

func parentIdToParentName(parentId int) string {

	parentName := ""

	if parentId > 0 {
		if parentId == 1 {
			parentName = "Resolved"
		} else if parentId == 2 {
			parentName = "False Positive"
		} else if parentId == 3 {
			parentName = "No Action"
		} else if parentId == 4 {
			parentName = "Duplicate"
		}
	}
	return parentName
}

func parentNameToParentId(parentName string) int {

	parentId := -1

	if parentName != "" {
		if parentName == "Resolved" {
			parentId = 1
		} else if parentName == "False Positive" {
			parentId = 2
		} else if parentName == "No Action" {
			parentId = 3
		} else if parentName == "Duplicate" {
			parentId = 4
		}
	}
	return parentId
}

func resourceSumologicCSEInsightsResolutionDelete(d *schema.ResourceData, meta interface{}) error {
	c := meta.(*Client)

	id, _ := strconv.Atoi(d.Id())
	return c.DeleteCSEInsightsResolution(id)

}

func resourceSumologicCSEInsightsResolutionCreate(d *schema.ResourceData, meta interface{}) error {
	c := meta.(*Client)

	if d.Id() == "" {
		id, err := c.CreateCSEInsightsResolution(CSEInsightsResolutionPost{
			Name:        d.Get("name").(string),
			Description: d.Get("description").(string),
			ParentId:    parentNameToParentId(d.Get("parent").(string)),
		})

		if err != nil {
			return err
		}
		log.Printf("[INFO] got id: %d", id)
		d.SetId(strconv.Itoa(id))
	}

	return resourceSumologicCSEInsightsResolutionUpdate(d, meta)
}

func resourceSumologicCSEInsightsResolutionUpdate(d *schema.ResourceData, meta interface{}) error {
	CSEInsightsResolutionPost, err := resourceToCSEInsightsResolution(d)
	if err != nil {
		return err
	}

	c := meta.(*Client)
	if err = c.UpdateCSEInsightsResolution(CSEInsightsResolutionPost); err != nil {
		return err
	}

	return resourceSumologicCSEInsightsResolutionRead(d, meta)
}

func resourceToCSEInsightsResolution(d *schema.ResourceData) (CSEInsightsResolutionPost, error) {
	id, err := strconv.Atoi(d.Id())
	if err != nil {
		return CSEInsightsResolutionPost{}, err
	}

	return CSEInsightsResolutionPost{
		ID:          id,
		Description: d.Get("description").(string),
	}, nil
}
