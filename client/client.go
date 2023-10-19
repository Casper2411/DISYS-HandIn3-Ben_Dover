package main

//go run client/client.go -cPort 8080 -sPort 5454 -name ""

//MÅSKE IMPLEMENTER TO GOROUTINES - den ene lytter hele tiden i streamen, og den anden holder scanneren åben til at skrive beskeder

import (
	"bufio"
	"context"
	"flag"
	"fmt"
	//"go/scanner"
	//"io"
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
	clientPort = flag.Int("cPort", 0, "client port number")
	serverPort = flag.Int("sPort", 0, "server port number (should match the port used for the server)")
	clientName = flag.String("name", "X", "client name")
)

var waitGroup *sync.WaitGroup

func main() {
	waitGroup = &sync.WaitGroup{}
	done := make(chan int)

	// Parse the flags to get the port for the client
	flag.Parse()

	// Create a client
	client := &Client{
		name:       *clientName,
		portNumber: *clientPort,
	}

	serverConnection, _ := connectToServer()

	participant := &proto.Participant{
		Id: strconv.Itoa(rand.Intn(100)),
		Name: client.name,
	}

	connectParticipant(participant, serverConnection)

	waitGroup.Add(1)
	go func(){
		defer waitGroup.Done()

		scanner := bufio.NewScanner(os.Stdin)
		timestamp := time.Now()
		messageId := client.name

		for scanner.Scan(){
			inputMessage := scanner.Text()
			
			chatMessage := &proto.ChatMessage{
				Id: messageId,
				Participant: participant,
				Message: inputMessage,
				Timestamp: timestamp.String(),
			}

			_, err := serverConnection.SendChatMessage(context.Background(), chatMessage)
			if err!=nil{
				log.Println("Connection to chatserver failed")
				break
			}
		}
	}()

	go func(){
		waitGroup.Wait()
		close(done)
	}()

	<-done

	// Wait for the client (user) to ask for the time
	/*go listenToStream(client, serverConnection)
	go sendChat(client, serverConnection)

	for {

	}*/

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

func connectParticipant(participant *proto.Participant, client proto.StreamingServiceClient) error{
	var streamError error

	stream, err := client.GetChatMessageStreaming(context.Background(), &proto.Connect{
		Participant: participant,
		Active: true,
	})

	if err!=nil {
		log.Println("Connection to chatserver failed")
		return(err)
	}

	waitGroup.Add(1)

	go func(chatStream proto.StreamingService_GetChatMessageStreamingClient){
		defer waitGroup.Done()

		for{
			chatMessage, err := chatStream.Recv()
			if err!=nil{
				streamError = fmt.Errorf("error reading message: %v", err)
				break
			}
			log.Printf("%s : %s", chatMessage.Id, chatMessage.Message)
		} 
	}(stream)
	
	return streamError
}

/*func listenToStream(client *Client, serverConnection proto.StreamingServiceClient) {
	// Connect to the server
	

	stream, err := serverConnection.GetChatMessageStreaming(context.Background(), &proto.PublishChatMessage{
		ClientId: int64(client.id),
		Message: fmt.Sprintf("Participant %d  joined Chitty-Chat at Lamport time L", client.id),
	})
	if err != nil {
		log.Printf(err.Error())
	}

	for{
		resp, err := stream.Recv()
		if err == io.EOF {
			return
		} else if err == nil {
			valStr := fmt.Sprintf("Response\n Part: %s \n Val: %s", resp.ServerName, resp.Message)
			log.Println(valStr)
		}

	}
}

//skal streamen måske med som argument?
func sendChat(client *Client, serverConnection proto.StreamingServiceClient){
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
}*/


