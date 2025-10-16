package vultr

import (
	"context"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"io/fs"
	"os"
	"runtime"

	"github.com/hashicorp/packer-plugin-sdk/multistep"
	"github.com/hashicorp/packer-plugin-sdk/packer"
	"github.com/hashicorp/packer-plugin-sdk/uuid"
	"github.com/vultr/govultr/v3"
	"golang.org/x/crypto/ssh"
)

var (
	rsaBits        int         = 2048
	fileMode       int         = 0600
	sshKeyFileMode fs.FileMode = os.FileMode(fileMode)
)

type stepCreateSSHKey struct {
	Debug        bool
	DebugKeyPath string

	client *govultr.Client

	SSHKeyID string
}

// Run provides the step create SSH key run functionality
func (s *stepCreateSSHKey) Run(_ context.Context, state multistep.StateBag) multistep.StepAction {
	ui := state.Get("ui").(packer.Ui)
	config := state.Get("config").(*Config)

	if !config.createTempSSHPair {
		return multistep.ActionContinue
	}

	ui.Say("Creating temporary SSH key...")

	priv, err := rsa.GenerateKey(rand.Reader, rsaBits)
	if err != nil {
		errOut := fmt.Errorf("error creating temporary SSH key: %s", err)
		state.Put("error", errOut)
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
		errOut := fmt.Errorf("error creating temporary SSH key: %s", err)
		state.Put("error", errOut)
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
	key, _, err := s.client.SSHKey.Create(context.Background(), sshKeyReq)
	if err != nil {
		errOut := fmt.Errorf("error creating temporary SSH key: %s", err)
		state.Put("error", errOut)
		ui.Error(errOut.Error())
		return multistep.ActionHalt
	}

	s.SSHKeyID = key.ID

	state.Put("temp_ssh_key_id", key.ID)

	// If we're in debug mode, output the private key to the working directory.
	if s.Debug {
		ui.Say(fmt.Sprintf("saving key for debug purposes: %s", s.DebugKeyPath))
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
			if err := f.Chmod(sshKeyFileMode); err != nil {
				state.Put("error", fmt.Errorf("error setting permissions of debug key: %s", err))
				return multistep.ActionHalt
			}
		}
	}
	return multistep.ActionContinue
}

// Cleanup provides the step create SSH key cleanup functionality
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
