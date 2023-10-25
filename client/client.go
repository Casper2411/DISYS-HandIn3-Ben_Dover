package main

//go run client/client.go -cPort 8080 -sPort 5454 -name ""

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
	name       string
	portNumber int
}

//var client proto.StreamingServiceClient

var (
	serverPort = flag.Int("sPort", 5454, "server port number (should match the port used for the server)")
	clientName = flag.String("name", "X", "client name")
)

var waitGroup *sync.WaitGroup
var LamportTimestamp int64

func main() {
	waitGroup = &sync.WaitGroup{}
	done := make(chan int)
	LamportTimestamp = 1

	// Parse the flags to get the port for the client
	flag.Parse()

	// Create a client
	client := &Client{
		name: *clientName,
	}

	serverConnection, _ := connectToServer()

	participant := &proto.Participant{
		Id:   strconv.Itoa(rand.Intn(100)),
		Name: client.name,
	}

	connectParticipant(participant, serverConnection)
	time.Sleep(1 * time.Second)

	waitGroup.Add(1)
	go func() {
		defer waitGroup.Done()

		chatMessage := &proto.ChatMessage{
			Id:          "New participant",
			Participant: participant,
			Message:     "Participant " + participant.Name + " joined the chat! Say hello!",
			Timestamp:   LamportTimestamp,
		}

		_, err := serverConnection.SendChatMessage(context.Background(), chatMessage)
		if err != nil {
			log.Println("Connection to chatserver failed")
		}

		scanner := bufio.NewScanner(os.Stdin)
		//timestamp := time.Now()
		messageId := client.name

		for scanner.Scan() {
			inputMessage := scanner.Text()
			LamportTimestamp = LamportTimestamp + 1
			chatMessage := &proto.ChatMessage{
				Id:          messageId,
				Participant: participant,
				Message:     inputMessage,
				Timestamp:   LamportTimestamp,
			}

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

		for {
			chatMessage, err := chatStream.Recv()
			fmt.Println(len(chatMessage.String()))
			if len(chatMessage.String()) > 180 {
				fmt.Println(len(chatMessage.String()), "wwwoowoowowo")
				streamError = fmt.Errorf("Message too long %v", err)
				break
			}
			if err != nil {
				streamError = fmt.Errorf("error reading message: %v", err)
				break
			}
			LamportTimestamp = Max(LamportTimestamp, chatMessage.Timestamp) + 1
			log.Printf("Lamport time: %d | %s : %s", LamportTimestamp, chatMessage.Id, chatMessage.Message)
		}
	}(stream)

	return streamError
}

func Max(x, y int64) int64 {
	if x < y {
		return y
	}
	return x
}
