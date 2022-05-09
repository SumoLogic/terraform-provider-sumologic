package sumologic

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

func TestAccSumologicMonitorsLibraryFolder_fgpSchemaValidations(t *testing.T) {
	config01 := `
resource "sumologic_monitor_folder" "test_monitorfolder" {
	name        = "terraform_test_monitorfolder"
	description = "terraform_test_monitorfolder_desc"
	obj_permission {
		subject_type = "foo_invalid_subject_type"
		subject_id = "dummyID_01"
		permissions = ["Create","Read","Update","Delete"] 
	}
	obj_permission {
		subject_type = "role"
		subject_id = "dummyID_02"
		permissions = ["Create", "Read"]
	}
}`
	expectedError01 := regexp.MustCompile(
		".*expected obj_permission.0.subject_type to be one of \\[role org], got foo_invalid_subject_type.*")

	resource.Test(t, resource.TestCase{
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckMonitorsLibraryFolderDestroy(),
		Steps: []resource.TestStep{
			{
				Config:      config01,
				PlanOnly:    true,
				ExpectError: expectedError01,
			},
		},
	})

	config02 := `
resource "sumologic_monitor_folder" "test_monitorfolder" {
		name        = "terraform_test_monitorfolder"
		description = "terraform_test_monitorfolder_desc"
		obj_permission {
			subject_type = "role"
			subject_id = "dummyID_01"
			permissions = ["Create","Read","Update","Delete"] 
		}
		obj_permission {
			subject_type = "role"
			subject_id = "dummyID_02"
			permissions = ["Create", "Read", "Invalid_Perm"]
		}
}`
	expectedError02 := regexp.MustCompile(
		".*expected obj_permission.1.permissions.1 to be one of \\[Create Read Update Delete Manage], got Invalid_Perm.*")

	resource.Test(t, resource.TestCase{
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckMonitorsLibraryFolderDestroy(),
		Steps: []resource.TestStep{
			{
				Config:      config02,
				PlanOnly:    true,
				ExpectError: expectedError02,
			},
		},
	})

}

func TestAccSumologicMonitorsLibraryFolder_createWithFGP(t *testing.T) {

	testNameSuffix := acctest.RandString(16)
	tfResourceKey := "sumologic_monitor_folder.test_monitorfolder"
	testName := "terraform_test_monitorfolder_" + testNameSuffix

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckMonitorsLibraryFolderDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccSumologicMonitorsLibraryFolder(testNameSuffix),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckMonitorsLibraryFolderExists(tfResourceKey, t),
					testAccCheckMonitorsLibraryFolderAttributes(tfResourceKey),

					resource.TestCheckResourceAttr("sumologic_monitor_folder.test_monitorfolder", "name", testName),
					resource.TestCheckResourceAttr("sumologic_monitor_folder.test_monitorfolder", "description",
						"terraform_test_monitorfolder_desc"),

					resource.TestCheckResourceAttr("sumologic_monitor_folder.test_monitorfolder",
						"obj_permission.#", "2"),
					testAccCheckMonitorsLibraryFolderFGPBackend(tfResourceKey, t, genExpectedPermStmts),
				),
			},
		},
	})
}

func TestAccSumologicMonitorsLibraryFolder_updateWithFGP(t *testing.T) {

	testNameSuffix := acctest.RandString(16)
	tfResourceKey := "sumologic_monitor_folder.test_monitorfolder"
	testName := "terraform_test_monitorfolder_" + testNameSuffix

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckMonitorsLibraryFolderDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccSumologicMonitorsLibraryFolder(testNameSuffix),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckMonitorsLibraryFolderExists(tfResourceKey, t),
					testAccCheckMonitorsLibraryFolderAttributes(tfResourceKey),

					resource.TestCheckResourceAttr("sumologic_monitor_folder.test_monitorfolder", "name", testName),
					resource.TestCheckResourceAttr("sumologic_monitor_folder.test_monitorfolder", "description",
						"terraform_test_monitorfolder_desc"),

					resource.TestCheckResourceAttr("sumologic_monitor_folder.test_monitorfolder",
						"obj_permission.#", "2"),
					testAccCheckMonitorsLibraryFolderFGPBackend(tfResourceKey, t, genExpectedPermStmts),
				),
			},
			{
				Config: testAccSumologicMonitorsLibraryFolderForUpdate(testNameSuffix),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckMonitorsLibraryFolderExists(tfResourceKey, t),
					testAccCheckMonitorsLibraryFolderAttributes(tfResourceKey),

					resource.TestCheckResourceAttr("sumologic_monitor_folder.test_monitorfolder", "name", testName),
					resource.TestCheckResourceAttr("sumologic_monitor_folder.test_monitorfolder", "description",
						"terraform_test_monitorfolder_desc update test"),
					// updated desc

					resource.TestCheckResourceAttr("sumologic_monitor_folder.test_monitorfolder",
						"obj_permission.#", "1"),
					// 1, instead of 2
					testAccCheckMonitorsLibraryFolderFGPBackend(tfResourceKey, t, genExpectedPermStmtsForUpdate),
				),
			},
		},
	})
}

func TestAccSumologicMonitorsLibraryFolder_driftingCorrectionFGP(t *testing.T) {

	testNameSuffix := acctest.RandString(16)
	tfResourceKey := "sumologic_monitor_folder.test_monitorfolder"
	testName := "terraform_test_monitorfolder_" + testNameSuffix

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckMonitorsLibraryFolderDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccSumologicMonitorsLibraryFolder(testNameSuffix),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckMonitorsLibraryFolderExists(tfResourceKey, t),
					testAccCheckMonitorsLibraryFolderAttributes(tfResourceKey),

					resource.TestCheckResourceAttr("sumologic_monitor_folder.test_monitorfolder", "name", testName),
					resource.TestCheckResourceAttr("sumologic_monitor_folder.test_monitorfolder", "description",
						"terraform_test_monitorfolder_desc"),

					resource.TestCheckResourceAttr("sumologic_monitor_folder.test_monitorfolder",
						"obj_permission.#", "2"),
					testAccCheckMonitorsLibraryFolderFGPBackend(tfResourceKey, t, genExpectedPermStmts),
					// Emulating Drifting at the Backend
					testAccEmulateFGPDrifting(t),
				),
				// "After applying this step and refreshing, the plan was not empty"
				// Non-Empty Plan would occur, after the above step that emulates FGP drifting
				ExpectNonEmptyPlan: true,
			},
			// the following Test Step emulates running "terraform apply" again.
			// This step would detect and correct Drifting
			{
				Config: testAccSumologicMonitorsLibraryFolder(testNameSuffix),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckMonitorsLibraryFolderExists(tfResourceKey, t),
					testAccCheckMonitorsLibraryFolderAttributes(tfResourceKey),

					resource.TestCheckResourceAttr("sumologic_monitor_folder.test_monitorfolder", "name", testName),
					resource.TestCheckResourceAttr("sumologic_monitor_folder.test_monitorfolder", "description",
						"terraform_test_monitorfolder_desc"),

					resource.TestCheckResourceAttr("sumologic_monitor_folder.test_monitorfolder",
						"obj_permission.#", "2"),
					testAccCheckMonitorsLibraryFolderFGPBackend(tfResourceKey, t, genExpectedPermStmts),
				),
			},
		},
	})

}

func testAccCheckMonitorsLibraryFolderDestroy() resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client := testAccProvider.Meta().(*Client)
		for _, r := range s.RootModule().Resources {
			id := r.Primary.ID
			u, err := client.GetMonitorsLibraryFolder(id)
			if err != nil {
				return fmt.Errorf("Encountered an error: " + err.Error())
			}
			if u != nil {
				return fmt.Errorf("monitorsLibraryFolder %s still exists", id)
			}
		}
		return nil
	}
}

func getResourceID(s *terraform.State, name string) (string, error) {
	rs, ok := s.RootModule().Resources[name]
	if !ok {
		return "", fmt.Errorf("getResourceID: ok = %t. Resource not found: %s", ok, name)
	}
	if rs.Primary.ID == "" {
		return "", fmt.Errorf("getResourceID: ID is not set")
	}
	return rs.Primary.ID, nil
}

func testAccCheckMonitorsLibraryFolderExists(name string, t *testing.T) resource.TestCheckFunc {
	return func(s *terraform.State) error {

		id, resIdErr := getResourceID(s, name)
		if resIdErr != nil {
			return resIdErr
		}

		client := testAccProvider.Meta().(*Client)
		folder, err := client.GetMonitorsLibraryFolder(id)
		if err != nil {
			return err
		}
		if folder == nil {
			return fmt.Errorf("MonitorsLibraryFolder %s not found", id)
		}

		return nil
	}
}

func testAccCheckMonitorsLibraryFolderFGPBackend(
	name string,
	t *testing.T,
	expectedFGPFunc func(*terraform.State, string) ([]CmfFgpPermStatement, error),
) resource.TestCheckFunc {

	return func(s *terraform.State) error {
		targetId, resIdErr := getResourceID(s, name)
		if resIdErr != nil {
			return resIdErr
		}

		expectedPermStmts, resIdErr := expectedFGPFunc(s, targetId)
		if resIdErr != nil {
			return resIdErr
		}

		client := testAccProvider.Meta().(*Client)

		fgpResult, fgpErr := client.GetCmfFgp("monitors", targetId)
		if fgpErr != nil {
			return fgpErr
		}

		if !CmfFgpPermStmtSetEqual(fgpResult.PermissionStatements, expectedPermStmts) {
			return fmt.Errorf("Permission Statements are different:\n  %+v\n  %+v\n",
				fgpResult.PermissionStatements, expectedPermStmts)
		}

		return nil
	}
}

func testAccEmulateFGPDrifting(
	t *testing.T,
	// expectedFGPFunc func(*terraform.State, string) ([]CmfFgpPermStatement, error),
) resource.TestCheckFunc {

	return func(s *terraform.State) error {

		folderTargetId, resIdErr := getResourceID(s, "sumologic_monitor_folder.test_monitorfolder")
		if resIdErr != nil {
			return resIdErr
		}
		role01Id, resIdErr := getResourceID(s, "sumologic_role.tf_test_role_01")
		if resIdErr != nil {
			return resIdErr
		}
		role02Id, resIdErr := getResourceID(s, "sumologic_role.tf_test_role_02")
		if resIdErr != nil {
			return resIdErr
		}
		role03Id, resIdErr := getResourceID(s, "sumologic_role.tf_test_role_03")
		if resIdErr != nil {
			return resIdErr
		}

		client := testAccProvider.Meta().(*Client)
		expectedReadPermStmts := []CmfFgpPermStatement{
			{SubjectType: "role", SubjectId: role01Id, TargetId: folderTargetId,
				Permissions: []string{"Create", "Read", "Update"}},
			{SubjectType: "role", SubjectId: role03Id, TargetId: folderTargetId,
				Permissions: []string{"Read"}},
		}
		// using an empty Permissions array to achieve the effect of FGP Revocation
		setFGPPermStmts := append(expectedReadPermStmts,
			CmfFgpPermStatement{SubjectType: "role", SubjectId: role02Id, TargetId: folderTargetId,
				Permissions: []string{}})

		_, setFgpErr := client.SetCmfFgp("monitors", CmfFgpRequest{
			PermissionStatements: setFGPPermStmts})
		if setFgpErr != nil {
			return setFgpErr
		}

		readfgpResult, readFgpErr := client.GetCmfFgp("monitors", folderTargetId)
		if readFgpErr != nil {
			return readFgpErr
		}

		var expectedPermStmts []CmfFgpPermStatement
		expectedPermStmts = append(expectedPermStmts,
			CmfFgpPermStatement{
				SubjectId:   role01Id,
				SubjectType: "role",
				TargetId:    folderTargetId,
				Permissions: []string{"Create", "Read", "Update"},
			},
			CmfFgpPermStatement{
				SubjectId:   role03Id,
				SubjectType: "role",
				TargetId:    folderTargetId,
				Permissions: []string{"Read"},
			},
		)

		if !CmfFgpPermStmtSetEqual(readfgpResult.PermissionStatements, expectedPermStmts) {
			return fmt.Errorf("Permission Statements are different:\n  %+v\n  %+v\n",
				readfgpResult.PermissionStatements, expectedPermStmts)
		}
		return nil
	}
}

func testAccCheckMonitorsLibraryFolderAttributes(name string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		f := resource.ComposeTestCheckFunc(
			resource.TestCheckResourceAttrSet(name, "name"),
			resource.TestCheckResourceAttrSet(name, "parent_id"),
			resource.TestCheckResourceAttrSet(name, "description"),

			resource.TestCheckResourceAttrSet(name, "version"),
			resource.TestCheckResourceAttrSet(name, "is_locked"),
			resource.TestCheckResourceAttrSet(name, "is_mutable"),
			resource.TestCheckResourceAttrSet(name, "is_system"),

			resource.TestCheckResourceAttrSet(name, "modified_at"),
			resource.TestCheckResourceAttrSet(name, "modified_by"),
			resource.TestCheckResourceAttrSet(name, "created_by"),
			resource.TestCheckResourceAttrSet(name, "created_at"),

			resource.TestCheckResourceAttrSet(name, "content_type"),
			resource.TestCheckResourceAttrSet(name, "type"),
		)
		return f(s)
	}
}

func testAccSumologicMonitorsLibraryFolder(testNameSuffix string) string {
	return fmt.Sprintf(`
resource "sumologic_role" "tf_test_role_01" {
	name        = "tf_test_role_01_%s"
	description = "Testing resource sumologic_role"
	capabilities = [
		"viewAlerts",
		"viewMonitorsV2",
		"manageMonitorsV2"
	]
}	
resource "sumologic_role" "tf_test_role_02" {
	name        = "tf_test_role_02_%s"
	description = "Testing resource sumologic_role"
	capabilities = [
		"viewAlerts",
		"viewMonitorsV2",
		"manageMonitorsV2"
	]
}
resource "sumologic_role" "tf_test_role_03" {
	name        = "tf_test_role_03_%s"
	description = "Testing resource sumologic_role"
	capabilities = [
		"viewAlerts",
		"viewMonitorsV2",
		"manageMonitorsV2"
	]
}

resource "sumologic_monitor_folder" "test_monitorfolder" {
		name        = "terraform_test_monitorfolder_%s"
		description = "terraform_test_monitorfolder_desc"
		obj_permission {
		  subject_type = "role"
		  subject_id = sumologic_role.tf_test_role_01.id 
		  permissions = ["Create","Read","Update","Delete"] 
		}
		obj_permission {
		  subject_type = "role"
		  subject_id = sumologic_role.tf_test_role_02.id
		  permissions = ["Create", "Read"]
		}
}
`, testNameSuffix, testNameSuffix, testNameSuffix, testNameSuffix)
}

func testAccSumologicMonitorsLibraryFolderForUpdate(testNameSuffix string) string {
	return fmt.Sprintf(`
resource "sumologic_role" "tf_test_role_01" {
	name        = "tf_test_role_01_%s"
	description = "Testing resource sumologic_role"
	capabilities = [
		"viewAlerts",
		"viewMonitorsV2",
		"manageMonitorsV2"
	]
}	
resource "sumologic_role" "tf_test_role_02" {
	name        = "tf_test_role_02_%s"
	description = "Testing resource sumologic_role"
	capabilities = [
		"viewAlerts",
		"viewMonitorsV2",
		"manageMonitorsV2"
	]
}
resource "sumologic_role" "tf_test_role_03" {
	name        = "tf_test_role_03_%s"
	description = "Testing resource sumologic_role"
	capabilities = [
		"viewAlerts",
		"viewMonitorsV2",
		"manageMonitorsV2"
	]
}

resource "sumologic_monitor_folder" "test_monitorfolder" {
		name        = "terraform_test_monitorfolder_%s"
		description = "terraform_test_monitorfolder_desc update test"
		obj_permission {
		  subject_type = "role"
		  subject_id = sumologic_role.tf_test_role_01.id 
		  permissions = ["Create","Read","Update"] 
		  // "Delete" permission is removed here. 
		}
		// permission for tf_test_role_02 is removed here. 
}
`, testNameSuffix, testNameSuffix, testNameSuffix, testNameSuffix)
}

func genExpectedPermStmts(s *terraform.State, targetId string) ([]CmfFgpPermStatement, error) {
	role01Id, resIdErr := getResourceID(s, "sumologic_role.tf_test_role_01")
	if resIdErr != nil {
		return nil, resIdErr
	}
	role02Id, resIdErr := getResourceID(s, "sumologic_role.tf_test_role_02")
	if resIdErr != nil {
		return nil, resIdErr
	}

	var expectedPermStmts []CmfFgpPermStatement
	expectedPermStmts = append(expectedPermStmts,
		CmfFgpPermStatement{
			SubjectId:   role01Id,
			SubjectType: "role",
			TargetId:    targetId,
			Permissions: []string{"Create", "Read", "Update", "Delete"},
		},
		CmfFgpPermStatement{
			SubjectId:   role02Id,
			SubjectType: "role",
			TargetId:    targetId,
			Permissions: []string{"Create", "Read"},
		},
	)
	return expectedPermStmts, nil
}

func genExpectedPermStmtsForUpdate(s *terraform.State, targetId string) ([]CmfFgpPermStatement, error) {
	role01Id, resIdErr := getResourceID(s, "sumologic_role.tf_test_role_01")
	if resIdErr != nil {
		return nil, resIdErr
	}

	var expectedPermStmts []CmfFgpPermStatement
	expectedPermStmts = append(expectedPermStmts,
		CmfFgpPermStatement{
			SubjectId:   role01Id,
			SubjectType: "role",
			TargetId:    targetId,
			Permissions: []string{"Create", "Read", "Update"},
		},
	)
	return expectedPermStmts, nil
}
