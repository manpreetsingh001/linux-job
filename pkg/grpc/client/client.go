package client

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	conf "github.office.opendns.com/quadra/linux-job/config"
	pb "github.office.opendns.com/quadra/linux-job/pkg/grpc/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"io/ioutil"
)

func loadTLSCredentials(config conf.Config) (credentials.TransportCredentials, error) {
	// Load certificate of the CA who signed server's certificate
	ServerCA, err := ioutil.ReadFile(config.ServerCA)
	if err != nil {
		return nil, err
	}
	certPool := x509.NewCertPool()
	if !certPool.AppendCertsFromPEM(ServerCA) {
		return nil, fmt.Errorf("failed to add server CA's certificate %v", ServerCA)
	}
	// Load client's certificate and private key
	clientCert, err := tls.LoadX509KeyPair(config.ClientCertificate, config.ClientKey)
	if err != nil {
		return nil, err
	}
	// Create the credentials and return it
	tlsConfig := &tls.Config{
		Certificates: []tls.Certificate{clientCert},
		RootCAs:      certPool,
		MinVersion:   tls.VersionTLS13,
	}
	return credentials.NewTLS(tlsConfig), nil
}

func NewWorkerClient(config conf.Config) (pb.WorkerServiceClient, error) {
	tlsCredentials, err := loadTLSCredentials(config)
	if err != nil {
		return nil, err
	}
	conn, err := grpc.Dial(
		config.ServerAddress,
		grpc.WithTransportCredentials(tlsCredentials),
	)
	if err != nil {
		return nil, err
	}
	return pb.NewWorkerServiceClient(conn), nil
}

