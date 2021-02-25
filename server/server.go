package main

import (
	"CalculatorService/proto"
	"google.golang.org/grpc"
	"io"
	"log"
	"net"
)

type Server struct {
	proto.UnimplementedCalculatorServiceServer
}


func (s *Server) PrimeNumberDecomposition(req *proto.IntRequest,stream proto.CalculatorService_PrimeNumberDecompositionServer)  error{
	number := int(req.Number)
	for number%2 == 0 {
		res := &proto.IntResponse{Number: 2, OptionalError: ""}
		if err := stream.Send(res); err != nil {
			log.Fatalf("error while sending stream responses: %v", err.Error())
		}
		number = number / 2
	}


	for i := 3; i*i <= number; i = i + 2 {
		// while i divides n, append i and divide n
		for number%i == 0 {
			res := &proto.IntResponse{Number: int64(i), OptionalError: ""}
			if err := stream.Send(res); err != nil {
				log.Fatalf("error while sending stream responses: %v", err.Error())
			}
			number = number / i
		}
	}

	if number > 2 {
		res := &proto.IntResponse{Number: int64(number), OptionalError: ""}
		if err := stream.Send(res); err != nil {
			log.Fatalf("error while sending stream responses: %v", err.Error())
		}
	}

	return nil
}

func (s *Server) ComputeAverage(stream proto.CalculatorService_ComputeAverageServer)  error{
	var avg float64
	a := make([]int64, 10)

	for {
		req, err := stream.Recv()
		if err == io.EOF {
			for i := 0; i < len(a); i++ {

				avg += float64(a[i])
			}

			avg := (float64(avg)) / (float64(len(a)))
			return stream.SendAndClose(&proto.FloatResponse{
				Number: avg,
			})

		}
		if err != nil {
			log.Fatalf("Error while reading avg_client stream: %v", err)
		}

		num := req.Number
		a = append(a, num)
	}
}

func main() {
	l, err := net.Listen("tcp", "0.0.0.0:59751")
	if err != nil {
		log.Fatalf("Failed to listen:%v", err)
	}
	s := grpc.NewServer()
	proto.RegisterCalculatorServiceServer(s, &Server{})
	log.Println("Server is running on port:59751")
	if err := s.Serve(l); err != nil {
		log.Fatalf("failed to serve:%v", err)
	}
}
