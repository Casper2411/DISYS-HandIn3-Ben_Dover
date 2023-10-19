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

type Connection struct {
	stream proto.StreamingService_GetChatMessageStreamingServer
	id     string
	active bool
	error  chan error
}

// Struct that will be used to represent the Server.
type Server struct {
	proto.UnimplementedStreamingServiceServer // Necessary
	name                             string
	port                             int
	connection []*Connection
}

// Used to get the user-defined port for the server from the command line
var port = flag.Int("port", 0, "server port number")

func main() {
	// Get the port from the command line when the server is run
	flag.Parse()

	var connections []*Connection

	// Create a server struct
	server := &Server{
		name: "serverName",
		port: *port,
		connection: connections,
	}

	// Start the server
	startServer(server)

	// Keep the server running until it is manually quit
	/*for {

	}*/
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
	/*log.Printf("Client with ID %d asked for the time\n", in.ClientId)
	return &proto.{
		Time:       time.Now().String(),
		ServerName: c.name,
	}, nil*/
	//log.Println("Fetch data streaming")

	//for i := 0; i < 10; i++ {
		//value := randStringBytes(5)

		conn := &Connection{
			stream: chatStream,
			id:		connection.Participant.Id,
			active: true,
			error: make(chan error),
		}
		
		server.connection = append(server.connection, conn)

		/*resp := proto.BroadcastChatMessage{
			ServerName: strconv.Itoa(14),
			Message: conn.Message,
		}

		if err := srv.Send(&resp); err != nil {
			log.Println("error generating response")
			return err
		}
	//}*/

	return <-conn.error
}

//SendChatMessage(context.Context, *ChatMessage) (*Empty, error)
func (server *Server) SendChatMessage(ctx context.Context, chatMessage *proto.ChatMessage) (*proto.Empty, error){
	waitGroup := sync.WaitGroup{}
	done := make(chan int)

	for _, connection := range server.connection{
		waitGroup.Add(1)

		go func(chatMessage *proto.ChatMessage, connection *Connection){
			defer waitGroup.Done()

			if(connection.active){
				err := connection.stream.Send(chatMessage)
				log.Println("Sending message from ", connection.id)
				
				if err!=nil {
					log.Println("Error: ", err)
					connection.active = false
					connection.error <- err
				}
			}
		}(chatMessage, connection)
	}

	go func(){
		waitGroup.Wait()
		close(done)
	}()

	<-done
	log.Println("end of method test")
	return &proto.Empty{}, nil
}


/*func randStringBytes(n int) string {
	const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

	b := make([]byte, n)
	for i := range b {
		b[i] = letterBytes[rand.Intn(len(letterBytes))]
	}
	return string(b)
}*/
