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

//Struct to represent a Connection from a client to the server
type Connection struct {
	stream proto.StreamingService_GetChatMessageStreamingServer
	id     string
	active bool
	error  chan error
}

// Struct to represent the server
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
	//We start at lamport time 1
	LamportTimestamp = 1

	//A slice of connections
	var connections []*Connection

	// Create a server struct with the port and the above slice of connections
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

//From the proto-file: GetChatMessageStreaming(*Connect, StreamingService_GetChatMessageStreamingServer) error
func (server *Server) GetChatMessageStreaming(connection *proto.Connect, chatStream proto.StreamingService_GetChatMessageStreamingServer) error {

	//Create new Connection struct object
	conn := &Connection{
		stream: chatStream,
		id:     connection.Participant.Id,
		active: true,
		error:  make(chan error),
	}

	//Append the Connection to the server slice of active connections
	server.connection = append(server.connection, conn)

	//This channel can only be emptied when context is done, 
	//which means the client has ended the connection
	<- chatStream.Context().Done()
	waitGroup := sync.WaitGroup{}
	LamportTimestamp = LamportTimestamp + 1
	
	//Remove the now inactive client-connection from the server's slice of connections
	server.connection = removeElement(server.connection, conn)
	
	//Broadcast to remaining clients that participant has left
	server.sendToAllConnections(&waitGroup, &proto.ChatMessage{Id: "Participant leaving", Participant: connection.Participant, Message: "Participant " + connection.Participant.Name + " is leaving", Timestamp: LamportTimestamp})
	log.Printf("Lamport time %d | Broadcasting that " + connection.Participant.Name + " has left the chat", LamportTimestamp)
	return <-conn.error
}

//From proto-file: SendChatMessage(context.Context, *ChatMessage) (*Empty, error)
func (server *Server) SendChatMessage(ctx context.Context, chatMessage *proto.ChatMessage) (*proto.Empty, error) {
	waitGroup := sync.WaitGroup{}
	done := make(chan int)

	//Compare the server's timestamp to the recieved ChatMessage's timestamp
	LamportTimestamp = Max(LamportTimestamp, chatMessage.Timestamp) + 1
	log.Printf("Lamport time %d | Receiving message from %s", LamportTimestamp, chatMessage.Participant.Name)
	
	//Further increment the server's timestamp
	LamportTimestamp = LamportTimestamp + 1
	log.Printf("Lamport time %d | Broadcasting message from %s", LamportTimestamp, chatMessage.Participant.Name)
	
	//Update the ChatMessage's timestamp before broadcasting it
	chatMessage.Timestamp = LamportTimestamp 

	//Broadcast the ChatMessage by sending the message to all connections in the connection-slice
	server.sendToAllConnections(&waitGroup, chatMessage)

	go func() {
		waitGroup.Wait()
		close(done)
	}()

	<-done
	return &proto.Empty{}, nil
}

func (server *Server) sendToAllConnections(waitGroup *sync.WaitGroup, chatMessage *proto.ChatMessage) {
	//For each connection in the server's connection-slice, send the ChatMessage into their stream
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
		}(chatMessage, connection) //The go func takes the ChatMessage and Connection as arguments
	}
}

//Helper method used to find the maximum of two lamport timestamps
func Max(x, y int64) int64 {
	if x < y {
		return y
	}
	return x
}

//Helper method used to remove a specific element from a clice
func removeElement(slice []*Connection, element *Connection) []*Connection {
    for i, v := range slice {
        if v == element {
            return append(slice[:i], slice[i+1:]...)
        }
    }
    return slice // Element not found
}

