package main

import (
	"fmt"
	"github.com/google/uuid"
	"google.golang.org/grpc"
	"io"
	"log"
	"net"
	"os"
	"sync"

	pb "github.com/samuelbagattin/grpc-chat/proto"
)

const (
	port = ":50051"
)

var connections = map[uuid.UUID]pb.MessagesService_ChatServer{}

func broadcastMessage(message *pb.Message) error {
	var wg sync.WaitGroup
	for _, str := range connections {
		wg.Add(1)
		go func(chatServer pb.MessagesService_ChatServer) {
			errSend := chatServer.Send(message)
			if errSend != nil {
				log.Printf("Error while broadcasting message : %v\n", errSend)
			}
			wg.Done()
		}(str)
	}
	wg.Wait()
	return nil
}

func addNewClient(chatServer pb.MessagesService_ChatServer) (*uuid.UUID, error) {
	var connectionUid = uuid.New()
	err := broadcastMessage(&pb.Message{
		Name:    "Inconnu",
		Content: "A rejoint le chat",
	})
	if err != nil {
		return nil, err
	}
	connections[connectionUid] = chatServer
	return &connectionUid, nil
}

// server is used to implement helloworld.GreeterServer.
type server struct {
	pb.UnimplementedMessagesServiceServer
}

func (s *server) Chat(stream pb.MessagesService_ChatServer) error {
	uid, err := addNewClient(stream)
	if err != nil {
		return err
	}
	for {
		in, err := stream.Recv()
		if err == io.EOF {
			return nil
		}
		if err != nil {
			delete(connections, *uid)
			return err
		}

		errBr := broadcastMessage(in)
		if errBr != nil {
			delete(connections, *uid)
			return errBr
		}
	}
}

func main() {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%s", os.Args[1]))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	pb.RegisterMessagesServiceServer(s, &server{})
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
