package main

import (
	"context"
	"fmt"
	"log"

	"github.com/SoraDaibu/grpc-go-course/greet/greetpb"

	"google.golang.org/grpc"
)

func main() {
	fmt.Printf("Hello I'm a client\n")
	cc, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("could not connect %v", err)
	}
	defer cc.Close()

	c := greetpb.NewGreetServiceClient(cc)
	// fmt.Printf("Created client: %f", c)
	doUnary(c)
}

func doUnary(c greetpb.GreetServiceClient) {
	fmt.Println("Starting to do an Unary RPC ...")
	req := &greetpb.GreetRequest{
		Greeting: &greetpb.Greeting{
			FirstName: "Sora",
			LastName:  "Daibu",
		},
	}
	res, err := c.Greet(context.Background(), req)
	if err != nil {
		log.Fatalf("error while ccalling Greet RPC: %v", req)
	}
	log.Printf("Response from Greet: %v", res.Result)
}
