package main

import (
	"flag"
	"fmt"
	"os"
	"strconv"
	"strings"

	"code.cloudfoundry.org/cli/plugin"
)

const (
	addCfmrNetworkPolicyCommand = "add-cfmr-network-policy"
)

type AddCfmrNetworkPolicyPlugin struct{}

type CommandArgs struct {
	command        string
	sourceApp      string
	destinationApp string
	port           int
	protocol       string
}

func (c *AddCfmrNetworkPolicyPlugin) Run(cliConnection plugin.CliConnection, args []string) {
	ca := validateAndParseArgs(args)
	fmt.Printf("CommandArgs:%+v\n", ca)
}

func validateAndParseArgs(args []string) CommandArgs {
	ca := CommandArgs{}
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

	flagSet := flag.NewFlagSet(addCfmrNetworkPolicyCommand, flag.ExitOnError)
	fmt.Println("Parsing Command Arguments...")

	destinationApp := flagSet.String(
		"destination-app",
		"",
		"destination application which is to be exposed by service",
	)

	port := flagSet.String(
		"port",
		"",
		"port on which destination app will be exposed",
	)

	protocol := flagSet.String(
		"protocol",
		"",
		"protocol on which destination app will be exposed",
	)

	err := flagSet.Parse(args[2:])
	if err != nil {
		fmt.Println("ERROR:>")
		fmt.Println(err)
	}

	ca.command = strings.TrimSpace(args[0])
	ca.sourceApp = strings.TrimSpace(args[1])
	ca.destinationApp = strings.TrimSpace(*destinationApp)
	ca.port, err = strconv.Atoi(strings.TrimSpace(*port))
	if err != nil {
		fmt.Println("port should be a number")
		os.Exit(0)
	}
	ca.protocol = strings.TrimSpace(*protocol)

	if ca.destinationApp == "" {
		fmt.Println("destination app name is required")
		os.Exit(0)
	}

	if ca.port == 0 {
		fmt.Println("port number is required")
		os.Exit(0)
	}

	if ca.protocol == "" {
		fmt.Println("protocol is required")
		os.Exit(0)
	}

	return ca
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
