package sumologic

import (
	"fmt"
	"log"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceSumologicLocalFileSource() *schema.Resource {
	localFileSource := resourceSumologicSource()
	localFileSource.Create = resourceSumologicLocalFileSourceCreate
	localFileSource.Read = resourceSumologicLocalFileSourceRead
	localFileSource.Update = resourceSumologicLocalFileSourceUpdate
	localFileSource.Importer = &schema.ResourceImporter{
		State: resourceSumologicSourceImport,
	}

	localFileSource.Schema["path_expression"] = &schema.Schema{
		Type:     schema.TypeString,
		Required: true,
	}

	localFileSource.Schema["encoding"] = &schema.Schema{
		Type:     schema.TypeString,
		Optional: true,
		Default:  "UTF-8",
	}

	localFileSource.Schema["deny_list"] = &schema.Schema{
		Type:     schema.TypeSet,
		Optional: true,
		Elem:     &schema.Schema{Type: schema.TypeString},
	}

	return localFileSource
}

func resourceSumologicLocalFileSourceCreate(d *schema.ResourceData, meta interface{}) error {
	c := meta.(*Client)

	if d.Id() == "" {
		source := resourceToLocalFileSource(d)

		id, err := c.CreateLocalFileSource(source, d.Get("collector_id").(int))

		if err != nil {
			return err
		}

		d.SetId(strconv.Itoa(id))
	}

	return resourceSumologicLocalFileSourceRead(d, meta)
}

func resourceSumologicLocalFileSourceUpdate(d *schema.ResourceData, meta interface{}) error {
	c := meta.(*Client)

	source := resourceToLocalFileSource(d)

	err := c.UpdateLocalFileSource(source, d.Get("collector_id").(int))

	if err != nil {
		return err
	}

	return resourceSumologicLocalFileSourceRead(d, meta)
}

func resourceToLocalFileSource(d *schema.ResourceData) LocalFileSource {
	rawDenyList := d.Get("deny_list").(*schema.Set).List()
	var denylist []string
	for _, j := range rawDenyList {
		denylist = append(denylist, j.(string))
	}
	source := resourceToSource(d)
	source.Type = "LocalFile"

	localFileSource := LocalFileSource{
		Source:         source,
		PathExpression: d.Get("path_expression").(string),
		Encoding:       d.Get("encoding").(string),
		DenyList:       denylist,
	}

	return localFileSource
}

func resourceSumologicLocalFileSourceRead(d *schema.ResourceData, meta interface{}) error {
	c := meta.(*Client)

	id, _ := strconv.Atoi(d.Id())
	source, err := c.GetLocalFileSource(d.Get("collector_id").(int), id)

	if err != nil {
		return err
	}

	if source == nil {
		log.Printf("[WARN] LocalFile source not found, removing from state: %v - %v", id, err)
		d.SetId("")

		return nil
	}

	if err := resourceSumologicSourceRead(d, source.Source); err != nil {
		return fmt.Errorf("%s", err)
	}
	d.Set("path_expression", source.PathExpression)
	d.Set("encoding", source.Encoding)
	d.Set("deny_list", source.DenyList)

	return nil
}
