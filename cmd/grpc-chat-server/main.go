package main

import (
	"fmt"
	"github.com/google/uuid"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"io"
	"log"
	"net"
	"os"
	"strings"
	"sync"

	pb "github.com/SamuelBagattin/grpc-chat/proto"
)

type ClientEntry struct {
	HasRegistered       bool
	UserName            string
	ClientUid           string
	SendMessagesChannel chan SendingClient
}

type Message struct {
	UserName string
	Content  string
}

type SendingClient struct {
	*Message
	ClientEntry
}

var connections = make(map[string]ClientEntry)

func broadcast(action func(entry ClientEntry) error) error {
	errorsChan := make(chan error, 1)
	var wg sync.WaitGroup
	for _, str := range connections {
		wg.Add(1)
		go func(clientEntry ClientEntry) {
			errorsChan <- action(clientEntry)
			wg.Done()
		}(str)
	}
	wg.Wait()
	//if err := <-errorsChan; err != nil {
	//	return err
	//}
	return nil
}

func addNewClient(newcomersChan chan ClientEntry, userName string) (*string, error) {
	var connectionUid = uuid.New().String()
	newcomer := ClientEntry{
		HasRegistered:       true,
		ClientUid:           connectionUid,
		UserName:            userName,
		SendMessagesChannel: make(chan SendingClient),
	}
	newcomersChan <- newcomer
	connections[connectionUid] = newcomer
	return &connectionUid, nil
}

// server is used to implement helloworld.GreeterServer.
type server struct {
	pb.UnimplementedMessagesServiceServer
	NewcomersChan   chan ClientEntry
	NewMessagesChan chan SendingClient
}

func GetRegisteredClient(clientUid string) (*ClientEntry, error) {
	sender := connections[clientUid]
	if sender.HasRegistered == false {
		return nil, status.Error(codes.Unauthenticated, fmt.Sprintf("client not registered : %s", clientUid))
	}
	return &sender, nil
}

func (s *server) SubscribeToNewcomers(request *pb.SubscribeToNewcomersRequest, stream pb.MessagesService_SubscribeToNewcomersServer) error {
	requester, clientError := GetRegisteredClient(request.GetClientUid())
	if clientError != nil {
		return clientError
	}

	for {
		newcomer := <-s.NewcomersChan
		if newcomer == *requester {
			continue
		}
		err := stream.Send(&pb.Newcomer{Name: newcomer.UserName})
		if err != nil {
			return status.Errorf(codes.Internal, "Cannot send newcomer %s to %s", newcomer, connections[request.GetClientUid()].UserName)
		}
	}
}

func (s *server) SendMessages(stream pb.MessagesService_SendMessagesServer) error {
	for {
		in, err := stream.Recv()
		sender, clientError := GetRegisteredClient(in.GetClientUid())
		if clientError != nil {
			return clientError
		}
		if err == io.EOF {
			return stream.SendAndClose(&pb.Ok{Ok: true})
		}
		if err != nil {
			return err
		}

		sendingClient := SendingClient{
			Message: &Message{
				UserName: sender.UserName,
				Content:  strings.Trim(in.GetContent(), "\n"),
			},
			ClientEntry: *sender,
		}
		for _, entry := range connections {
			go func(channel chan SendingClient) {
				channel <- sendingClient
			}(entry.SendMessagesChannel)
		}
	}

}

func (s *server) ReceiveMessages(infos *pb.ClientInfos, stream pb.MessagesService_ReceiveMessagesServer) error {
	sender, clientError := GetRegisteredClient(infos.GetClientUid())
	if clientError != nil {
		return clientError
	}
	for {
		newMessage := <-sender.SendMessagesChannel
		if newMessage.ClientUid == sender.ClientUid {
			continue
		}
		err := stream.Send(&pb.ReceiveMessageResponse{
			Name:    newMessage.Message.UserName,
			Content: newMessage.Content,
		})
		if err != nil {
			return err
		}
	}
}

func (s *server) RegisterClient(_ context.Context, info *pb.RegisterClientRequest) (*pb.ClientInfos, error) {
	uid, err := addNewClient(s.NewcomersChan, info.GetUserName())
	if err != nil {
		log.Printf("%v", err)
	}
	return &pb.ClientInfos{ClientUid: *uid}, nil
}

func main() {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%s", os.Args[1]))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	pb.RegisterMessagesServiceServer(s, &server{
		NewcomersChan:   make(chan ClientEntry, 1),
		NewMessagesChan: make(chan SendingClient),
	})
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
