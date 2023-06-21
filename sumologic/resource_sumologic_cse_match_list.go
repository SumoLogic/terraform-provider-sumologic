package sumologic

import (
	"fmt"
	"log"
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

func getCSEMatchListItemsInMatchList(id string, meta interface{}) *CSEMatchListItemsInMatchListGet {
	c := meta.(*Client)

	CSEMatchListItems, err := c.GetCSEMatchListItemsInMatchList(id)
	if err != nil {
		log.Printf("[WARN] CSE Match List items not found when looking by match list id: %s, err: %v", id, err)
	}
	return CSEMatchListItems
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

	CSEMatchListItems := getCSEMatchListItemsInMatchList(id, meta)

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
	println("delete: started resourceSumologicCSEMatchListDelete")

	c := meta.(*Client)
	err := c.DeleteCSEMatchList(d.Id())

	println("delete: finished resourceSumologicCSEMatchListDelete")

	return err

}

func resourceSumologicCSEMatchListCreate(d *schema.ResourceData, meta interface{}) error {
	println("create: starting resourceSumologicCSEMatchListCreate")

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
		//itemsData := d.Get("items").([]interface{})
		itemsData := d.Get("items").(*schema.Set).List()
		var items []CSEMatchListItemPost
		for _, data := range itemsData {
			item := resourceToCSEMatchListItem([]interface{}{data})
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

	println("create: finished resourceSumologicCSEMatchListCreate")

	return resourceSumologicCSEMatchListRead(d, meta)
}

func resourceToCSEMatchListItem(data interface{}) CSEMatchListItemPost {
	itemsSlice := data.([]interface{})
	item := CSEMatchListItemPost{}
	if len(itemsSlice) > 0 {
		itemObj := itemsSlice[0].(map[string]interface{})
		// for _, map := range itemObj {
		// 	println("itemObj id: ", )
		// }
		// b, err := json.Marshal(itemObj)
		// if err != nil {
		// 	log.Fatal(err)
		// }
		// fmt.Println("b: ", string(b))
		item.ID = itemObj["id"].(string)
		item.Description = itemObj["description"].(string)
		item.Active = true
		item.Expiration = itemObj["expiration"].(string)
		item.Value = itemObj["value"].(string)
		// println("resourceToCSEMatchListItem itemid1: " + itemObj["id"].(string))
		// println("resourceToCSEMatchListItem itemid2: " + item.ID)
	}
	//println("resourceToCSEMatchListItem returning id", item.ID)
	return item //, item.ID
}

func resourceSumologicCSEMatchListUpdate(d *schema.ResourceData, meta interface{}) error {
	println("update: starting resourceSumologicCSEMatchListUpdate")
	CSEMatchListPost, err := resourceToCSEMatchList(d)
	if err != nil {
		return err
	}

	c := meta.(*Client)
	if err = c.UpdateCSEMatchList(CSEMatchListPost); err != nil {
		return err
	}

	//Match list items
	//itemsData := d.Get("items").([]interface{})
	itemsData := d.Get("items").(*schema.Set).List()
	var items []CSEMatchListItemPost
	for _, data := range itemsData {
		//println("in itemsData, id is")
		item := resourceToCSEMatchListItem([]interface{}{data}) //TODO remove id return
		// item.ID = ""
		items = append(items, item)
		// println("appending id to itemsID: " + id)
		// itemIds = append(itemIds, id)
	}

	// itemIdsString := ""
	// for i := 0; i < len(itemIds); i++ {
	// 	itemIdsString += itemIds[i] + " "
	// }
	// println("update: itemIds are " + itemIdsString)

	CSEMatchListItems := getCSEMatchListItemsInMatchList(d.Id(), meta)
	var itemIds []string
	for _, item := range CSEMatchListItems.CSEMatchListItemsGetObjects {
		//println("current id: ", item.ID)
		itemIds = append(itemIds, item.ID)
	}

	currentItemCount := CSEMatchListItems.Total
	newItemCount := len(itemsData) + currentItemCount

	println(fmt.Sprintf("current count %d, new count %d", currentItemCount, newItemCount))

	if len(items) > 0 {
		err = c.CreateCSEMatchListItems(items, d.Id())
		if err != nil {
			log.Printf("[WARN] An error occurred while adding match list items to match list id: %s, err: %v", d.Id(), err)
		}

	}

	CSEMatchListItems = getCSEMatchListItemsInMatchList(d.Id(), meta)

	// Waits until all new items have finished being indexed in ES
	println(fmt.Sprintf("total %d, len %d, goal %d", CSEMatchListItems.Total, len(CSEMatchListItems.CSEMatchListItemsGetObjects), newItemCount))
	for CSEMatchListItems.Total < newItemCount {
		println(fmt.Sprintf("update: only got %d/%d items, sleeping for 3 seconds", CSEMatchListItems.Total, newItemCount))
		time.Sleep(3 * time.Second)
		CSEMatchListItems = getCSEMatchListItemsInMatchList(d.Id(), meta)
	}

	println(fmt.Sprintf("update: checking %d items for deletion", len(CSEMatchListItems.CSEMatchListItemsGetObjects)))

	// Delete all the old items
	for _, t := range CSEMatchListItems.CSEMatchListItemsGetObjects {
		//println("update: checking id for deletion: " + t.ID)
		if contains(itemIds, t.ID) {
			//println("update: itemIds contains " + t.ID)
			err = c.DeleteCSEMatchListItem(t.ID)
			if err != nil {
				log.Printf("[WARN] An error occurred deleting match list item with id: %s, err: %v", t.ID, err)
			}
		}
	}

	println("update: finished deleting, now getting items")

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

	println("update: finished resourceSumologicCSEMatchListUpdate")

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
