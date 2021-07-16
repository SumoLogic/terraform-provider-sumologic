package sumologic

import (
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceSumologicSamlConfiguration() *schema.Resource {
	return &schema.Resource{
		Create: resourceSumologicSamlConfigurationCreate,
		Read:   resourceSumologicSamlConfigurationRead,
		Update: resourceSumologicSamlConfigurationUpdate,
		Delete: resourceSumologicSamlConfigurationDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"sp_initiated_login_path": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "",
			},
			"configuration_name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"issuer": {
				Type:     schema.TypeString,
				Required: true,
			},
			"sp_initiated_login_enabled": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			"authn_request_url": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "",
			},
			"x509cert1": {
				Type:     schema.TypeString,
				Required: true,
			},
			"x509cert2": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "",
			},
			"x509cert3": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "",
			},
			"on_demand_provisioning_enabled": {
				Type:     schema.TypeList,
				MaxItems: 1,
				Optional: true,
				Elem: &schema.Resource{
					Schema: getOnDemandProvisioningEnabledSchema(),
				},
			},
			"roles_attribute": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "",
			},
			"logout_enabled": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			"logout_url": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "",
			},
			"email_attribute": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "",
			},
			"debug_mode": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			"sign_authn_request": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			"disable_requested_authn_context": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			"is_redirect_binding": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			"certificate": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"assertion_consumer_url": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"entity_id": {
                Type:     schema.TypeString,
                Computed: true,
            },
		},
	}
}

func getOnDemandProvisioningEnabledSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"first_name_attribute": {
			Type:     schema.TypeString,
			Optional: true,
			Default:  "",
		},
		"last_name_attribute": {
			Type:     schema.TypeString,
			Optional: true,
			Default:  "",
		},
		"on_demand_provisioning_roles": {
			Type:     schema.TypeList,
			Required: true,
			Elem: &schema.Schema{
				Type: schema.TypeString,
			},
		},
	}
}

func resourceSumologicSamlConfigurationRead(d *schema.ResourceData, meta interface{}) error {
	c := meta.(*Client)

	id := d.Id()
	samlConfiguration, err := c.GetSamlConfiguration(id)

	if err != nil {
		return err
	}

	if samlConfiguration == nil {
		log.Printf("[WARN] SamlConfiguration not found, removing from state: %v - %v", id, err)
		d.SetId("")
		return nil
	}

	setSamlConfiguration(d, samlConfiguration)
	return nil
}

func resourceSumologicSamlConfigurationCreate(d *schema.ResourceData, meta interface{}) error {
	c := meta.(*Client)

	if d.Id() == "" {
		samlConfiguration := resourceToSamlConfiguration(d)

		createdSamlConfiguration, err := c.CreateSamlConfiguration(samlConfiguration)
		if err != nil {
			return err
		}

		d.SetId(createdSamlConfiguration.ID)
	}

	return resourceSumologicSamlConfigurationRead(d, meta)
}

func resourceSumologicSamlConfigurationDelete(d *schema.ResourceData, meta interface{}) error {
	c := meta.(*Client)

	return c.DeleteSamlConfiguration(d.Id())
}

func resourceSumologicSamlConfigurationUpdate(d *schema.ResourceData, meta interface{}) error {
	c := meta.(*Client)

	samlConfiguration := resourceToSamlConfiguration(d)

	err := c.UpdateSamlConfiguration(d.Id(), samlConfiguration)
	if err != nil {
		return err
	}

	return resourceSumologicSamlConfigurationRead(d, meta)
}

func setSamlConfiguration(d *schema.ResourceData, samlConfiguration *SamlConfiguration) {
	d.Set("sp_initiated_login_path", samlConfiguration.SpInitiatedLoginPath)
	d.Set("configuration_name", samlConfiguration.ConfigurationName)
	d.Set("issuer", samlConfiguration.Issuer)
	d.Set("sp_initiated_login_enabled", samlConfiguration.SpInitiatedLoginEnabled)
	d.Set("authn_request_url", samlConfiguration.AuthnRequestUrl)
	d.Set("x509cert1", samlConfiguration.X509cert1)
	d.Set("x509cert2", samlConfiguration.X509cert2)
	d.Set("x509cert3", samlConfiguration.X509cert3)
	setOnDemandProvisioningEnabled(d, samlConfiguration.OnDemandProvisioningEnabled)
	d.Set("roles_attribute", samlConfiguration.RolesAttribute)
	d.Set("logout_enabled", samlConfiguration.LogoutEnabled)
	d.Set("logout_url", samlConfiguration.LogoutUrl)
	d.Set("email_attribute", samlConfiguration.EmailAttribute)
	d.Set("debug_mode", samlConfiguration.DebugMode)
	d.Set("sign_authn_request", samlConfiguration.SignAuthnRequest)
	d.Set("disable_requested_authn_context", samlConfiguration.DisableRequestedAuthnContext)
	d.Set("is_redirect_binding", samlConfiguration.IsRedirectBinding)
	d.Set("assertion_consumer_url", samlConfiguration.AssertionConsumerUrl)
	d.Set("entity_id", samlConfiguration.EntityId)

	d.Set("certificate", samlConfiguration.Certificate)
}

func getTerraformOnDemandProvisioningEnabled(onDemandProvisioningEnabled *OnDemandProvisioningEnabled) []map[string]interface{} {
	tfOnDemandProvisioningEnabledMap := make(map[string]interface{})
	tfOnDemandProvisioningEnabledMap["first_name_attribute"] = onDemandProvisioningEnabled.FirstNameAttribute
	tfOnDemandProvisioningEnabledMap["last_name_attribute"] = onDemandProvisioningEnabled.LastNameAttribute
	tfOnDemandProvisioningEnabledMap["on_demand_provisioning_roles"] = onDemandProvisioningEnabled.OnDemandProvisioningRoles

	tfOnDemandProvisioningEnabled := make([]map[string]interface{}, 1)
	tfOnDemandProvisioningEnabled[0] = tfOnDemandProvisioningEnabledMap
	return tfOnDemandProvisioningEnabled
}

func resourceToSamlConfiguration(d *schema.ResourceData) SamlConfiguration {
	var samlConfiguration SamlConfiguration
	samlConfiguration.SpInitiatedLoginPath = d.Get("sp_initiated_login_path").(string)
	samlConfiguration.ConfigurationName = d.Get("configuration_name").(string)
	samlConfiguration.Issuer = d.Get("issuer").(string)
	samlConfiguration.SpInitiatedLoginEnabled = d.Get("sp_initiated_login_enabled").(bool)
	samlConfiguration.AuthnRequestUrl = d.Get("authn_request_url").(string)
	samlConfiguration.X509cert1 = d.Get("x509cert1").(string)
	samlConfiguration.X509cert2 = d.Get("x509cert2").(string)
	samlConfiguration.X509cert3 = d.Get("x509cert3").(string)
	if val, ok := d.GetOk("on_demand_provisioning_enabled"); ok {
		obj := val.([]interface{})[0]
		samlConfiguration.OnDemandProvisioningEnabled = getOnDemandProvisioningEnabled(obj.(map[string]interface{}))
	}
	samlConfiguration.RolesAttribute = d.Get("roles_attribute").(string)
	samlConfiguration.LogoutEnabled = d.Get("logout_enabled").(bool)
	samlConfiguration.LogoutUrl = d.Get("logout_url").(string)
	samlConfiguration.EmailAttribute = d.Get("email_attribute").(string)
	samlConfiguration.DebugMode = d.Get("debug_mode").(bool)
	samlConfiguration.SignAuthnRequest = d.Get("sign_authn_request").(bool)
	samlConfiguration.DisableRequestedAuthnContext = d.Get("disable_requested_authn_context").(bool)
	samlConfiguration.IsRedirectBinding = d.Get("is_redirect_binding").(bool)
	return samlConfiguration
}

func setOnDemandProvisioningEnabled(d *schema.ResourceData, obj *OnDemandProvisioningEnabled) {
	// The API responds with an empty OnDemandProvisioningEnabled object even if it wasn't provided in the request.
	// If we set the state with that empty object, then subsequent `terraform plan` will show an update is needed
	// if it's not specified in the configuration. We also can't set an empty OnDemandProvisioningEnabled as the
	// default since it'll result in an invalid request body. For this reason, don't set the state if an empty
	// OnDemandProvisioningEnabled is returned.
	if !isOnDemandProvisioningEnabledNilOrEmpty(obj) {
		onDemandProvisioningEnabled := getTerraformOnDemandProvisioningEnabled(obj)
		if err := d.Set("on_demand_provisioning_enabled", onDemandProvisioningEnabled); err != nil {
			log.Printf("[ERROR] in setting on_demand_provisioning_enabled: %v", err)
		}
	}
}

func isOnDemandProvisioningEnabledNilOrEmpty(obj *OnDemandProvisioningEnabled) bool {
	if obj == nil {
		return true
	}

	return obj.FirstNameAttribute == "" && obj.LastNameAttribute == "" && len(obj.OnDemandProvisioningRoles) == 0
}

func getOnDemandProvisioningEnabled(resourceMap map[string]interface{}) *OnDemandProvisioningEnabled {
	var onDemandProvisioningEnabled OnDemandProvisioningEnabled
	onDemandProvisioningEnabled.FirstNameAttribute = resourceMap["first_name_attribute"].(string)
	onDemandProvisioningEnabled.LastNameAttribute = resourceMap["last_name_attribute"].(string)
	// https://stackoverflow.com/questions/37329246/how-to-convert-string-from-interface-to-string-in-golang
	rolesCollection := resourceMap["on_demand_provisioning_roles"].([]interface{})
	roles := make([]string, len(rolesCollection))
	for i, v := range rolesCollection {
		roles[i] = v.(string)
	}
	onDemandProvisioningEnabled.OnDemandProvisioningRoles = roles
	return &onDemandProvisioningEnabled
}
