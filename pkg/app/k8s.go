package app

import (
	servertypes "github.com/tensorchord/envd-server/api/types"
	"github.com/tensorchord/envd-server/client"
	"github.com/urfave/cli/v2"

	"github.com/tensorchord/envd/pkg/home"
)

var CommandK8s = &cli.Command{
	Name:     "k8s",
	Category: CategoryBasic,
	Hidden:   true,
	Usage:    "TestK8s",
	Action:   k8s,
}

func k8s(clicontext *cli.Context) error {
	ac, err := home.GetManager().AuthGetCurrent()
	if err != nil {
		return err
	}
	it := ac.IdentityToken
	c, err := client.NewClientWithOpts(client.FromEnv)
	if err != nil {
		return err
	}
	req := servertypes.EnvironmentCreateRequest{
		IdentityToken: it,
		Image:         "gaocegege/test-envd",
	}
	_, err = c.EnvironmentCreate(clicontext.Context, req)
	if err != nil {
		return err
	}

	return nil
}
