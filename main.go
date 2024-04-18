package main

import (
	"github.com/SumoLogic/terraform-provider-sumologic/sumologic"
	"github.com/hashicorp/terraform-plugin-sdk/plugin"
	"log"
)

var version string // provider version is passed as compile time argument
var defaultVersion = "dev"

func main() {
	// Remove any date and time prefix in log package function output to
	// prevent duplicate timestamp and incorrect log level setting
	// See: https://developer.hashicorp.com/terraform/plugin/log/writing#duplicate-timestamp-and-incorrect-level-messages
	log.SetFlags(log.Flags() &^ (log.Ldate | log.Ltime))
	if version == "" {
		sumologic.ProviderVersion = defaultVersion
	} else {
		sumologic.ProviderVersion = version
	}
	plugin.Serve(&plugin.ServeOpts{
		ProviderFunc: sumologic.Provider,
	})
}
