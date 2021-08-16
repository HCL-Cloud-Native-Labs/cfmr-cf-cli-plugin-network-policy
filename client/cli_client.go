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
package client

import "code.cloudfoundry.org/cli/plugin"

type CliClient struct {
	plugin.CliConnection
}

func NewCliClient(cliConn plugin.CliConnection) *CliClient {
	return &CliClient{
		CliConnection: cliConn,
	}
}

func (cliClient *CliClient) GetAppGUID(appName string) (string, error) {
	appModel, err := cliClient.GetApp(appName)
	if err != nil {
		return "", err
	}
	return appModel.Guid, nil
}
