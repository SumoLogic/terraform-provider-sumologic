package main

import (
	"github.com/SumoLogic/terraform-provider-sumologic/sumologic"
	"github.com/hashicorp/terraform-plugin-sdk/plugin"
)

var version string // provider version is passed as compile time argument
var defaultVersion = "dev"

func main() {
	if version == "" {
		sumologic.ProviderVersion = defaultVersion
	} else {
		sumologic.ProviderVersion = version
	}
	plugin.Serve(&plugin.ServeOpts{
		ProviderFunc: sumologic.Provider,
	})
}
