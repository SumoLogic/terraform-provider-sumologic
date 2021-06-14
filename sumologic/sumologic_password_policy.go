package sumologic

import (
	"encoding/json"
)

func (s *Client) GetPasswordPolicy() (*PasswordPolicy, error) {
	url := "v1/passwordPolicy"

	data, _, err := s.Get(url, false)
	if err != nil {
		return nil, err
	}

	var passwordPolicy PasswordPolicy
	err = json.Unmarshal(data, &passwordPolicy)
	if err != nil {
		return nil, err
	}

	return &passwordPolicy, nil
}

func (s *Client) ResetPasswordPolicy() error {
	url := "v1/passwordPolicy"

	// Since password policy cannot be deleted, we just reset it back to the default by passing an empty request body.
	data, err := s.Put(url, map[string]string{}, false)
	if err != nil {
		return err
	}

	var resetPasswordPolicy PasswordPolicy
	err = json.Unmarshal(data, &resetPasswordPolicy)
	return err
}

func (s *Client) UpdatePasswordPolicy(passwordPolicy PasswordPolicy) (*PasswordPolicy, error) {
	url := "v1/passwordPolicy"

	data, err := s.Put(url, passwordPolicy, false)
	if err != nil {
		return nil, err
	}

	var updatedPasswordPolicy PasswordPolicy
	err = json.Unmarshal(data, &updatedPasswordPolicy)
	if err != nil {
		return nil, err
	}

	return &updatedPasswordPolicy, err
}

type PasswordPolicy struct {
	MinLength                      int  `json:"minLength"`
	MaxLength                      int  `json:"maxLength"`
	MustContainLowercase           bool `json:"mustContainLowercase"`
	MustContainUppercase           bool `json:"mustContainUppercase"`
	MustContainDigits              bool `json:"mustContainDigits"`
	MustContainSpecialChars        bool `json:"mustContainSpecialChars"`
	MaxPasswordAgeInDays           int  `json:"maxPasswordAgeInDays"`
	MinUniquePasswords             int  `json:"minUniquePasswords"`
	AccountLockoutThreshold        int  `json:"accountLockoutThreshold"`
	FailedLoginResetDurationInMins int  `json:"failedLoginResetDurationInMins"`
	AccountLockoutDurationInMins   int  `json:"accountLockoutDurationInMins"`
	RequireMfa                     bool `json:"requireMfa"`
	RememberMfa                    bool `json:"rememberMfa"`
}
