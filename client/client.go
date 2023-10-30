package main

//go run client/client.go -sPort 5454 -name "..."

import (
	"bufio"
	"context"
	"flag"
	"fmt"
	"log"
	"math/rand"
	"os"
	proto "simpleGuide/grpc"
	"strconv"
	"sync"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type Client struct {
	name string
}

var (
	serverPort = flag.Int("sPort", 5454, "server port number (should match the port used for the server)")
	clientName = flag.String("name", "X", "client name")
)

var waitGroup *sync.WaitGroup
var LamportTimestamp int64

func main() {
	waitGroup = &sync.WaitGroup{}
	done := make(chan int)
	//We start at lamport time 1
	LamportTimestamp = 1

	// Parse the flags to get the port for the client
	flag.Parse()

	// Create a client with name from flag
	client := &Client{
		name: *clientName,
	}

	//Connect to server, returns a proto.StreamingServiceClient
	serverConnection, _ := connectToServer()

	//Create a new proto.Participant-message with random id & name
	participant := &proto.Participant{
		Id:   strconv.Itoa(rand.Intn(100)),
		Name: client.name,
	}

	//Uses the participant object to create a proto.Connect-message to the server, 
	//which in return sends a stream
	connectParticipant(participant, serverConnection)
	time.Sleep(1 * time.Second)

	waitGroup.Add(1)
	go func() {
		defer waitGroup.Done()

		//Creates a new proto.ChatMessage-message saying that the participant has joined the chat
		chatMessage := &proto.ChatMessage{
			Id:          "New participant",
			Participant: participant,
			Message:     "Participant " + participant.Name + " joined the chat!",
			Timestamp:   LamportTimestamp,
		}

		//Sends the above message to the server, and the server returns a proto.Empty-message
		_, err := serverConnection.SendChatMessage(context.Background(), chatMessage)
		if err != nil {
			log.Println("Connection to chatserver failed")
		}

		//Open up the scanner to wait for client's messages
		scanner := bufio.NewScanner(os.Stdin)
		messageId := client.name

		//Infinite for loop that keeps the scanner open until terminal is closed
		for scanner.Scan() {
			inputMessage := scanner.Text()
			LamportTimestamp = LamportTimestamp + 1
			//Creates a new proto.ChatMessage with input from user
			chatMessage := &proto.ChatMessage{
				Id:          messageId,
				Participant: participant,
				Message:     inputMessage,
				Timestamp:   LamportTimestamp,
			}

			//Sends the above ChatMessage to the server, which returns a proto.Empty-message
			_, err := serverConnection.SendChatMessage(context.Background(), chatMessage)
			if err != nil {
				log.Println("Connection to chatserver failed")
				break
			}
		}
	}()

	go func() {
		waitGroup.Wait()
		close(done)
	}()

	<-done
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

func connectParticipant(participant *proto.Participant, client proto.StreamingServiceClient) error {
	var streamError error

	//Use info from Participant to make a proto.Connect-message that we send to server
	//Server returns a stream of ChatMessages
	stream, err := client.GetChatMessageStreaming(context.Background(), &proto.Connect{
		Participant: participant,
		Active:      true,
	})

	if err != nil {
		log.Println("Connection to chatserver failed")
		return (err)
	}

	waitGroup.Add(1)

	go func(chatStream proto.StreamingService_GetChatMessageStreamingClient) {
		defer waitGroup.Done()

		//Infinite for loop to listen to the stream and print out the received ChatMessages
		for {
			chatMessage, err := chatStream.Recv()
			if len(chatMessage.String()) > 180 {
				streamError = fmt.Errorf("message too long %v", err)
				break
			}
			if err != nil {
				streamError = fmt.Errorf("error reading message: %v", err)
				break
			}
			//Update the client's lamport timestamp
			LamportTimestamp = Max(LamportTimestamp, chatMessage.Timestamp) + 1
			log.Printf("Lamport time: %d | %s : %s", LamportTimestamp, chatMessage.Id, chatMessage.Message)
		}
	}(stream) //the go func takes the stream as an argument

	return streamError
}

//Helper method used to find the maximum of two lamport timestamps
func Max(x, y int64) int64 {
	if x < y {
		return y
	}
	return x
}
