package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/muratbekauyez/go-final/primenumber/primepb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func main() {
	log.Print("Starting client")

	if len(os.Args) != 2 {
		log.Fatal("Missing number to decompose")
	}

	con, err := grpc.Dial("0.0.0.0:50051", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("could not connect: %v", err)
	}
	defer con.Close()

	c := primepb.NewPrimeServiceClient(con)

	i, err := strconv.Atoi(os.Args[1])
	if err != nil {
		log.Fatalf("Unable to parse argument number")
	}

	req := &primepb.PrimeRequest{
		Number: int64(i),
	}

	ctx, cancel := context.WithTimeout(context.Background(), 1000*time.Millisecond)
	defer cancel()

	stream, err := c.Prime(ctx, req)
	if err != nil {
		log.Fatalf("Error while requesting gRPC server: %v", err)
	}

	for {
		res, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			respErr, ok := status.FromError(err)
			if ok {
				if respErr.Code() == codes.InvalidArgument {
					fmt.Println("Number must be greater than 1")
				} else if respErr.Code() == codes.DeadlineExceeded {
					fmt.Println("Timeout")
				} else {
					log.Printf("Unknown error: %v", respErr.Message())
				}
			} else {
				log.Fatalf("Error with gRPC response: %v", err)
			}
			break
		}

		fmt.Println(res.GetPrime())
	}
}
