package main

import (
	conf "github.office.opendns.com/quadra/linux-job/config"
	"github.office.opendns.com/quadra/linux-job/pkg/grpc/client/command"
	"os"
)

func main() {
	config := conf.NewConfig()
	config.ServerAddress = "localhost:5050"
	config.ServerCA = "certs/ca.crt"
	config.ClientCertificate = "certs/client_a.crt"
	config.ClientKey = "certs/client_a.key"
	err := command.Execute(config, os.Args[1:])
	if err != nil {
		os.Stdout.WriteString(err.Error())
		os.Exit(-1)
	}
	os.Exit(0)
}

