package main

import (
	"context"
	"fmt"
	"log"
	"net"

	"github.com/SoraDaibu/grpc-go-course/calculator/calculatorpb"

	"google.golang.org/grpc"
)

type server struct {
}

func (*server) Calculate(ctx context.Context, req *calculatorpb.CalculateRequest) (*calculatorpb.CalculateResponse, error) {
	fmt.Printf("Calculate func in calc_server.go was invoked with %v\n", req)
	firstNum := req.GetCalculator().GetFirstNum()
	secondNum := req.GetCalculator().GetSecondNum()
	result := firstNum + secondNum
	res := &calculatorpb.CalculateResponse{
		Result: result,
	}
	return res, nil
}

func (*server) DecomposeManyTimes(req *calculatorpb.DecomposeRequest, stream calculatorpb.CalculatorService_DecomposeManyTimesServer) error {
	fmt.Printf("DecomposeManyTimes fuc was incoked with %v\n", req)
	N := req.GetNumber()
	var k int32 = 2
	for N > 1 {
		// if k evenly divides into N
		if N%k == 0 {
			fmt.Printf("n = %v\n", k)
			// divide N by k so that we have the rest of the number left.
			stream.Send(&calculatorpb.DecomposeResponse{
				DecomposedNum: k,
			})
			N = N / k
		} else {
			k++
			fmt.Printf("k: has increased to %v\n", k)
		}
	}
	return nil
}

func main() {
	fmt.Println("¡Hola Mundo!")

	// 0.0.0.0:500051 is: a default port for gRPC (I believe)
	lis, err := net.Listen("tcp", "0.0.0.0:50051")
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	s := grpc.NewServer()
	calculatorpb.RegisterCalculatorServiceServer(s, &server{})

	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
