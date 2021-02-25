package main

import (
	"context"
	"fmt"
	"google.golang.org/grpc"
	"CalculatorService/proto"
	"io"
	"log"
)

func main() {

	fmt.Println("Client started")

	conn, err := grpc.Dial("localhost:59751", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("could not connect: %v", err)
	}
	defer conn.Close()

	c := proto.NewCalculatorServiceClient(conn)
	getDecomposNums(c)
}

func getDecomposNums(c proto.CalculatorServiceClient) {
	ctx := context.Background()
	req := &proto.IntRequest{Number: 120}
	stream, err := c.PrimeNumberDecomposition(ctx, req)
	if err != nil {
		log.Fatalf("error while calling PrimeNumDecompos RPC %v", err)
	}
	defer stream.CloseSend()

LOOP:
	for {
		res, err := stream.Recv()
		if err != nil {
			if err == io.EOF {
				// we've reached the end of the stream
				break LOOP
			}
			log.Fatalf("error while reciving from PrimeNumDecompos RPC %v", err)
		}
		log.Printf("response from PrimeNumDecompos:%v \n", res.Number)
	}

}
