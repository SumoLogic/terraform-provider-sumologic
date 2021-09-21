package sumologic

import (
	"errors"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"log"
)

func resourceSumologicCSENetworkBlock() *schema.Resource {
	return &schema.Resource{
		Create: resourceSumologicCSENetworkBlockCreate,
		Read:   resourceSumologicCSENetworkBlockRead,
		Delete: resourceSumologicCSENetworkBlockDelete,
		Update: resourceSumologicCSENetworkBlockUpdate,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"address_block": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: false,
			},
			"label": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: false,
			},
			"internal": {
				Type:     schema.TypeBool,
				Optional: true,
				ForceNew: false,
				Default:  true,
			},
			"suppresses_signals": {
				Type:     schema.TypeBool,
				Optional: true,
				ForceNew: false,
				Default:  false,
			},
		},
	}
}

func resourceSumologicCSENetworkBlockRead(d *schema.ResourceData, meta interface{}) error {
	c := meta.(*Client)

	var cseNetworkBlock *CSENetworkBlock
	id := d.Id()

	cseNetworkBlock, err := c.GetCSENetworkBlock(id)
	if err != nil {
		log.Printf("[WARN] CSE Network Block not found when looking by id: %s, err: %v", id, err)

	}

	if cseNetworkBlock == nil {
		log.Printf("[WARN] CSE Network Block not found, removing from state: %v - %v", id, err)
		d.SetId("")
		return nil
	}

	d.Set("address_block", cseNetworkBlock.AddressBlock)
	d.Set("label", cseNetworkBlock.Label)
	d.Set("internal", cseNetworkBlock.Internal)
	d.Set("suppresses_signals", cseNetworkBlock.SuppressesSignals)

	return nil
}

func resourceSumologicCSENetworkBlockDelete(d *schema.ResourceData, meta interface{}) error {
	c := meta.(*Client)

	id := d.Id()
	return c.DeleteCSENetworkBlock(id)

}

func resourceSumologicCSENetworkBlockCreate(d *schema.ResourceData, meta interface{}) error {
	c := meta.(*Client)

	if d.Id() == "" {
		id, err := c.CreateCSENetworkBlock(CSENetworkBlock{
			AddressBlock:      d.Get("address_block").(string),
			Label:             d.Get("label").(string),
			Internal:          d.Get("internal").(bool),
			SuppressesSignals: d.Get("suppresses_signals").(bool),
		})

		if err != nil {
			return err
		}
		log.Printf("[INFO] got id: %s", id)
		d.SetId(id)
	}

	return resourceSumologicCSENetworkBlockUpdate(d, meta)
}

func resourceSumologicCSENetworkBlockUpdate(d *schema.ResourceData, meta interface{}) error {
	cseNetworkBlock, err := resourceToCSENetworkBlock(d)
	if err != nil {
		return err
	}

	c := meta.(*Client)
	if err = c.UpdateCSENetworkBlock(cseNetworkBlock); err != nil {
		return err
	}

	return resourceSumologicCSENetworkBlockRead(d, meta)
}

func resourceToCSENetworkBlock(d *schema.ResourceData) (CSENetworkBlock, error) {
	id := d.Id()
	if id == "" {
		return CSENetworkBlock{}, errors.New("network block id not specified")
	}

	return CSENetworkBlock{
		ID:                id,
		AddressBlock:      d.Get("address_block").(string),
		Label:             d.Get("label").(string),
		Internal:          d.Get("internal").(bool),
		SuppressesSignals: d.Get("suppresses_signals").(bool),
	}, nil
}
