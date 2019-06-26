package vultr

import (
	"net"

	"github.com/JamesClonk/vultr/lib"
	"github.com/hashicorp/packer/helper/multistep"
	"golang.org/x/crypto/ssh"
)

func commHost(state multistep.StateBag) (string, error) {
	return state.Get("server").(lib.Server).MainIP, nil
}

func keyboardInteractive(password string) ssh.KeyboardInteractiveChallenge {
	return func(user, instruction string, questions []string, echos []bool) ([]string, error) {
		answers := make([]string, len(questions))
		for i := range questions {
			answers[i] = password
		}
		return answers, nil
	}
}

func sshConfig(state multistep.StateBag) (*ssh.ClientConfig, error) {
	c := state.Get("config").(Config)
	server := state.Get("server").(lib.Server)

	config := &ssh.ClientConfig{
		User: c.SSHUsername,
		Auth: make([]ssh.AuthMethod, 0),
		HostKeyCallback: func(hostname string, remote net.Addr, key ssh.PublicKey) error {
			return nil // accept anything
		},
	}

	if c.OSID == SnapshotOSID || c.OSID == CustomOSID {
		config.Auth = append(config.Auth, ssh.Password(c.SSHPassword), keyboardInteractive(c.SSHPassword))
	} else {
		config.Auth = append(config.Auth, ssh.Password(server.DefaultPassword), keyboardInteractive(server.DefaultPassword))
	}

	return config, nil
}
