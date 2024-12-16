package main

import (
	"context"
	"fmt"
	proto "huge-file/proto"
	"io"
	"os"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var client proto.StreamUploadClient

func main() {
	trpcConn, err := grpc.NewClient("localhost:9000", grpc.WithTransportCredentials(insecure.NewCredentials()))

	if err != nil {
		fmt.Println("err: ", err)
		os.Exit(1)
	}
	client = proto.NewStreamUploadClient(trpcConn)

	mb := 1024 * 1024 * 2
	uploadFile("./syslog", mb)
}

func uploadFile(path string, batchSize int) {
	t := time.Now()

	file, err := os.Open(path)

	if err != nil {
		fmt.Println("Error to open the file: ", err)

		os.Exit(1)
	}

	buf := make([]byte, batchSize)

	batchNumber := 1

	stream, err := client.Upload(context.TODO())

	if err != nil {
		fmt.Println("Error to open the stream: ", err)

		os.Exit(1)
	}

	for {
		num, err := file.Read(buf)
		if err == io.EOF {
			fmt.Println("eof")
			break
		}

		if err != nil {
			fmt.Println("Err: ", err)
			return
		}

		chunk := buf[:num]

		err = stream.Send(&proto.UploadRequest{FilePath: path, Chunks: chunk})

		if err != nil {
			fmt.Println("err: ", err)

			return
		}

		batchNumber++
	}

	res, err := stream.CloseAndRecv()

	fmt.Println(err)

	fmt.Println(res)
	fmt.Println(time.Since(t))
}
