package main

import (
	"flag"
	conf "github.office.opendns.com/quadra/linux-job/config"
	"github.office.opendns.com/quadra/linux-job/pkg/grpc/server"
	"log"
)


func main() {
	// TODO file names handle
	config := conf.NewConfig()
	flag.StringVar(&config.ServerAddress, "host", "localhost:5050", "host:port")
	flag.StringVar(&config.ClientCA, "ca", "certs/ca.crt", "client ca path")
	flag.StringVar(&config.ServerCertificate, "cert", "certs/server.crt", "server cert path")
	flag.StringVar(&config.ServerKey, "key", "certs/server.key", "server key path")
	flag.Parse()
	if err := server.StartServer(config); err != nil {
		log.Fatalf("fail to start  grpc server, %v", err)
	}
}
