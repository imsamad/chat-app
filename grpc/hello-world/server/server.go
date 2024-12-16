package main

import (
	"context"
	"fmt"
	"io"
	"net"
	"strconv"

	proto "grpc/protoc" // Importing the generated gRPC code

	"google.golang.org/grpc"
)

// server struct implements the generated HelloWorldServer interface
type server struct {
	proto.UnimplementedHelloWorldServer // Embedding to fulfill the interface requirements
}

func main() {
	// Step 1: Create a TCP listener on port 9000
	listener, tcpErr := net.Listen("tcp", ":9000")
	if tcpErr != nil {
		// Log and exit if the listener fails to start
		fmt.Println("Failed to start server listener, reason:", tcpErr)
		return
	}
	fmt.Println("Server is listening on port 9000")

	// Step 2: Create a new gRPC server instance
	srv := grpc.NewServer()

	// Step 3: Register the HelloWorldServer implementation with the gRPC server
	proto.RegisterHelloWorldServer(srv, &server{})
	// fmt.Println("gRPC server registered")

	// Step 4: Register server reflection to make it easier to explore with tools like `grpc_cli`
	// reflection.Register(srv)
	// fmt.Println("gRPC server reflection enabled")

	// Step 5: Start serving requests
	err := srv.Serve(listener)
	if err != nil {
		// Log any error that occurs while serving
		fmt.Println("Failed to start serving gRPC requests, reason:", err)
	}
}

// ServerReply handles incoming requests and sends back responses
func (s *server) Default(ctx context.Context, req *proto.HelloRequst) (*proto.HelloResponse, error) {
	// Log the received request for debugging purposes
	fmt.Printf("Received request from client: %s\n", req.Req)

	// Construct a response and log the outgoing response
	response := &proto.HelloResponse{
		Res: "string", // Example static response, can be dynamic
	}
	fmt.Printf("Sending response to client: %s\n", response.Res)

	return response, nil
}

func (s *server) ClientStream(stream proto.HelloWorld_ClientStreamServer) error {
	total := 0
	for {
		req, err := stream.Recv()

		if err == io.EOF {
			fmt.Println("no more to read: ", err)
			return stream.SendAndClose(&proto.HelloResponse{
				Res: "total messages sent: " + strconv.Itoa(total),
			})
		}
		if err != nil {
			return err
		}
		total++
		fmt.Printf("%dth HelloMessage recieved from client %v\n", total, req)
	}
}

func (s *server) ServerStream(req *proto.HelloRequst, stream proto.HelloWorld_ServerStreamServer) error {
	fmt.Println("req: ", req)
	replies := []*proto.HelloResponse{
		{Res: "response 1"},
		{Res: "response 2"},
		{Res: "response 3"},
		{Res: "response 4"},
	}

	for _, reply := range replies {
		err := stream.Send(reply)

		if err != nil {
			fmt.Println("error on while sending : ", err)
			return err
		}

	}
	return nil
}

func (s *server) Bidirectional(stream proto.HelloWorld_BidirectionalServer) error {
	for i := 0; i < 10; i++ {
		err := stream.Send(&proto.HelloResponse{Res: "res: " + strconv.Itoa(i)})
		if err != nil {
			fmt.Println("unable to send data from server: resons ", err)
			return err
		}
	}

	for {
		req, err := stream.Recv()
		fmt.Println("err :", err)
		if err == io.EOF {
			fmt.Println("breakL ")
			break
		}

		if err != nil {

			fmt.Println("err while recv stream: ", err)
			// return err
		}

		fmt.Println("message recieved: ", req)
	}
	return nil
}
