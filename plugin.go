/*
###############################################################################
# Licensed Materials - Property of IBM
# Copyright IBM Corporation 2020, 2021. All Rights Reserved
# US Government Users Restricted Rights -
# Use, duplication or disclosure restricted by GSA ADP Schedule Contract with IBM Corp.
#
# This is an internal component, bundled with an official IBM product.
# Please refer to that particular license for additional information.
# ###############################################################################
*/
package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"strconv"
	"strings"

	"code.ibm.com/cfmr-cf-cli-plugin-network-policy/client"

	"code.cloudfoundry.org/cli/plugin"
)

const (
	addCfmrNetworkPolicyCommand = "add-cfmr-network-policy"
	networkPolicyServiceBroker  = "network-policy"
	networkPolicyServicePlan    = "c2c"
)

type AddCfmrNetworkPolicyPlugin struct{}

type CommandArgs struct {
	command        string
	sourceApp      string
	destinationApp string
	port           int
	protocol       string
}

type NetworkPolicyServiceConfigParams struct {
	SourceGUID         string        `json:"source-guid"`
	DestinationAppName string        `json:"destination-appname"`
	DestinationGUID    string        `json:"destination-guid"`
	Ports              []ServicePort `json:"ports"`
}

type ServicePort struct {
	Name       string `json:"name"`
	Port       int    `json:"port"`
	TargetPort int    `json:"targetport"`
	Protocol   string `json:"protocol"`
}

func (c *AddCfmrNetworkPolicyPlugin) Run(cliConnection plugin.CliConnection, args []string) {
	fmt.Println("args:", args)
	ca := parseAndValidateArgs(args)
	fmt.Printf("CommandArgs:%+v\n", ca)
	cliClient := client.NewCliClient(cliConnection)
	createNetworkPolicy(cliClient, ca)
}

func parseAndValidateArgs(args []string) CommandArgs {
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

func createNetworkPolicy(cliClient *client.CliClient, ca CommandArgs) {
	fmt.Println("Fetching GUID for", ca.sourceApp)
	sourceGUID, err := cliClient.GetAppGUID(ca.sourceApp)
	if err != nil {
		fmt.Println("Unable to fetch guid for app", ca.sourceApp, " \nERROR:", err)
		os.Exit(0)
	}

	fmt.Println("sourceAppGUID:", sourceGUID)

	fmt.Println("Fetching GUID for", ca.destinationApp)
	destinationGUID, err := cliClient.GetAppGUID(ca.destinationApp)
	if err != nil {
		fmt.Println("Unable to fetch guid for app", ca.destinationApp, " \nERROR:", err)
		os.Exit(0)
	}

	fmt.Println("destinationGUID:", destinationGUID)

	serviceArgs := []string{"create-service", networkPolicyServiceBroker, networkPolicyServicePlan}
	serviceName := ca.sourceApp + "2" + ca.destinationApp
	serviceArgs = append(serviceArgs, serviceName)
	serviceArgs = append(serviceArgs, "-c")
	serviceConfigParams := NetworkPolicyServiceConfigParams{
		SourceGUID:         sourceGUID,
		DestinationAppName: ca.destinationApp,
		DestinationGUID:    destinationGUID,
		Ports: []ServicePort{
			{
				Name:       "port01",
				Port:       ca.port,
				TargetPort: ca.port,
				Protocol:   ca.protocol,
			},
		},
	}

	serviceConfigParamsJSON, err := json.Marshal(serviceConfigParams)
	if err != nil {
		fmt.Println("Unable to unmarshal network policy service configuration parameters", " \nERROR:", err)
		os.Exit(0)
	}
	serviceArgs = append(serviceArgs, string(serviceConfigParamsJSON))
	fmt.Println("serviceArgs", serviceArgs)
	_, err = cliClient.CliCommand(serviceArgs...)
	if err != nil {
		fmt.Println("Unable to create network policy service", " \nERROR:", err)
		os.Exit(0)
	}
	fmt.Printf("Network policy service '%s' created successfully!!\n", serviceName)
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
