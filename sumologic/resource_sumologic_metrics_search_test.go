package sumologic

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

func TestAccSumologicMetricsSearch_basic(t *testing.T) {
	var metricsSearch MetricsSearch
	title := "TF Import Metrics Search Test"
	description := "TF Import Metrics Search Test Description"
	logQuery := ""
	literalRangeName := "today"
	desiredQuantizationInSecs := 0

	metricsQuery := []MetricsSearchQuery{
		{
			RowId: "A",
			Query: "metric=cpu_idle | avg",
		},
	}

	tfResourceName := "tf_import_metrics_search_test"
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckMetricsSearchDestroy(metricsSearch),
		Steps: []resource.TestStep{
			{
				Config: testAccSumologicMetricsSearch(tfResourceName, title, description, logQuery, desiredQuantizationInSecs,
					metricsQuery, literalRangeName),
			},
			{
				ResourceName:      fmt.Sprintf("sumologic_metrics_search.%s", tfResourceName),
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccCheckMetricsSearchDestroy(metricsSearch MetricsSearch) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client := testAccProvider.Meta().(*Client)
		for _, r := range s.RootModule().Resources {
			id := r.Primary.ID
			search, err := client.GetMetricsSearch(id)
			if err != nil {
				return fmt.Errorf("Encountered an error: " + err.Error())
			}
			if search != nil {
				return fmt.Errorf("MetricsSearch %s still exists", id)
			}
		}
		return nil
	}
}

func testAccSumologicMetricsSearch(tfResourceName string, title string, description string, logQuery string,
	desiredQuantizationInSecs int, metricsSearchQuery []MetricsSearchQuery, literalRangeName string) string {

	return fmt.Sprintf(`
	data "sumologic_personal_folder" "personalFolder" {}

	resource "sumologic_metrics_search" "%s" {
		title = "%s"
		description = "%s"
		log_query = "%s"
		parent_id = data.sumologic_personal_folder.personalFolder.id
		desired_quantization_in_secs = %d
		metrics_queries {
			row_id = "%s"
			query = "%s"
		}
		time_range {
			begin_bounded_time_range {
				from {
					literal_time_range {
						range_name = "%s"
					}
				}
			}
		}
	}
	`, tfResourceName, title, description, logQuery, desiredQuantizationInSecs,
		metricsSearchQuery[0].RowId, metricsSearchQuery[0].Query,
		literalRangeName)
}
