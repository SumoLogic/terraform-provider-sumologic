package main

import (
	"github.com/SumoLogic/terraform-provider-sumologic/sumologic"
	"github.com/hashicorp/terraform-plugin-sdk/plugin"
)

func main() {
	plugin.Serve(&plugin.ServeOpts{
		ProviderFunc: sumologic.Provider,
	})
}
