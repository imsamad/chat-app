package main

import (
	"fmt"
	proto "huge-file/proto"
	"io"
	"net"
	"os"

	"google.golang.org/grpc"
)

type server struct {
	proto.UnimplementedStreamUploadServer
}

func main() {
	listener, tcpErr := net.Listen("tcp", ":9000")

	if tcpErr != nil {
		fmt.Println("Error to start tcp server: ", tcpErr)
		return
	}

	srv := grpc.NewServer()

	proto.RegisterStreamUploadServer(srv, &server{})

	if e := srv.Serve(listener); e != nil {
		fmt.Println("error to serve the service: ", e)
		os.Exit(1)
	}
}

func (s *server) Upload(stream proto.StreamUpload_UploadServer) error {
	var fileBytes []byte
	var fileSize int64 = 0

	for {
		req, err := stream.Recv()

		if err != nil {
			fmt.Println("error while reading messages: ", err)

		}

		if err == io.EOF {
			break
		}
		chunk := req.GetChunks()
		fileBytes = append(fileBytes, chunk...)
		fileSize += int64(len(chunk))
	}

	f, err := os.Create("./abc.bin")

	if err != nil {
		fmt.Println("errior while creating file: ", err)
	}

	defer f.Close()

	_, err = f.Write(fileBytes)

	if err != nil {
		fmt.Println("error while writing to bytes to file: ", err)
	}

	return stream.SendAndClose(&proto.UploadResponse{FileSize: fileSize, Message: "Filed saved!"})
}
