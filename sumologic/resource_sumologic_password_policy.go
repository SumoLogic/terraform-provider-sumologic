package sumologic

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

var DefaultPasswordPolicy = PasswordPolicy{
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

func resourceSumologicPasswordPolicy() *schema.Resource {
	return &schema.Resource{
		Create: resourceSumologicPasswordPolicyCreate,
		Read:   resourceSumologicPasswordPolicyRead,
		Update: resourceSumologicPasswordPolicyUpdate,
		Delete: resourceSumologicPasswordPolicyDelete,

		Schema: map[string]*schema.Schema{
			"min_length": {
				Type:     schema.TypeInt,
				Optional: true,
				Default:  DefaultPasswordPolicy.MinLength,
			},
			"max_length": {
				Type:     schema.TypeInt,
				Optional: true,
				Default:  DefaultPasswordPolicy.MaxLength,
			},
			"must_contain_lowercase": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  DefaultPasswordPolicy.MustContainLowercase,
			},
			"must_contain_uppercase": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  DefaultPasswordPolicy.MustContainUppercase,
			},
			"must_contain_digits": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  DefaultPasswordPolicy.MustContainDigits,
			},
			"must_contain_special_chars": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  DefaultPasswordPolicy.MustContainSpecialChars,
			},
			"max_password_age_in_days": {
				Type:     schema.TypeInt,
				Optional: true,
				Default:  DefaultPasswordPolicy.MaxPasswordAgeInDays,
			},
			"min_unique_passwords": {
				Type:     schema.TypeInt,
				Optional: true,
				Default:  DefaultPasswordPolicy.MinUniquePasswords,
			},
			"account_lockout_threshold": {
				Type:     schema.TypeInt,
				Optional: true,
				Default:  DefaultPasswordPolicy.AccountLockoutThreshold,
			},
			"failed_login_reset_duration_in_mins": {
				Type:     schema.TypeInt,
				Optional: true,
				Default:  DefaultPasswordPolicy.FailedLoginResetDurationInMins,
			},
			"account_lockout_duration_in_mins": {
				Type:     schema.TypeInt,
				Optional: true,
				Default:  DefaultPasswordPolicy.AccountLockoutDurationInMins,
			},
			"require_mfa": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  DefaultPasswordPolicy.RequireMfa,
			},
			"remember_mfa": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  DefaultPasswordPolicy.RememberMfa,
			},
		},
	}
}

func resourceSumologicPasswordPolicyRead(d *schema.ResourceData, meta interface{}) error {
	c := meta.(*Client)

	passwordPolicy, err := c.GetPasswordPolicy()
	if err != nil {
		return err
	}

	setPasswordPolicyResource(d, passwordPolicy)
	return nil
}

func resourceSumologicPasswordPolicyCreate(d *schema.ResourceData, meta interface{}) error {
	// Since password policy can only be set and not created, we just update the password policy with the given fields.
	err := resourceSumologicPasswordPolicyUpdate(d, meta)
	if err != nil {
		return err
	}

	d.SetId("passwordPolicy")
	return nil
}

func resourceSumologicPasswordPolicyDelete(d *schema.ResourceData, meta interface{}) error {
	c := meta.(*Client)
	return c.ResetPasswordPolicy()
}

func resourceSumologicPasswordPolicyUpdate(d *schema.ResourceData, meta interface{}) error {
	passwordPolicy := resourceToPasswordPolicy(d)

	c := meta.(*Client)
	updatedPasswordPolicy, err := c.UpdatePasswordPolicy(passwordPolicy)
	if err != nil {
		return err
	}

	setPasswordPolicyResource(d, updatedPasswordPolicy)
	return nil
}

func resourceToPasswordPolicy(d *schema.ResourceData) PasswordPolicy {
	return PasswordPolicy{
		MinLength:                      d.Get("min_length").(int),
		MaxLength:                      d.Get("max_length").(int),
		MustContainLowercase:           d.Get("must_contain_lowercase").(bool),
		MustContainUppercase:           d.Get("must_contain_uppercase").(bool),
		MustContainDigits:              d.Get("must_contain_digits").(bool),
		MustContainSpecialChars:        d.Get("must_contain_special_chars").(bool),
		MaxPasswordAgeInDays:           d.Get("max_password_age_in_days").(int),
		MinUniquePasswords:             d.Get("min_unique_passwords").(int),
		AccountLockoutThreshold:        d.Get("account_lockout_threshold").(int),
		FailedLoginResetDurationInMins: d.Get("failed_login_reset_duration_in_mins").(int),
		AccountLockoutDurationInMins:   d.Get("account_lockout_duration_in_mins").(int),
		RequireMfa:                     d.Get("require_mfa").(bool),
		RememberMfa:                    d.Get("remember_mfa").(bool),
	}
}

func setPasswordPolicyResource(d *schema.ResourceData, passwordPolicy *PasswordPolicy) {
	d.Set("min_length", passwordPolicy.MinLength)
	d.Set("max_length", passwordPolicy.MaxLength)
	d.Set("must_contain_lowercase", passwordPolicy.MustContainLowercase)
	d.Set("must_contain_uppercase", passwordPolicy.MustContainUppercase)
	d.Set("must_contain_digits", passwordPolicy.MustContainDigits)
	d.Set("must_contain_special_chars", passwordPolicy.MustContainSpecialChars)
	d.Set("max_password_age_in_days", passwordPolicy.MaxPasswordAgeInDays)
	d.Set("min_unique_passwords", passwordPolicy.MinUniquePasswords)
	d.Set("account_lockout_threshold", passwordPolicy.AccountLockoutThreshold)
	d.Set("failed_login_reset_duration_in_mins", passwordPolicy.FailedLoginResetDurationInMins)
	d.Set("account_lockout_duration_in_mins", passwordPolicy.AccountLockoutDurationInMins)
	d.Set("require_mfa", passwordPolicy.RequireMfa)
	d.Set("remember_mfa", passwordPolicy.RememberMfa)
}
