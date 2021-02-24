---
layout: "sumologic"
page_title: "SumoLogic: sumologic_password_policy"
description: |-
  Sets the Sumologic Password Policy
---

# sumologic_password_policy
Sets the [Sumologic Password Policy][1]. Since there is only a single password policy for an organization,
please ensure that only a single instance of such resource is defined.
The behavior for defining more than one password policy resources is undefined.

## Example Usage
```hcl
resource "sumologic_password_policy" "examplePasswordPolicy" {
  min_length = 8
  max_length = 128
  must_contain_lowercase = true
  must_contain_uppercase = true
  must_contain_digits = true
  must_contain_special_chars = true
  max_password_age_in_days = 365
  min_unique_passwords = 10
  account_lockout_threshold = 6
  failed_login_reset_duration_in_mins = 10
  account_lockout_duration_in_mins = 30
  require_mfa = false
  remember_mfa = true
}
```

## Argument reference

The following arguments are supported:

- `min_length` - (Optional) The minimum length of the password. Defaults to 8.
- `max_length` - (Optional) The maximum length of the password. Defaults to 128.
- `must_contain_lowercase` - (Optional) If the password must contain lower case characters. Defaults to true.
- `must_contain_uppercase` - (Optional) If the password must contain upper case characters. Defaults to true.
- `must_contain_digits` - (Optional) If the password must contain digits. Defaults to true.
- `must_contain_special_chars` - (Optional) If the password must contain special characters. Defaults to true.
- `max_password_age_in_days` - (Optional) Maximum number of days that a password can be used before user is required to change it. Put -1 if the user should not have to change their password. Defaults to 365.
- `min_unique_passwords` - (Optional) The minimum number of unique new passwords that a user must use before an old password can be reused. Defaults to 10.
- `account_lockout_threshold` - (Optional) Number of failed login attempts allowed before account is locked-out. Defaults to 6.
- `failed_login_reset_duration_in_mins` - (Optional) The duration of time in minutes that must elapse from the first failed login attempt after which failed login count is reset to 0. Defaults to 10.
- `account_lockout_duration_in_mins` - (Optional) The duration of time in minutes that a locked-out account remained locked before getting unlocked automatically. Defaults to 30.
- `require_mfa` - (Optional) If MFA should be required to log in. Defaults to false.
- `remember_mfa` - (Optional) If MFA should be remembered on the browser. Defaults to true.

The default values for these arguments are described in [Password Policy API][2].

[1]: https://help.sumologic.com/Manage/Security/Set-the-Password-Policy
[2]: https://api.sumologic.com/docs/#operation/setPasswordPolicy
