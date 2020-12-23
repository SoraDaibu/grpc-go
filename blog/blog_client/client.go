package main

import (
	"context"
	"fmt"
	"io"
	"log"

	"github.com/SoraDaibu/grpc-go-course/blog/blogpb"
	"google.golang.org/grpc"
)

func main() {
	fmt.Println("Blog Client")

	opts := grpc.WithInsecure()
	cc, err := grpc.Dial("localhost:50051", opts)
	if err != nil {
		log.Fatalf("could not connect %v\n", err)
	}
	defer cc.Close()

	c := blogpb.NewBlogServiceClient(cc)

	// create blog
	fmt.Println("Creating the blog")
	blog := &blogpb.Blog{
		AuthorId: "Sora",
		Title:    "How to become an engineer",
		Content:  "Start with prog8 and create your own portfolio!",
	}
	createBlogRes, err := c.CreateBlog(context.Background(), &blogpb.CreateBlogRequest{Blog: blog})
	if err != nil {
		log.Fatalf("Failed to create blog with unexpected error: %v\n", err)
	}
	fmt.Printf("Blog has been created: %v\n", createBlogRes)
	blogId := createBlogRes.GetBlog().GetId()
	// read Blog
	fmt.Println("Reading the blog")

	_, readErr := c.ReadBlog(context.Background(), &blogpb.ReadBlogRequest{BlogId: "5fe1b76c4b14da4022e3c15f"})

	if readErr != nil {
		fmt.Printf("Error while reading: %v\n", readErr)
	}

	readBlogReq := &blogpb.ReadBlogRequest{BlogId: blogId}
	readBlogRes, readBlogErr := c.ReadBlog(context.Background(), readBlogReq)
	if readBlogErr != nil {
		fmt.Printf("Error while reading: %v\n", readBlogErr)
	}

	fmt.Printf("Blog was read fine: %v\n", readBlogRes)

	// update Blog
	fmt.Println("Updating the blog")
	newBlog := &blogpb.Blog{
		Id:       blogId,
		AuthorId: "Changed Author",
		Title:    "(Editted) How to become an engineer in Spain",
		Content:  "Start con prog8 y crear tu portafolio!",
	}
	updateRes, updateErr := c.UpdateBlog(context.Background(), &blogpb.UpdateBlogRequest{Blog: newBlog})
	if updateErr != nil {
		fmt.Printf("Error while updating: %v\n", updateErr)
	}
	fmt.Printf("Blog was updated: %v\n", updateRes)

	// delete a blog
	fmt.Println("Deleting the blog")
	deleteRes, deleteErr := c.DeleteBlog(context.Background(), &blogpb.DeleteBlogRequest{BlogId: blogId})
	if deleteErr != nil {
		fmt.Printf("Error while deleting: %v\n", deleteErr)
	}
	fmt.Printf("Blog was deleted: %v\n", deleteRes)

	// list Blogs

	stream, err := c.ListBlog(context.Background(), &blogpb.ListBlogRequest{})

	if err != nil {
		log.Fatalf("error while calling ListBlog PRC: %v", err)
	}
	for {
		res, err := stream.Recv()
		if err == io.EOF {
			// we've reached the end of the stream
			break
		}
		if err != nil {
			log.Fatalf("error while reading stream: %v", err)
		}
		log.Println(res.GetBlog())
	}
}
