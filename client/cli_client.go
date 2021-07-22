package client

import "code.cloudfoundry.org/cli/plugin"

type CliClient struct {
	cf plugin.CliConnection
}

func NewCliClient(cf plugin.CliConnection) *CliClient {
	return &CliClient{
		cf: cf,
	}
}

func (cliClient *CliClient) GetAppGUID(appName string) (string, error) {
	appModel, err := cliClient.cf.GetApp(appName)
	if err != nil {
		return "", err
	}
	return appModel.Guid, nil
}
