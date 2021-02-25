package server

import (
	//"context"
	"log"
	"net"
	"fmt"
	calculator "../averagepb"
	"google.golang.org/grpc"
	//"time"
	"io"
)

type server struct {}

func (*server) ComputeAvg(stream computeAverage.ComputeAverageService_ComputeAvgServer) error {
	fmt.Println("Invoking ComputeAvg")
	counter := float64(0)
	sum := float64(0)
	for {
		req, err := stream.Recv()
		if err == io.EOF {
			fmt.Printf("Sending back result: %v\n", sum / counter)
			return stream.SendAndClose(&computeAverage.ComputeAverageResponse {
				Average : sum / counter,
			})
		}
		if err != nil {
			log.Fatalf("Error while reading client stream %v", err)
		}
		counter++
		sum += float64(req.GetVal())
	}
}

func main() {
	fmt.Println("Hello World")

	lis, err := net.Listen("tcp", "0.0.0.0:50051")
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	s := grpc.NewServer()

	computeAverage.RegisterComputeAverageServiceServer(s, &server{})

	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}

}
© 2021 GitHub, Inc.
