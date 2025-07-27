package sumologic

import (
	"fmt"
	"log"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func resourceSumologicLocalWindowsEventLogSource() *schema.Resource {
	LocalWindowsEventLogSource := resourceSumologicSource()
	LocalWindowsEventLogSource.Create = resourceSumologicLocalWindowsEventLogSourceCreate
	LocalWindowsEventLogSource.Read = resourceSumologicLocalWindowsEventLogSourceRead
	LocalWindowsEventLogSource.Update = resourceSumologicLocalWindowsEventLogSourceUpdate
	LocalWindowsEventLogSource.Importer = &schema.ResourceImporter{
		State: resourceSumologicSourceImport,
	}

	// Windows Event Log specific fields
	LocalWindowsEventLogSource.Schema["log_names"] = &schema.Schema{
		Type:        schema.TypeList,
		Required:    true,
		Elem:        &schema.Schema{Type: schema.TypeString},
		Description: "List of Windows log types to collect (e.g., Security, Application, System)",
	}

	LocalWindowsEventLogSource.Schema["render_messages"] = &schema.Schema{
		Type:        schema.TypeBool,
		Optional:    true,
		Default:     true,
		Description: "When using legacy format, indicates if full event messages are collected",
	}

	LocalWindowsEventLogSource.Schema["event_format"] = &schema.Schema{
		Type:         schema.TypeInt,
		Optional:     true,
		Default:      0,
		ValidateFunc: validation.IntInSlice([]int{0, 1}),
		Description:  "0 for legacy format (XML), 1 for JSON format",
	}

	LocalWindowsEventLogSource.Schema["event_message"] = &schema.Schema{
		Type:         schema.TypeInt,
		Optional:     true,
		ValidateFunc: validation.IntInSlice([]int{0, 1, 2}),
		Description:  "0 for complete message, 1 for message title, 2 for metadata only. Required if event_format is 0",
	}

	LocalWindowsEventLogSource.Schema["deny_list"] = &schema.Schema{
		Type:        schema.TypeString,
		Optional:    true,
		Description: "Comma-separated list of event IDs to deny",
	}

	LocalWindowsEventLogSource.Schema["allow_list"] = &schema.Schema{
		Type:        schema.TypeString,
		Optional:    true,
		Description: "Comma-separated list of event IDs to allow",
	}

	return LocalWindowsEventLogSource
}

func resourceSumologicLocalWindowsEventLogSourceCreate(d *schema.ResourceData, meta interface{}) error {
	c := meta.(*Client)

	if d.Id() == "" {
		source := resourceToLocalWindowsEventLogSource(d)
		collectorID := d.Get("collector_id").(int)

		id, err := c.CreateLocalWindowsEventLogSource(source, collectorID)
		if err != nil {
			return err
		}

		d.SetId(strconv.Itoa(id))
	}

	return resourceSumologicLocalWindowsEventLogSourceRead(d, meta)
}

func resourceSumologicLocalWindowsEventLogSourceUpdate(d *schema.ResourceData, meta interface{}) error {
	c := meta.(*Client)

	source := resourceToLocalWindowsEventLogSource(d)

	err := c.UpdateLocalWindowsEventLogSource(source, d.Get("collector_id").(int))

	if err != nil {
		return err
	}

	return resourceSumologicLocalWindowsEventLogSourceRead(d, meta)
}

func resourceToLocalWindowsEventLogSource(d *schema.ResourceData) LocalWindowsEventLogSource {

	source := resourceToSource(d)
	source.Type = "LocalWindowsEventLog"

	LocalWindowsEventLogSource := LocalWindowsEventLogSource{
		Source:         source,
		LogNames:       d.Get("log_names").([]interface{}),
		RenderMessages: d.Get("render_messages").(bool),
		EventFormat:    d.Get("event_format").(int),
	}

	// Handle optional deny_list
	if DenyList, ok := d.GetOk("deny_list"); ok {
		LocalWindowsEventLogSource.DenyList = DenyList.(string)
	}

	// Handle optional allow_list
	if AllowList, ok := d.GetOk("allow_list"); ok {
		LocalWindowsEventLogSource.AllowList = AllowList.(string)
	}

	// Handle optional event_message field
	if eventMessage, ok := d.GetOk("event_message"); ok {
		eventMessageInt := eventMessage.(int)
		LocalWindowsEventLogSource.EventMessage = &eventMessageInt
	}

	return LocalWindowsEventLogSource

}

func resourceSumologicLocalWindowsEventLogSourceRead(d *schema.ResourceData, meta interface{}) error {
	c := meta.(*Client)

	id, _ := strconv.Atoi(d.Id())
	source, err := c.GetLocalWindowsEventLogSource(d.Get("collector_id").(int), id)

	if err != nil {
		return err
	}

	if source == nil {
		log.Printf("[WARN] Local Windows Event Log source not found, removing from state: %v - %v", id, err)
		d.SetId("")
		return nil
	}

	if err := resourceSumologicSourceRead(d, source.Source); err != nil {
		return fmt.Errorf("%s", err)
	}
	d.Set("log_names", source.LogNames)
	d.Set("render_messages", source.RenderMessages)
	d.Set("event_format", source.EventFormat)
	d.Set("deny_list", source.DenyList)
	d.Set("allow_list", source.AllowList)
	d.Set("event_message", source.EventMessage)

	return nil
}
