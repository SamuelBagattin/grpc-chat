package main

import (
	"bufio"
	"fmt"
	pb "github.com/samuelbagattin/grpc-chat/proto"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"io"
	"log"
	"os"
)

const (
	address     = "localhost:50051"
	defaultName = "world"
)

var userName string

func main() {
	userName = os.Args[1]
	// Set up a connection to the server.
	conn, err := grpc.Dial(address, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer func(conn *grpc.ClientConn) {
		err := conn.Close()
		if err != nil {
			panic(err)
		}
	}(conn)
	c := pb.NewMessagesServiceClient(conn)

	stream, err := c.Chat(context.Background())
	receiveMsg := make(chan *pb.Message)
	cl := make(chan bool)

	go receiveMessages(stream, receiveMsg, cl)
	go processMessages(receiveMsg, cl)
	go sendMessages(stream)

	<-cl
	if errClose := stream.CloseSend(); errClose != nil {
		panic(errClose)
	}

}

func receiveMessages(stream pb.MessagesService_ChatClient, receive chan *pb.Message, cl chan bool) {
	for {
		in, err := stream.Recv()
		if err == io.EOF {
			// read done.
			cl <- true
			return
		}
		if err != nil {
			log.Fatalf("Failed to receive a note : %v", err)
		}
		receive <- in
	}
}

func processMessages(receive chan *pb.Message, cl chan bool) {
	for {
		select {
		case msg := <-receive:
			fmt.Printf("%s: %s\n", msg.GetName(), msg.GetContent())
		case <-cl:
			fmt.Printf("closing connection")
		}
	}
}

func sendMessages(stream pb.MessagesService_ChatClient) {
	for {
		reader := bufio.NewReader(os.Stdin)
		readString, err := reader.ReadString('\n')
		if err != nil {
			log.Printf("%v", err)
		}
		print("\r                \r")
		print("\r            \r")
		print("\r           \r")
		print("\r         \r")
		print("\r\r")
		if err := stream.Send(&pb.Message{
			Name:    userName,
			Content: readString,
		}); err != nil {
			log.Fatalf("Failed to send a note: %v", err)
		}
	}

}
