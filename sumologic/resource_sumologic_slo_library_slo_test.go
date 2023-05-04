package sumologic

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"strconv"
	"strings"
	"testing"
)

var allExampleSlos = []func(testName string) string{
	exampleLogsWindowThresholdSlo,
	exampleLogsWindowRatioSlo,
	exampleLogsRequestThresholdSlo,
	exampleLogsRequestRatioSlo,
	exampleMetricsWindowThresholdSlo,
	exampleMetricsWindowRatioSlo,
	exampleMetricsRequestRatioSlo,
}

func testAccCheckSloLibrarySloDestroy(sloLibrarySlo SLOLibrarySLO) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client := testAccProvider.Meta().(*Client)
		for _, r := range s.RootModule().Resources {
			id := r.Primary.ID
			u, err := client.SLORead(id)
			if err != nil {
				return fmt.Errorf("Encountered an error: " + err.Error())
			}
			if u != nil {
				return fmt.Errorf("SloLibrarySlo %s still exists", id)
			}
		}
		return nil
	}
}

func testAccCheckSloLibrarySloExists(name string, sloLibrarySlo *SLOLibrarySLO, t *testing.T) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[name]
		if !ok {
			//need this so that we don't get an unused import error for strconv in some cases
			return fmt.Errorf("Error = %s. SloLibrarySlo not found: %s", strconv.FormatBool(ok), name)
		}

		//need this so that we don't get an unused import error for strings in some cases
		if strings.EqualFold(rs.Primary.ID, "") {
			return fmt.Errorf("SloLibrarySlo ID is not set")
		}

		id := rs.Primary.ID
		client := testAccProvider.Meta().(*Client)
		newSloLibrarySlo, err := client.SLORead(id)
		if err != nil {
			return fmt.Errorf("SloLibrarySlo %s not found", id)
		}
		sloLibrarySlo = newSloLibrarySlo
		return nil
	}
}

func TestAccSumologicSloLibrarySlo_create_all_slo_types(t *testing.T) {
	var sloLibrarySlo SLOLibrarySLO
	for _, sloConfig := range allExampleSlos {
		testNameSuffix := acctest.RandString(16)

		testName := "terraform_test_slo_" + testNameSuffix

		resource.Test(t, resource.TestCase{
			PreCheck:     func() { testAccPreCheck(t) },
			Providers:    testAccProviders,
			CheckDestroy: testAccCheckSloLibrarySloDestroy(sloLibrarySlo),
			Steps: []resource.TestStep{
				{
					Config: sloConfig(testName),
					Check: resource.ComposeTestCheckFunc(
						testAccCheckSloLibrarySloExists("sumologic_slo.test", &sloLibrarySlo, t),
						resource.TestCheckResourceAttr("sumologic_slo.test", "name", testName),
					),
				},
			},
		})
	}
}

func exampleLogsWindowThresholdSlo(testName string) string {
	var resourceText = fmt.Sprintf(`resource "sumologic_slo_folder" "tf_slo_folder" {
  name        = "slo-tf-test-folder"
  description = "folder for SLO created for testing"
}

resource "sumologic_slo" "test" {
  name        = "%s"
  description = "per minute login error rate over rolling 1 day"
  parent_id   = sumologic_slo_folder.tf_slo_folder.id
  signal_type = "Error"
  service     = "auth"
  application = "login"
  compliance {
      compliance_type = "Rolling"
      size            = "1d"
      target          = 95
      timezone        = "Asia/Kolkata"
  }
  tags = {
    team = "metrics"
    application = "sumologic"
  }
  indicator {
    window_based_evaluation {
      op         = "LessThan"
      query_type = "Logs"
      size       = "1m"
      threshold  = 99.0
      aggregation = "Avg"
      queries {
        query_group_type = "Threshold"
        query_group {
          row_id        = "A"
          query         = "(deployment=api ApiHttpRequestLoggingFilter Request) | parse \"Request: '* *', QueryString: '*' Status: '*' (took: *ms)\" as method,url,queryString,status,timeTaken | toInt(status) | where status >= 200 and status < 300"
          use_row_count = false
          field = "status"
        }
      }
    }
  }
}`, testName)
	return resourceText
}

func exampleLogsWindowRatioSlo(testName string) string {
	var resourceText = fmt.Sprintf(`resource "sumologic_slo_folder" "tf_slo_folder" {
  name        = "slo-tf-test-folder"
  description = "folder for SLO created for testing"
}

resource "sumologic_slo" "test" {
  name        = "%s"
  description = "per minute login error rate over rolling 1 day"
  parent_id   = sumologic_slo_folder.tf_slo_folder.id
  signal_type = "Error"
  service     = "auth"
  application = "login"
  compliance {
      compliance_type = "Rolling"
      size            = "1d"
      target          = 95
      timezone        = "Asia/Kolkata"
  }
  tags = {
    team = "metrics"
    application = "sumologic"
  }
  indicator {
    window_based_evaluation {
      op         = "LessThan"
      query_type = "Logs"
      size       = "1m"
      threshold  = 99.0
      queries {
        query_group_type = "Successful"
        query_group {
          row_id        = "A"
          query         = "(deployment=api ApiHttpRequestLoggingFilter Request) | parse \"Request: '* *', QueryString: '*' Status: '*' (took: *ms)\" as method,url,queryString,status,timeTaken | toInt(status) | where status >= 200 and status < 300"
          use_row_count = false
          field = "status"
        }
      }
      queries {
        query_group_type = "Total"
        query_group {
          row_id        = "A"
          query         = "(deployment=api ApiHttpRequestLoggingFilter Request) | parse \"Request: '* *', QueryString: '*' Status: '*' (took: *ms)\" as method,url,queryString,status,timeTaken | toInt(status) | where status >= 200 and status < 300"
          use_row_count = false
          field = "status"
        }
      }
    }
  }
}`, testName)
	return resourceText
}

func exampleLogsRequestThresholdSlo(testName string) string {
	var resourceText = fmt.Sprintf(`resource "sumologic_slo_folder" "tf_slo_folder" {
  name        = "slo-tf-test-folder"
  description = "folder for SLO created for testing"
}

resource "sumologic_slo" "test" {
  name        = "%s"
  description = "per minute login error rate over rolling 1 day"
  parent_id   = sumologic_slo_folder.tf_slo_folder.id
  signal_type = "Error"
  service     = "auth"
  application = "login"
  compliance {
      compliance_type = "Rolling"
      size            = "1d"
      target          = 95
      timezone        = "Asia/Kolkata"
  }
  tags = {
    team = "metrics"
    application = "sumologic"
  }
  indicator {
    request_based_evaluation {
      op         = "LessThan"
      query_type = "Logs"
      threshold  = 99.0
      queries {
        query_group_type = "Threshold"
        query_group {
          row_id        = "A"
          query         = "(deployment=api ApiHttpRequestLoggingFilter Request) | parse \"Request: '* *', QueryString: '*' Status: '*' (took: *ms)\" as method,url,queryString,status,timeTaken | toInt(status) | where status >= 200 and status < 300"
          use_row_count = false
          field = "status"
        }
      }
    }
  }
}`, testName)
	return resourceText
}

func exampleLogsRequestRatioSlo(testName string) string {
	var resourceText = fmt.Sprintf(`resource "sumologic_slo_folder" "tf_slo_folder" {
  name        = "slo-tf-test-folder"
  description = "folder for SLO created for testing"
}

resource "sumologic_slo" "test" {
  name        = "%s"
  description = "per minute login error rate over rolling 1 day"
  parent_id   = sumologic_slo_folder.tf_slo_folder.id
  signal_type = "Error"
  service     = "auth"
  application = "login"
  compliance {
      compliance_type = "Rolling"
      size            = "1d"
      target          = 95
      timezone        = "Asia/Kolkata"
  }
  tags = {
    team = "metrics"
    application = "sumologic"
  }
  indicator {
    request_based_evaluation {
      query_type = "Logs"
      queries {
        query_group_type = "Successful"
        query_group {
          row_id        = "A"
          query         = "(deployment=api ApiHttpRequestLoggingFilter Request) | parse \"Request: '* *', QueryString: '*' Status: '*' (took: *ms)\" as method,url,queryString,status,timeTaken | toInt(status) | where status >= 200 and status < 300"
          use_row_count = false
          field = "status"
        }
      }
	  queries {
        query_group_type = "Total"
        query_group {
          row_id        = "A"
          query         = "(deployment=api ApiHttpRequestLoggingFilter Request) | parse \"Request: '* *', QueryString: '*' Status: '*' (took: *ms)\" as method,url,queryString,status,timeTaken | toInt(status) | where status >= 200 and status < 300"
          use_row_count = false
          field = "status"
        }
      }
    }
  }
}`, testName)
	return resourceText
}

func exampleMetricsWindowThresholdSlo(testName string) string {
	var resourceText = fmt.Sprintf(`resource "sumologic_slo_folder" "tf_slo_folder" {
  name        = "slo-tf-test-folder"
  description = "folder for SLO created for testing"
}

resource "sumologic_slo" "test" {
  name        = "%s"
  description = "per minute login error rate over rolling 1 day"
  parent_id   = sumologic_slo_folder.tf_slo_folder.id
  signal_type = "Error"
  service     = "auth"
  application = "login"
  compliance {
      compliance_type = "Rolling"
      size            = "1d"
      target          = 95
      timezone        = "Asia/Kolkata"
  }
  tags = {
    team = "metrics"
    application = "sumologic"
  }
  indicator {
    window_based_evaluation {
      op         = "LessThan"
      query_type = "Metrics"
      size       = "1m"
      threshold  = 99.0
      aggregation = "Avg"
      queries {
        query_group_type = "Threshold"
        query_group {
          row_id        = "A"
          query         = "example"
          use_row_count = false
        }
      }
    }
  }
}`, testName)
	return resourceText
}

func exampleMetricsWindowRatioSlo(testName string) string {
	var resourceText = fmt.Sprintf(`resource "sumologic_slo_folder" "tf_slo_folder" {
  name        = "slo-tf-test-folder"
  description = "folder for SLO created for testing"
}

resource "sumologic_slo" "test" {
  name        = "%s"
  description = "per minute login error rate over rolling 1 day"
  parent_id   = sumologic_slo_folder.tf_slo_folder.id
  signal_type = "Error"
  service     = "auth"
  application = "login"
  compliance {
      compliance_type = "Rolling"
      size            = "1d"
      target          = 95
      timezone        = "Asia/Kolkata"
  }
  tags = {
    team = "metrics"
    application = "sumologic"
  }
  indicator {
    window_based_evaluation {
      op         = "LessThan"
      query_type = "Metrics"
      size       = "1m"
      threshold  = 99.0
      queries {
        query_group_type = "Successful"
        query_group {
          row_id        = "A"
          query         = "example"
          use_row_count = true
        }
      }
	  queries {
        query_group_type = "Total"
        query_group {
          row_id        = "A"
          query         = "example"
          use_row_count = false
        }
      }
    }
  }
}`, testName)
	return resourceText
}

func exampleMetricsRequestRatioSlo(testName string) string {
	var resourceText = fmt.Sprintf(`resource "sumologic_slo_folder" "tf_slo_folder" {
  name        = "slo-tf-test-folder"
  description = "folder for SLO created for testing"
}

resource "sumologic_slo" "test" {
  name        = "%s"
  description = "per minute login error rate over rolling 1 day"
  parent_id   = sumologic_slo_folder.tf_slo_folder.id
  signal_type = "Error"
  service     = "auth"
  application = "login"
  compliance {
      compliance_type = "Rolling"
      size            = "1d"
      target          = 95
      timezone        = "Asia/Kolkata"
  }
  tags = {
    team = "metrics"
    application = "sumologic"
  }
  indicator {
    request_based_evaluation {
      query_type = "Metrics"
      queries {
        query_group_type = "Successful"
        query_group {
          row_id        = "A"
          query         = "example"
          use_row_count = false
        }
      }
	  queries {
        query_group_type = "Total"
        query_group {
          row_id        = "A"
          query         = "example"
          use_row_count = false
        }
      }
    }
  }
}`, testName)
	return resourceText
}
