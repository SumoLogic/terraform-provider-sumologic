package sumologic

import (
	"fmt"
	"strconv"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

func getRandomizedParamsForHierarchy() (string, HierarchyFilteringClause, Level) {
	testNextLevel := Level{
		EntityType:               "node",
		NextLevelsWithConditions: []LevelWithCondition{},
	}

	name := acctest.RandomWithPrefix("tf-acc-test")
	filter := HierarchyFilteringClause{
		Key:   acctest.RandomWithPrefix("tf-acc-test"),
		Value: acctest.RandomWithPrefix("tf-acc-test"),
	}
	level := Level{
		EntityType: "cluster",
		NextLevelsWithConditions: []LevelWithCondition{
			{
				Condition: acctest.RandomWithPrefix("tf-acc-test"),
				Level: Level{
					EntityType:               "namespace",
					NextLevelsWithConditions: []LevelWithCondition{},
				},
			},
		},
		NextLevel: &testNextLevel,
	}

	return name, filter, level
}

func TestAccSumologicHierarchy_basic(t *testing.T) {
	var hierarchy Hierarchy

	testName, testFilter, testLevel := getRandomizedParamsForHierarchy()

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckHierarchyDestroy(hierarchy),
		Steps: []resource.TestStep{
			{
				Config: testAccCheckSumologicHierarchyConfigImported(testName, testFilter, testLevel),
			},
			{
				ResourceName:      "sumologic_hierarchy.test",
				ImportState:       true,
				ImportStateVerify: false,
			},
		},
	})
}

func TestAccSumologicHierarchy_create(t *testing.T) {
	var hierarchy Hierarchy
	testName, testFilter, testLevel := getRandomizedParamsForHierarchy()

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckHierarchyDestroy(hierarchy),
		Steps: []resource.TestStep{
			{
				Config: testAccSumologicHierarchy(testName, testFilter, testLevel),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckHierarchyExists("sumologic_hierarchy.test", &hierarchy, t),
					resource.TestCheckResourceAttrSet("sumologic_hierarchy.test", "id"),
					resource.TestCheckResourceAttr("sumologic_hierarchy.test", "name", testName),
					resource.TestCheckResourceAttr("sumologic_hierarchy.test", "filter.0.key", testFilter.Key),
					resource.TestCheckResourceAttr("sumologic_hierarchy.test", "filter.0.value", testFilter.Value),
					resource.TestCheckResourceAttr("sumologic_hierarchy.test", "level.0.entity_type",
						testLevel.EntityType),
					resource.TestCheckResourceAttr("sumologic_hierarchy.test",
						"level.0.next_levels_with_conditions.0.condition",
						testLevel.NextLevelsWithConditions[0].Condition),
					resource.TestCheckResourceAttr("sumologic_hierarchy.test",
						"level.0.next_levels_with_conditions.0.level.0.entity_type",
						testLevel.NextLevelsWithConditions[0].Level.EntityType),
					resource.TestCheckResourceAttr("sumologic_hierarchy.test",
						"level.0.next_level.0.entity_type",
						testLevel.NextLevel.EntityType),
				),
			},
		},
	})
}

func TestAccSumologicHierarchy_update(t *testing.T) {
	var hierarchy Hierarchy
	testName, testFilter, testLevel := getRandomizedParamsForHierarchy()

	testUpdatedName, testUpdatedFilter, testUpdatedLevel := getRandomizedParamsForHierarchy()

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckHierarchyDestroy(hierarchy),
		Steps: []resource.TestStep{
			{
				Config: testAccSumologicHierarchy(testName, testFilter, testLevel),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckHierarchyExists("sumologic_hierarchy.test", &hierarchy, t),
					resource.TestCheckResourceAttrSet("sumologic_hierarchy.test", "id"),
					resource.TestCheckResourceAttr("sumologic_hierarchy.test", "name", testName),
					resource.TestCheckResourceAttr("sumologic_hierarchy.test", "filter.0.key", testFilter.Key),
					resource.TestCheckResourceAttr("sumologic_hierarchy.test", "filter.0.value", testFilter.Value),
					resource.TestCheckResourceAttr("sumologic_hierarchy.test", "level.0.entity_type",
						testLevel.EntityType),
					resource.TestCheckResourceAttr("sumologic_hierarchy.test",
						"level.0.next_levels_with_conditions.0.condition",
						testLevel.NextLevelsWithConditions[0].Condition),
					resource.TestCheckResourceAttr("sumologic_hierarchy.test",
						"level.0.next_levels_with_conditions.0.level.0.entity_type",
						testLevel.NextLevelsWithConditions[0].Level.EntityType),
					resource.TestCheckResourceAttr("sumologic_hierarchy.test",
						"level.0.next_level.0.entity_type",
						testLevel.NextLevel.EntityType),
				),
			},
			{
				Config: testAccSumologicHierarchyUpdate(testUpdatedName, testUpdatedFilter, testUpdatedLevel),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckHierarchyExists("sumologic_hierarchy.test", &hierarchy, t),
					resource.TestCheckResourceAttrSet("sumologic_hierarchy.test", "id"),
					resource.TestCheckResourceAttr("sumologic_hierarchy.test", "name", testUpdatedName),
					resource.TestCheckResourceAttr("sumologic_hierarchy.test", "filter.0.key", testUpdatedFilter.Key),
					resource.TestCheckResourceAttr("sumologic_hierarchy.test", "filter.0.value", testUpdatedFilter.Value),
					resource.TestCheckResourceAttr("sumologic_hierarchy.test", "level.0.entity_type",
						testUpdatedLevel.EntityType),
					resource.TestCheckResourceAttr("sumologic_hierarchy.test",
						"level.0.next_levels_with_conditions.0.condition",
						testUpdatedLevel.NextLevelsWithConditions[0].Condition),
					resource.TestCheckResourceAttr("sumologic_hierarchy.test",
						"level.0.next_levels_with_conditions.0.level.0.entity_type",
						testUpdatedLevel.NextLevelsWithConditions[0].Level.EntityType),
					resource.TestCheckResourceAttr("sumologic_hierarchy.test",
						"level.0.next_level.0.entity_type",
						testUpdatedLevel.NextLevel.EntityType),
				),
			},
		},
	})
}

func testAccCheckHierarchyDestroy(hierarchy Hierarchy) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client := testAccProvider.Meta().(*Client)
		for _, r := range s.RootModule().Resources {
			id := r.Primary.ID
			u, err := client.GetHierarchy(id)
			if err != nil {
				return fmt.Errorf("Encountered an error: " + err.Error())
			}
			if u != nil {
				return fmt.Errorf("Hierarchy %s still exists", id)
			}
		}
		return nil
	}
}

func testAccCheckHierarchyExists(name string, hierarchy *Hierarchy, t *testing.T) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[name]
		if !ok {
			//need this so that we don't get an unused import error for strconv in some cases
			return fmt.Errorf("Error = %s. Hierarchy not found: %s", strconv.FormatBool(ok), name)
		}

		//need this so that we don't get an unused import error for strings in some cases
		if strings.EqualFold(rs.Primary.ID, "") {
			return fmt.Errorf("Hierarchy ID is not set")
		}

		id := rs.Primary.ID
		c := testAccProvider.Meta().(*Client)
		newHierarchy, err := c.GetHierarchy(id)
		if err != nil {
			return fmt.Errorf("Hierarchy %s not found", id)
		}
		hierarchy = newHierarchy
		return nil
	}
}

func testAccCheckSumologicHierarchyConfigImported(name string, filter HierarchyFilteringClause, level Level) string {
	return fmt.Sprintf(`
        resource "sumologic_hierarchy" "test" {
            name = "%s"
            filter {
                key = "%s"
                value = "%s"
            }
            level {
                entity_type = "cluster"
                next_levels_with_conditions {
                    condition = "%s"
                    level {
                        entity_type = "namespace"
                    }
                }
                next_level {
                    entity_type = "node"
                }
            }
        }`, name, filter.Key, filter.Value, level.NextLevelsWithConditions[0].Condition)
}

func testAccSumologicHierarchy(name string, filter HierarchyFilteringClause, level Level) string {
	return fmt.Sprintf(`
        resource "sumologic_hierarchy" "test" {
            name = "%s"
            filter {
                key = "%s"
                value = "%s"
            }
            level {
                entity_type = "%s"
                next_levels_with_conditions {
                    condition = "%s"
                    level {
                        entity_type = "%s"
                    }
                }
                next_level {
                    entity_type = "%s"
                }
            }
        }`,
		name, filter.Key, filter.Value, level.EntityType, level.NextLevelsWithConditions[0].Condition,
		level.NextLevelsWithConditions[0].Level.EntityType, level.NextLevel.EntityType)
}

func testAccSumologicHierarchyUpdate(name string, filter HierarchyFilteringClause, level Level) string {
	return fmt.Sprintf(`
        resource "sumologic_hierarchy" "test" {
            name = "%s"
            filter {
                key = "%s"
                value = "%s"
            }
            level {
                entity_type = "%s"
                next_levels_with_conditions {
                    condition = "%s"
                    level {
                        entity_type = "%s"
                    }
                }
                next_level {
                    entity_type = "%s"
                }
            }
        }`,
		name, filter.Key, filter.Value, level.EntityType, level.NextLevelsWithConditions[0].Condition,
		level.NextLevelsWithConditions[0].Level.EntityType, level.NextLevel.EntityType)
}
