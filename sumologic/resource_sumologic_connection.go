package sumologic

import (
	"encoding/json"
	"fmt"
	"log"
	"regexp"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

func resourceSumologicConnection() *schema.Resource {
	nameValidation := `^([a-zA-Z0-9 +%\-@.,_()]+)$`
	return &schema.Resource{
		Create: resourceSumologicConnectionCreate,
		Read:   resourceSumologicConnectionRead,
		Delete: resourceSumologicConnectionDelete,
		Update: resourceSumologicConnectionUpdate,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"type": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
				// Only WebhookConnection is supported right now
				ValidateFunc: validation.StringInSlice([]string{"WebhookConnection"}, false),
			},
			"name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: false,
				ValidateFunc: validation.All(
					validation.StringLenBetween(1, 128),
					validation.StringMatch(regexp.MustCompile(nameValidation), fmt.Sprintf("Must match regex %s", nameValidation)),
				),
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					newJSON, _ := normalizeJsonString(new)
					oldJSON, _ := normalizeJsonString(old)
					return newJSON == oldJSON
				},
			},
			"description": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     false,
				ValidateFunc: validation.StringLenBetween(0, 1024),
				Default:      "",
			},
			"url": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     false,
				ValidateFunc: validation.IsURLWithHTTPorHTTPS,
			},
			"headers": {
				Type: schema.TypeMap,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Optional: true,
				ForceNew: false,
				Default:  map[string]interface{}{},
			},
			"custom_headers": {
				Type: schema.TypeMap,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Optional: true,
				ForceNew: false,
				Default:  map[string]interface{}{},
			},
			"default_payload": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     false,
				ValidateFunc: validation.StringIsJSON,
			},
			"webhook_type": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     false,
				ValidateFunc: validation.StringInSlice([]string{"AWSLambda", "Azure", "Datadog", "HipChat", "PagerDuty", "Slack", "Webhook", "NewRelic", "Jira", "Opsgenie"}, false),
				Default:      "Webhook",
			},
		},
	}
}

// Takes a value containing JSON string and passes it through
// the JSON parser to normalize it, returns either a parsing
// error or normalized JSON string.
func normalizeJsonString(jsonString interface{}) (string, error) {
	var j interface{}

	if jsonString == nil || jsonString.(string) == "" {
		return "", nil
	}

	s := jsonString.(string)

	err := json.Unmarshal([]byte(s), &j)
	if err != nil {
		return s, err
	}

	bytes, _ := json.Marshal(j)
	return string(bytes[:]), nil
}

func resourceSumologicConnectionRead(d *schema.ResourceData, meta interface{}) error {
	log.Println("====Begin Connection Read====")

	c := meta.(*Client)
	id := d.Id()
	log.Printf("Connection Id from schema: %s", id)

	connection, err := c.GetConnection(id)
	if err != nil {
		return err
	}

	if connection == nil {
		log.Printf("Connection not found, removing from state: %v - %v", id, err)
		d.SetId("")
		return nil
	}

	// Write the newly read connection into the schema
	d.Set("type", connection.Type)
	d.Set("name", connection.Name)
	d.Set("description", connection.Description)
	d.Set("url", connection.URL)
	if err := d.Set("headers", headersToMap(connection.Headers)); err != nil {
		return fmt.Errorf("error setting headers for resource %s: %s", d.Id(), err)
	}
	if err := d.Set("custom_headers", headersToMap(connection.CustomHeaders)); err != nil {
		return fmt.Errorf("error setting custom headers for resource %s: %s", d.Id(), err)
	}
	d.Set("default_payload", connection.DefaultPayload)
	d.Set("webhook_type", connection.WebhookType)
	d.SetId(connection.ID)

	log.Println("====End Connection Read====")
	return nil
}

func resourceSumologicConnectionDelete(d *schema.ResourceData, meta interface{}) error {
	log.Println("====Begin Connection Delete====")
	connectionType := d.Get("type").(string)
	log.Printf("Deleting Connection Id: %s of type: %s", d.Id(), connectionType)
	c := meta.(*Client)
	err := c.DeleteConnection(d.Id(), connectionType)
	log.Println("====End Connection Delete====")
	return err
}

func resourceSumologicConnectionCreate(d *schema.ResourceData, meta interface{}) error {
	log.Println("====Begin Connection Create====")
	c := meta.(*Client)

	if d.Id() == "" {
		connection := resourceToConnection(d)
		printConnection(connection)

		id, err := c.CreateConnection(connection)

		if err != nil {
			return err
		}

		log.Println("Saving Id to state...")
		d.SetId(id)
		log.Printf("ConnectionId: %s", id)
	}

	log.Println("====End Connection Create====")
	return resourceSumologicConnectionRead(d, meta)
}

func resourceSumologicConnectionUpdate(d *schema.ResourceData, meta interface{}) error {
	log.Println("====Begin Connection Update====")

	c := meta.(*Client)

	connection := resourceToConnection(d)
	printConnection(connection)

	err := c.UpdateConnection(connection)
	log.Println("====End Connection Update====")
	return err
}

func resourceToConnection(d *schema.ResourceData) Connection {
	log.Println("Loading data from schema to Connection struct...")

	var connection Connection
	connection.ID = d.Id()
	connection.Type = d.Get("type").(string)
	connection.Name = d.Get("name").(string)
	connection.Description = d.Get("description").(string)
	connection.URL = d.Get("url").(string)
	connection.Headers = mapToHeaders(d.Get("headers").(map[string]interface{}))
	connection.CustomHeaders = mapToHeaders(d.Get("custom_headers").(map[string]interface{}))
	connection.DefaultPayload = d.Get("default_payload").(string)
	connection.WebhookType = d.Get("webhook_type").(string)

	return connection
}

func mapToHeaders(in map[string]interface{}) []Headers {
	headers := []Headers{}
	for k, v := range in {
		headers = append(headers, Headers{Name: k, Value: v.(string)})
	}

	return headers
}

func headersToMap(in []Headers) map[string]interface{} {
	headerMap := map[string]interface{}{}
	for _, header := range in {
		headerMap[header.Name] = header.Value
	}

	return headerMap
}

func printConnection(connection Connection) {
	log.Println("Connection values:")
	log.Printf("Type: %s", connection.Type)
	log.Printf("Name: %s", connection.Name)
	log.Printf("Description: %s", connection.Description)
	log.Printf("URL: %s", connection.URL)
	log.Printf("Headers: %s", connection.Headers)
	log.Printf("CustomHeaders: %s", connection.CustomHeaders)
	log.Printf("DefaultPayload: %s", connection.DefaultPayload)
	log.Printf("WebhookType: %s", connection.WebhookType)
}
