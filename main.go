package main

import (
	"github.com/ivandeex/packer-builder-vultr/vultr"
	"github.com/hashicorp/packer/packer/plugin"
)

func main() {
	server, err := plugin.Server()
	if err != nil {
		panic(err)
	}
	err = server.RegisterBuilder(new(vultr.Builder))
	if err != nil {
		panic(err)
	}
	server.Serve()
}
