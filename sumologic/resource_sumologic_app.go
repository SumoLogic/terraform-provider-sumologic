package sumologic

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceSumologicApp() *schema.Resource {
	return &schema.Resource{
		Create: resourceSumologicAppCreate,
		Read:   resourceSumologicAppRead,
		Delete: resourceSumologicAppDelete,
		Update: resourceSumologicAppUpdate,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"uuid": {
				Type:     schema.TypeString,
				Required: true,
			},
			"version": {
				Type:     schema.TypeString,
				Required: true,
			},
			"parameters": {
				Type:     schema.TypeMap,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
		},
	}
}

func resourceSumologicAppCreate(d *schema.ResourceData, meta interface{}) error {
	c := meta.(*Client)
	if d.Id() == "" {
		uuid := d.Get("uuid").(string)
		version := d.Get("version").(string)
		parameters := d.Get("parameters").(map[string]interface{})

		appInstallPayload := AppInstallPayload{
			VERSION:    version,
			PARAMETERS: parameters,
		}

		log.Println("=====================================================================")
		log.Printf("Installing app; uuid: %+v, version: %+v\n", uuid, version)
		log.Println("=====================================================================")

		appInstanceId, err := c.CreateAppInstance(uuid, appInstallPayload)
		if err != nil {
			return err
		}
		d.SetId(appInstanceId)
	}

	return resourceSumologicAppRead(d, meta)
}

func resourceSumologicAppRead(d *schema.ResourceData, meta interface{}) error {
	c := meta.(*Client)

	id := d.Id()
	appInstance, err := c.GetAppInstance(id)
	log.Println("=====================================================================")
	log.Printf("Read app instance: %+v\n", appInstance)
	log.Println("=====================================================================")
	if err != nil {
		return err
	}

	if appInstance == nil {
		log.Printf("[WARN] AppInstance not found, removing from state: %v - %v", id, err)
		d.SetId("")
		return nil
	}

	var parameters map[string]interface{}
	if err := json.Unmarshal([]byte(appInstance.CONFIGURATIONBLOB), &parameters); err != nil {
		return err
	}
	d.Set("uuid", appInstance.UUID)
	d.Set("version", appInstance.VERSION)
	d.Set("parameters", parameters)
	d.SetId(appInstance.ID)

	return nil
}

func resourceSumologicAppDelete(d *schema.ResourceData, meta interface{}) error {
	c := meta.(*Client)
	uuid := d.Get("uuid").(string)
	log.Printf("Uninstalling app: %+v\n", uuid)
	return c.DeleteAppInstance(uuid)
}

func resourceSumologicAppUpdate(d *schema.ResourceData, meta interface{}) error {
	c := meta.(*Client)

	uuid := d.Get("uuid").(string)

	// ensure that uuid matches with already installed instance's uuid
	appInstance, err := c.GetAppInstance(d.Id())
	if err != nil {
		return err
	}
	if uuid == appInstance.UUID {
		version := d.Get("version").(string)
		if version == appInstance.VERSION {
			return nil
		}
		parameters := d.Get("parameters").(map[string]interface{})
		appInstallPayload := AppInstallPayload{
			VERSION:    version,
			PARAMETERS: parameters,
		}

		log.Println("=====================================================================")
		log.Printf("Upgrading app; uuid: %+v, version: %+v\n", uuid, version)
		log.Println("=====================================================================")

		_, err := c.UpdateAppInstance(uuid, appInstallPayload)

		if err != nil {
			return err
		}
		return resourceSumologicAppRead(d, meta)
	}

	return fmt.Errorf("uuid is incorrect")
}
