# Chitty Chat How to run
Group: Ben Dover

To start the program you have to run this command to start the server:
```console
go run server/server.go
```

and then go to a new terminal, and write this command to create a client:
```console
go run client/client.go -name ""
```
Here You can input a name on the argument ``-name`` by writing it in the ``""``.
You can then create multible clients by repeating the same Command mulitble times, in new terminals.

## How to write in the chat
When you have created one or more clients, it is as simple as writing in the terminal on the clients, and the chats will then be distributed to all Clients, with a Lamport Timestamp, the name of the Client, and the message.

## How to leave the chat
If a Client wants to leave, you can either close the terminal, or use the shortrcut ctrl+c to close the process.
This should provide a "Client left" message to all the remaining clients.

## Running Chitty Chat with different server ports
It is possible to specify which serverport the server should run at. This is done by giving a Flag with the commandline specifying which port. Under here are examples of how to make a server and a client on the server port 5454.
```console
go run server/server.go -port 5454
```
```console
go run client/client.go -sPort 5454 -name ""
```
