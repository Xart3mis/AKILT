package consumer

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"io/ioutil"

	pb "github.com/Xart3mis/GoHkarComms/client_data_pb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

func Init(address string) (pb.ConsumerClient, error) {
	// certFile, err := filepath.Abs("./../../../cert/cert.pem")
	// if err != nil {
	// 	return nil, err
	// }

	creds, err := loadTLSCredentials()
	if err != nil {
		return nil, err
	}

	conn, err := grpc.Dial(address, grpc.WithTransportCredentials(creds))
	if err != nil {
		return nil, err
	}

	return pb.NewConsumerClient(conn), nil
}

func loadTLSCredentials() (credentials.TransportCredentials, error) {
	// Load certificate of the CA who signed server's certificate
	pemServerCA, err := ioutil.ReadFile("C:\\Users\\ghost\\OneDrive\\Documents\\Code\\GoHkar\\cert\\ca-cert.pem")
	if err != nil {
		return nil, err
	}

	certPool := x509.NewCertPool()
	if !certPool.AppendCertsFromPEM(pemServerCA) {
		return nil, fmt.Errorf("failed to add server CA's certificate")
	}

	// Create the credentials and return it
	config := &tls.Config{
		RootCAs:            certPool,
		InsecureSkipVerify: true,
	}

	return credentials.NewTLS(config), nil
}

func UpdateClients(consumer pb.ConsumerClient, sys_id string) (*pb.ClientDataResponse, error) {
	resp, err := consumer.UpdateClients(context.Background(), &pb.ClientDataRequest{ClientId: sys_id})
	if err != nil {
		return nil, err
	}

	return resp, nil
}
