package main

//go run client/client.go -cPort 8080 -sPort 5454

import (
	"bufio"
	"context"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	proto "simpleGuide/grpc"
	"strconv"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type Client struct {
	id         int
	portNumber int
}

var (
	clientPort = flag.Int("cPort", 0, "client port number")
	serverPort = flag.Int("sPort", 0, "server port number (should match the port used for the server)")
)

func main() {
	// Parse the flags to get the port for the client
	flag.Parse()

	// Create a client
	client := &Client{
		id:         1,
		portNumber: *clientPort,
	}

	// Wait for the client (user) to ask for the time
	go waitForChatMessages(client)

	for {

	}
}

func waitForChatMessages(client *Client) {
	// Connect to the server
	serverConnection, _ := connectToServer()

	stream, err := serverConnection.GetChatMessageStreaming(context.Background(), &proto.PublishChatMessage{
		ClientId: int64(client.id),
		Message: fmt.Sprintf("Participant %d  joined Chitty-Chat at Lamport time L", client.id),
	})
	if err != nil {
		log.Printf(err.Error())
	}

	resp, err := stream.Recv()
			if err == io.EOF {
				return
			} else if err == nil {
				valStr := fmt.Sprintf("Response\n Part: %s \n Val: %s", resp.ServerName, resp.Message)
				log.Println(valStr)
			}

	// Wait for input in the client terminal
	
	for {
		scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		
		input := scanner.Text()
		//log.Printf("Client asked for time with input: %s\n", input)

		// Send a message
		stream, err := serverConnection.GetChatMessageStreaming(context.Background(), &proto.PublishChatMessage{
			ClientId: int64(client.id),
			Message: input,
		})
		if err != nil {
			log.Printf(err.Error())
		}

		
			resp, err := stream.Recv()
			if err == io.EOF {
				return
			} else if err == nil {
				valStr := fmt.Sprintf("Response\n Part: %s \n Val: %s", resp.ServerName, resp.Message)
				log.Println(valStr)
			}
	
			if err != nil {
				log.Printf(err.Error())
			}
	
		
	}
	}
}

func connectToServer() (proto.StreamingServiceClient, error) {
	// Dial the server at the specified port.
	conn, err := grpc.Dial("localhost:"+strconv.Itoa(*serverPort), grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("Could not connect to port %d", *serverPort)
	} else {
		log.Printf("Connected to the server at port %d\n", *serverPort)
	}
	return proto.NewStreamingServiceClient(conn), nil
}
