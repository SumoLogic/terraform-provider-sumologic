package sumologic

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"log"
	"time"
)

func resourceSumologicCSEMatchList() *schema.Resource {
	return &schema.Resource{
		Create: resourceSumologicCSEMatchListCreate,
		Read:   resourceSumologicCSEMatchListRead,
		Delete: resourceSumologicCSEMatchListDelete,
		Update: resourceSumologicCSEMatchListUpdate,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"default_ttl": {
				Type:     schema.TypeInt,
				Optional: true,
				ForceNew: false,
			},
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
			"target_column": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: false,
			},
			"created": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"created_by": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"last_updated": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"last_updated_by": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"items": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"description": {
							Type:     schema.TypeString,
							Required: true,
							ForceNew: false,
						},
						"expiration": {
							Type:     schema.TypeString,
							Optional: true,
							ForceNew: false,
						},
						"value": {
							Type:     schema.TypeString,
							Required: true,
							ForceNew: false,
						},
					},
				},
			},
		},
	}
}

func resourceSumologicCSEMatchListRead(d *schema.ResourceData, meta interface{}) error {
	c := meta.(*Client)

	var CSEMatchList *CSEMatchListGet
	id := d.Id()

	CSEMatchList, err := c.GetCSEMatchList(id)
	if err != nil {
		log.Printf("[WARN] CSE Match List not found when looking by id: %s, err: %v", id, err)

	}

	if CSEMatchList == nil {
		log.Printf("[WARN] CSE Match List not found, removing from state: %v - %v", id, err)
		d.SetId("")
		return nil
	}

	d.Set("name", CSEMatchList.Name)
	d.Set("default_ttl", CSEMatchList.DefaultTtl)
	d.Set("description", CSEMatchList.Description)
	d.Set("name", CSEMatchList.Name)
	d.Set("target_column", CSEMatchList.TargetColumn)
	d.Set("created", CSEMatchList.Created)
	d.Set("created_by", CSEMatchList.CreatedBy)
	d.Set("last_updated", CSEMatchList.LastUpdated)
	d.Set("last_updated_by", CSEMatchList.LastUpdatedBy)

	//items
	var CSEMatchListItems *CSEMatchListItemsInMatchListGet

	CSEMatchListItems, err2 := c.GetCSEMatchListItemsInMatchList(id)
	if err2 != nil {
		log.Printf("[WARN] CSE Match List items not found when looking by match list id: %s, err: %v", id, err2)
	}
	if CSEMatchListItems == nil {
		d.Set("items", nil)
	} else {
		setItems(d, CSEMatchListItems.CSEMatchListItemsGetObjects)
	}

	return nil
}

func setItems(d *schema.ResourceData, items []CSEMatchListItemGet) {

	var its []map[string]interface{}

	for _, t := range items {
		item := map[string]interface{}{
			"id":          t.ID,
			"description": t.Meta.Description,
			"expiration":  t.Expiration,
			"value":       t.Value,
		}
		its = append(its, item)
	}

	d.Set("items", its)

}

func resourceSumologicCSEMatchListDelete(d *schema.ResourceData, meta interface{}) error {
	c := meta.(*Client)
	err := c.DeleteCSEMatchList(d.Id())
	return err

}

func resourceSumologicCSEMatchListCreate(d *schema.ResourceData, meta interface{}) error {
	c := meta.(*Client)

	if d.Id() == "" {
		id, err := c.CreateCSEMatchList(CSEMatchListPost{
			Active:       true,
			DefaultTtl:   d.Get("default_ttl").(int),
			Description:  d.Get("description").(string),
			Name:         d.Get("name").(string),
			TargetColumn: d.Get("target_column").(string),
		})

		if err != nil {
			return err
		}
		d.SetId(id)

		//Match list items
		itemsData := d.Get("items").([]interface{})
		var items []CSEMatchListItemPost
		for _, data := range itemsData {
			item, _ := resourceToCSEMatchListItem([]interface{}{data})
			items = append(items, item)

		}

		if len(items) > 0 {
			err2 := c.CreateCSEMatchListItems(items, id)
			if err2 != nil {
				log.Printf("[WARN] An error occurred while adding match list items to match list id: %s, err: %v", id, err2)
			}

		}

		createStateConf := &resource.StateChangeConf{
			Target: []string{
				fmt.Sprint(len(items)),
			},
			Refresh: func() (interface{}, string, error) {
				resp, err := c.GetCSEMatchListItemsInMatchList(d.Id())
				if err != nil {
					return 0, "", err
				}
				return resp, fmt.Sprint(resp.Total), nil
			},
			Timeout:                   d.Timeout(schema.TimeoutCreate),
			Delay:                     10 * time.Second,
			MinTimeout:                5 * time.Second,
			ContinuousTargetOccurence: 1,
		}

		_, err = createStateConf.WaitForState()
		if err != nil {
			return fmt.Errorf("error waiting for match list (%s) to be created: %s", d.Id(), err)
		}

	}

	return resourceSumologicCSEMatchListRead(d, meta)
}

func resourceToCSEMatchListItem(data interface{}) (CSEMatchListItemPost, string) {
	itemsSlice := data.([]interface{})
	item := CSEMatchListItemPost{}
	if len(itemsSlice) > 0 {
		itemObj := itemsSlice[0].(map[string]interface{})
		item.ID = itemObj["id"].(string)
		item.Description = itemObj["description"].(string)
		item.Active = true
		item.Expiration = itemObj["expiration"].(string)
		item.Value = itemObj["value"].(string)
	}
	return item, item.ID
}

func resourceSumologicCSEMatchListUpdate(d *schema.ResourceData, meta interface{}) error {
	CSEMatchListPost, err := resourceToCSEMatchList(d)
	if err != nil {
		return err
	}

	c := meta.(*Client)
	if err = c.UpdateCSEMatchList(CSEMatchListPost); err != nil {
		return err
	}

	//Match list items
	itemsData := d.Get("items").([]interface{})
	var itemIds []string
	var items []CSEMatchListItemPost
	for _, data := range itemsData {
		item, id := resourceToCSEMatchListItem([]interface{}{data})
		item.ID = ""
		items = append(items, item)
		itemIds = append(itemIds, id)

	}

	if len(items) > 0 {
		err2 := c.CreateCSEMatchListItems(items, d.Id())
		if err2 != nil {
			log.Printf("[WARN] An error occurred while adding match list items to match list id: %s, err: %v", d.Id(), err2)
		}

	}

	var CSEMatchListItems *CSEMatchListItemsInMatchListGet

	CSEMatchListItems, err2 := c.GetCSEMatchListItemsInMatchList(d.Id())
	if err2 != nil {
		log.Printf("[WARN] CSE Match List items not found when looking by match list id: %s, err: %v", d.Id(), err2)
	}
	if CSEMatchListItems != nil {

		for _, t := range CSEMatchListItems.CSEMatchListItemsGetObjects {
			if !contains(itemIds, t.ID) {
				err3 := c.DeleteCSEMatchListItem(t.ID)
				if err3 != nil {
					log.Printf("[WARN] An error occurred deleting match list item with id: %s, err: %v", t.ID, err3)
				}
			}
		}
	}

	createStateConf := &resource.StateChangeConf{
		Target: []string{
			fmt.Sprint(len(items)),
		},
		Refresh: func() (interface{}, string, error) {
			resp, err := c.GetCSEMatchListItemsInMatchList(d.Id())
			if err != nil {
				return 0, "", err
			}
			return resp, fmt.Sprint(resp.Total), nil
		},
		Timeout:                   d.Timeout(schema.TimeoutUpdate),
		Delay:                     10 * time.Second,
		MinTimeout:                5 * time.Second,
		ContinuousTargetOccurence: 1,
	}

	_, err = createStateConf.WaitForState()
	if err != nil {
		return fmt.Errorf("error waiting for match list (%s) to be updated: %s", d.Id(), err)
	}

	return resourceSumologicCSEMatchListRead(d, meta)
}

func contains(slice []string, item string) bool {
	set := make(map[string]struct{}, len(slice))
	for _, s := range slice {
		set[s] = struct{}{}
	}

	_, ok := set[item]
	return ok
}

func resourceToCSEMatchList(d *schema.ResourceData) (CSEMatchListPost, error) {
	id := d.Id()
	if id == "" {
		return CSEMatchListPost{}, nil
	}

	return CSEMatchListPost{
		ID:           id,
		Active:       true,
		DefaultTtl:   d.Get("default_ttl").(int),
		Description:  d.Get("description").(string),
		Name:         d.Get("name").(string),
		TargetColumn: d.Get("target_column").(string),
	}, nil
}
