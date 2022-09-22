package app

import (
	"github.com/cockroachdb/errors"
	"github.com/tensorchord/envd/pkg/home"
	"github.com/tensorchord/envd/pkg/ssh"
	sshconfig "github.com/tensorchord/envd/pkg/ssh/config"
	"github.com/urfave/cli/v2"
)

var CommandSSH = &cli.Command{
	Name:     "ssh",
	Category: CategoryBasic,
	Hidden:   true,
	Usage:    "TestK8s",
	Flags: []cli.Flag{
		&cli.PathFlag{
			Name:    "private-key",
			Usage:   "Path to the private key",
			Aliases: []string{"k"},
			Value:   sshconfig.GetPrivateKeyOrPanic(),
			Hidden:  true,
		},
		&cli.PathFlag{
			Name:    "public-key",
			Usage:   "Path to the public key",
			Aliases: []string{"pubk"},
			Value:   sshconfig.GetPublicKeyOrPanic(),
			Hidden:  true,
		},
	},
	Action: sshc,
}

func sshc(clicontext *cli.Context) error {
	ac, err := home.GetManager().AuthGetCurrent()
	if err != nil {
		return err
	}
	it := ac.IdentityToken

	opt := ssh.DefaultOptions()
	opt.User = it
	opt.PrivateKeyPath = clicontext.Path("private-key")
	opt.Port = 2222
	sshClient, err := ssh.NewClient(opt)
	if err != nil {
		return errors.Wrap(err, "failed to create the ssh client")
	}
	if err := sshClient.Attach(); err != nil {
		return errors.Wrap(err, "failed to attach to the container")
	}
	return nil
}
