package sumologic

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

func TestAccSumologicMutingSchedulesLibraryMutingSchedule_basic(t *testing.T) {
	var mutingSchedulesLibraryMutingSchedule MutingSchedulesLibraryMutingSchedule
	testNameSuffix := acctest.RandString(16)

	testName := "terraform_test_muting_schedule_" + testNameSuffix

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckMutingSchedulesLibraryMutingScheduleDestroy(mutingSchedulesLibraryMutingSchedule),
		Steps: []resource.TestStep{
			{
				Config: testAccSumologicMutingSchedulesLibraryMutingSchedule(testName),
			},
			{
				ResourceName:      "sumologic_muting_schedule.test",
				ImportState:       true,
				ImportStateVerify: false,
			},
		},
	})
}
func TestAccSumologicMutingSchedulesLibraryMutingSchedule_create(t *testing.T) {
	var mutingSchedulesLibraryMutingSchedule MutingSchedulesLibraryMutingSchedule
	testNameSuffix := acctest.RandString(16)
	tomorrow := time.Now().AddDate(0, 0, 1)

	testName := "terraform_test_muting_schedule_" + testNameSuffix
	testDescription := "terraform_test_muting_schedule_description"
	testType := "MutingSchedulesLibraryMutingSchedule"
	testContentType := "MutingSchedule"
	testMonitor := MonitorScope{
		All: true,
	}
	testSchedule := ScheduleDefinition{
		TimeZone:  "America/Los_Angeles",
		StartDate: tomorrow.Format("2006-01-02"),
		StartTime: "00:00",
		Duration:  40,
		RRule:     "FREQ=DAILY;INTERVAL=1;BYHOUR=9,10",
	}
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckMutingSchedulesLibraryMutingScheduleDestroy(mutingSchedulesLibraryMutingSchedule),
		Steps: []resource.TestStep{
			{
				Config: testAccSumologicMutingSchedulesLibraryMutingSchedule(testNameSuffix),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckMutingSchedulesLibraryMutingScheduleExists("sumologic_muting_schedule.test", &mutingSchedulesLibraryMutingSchedule, t),
					testAccCheckMutingSchedulesLibraryMutingScheduleAttributes("sumologic_muting_schedule.test"),
					resource.TestCheckResourceAttr("sumologic_muting_schedule.test", "name", testName),
					resource.TestCheckResourceAttr("sumologic_muting_schedule.test", "type", testType),
					resource.TestCheckResourceAttr("sumologic_muting_schedule.test", "description", testDescription),
					resource.TestCheckResourceAttr("sumologic_muting_schedule.test", "content_type", testContentType),
					resource.TestCheckResourceAttr("sumologic_muting_schedule.test", "monitor.0.all", strconv.FormatBool(testMonitor.All)),
					resource.TestCheckResourceAttr("sumologic_muting_schedule.test", "schedule.0.timezone", testSchedule.TimeZone),
					resource.TestCheckResourceAttr("sumologic_muting_schedule.test", "schedule.0.start_date", testSchedule.StartDate),
					resource.TestCheckResourceAttr("sumologic_muting_schedule.test", "schedule.0.start_time", testSchedule.StartTime),
					resource.TestCheckResourceAttr("sumologic_muting_schedule.test", "schedule.0.rrule", testSchedule.RRule),
				),
			},
		},
	})
}

func TestAccSumologicMutingSchedulesLibraryMutingScheduleWithNotificationGroup_create(t *testing.T) {
	var mutingSchedulesLibraryMutingSchedule MutingSchedulesLibraryMutingSchedule
	testNameSuffix := acctest.RandString(16)
	tomorrow := time.Now().AddDate(0, 0, 1)

	testName := "terraform_test_muting_schedule_" + testNameSuffix
	testDescription := "terraform_test_muting_schedule_description"
	testType := "MutingSchedulesLibraryMutingSchedule"
	testContentType := "MutingSchedule"
	testMonitor := MonitorScope{
		All: true,
	}
	testSchedule := ScheduleDefinition{
		TimeZone:  "America/Los_Angeles",
		StartDate: tomorrow.Format("2006-01-02"),
		StartTime: "00:00",
		Duration:  40,
		RRule:     "FREQ=DAILY;INTERVAL=1;BYHOUR=9,10",
	}
	testNotificationGroup := []NotificationGroupDefinition{{
		GroupKey:    "host",
		GroupValues: []string{"localhost", "127.0.0.1"},
	}}
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckMutingSchedulesLibraryMutingScheduleDestroy(mutingSchedulesLibraryMutingSchedule),
		Steps: []resource.TestStep{
			{
				Config: testAccSumologicMutingSchedulesLibraryMutingScheduleWithNotificationGroups(testNameSuffix),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckMutingSchedulesLibraryMutingScheduleExists("sumologic_muting_schedule.test", &mutingSchedulesLibraryMutingSchedule, t),
					testAccCheckMutingSchedulesLibraryMutingScheduleAttributes("sumologic_muting_schedule.test"),
					resource.TestCheckResourceAttr("sumologic_muting_schedule.test", "name", testName),
					resource.TestCheckResourceAttr("sumologic_muting_schedule.test", "type", testType),
					resource.TestCheckResourceAttr("sumologic_muting_schedule.test", "description", testDescription),
					resource.TestCheckResourceAttr("sumologic_muting_schedule.test", "content_type", testContentType),
					resource.TestCheckResourceAttr("sumologic_muting_schedule.test", "monitor.0.all", strconv.FormatBool(testMonitor.All)),
					resource.TestCheckResourceAttr("sumologic_muting_schedule.test", "schedule.0.timezone", testSchedule.TimeZone),
					resource.TestCheckResourceAttr("sumologic_muting_schedule.test", "schedule.0.start_date", testSchedule.StartDate),
					resource.TestCheckResourceAttr("sumologic_muting_schedule.test", "schedule.0.start_time", testSchedule.StartTime),
					resource.TestCheckResourceAttr("sumologic_muting_schedule.test", "schedule.0.rrule", testSchedule.RRule),
					resource.TestCheckResourceAttr("sumologic_muting_schedule.test", "notification_groups.0.group_key", testNotificationGroup[0].GroupKey),
					resource.TestCheckResourceAttr("sumologic_muting_schedule.test", "notification_groups.0.group_values.0", testNotificationGroup[0].GroupValues[0]),
					resource.TestCheckResourceAttr("sumologic_muting_schedule.test", "notification_groups.0.group_values.1", testNotificationGroup[0].GroupValues[1]),
				),
			},
		},
	})
}

func TestAccSumologicMutingSchedulesLibraryMutingSchedule_update(t *testing.T) {
	var mutingSchedulesLibraryMutingSchedule MutingSchedulesLibraryMutingSchedule
	testNameSuffix := acctest.RandString(16)
	tomorrow := time.Now().AddDate(0, 0, 1)
	testName := "terraform_test_muting_schedule_" + testNameSuffix
	testDescription := "terraform_test_muting_schedule_description"
	testType := "MutingSchedulesLibraryMutingSchedule"
	testContentType := "MutingSchedule"

	testMonitor := MonitorScope{
		All: true,
	}
	testSchedule := ScheduleDefinition{
		TimeZone:  "America/Los_Angeles",
		StartDate: tomorrow.Format("2006-01-02"),
		StartTime: "00:00",
		Duration:  40,
		RRule:     "FREQ=DAILY;INTERVAL=1;BYHOUR=9,10",
	}

	// updated fields
	testUpdatedName := "terraform_test_muting_schedule_" + testNameSuffix
	testUpdatedDescription := "terraform_test_muting_schedule_description"
	testUpdatedType := "MutingSchedulesLibraryMutingSchedule"
	testUpdatedContentType := "MutingSchedule"
	testUpdateMonitor := MonitorScope{
		All: true,
	}
	testUpdateSchedule := ScheduleDefinition{
		TimeZone:  "America/Los_Angeles",
		StartDate: tomorrow.Format("2006-01-02"),
		StartTime: "01:00",
		Duration:  50,
		RRule:     "FREQ=DAILY;INTERVAL=1",
	}

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckMutingSchedulesLibraryMutingScheduleDestroy(mutingSchedulesLibraryMutingSchedule),
		Steps: []resource.TestStep{
			{
				Config: testAccSumologicMutingSchedulesLibraryMutingSchedule(testNameSuffix),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckMutingSchedulesLibraryMutingScheduleExists("sumologic_muting_schedule.test", &mutingSchedulesLibraryMutingSchedule, t),
					testAccCheckMutingSchedulesLibraryMutingScheduleAttributes("sumologic_muting_schedule.test"),
					resource.TestCheckResourceAttr("sumologic_muting_schedule.test", "name", testName),
					resource.TestCheckResourceAttr("sumologic_muting_schedule.test", "type", testType),
					resource.TestCheckResourceAttr("sumologic_muting_schedule.test", "description", testDescription),
					resource.TestCheckResourceAttr("sumologic_muting_schedule.test", "content_type", testContentType),
					resource.TestCheckResourceAttr("sumologic_muting_schedule.test", "monitor.0.all", strconv.FormatBool(testMonitor.All)),
					resource.TestCheckResourceAttr("sumologic_muting_schedule.test", "schedule.0.timezone", testSchedule.TimeZone),
					resource.TestCheckResourceAttr("sumologic_muting_schedule.test", "schedule.0.start_date", testSchedule.StartDate),
					resource.TestCheckResourceAttr("sumologic_muting_schedule.test", "schedule.0.start_time", testSchedule.StartTime),
					resource.TestCheckResourceAttr("sumologic_muting_schedule.test", "schedule.0.rrule", testSchedule.RRule),
				),
			},
			{
				Config: testAccSumologicMutingSchedulesLibraryMutingScheduleUpdate(testNameSuffix),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("sumologic_muting_schedule.test", "name", testUpdatedName),
					resource.TestCheckResourceAttr("sumologic_muting_schedule.test", "type", testUpdatedType),
					resource.TestCheckResourceAttr("sumologic_muting_schedule.test", "description", testUpdatedDescription),
					resource.TestCheckResourceAttr("sumologic_muting_schedule.test", "content_type", testUpdatedContentType),
					resource.TestCheckResourceAttr("sumologic_muting_schedule.test", "monitor.0.all", strconv.FormatBool(testUpdateMonitor.All)),
					resource.TestCheckResourceAttr("sumologic_muting_schedule.test", "schedule.0.timezone", testUpdateSchedule.TimeZone),
					resource.TestCheckResourceAttr("sumologic_muting_schedule.test", "schedule.0.start_date", testUpdateSchedule.StartDate),
					resource.TestCheckResourceAttr("sumologic_muting_schedule.test", "schedule.0.start_time", testUpdateSchedule.StartTime),
					resource.TestCheckResourceAttr("sumologic_muting_schedule.test", "schedule.0.duration", strconv.FormatInt((int64)(testUpdateSchedule.Duration), 10)),
					resource.TestCheckResourceAttr("sumologic_muting_schedule.test", "schedule.0.rrule", testUpdateSchedule.RRule),
				),
			},
		},
	})
}

func TestAccSumologicMutingSchedulesLibraryMutingSchedule_monitorScopeValidations(t *testing.T) {
	var mutingSchedulesLibraryMutingSchedule MutingSchedulesLibraryMutingSchedule
	testNameSuffix := acctest.RandString(16)
	testName := "terraform_test_muting_schedule_" + testNameSuffix
	config := testAccSumologicMutingSchedulesLibraryMutingScheduleBadMonitorScope(testName)
	expectedError := regexp.MustCompile("An argument named \"abc\" is not expected here. Did you mean \"all\"?")
	resource.Test(t, resource.TestCase{
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckMutingSchedulesLibraryMutingScheduleDestroy(mutingSchedulesLibraryMutingSchedule),
		Steps: []resource.TestStep{
			{
				Config:      config,
				PlanOnly:    true,
				ExpectError: expectedError,
			},
		},
	})
}

func testAccCheckMutingSchedulesLibraryMutingScheduleDestroy(mutingSchedulesLibraryMutingSchedule MutingSchedulesLibraryMutingSchedule) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client := testAccProvider.Meta().(*Client)
		for _, r := range s.RootModule().Resources {
			id := r.Primary.ID
			u, err := client.MutingSchedulesRead(id)
			if err != nil {
				return fmt.Errorf("Encountered an error: " + err.Error())
			}
			if u != nil {
				return fmt.Errorf("MutingSchedulesLibraryMutingSchedule %s still exists", id)
			}
		}
		return nil
	}
}

func testAccCheckMutingSchedulesLibraryMutingScheduleExists(name string, mutingSchedulesLibraryMutingSchedule *MutingSchedulesLibraryMutingSchedule, t *testing.T) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[name]
		if !ok {
			return fmt.Errorf("Error = %s. MutingSchedulesLibraryMutingSchedule not found: %s", strconv.FormatBool(ok), name)
		}
		if strings.EqualFold(rs.Primary.ID, "") {
			return fmt.Errorf("MutingSchedulesLibraryMutingSchedule ID is not set")
		}

		id := rs.Primary.ID
		client := testAccProvider.Meta().(*Client)
		newMutingSchedulesLibraryMutingSchedule, err := client.MutingSchedulesRead(id)
		if err != nil {
			return fmt.Errorf("MutingSchedulesLibraryMutingSchedule %s not found", id)
		}
		mutingSchedulesLibraryMutingSchedule = newMutingSchedulesLibraryMutingSchedule
		return nil
	}
}

func testAccCheckMutingSchedulesLibraryMutingScheduleAttributes(name string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		f := resource.ComposeTestCheckFunc(
			resource.TestCheckResourceAttrSet(name, "modified_at"),
			resource.TestCheckResourceAttrSet(name, "created_by"),
			resource.TestCheckResourceAttrSet(name, "name"),
			resource.TestCheckResourceAttrSet(name, "parent_id"),
			resource.TestCheckResourceAttrSet(name, "is_system"),
			resource.TestCheckResourceAttrSet(name, "is_mutable"),
			resource.TestCheckResourceAttrSet(name, "type"),
			resource.TestCheckResourceAttrSet(name, "version"),
			resource.TestCheckResourceAttrSet(name, "description"),
			resource.TestCheckResourceAttrSet(name, "modified_by"),
			resource.TestCheckResourceAttrSet(name, "created_at"),
			resource.TestCheckResourceAttrSet(name, "content_type"),
		)
		return f(s)
	}
}

func testAccSumologicMutingSchedulesLibraryMutingSchedule(testName string) string {
	tomorrow := time.Now().AddDate(0, 0, 1)
	startDate := tomorrow.Format("2006-01-02")
	return fmt.Sprintf(`
    resource "sumologic_muting_schedule" "test" {
	name = "terraform_test_muting_schedule_%s"
	description = "terraform_test_muting_schedule_description"
	type = "MutingSchedulesLibraryMutingSchedule"
	content_type = "MutingSchedule"
	monitor {
		ids = []
		all = true
	  }
	schedule  {
		timezone = "America/Los_Angeles"
		start_date = "%s"
		start_time = "00:00"
		duration = 40
		rrule = "FREQ=DAILY;INTERVAL=1;BYHOUR=9,10"
	  }
}
`, testName, startDate)
}

func testAccSumologicMutingSchedulesLibraryMutingScheduleUpdate(testName string) string {
	tomorrow := time.Now().AddDate(0, 0, 1)
	startDate := tomorrow.Format("2006-01-02")
	return fmt.Sprintf(`
   resource "sumologic_muting_schedule" "test" {
	name = "terraform_test_muting_schedule_%s"
	description = "terraform_test_muting_schedule_description"
	type = "MutingSchedulesLibraryMutingSchedule"
	content_type = "MutingSchedule"
	monitor {
		ids = []
		all = true
	  }
	schedule  {
	timezone = "America/Los_Angeles"
	start_date = "%s"
	start_time = "01:00"
	duration = 50
	rrule = "FREQ=DAILY;INTERVAL=1"
	}
}
`, testName, startDate)
}

func testAccSumologicMutingSchedulesLibraryMutingScheduleBadMonitorScope(testName string) string {
	tomorrow := time.Now().AddDate(0, 0, 1)
	startDate := tomorrow.Format("2006-01-02")
	return fmt.Sprintf(`
   resource "sumologic_muting_schedule" "test" {
	name = "terraform_test_muting_schedule_%s"
	description = "terraform_test_muting_schedule_description"
	type = "MutingSchedulesLibraryMutingSchedule"
	content_type = "MutingSchedule"
	monitor {
		abc="not right"
	  }
	schedule  {
	timezone = "America/Los_Angeles"
	start_date = "%s"
	start_time = "01:00"
	duration = 50
	rrule = "FREQ=DAILY;INTERVAL=1"
	}
}
`, testName, startDate)
}

func testAccSumologicMutingSchedulesLibraryMutingScheduleWithNotificationGroups(testName string) string {
	tomorrow := time.Now().AddDate(0, 0, 1)
	startDate := tomorrow.Format("2006-01-02")
	return fmt.Sprintf(`
    resource "sumologic_muting_schedule" "test" {
	name = "terraform_test_muting_schedule_%s"
	description = "terraform_test_muting_schedule_description"
	type = "MutingSchedulesLibraryMutingSchedule"
	content_type = "MutingSchedule"
	monitor {
		ids = []
		all = true
	  }
	schedule  {
		timezone = "America/Los_Angeles"
		start_date = "%s"
		start_time = "00:00"
		duration = 40
		rrule = "FREQ=DAILY;INTERVAL=1;BYHOUR=9,10"
	  }
	notification_groups {
		group_key = "host"
		group_values =["localhost","127.0.0.1"]
	}  
}
`, testName, startDate)
}
