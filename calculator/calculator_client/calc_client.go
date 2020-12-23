package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"time"

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
	// doServerStreaming(c)
	// doClientStreaming(c)
	doBiDiStreaming(c)
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

func doClientStreaming(c calculatorpb.CalculatorServiceClient) {
	fmt.Println("Starting to do an Client Streaming RPC ...")

	requests := []int32{140, 120, 120, 120, 120, 120, 120}

	stream, err := c.ComputeAverage(context.Background())
	if err != nil {
		log.Fatalf("error while calling ComputeAverage:%v", err)
	}

	for _, req := range requests {
		fmt.Printf("Sending req: %v\n", req)
		stream.Send(&calculatorpb.ComputeAverageRequest{
			Number: req,
		})
		time.Sleep(500 * time.Millisecond)
	}

	res, err := stream.CloseAndRecv()
	if err != nil {
		log.Fatalf("error while receiving response from ComputeAverage:%v", err)
	}
	fmt.Printf("ComputeAverage response: %v\n", res)
}

func doBiDiStreaming(c calculatorpb.CalculatorServiceClient) {
	fmt.Println("Starting to do an BiDi Streaming RPC ...")

	// we create a stream by invoking the client
	stream, err := c.FindMaximum(context.Background())
	if err != nil {
		log.Fatalf("Errror while creating stream: %v", err)
		return
	}

	waitc := make(chan struct{})

	// we send a bunch of integers to the server (go routine)
	go func() {
		// func to send a bunch of integers
		num := []int32{1, 14, 1200, 12, 1100, 1500, 3150, 1000}
		for _, req := range num {
			fmt.Printf("Sending num: %v\n", req)
			stream.Send(&calculatorpb.FindMaximumRequest{
				Number: req,
			})
			time.Sleep(500 * time.Millisecond)
		}
	}()

	// we receive a bunch of ints from the server (go routine)
	go func() {
		// func to receive a bunch of ints
		for {
			res, err := stream.Recv()
			if err == io.EOF {
				break
			}
			if err != nil {
				log.Fatalf("Error while receiving: %v", err)
				break
			}
			fmt.Printf("Received: %v !!!!!!!!\n", res.GetCurrentMaxNum())
		}
		close(waitc)
	}()

	// block until everything is done
	fmt.Printf("waitc is: %v\n", waitc)
	<-waitc

}
