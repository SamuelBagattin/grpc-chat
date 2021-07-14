package main

import (
	"bufio"
	"fmt"
	pb "github.com/SamuelBagattin/grpc-chat/proto"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"io"
	"log"
	"os"
)

var (
	userName      string
	serverAddress string
)

func main() {
	serverAddress = os.Args[1]
	userName = os.Args[2]
	fmt.Print("Connecting...")
	// Set up a connection to the server.
	conn, err := grpc.Dial(serverAddress, grpc.WithInsecure(), grpc.WithBlock())
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
	print("\rRegistering client")
	clientUid := registerClient(c)
	go subscribeToNewcomers(c, clientUid)

	print("\rOpening SendMessages rpc")
	sendMessagesStream, sendMessageErr := c.SendMessages(context.Background())
	if sendMessageErr != nil {
		log.Fatalf("%v", sendMessageErr)
	}
	print("\rOpening ReceiveMessages rpc")
	receiveMessagesStream, receiveMessagesErr := c.ReceiveMessages(context.Background(), &pb.ClientInfos{ClientUid: clientUid})
	if receiveMessagesErr != nil {
		log.Fatalf("%v", receiveMessagesErr)
	}

	receiveMsg := make(chan *pb.ReceiveMessageResponse)
	cl := make(chan bool)
	fmt.Println("\rConnected !")
	go receiveMessages(receiveMessagesStream, receiveMsg, cl)
	go processMessages(receiveMsg, cl)
	go sendMessages(sendMessagesStream, clientUid)

	<-cl
	if errClose := sendMessagesStream.CloseSend(); errClose != nil {
		panic(errClose)
	}

}

func subscribeToNewcomers(c pb.MessagesServiceClient, clientUid string) {
	newcomersStream, newcomersError := c.SubscribeToNewcomers(context.Background(), &pb.SubscribeToNewcomersRequest{ClientUid: clientUid})
	if newcomersError != nil {
		log.Fatalf("Failed to subscribe to newcomers : %v", newcomersError)
	}
	for {
		in, newcomersStreamError := newcomersStream.Recv()
		if newcomersStreamError == io.EOF {
			// read done.
			log.Fatalf("closing connection")
		}
		if newcomersStreamError != nil {
			log.Fatalf("Failed to receive a newcomer : %v", newcomersStreamError)
		}
		fmt.Printf("%s est arrivé dans la salle !", in.GetName())
	}
}

// Returns the clientUid
func registerClient(c pb.MessagesServiceClient) string {
	//registerClientContext, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	client, registerClientError := c.RegisterClient(context.Background(), &pb.RegisterClientRequest{UserName: userName})
	if registerClientError != nil {
		log.Fatalf("Cannot register client : %v", registerClientError)
	}
	//cancel()
	return client.GetClientUid()
}

func receiveMessages(stream pb.MessagesService_ReceiveMessagesClient, receive chan *pb.ReceiveMessageResponse, cl chan bool) {
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

func processMessages(receive chan *pb.ReceiveMessageResponse, cl chan bool) {
	for {
		select {
		case msg := <-receive:
			fmt.Printf("%s: %s\n", msg.GetName(), msg.GetContent())
		case <-cl:
			fmt.Printf("closing connection")
		}
	}
}

func sendMessages(stream pb.MessagesService_SendMessagesClient, clientUid string) {
	for {
		reader := bufio.NewReader(os.Stdin)
		readString, err := reader.ReadString('\n')
		if err != nil {
			if err == io.EOF {
				continue
			}
			log.Printf("%v", err)
		}
		if err := stream.Send(&pb.SendMessageRequest{
			ClientUid: clientUid,
			Content:   readString,
		}); err != nil {
			log.Fatalf("Failed to send : %v", err)
		}
	}

}
