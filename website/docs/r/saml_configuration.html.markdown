---
layout: "sumologic"
page_title: "SumoLogic: sumologic_saml_configuration"
description: |-
  Provides a Sumologic SAML Configuration
---

# sumologic_saml_configuration
Provides a [Sumologic SAML Configuration][1].

## Example Usage
```hcl
resource "sumologic_saml_configuration" "exampleSamlConfiguration" {
  sp_initiated_login_path = ""
  configuration_name = "SumoLogic"
  issuer = "http://www.okta.com/abxcseyuiwelflkdjh"
  sp_initiated_login_enabled = false
  authn_request_url = ""
  x509cert1 = "string"
  x509cert2 = ""
  x509cert3 = ""
  on_demand_provisioning_enabled {
    first_name_attribute = "firstName"
    last_name_attribute = "lastName"
    on_demand_provisioning_roles = ["Administrator"]
  }
  roles_attribute = "Administrator"
  logout_enabled = false
  logout_url = ""
  email_attribute = ""
  debug_mode = false
  sign_authn_request = false
  disable_requested_authn_context = false
  is_redirect_binding = false
}
```

## Argument reference

The following arguments are supported:

- `sp_initiated_login_path` - (Optional) The identifier used to generate a unique URL for user login. Defaults to "".
- `configuration_name` - (Required) Name of the SSO policy or another name used to describe the policy internally.
- `issuer` - (Required) The unique URL assigned to the organization by the SAML Identity Provider.
- `sp_initiated_login_enabled` - (Optional) True if Sumo Logic redirects users to your identity provider with a SAML AuthnRequest when signing in. Defaults to false.
- `authn_request_url` - (Optional) The URL that the identity provider has assigned for Sumo Logic to submit SAML authentication requests to the identity provider. Defaults to "".
- `x509cert1` - (Required) The certificate is used to verify the signature in SAML assertions.
- `x509cert2` - (Optional) The backup certificate used to verify the signature in SAML assertions when x509cert1 expires. Defaults to "".
- `x509cert3` - (Optional) The backup certificate used to verify the signature in SAML assertions when x509cert1 expires and x509cert2 is empty. Defaults to "".
- `on_demand_provisioning_enabled` - (Block List, Max: 1, Optional) The configuration for on-demand provisioning. See [on_demand_provisioning_enabled schema](#schema-for-on_demand_provisioning_enabled) for details.
- `roles_attribute` - (Optional) The role that Sumo Logic will assign to users when they sign in. Defaults to "".
- `logout_enabled` - (Optional) True if users are redirected to a URL after signing out of Sumo Logic. Defaults to false.
- `logout_url` - (Optional) The URL that users will be redirected to after signing out of Sumo Logic. Defaults to "".
- `email_attribute` - (Optional) The email address of the new user account. Defaults to "".
- `debug_mode` - (Optional) True if additional details are included when a user fails to sign in. Defaults to false.
- `sign_authn_request` - (Optional) True if Sumo Logic will send signed Authn requests to the identity provider. Defaults to false.
- `disable_requested_authn_context` - (Optional) True if Sumo Logic will include the RequestedAuthnContext element of the SAML AuthnRequests it sends to the identity provider. Defaults to false.
- `is_redirect_binding` - (Optional) True if the SAML binding is of HTTP Redirect type. Defaults to false.

### Schema for `on_demand_provisioning_enabled`
- `first_name_attribute` - (Optional) First name attribute of the new user account. Defaults to "".
- `last_name_attribute` - (Optional) Last name attribute of the new user account. Defaults to "".
- `on_demand_provisioning_roles` - (Required) List of Sumo Logic RBAC roles to be assigned when user accounts are provisioned.

## Attributes reference

The following attributes are exported:

- `id` - Unique identifier for the SAML Configuration.
- `certificate` - Authentication Request Signing Certificate for the user.
- `assertion_consumer_url` - The URL on Sumo Logic where the IdP will redirect to with its authentication response.

## Import
SAML Configuration can be imported using the SAML configuration id, e.g.:
```hcl
terraform import sumologic_saml_configuration.example 00000000454A5979
```

[1]: https://help.sumologic.com/Manage/Security/SAML/01-Set-Up-SAML-for-Single-Sign-On
