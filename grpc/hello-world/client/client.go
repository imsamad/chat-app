package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	proto "grpc/protoc"
)

var client proto.HelloWorldClient

func main() {
	trpcConn, err := grpc.NewClient("localhost:9000", grpc.WithTransportCredentials(insecure.NewCredentials()))

	if err != nil {
		fmt.Println("err: ", err)
	}

	client = proto.NewHelloWorldClient(trpcConn)

	r := chi.NewRouter()

	r.Get("/", defaultMethod)
	r.Get("/client", clientStreaming)
	r.Get("/server", serverStreaming)
	r.Get("/bi", bidirectional)

	http.ListenAndServe(":3000", r)
}

func defaultMethod(w http.ResponseWriter, r *http.Request) {
	res, err := client.Default(context.TODO(), &proto.HelloRequst{Req: "Hello from client"})
	if err != nil {
		fmt.Println("err: ", err)
	} else {
		fmt.Println(res)
	}
	json.NewEncoder(w).Encode(res)
}

func clientStreaming(w http.ResponseWriter, r *http.Request) {
	req := []*proto.HelloRequst{
		{Req: "req 1"},
		{Req: "req 2"},
		{Req: "req 3"},
		{Req: "req 4"},
	}

	stream, err := client.ClientStream(context.TODO())

	if err != nil {
		fmt.Println("err: ", err)
		w.WriteHeader(http.StatusConflict)
		json.NewEncoder(w).Encode("what i this")
		return
	}

	for _, re := range req {
		err := stream.Send(re)

		if err != nil {
			fmt.Println("request not fulfilled")
			return
		}
	}

	response, err := stream.CloseAndRecv()

	if err != nil {
		fmt.Println("err:", err)
		return
	}

	json.NewEncoder(w).Encode(response)

}

func serverStreaming(w http.ResponseWriter, _ *http.Request) {
	stream, err := client.ServerStream(context.TODO(), &proto.HelloRequst{Req: "i am sending request"})

	if err != nil {
		fmt.Println("err while creating stream obejct:", err)
		return
	}

	var messages []*proto.HelloResponse
	count := 0

	for {
		msg, err := stream.Recv()

		if err == io.EOF {
			break
		}

		fmt.Println("replies: ", msg)
		count++

		messages = append(messages, msg)
	}

	json.NewEncoder(w).Encode(messages)

}

func bidirectional(w http.ResponseWriter, r *http.Request) {

	stream, err := client.Bidirectional(context.TODO())
	if err != nil {
		fmt.Println("Err: ", err)
	}

	var messages []*proto.HelloResponse
	for i := 0; i < 10; i++ {
		err := stream.Send(&proto.HelloRequst{Req: "req no: " + strconv.Itoa(i)})
		if err != nil {
			fmt.Println("Err: ", err)
		}
	}

	if err := stream.CloseSend(); err != nil {
		fmt.Println("error while closing send: ", err)
	}
	for {
		res, err := stream.Recv()
		fmt.Println("res bidirectional: ", res)
		if err == io.EOF {
			break
		}
		if err != nil {
			fmt.Println("err: ", err)
		}
		messages = append(messages, res)
	}

	json.NewEncoder(w).Encode(messages)
}
