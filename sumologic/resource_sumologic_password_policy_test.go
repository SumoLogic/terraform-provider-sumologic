package sumologic

import (
	"fmt"
	"strconv"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

func TestAccPasswordPolicy_create(t *testing.T) {
	passwordPolicy := PasswordPolicy{
		MinLength:                      9,
		MaxLength:                      50,
		MustContainLowercase:           false,
		MustContainUppercase:           true,
		MustContainDigits:              false,
		MustContainSpecialChars:        true,
		MaxPasswordAgeInDays:           364,
		MinUniquePasswords:             4,
		AccountLockoutThreshold:        5,
		FailedLoginResetDurationInMins: 3,
		AccountLockoutDurationInMins:   30,
		RequireMfa:                     false,
		RememberMfa:                    false,
	}

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckPasswordPolicyDestroy(),
		Steps: []resource.TestStep{
			{
				Config: newPasswordPolicyConfig("test", &passwordPolicy),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckPasswordPolicyExists("sumologic_password_policy.test"),
					testPasswordPolicyCheckResourceAttr("sumologic_password_policy.test", &passwordPolicy),
				),
			},
		},
	})
}

func TestAccPasswordPolicy_update(t *testing.T) {
	passwordPolicy := PasswordPolicy{
		MinLength:                      10,
		MaxLength:                      51,
		MustContainLowercase:           true,
		MustContainUppercase:           true,
		MustContainDigits:              false,
		MustContainSpecialChars:        true,
		MaxPasswordAgeInDays:           364,
		MinUniquePasswords:             4,
		AccountLockoutThreshold:        5,
		FailedLoginResetDurationInMins: 3,
		AccountLockoutDurationInMins:   30,
		RequireMfa:                     false,
		RememberMfa:                    false,
	}

	updatedPasswordPolicy := PasswordPolicy{
		MinLength:                      12,
		MaxLength:                      52,
		MustContainLowercase:           false,
		MustContainUppercase:           false,
		MustContainDigits:              false,
		MustContainSpecialChars:        true,
		MaxPasswordAgeInDays:           365,
		MinUniquePasswords:             5,
		AccountLockoutThreshold:        6,
		FailedLoginResetDurationInMins: 5,
		AccountLockoutDurationInMins:   31,
		RequireMfa:                     true,
		RememberMfa:                    true,
	}

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckPasswordPolicyDestroy(),
		Steps: []resource.TestStep{
			{
				Config: newPasswordPolicyConfig("test", &passwordPolicy),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckPasswordPolicyExists("sumologic_password_policy.test"),
					testPasswordPolicyCheckResourceAttr("sumologic_password_policy.test", &passwordPolicy),
				),
			},
			{
				Config: newPasswordPolicyConfig("test", &updatedPasswordPolicy),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckPasswordPolicyExists("sumologic_password_policy.test"),
					testPasswordPolicyCheckResourceAttr("sumologic_password_policy.test", &updatedPasswordPolicy),
				),
			},
		},
	})
}

func testAccCheckPasswordPolicyDestroy() resource.TestCheckFunc {
	// This check will break if we change the default values for password policy
	return func(s *terraform.State) error {
		client := testAccProvider.Meta().(*Client)
		for _, rs := range s.RootModule().Resources {
			if rs.Type != "sumologic_password_policy" {
				continue
			}

			passwordPolicy, err := client.GetPasswordPolicy()
			if err != nil {
				return fmt.Errorf("Encountered an error: " + err.Error())
			}

			if (*passwordPolicy) != newDefaultPasswordPolicy() {
				return fmt.Errorf("Password policy wasn't reset properly")
			}
		}
		return nil
	}
}

func testAccCheckPasswordPolicyExists(name string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		_, ok := s.RootModule().Resources[name]
		if !ok {
			return fmt.Errorf("Password policy not found: %s", name)
		}
		return nil
	}
}

func testPasswordPolicyCheckResourceAttr(resourceName string, passwordPolicy *PasswordPolicy) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		f := resource.ComposeTestCheckFunc(
			resource.TestCheckResourceAttr(resourceName, "min_length", strconv.Itoa(passwordPolicy.MinLength)),
			resource.TestCheckResourceAttr(resourceName, "max_length", strconv.Itoa(passwordPolicy.MaxLength)),
			resource.TestCheckResourceAttr(resourceName, "must_contain_lowercase", strconv.FormatBool(passwordPolicy.MustContainLowercase)),
			resource.TestCheckResourceAttr(resourceName, "must_contain_uppercase", strconv.FormatBool(passwordPolicy.MustContainUppercase)),
			resource.TestCheckResourceAttr(resourceName, "must_contain_digits", strconv.FormatBool(passwordPolicy.MustContainDigits)),
			resource.TestCheckResourceAttr(resourceName, "must_contain_special_chars", strconv.FormatBool(passwordPolicy.MustContainSpecialChars)),
			resource.TestCheckResourceAttr(resourceName, "max_password_age_in_days", strconv.Itoa(passwordPolicy.MaxPasswordAgeInDays)),
			resource.TestCheckResourceAttr(resourceName, "min_unique_passwords", strconv.Itoa(passwordPolicy.MinUniquePasswords)),
			resource.TestCheckResourceAttr(resourceName, "account_lockout_threshold", strconv.Itoa(passwordPolicy.AccountLockoutThreshold)),
			resource.TestCheckResourceAttr(resourceName, "failed_login_reset_duration_in_mins", strconv.Itoa(passwordPolicy.FailedLoginResetDurationInMins)),
			resource.TestCheckResourceAttr(resourceName, "account_lockout_duration_in_mins", strconv.Itoa(passwordPolicy.AccountLockoutDurationInMins)),
			resource.TestCheckResourceAttr(resourceName, "require_mfa", strconv.FormatBool(passwordPolicy.RequireMfa)),
			resource.TestCheckResourceAttr(resourceName, "remember_mfa", strconv.FormatBool(passwordPolicy.RememberMfa)),
		)
		return f(s)
	}
}

func newPasswordPolicyConfig(label string, passwordPolicy *PasswordPolicy) string {
	return fmt.Sprintf(`
resource "sumologic_password_policy" "%s" {
	min_length = %d
	max_length = %d
	must_contain_lowercase = %t
	must_contain_uppercase = %t
	must_contain_digits = %t
	must_contain_special_chars = %t
	max_password_age_in_days = %d
	min_unique_passwords = %d
	account_lockout_threshold = %d
	failed_login_reset_duration_in_mins = %d
	account_lockout_duration_in_mins = %d
	require_mfa = %t
	remember_mfa = %t
}`, label,
		passwordPolicy.MinLength,
		passwordPolicy.MaxLength,
		passwordPolicy.MustContainLowercase,
		passwordPolicy.MustContainUppercase,
		passwordPolicy.MustContainDigits,
		passwordPolicy.MustContainSpecialChars,
		passwordPolicy.MaxPasswordAgeInDays,
		passwordPolicy.MinUniquePasswords,
		passwordPolicy.AccountLockoutThreshold,
		passwordPolicy.FailedLoginResetDurationInMins,
		passwordPolicy.AccountLockoutDurationInMins,
		passwordPolicy.RequireMfa,
		passwordPolicy.RememberMfa)
}

func newDefaultPasswordPolicy() PasswordPolicy {
	return PasswordPolicy{
		MinLength:                      8,
		MaxLength:                      128,
		MustContainLowercase:           true,
		MustContainUppercase:           true,
		MustContainDigits:              true,
		MustContainSpecialChars:        true,
		MaxPasswordAgeInDays:           365,
		MinUniquePasswords:             10,
		AccountLockoutThreshold:        6,
		FailedLoginResetDurationInMins: 10,
		AccountLockoutDurationInMins:   30,
		RequireMfa:                     false,
		RememberMfa:                    true,
	}
}
