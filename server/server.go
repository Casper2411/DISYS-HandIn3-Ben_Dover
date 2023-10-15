package main

//go run server/server.go -port 5454

import (
	"flag"
	"log"
	"math/rand"
	"net"
	proto "simpleGuide/grpc"
	"strconv"
	"google.golang.org/grpc"
)

// Struct that will be used to represent the Server.
type Server struct {
	proto.UnimplementedStreamingServiceServer // Necessary
	name                             string
	port                             int
}

// Used to get the user-defined port for the server from the command line
var port = flag.Int("port", 0, "server port number")

func main() {
	// Get the port from the command line when the server is run
	flag.Parse()

	// Create a server struct
	server := &Server{
		name: "serverName",
		port: *port,
	}

	// Start the server
	go startServer(server)

	// Keep the server running until it is manually quit
	for {

	}
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

func (c *Server) GetChatMessageStreaming(req *proto.PublishChatMessage, srv proto.StreamingService_GetChatMessageStreamingServer) error {
	/*log.Printf("Client with ID %d asked for the time\n", in.ClientId)
	return &proto.{
		Time:       time.Now().String(),
		ServerName: c.name,
	}, nil*/
	//log.Println("Fetch data streaming")

	//for i := 0; i < 10; i++ {
		//value := randStringBytes(5)

		resp := proto.BroadcastChatMessage{
			ServerName: strconv.Itoa(14),
			Message: req.Message,
		}

		if err := srv.Send(&resp); err != nil {
			log.Println("error generating response")
			return err
		}
	//}

	return nil
}

func randStringBytes(n int) string {
	const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

	b := make([]byte, n)
	for i := range b {
		b[i] = letterBytes[rand.Intn(len(letterBytes))]
	}
	return string(b)
}
