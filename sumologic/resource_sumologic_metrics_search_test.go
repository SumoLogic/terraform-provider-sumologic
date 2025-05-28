package sumologic

import (
	"fmt"
	"strconv"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
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
					metricsQuery[0], literalRangeName),
			},
			{
				ResourceName:      fmt.Sprintf("sumologic_metrics_search.%s", tfResourceName),
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccSumologicMetricsSearch_create(t *testing.T) {
	testNameSuffix := acctest.RandString(16)

	// create config
	var metricsSearch MetricsSearch
	title := "terraform_test_metrics_search_" + testNameSuffix
	description := "Test metrics search description"
	logQuery := ""
	literalRangeName := "today"
	desiredQuantizationInSecs := 0

	metricsQuery := MetricsSearchQuery{
		RowId: "A",
		Query: "metric=cpu_idle | avg",
	}

	tfResourceName := "tf_create_metrics_search_test"
	tfSearchResource := fmt.Sprintf("sumologic_metrics_search.%s", tfResourceName)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckMetricsSearchDestroy(metricsSearch),
		Steps: []resource.TestStep{
			{
				Config: testAccSumologicMetricsSearch(tfResourceName, title, description, logQuery, desiredQuantizationInSecs,
					metricsQuery, literalRangeName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckMetricsSearchExists(tfSearchResource, &metricsSearch, t),
					resource.TestCheckResourceAttr(tfSearchResource,
						"title", title),
					resource.TestCheckResourceAttr(tfSearchResource,
						"description", description),
					// timerange
					resource.TestCheckResourceAttr(tfSearchResource, "time_range.#", "1"),
					resource.TestCheckResourceAttr(tfSearchResource,
						"time_range.0.begin_bounded_time_range.0.from.0.literal_time_range.0.range_name",
						literalRangeName),

					resource.TestCheckResourceAttr(tfSearchResource, "desired_quantization_in_secs", "0"),

					// metrics query
					resource.TestCheckResourceAttr(tfSearchResource, "metrics_queries.#", "1"),
					resource.TestCheckResourceAttr(tfSearchResource, "metrics_queries.0.row_id", metricsQuery.RowId),
					resource.TestCheckResourceAttr(tfSearchResource, "metrics_queries.0.query", metricsQuery.Query),
				),
			},
		},
	})
}

func TestAccSumologicMetricsSearch_update(t *testing.T) {
	testNameSuffix := acctest.RandString(16)

	// create config
	var metricsSearch MetricsSearch
	title := "terraform_test_metrics_search_" + testNameSuffix
	description := "Test metrics search description"
	logQuery := ""
	literalRangeName := "today"
	desiredQuantizationInSecs := 0

	metricsQuery := MetricsSearchQuery{
		RowId: "A",
		Query: "metric=cpu_idle | avg",
	}

	newMetricsQuery := []MetricsSearchQuery{
		{
			RowId: "B",
			Query: "metric=cpu_total | avg",
		},
	}

	tfResourceName := "tf_update_metrics_search_test"
	tfSearchResource := fmt.Sprintf("sumologic_metrics_search.%s", tfResourceName)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckMetricsSearchDestroy(metricsSearch),
		Steps: []resource.TestStep{
			{
				Config: testAccSumologicMetricsSearch(tfResourceName, title, description, logQuery, desiredQuantizationInSecs,
					metricsQuery, literalRangeName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckMetricsSearchExists(tfSearchResource, &metricsSearch, t),
					resource.TestCheckResourceAttr(tfSearchResource,
						"title", title),
					resource.TestCheckResourceAttr(tfSearchResource,
						"description", description),
					// timerange
					resource.TestCheckResourceAttr(tfSearchResource, "time_range.#", "1"),
					resource.TestCheckResourceAttr(tfSearchResource,
						"time_range.0.begin_bounded_time_range.0.from.0.literal_time_range.0.range_name",
						literalRangeName),

					resource.TestCheckResourceAttr(tfSearchResource, "desired_quantization_in_secs", "0"),

					// metrics query
					resource.TestCheckResourceAttr(tfSearchResource, "metrics_queries.#", "1"),
					resource.TestCheckResourceAttr(tfSearchResource, "metrics_queries.0.row_id", metricsQuery.RowId),
					resource.TestCheckResourceAttr(tfSearchResource, "metrics_queries.0.query", metricsQuery.Query),
				),
			},
			{
				Config: testAccSumologicUpdatedMetricsSearch(tfResourceName, title, description, logQuery, desiredQuantizationInSecs,
					newMetricsQuery, literalRangeName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckMetricsSearchExists(tfSearchResource, &metricsSearch, t),
					resource.TestCheckResourceAttr(tfSearchResource,
						"title", title),
					resource.TestCheckResourceAttr(tfSearchResource,
						"description", description),
					// timerange
					resource.TestCheckResourceAttr(tfSearchResource, "time_range.#", "1"),
					resource.TestCheckResourceAttr(tfSearchResource,
						"time_range.0.begin_bounded_time_range.0.from.0.literal_time_range.0.range_name",
						literalRangeName),

					resource.TestCheckResourceAttr(tfSearchResource, "desired_quantization_in_secs", "0"),

					// metrics query
					resource.TestCheckResourceAttr(tfSearchResource, "metrics_queries.#", "1"),
					resource.TestCheckResourceAttr(tfSearchResource, "metrics_queries.0.row_id", newMetricsQuery[0].RowId),
					resource.TestCheckResourceAttr(tfSearchResource, "metrics_queries.0.query", newMetricsQuery[0].Query),
				),
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

func testAccCheckMetricsSearchExists(name string, metricsSearch *MetricsSearch, t *testing.T) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[name]
		if !ok {
			return fmt.Errorf("Error = %s. MetricsSearch not found: %s", strconv.FormatBool(ok), name)
		}

		if strings.EqualFold(rs.Primary.ID, "") {
			return fmt.Errorf("MetricsSearch ID is not set")
		}

		id := rs.Primary.ID
		client := testAccProvider.Meta().(*Client)
		newMetricsSearch, err := client.GetMetricsSearch(id)
		if err != nil {
			return fmt.Errorf("MetricsSearch (id=%s) not found", id)
		}
		metricsSearch = newMetricsSearch
		return nil
	}
}

func testAccSumologicMetricsSearch(tfResourceName string, title string, description string, logQuery string,
	desiredQuantizationInSecs int, metricsSearchQuery MetricsSearchQuery, literalRangeName string) string {

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
		metricsSearchQuery.RowId, metricsSearchQuery.Query,
		literalRangeName)
}

func testAccSumologicUpdatedMetricsSearch(tfResourceName string, title string, description string, logQuery string,
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
