package main

import (
	"fs/service/config"
	"fs/service/rest"
)

var (
	gitCommitHash, builtAt, builtBy, builtOn string
)

// main loads config, creates the servers and starts them if needed
func main() {
	config.Init(gitCommitHash, builtAt, builtBy, builtOn)
	
	server, err := rest.NewServer()
	if err != nil {
		panic(err)
	}
	server.RunServer()
}
