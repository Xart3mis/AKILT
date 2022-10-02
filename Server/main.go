package main

import (
	"bufio"
	"log"
	"net"
	"os"

	pb "github.com/Xart3mis/GoHkarComms/client_data_pb"

	"golang.org/x/net/context"
	"google.golang.org/grpc"
)

type server struct {
	pb.ConsumerServer
}

var on_screen_text string = ""

func main() {
	go func() {
		scanner := bufio.NewScanner(os.Stdin)
		for scanner.Scan() {
			on_screen_text = scanner.Text()
		}
		if err := scanner.Err(); err != nil {
			os.Exit(1)
		}
	}()

	s := grpc.NewServer()
	lis, err := net.Listen("tcp", ":8000")
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}
	pb.RegisterConsumerServer(s, &server{})

	err = s.Serve(lis)
	if err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}

}

func (s *server) UpdateClients(ctx context.Context, in *pb.ClientDataRequest) (*pb.ClientDataResponse, error) {
	return &pb.ClientDataResponse{ClientData: map[string]*pb.ClientData{
		in.ClientId: {
			OnScreenText: on_screen_text,
			ShouldUpdate: len(on_screen_text) > 0,
		},
	},
	}, nil
}
