package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"google.golang.org/grpc"

	pb "github.com/daodao97/egin/rpc/pb"
)

func main() {

	r := gin.Default()

	r.GET("/rpc/hello", func(c *gin.Context) {
		sayHello(c)
	})

	// Run http server
	if err := r.Run(":8052"); err != nil {
		log.Fatalf("could not run server: %v", err)
	}
}

func sayHello(c *gin.Context) {

	// Set up a connection to the server.
	conn, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}

	defer conn.Close()

	client := pb.NewHelloClient(conn)
	name := c.DefaultQuery("name", "战士上战场")
	req := &pb.HelloRequest{Name: name}
	res, err := client.SayHello(c, req)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"result": fmt.Sprint(res.Message),
		"code":   fmt.Sprint(res.Code),
	})

}
