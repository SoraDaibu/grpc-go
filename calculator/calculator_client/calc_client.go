package main

import (
	"context"
	"fmt"
	"io"
	"log"

	"github.com/SoraDaibu/grpc-go-course/calculator/calculatorpb"

	"google.golang.org/grpc"
)

func main() {
	fmt.Printf("Hello I am a calc_client!!\n")
	cc, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("could not connect %v", err)
	}
	defer cc.Close()

	c := calculatorpb.NewCalculatorServiceClient(cc)
	// doUnary(c)

	doServerStreaming(c)
}

func doUnary(c calculatorpb.CalculatorServiceClient) {
	fmt.Println("Starting to do an Unary RPC ...")
	req := &calculatorpb.CalculateRequest{
		Calculator: &calculatorpb.Calculator{
			FirstNum:  10,
			SecondNum: 3,
		},
	}
	res, err := c.Calculate(context.Background(), req)
	if err != nil {
		log.Fatalf("error while calling Calculate: %v", req)
	}
	log.Printf("Response from Calculator: %v", res.Result)
}

func doServerStreaming(c calculatorpb.CalculatorServiceClient) {
	fmt.Println("Starting to do an Decompose Server Streaming RPC ...")

	req := &calculatorpb.DecomposeRequest{
		Number: 1204567,
	}
	resStream, err := c.DecomposeManyTimes(context.Background(), req)
	if err != nil {
		log.Fatalf("error while calling DecomposeManyTimes PRC: %v", err)
	}
	for {
		msg, err := resStream.Recv()
		if err == io.EOF {
			// we've reached the end of the stream
			break
		}
		if err != nil {
			log.Fatalf("error while reading stream: %v", err)
		}
		log.Printf("Response from DecomposeManyTimes: %v", msg.GetDecomposedNum())
	}
}
