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
