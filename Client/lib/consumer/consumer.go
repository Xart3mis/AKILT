package consumer

import (
	"context"
	"log"

	pb "github.com/Xart3mis/GoHkarComms/client_data_pb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func Init(address string) pb.ConsumerClient {
	conn, err := grpc.Dial(address, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("Error while making connection, %v", err)
	}

	return pb.NewConsumerClient(conn)
}

func UpdateClients(consumer pb.ConsumerClient, sys_id string) (*pb.ClientDataResponse, error) {
	resp, err := consumer.UpdateClients(context.Background(), &pb.ClientDataRequest{ClientId: sys_id})
	if err != nil {
		return nil, err
	}

	return resp, nil
}
