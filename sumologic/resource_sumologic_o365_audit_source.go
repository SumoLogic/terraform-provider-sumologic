package sumologic

import (
	"fmt"
	"log"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func resourceSumologicO365AuditSource() *schema.Resource {
	o365Source := resourceSumologicSource()
	o365Source.Create = resourceSumologicO365AuditSourceCreate
	o365Source.Read = resourceSumologicO365AuditSourceRead
	o365Source.Update = resourceSumologicO365AuditSourceUpdate
	o365Source.Importer = &schema.ResourceImporter{
		State: resourceSumologicSourceImport,
	}

	o365Source.Schema["message_per_request"] = &schema.Schema{
		Type:     schema.TypeBool,
		Optional: true,
		Default:  false,
	}
	o365Source.Schema["url"] = &schema.Schema{
		Type:     schema.TypeString,
		Computed: true,
	}
	o365Source.Schema["content_type"] = &schema.Schema{
		Type:     schema.TypeString,
		Computed: true,
	}
	o365Source.Schema["third_party_ref"] = &schema.Schema{
		Type:     schema.TypeList,
		Required: true,
		ForceNew: true,
		MaxItems: 1,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"resources": {
					Type:     schema.TypeList,
					Required: true,
					ForceNew: true,
					MinItems: 1,
					MaxItems: 1,
					Elem: &schema.Resource{
						Schema: map[string]*schema.Schema{
							"service_type": {
								Type:         schema.TypeString,
								Required:     true,
								ValidateFunc: validation.StringInSlice([]string{"O365AuditNotification"}, false),
							},
							"path": {
								Type:     schema.TypeList,
								Required: true,
								ForceNew: true,
								MinItems: 1,
								MaxItems: 1,
								Elem: &schema.Resource{
									Schema: map[string]*schema.Schema{
										"type": {
											Type:         schema.TypeString,
											Required:     true,
											ValidateFunc: validation.StringInSlice([]string{"O365NotificationPath"}, false),
										},
										"workload": {
											Type:     schema.TypeString,
											Required: true,
											ValidateFunc: validation.StringInSlice([]string{
												"Audit.Exchange",
												"Audit.AzureActiveDirectory",
												"Audit.SharePoint",
												"Audit.General",
												"DLP.All",
											}, false),
										},
										"region": {
											Type:     schema.TypeString,
											Required: true,
											ValidateFunc: validation.StringInSlice([]string{
												"Commercial",
												"GCC",
												"GCC High",
											}, false),
										},
									},
								},
							},
							"authentication": {
								Type:     schema.TypeList,
								Required: true,
								ForceNew: true,
								MinItems: 1,
								MaxItems: 1,
								Elem: &schema.Resource{
									Schema: map[string]*schema.Schema{
										"type": {
											Type:         schema.TypeString,
											Required:     true,
											ValidateFunc: validation.StringInSlice([]string{"O365AppRegistrationAuthentication"}, false),
										},
										"tenant_id": {
											Type:     schema.TypeString,
											Required: true,
										},
										"client_id": {
											Type:     schema.TypeString,
											Required: true,
										},
										"client_secret": {
											Type:      schema.TypeString,
											Required:  true,
											Sensitive: true,
											DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
												// API returns masked value, suppress diff if new value is set
												return new != "" && old == "********"
											},
										},
									},
								},
							},
						},
					},
				},
			},
		},
	}

	return o365Source
}

func resourceSumologicO365AuditSourceCreate(d *schema.ResourceData, meta interface{}) error {
	c := meta.(*Client)

	if d.Id() == "" {
		source := resourceToO365AuditSource(d)

		id, err := c.CreateHTTPSource(source, d.Get("collector_id").(int))

		if err != nil {
			return err
		}

		d.SetId(strconv.Itoa(id))
	}

	return resourceSumologicO365AuditSourceRead(d, meta)
}

func resourceSumologicO365AuditSourceUpdate(d *schema.ResourceData, meta interface{}) error {
	c := meta.(*Client)

	source := resourceToO365AuditSource(d)

	err := c.UpdateHTTPSource(source, d.Get("collector_id").(int))

	if err != nil {
		return err
	}

	return resourceSumologicO365AuditSourceRead(d, meta)
}

func resourceToO365AuditSource(d *schema.ResourceData) HTTPSource {
	source := resourceToSource(d)
	source.Type = "HTTP"

	httpSource := HTTPSource{
		Source:            source,
		MessagePerRequest: d.Get("message_per_request").(bool),
	}

	thirdPartyRefList := d.Get("third_party_ref").([]interface{})
	if len(thirdPartyRefList) > 0 {
		httpSource.ThirdPartyRef = &HTTPThirdPartyRef{}

		thirdPartyRef := thirdPartyRefList[0].(map[string]interface{})
		resourcesList := thirdPartyRef["resources"].([]interface{})

		for _, res := range resourcesList {
			resourceMap := res.(map[string]interface{})

			pathList := resourceMap["path"].([]interface{})
			pathMap := pathList[0].(map[string]interface{})
			path := HTTPPath{
				Type:     pathMap["type"].(string),
				Workload: pathMap["workload"].(string),
				Region:   pathMap["region"].(string),
			}

			authList := resourceMap["authentication"].([]interface{})
			authMap := authList[0].(map[string]interface{})
			auth := HTTPAuthentication{
				Type:         authMap["type"].(string),
				TenantId:     authMap["tenant_id"].(string),
				ClientId:     authMap["client_id"].(string),
				ClientSecret: authMap["client_secret"].(string),
			}

			resource := HTTPResource{
				ServiceType:    resourceMap["service_type"].(string),
				Path:           path,
				Authentication: auth,
			}

			httpSource.ThirdPartyRef.Resources = append(httpSource.ThirdPartyRef.Resources, resource)
		}
	}

	return httpSource
}

func resourceSumologicO365AuditSourceRead(d *schema.ResourceData, meta interface{}) error {
	c := meta.(*Client)

	id, _ := strconv.Atoi(d.Id())
	source, err := c.GetHTTPSource(d.Get("collector_id").(int), id)

	if err != nil {
		return err
	}

	if source == nil {
		log.Printf("[WARN] O365 Audit source not found, removing from state: %v - %v", id, err)
		d.SetId("")

		return nil
	}

	if err := resourceSumologicSourceRead(d, source.Source); err != nil {
		return fmt.Errorf("%s", err)
	}
	d.Set("message_per_request", source.MessagePerRequest)
	d.Set("url", source.URL)

	if source.ThirdPartyRef != nil && len(source.ThirdPartyRef.Resources) > 0 {
		thirdPartyRef := flattenO365HTTPThirdPartyRef(*source.ThirdPartyRef)
		if err := d.Set("third_party_ref", thirdPartyRef); err != nil {
			return err
		}
	}

	return nil
}

func flattenO365HTTPThirdPartyRef(thirdPartyRef HTTPThirdPartyRef) []map[string]interface{} {
	var resources []map[string]interface{}

	for _, resource := range thirdPartyRef.Resources {
		path := []map[string]interface{}{
			{
				"type":     resource.Path.Type,
				"workload": resource.Path.Workload,
				"region":   resource.Path.Region,
			},
		}

		auth := []map[string]interface{}{
			{
				"type":          resource.Authentication.Type,
				"tenant_id":     resource.Authentication.TenantId,
				"client_id":     resource.Authentication.ClientId,
				"client_secret": resource.Authentication.ClientSecret,
			},
		}

		resourceMap := map[string]interface{}{
			"service_type":   resource.ServiceType,
			"path":           path,
			"authentication": auth,
		}

		resources = append(resources, resourceMap)
	}

	return []map[string]interface{}{
		{
			"resources": resources,
		},
	}
}
