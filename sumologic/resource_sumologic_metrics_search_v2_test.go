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

func TestAccSumologicMetricsSearchV2_basic(t *testing.T) {
	var metricsSearchV2 MetricsSearchV2
	title := "TF Import Metrics Search V2 Test"
	description := "TF Import Metrics Search V2 Test Description"
	literalRangeName := "today"

	metricsQuery := []MetricsSearchQueryV2{
		{
			QueryKey:    "A",
			QueryString: "metric=cpu_idle | avg",
			QueryType:   "Metrics",
		},
	}

	tfResourceName := "tf_import_metrics_search_v2_test"
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckMetricsSearchV2Destroy(metricsSearchV2),
		Steps: []resource.TestStep{
			{
				Config: testAccSumologicMetricsSearchV2(tfResourceName, title, description,
					metricsQuery[0], literalRangeName),
			},
			{
				ResourceName:      fmt.Sprintf("sumologic_metrics_search_v2.%s", tfResourceName),
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccSumologicMetricsSearchV2_create(t *testing.T) {
	testNameSuffix := acctest.RandString(16)

	// create config
	var metricsSearchV2 MetricsSearchV2
	title := "terraform_test_metrics_search_v2_" + testNameSuffix
	description := "Test metrics search v2 description"
	literalRangeName := "today"

	metricsQuery := MetricsSearchQueryV2{
		QueryKey:    "A",
		QueryString: "metric=cpu_idle | avg",
		QueryType:   "Metrics",
	}

	tfResourceName := "tf_create_metrics_search_v2_test"
	tfSearchResource := fmt.Sprintf("sumologic_metrics_search_v2.%s", tfResourceName)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckMetricsSearchV2Destroy(metricsSearchV2),
		Steps: []resource.TestStep{
			{
				Config: testAccSumologicMetricsSearchV2(tfResourceName, title, description,
					metricsQuery, literalRangeName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckMetricsSearchV2Exists(tfSearchResource, &metricsSearchV2, t),
					resource.TestCheckResourceAttr(tfSearchResource,
						"title", title),
					resource.TestCheckResourceAttr(tfSearchResource,
						"description", description),
					// timerange
					resource.TestCheckResourceAttr(tfSearchResource, "time_range.#", "1"),
					resource.TestCheckResourceAttr(tfSearchResource,
						"time_range.0.begin_bounded_time_range.0.from.0.literal_time_range.0.range_name",
						literalRangeName),

					// metrics query
					resource.TestCheckResourceAttr(tfSearchResource, "queries.#", "1"),
					resource.TestCheckResourceAttr(tfSearchResource, "queries.0.query_key", metricsQuery.QueryKey),
					resource.TestCheckResourceAttr(tfSearchResource, "queries.0.query_string", metricsQuery.QueryString),
					resource.TestCheckResourceAttr(tfSearchResource, "queries.0.query_type", metricsQuery.QueryType),
				),
			},
		},
	})
}

func TestAccSumologicMetricsSearchV2_update(t *testing.T) {
	testNameSuffix := acctest.RandString(16)

	// create config
	var metricsSearchV2 MetricsSearchV2
	title := "terraform_test_metrics_search_v2_" + testNameSuffix
	description := "Test metrics search v2 description"
	literalRangeName := "today"

	metricsQuery := MetricsSearchQueryV2{
		QueryKey:    "A",
		QueryString: "metric=cpu_idle | avg",
		QueryType:   "Metrics",
	}

	newMetricsQuery := []MetricsSearchQueryV2{
		{
			QueryKey:    "B",
			QueryString: "metric=cpu_total | avg",
			QueryType:   "Metrics",
		},
	}

	tfResourceName := "tf_update_metrics_search_v2_test"
	tfSearchResource := fmt.Sprintf("sumologic_metrics_search_v2.%s", tfResourceName)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckMetricsSearchV2Destroy(metricsSearchV2),
		Steps: []resource.TestStep{
			{
				Config: testAccSumologicMetricsSearchV2(tfResourceName, title, description,
					metricsQuery, literalRangeName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckMetricsSearchV2Exists(tfSearchResource, &metricsSearchV2, t),
					resource.TestCheckResourceAttr(tfSearchResource,
						"title", title),
					resource.TestCheckResourceAttr(tfSearchResource,
						"description", description),
					// timerange
					resource.TestCheckResourceAttr(tfSearchResource, "time_range.#", "1"),
					resource.TestCheckResourceAttr(tfSearchResource,
						"time_range.0.begin_bounded_time_range.0.from.0.literal_time_range.0.range_name",
						literalRangeName),

					// metrics query
					resource.TestCheckResourceAttr(tfSearchResource, "queries.#", "1"),
					resource.TestCheckResourceAttr(tfSearchResource, "queries.0.query_key", metricsQuery.QueryKey),
					resource.TestCheckResourceAttr(tfSearchResource, "queries.0.query_string", metricsQuery.QueryString),
					resource.TestCheckResourceAttr(tfSearchResource, "queries.0.query_type", metricsQuery.QueryType),
				),
			},
			{
				Config: testAccSumologicUpdatedMetricsSearchV2(tfResourceName, title, description,
					newMetricsQuery, literalRangeName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckMetricsSearchV2Exists(tfSearchResource, &metricsSearchV2, t),
					resource.TestCheckResourceAttr(tfSearchResource,
						"title", title),
					resource.TestCheckResourceAttr(tfSearchResource,
						"description", description),
					// timerange
					resource.TestCheckResourceAttr(tfSearchResource, "time_range.#", "1"),
					resource.TestCheckResourceAttr(tfSearchResource,
						"time_range.0.begin_bounded_time_range.0.from.0.literal_time_range.0.range_name",
						literalRangeName),

					// metrics query
					resource.TestCheckResourceAttr(tfSearchResource, "queries.#", "1"),
					resource.TestCheckResourceAttr(tfSearchResource, "queries.0.query_key", newMetricsQuery[0].QueryKey),
					resource.TestCheckResourceAttr(tfSearchResource, "queries.0.query_string", newMetricsQuery[0].QueryString),
					resource.TestCheckResourceAttr(tfSearchResource, "queries.0.query_type", newMetricsQuery[0].QueryType),
				),
			},
		},
	})
}

func testAccCheckMetricsSearchV2Destroy(metricsSearchV2 MetricsSearchV2) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client := testAccProvider.Meta().(*Client)
		for _, r := range s.RootModule().Resources {
			id := r.Primary.ID
			search, err := client.GetMetricsSearchV2(id)
			if err != nil {
				return fmt.Errorf("Encountered an error: " + err.Error())
			}
			if search != nil {
				return fmt.Errorf("MetricsSearchV2 %s still exists", id)
			}
		}
		return nil
	}
}

func testAccCheckMetricsSearchV2Exists(name string, metricsSearchV2 *MetricsSearchV2, t *testing.T) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[name]
		if !ok {
			return fmt.Errorf("Error = %s. MetricsSearchV2 not found: %s", strconv.FormatBool(ok), name)
		}

		if strings.EqualFold(rs.Primary.ID, "") {
			return fmt.Errorf("MetricsSearchV2 ID is not set")
		}

		id := rs.Primary.ID
		client := testAccProvider.Meta().(*Client)
		newMetricsSearchV2, err := client.GetMetricsSearchV2(id)
		if err != nil {
			return fmt.Errorf("MetricsSearchV2 (id=%s) not found", id)
		}
		metricsSearchV2 = newMetricsSearchV2
		return nil
	}
}

func testAccSumologicMetricsSearchV2(tfResourceName string, title string, description string,
	metricsSearchQuery MetricsSearchQueryV2, literalRangeName string) string {

	return fmt.Sprintf(`
	data "sumologic_personal_folder" "personalFolder" {}

	resource "sumologic_metrics_search_v2" "%s" {
		title = "%s"
		description = "%s"
		folder_id = data.sumologic_personal_folder.personalFolder.id
		queries {
			query_key = "%s"
			query_string = "%s"
			query_type = "%s"
			metrics_query_mode = "Advanced"
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
	`, tfResourceName, title, description,
		metricsSearchQuery.QueryKey, metricsSearchQuery.QueryString,
		metricsSearchQuery.QueryType, literalRangeName)
}

func testAccSumologicUpdatedMetricsSearchV2(tfResourceName string, title string, description string,
	metricsSearchQuery []MetricsSearchQueryV2, literalRangeName string) string {

	return fmt.Sprintf(`
	data "sumologic_personal_folder" "personalFolder" {}

	resource "sumologic_metrics_search_v2" "%s" {
		title = "%s"
		description = "%s"
		folder_id = data.sumologic_personal_folder.personalFolder.id
		queries {
			query_key = "%s"
			query_string = "%s"
			query_type = "%s"
			metrics_query_mode = "Advanced"
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
	`, tfResourceName, title, description,
		metricsSearchQuery[0].QueryKey, metricsSearchQuery[0].QueryString,
		metricsSearchQuery[0].QueryType, literalRangeName)
}
