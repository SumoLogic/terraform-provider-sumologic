package sumologic

import (
	"fmt"
	"log"
	"regexp"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func resourceSumologicServiceAccount() *schema.Resource {
	return &schema.Resource{
		Create: resourceSumologicServiceAccountCreate,
		Read:   resourceSumologicServiceAccountRead,
		Update: resourceSumologicServiceAccountUpdate,
		Delete: resourceSumologicServiceAccountDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringLenBetween(1, 128),
			},
			"email": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.All(validation.StringLenBetween(1, 255), validation.StringMatch(regexp.MustCompile("@"), "must be an email address")),
			},
			"role_ids": {
				Type:     schema.TypeSet,
				Required: true,
				MinItems: 1,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"is_active": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  true,
			},
			"transfer_to": {
				Type:          schema.TypeString,
				Optional:      true,
				ConflictsWith: []string{"delete_content"},
			},
			"delete_content": {
				Type:          schema.TypeBool,
				Optional:      true,
				Default:       false,
				ConflictsWith: []string{"transfer_to"},
			},
		},
	}
}

func resourceSumologicServiceAccountRead(d *schema.ResourceData, meta interface{}) error {
	c := meta.(*Client)

	id := d.Id()
	account, err := c.GetServiceAccount(id)
	if err != nil {
		return err
	}
	if account == nil {
		log.Printf("[WARN] Service account not found, removing from state: %v", id)
		d.SetId("")
		return nil
	}

	d.Set("name", account.Name)
	d.Set("email", account.Email)
	if err := d.Set("role_ids", account.RoleIds); err != nil {
		return fmt.Errorf("error setting role ids for resource %s: %s", d.Id(), err)
	}
	d.Set("is_active", account.IsActive)
	return nil
}

func resourceSumologicServiceAccountCreate(d *schema.ResourceData, meta interface{}) error {
	c := meta.(*Client)

	id, err := c.CreateServiceAccount(ServiceAccount{
		Name:    d.Get("name").(string),
		Email:   d.Get("email").(string),
		RoleIds: serviceAccountRoleIds(d),
	})
	if err != nil {
		return err
	}
	d.SetId(id)

	// The API always creates an active service account; deactivating is an update.
	if !d.Get("is_active").(bool) {
		if err := c.UpdateServiceAccount(id, serviceAccountUpdateFromConfig(d)); err != nil {
			return err
		}
	}

	return resourceSumologicServiceAccountRead(d, meta)
}

func resourceSumologicServiceAccountUpdate(d *schema.ResourceData, meta interface{}) error {
	c := meta.(*Client)

	if d.HasChanges("name", "email", "role_ids", "is_active") {
		if err := c.UpdateServiceAccount(d.Id(), serviceAccountUpdateFromConfig(d)); err != nil {
			return err
		}
	}

	return resourceSumologicServiceAccountRead(d, meta)
}

func resourceSumologicServiceAccountDelete(d *schema.ResourceData, meta interface{}) error {
	c := meta.(*Client)

	err := c.DeleteServiceAccount(d.Id(), d.Get("transfer_to").(string), d.Get("delete_content").(bool))
	if isNotFoundError(err) {
		return nil
	}
	return err
}

func serviceAccountRoleIds(d *schema.ResourceData) []string {
	raw := d.Get("role_ids").(*schema.Set).List()
	roleIds := make([]string, len(raw))
	for i, v := range raw {
		roleIds[i] = v.(string)
	}
	return roleIds
}

func serviceAccountUpdateFromConfig(d *schema.ResourceData) ServiceAccountUpdate {
	return ServiceAccountUpdate{
		Name:     d.Get("name").(string),
		Email:    d.Get("email").(string),
		RoleIds:  serviceAccountRoleIds(d),
		IsActive: d.Get("is_active").(bool),
	}
}
