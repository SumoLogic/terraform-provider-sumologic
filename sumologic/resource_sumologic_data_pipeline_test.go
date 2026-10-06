package sumologic

import (
	"fmt"
	"strconv"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccSumologicDataPipeline_createAndUpdate(t *testing.T) {
	skipDataPipelineTest(t)
	var pipeline DataPipeline
	var pipelineID string
	resourceName := "sumologic_data_pipeline.test"
	name := acctest.RandomWithPrefix("tf-data-pipeline-test")
	renamedName := name + "-renamed"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckDataPipelineDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccDataPipelineConfig(name, "", "_sourceCategory=tf-provider-test", true),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDataPipelineExists(resourceName, &pipeline),
					resource.TestCheckResourceAttr(resourceName, "name", name),
					resource.TestCheckResourceAttr(resourceName, "pipeline_type", "route_based"),
					resource.TestCheckResourceAttr(resourceName, "route_expression", "_sourceCategory=tf-provider-test"),
					resource.TestCheckResourceAttr(resourceName, "is_enabled", "true"),
					resource.TestCheckResourceAttr(resourceName, "node.#", "2"),
					resource.TestCheckResourceAttr(resourceName, "node.0.name", "Routing Expression"),
					resource.TestCheckResourceAttr(resourceName, "node.0.node_type", "source"),
					resource.TestCheckResourceAttr(resourceName, "node.0.output.0.target", "Sumo Logic"),
					resource.TestCheckResourceAttr(resourceName, "node.1.name", "Sumo Logic"),
					resource.TestCheckResourceAttr(resourceName, "node.1.node_type", "destination"),
					resource.TestCheckResourceAttrSet(resourceName, "version"),
				),
			},
			{
				Config: testAccDataPipelineConfig(name, "", "_sourceCategory=tf-provider-test-updated", false),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDataPipelineExists(resourceName, &pipeline),
					resource.TestCheckResourceAttr(resourceName, "state", "published"),
					resource.TestCheckResourceAttr(resourceName, "route_expression", "_sourceCategory=tf-provider-test-updated"),
					resource.TestCheckResourceAttr(resourceName, "is_enabled", "false"),
					func(s *terraform.State) error {
						pipelineID = pipeline.ID
						return nil
					},
				),
			},
			{
				Config: testAccDataPipelineConfig(renamedName, "renamed via updatePipelineMetadata", "_sourceCategory=tf-provider-test-updated", false),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDataPipelineExists(resourceName, &pipeline),
					resource.TestCheckResourceAttr(resourceName, "name", renamedName),
					resource.TestCheckResourceAttr(resourceName, "description", "renamed via updatePipelineMetadata"),
					func(s *terraform.State) error {
						if pipeline.ID != pipelineID {
							return fmt.Errorf("expected rename to update the pipeline in place, but id changed from %s to %s", pipelineID, pipeline.ID)
						}
						return nil
					},
				),
			},
		},
	})
}

func TestAccSumologicDataPipeline_nodeDiffing(t *testing.T) {
	skipDataPipelineTest(t)
	var pipeline DataPipeline
	resourceName := "sumologic_data_pipeline.test"
	name := acctest.RandomWithPrefix("tf-data-pipeline-nodediff-test")

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckDataPipelineDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccDataPipelineNodeDiffConfigBaseline(name),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckDataPipelineExists(resourceName, &pipeline),
					resource.TestCheckResourceAttr(resourceName, "node.#", "3"),
					testAccCheckDataPipelineHasNode(resourceName, "Routing Expression", "source"),
					testAccCheckDataPipelineHasNode(resourceName, "Triage", "router"),
					testAccCheckDataPipelineHasNode(resourceName, "Sumo Logic", "destination"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccDataPipelineNodeDiffConfigNodeAdded(name),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckDataPipelineExists(resourceName, &pipeline),
					resource.TestCheckResourceAttr(resourceName, "node.#", "4"),
					testAccCheckDataPipelineHasNode(resourceName, "Routing Expression", "source"),
					testAccCheckDataPipelineHasNode(resourceName, "Triage", "router"),
					testAccCheckDataPipelineHasNode(resourceName, "Secondary Triage", "router"),
					testAccCheckDataPipelineHasNode(resourceName, "Sumo Logic", "destination"),
				),
			},
			{
				Config: testAccDataPipelineNodeDiffConfigBaseline(name),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckDataPipelineExists(resourceName, &pipeline),
					resource.TestCheckResourceAttr(resourceName, "node.#", "3"),
					testAccCheckDataPipelineHasNode(resourceName, "Routing Expression", "source"),
					testAccCheckDataPipelineHasNode(resourceName, "Triage", "router"),
					testAccCheckDataPipelineHasNode(resourceName, "Sumo Logic", "destination"),
				),
			},
			{
				Config: testAccDataPipelineNodeDiffConfigReordered(name),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckDataPipelineExists(resourceName, &pipeline),
					resource.TestCheckResourceAttr(resourceName, "node.#", "3"),
					testAccCheckDataPipelineHasNode(resourceName, "Routing Expression", "source"),
					testAccCheckDataPipelineHasNode(resourceName, "Triage", "router"),
					testAccCheckDataPipelineHasNode(resourceName, "Sumo Logic", "destination"),
				),
			},
		},
	})
}

func TestAccSumologicDataPipeline_processingGroup(t *testing.T) {
	skipDataPipelineTest(t)
	var pipeline DataPipeline
	resourceName := "sumologic_data_pipeline.test"
	name := acctest.RandomWithPrefix("tf-data-pipeline-processinggroup-test")

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckDataPipelineDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccDataPipelineProcessingGroupConfigBaseline(name),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckDataPipelineExists(resourceName, &pipeline),
					testAccCheckDataPipelineHasNode(resourceName, "Log Normalization", "processing_group"),
					testAccCheckDataPipelineHasProcessor(resourceName, "Log Normalization", "Extract Client IP", "REGEX_PARSE"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccDataPipelineProcessingGroupConfigProcessorAdded(name),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckDataPipelineExists(resourceName, &pipeline),
					testAccCheckDataPipelineHasProcessor(resourceName, "Log Normalization", "Extract Client IP", "REGEX_PARSE"),
					testAccCheckDataPipelineHasProcessor(resourceName, "Log Normalization", "Extract Timestamp", "REGEX_PARSE"),
				),
			},
			{
				Config: testAccDataPipelineProcessingGroupConfigBaseline(name),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckDataPipelineExists(resourceName, &pipeline),
					testAccCheckDataPipelineHasProcessor(resourceName, "Log Normalization", "Extract Client IP", "REGEX_PARSE"),
				),
			},
			{
				Config: testAccDataPipelineProcessingGroupConfigReordered(name),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckDataPipelineExists(resourceName, &pipeline),
					testAccCheckDataPipelineHasProcessor(resourceName, "Log Normalization", "Extract Client IP", "REGEX_PARSE"),
					testAccCheckDataPipelineHasProcessor(resourceName, "Log Normalization", "Extract Timestamp", "REGEX_PARSE"),
				),
			},
		},
	})
}

func testAccCheckDataPipelineHasProcessor(resourceName, nodeName, processorName, processorType string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("resource not found: %s", resourceName)
		}

		nodeCount, err := strconv.Atoi(rs.Primary.Attributes["node.#"])
		if err != nil {
			return fmt.Errorf("node.# is not a number: %w", err)
		}

		for i := 0; i < nodeCount; i++ {
			if rs.Primary.Attributes[fmt.Sprintf("node.%d.name", i)] != nodeName {
				continue
			}

			processorCount, err := strconv.Atoi(rs.Primary.Attributes[fmt.Sprintf("node.%d.processor.#", i)])
			if err != nil {
				return fmt.Errorf("node.%d.processor.# is not a number: %w", i, err)
			}

			for j := 0; j < processorCount; j++ {
				if rs.Primary.Attributes[fmt.Sprintf("node.%d.processor.%d.name", i, j)] != processorName {
					continue
				}
				if got := rs.Primary.Attributes[fmt.Sprintf("node.%d.processor.%d.processor_type", i, j)]; got != processorType {
					return fmt.Errorf("processor %q processor_type = %q, want %q", processorName, got, processorType)
				}
				if rs.Primary.Attributes[fmt.Sprintf("node.%d.processor.%d.id", i, j)] == "" {
					return fmt.Errorf("processor %q has an empty id", processorName)
				}
				return nil
			}
			return fmt.Errorf("processor %q not found on node %q among %d processors", processorName, nodeName, processorCount)
		}
		return fmt.Errorf("node %q not found among %d nodes", nodeName, nodeCount)
	}
}

func testAccDataPipelineProcessingGroupConfigBaseline(name string) string {
	return fmt.Sprintf(`
resource "sumologic_data_pipeline" "test" {
  name             = %[1]q
  route_expression = "_sourceCategory=tf-provider-processinggroup-test"

  node {
    name      = "Routing Expression"
    node_type = "source"

    output {
      target = "Log Normalization"
    }
  }

  node {
    name      = "Log Normalization"
    node_type = "processing_group"

    processor {
      name           = "Extract Client IP"
      processor_type = "REGEX_PARSE"
      config = jsonencode({
        fieldToParseFrom = "_raw"
        regexPattern     = "(?<clientip>[\\d.]+)"
      })
    }

    output {
      target = "Sumo Logic"
    }
  }

  node {
    name      = "Sumo Logic"
    node_type = "destination"
  }
}
`, name)
}

func testAccDataPipelineProcessingGroupConfigProcessorAdded(name string) string {
	return fmt.Sprintf(`
resource "sumologic_data_pipeline" "test" {
  name             = %[1]q
  route_expression = "_sourceCategory=tf-provider-processinggroup-test"

  node {
    name      = "Routing Expression"
    node_type = "source"

    output {
      target = "Log Normalization"
    }
  }

  node {
    name      = "Log Normalization"
    node_type = "processing_group"

    processor {
      name           = "Extract Client IP"
      processor_type = "REGEX_PARSE"
      config = jsonencode({
        fieldToParseFrom = "_raw"
        regexPattern     = "(?<clientip>[\\d.]+)"
      })
    }

    processor {
      name           = "Extract Timestamp"
      processor_type = "REGEX_PARSE"
      config = jsonencode({
        fieldToParseFrom = "_raw"
        regexPattern     = "\\[(?<timestamp>[^\\]]+)\\]"
      })
    }

    output {
      target = "Sumo Logic"
    }
  }

  node {
    name      = "Sumo Logic"
    node_type = "destination"
  }
}
`, name)
}

func testAccDataPipelineProcessingGroupConfigReordered(name string) string {
	return fmt.Sprintf(`
resource "sumologic_data_pipeline" "test" {
  name             = %[1]q
  route_expression = "_sourceCategory=tf-provider-processinggroup-test"

  node {
    name      = "Routing Expression"
    node_type = "source"

    output {
      target = "Log Normalization"
    }
  }

  node {
    name      = "Log Normalization"
    node_type = "processing_group"

    processor {
      name           = "Extract Timestamp"
      processor_type = "REGEX_PARSE"
      config = jsonencode({
        fieldToParseFrom = "_raw"
        regexPattern     = "\\[(?<timestamp>[^\\]]+)\\]"
      })
    }

    processor {
      name           = "Extract Client IP"
      processor_type = "REGEX_PARSE"
      config = jsonencode({
        fieldToParseFrom = "_raw"
        regexPattern     = "(?<clientip>[\\d.]+)"
      })
    }

    output {
      target = "Sumo Logic"
    }
  }

  node {
    name      = "Sumo Logic"
    node_type = "destination"
  }
}
`, name)
}

func testAccCheckDataPipelineHasNode(resourceName, nodeName, nodeType string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("resource not found: %s", resourceName)
		}

		count, err := strconv.Atoi(rs.Primary.Attributes["node.#"])
		if err != nil {
			return fmt.Errorf("node.# is not a number: %w", err)
		}

		for i := 0; i < count; i++ {
			if rs.Primary.Attributes[fmt.Sprintf("node.%d.name", i)] != nodeName {
				continue
			}
			if got := rs.Primary.Attributes[fmt.Sprintf("node.%d.node_type", i)]; got != nodeType {
				return fmt.Errorf("node %q node_type = %q, want %q", nodeName, got, nodeType)
			}
			if rs.Primary.Attributes[fmt.Sprintf("node.%d.id", i)] == "" {
				return fmt.Errorf("node %q has an empty id", nodeName)
			}
			return nil
		}
		return fmt.Errorf("node %q not found among %d nodes", nodeName, count)
	}
}

func testAccDataPipelineNodeDiffConfigBaseline(name string) string {
	return fmt.Sprintf(`
resource "sumologic_data_pipeline" "test" {
  name             = %[1]q
  route_expression = "_sourceCategory=tf-provider-nodediff-test"

  node {
    name      = "Routing Expression"
    node_type = "source"

    output {
      target = "Triage"
    }
  }

  node {
    name      = "Triage"
    node_type = "router"

    output {
      target = "Sumo Logic"
    }
  }

  node {
    name      = "Sumo Logic"
    node_type = "destination"
  }
}
`, name)
}

func testAccDataPipelineNodeDiffConfigNodeAdded(name string) string {
	return fmt.Sprintf(`
resource "sumologic_data_pipeline" "test" {
  name             = %[1]q
  route_expression = "_sourceCategory=tf-provider-nodediff-test"

  node {
    name      = "Routing Expression"
    node_type = "source"

    output {
      target = "Triage"
    }
  }

  node {
    name      = "Triage"
    node_type = "router"

    output {
      target = "Secondary Triage"
    }
  }

  node {
    name      = "Secondary Triage"
    node_type = "router"

    output {
      target = "Sumo Logic"
    }
  }

  node {
    name      = "Sumo Logic"
    node_type = "destination"
  }
}
`, name)
}

func testAccDataPipelineNodeDiffConfigReordered(name string) string {
	return fmt.Sprintf(`
resource "sumologic_data_pipeline" "test" {
  name             = %[1]q
  route_expression = "_sourceCategory=tf-provider-nodediff-test"

  node {
    name      = "Sumo Logic"
    node_type = "destination"
  }

  node {
    name      = "Triage"
    node_type = "router"

    output {
      target = "Sumo Logic"
    }
  }

  node {
    name      = "Routing Expression"
    node_type = "source"

    output {
      target = "Triage"
    }
  }
}
`, name)
}

func testAccCheckDataPipelineExists(name string, pipeline *DataPipeline) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[name]
		if !ok {
			return fmt.Errorf("data pipeline resource not found: %s", name)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("data pipeline id is not set")
		}

		c := testAccProvider.Meta().(*Client)
		found, err := c.GetDataPipeline(rs.Primary.ID)
		if err != nil {
			return err
		}
		if found == nil {
			return fmt.Errorf("data pipeline %s not found", rs.Primary.ID)
		}

		*pipeline = *found
		return nil
	}
}

func testAccCheckDataPipelineDestroy(s *terraform.State) error {
	c := testAccProvider.Meta().(*Client)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "sumologic_data_pipeline" {
			continue
		}

		pipeline, err := c.GetDataPipeline(rs.Primary.ID)
		if err != nil {
			return fmt.Errorf("encountered an error: %w", err)
		}
		if pipeline != nil {
			return fmt.Errorf("data pipeline %s still exists", rs.Primary.ID)
		}
	}
	return nil
}

func testAccDataPipelineConfig(name, description, routeExpression string, isEnabled bool) string {
	return fmt.Sprintf(`
resource "sumologic_data_pipeline" "test" {
  name             = "%s"
  description      = "%s"
  route_expression = "%s"
  is_enabled       = %t

  node {
    name      = "Routing Expression"
    node_type = "source"

    output {
      target = "Sumo Logic"
    }
  }

  node {
    name      = "Sumo Logic"
    node_type = "destination"
  }
}
`, name, description, routeExpression, isEnabled)
}
