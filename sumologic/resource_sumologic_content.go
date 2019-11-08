package sumologic

import (
	"encoding/json"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"log"
	"reflect"
)

func resourceSumologicContent() *schema.Resource {
	return &schema.Resource{
		Create: resourceSumologicContentCreate,
		Read:   resourceSumologicContentRead,
		Delete: resourceSumologicContentDelete,
		Update: resourceSumologicContentUpdate,
		//		Importer: &schema.ResourceImporter{
		//			State: resourceSumologicContentImport,
		//		},

		Schema: map[string]*schema.Schema{
			"parent_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"config": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: false,
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					log.Println("====Begin Config Comparison====")
					//convert the json strings to content structs for comparison
					var a Content
					var b Content
					err := json.Unmarshal([]byte(old), &a)
					log.Println("old config:")
					log.Println(a)
					log.Println(err)

					err = json.Unmarshal([]byte(new), &b)
					log.Println("new config:")
					log.Println(b)
					log.Println(err)

					//Set the Children element for each content object to an empty array
					a.Children = []Content{}
					b.Children = []Content{}
					//Make the comparison
					result := reflect.DeepEqual(a, b)
					log.Printf("Equivalent: %t", result)
					log.Println("====End Config Comparison====")
					return result
				},
			},
		},
	}
}

func resourceSumologicContentRead(d *schema.ResourceData, meta interface{}) error {
	log.Println("====Begin Content Read====")

	c := meta.(*Client)
	//retrieve the content Id from the state
	id := d.Id()
	log.Printf("Search for Content Id: %s", id)

	log.Println("Looking up content...")
	content, err := c.GetContent(id)

	//Error retrieving content
	if err != nil {
		return err
	}

	if content == nil {
		log.Printf("[WARN] Content not found, removing from state: %v - %v", id, err)
		d.SetId("")
		return nil
	}

	log.Println("Read Values:")
	log.Printf("ParentId: %s", content.ParentId)
	log.Printf("Config: %s", content.Config)
	log.Printf("Name: %s", content.Name)

	// Write the newly read content object into the schema
	d.Set("config", content.Config)

	log.Println("====End Content Read====")
	return nil
}

func resourceSumologicContentDelete(d *schema.ResourceData, meta interface{}) error {
	log.Println("====Begin Content Delete====")
	log.Printf("Deleting Content Id: %s", d.Id())
	c := meta.(*Client)
	log.Println("====End Content Delete====")
	return c.DeleteContent(d.Id())
}

func resourceSumologicContentCreate(d *schema.ResourceData, meta interface{}) error {
	log.Println("====Begin Content Create====")
	c := meta.(*Client)

	//If there is no id in the state, then we need to create the object
	if d.Id() == "" {

		//Load all the data we have from the schema into a Content Struct
		content := resourceToContent(d)
		log.Println("Newly populated content values:")
		log.Printf("ParentId: %s", content.ParentId)
		log.Printf("Config: %s", content.Config)

		//Call create content with our newly populated struct
		id, err := c.CreateContent(*content)

		//Error during CreateContent
		if err != nil {
			return err
		}

		log.Println("Saving Id to state...")
		d.SetId(id)
		log.Printf("ContentId: %s", id)
		log.Printf("ContentType: %s", content.Type)

	}

	log.Println("====End Content Create====")

	//After creating an object, we read it again to make sure the state is properly saved
	return resourceSumologicContentRead(d, meta)
}

func resourceSumologicContentUpdate(d *schema.ResourceData, meta interface{}) error {
	log.Println("====Begin Content Update====")

	c := meta.(*Client)

	//Load all data from the schema into a Content Struct
	content := resourceToContent(d)

	log.Printf("Parent Id: %s", content.ParentId)
	log.Printf("Content Id: %s", content.ID)

	//Due to API limitations, updating contentobjects means they must be deleted & remade
	log.Printf("Deleting Content with Id: %s", content.ID)
	err := c.DeleteContent(content.ID)

	//error during delete operation
	if err != nil {
		return err
	}

	// reset the Id and remake the object with new config
	d.SetId("")

	log.Println("Remaking Deleted Content...")
	return resourceSumologicContentCreate(d, meta)
}

func resourceToContent(d *schema.ResourceData) *Content {
	log.Println("Loading data from schema to Content struct...")
	var content Content

	_ = json.Unmarshal([]byte(d.Get("config").(string)), &content)

	content.Children = []Content{}
	content.ParentId = d.Get("parent_id").(string)
	content.Config = d.Get("config").(string)
	content.ID = d.Id()

	return &content
}

//func resourceSumologicContentImport(d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
//	if err := resourceSumologicContentRead(d, m); err != nil {
//		return nil, err
//	}
//	return []*schema.ResourceData{d}, nil
//}
