package sumologic

import (
	"os"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

var testAccProviders map[string]*schema.Provider
var testAccProvider *schema.Provider

func init() {
	testAccProvider = Provider()
	testAccProviders = map[string]*schema.Provider{
		"sumologic": testAccProvider,
	}
}

func TestProvider(t *testing.T) {
	if err := Provider().InternalValidate(); err != nil {
		t.Fatalf("err: %s", err)
	}
}

func TestProvider_impl(t *testing.T) {
	var _ *schema.Provider = Provider()
}

func SkipCseTest(t *testing.T) {
	if strings.ToLower(os.Getenv("SKIP_CSE_TESTS")) == "true" {
		t.Skip("Skipping CSE Test")
	}
}

func testAccPreCheck(t *testing.T) {
	if os.Getenv("SUMOLOGIC_ACCESSKEY") == "" {
		t.Fatal("SUMOLOGIC_ACCESSKEY must be set for acceptance tests")
	}
	if os.Getenv("SUMOLOGIC_ACCESSID") == "" {
		t.Fatal("SUMOLOGIC_ACCESSID must be set for acceptance tests")
	}
	if os.Getenv("SUMOLOGIC_ENVIRONMENT") == "" && os.Getenv("SUMOLOGIC_BASE_URL") == "" {
		t.Fatal("SUMOLOGIC_ENVIRONMENT must be set for acceptance tests")
	}
}

func testAccPreCheckWithAWS(t *testing.T) {
	testAccPreCheck(t)
	if v := os.Getenv("SUMOLOGIC_TEST_ROLE_ARN"); v == "" {
		t.Fatal("SUMOLOGIC_TEST_ROLE_ARN must be set for polling source acceptance tests")
	}
	if v := os.Getenv("SUMOLOGIC_TEST_BUCKET_NAME"); v == "" {
		t.Fatal("SUMOLOGIC_TEST_BUCKET_NAME must be set for polling source acceptance tests")
	}
}
