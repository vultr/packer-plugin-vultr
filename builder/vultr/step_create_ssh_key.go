package vultr

import (
	"context"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"os"
	"runtime"

	"github.com/hashicorp/packer-plugin-sdk/multistep"
	"github.com/hashicorp/packer-plugin-sdk/packer"
	"github.com/hashicorp/packer-plugin-sdk/uuid"
	"github.com/vultr/govultr/v2"
	"golang.org/x/crypto/ssh"
)

type stepCreateSSHKey struct {
	Debug        bool
	DebugKeyPath string

	client *govultr.Client

	SSHKeyID string
}

func (s *stepCreateSSHKey) Run(_ context.Context, state multistep.StateBag) multistep.StepAction {
	ui := state.Get("ui").(packer.Ui)
	config := state.Get("config").(*Config)

	if !config.create_temp_ssh_pair {
		return multistep.ActionContinue
	}

	ui.Say("Creating temporary SSH key...")

	priv, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		err := fmt.Errorf("Error creating temporary SSH key: %s", err)
		state.Put("error", err)
		ui.Error(err.Error())
		return multistep.ActionHalt
	}

	privBlk := pem.Block{
		Type:    "RSA PRIVATE KEY",
		Headers: nil,
		Bytes:   x509.MarshalPKCS1PrivateKey(priv),
	}

	pub, err := ssh.NewPublicKey(&priv.PublicKey)
	if err != nil {
		err := fmt.Errorf("Error creating temporary SSH key: %s", err)
		state.Put("error", err)
		ui.Error(err.Error())
		return multistep.ActionHalt
	}
	config.Comm.SSHPrivateKey = pem.EncodeToMemory(&privBlk)
	config.Comm.SSHPublicKey = ssh.MarshalAuthorizedKey(pub)

	name := fmt.Sprintf("packer-%s", uuid.TimeOrderedUUID())

	sshKeyReq := &govultr.SSHKeyReq{
		Name:   name,
		SSHKey: string(config.Comm.SSHPublicKey),
	}
	key, err := s.client.SSHKey.Create(context.Background(), sshKeyReq)
	if err != nil {
		err := fmt.Errorf("Error creating temporary SSH key: %s", err)
		state.Put("error", err)
		ui.Error(err.Error())
		return multistep.ActionHalt
	}

	s.SSHKeyID = key.ID

	state.Put("temp_ssh_key_id", key.ID)

	// If we're in debug mode, output the private key to the working directory.
	if s.Debug {
		ui.Message(fmt.Sprintf("saving key for debug purposes: %s", s.DebugKeyPath))
		f, err := os.Create(s.DebugKeyPath)
		if err != nil {
			state.Put("error", fmt.Errorf("error saving debug key: %s", err))
			return multistep.ActionHalt
		}

		// Write out the key
		err = pem.Encode(f, &privBlk)
		defer f.Close()
		if err != nil {
			state.Put("error", fmt.Errorf("error saving debug key: %s", err))
			return multistep.ActionHalt
		}

		// Chmod it so that it is SSH ready
		if runtime.GOOS != "windows" {
			if err := f.Chmod(0600); err != nil {
				state.Put("error", fmt.Errorf("error setting permissions of debug key: %s", err))
				return multistep.ActionHalt
			}
		}
	}
	return multistep.ActionContinue
}

func (s *stepCreateSSHKey) Cleanup(state multistep.StateBag) {
	if s.SSHKeyID == "" {
		return
	}

	ui := state.Get("ui").(packer.Ui)
	ui.Say("Deleting temporary SSH key...")

	err := s.client.SSHKey.Delete(context.TODO(), s.SSHKeyID)
	if err != nil {
		ui.Error(fmt.Sprintf("error deleting temporary SSH key (%s) - please delete the key manually: %s", s.SSHKeyID, err))
	}
}
