package sumologic

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

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
				Default:  8,
			},
			"max_length": {
				Type:     schema.TypeInt,
				Optional: true,
				Default:  128,
			},
			"must_contain_lowercase": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  true,
			},
			"must_contain_uppercase": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  true,
			},
			"must_contain_digits": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  true,
			},
			"must_contain_special_chars": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  true,
			},
			"max_password_age_in_days": {
				Type:     schema.TypeInt,
				Optional: true,
				Default:  365,
			},
			"min_unique_passwords": {
				Type:     schema.TypeInt,
				Optional: true,
				Default:  10,
			},
			"account_lockout_threshold": {
				Type:     schema.TypeInt,
				Optional: true,
				Default:  6,
			},
			"failed_login_reset_duration_in_mins": {
				Type:     schema.TypeInt,
				Optional: true,
				Default:  10,
			},
			"account_lockout_duration_in_mins": {
				Type:     schema.TypeInt,
				Optional: true,
				Default:  30,
			},
			"require_mfa": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			"remember_mfa": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  true,
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
	return c.DeletePasswordPolicy()
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
