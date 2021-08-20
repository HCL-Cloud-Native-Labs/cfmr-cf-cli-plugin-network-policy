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

	"code.cloudfoundry.org/cli/plugin"
	"code.ibm.com/cfmr-cf-cli-plugin-network-policy/client"
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
	ports          []int
	protocols      []string
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
	// fmt.Printf("%#v\n", args)
	ca := parseAndValidateArgs(args)
	cliClient := client.NewCliClient(cliConnection)
	createNetworkPolicy(cliClient, ca)
}

func parseAndValidateArgs(args []string) CommandArgs {
	ca := CommandArgs{ports: []int{}, protocols: []string{}}
	if len(args) == 1 && args[0] == "CLI-MESSAGE-UNINSTALL" {
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
	// fmt.Println("Parsing Command Arguments...")

	destinationApp := flagSet.String(
		"destination-app",
		"",
		"destination application which is to be exposed by service",
	)

	port := flagSet.String(
		"port",
		"8080",
		"ports on which destination app will be exposed",
	)

	protocol := flagSet.String(
		"protocol",
		"tcp",
		"protocols on which destination app will be exposed",
	)

	err := flagSet.Parse(args[2:])
	if err != nil {
		fmt.Println("ERROR:", err)
	}

	if *destinationApp == "" {
		fmt.Println("destination app name is required")
		os.Exit(0)
	}

	ca.command = strings.TrimSpace(args[0])
	ca.sourceApp = strings.TrimSpace(args[1])
	ca.destinationApp = strings.TrimSpace(*destinationApp)
	validateAndPopulatePortAndProtocol(*port, *protocol, &ca)
	return ca
}

func validateAndPopulatePortAndProtocol(port string, protocol string, ca *CommandArgs) {
	ports := strings.Split(port, ",")

	for _, p := range ports {
		prt, err := strconv.Atoi(strings.TrimSpace(p))
		if err != nil {
			fmt.Println("port should be a number")
			os.Exit(0)
		}
		ca.ports = append(ca.ports, prt)
	}

	// validate and populate protocols
	protocols := strings.Split(protocol, ",")
	if len(protocols) > len(ports) {
		fmt.Println("protocol and port mismatched")
		os.Exit(0)
	}

	for _, p := range protocols {
		prot := strings.TrimSpace(p)
		switch {
		case prot == "":
			ca.protocols = append(ca.protocols, "tcp")
		case prot != "tcp" && prot != "udp":
			fmt.Println("Invalid protocol, valid values are (tcp | udp)")
			os.Exit(0)
		default:
			ca.protocols = append(ca.protocols, prot)

		}
	}

	// Set default protocol for all other ports
	for i := 0; i < len(ports)-len(protocols); i++ {
		ca.protocols = append(ca.protocols, "tcp")
	}
}

func createNetworkPolicy(cliClient *client.CliClient, ca CommandArgs) {
	fmt.Println("Fetching GUID for", ca.sourceApp, "app")
	sourceGUID, err := cliClient.GetAppGUID(ca.sourceApp)
	if err != nil {
		fmt.Println(err)
		os.Exit(0)
	}

	fmt.Println("GUID for", ca.sourceApp, "app is", sourceGUID)

	fmt.Println("Fetching GUID for", ca.destinationApp, "app")
	destinationGUID, err := cliClient.GetAppGUID(ca.destinationApp)
	if err != nil {
		fmt.Println(err)
		os.Exit(0)
	}

	fmt.Println("GUID for", ca.destinationApp, "app is", destinationGUID)

	serviceArgs := []string{"create-service", networkPolicyServiceBroker, networkPolicyServicePlan}
	serviceName := ca.sourceApp + "-" + ca.destinationApp
	serviceArgs = append(serviceArgs, serviceName)
	serviceArgs = append(serviceArgs, "-c")
	serviceConfigParams := NetworkPolicyServiceConfigParams{
		SourceGUID:         sourceGUID,
		DestinationAppName: ca.destinationApp,
		DestinationGUID:    destinationGUID,
		Ports:              []ServicePort{},
	}

	for i, port := range ca.ports {
		serviceConfigParams.Ports = append(serviceConfigParams.Ports, ServicePort{
			Name:       fmt.Sprintf("port%02d", i+1),
			Port:       port,
			TargetPort: port,
			Protocol:   ca.protocols[i],
		})
	}

	serviceConfigParamsJSON, err := json.Marshal(serviceConfigParams)
	if err != nil {
		fmt.Println("Unable to unmarshal network policy configuration parameters", " \nERROR:", err)
		os.Exit(0)
	}
	serviceArgs = append(serviceArgs, string(serviceConfigParamsJSON))
	// fmt.Println("serviceArgs", serviceArgs)
	_, err = cliClient.CliCommand(serviceArgs...)
	if err != nil {
		fmt.Println("Unable to create network policy", " \nERROR:", err)
		os.Exit(0)
	}
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
