package client

import (
	"context"
	"fmt"
	"log"
	"google.golang.org/grpc"
	computeAverage "../pb"
	//"io"
	"time"
)

func main() {
	fmt.Println("Hello I'm the client")
	cc, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("could not connect: %v", err)
	}

	defer cc.Close()

	c:= computeAverage.NewComputeAverageServiceClient(cc)
	//doUnary(c)
	arr := []int32{1,2,3,4}
	GetAverageOfArr(arr,c)


}

func GetAverageOfArr(arr []int32, c computeAverage.ComputeAverageServiceClient ) {
	stream, err := c.ComputeAvg(context.Background())
	if err != nil {
		log.Fatalf("error while calling LongGreet RPC: %v", err)
	}
	for _, val := range arr {
		req := &computeAverage.ComputeAverageRequest {
			Val : val,
		}
		fmt.Println("Sending req %v", req)
		stream.Send(req)
		time.Sleep(100 * time.Millisecond)
	}

	res, err := stream.CloseAndRecv()
	if err != nil {
		log.Fatalf("error while recieving response from computeAvg &v", err)
	}
	fmt.Printf("computeAvg Response: %v\n",res)
}