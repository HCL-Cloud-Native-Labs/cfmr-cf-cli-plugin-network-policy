package client

import "code.cloudfoundry.org/cli/plugin"

type CliClient struct {
	cliConnection plugin.CliConnection
}

func NewCliClient(cliConn plugin.CliConnection) *CliClient {
	return &CliClient{
		cliConnection: cliConn,
	}
}

func (cliClient *CliClient) GetAppGUID(appName string) (string, error) {
	appModel, err := cliClient.cliConnection.GetApp(appName)
	if err != nil {
		return "", err
	}
	return appModel.Guid, nil
}
