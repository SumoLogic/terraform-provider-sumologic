package sumologic

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccDataSourceSumologicDataPipelineTest_regexParse(t *testing.T) {
	skipDataPipelineTest(t)
	dataSourceName := "data.sumologic_data_pipeline_test.test"

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceDataPipelineTestConfig(`\n`),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(dataSourceName, "id"),
					resource.TestCheckResourceAttr(dataSourceName, "message.#", "2"),
					resource.TestCheckResourceAttr(dataSourceName, "message.0.message_index", "0"),
					resource.TestCheckResourceAttr(dataSourceName, "message.0.input", "line one"),
					resource.TestCheckResourceAttr(dataSourceName, "message.0.output", `{"_raw":"line one","w":"one"}`),
					resource.TestCheckResourceAttr(dataSourceName, "message.1.input", "line two"),
					resource.TestCheckResourceAttr(dataSourceName, "message.1.output", `{"_raw":"line two","w":"two"}`),
					resource.TestCheckResourceAttr(dataSourceName, "message.0.node_result.0.node_name", "Parse"),
					resource.TestCheckResourceAttr(dataSourceName, "message.0.node_result.0.node_type", "processing_group"),
					resource.TestCheckResourceAttr(dataSourceName, "message.0.node_result.0.matched", "true"),
					resource.TestCheckResourceAttr(dataSourceName, "message.0.node_result.0.processor_result.0.processor_name", "extract-last-word"),
					resource.TestCheckResourceAttr(dataSourceName, "message.0.node_result.0.processor_result.0.processor_type", "REGEX_PARSE"),
					resource.TestCheckResourceAttr(dataSourceName, "message.0.node_result.0.processor_result.0.message_dropped", "false"),
					resource.TestCheckResourceAttr(dataSourceName, "message.0.edge_result.#", "2"),
					resource.TestCheckResourceAttr(dataSourceName, "message.0.edge_result.0.source_node", "Routing Expression"),
					resource.TestCheckResourceAttr(dataSourceName, "message.0.edge_result.0.target_node", "Parse"),
					resource.TestCheckResourceAttr(dataSourceName, "message.0.edge_result.1.target_node", "Sumo Logic"),
				),
			},
			{
				Config: testAccDataSourceDataPipelineTestConfig(`(empty)`),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(dataSourceName, "message.#", "1"),
					resource.TestCheckResourceAttr(dataSourceName, "message.0.input", "line one\nline two"),
				),
			},
		},
	})
}

func testAccDataSourceDataPipelineTestConfig(separator string) string {
	return `
data "sumologic_data_pipeline_test" "test" {
  name             = "tf-data-pipeline-sample-test"
  route_expression = "_sourceCategory=tf-provider-test"
  sample_data      = "line one\nline two"

  message_separator = "` + escapeDataPipelineTestSeparator(separator) + `"

  node {
    name      = "Routing Expression"
    node_type = "source"
    output {
      target = "Parse"
    }
  }

  node {
    name      = "Parse"
    node_type = "processing_group"
    processor {
      name           = "extract-last-word"
      processor_type = "REGEX_PARSE"
      config = jsonencode({
        fieldToParseFrom = "_raw"
        regexPattern     = "(?<w>\\w+)$"
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
`
}

func escapeDataPipelineTestSeparator(separator string) string {
	escaped := ""
	for _, r := range separator {
		if r == '\\' {
			escaped += `\\`
			continue
		}
		escaped += string(r)
	}
	return escaped
}
