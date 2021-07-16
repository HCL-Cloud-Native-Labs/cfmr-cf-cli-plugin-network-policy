package main

import (
	"fmt"
	"os"

	"code.cloudfoundry.org/cli/plugin"
)

const (
	addCfmrNetworkPolicyCommand = "add-cfmr-network-policy"
)

type AddCfmrNetworkPolicyPlugin struct{}

func (c *AddCfmrNetworkPolicyPlugin) Run(cliConnection plugin.CliConnection, args []string) {
	if len(args) == 1 && args[0] == "CLI-MESSAGE-UNINSTALL" {
		// someone's uninstalling the plugin, but we don't need to clean up
		fmt.Println("Uninstalling plugin, but we don't need to clean up")
		os.Exit(0)
	}

	if len(args) < 1 {
		fmt.Printf("Expected at least 1 argument, but got %d.", len(args))
		os.Exit(0)
	}

	if len(args) == 1 && args[0] == addCfmrNetworkPolicyCommand {
		fmt.Println("Source app name is required")
		os.Exit(0)
	}

	fmt.Println("Hello World")

}

func (c *AddCfmrNetworkPolicyPlugin) GetMetadata() plugin.PluginMetadata {
	return plugin.PluginMetadata{
		Name: "add-cfmr-network-policy-plugin",
		Version: plugin.VersionType{
			Major: 1,
			Minor: 0,
			Build: 0,
		},
		MinCliVersion: plugin.VersionType{
			Major: 6,
			Minor: 7,
			Build: 0,
		},
		Commands: []plugin.Command{
			{
				Name:     "add-cfmr-network-policy",
				HelpText: "Create policy to allow direct network traffic from one app to another",
				UsageDetails: plugin.Usage{
					Usage: "cf add-network-policy SOURCE_APP --destination-app DESTINATION_APP --port PORT --protocol (tcp | udp) ",
					Options: map[string]string{
						"-destination-app": "Destination app name",
						"-port":            "Port number",
						"-protocol":        "Protocol (tcp | udp)",
					},
				},
			},
		},
	}
}

func main() {
	plugin.Start(new(AddCfmrNetworkPolicyPlugin))
}
