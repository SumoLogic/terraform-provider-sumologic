package sumologic

import (
	"fmt"
	"strconv"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

// If the tests fail due to invalid / expired certificate, refer to SUMO-187643.
var testCertificate = `
-----BEGIN CERTIFICATE-----
MIIDqDCCApACCQClhLh7fXVW7TANBgkqhkiG9w0BAQUFADCBlTELMAkGA1UEBhMC
dXMxEzARBgNVBAgMCmNhbGlmb3JuaWExFTATBgNVBAcMDHJlZHdvb2QgY2l0eTEN
MAsGA1UECgwEc3VtbzEUMBIGA1UECwwLZW5naW5lZXJpbmcxDTALBgNVBAMMBHN1
bW8xJjAkBgkqhkiG9w0BCQEWF2VtaWNoYWVsaUBzdW1vbG9naWMuY29tMB4XDTIy
MDMwNzE5MjI0OVoXDTI0MDMwNjE5MjI0OVowgZUxCzAJBgNVBAYTAnVzMRMwEQYD
VQQIDApjYWxpZm9ybmlhMRUwEwYDVQQHDAxyZWR3b29kIGNpdHkxDTALBgNVBAoM
BHN1bW8xFDASBgNVBAsMC2VuZ2luZWVyaW5nMQ0wCwYDVQQDDARzdW1vMSYwJAYJ
KoZIhvcNAQkBFhdlbWljaGFlbGlAc3Vtb2xvZ2ljLmNvbTCCASIwDQYJKoZIhvcN
AQEBBQADggEPADCCAQoCggEBANVb2eHwr63qc7ov/9WSNZskGY3efSsuh+WlM66A
+LBpf7jvpSivB8v3TUaRth8MGGYddQlYug9kA4ruUeF6yVfNe6YdcWT8THFUvaLu
byi1KzqL7nABldHDdbWoskrIf8LK0CEw3qfIS17kBtwv/BfqF97jQ1JMT1ZdFXLO
q4D0Ezy8YOw4MpbpzrkFtseM/m5KS4WjJLWcCEUiAUCGZCsvgOr7WrpaaKmCasMJ
TcIZe6NiA73FChUx0wU0B+T7CBHmWIoQO8vt8YptSBhPF6B4Hp5M7mAzkvMu0DjX
J/DQHNPtGWtZ+PN95qIrMyCYYGK6jyISd/n5/jM/ASVaWXsCAwEAATANBgkqhkiG
9w0BAQUFAAOCAQEAGqxCY9F6S5jHwiU9OUguvM7u0At7wrnSohy2ZiPtN8kQvFOm
SqC7LGqVoCVGkXsnravEqZZ0r5Q3AkTZQzLBPJ5VqdCcbZG19+o+a8pOpyTbhBGz
guxUKwn7ZFJnQ6Ex7/CA1gRffIMborfkmVkrfN3GcAlnJKYIAhFZuf7iReKW/nKc
wpXwzwjwJg8/EecfCo8m7bxc42npc1d+ESqttsgZhJC+okvI11bvRHybJ/B+R69A
2fuOtYcdC3tdkLu+byPf9HWDPlWTbJlHt++qxOw0QPZHoZbL3knLzRQF5PxS5T+h
IaRGjCIgVK9gd3lhPwChyObu8r688zDiL3Cv8Q==
-----END CERTIFICATE-----`

func TestAccSumologicSamlConfiguration_basic(t *testing.T) {
	samlConfiguration := SamlConfiguration{
		SpInitiatedLoginPath:    "",
		ConfigurationName:       "test",
		Issuer:                  "test",
		SpInitiatedLoginEnabled: false,
		AuthnRequestUrl:         "",
		X509cert1:               testCertificate,
		X509cert2:               "",
		X509cert3:               "",
		OnDemandProvisioningEnabled: &OnDemandProvisioningEnabled{
			FirstNameAttribute:        "",
			LastNameAttribute:         "",
			OnDemandProvisioningRoles: []string{"Administrator"},
		},
		RolesAttribute:               "",
		LogoutEnabled:                false,
		LogoutUrl:                    "",
		EmailAttribute:               "",
		DebugMode:                    false,
		SignAuthnRequest:             false,
		DisableRequestedAuthnContext: false,
		IsRedirectBinding:            false,
	}

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckSamlConfigurationDestroy,
		Steps: []resource.TestStep{
			{
				Config: newSamlConfigurationConfig("tf_saml_import_test", &samlConfiguration),
			},
			{
				ResourceName:      "sumologic_saml_configuration.tf_saml_import_test",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccSumologicSamlConfiguration_create(t *testing.T) {
	samlConfiguration := SamlConfiguration{
		SpInitiatedLoginPath:    "",
		ConfigurationName:       "tf_test",
		Issuer:                  "tf_test",
		SpInitiatedLoginEnabled: false,
		AuthnRequestUrl:         "https://myurl.com",
		X509cert1:               testCertificate,
		X509cert2:               "",
		X509cert3:               "",
		OnDemandProvisioningEnabled: &OnDemandProvisioningEnabled{
			FirstNameAttribute:        "testFirstName",
			LastNameAttribute:         "testSecondName",
			OnDemandProvisioningRoles: []string{"Administrator"},
		},
		RolesAttribute:               "roleAttr",
		LogoutEnabled:                false,
		LogoutUrl:                    "https://myurl.com/logout",
		EmailAttribute:               "",
		DebugMode:                    true,
		SignAuthnRequest:             false,
		DisableRequestedAuthnContext: true,
		IsRedirectBinding:            true,
	}

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckSamlConfigurationDestroy,
		Steps: []resource.TestStep{
			{
				Config: newSamlConfigurationConfig("create_test", &samlConfiguration),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSamlConfigurationExists("sumologic_saml_configuration.create_test"),
					testSamlConfigurationCheckResourceAttr("sumologic_saml_configuration.create_test", &samlConfiguration),
				),
			},
		},
	})
}

func TestAccSumologicSamlConfiguration_update(t *testing.T) {
	samlConfiguration := SamlConfiguration{
		SpInitiatedLoginPath:    "",
		ConfigurationName:       "tf_test2",
		Issuer:                  "tf_test2",
		SpInitiatedLoginEnabled: false,
		AuthnRequestUrl:         "https://myurl.com",
		X509cert1:               testCertificate,
		X509cert2:               "",
		X509cert3:               "",
		OnDemandProvisioningEnabled: &OnDemandProvisioningEnabled{
			FirstNameAttribute:        "testFirstName",
			LastNameAttribute:         "testSecondName",
			OnDemandProvisioningRoles: []string{"Administrator"},
		},
		RolesAttribute:               "roleAttr",
		LogoutEnabled:                false,
		LogoutUrl:                    "https://myurl.com/logout",
		EmailAttribute:               "",
		DebugMode:                    false,
		SignAuthnRequest:             false,
		DisableRequestedAuthnContext: false,
		IsRedirectBinding:            false,
	}

	updatedSamlConfigurationConfig := fmt.Sprintf(`
resource "sumologic_saml_configuration" "update_test" {
  configuration_name = "tf_test3"
  issuer = "tf_test3"
  x509cert1 = <<EOT
%s
EOT
}`, testCertificate)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckSamlConfigurationDestroy,
		Steps: []resource.TestStep{
			{
				Config: newSamlConfigurationConfig("update_test", &samlConfiguration),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSamlConfigurationExists("sumologic_saml_configuration.update_test"),
					testSamlConfigurationCheckResourceAttr("sumologic_saml_configuration.update_test", &samlConfiguration),
				),
			},
			{
				Config: updatedSamlConfigurationConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSamlConfigurationExists("sumologic_saml_configuration.update_test"),
					resource.TestCheckResourceAttr("sumologic_saml_configuration.update_test", "configuration_name", "tf_test3"),
					resource.TestCheckResourceAttr("sumologic_saml_configuration.update_test", "issuer", "tf_test3"),
				),
			},
		},
	})
}

func testAccCheckSamlConfigurationDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*Client)
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "sumologic_saml_configuration" {
			continue
		}

		id := rs.Primary.ID
		samlConfiguration, err := client.GetSamlConfiguration(id)
		if err != nil {
			return fmt.Errorf("Encountered an error: " + err.Error())
		}

		if samlConfiguration != nil {
			return fmt.Errorf("Saml Configuration (ID=%s) still exists", id)
		}
	}

	return nil
}

func testAccCheckSamlConfigurationExists(name string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[name]
		if !ok {
			return fmt.Errorf("Error = %s. Saml Configuration not found: %s", strconv.FormatBool(ok), name)
		}

		if strings.EqualFold(rs.Primary.ID, "") {
			return fmt.Errorf("Saml Configuration ID is not set")
		}

		id := rs.Primary.ID
		client := testAccProvider.Meta().(*Client)
		_, err := client.GetSamlConfiguration(id)
		if err != nil {
			return fmt.Errorf("Saml Configuration (id=%s) not found", id)
		}
		assertion_consumer_url := rs.Primary.Attributes["assertion_consumer_url"]
		if strings.EqualFold(assertion_consumer_url, "") {
			return fmt.Errorf("Assertion Consumer URL not found for Saml Configuration (id=%s)", id)
		}
		entity_id := rs.Primary.Attributes["entity_id"]
		if strings.EqualFold(entity_id, "") {
			return fmt.Errorf("Entity Id not set for Saml Configuration (id=%s)", id)
		}
		return nil
	}
}

func testSamlConfigurationCheckResourceAttr(resourceName string, samlConfiguration *SamlConfiguration) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		f := resource.ComposeTestCheckFunc(
			resource.TestCheckResourceAttr(resourceName, "sp_initiated_login_path", samlConfiguration.SpInitiatedLoginPath),
			resource.TestCheckResourceAttr(resourceName, "configuration_name", samlConfiguration.ConfigurationName),
			resource.TestCheckResourceAttr(resourceName, "issuer", samlConfiguration.Issuer),
			resource.TestCheckResourceAttr(resourceName, "sp_initiated_login_enabled", strconv.FormatBool(samlConfiguration.SpInitiatedLoginEnabled)),
			resource.TestCheckResourceAttr(resourceName, "authn_request_url", samlConfiguration.AuthnRequestUrl),
			resource.TestCheckResourceAttrSet(resourceName, "x509cert1"),
			resource.TestCheckResourceAttr(resourceName, "x509cert2", samlConfiguration.X509cert2),
			resource.TestCheckResourceAttr(resourceName, "x509cert3", samlConfiguration.X509cert3),
			resource.TestCheckResourceAttr(resourceName, "on_demand_provisioning_enabled.0.first_name_attribute", samlConfiguration.OnDemandProvisioningEnabled.FirstNameAttribute),
			resource.TestCheckResourceAttr(resourceName, "on_demand_provisioning_enabled.0.last_name_attribute", samlConfiguration.OnDemandProvisioningEnabled.LastNameAttribute),
			resource.TestCheckResourceAttr(resourceName, "on_demand_provisioning_enabled.0.on_demand_provisioning_roles.0", samlConfiguration.OnDemandProvisioningEnabled.OnDemandProvisioningRoles[0]),
			resource.TestCheckResourceAttr(resourceName, "roles_attribute", samlConfiguration.RolesAttribute),
			resource.TestCheckResourceAttr(resourceName, "logout_enabled", strconv.FormatBool(samlConfiguration.LogoutEnabled)),
			resource.TestCheckResourceAttr(resourceName, "logout_url", samlConfiguration.LogoutUrl),
			resource.TestCheckResourceAttr(resourceName, "email_attribute", samlConfiguration.EmailAttribute),
			resource.TestCheckResourceAttr(resourceName, "debug_mode", strconv.FormatBool(samlConfiguration.DebugMode)),
			resource.TestCheckResourceAttr(resourceName, "sign_authn_request", strconv.FormatBool(samlConfiguration.SignAuthnRequest)),
			resource.TestCheckResourceAttr(resourceName, "disable_requested_authn_context", strconv.FormatBool(samlConfiguration.DisableRequestedAuthnContext)),
			resource.TestCheckResourceAttr(resourceName, "is_redirect_binding", strconv.FormatBool(samlConfiguration.IsRedirectBinding)),
		)
		return f(s)
	}
}

func newSamlConfigurationConfig(label string, samlConfiguration *SamlConfiguration) string {
	return fmt.Sprintf(`
resource "sumologic_saml_configuration" "%s" {
  sp_initiated_login_path = "%s"
  configuration_name = "%s"
  issuer = "%s"
  sp_initiated_login_enabled = %t
  authn_request_url = "%s"
  x509cert1 = <<EOT
%s
EOT
  x509cert2 = "%s"
  x509cert3 = "%s"
  on_demand_provisioning_enabled {
    first_name_attribute = "%s"
    last_name_attribute = "%s"
    on_demand_provisioning_roles = ["%s"]
  }
  roles_attribute = "%s"
  logout_enabled = %t
  logout_url = "%s"
  email_attribute = "%s"
  debug_mode = %t
  sign_authn_request = %t
  disable_requested_authn_context = %t
  is_redirect_binding = %t
}`, label,
		samlConfiguration.SpInitiatedLoginPath,
		samlConfiguration.ConfigurationName,
		samlConfiguration.Issuer,
		samlConfiguration.SpInitiatedLoginEnabled,
		samlConfiguration.AuthnRequestUrl,
		samlConfiguration.X509cert1,
		samlConfiguration.X509cert2,
		samlConfiguration.X509cert3,
		samlConfiguration.OnDemandProvisioningEnabled.FirstNameAttribute,
		samlConfiguration.OnDemandProvisioningEnabled.LastNameAttribute,
		samlConfiguration.OnDemandProvisioningEnabled.OnDemandProvisioningRoles[0], // For simplicity, assume only one role
		samlConfiguration.RolesAttribute,
		samlConfiguration.LogoutEnabled,
		samlConfiguration.LogoutUrl,
		samlConfiguration.EmailAttribute,
		samlConfiguration.DebugMode,
		samlConfiguration.SignAuthnRequest,
		samlConfiguration.DisableRequestedAuthnContext,
		samlConfiguration.IsRedirectBinding)
}
