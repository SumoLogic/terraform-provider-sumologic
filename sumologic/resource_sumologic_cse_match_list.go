package sumologic

import (
	"fmt"
	"log"
	"regexp"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
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
				Type:     schema.TypeSet,
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

	// Determine whether the target column is defined using its ID or its name
	definedTargetColumnIsId, _ := regexp.MatchString("^-?[0-9]*$", d.Get("target_column").(string))
	definedTargetColumnIsName := !definedTargetColumnIsId

	CSEMatchList, err := c.GetCSEMatchList(id, definedTargetColumnIsName)
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

	CSEMatchListItems, err := c.GetCSEMatchListItemsInMatchList(id)
	if err != nil {
		log.Printf("[WARN] CSE Match List items not found when looking by match list id: %s, err: %v", id, err)
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
			return fmt.Errorf("[ERROR] An error occurred converting resource to match list with id %s, err: %v", d.Id(), err)
		}
		d.SetId(id)

		itemsData := d.Get("items").(*schema.Set).List()
		var items []CSEMatchListItemPost
		for _, data := range itemsData {
			item := resourceToCSEMatchListItem([]interface{}{data})
			items = append(items, item)
		}

		if len(items) > 0 {
			err = c.CreateCSEMatchListItems(items, id)
			if err != nil {
				return fmt.Errorf("[ERROR] An error occurred while adding match list items to match list with id %s, err: %v", id, err)
			}
		}

		createStateConf := &resource.StateChangeConf{
			Target: []string{
				fmt.Sprint(len(items)),
			},
			Refresh: func() (interface{}, string, error) {
				resp, err := c.GetCSEMatchListItemsInMatchList(d.Id())
				if err != nil {
					log.Printf("[ERROR] CSE Match List items not found when looking by match list id %s, err: %v", d.Id(), err)
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
			return fmt.Errorf("[ERROR] error waiting for match list with id %s to be created: %s", d.Id(), err)
		}

	}

	return resourceSumologicCSEMatchListRead(d, meta)
}

func resourceToCSEMatchListItem(data interface{}) CSEMatchListItemPost {
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
	return item
}

func resourceSumologicCSEMatchListUpdate(d *schema.ResourceData, meta interface{}) error {
	CSEMatchListPost, err := resourceToCSEMatchList(d)
	if err != nil {
		return fmt.Errorf("[ERROR] An error occurred converting resource to match list with id %s, err: %v", d.Id(), err)
	}

	c := meta.(*Client)
	if err = c.UpdateCSEMatchList(CSEMatchListPost); err != nil {
		return fmt.Errorf("[ERROR] An error occurred updating match list with id %s, err: %v", d.Id(), err)
	}

	itemsData := d.Get("items").(*schema.Set).List()
	var newItems []CSEMatchListItemPost
	var newItemIds []string
	for _, data := range itemsData {
		item := resourceToCSEMatchListItem([]interface{}{data})
		newItems = append(newItems, item)
		newItemIds = append(newItemIds, item.ID)
	}

	CSEMatchListItems, err := c.GetCSEMatchListItemsInMatchList(d.Id())
	if err != nil {
		return fmt.Errorf("[ERROR] CSE Match List items not found when looking by match list id %s, err: %v", d.Id(), err)
	}

	var deleteItemIds []string
	var updateItemIds []string

	// Compare currently existing match list items with the new items to determine if they should be deleted or updated
	for _, item := range CSEMatchListItems.CSEMatchListItemsGetObjects {
		var oldItemId = item.ID
		if contains(newItemIds, oldItemId) {
			updateItemIds = append(updateItemIds, oldItemId)
		} else {
			deleteItemIds = append(deleteItemIds, oldItemId)
		}
	}

	var updateItems []CSEMatchListItemPost
	var addItems []CSEMatchListItemPost

	// Any new items that are not updates to existing items should be added instead
	for _, newItem := range newItems {
		if contains(updateItemIds, newItem.ID) {
			updateItems = append(updateItems, newItem)
		} else {
			addItems = append(addItems, newItem)
		}
	}

	// Delete old items
	for _, oldItem := range CSEMatchListItems.CSEMatchListItemsGetObjects {
		if contains(deleteItemIds, oldItem.ID) {
			err = c.DeleteCSEMatchListItem(oldItem.ID)
			if err != nil {
				return fmt.Errorf("[ERROR] An error occurred while deleting match list item with id %s, err: %v", oldItem.ID, err)
			}
		}
	}

	// Update old items with new items
	for _, updateItem := range updateItems {
		err = c.UpdateCSEMatchListItem(updateItem)
		if err != nil {
			return fmt.Errorf("[ERROR] An error occurred while updating match list item with id %s, err: %v", updateItem.ID, err)
		}
	}

	//Add new items
	if len(addItems) > 0 {
		err = c.CreateCSEMatchListItems(addItems, d.Id())
		if err != nil {
			return fmt.Errorf("[ERROR] An error occurred while adding match list items to match list with id %s, err: %v", d.Id(), err)
		}
	}

	// Wait for update to finish
	updateStateConf := &resource.StateChangeConf{
		Target: []string{
			fmt.Sprint(len(newItems)),
		},
		Refresh: func() (interface{}, string, error) {
			resp, err := c.GetCSEMatchListItemsInMatchList(d.Id())
			if err != nil {
				log.Printf("[ERROR] CSE Match List items not found when looking by match list id %s, err: %v", d.Id(), err)
				return 0, "", err
			}
			return resp, fmt.Sprint(resp.Total), nil
		},
		Timeout:                   d.Timeout(schema.TimeoutUpdate),
		Delay:                     10 * time.Second,
		MinTimeout:                5 * time.Second,
		ContinuousTargetOccurence: 1,
	}

	_, err = updateStateConf.WaitForState()
	if err != nil {
		return fmt.Errorf("[ERROR] Error waiting for match list with id %s to be updated: %s", d.Id(), err)
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
