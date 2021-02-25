package main

import (
	"log"
	"net"

	"github.com/muratbekauyez/prime/primepb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/reflection"
	"google.golang.org/grpc/status"
)

type server struct{}

func (*server) Prime(req *primepb.PrimeRequest, stream primepb.PrimeService_PrimeServer) error {
	log.Printf("Received Prime gRPC request with number to decompose: %d", req.GetNumber())

	n := req.GetNumber()
	if n < 2 {
		return status.Errorf(codes.InvalidArgument, "number must be greater than 1")
	}

	k := int64(2)

	for n > 1 {
		if n%k == 0 {
			res := &primepb.PrimeResponse{
				Prime: k,
			}
			err := stream.Send(res)
			if err != nil {
				log.Fatalf("Unable to send gRPC response: %v", err)
			}
			n = n / k
		} else {
			k = k + 1
		}
	}

	return nil
}

func main() {
	l, err := net.Listen("tcp", "0.0.0.0:50051")
	if err != nil {
		log.Fatalf("Failed to listen:%v", err)
	}
	s := grpc.NewServer()
	primepb.RegisterPrimeServiceServer(s, &server{})
	log.Println("Server is running on port:50051")
	if err = s.Serve(l); err != nil {
		log.Fatalf("failed to serve:%v", err)
	}
}