package server

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	conf "github.office.opendns.com/quadra/linux-job/config"
	pb "github.office.opendns.com/quadra/linux-job/pkg/grpc/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"io/ioutil"
	"net"
)

func loadTLSCredentials(conf conf.Config) (credentials.TransportCredentials, error) {
	ClientCA, err := ioutil.ReadFile(conf.ClientCA)
	if err != nil {
		return nil, err
	}
	// Certification pool to append client CA's crt
	certPool := x509.NewCertPool()
	if !certPool.AppendCertsFromPEM(ClientCA) {
		return nil, fmt.Errorf("failed to add client CA's certificate")
	}
	// Load server's certificate and private key
	serverCert, err := tls.LoadX509KeyPair(conf.ServerCertificate, conf.ServerKey)
	if err != nil {
		return nil, err
	}
	// cert client verification
	config := &tls.Config{
		Certificates: []tls.Certificate{serverCert},
		ClientAuth:   tls.RequireAndVerifyClientCert,
		ClientCAs:    certPool,
		MinVersion:   tls.VersionTLS13,
	}
	return credentials.NewTLS(config), nil
}


func createServer(conf conf.Config, cred credentials.TransportCredentials) (*grpc.Server, net.Listener, error) {
	l, err := net.Listen("tcp", conf.ServerAddress)
	if err != nil {
		return nil, nil, err
	}
	w := NewWorkerServer()

	grpcServer := grpc.NewServer(
		grpc.Creds(cred),
		grpc.UnaryInterceptor(w.UnaryAuthInterceptor),
		grpc.StreamInterceptor(w.StreamAuthInterceptor),
	)
	pb.RegisterWorkerServiceServer(grpcServer, NewWorkerServer())
	return grpcServer, l, nil
}

func StartServer(conf conf.Config) error {
	cred, err := loadTLSCredentials(conf)
	if err != nil {
		return err
	}
	serv, lis, err := createServer(conf, cred)
	if err != nil {
		return err
	}
	defer lis.Close()
	if err := serv.Serve(lis); err != nil {
		return err
	}
	return nil
}