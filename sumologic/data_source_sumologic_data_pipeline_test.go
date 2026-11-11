package sumologic

import (
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func skipDataPipelineTest(t *testing.T) {
	if strings.ToLower(os.Getenv("SKIP_DATA_PIPELINE_TESTS")) == "true" {
		t.Skip("Skipping Data Pipeline Test")
	}
}

func dataPipelineFixtureID(t *testing.T) string {
	id := os.Getenv("SUMOLOGIC_TEST_PIPELINE_ID")
	if id == "" {
		t.Skip("SUMOLOGIC_TEST_PIPELINE_ID must be set to run data pipeline data source tests")
	}
	return id
}

func TestAccDataSourceSumologicDataPipeline_byID(t *testing.T) {
	skipDataPipelineTest(t)
	id := dataPipelineFixtureID(t)

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceDataPipelineByIDConfig(id),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.sumologic_data_pipeline.test", "id", id),
					resource.TestCheckResourceAttrSet("data.sumologic_data_pipeline.test", "name"),
					resource.TestCheckResourceAttrSet("data.sumologic_data_pipeline.test", "customer_id"),
					resource.TestCheckResourceAttrSet("data.sumologic_data_pipeline.test", "version"),
					resource.TestCheckResourceAttrSet("data.sumologic_data_pipeline.test", "created_at"),
					resource.TestCheckResourceAttr("data.sumologic_data_pipeline.test", "pipeline_type", "route_based"),
					resource.TestCheckResourceAttrSet("data.sumologic_data_pipeline.test", "node.#"),
					resource.TestCheckResourceAttrSet("data.sumologic_data_pipeline.test", "node.0.name"),
					resource.TestCheckResourceAttrSet("data.sumologic_data_pipeline.test", "node.0.node_type"),
				),
			},
		},
	})
}

func TestAccDataSourceSumologicDataPipeline_byName(t *testing.T) {
	skipDataPipelineTest(t)
	dataPipelineFixtureID(t)

	name := os.Getenv("SUMOLOGIC_TEST_PIPELINE_NAME")
	if name == "" {
		t.Skip("SUMOLOGIC_TEST_PIPELINE_NAME must be set to run the by-name data pipeline test")
	}

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceDataPipelineByNameConfig(name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.sumologic_data_pipeline.by_name", "name", name),
					resource.TestCheckResourceAttrSet("data.sumologic_data_pipeline.by_name", "id"),
					resource.TestCheckResourceAttrSet("data.sumologic_data_pipeline.by_name", "node.#"),
				),
			},
		},
	})
}

func testAccDataSourceDataPipelineByIDConfig(id string) string {
	return fmt.Sprintf(`
data "sumologic_data_pipeline" "test" {
  id = "%s"
}
`, id)
}

func testAccDataSourceDataPipelineByNameConfig(name string) string {
	return fmt.Sprintf(`
data "sumologic_data_pipeline" "by_name" {
  name = "%s"
}
`, name)
}
