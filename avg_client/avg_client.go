package main

import (
	"context"
	"CalculatorService/proto"
	"fmt"
	"google.golang.org/grpc"
	"log"
	"time"
)

func main() {

	fmt.Println("Client started")

	conn, err := grpc.Dial("localhost:59751", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("could not connect: %v", err)
	}
	defer conn.Close()

	c := proto.NewCalculatorServiceClient(conn)
	getAverage(c)
}

func getAverage(c proto.CalculatorServiceClient) {

	requests := []*proto.IntRequest{
		{
			Number: 1,
		},
		{
			Number: 2,
		},
		{
			Number: 3,
		},
		{
			Number: 4,
		},
	}

	ctx := context.Background()
	stream, err := c.ComputeAverage(ctx)
	if err != nil {
		log.Fatalf("error while calling ComputeAverage: %v", err)
	}

	for _, req := range requests {
		fmt.Printf("Sending requests with num: %v\n", req.Number)
		stream.Send(req)
		time.Sleep(1000 * time.Millisecond)
	}

	res, err := stream.CloseAndRecv()
	if err != nil {
		log.Fatalf("error while receiving response from ComputeAverage: %v", err)
	}
	fmt.Printf("ComputeAverage Response: %v\n", res.Number)
}


