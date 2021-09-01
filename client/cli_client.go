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

import (
	"fmt"

	"code.cloudfoundry.org/cli/plugin"
)

type AppNotFoundError struct {
	Name string
	Err  error
}

func (e *AppNotFoundError) Error() string {
	return fmt.Sprintf("The application %s cannot be found.This could be a result of the application not being visible within the current CF CLI target context or the current user lacking permission to list and view the application names provided or a typo within the provided application names. Please review and try again", e.Name)
}

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
		if err.Error() == fmt.Sprintf("App %s not found", appName) {
			return "", &AppNotFoundError{
				Name: appName,
				Err:  err,
			}
		}
		return "", fmt.Errorf("%s: %w", fmt.Sprintf("Unable to fetch guid for app %s", appName), err)
	}

	return appModel.Guid, nil
}
