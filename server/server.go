package main

//go run server/server.go -port 5454

import (
	"context"
	"flag"
	"log"
	"net"
	proto "simpleGuide/grpc"
	"strconv"
	"sync"
	"google.golang.org/grpc"
)

//Struct that will be used to represent a Connection from a client to the server
type Connection struct {
	stream proto.StreamingService_GetChatMessageStreamingServer
	id     string
	active bool
	error  chan error
}

// Struct that will be used to represent the Server.
type Server struct {
	proto.UnimplementedStreamingServiceServer // Necessary
	name                                      string
	port                                      int
	connection                                []*Connection
}

// Used to get the user-defined port for the server from the command line
var port = flag.Int("port", 5454, "server port number")
var LamportTimestamp int64

func main() {
	// Get the port from the command line when the server is run
	flag.Parse()
	LamportTimestamp = 1

	var connections []*Connection

	// Create a server struct
	server := &Server{
		name:       "serverName",
		port:       *port,
		connection: connections,
	}

	// Start the server
	startServer(server)
}

func startServer(server *Server) {

	// Create a new grpc server
	grpcServer := grpc.NewServer()

	// Make the server listen at the given port (convert int port to string)
	listener, err := net.Listen("tcp", ":"+strconv.Itoa(server.port))

	if err != nil {
		log.Fatalf("Could not create the server %v", err)
	}
	log.Printf("Started server at port: %d\n", server.port)

	// Register the grpc server and serve its listener
	proto.RegisterStreamingServiceServer(grpcServer, server)
	serveError := grpcServer.Serve(listener)
	if serveError != nil {
		log.Fatalf("Could not serve listener")
	}
}

//GetChatMessageStreaming(*Connect, StreamingService_GetChatMessageStreamingServer) error

func (server *Server) GetChatMessageStreaming(connection *proto.Connect, chatStream proto.StreamingService_GetChatMessageStreamingServer) error {

	conn := &Connection{
		stream: chatStream,
		id:     connection.Participant.Id,
		active: true,
		error:  make(chan error),
	}

	server.connection = append(server.connection, conn)

	<- chatStream.Context().Done()
	waitGroup := sync.WaitGroup{}
	LamportTimestamp = LamportTimestamp + 1
	server.connection = removeElement(server.connection, conn)
	server.sendToAllConnections(&waitGroup, &proto.ChatMessage{Id: "Participant leaving", Participant: connection.Participant, Message: "Participant " + connection.Participant.Name + " is leaving", Timestamp: LamportTimestamp})
	log.Printf("Lamport time %d | Broadcasting that " + connection.Participant.Name + " has left the chat", LamportTimestamp)
	return <-conn.error
}

// SendChatMessage(context.Context, *ChatMessage) (*Empty, error)
func (server *Server) SendChatMessage(ctx context.Context, chatMessage *proto.ChatMessage) (*proto.Empty, error) {
	waitGroup := sync.WaitGroup{}
	done := make(chan int)

	LamportTimestamp = Max(LamportTimestamp, chatMessage.Timestamp) + 1
	log.Printf("Lamport time %d | Receiving message from %s", LamportTimestamp, chatMessage.Participant.Name)
	LamportTimestamp = LamportTimestamp + 1
	log.Printf("Lamport time %d | Broadcasting message from %s", LamportTimestamp, chatMessage.Participant.Name)
	chatMessage.Timestamp = LamportTimestamp 

	server.sendToAllConnections(&waitGroup, chatMessage)

	go func() {
		waitGroup.Wait()
		close(done)
	}()

	<-done
	return &proto.Empty{}, nil
}

func (server *Server) sendToAllConnections(waitGroup *sync.WaitGroup, chatMessage *proto.ChatMessage) {
	for _, connection := range server.connection {
		waitGroup.Add(1)
	
		go func(chatMessage *proto.ChatMessage, connection *Connection) {
			defer waitGroup.Done()
	
			if connection.active {
				err := connection.stream.Send(chatMessage)
				
				if err != nil {
					log.Println("Error: ", err)
					connection.active = false
					connection.error <- err
				}
			}
		}(chatMessage, connection)
	}
}



func Max(x, y int64) int64 {
	if x < y {
		return y
	}
	return x
}

//From ChatGPT
func removeElement(slice []*Connection, element *Connection) []*Connection {
    for i, v := range slice {
        if v == element {
            return append(slice[:i], slice[i+1:]...)
        }
    }
    return slice // Element not found
}

