package sumologic

import (
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

var (
	validArgumentTypes = []string{
		"String",
		"Any",
		"Number",
		"Keyword",
	}
)

func resourceSumologicMacro() *schema.Resource {
	return &schema.Resource{
		Create: resourceSumologicMacroCreate,
		Read:   resourceSumologicMacroRead,
		Delete: resourceSumologicMacroDelete,
		Update: resourceSumologicMacroUpdate,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringLenBetween(1, 128),
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"definition": {
				Type:     schema.TypeString,
				Required: true,
			},
			"enabled": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  true,
			},
			"argument": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: getArgumentsSchema(),
				},
			},
			"argument_validation": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: getArgumentValidationsSchema(),
				},
			},
		},
	}
}

func getArgumentsSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"name": {
			Type:     schema.TypeString,
			Required: true,
		},
		"type": {
			Type:         schema.TypeString,
			Optional:     true,
			ValidateFunc: validation.StringInSlice(validArgumentTypes, false),
		},
	}
}

func getArgumentValidationsSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"eval_expression": {
			Type:     schema.TypeString,
			Required: true,
		},
		"error_message": {
			Type:     schema.TypeString,
			Optional: true,
		},
	}
}

func resourceToMacro(d *schema.ResourceData) Macro {
	var arguments []Argument
	if val, ok := d.GetOk("argument"); ok {
		tfArguments := val.([]interface{})
		for _, tfArgument := range tfArguments {
			argument := getArgument(tfArgument.(map[string]interface{}))
			arguments = append(arguments, argument)
		}
	}

	var argumentValidations []Argumentvalidation
	if val, ok := d.GetOk("argument_validation"); ok {
		tfArgumentValidations := val.([]interface{})
		for _, tfArgumentValidation := range tfArgumentValidations {
			argumentValidation := getArgumentValidation(tfArgumentValidation.(map[string]interface{}))
			argumentValidations = append(argumentValidations, argumentValidation)
		}
	}

	return Macro{
		ID:                  d.Id(),
		Name:                d.Get("name").(string),
		Description:         d.Get("description").(string),
		Definition:          d.Get("definition").(string),
		Enabled:             d.Get("enabled").(bool),
		Arguments:           arguments,
		ArgumentValidations: argumentValidations,
	}
}

func getArgument(tfArgument map[string]interface{}) Argument {
	var argument Argument

	if val, ok := tfArgument["name"]; ok {
		argument.Name = val.(string)
	}

	if val, ok := tfArgument["type"]; ok {
		argument.Type = val.(string)
	}
	return argument
}

func getArgumentValidation(tfArgumentValidation map[string]interface{}) Argumentvalidation {
	var argumentValidation Argumentvalidation

	if val, ok := tfArgumentValidation["eval_expression"]; ok {
		argumentValidation.EvalExpression = val.(string)
	}

	if val, ok := tfArgumentValidation["error_message"]; ok {
		argumentValidation.ErrorMessage = val.(string)
	}
	return argumentValidation
}

func setMacro(d *schema.ResourceData, macro *Macro) error {
	if err := d.Set("name", macro.Name); err != nil {
		return err
	}
	if err := d.Set("description", macro.Description); err != nil {
		return err
	}

	if err := d.Set("definition", macro.Definition); err != nil {
		return err
	}

	if err := d.Set("enabled", macro.Enabled); err != nil {
		return err
	}

	arguments := getTerraformArguments(macro.Arguments)
	if err := d.Set("argument", arguments); err != nil {
		return err
	}

	argumentValidations := getTerraformArgumentValidations(macro.ArgumentValidations)
	if err := d.Set("argument_validation", argumentValidations); err != nil {
		return err
	}

	log.Println("=====================================================================")
	log.Printf("name: %s\n", d.Get("name"))
	log.Printf("description: %s\n", d.Get("description"))
	log.Printf("definition: %s\n", d.Get("definition"))
	log.Printf("enabled: %s\n", d.Get("enabled"))
	log.Printf("arguments: %+v\n", d.Get("arguments"))
	log.Printf("argumentValidations: %+v\n", d.Get("argument_validation"))
	log.Println("=====================================================================")
	return nil
}

func getTerraformArguments(arguments []Argument) []map[string]interface{} {
	tfArguments := make([]map[string]interface{}, len(arguments))

	for i, argument := range arguments {
		tfArguments[i] = make(map[string]interface{})
		tfArguments[i]["name"] = argument.Name
		tfArguments[i]["type"] = argument.Type
	}
	return tfArguments
}

func getTerraformArgumentValidations(argumentValidations []Argumentvalidation) []map[string]interface{} {
	tfArgumentValidations := make([]map[string]interface{}, len(argumentValidations))

	for i, argumentValidation := range argumentValidations {
		tfArgumentValidations[i] = make(map[string]interface{})
		tfArgumentValidations[i]["eval_expression"] = argumentValidation.EvalExpression
		tfArgumentValidations[i]["error_message"] = argumentValidation.ErrorMessage
	}
	return tfArgumentValidations
}

func resourceSumologicMacroCreate(d *schema.ResourceData, meta interface{}) error {
	c := meta.(*Client)
	if d.Id() == "" {
		macro := resourceToMacro(d)
		log.Println("=====================================================================")
		log.Printf("Creating macro: %+v\n", macro)
		log.Println("=====================================================================")

		createdMacro, err := c.CreateMacro(macro)
		if err != nil {
			return err
		}
		d.SetId(createdMacro.ID)
	}

	return resourceSumologicMacroRead(d, meta)
}

func resourceSumologicMacroRead(d *schema.ResourceData, meta interface{}) error {
	c := meta.(*Client)

	id := d.Id()
	macro, err := c.GetMacro(id)
	log.Println("=====================================================================")
	log.Printf("Read macro: %+v\n", macro)
	log.Println("=====================================================================")
	if err != nil {
		return err
	}

	if macro == nil {
		log.Printf("[WARN] Macro not found, removing from state: %v - %v", id, err)
		d.SetId("")
		return nil
	}

	err = setMacro(d, macro)
	return err
}

func resourceSumologicMacroDelete(d *schema.ResourceData, meta interface{}) error {
	c := meta.(*Client)
	log.Printf("Deleting macro: %+v\n", d.Id())
	return c.DeleteMacro(d.Id())
}

func resourceSumologicMacroUpdate(d *schema.ResourceData, meta interface{}) error {
	macro := resourceToMacro(d)
	log.Println("=====================================================================")
	log.Printf("Updating macro: %+v\n", macro)
	log.Println("=====================================================================")

	c := meta.(*Client)
	err := c.UpdateMacro(macro)

	if err != nil {
		return err
	}

	return resourceSumologicMacroRead(d, meta)
}
