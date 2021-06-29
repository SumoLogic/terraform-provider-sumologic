package sumologic

import (
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

var testAccProviders map[string]terraform.ResourceProvider
var testAccProvider *schema.Provider

func init() {
	testAccProvider = Provider().(*schema.Provider)
	testAccProviders = map[string]terraform.ResourceProvider{
		"sumologic": testAccProvider,
	}
}

func TestProvider(t *testing.T) {
	if err := Provider().(*schema.Provider).InternalValidate(); err != nil {
		t.Fatalf("err: %s", err)
	}
}

func TestProvider_impl(t *testing.T) {
	var _ terraform.ResourceProvider = Provider()
}

func testAccPreCheck(t *testing.T) {
	if v := os.Getenv("SUMOLOGIC_ACCESSKEY"); v == "" {
		t.Fatal("SUMOLOGIC_ACCESSKEY must be set for acceptance tests")
	}
	if v := os.Getenv("SUMOLOGIC_ACCESSID"); v == "" {
		t.Fatal("SUMOLOGIC_ACCESSID must be set for acceptance tests")
	}
	if v := os.Getenv("SUMOLOGIC_ENVIRONMENT"); v == "" {
		t.Fatal("SUMOLOGIC_ENVIRONMENT must be set for acceptance tests")
	}
}

func testAccPreCheckWithAWS(t *testing.T) {
	testAccPreCheck(t)
	if v := os.Getenv("SUMOLOGIC_TEST_AWS_ID"); v == "" {
		t.Fatal("SUMOLOGIC_TEST_AWS_ID must be set for polling source acceptance tests")
	}
	if v := os.Getenv("SUMOLOGIC_TEST_AWS_KEY"); v == "" {
		t.Fatal("SUMOLOGIC_TEST_AWS_KEY must be set for polling source acceptance tests")
	}
	if v := os.Getenv("SUMOLOGIC_TEST_BUCKET_NAME"); v == "" {
		t.Fatal("SUMOLOGIC_TEST_BUCKET_NAME must be set for polling source acceptance tests")
	}
}
