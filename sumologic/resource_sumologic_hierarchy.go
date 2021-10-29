package sumologic

import (
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

var (
	maxDepth            = 6
	maxWidth            = 20
	maxHierarchyNameLen = 256
	maxFilterKeyLen     = 1024
	maxFilterValueLen   = 1024
	maxEntityTypeLen    = 1024
)

func resourceSumologicHierarchy() *schema.Resource {
	return &schema.Resource{
		Create: resourceSumologicHierarchyCreate,
		Read:   resourceSumologicHierarchyRead,
		Update: resourceSumologicHierarchyUpdate,
		Delete: resourceSumologicHierarchyDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringLenBetween(1, maxHierarchyNameLen),
			},

			"filter": {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: getFilterSchema(),
				},
			},

			"level": {
				Type:     schema.TypeList,
				Required: true,
				Elem: &schema.Resource{
					Schema: getLevelSchema(1),
				},
			},
		},
	}
}

func getFilterSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"key": {
			Type:         schema.TypeString,
			Required:     true,
			ValidateFunc: validation.StringLenBetween(1, maxFilterKeyLen),
		},
		"value": {
			Type:         schema.TypeString,
			Required:     true,
			ValidateFunc: validation.StringLenBetween(1, maxFilterValueLen),
		},
	}
}

func getLevelSchema(currentDepth int) map[string]*schema.Schema {
	levelSchema := map[string]*schema.Schema{
		"entity_type": {
			Type:         schema.TypeString,
			Required:     true,
			ValidateFunc: validation.StringLenBetween(1, maxEntityTypeLen),
		},
	}

	var nextLevelsWithConditionsSchema map[string]*schema.Schema
	if currentDepth < maxDepth {
		nextLevelsWithConditionsSchema = map[string]*schema.Schema{
			"next_levels_with_conditions": {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: maxWidth,
				Elem: &schema.Resource{
					Schema: getLevelWithConditionSchema(currentDepth + 1),
				},
			},
		}
	} else {
		nextLevelsWithConditionsSchema = map[string]*schema.Schema{
			"next_levels_with_conditions": {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 0,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
		}
	}
	for k, v := range nextLevelsWithConditionsSchema {
		levelSchema[k] = v
	}

	if currentDepth < maxDepth {
		nextLevelSchema := map[string]*schema.Schema{
			"next_level": {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: getLevelSchema(currentDepth + 1),
				},
			},
		}
		for k, v := range nextLevelSchema {
			levelSchema[k] = v
		}
	}

	return levelSchema
}

func getLevelWithConditionSchema(currentDepth int) map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"condition": {
			Type:     schema.TypeString,
			Required: true,
		},
		"level": {
			Type:     schema.TypeList,
			Required: true,
			MinItems: 1,
			MaxItems: 1,
			Elem: &schema.Resource{
				Schema: getLevelSchema(currentDepth),
			},
		},
	}
}

func resourceSumologicHierarchyCreate(d *schema.ResourceData, meta interface{}) error {
	c := meta.(*Client)

	if d.Id() == "" {
		hierarchy := resourceToHierarchy(d)

		log.Println("=====================================================================")
		log.Printf("Creating hierarchy: %+v\n", hierarchy)
		log.Println("=====================================================================")

		id, err := c.CreateHierarchy(hierarchy)
		if err != nil {
			return err
		}

		d.SetId(id)
	}

	return resourceSumologicHierarchyRead(d, meta)
}

func resourceSumologicHierarchyRead(d *schema.ResourceData, meta interface{}) error {
	c := meta.(*Client)

	id := d.Id()
	hierarchy, err := c.GetHierarchy(id)

	log.Println("=====================================================================")
	log.Printf("Read hierarchy: %+v\n", hierarchy)
	log.Println("=====================================================================")

	if err != nil {
		return err
	}

	if hierarchy == nil {
		log.Printf("[WARN] Hierarchy not found, removing from state: %v - %v", id, err)
		d.SetId("")
		return nil
	}

	d.Set("name", hierarchy.Name)

	filter := getTerraformFilter(hierarchy.Filter)
	if err := d.Set("filter", filter); err != nil {
		return err
	}

	level := getTerraformLevel(&hierarchy.Level)
	if err := d.Set("level", level); err != nil {
		return err
	}

	log.Println("=====================================================================")
	log.Printf("name: %+v\n", d.Get("name"))
	log.Printf("filter: %+v\n", d.Get("filter"))
	log.Printf("level: %+v\n", d.Get("level"))
	log.Println("=====================================================================")

	return nil
}

func resourceSumologicHierarchyDelete(d *schema.ResourceData, meta interface{}) error {
	c := meta.(*Client)

	log.Printf("Deleting dashboard: %+v\n", d.Id())

	return c.DeleteHierarchy(d.Id())
}

func resourceSumologicHierarchyUpdate(d *schema.ResourceData, meta interface{}) error {
	c := meta.(*Client)

	hierarchy := resourceToHierarchy(d)

	log.Println("=====================================================================")
	log.Printf("Updating hierarchy: %+v\n", hierarchy)
	log.Println("=====================================================================")

	err := c.UpdateHierarchy(hierarchy)
	if err != nil {
		return err
	}

	return resourceSumologicHierarchyRead(d, meta)
}

func resourceToHierarchy(d *schema.ResourceData) Hierarchy {
	var filter HierarchyFilteringClause
	if val, ok := d.GetOk("filter"); ok {
		rawFilter := val.([]interface{})[0]
		filter = getFilter(rawFilter.(map[string]interface{}))
	}

	rawLevel := d.Get("level").([]interface{})[0]
	level := getLevel(rawLevel.(map[string]interface{}))

	hierarchy := Hierarchy{}
	hierarchy.Name = d.Get("name").(string)
	hierarchy.ID = d.Id()
	hierarchy.Level = level
	if len(filter.Key) > 0 {
		hierarchy.Filter = &filter
	}

	return hierarchy
}

func getFilter(rawFilter map[string]interface{}) HierarchyFilteringClause {
	return HierarchyFilteringClause{
		Key:   rawFilter["key"].(string),
		Value: rawFilter["value"].(string),
	}
}

func getLevel(rawLevel map[string]interface{}) Level {
	level := Level{}
	level.EntityType = rawLevel["entity_type"].(string)

	rawNextLevelsWithConditions := rawLevel["next_levels_with_conditions"].([]interface{})
	nextLevelsWithConditions := []LevelWithCondition{}
	for _, rawLevelWithCondition := range rawNextLevelsWithConditions {
		levelWithCondition := getLevelWithCondition(rawLevelWithCondition.(map[string]interface{}))
		nextLevelsWithConditions = append(nextLevelsWithConditions, levelWithCondition)
	}
	level.NextLevelsWithConditions = nextLevelsWithConditions

	rawNextLevels := rawLevel["next_level"].([]interface{})
	if len(rawNextLevels) > 0 {
		nextLevel := getLevel(rawNextLevels[0].(map[string]interface{}))
		level.NextLevel = &nextLevel
	}

	return level
}

func getLevelWithCondition(rawLevelWithCondition map[string]interface{}) LevelWithCondition {
	rawLevel := rawLevelWithCondition["level"].([]interface{})[0]
	return LevelWithCondition{
		Condition: rawLevelWithCondition["condition"].(string),
		Level:     getLevel(rawLevel.(map[string]interface{})),
	}
}

func getTerraformFilter(filter *HierarchyFilteringClause) []map[string]interface{} {
	if filter == nil {
		return nil
	}

	tfFilter := []map[string]interface{}{}
	tfFilter = append(tfFilter, make(map[string]interface{}))

	tfFilter[0]["key"] = filter.Key
	tfFilter[0]["value"] = filter.Value

	return tfFilter
}

func getTerraformLevel(level *Level) []map[string]interface{} {
	if level == nil {
		return nil
	}

	tfLevel := []map[string]interface{}{}
	tfLevel = append(tfLevel, make(map[string]interface{}))

	tfLevel[0]["entity_type"] = level.EntityType
	tfLevel[0]["next_levels_with_conditions"] = getTerraformLevelsWithConditions(level.NextLevelsWithConditions)
	tfLevel[0]["next_level"] = getTerraformLevel(level.NextLevel)

	return tfLevel
}

func getTerraformLevelsWithConditions(levelsWithConditions []LevelWithCondition) []map[string]interface{} {
	tfLevelsWithConditions := []map[string]interface{}{}

	for _, levelWithCondition := range levelsWithConditions {
		tfLevelWithCondition := make(map[string]interface{})
		tfLevelWithCondition["condition"] = levelWithCondition.Condition
		tfLevelWithCondition["level"] = getTerraformLevel(&levelWithCondition.Level)
		tfLevelsWithConditions = append(tfLevelsWithConditions, tfLevelWithCondition)
	}

	return tfLevelsWithConditions
}
