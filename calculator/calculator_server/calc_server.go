package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"math"
	"net"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

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
	fmt.Printf("DecomposeManyTimes func was incoked with %v\n", req)
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

func (*server) ComputeAverage(stream calculatorpb.CalculatorService_ComputeAverageServer) error {
	fmt.Printf("ComputeAverage func was invoked with a client streaming req\n")
	//the returned value of req.GetNumber() is int32, therefore setting avenum int32 for now and setting float64 later inside of if(err == io.EOF)
	var ave_num int32
	var count float64
	for {
		req, err := stream.Recv()
		if err == io.EOF {
			// we have finished reading the client stream
			return stream.SendAndClose(&calculatorpb.ComputeAverageResponse{
				AveNum: float64(ave_num) / count,
			})
		}
		if err != nil {
			log.Fatal("Error while reading client stream: %v\n", err)
		}

		// Normal Process starts here
		count++
		ave_num += req.GetNumber()
	}
}

func (*server) FindMaximum(stream calculatorpb.CalculatorService_FindMaximumServer) error {
	fmt.Printf("FindMaximum func was invoked with a BiDi streaming req\n")
	var max int32
	for {
		req, err := stream.Recv()
		if err == io.EOF {
			return nil
		}
		if err != nil {
			log.Fatalf("Error while reading client stream: %v\n", err)
			return err
		}
		num := req.GetNumber()
		if num > max {
			max = num
			sendErr := stream.Send(&calculatorpb.FindMaximumResponse{
				CurrentMaxNum: max,
			})
			if sendErr != nil {
				log.Fatal("Error while sending maximum num data to client: %v\n", sendErr)
				return sendErr
			}
		}
	}
}

func (*server) SquareRoot(ctx context.Context, req *calculatorpb.SquareRootRequest) (*calculatorpb.SquareRootResponse, error) {
	fmt.Printf("SquareRoot func was invoked with an Unary req\n")
	number := req.GetNumber()
	if number < 0 {
		return nil, status.Errorf(
			codes.InvalidArgument,
			fmt.Sprintf("Received a negative number: %v", number),
		)
	}
	return &calculatorpb.SquareRootResponse{
		NumberRoot: math.Sqrt(float64(number)),
	}, nil
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
