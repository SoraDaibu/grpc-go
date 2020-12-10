package main

import (
	"fmt"
	"log"
	"net"

	"google.golang.org/grpc"
)

type server struct {
}

func main() {
	fmt.Println("¡Hola Mundo!")

	// 0.0.0.0:500051 is: a default port for gRPC (I believe)
	lis, err := net.Listen("tcp", "0.0.0.0:500051")
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	s := grpc.NewServer()
	calculatorpb.RegisterCaltulatorServiceServer(s, &server{})
}
