// package slurp defines an HTTP server that displays newsgroups and articles.
package main

import (
	"flag"
	"fmt"
	"log"

	"google.golang.org/grpc"

	pb "github.com/dominichamon/slurp/internal/api/slurp"
)

var (
	dport = flag.Int("dport", 3232, "The port on which the slurp daemon is listening")
)

func main() {
	log.Println("loading config")
	// TODO: load local config with subscriptions

	log.Printf("dialling daemon on port %d", *dport)
	conn, err := grpc.Dial(fmt.Sprintf(":%d", *dport), grpc.WithInsecure())
	if err != nil {
		log.Fatal(err)
	}
	client := pb.NewSlurpClient(conn)

	log.Println("serving...")
	if err := serve(client); err != nil {
		log.Fatal(err)
	}
}
