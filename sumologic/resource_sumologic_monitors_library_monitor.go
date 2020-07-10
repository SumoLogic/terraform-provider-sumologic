// ----------------------------------------------------------------------------
//
//     ***     AUTO GENERATED CODE    ***    AUTO GENERATED CODE     ***
//
// ----------------------------------------------------------------------------
//
//     This file is automatically generated by Sumo Logic and manual
//     changes will be clobbered when the file is regenerated. Do not submit
//     changes to this file.
//
// ----------------------------------------------------------------------------
package sumologic

import (
  "log"
  "github.com/hashicorp/terraform-plugin-sdk/helper/schema"
  
)

func resourceSumologicMonitorsLibraryMonitor() *schema.Resource {
    return &schema.Resource{
      Create: resourceSumologicMonitorsLibraryMonitorCreate,
      Read: resourceSumologicMonitorsLibraryMonitorRead,
      Update: resourceSumologicMonitorsLibraryMonitorUpdate,
      Delete: resourceSumologicMonitorsLibraryMonitorDelete,
      Importer: &schema.ResourceImporter{
        State: schema.ImportStatePassthrough,
      },

       Schema: map[string]*schema.Schema{
        ,
        "post_request_map": {
           Type: schema.TypeMap,
          Optional: true,
           Elem: &schema.Schema{
            Type: schema.TypeString,
            },
         },
    },
  }
}




func resourceToMonitorsLibraryMonitor(d *schema.ResourceData) MonitorsLibraryMonitor {
   rawNotifications := d.Get("notifications").([]interface{})
	notifications := make([]string, len(rawNotifications ))
	for i, v := range rawNotifications  {
		notifications[i] = v.(string)
	}
rawTriggers := d.Get("triggers").([]interface{})
	triggers := make([]string, len(rawTriggers ))
	for i, v := range rawTriggers  {
		triggers[i] = v.(string)
	}
rawQueries := d.Get("queries").([]interface{})
	queries := make([]string, len(rawQueries ))
	for i, v := range rawQueries  {
		queries[i] = v.(string)
	}
   
   return MonitorsLibraryMonitor{
    Name: d.Get("name").(string),
    ID: d.Id(),
    MonitorType: d.Get("monitor_type").(string),
    Description: d.Get("description").(string),
    Queries: queries,
    Notifications: notifications,
    Type: d.Get("type").(string),
    Triggers: triggers,
   }
 }