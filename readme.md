# Time Request gRPC Guide

- [gRPC guide](#Time-Request-gRPC-Guide)
    - [Setting Up Files](#setting-up-files)
    - [Populating Proto File](#populating-proto-file)
    - [Setting Up Server](#setting-up-server)
      - [Making the Server Run](#making-the-server-run)
      - [Implementing the Server Function](#implementing-the-server-function)
      - [Summary of the Process](#summary-of-the-process)
    - [Setting Up the Client](#setting-up-the-client)
    - [How to Run](#how-to-run)
    - [Important Notes](#important-notes)

This guide provides an overview of how to implement a simple scenario with grpc (see [here](https://github.com/PatrickMatthiesen/DSYS-gRPC-template) for a more complex one). We want to construct a scenario with two entities:
* A server
* A client

and with the functionality that:
* The client can ask the server for the current time.
* The server will respond with its server name and the current time.

NOTE: while going through the guide, refer to the files in the repo for more context and details.

## Setting Up Files
1. Open an empty folder in your editor.
2. Run `go mod init someName` where 'someName' can be anything relevant. This will create a `go.mod` file with your chosen name at the top.
3. Create a folder named grpc (or something else) with a file proto file inside (for example `proto.proto`).
4. Create a client folder with a `client.go` file.
5. Create a server folder with a `server.go` file.
6. Put both the client and server into the same package by stating the package at the top of the files (for example ``package main``).

## Populating Proto File
Inside the `proto.proto` file we need to set up the service, service functions, and message type(s) we need. The functions are used for communicating between the entities, while the message types are used for input/output in the functions. The service contains the functions.

In our case, in the proto file, we want to have a:
* TimeAsk service
* function for the client to ask the server for the time. 
* Type for the message the client will send to ask for the time.
* Type for the message the server will send back to the client with the time. 

Therefore, we add (where line 2-3 depends on the values above)
```proto
syntax = "proto3";

package simpleGuide;

option go_package = "grpc/proto";

message AskForTimeMessage {
  int64 clientId = 1;
}

message TimeMessage {
  string serverName = 1;
  string time = 2;
}

service TimeAsk {
  rpc AskForTime(AskForTimeMessage) returns (TimeMessage);
}
```
Then we run the command:

`protoc --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative grpc/proto.proto`

where the last part specifies the path to the proto file. This should create two files inside the grpc folder (proto.pb.go and proto_grpc.go).
These files contain autogenerated code based on the things we implemented in the proto file.

## Setting Up Server
NOTE: imports are not always explicitly mentioned here (such as `time`, etc.). They should be added to the imports at the top. If in doubt, refer to the files in the repo.

### Making the Server Run

Now we need to set up the server. We will let the user determine which port it should run on via the flag "-port". 

1. Create a server struct with the necessary fields. You can read more about why the unimplemented part is needed [here](https://stackoverflow.com/questions/69700899/grpc-error-about-missing-an-unimplemented-server-function).

```go
type Server struct {
	proto.UnimplementedTimeAskServer // Necessary
	name                             string
	port                             int
}
```

2. Add a port variable and create a main that parses the flag, sets up a server struct, starts the server, and ensures that main will keep running.

```go
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
```

NOTE: You will need to add `"google.golang.org/grpc"`to the imports section and then run `go mod tidy` for the above code to work (specifically to be able to use the NewServer function).

Now you should be able to start the server at some port like ``go run server/server.go -port 5454`` and get the output `some_time_stamp Started server at port: 5454` (while the server keeps running).

### Implementing the Server Function

We now need to implement the `AskForTime` function we defined in the proto file, so that the client can call it.

To do so, we have to find the corresponding function inside the `proto_grpc.pb.go` file and copy its signature (in my case on line 32):

```go
func (c *timeAskClient) AskForTime(ctx context.Context, in *AskForTimeMessage, opts ...grpc.CallOption) (*SendTime, error)
```

Copy this into the ``server.go`` file, turn it into a function, and change some things:
* Import the context package. 
* Change timeAskClient to Server, so that it matches our server struct definition.
* Add *proto.AskForTimeMessage instead of just AskForTimeMessage as the `in` value (to point to the right type in the proto file).
* Add *proto.SendTime instead of just SendTime as the output.
* Remove the `opts ...grpc.CallOption` from the input.
* Temporarily make it return ``nil, nil``

We now have:
```go
func (c *Server) AskForTime(ctx context.Context, in *proto.AskForTimeMessage) (*proto.TimeMessage, error) {
    return nil, nil
}
```
and there should not be any compile-errors. 

Now we will make the server actually return something and log which client is requesting the time. To do so, update the function to:
```go
func (c *Server) AskForTime(ctx context.Context, in *proto.AskForTimeMessage) (*proto.TimeMessage, error) {
	log.Printf("Client with ID %d asked for the time\n", in.ClientId)
	return &proto.TimeMessage{Time: time.Now().String()}, nil
}
```

### Summary of the Process
Brief overview of what we have done so far with the server:
1. Created a server struct.
2. Created a variable for getting the desired port from the user. 
3. Created a main which parses the port flag, creates an instance of the server struct, starts the server, and keeps running.
4. Implemented a function we specified in the proto file. 

NOTE: if you have more functions in your proto file for the `TimeAsk` service, then you also need to implement those.

## Setting Up the Client

Similar to the server create a client struct, a port variable, set up the struct in main, and ensure that the client keeps running:
```go
type Client struct {
	id int
}

var (
	port = flag.Int("port", 0, "client port number")
)

func main() {
	// Parse the flags to get the port for the client
	flag.Parse()

	// Create a client
	client := &Client{
		id: *port,
	}

	// Wait for the client (user) to ask for the time
	go waitForTimeRequest(client)

	for {
		
	}
}
```

Now we have to implement a function to connect to the server that can be used to obtain a connection, which can be used in the `waitForTimeRequest` function.
```go
func connectToServer() (proto.TimeAskClient, error) {
	// Dial the server at the specified port.
	conn, err := grpc.Dial("localhost:"+strconv.Itoa(*serverPort), grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("Could not connect to port %d", *serverPort)
	} else {
		log.Printf("Connected to the server at port %d\n", *serverPort)
	}
	return proto.NewTimeAskClient(conn), nil
}
```

Lastly, we can implement the actual function:

```go
func waitForTimeRequest(client *Client) {
	// Connect to the server
	serverConnection, _ := connectToServer()

	// Wait for input in the client terminal
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		input := scanner.Text()
		log.Printf("Client asked for time with input: %s\n", input)

		// Ask the server for the time
		timeReturnMessage, err := serverConnection.AskForTime(context.Background(), &proto.AskForTimeMessage{
			ClientId: int64(client.id),
		})

		if err != nil {
			log.Printf(err.Error())
		} else {
			log.Printf("Server %s says the time is %s\n", timeReturnMessage.ServerName, timeReturnMessage.Time)
		}
	}
}
```

We first set up a connection to the server with the helper function. Then we create a scanner and wait for any input (+ enter keyboard press) from the user, after which we get the time from the server. 

## How to Run

1. Run the server: `go run server/server.go -port 5454`.
2. In a different terminal, run the client: `go run client/client.go -cPort 8080 -sPort 5454`.
3. In the client terminal input something and press enter. You should now get the time and be able to see in the server terminal that the client requested the time.

Below is an example of running the server and client:

![Example of running the server and client](running_example.png?raw=true "Example of Running the Server and Client")

## Important Notes
* Keep track of naming. For example, if the name of the service is changed to YearAsk then we e.g. also have to change the return statement in ``connectToServer`` to instead use `proto.NewYearAskClient`. This also goes for message and function changes.
* Make sure you have imported `"google.golang.org/grpc"` in the server and ran `go mod tidy`. 
* Ensure that you are using * and & properly. If you feel a bit shaky about this, you might want to consider using IntelliJ with the Go plugin, since it provides better intellisense.
