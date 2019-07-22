package main

import (
	"log"

	"github.com/vultr/packer-builder-vultr/vultr"
	"github.com/hashicorp/packer/packer/plugin"
)

var version = "DEV"

func main() {
	log.Println("[INFO] Vultr builder version:", version)

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
