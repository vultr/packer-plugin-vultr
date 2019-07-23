package main

import (
	"github.com/hashicorp/packer/packer/plugin"
	"github.com/vultr/packer-builder-vultr/vultr"
)

func main() {
	server, err := plugin.Server()
	if err != nil {
		panic(err)
	}

	server.RegisterBuilder(new(vultr.Builder))
	server.Serve()
}
