package sumologic

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"log"
)

func resourceSumologicDataForwardingRule() *schema.Resource {
	return &schema.Resource{

		Create: resourceSumologicDataForwardingRuleCreate,
		Read:   resourceSumologicDataForwardingRuleRead,
		Update: resourceSumologicDataForwardingRuleUpdate,
		Delete: resourceSumologicDataForwardingRuleDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{

			"index_id": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringLenBetween(16, 16),
			},
			"destination_id": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringLenBetween(16, 16),
			},
			"enabled": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"file_prefix": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"payload_schema": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validatePayloadSchema,
			},
			"format": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validateFormat,
			},
		},
	}

}

func resourceSumologicDataForwardingRuleCreate(d *schema.ResourceData, meta interface{}) error {
	c := meta.(*Client)

	if d.Id() == "" {
		dataForwardingRule := resourceToDataForwardingRule(d)
		createdDataForwardingRule, err := c.CreateDataForwardingRule(dataForwardingRule)

		if err != nil {
			return err
		}

		d.SetId(createdDataForwardingRule.IndexId)
	}

	return resourceSumologicDataForwardingRuleRead(d, meta)
}

func resourceSumologicDataForwardingRuleUpdate(d *schema.ResourceData, meta interface{}) error {
	c := meta.(*Client)

	dataForwardingRule := resourceToDataForwardingRule(d)

	err := c.UpdateDataForwardingRule(dataForwardingRule)
	if err != nil {
		return err
	}

	return resourceSumologicDataForwardingRuleRead(d, meta)
}

func resourceSumologicDataForwardingRuleDelete(d *schema.ResourceData, meta interface{}) error {
	c := meta.(*Client)

	return c.DeleteDataForwardingRule(d.Id())
}

func resourceSumologicDataForwardingRuleRead(d *schema.ResourceData, meta interface{}) error {
	c := meta.(*Client)

	indexId := d.Get("index_id").(string)
	dataForwardingRule, err := c.getDataForwardingRule(indexId)

	if err != nil {
		return err
	}

	if dataForwardingRule == nil {
		log.Printf("[WARN] DataForwarding Rule (%s) not found, removing from state", indexId)
		d.SetId("")

		return nil
	}

	d.Set("index_id", indexId)
	d.Set("destination_id", dataForwardingRule.DestinationId)
	d.Set("enabled", dataForwardingRule.Enabled)
	d.Set("file_prefix", dataForwardingRule.FileFormat)
	d.Set("payload_schema", dataForwardingRule.PayloadSchema)
	d.Set("format", dataForwardingRule.Format)

	return nil
}

func validatePayloadSchema(val interface{}, key string) (warns []string, errs []error) {
	allowedValues := map[string]struct{}{
		"builtInFields": {},
		"allFields":     {},
		"raw":           {},
	}

	if strVal, ok := val.(string); ok {
		if _, valid := allowedValues[strVal]; !valid {
			errs = append(errs, fmt.Errorf("%q must be one of %v", key, []string{"builtInFields", "allFields", "raw"}))
		}
	} else {
		errs = append(errs, fmt.Errorf("%q must be a string", key))
	}

	return warns, errs
}

func validateFormat(val interface{}, key string) (warns []string, errs []error) {
	allowedValues := map[string]struct{}{
		"csv":  {},
		"json": {},
		"text": {},
	}

	if strVal, ok := val.(string); ok {
		if _, valid := allowedValues[strVal]; !valid {
			errs = append(errs, fmt.Errorf("%q must be one of %v", key, []string{"csv", "json", "text"}))
		}
	} else {
		errs = append(errs, fmt.Errorf("%q must be a string", key))
	}

	return warns, errs
}

func resourceToDataForwardingRule(d *schema.ResourceData) DataForwardingRule {
	return DataForwardingRule{
		IndexId:       d.Get("index_id").(string),
		DestinationId: d.Get("destination_id").(string),
		Enabled:       d.Get("enabled").(bool),
		FileFormat:    d.Get("file_prefix").(string),
		PayloadSchema: d.Get("payload_schema").(string),
		Format:        d.Get("format").(string),
	}
}
