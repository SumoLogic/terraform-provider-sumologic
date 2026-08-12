package sumologic

import (
	"context"
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/customdiff"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func resourceSumologicDataArchivingDestination() *schema.Resource {
	return &schema.Resource{
		Create: resourceSumologicDataArchivingDestinationCreate,
		Read:   resourceSumologicDataArchivingDestinationRead,
		Update: resourceSumologicDataArchivingDestinationUpdate,
		Delete: resourceSumologicDataArchivingDestinationDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		CustomizeDiff: customdiff.All(
			validateDataArchivingDestinationConfig,
		),

		Schema: map[string]*schema.Schema{
			"destination_name": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringLenBetween(1, 128),
			},
			"destination_config": {
				Type:     schema.TypeList,
				Required: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: dataArchivingDestinationConfigSchema(),
				},
			},
			"created_at": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"created_by": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"modified_at": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"modified_by": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func dataArchivingDestinationConfigSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"destination_type": {
			Type:         schema.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringInSlice([]string{"S3", "Syslog", "Hitachi", "RestAPI"}, false),
		},
		"description": {
			Type:     schema.TypeString,
			Optional: true,
		},
		"bucket_name": {
			Type:     schema.TypeString,
			Optional: true,
			ForceNew: true,
		},
		"region": {
			Type:     schema.TypeString,
			Optional: true,
		},
		"encrypted": {
			Type:     schema.TypeBool,
			Optional: true,
		},
		"enabled": {
			Type:     schema.TypeBool,
			Optional: true,
		},
		"invalidated_by_system": {
			Type:     schema.TypeBool,
			Computed: true,
		},
		"auth_config": {
			Type:     schema.TypeList,
			Optional: true,
			MaxItems: 1,
			Elem: &schema.Resource{
				Schema: s3ArchivingAuthConfigSchema(),
			},
		},
		"protocol": {
			Type:         schema.TypeString,
			Optional:     true,
			ValidateFunc: validation.StringInSlice([]string{"tcp", "udp"}, false),
		},
		"host": {
			Type:     schema.TypeString,
			Optional: true,
		},
		"port": {
			Type:         schema.TypeInt,
			Optional:     true,
			ValidateFunc: validation.IntBetween(1, 65535),
		},
		"token": {
			Type:      schema.TypeString,
			Optional:  true,
			Sensitive: true,
		},
		"url": {
			Type:     schema.TypeString,
			Optional: true,
		},
		"object_id": {
			Type:     schema.TypeString,
			Optional: true,
		},
		"username": {
			Type:     schema.TypeString,
			Optional: true,
		},
		"password": {
			Type:      schema.TypeString,
			Optional:  true,
			Sensitive: true,
		},
	}
}

func s3ArchivingAuthConfigSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"authentication_mode": {
			Type:         schema.TypeString,
			Required:     true,
			ValidateFunc: validation.StringInSlice([]string{"AccessKey", "RoleBased"}, false),
		},
		"access_key_id": {
			Type:     schema.TypeString,
			Optional: true,
		},
		"access_key_secret": {
			Type:      schema.TypeString,
			Optional:  true,
			Sensitive: true,
		},
		"role_arn": {
			Type:     schema.TypeString,
			Optional: true,
		},
	}
}

// An omitted bool is indistinguishable from false through d.Get, so presence of the
// S3 booleans has to be read from the raw config. Anything undeterminable (unknown
// values, destroy plans) counts as set so that planning never fails spuriously.
func destinationConfigAttrIsSet(d *schema.ResourceDiff, attr string) bool {
	raw := d.GetRawConfig()
	if !raw.IsKnown() || raw.IsNull() {
		return true
	}
	blocks := raw.GetAttr("destination_config")
	if !blocks.IsKnown() || blocks.IsNull() {
		return true
	}
	items := blocks.AsValueSlice()
	if len(items) == 0 || items[0].IsNull() {
		return true
	}
	val := items[0].GetAttr(attr)
	return !val.IsKnown() || !val.IsNull()
}

func validateDataArchivingDestinationConfig(_ context.Context, d *schema.ResourceDiff, _ interface{}) error {
	configs := d.Get("destination_config").([]interface{})
	if len(configs) == 0 {
		return nil
	}
	cfg := configs[0].(map[string]interface{})
	destType := cfg["destination_type"].(string)

	switch destType {
	case "S3":
		if d.Id() == "" {
			if cfg["bucket_name"].(string) == "" {
				return fmt.Errorf("bucket_name is required for S3 destination on create")
			}
			if cfg["region"].(string) == "" {
				return fmt.Errorf("region is required for S3 destination on create")
			}
		}
		if !destinationConfigAttrIsSet(d, "encrypted") {
			return fmt.Errorf("encrypted is required for S3 destination")
		}
		if !destinationConfigAttrIsSet(d, "enabled") {
			return fmt.Errorf("enabled is required for S3 destination")
		}
		authConfigs := cfg["auth_config"].([]interface{})
		if len(authConfigs) == 0 {
			return fmt.Errorf("auth_config is required for S3 destination")
		}
	case "Syslog":
		if cfg["protocol"].(string) == "" {
			return fmt.Errorf("protocol is required for Syslog destination")
		}
		if cfg["host"].(string) == "" {
			return fmt.Errorf("host is required for Syslog destination")
		}
		if cfg["port"].(int) == 0 {
			return fmt.Errorf("port is required for Syslog destination")
		}
	case "Hitachi":
		if cfg["url"].(string) == "" {
			return fmt.Errorf("url is required for Hitachi destination")
		}
		if cfg["object_id"].(string) == "" {
			return fmt.Errorf("object_id is required for Hitachi destination")
		}
		if cfg["username"].(string) == "" {
			return fmt.Errorf("username is required for Hitachi destination")
		}
		if d.Id() == "" && cfg["password"].(string) == "" {
			return fmt.Errorf("password is required for Hitachi destination on create")
		}
	case "RestAPI":
		if cfg["url"].(string) == "" {
			return fmt.Errorf("url is required for RestAPI destination")
		}
		// Only the update contract requires a username, so a destination created
		// without one cannot be updated afterwards.
		if d.Id() != "" && cfg["username"].(string) == "" {
			return fmt.Errorf("username is required to update a RestAPI destination")
		}
	}
	return nil
}

func resourceSumologicDataArchivingDestinationCreate(d *schema.ResourceData, meta interface{}) error {
	c := meta.(*Client)

	if d.Id() == "" {
		dest := expandDataArchivingDestination(d)
		created, err := c.CreateDataArchivingDestination(dest)
		if err != nil {
			return err
		}
		d.SetId(created.ID)
	}

	return resourceSumologicDataArchivingDestinationRead(d, meta)
}

func resourceSumologicDataArchivingDestinationRead(d *schema.ResourceData, meta interface{}) error {
	c := meta.(*Client)

	dest, err := c.GetDataArchivingDestination(d.Id())
	if err != nil {
		return err
	}
	if dest == nil {
		log.Printf("[WARN] Data archiving destination not found (id=%s), removing from state", d.Id())
		d.SetId("")
		return nil
	}

	d.Set("destination_name", dest.DestinationName)
	d.Set("created_at", dest.CreatedAt)
	d.Set("created_by", dest.CreatedBy)
	d.Set("modified_at", dest.ModifiedAt)
	d.Set("modified_by", dest.ModifiedBy)

	if err := d.Set("destination_config", flattenDataArchivingDestinationConfig(dest.DestinationConfig, d)); err != nil {
		return err
	}

	return nil
}

func resourceSumologicDataArchivingDestinationUpdate(d *schema.ResourceData, meta interface{}) error {
	c := meta.(*Client)
	dest := expandDataArchivingDestination(d)

	err := c.UpdateDataArchivingDestination(dest)
	if err != nil {
		return err
	}

	return resourceSumologicDataArchivingDestinationRead(d, meta)
}

func resourceSumologicDataArchivingDestinationDelete(d *schema.ResourceData, meta interface{}) error {
	c := meta.(*Client)
	return c.DeleteDataArchivingDestination(d.Id())
}

func expandDataArchivingDestination(d *schema.ResourceData) DataArchivingDestination {
	configs := d.Get("destination_config").([]interface{})
	cfg := configs[0].(map[string]interface{})

	config := DataArchivingDestinationConfig{
		DestinationType: cfg["destination_type"].(string),
	}

	switch config.DestinationType {
	case "S3":
		config.Description = cfg["description"].(string)
		config.BucketName = cfg["bucket_name"].(string)
		config.Region = cfg["region"].(string)
		encrypted := cfg["encrypted"].(bool)
		enabled := cfg["enabled"].(bool)
		config.Encrypted = &encrypted
		config.Enabled = &enabled
		if authConfigs := cfg["auth_config"].([]interface{}); len(authConfigs) > 0 {
			auth := authConfigs[0].(map[string]interface{})
			authConfig := &S3ArchivingAuthConfig{
				AuthenticationMode: auth["authentication_mode"].(string),
			}
			switch authConfig.AuthenticationMode {
			case "AccessKey":
				authConfig.AccessKeyId = auth["access_key_id"].(string)
				authConfig.AccessKeySecret = auth["access_key_secret"].(string)
			case "RoleBased":
				authConfig.RoleArn = auth["role_arn"].(string)
			}
			config.AuthConfig = authConfig
		}
	case "Syslog":
		config.Protocol = cfg["protocol"].(string)
		config.Host = cfg["host"].(string)
		config.Port = cfg["port"].(int)
		config.Token = cfg["token"].(string)
	case "Hitachi":
		config.URL = cfg["url"].(string)
		config.ObjectId = cfg["object_id"].(string)
		config.Username = cfg["username"].(string)
		config.Password = cfg["password"].(string)
	case "RestAPI":
		config.URL = cfg["url"].(string)
		config.ObjectId = cfg["object_id"].(string)
		config.Username = cfg["username"].(string)
		config.Password = cfg["password"].(string)
	}

	return DataArchivingDestination{
		ID:                d.Id(),
		DestinationName:   d.Get("destination_name").(string),
		DestinationConfig: config,
	}
}

func flattenDataArchivingDestinationConfig(cfg DataArchivingDestinationConfig, d *schema.ResourceData) []map[string]interface{} {
	result := map[string]interface{}{
		"destination_type":      cfg.DestinationType,
		"description":           cfg.Description,
		"bucket_name":           cfg.BucketName,
		"region":                cfg.Region,
		"invalidated_by_system": false,
		"protocol":              cfg.Protocol,
		"host":                  cfg.Host,
		"port":                  cfg.Port,
		"url":                   cfg.URL,
		"object_id":             cfg.ObjectId,
		"username":              cfg.Username,
	}

	if cfg.Encrypted != nil {
		result["encrypted"] = *cfg.Encrypted
	}
	if cfg.Enabled != nil {
		result["enabled"] = *cfg.Enabled
	}
	if cfg.InvalidatedBySystem != nil {
		result["invalidated_by_system"] = *cfg.InvalidatedBySystem
	}

	if cfg.AuthConfig != nil {
		result["auth_config"] = flattenS3ArchivingAuthConfig(cfg.AuthConfig, d)
	} else {
		result["auth_config"] = []map[string]interface{}{}
	}

	// Preserve sensitive fields from state since API returns masked values
	if existing := getExistingDestinationConfig(d); existing != nil {
		result["token"] = existing["token"]
		result["password"] = existing["password"]
	} else {
		result["token"] = cfg.Token
		result["password"] = cfg.Password
	}

	return []map[string]interface{}{result}
}

func flattenS3ArchivingAuthConfig(auth *S3ArchivingAuthConfig, d *schema.ResourceData) []map[string]interface{} {
	result := map[string]interface{}{
		"authentication_mode": auth.AuthenticationMode,
		"access_key_id":       auth.AccessKeyId,
		"role_arn":            auth.RoleArn,
	}

	// Preserve access_key_secret from state since API returns masked value
	if existing := getExistingAuthConfig(d); existing != nil {
		result["access_key_secret"] = existing["access_key_secret"]
	} else {
		result["access_key_secret"] = auth.AccessKeySecret
	}

	return []map[string]interface{}{result}
}

func getExistingDestinationConfig(d *schema.ResourceData) map[string]interface{} {
	if configs, ok := d.GetOk("destination_config"); ok {
		cfgList := configs.([]interface{})
		if len(cfgList) > 0 && cfgList[0] != nil {
			return cfgList[0].(map[string]interface{})
		}
	}
	return nil
}

func getExistingAuthConfig(d *schema.ResourceData) map[string]interface{} {
	if configs, ok := d.GetOk("destination_config"); ok {
		cfgList := configs.([]interface{})
		if len(cfgList) > 0 && cfgList[0] != nil {
			cfg := cfgList[0].(map[string]interface{})
			if authConfigs, ok := cfg["auth_config"].([]interface{}); ok && len(authConfigs) > 0 && authConfigs[0] != nil {
				return authConfigs[0].(map[string]interface{})
			}
		}
	}
	return nil
}
