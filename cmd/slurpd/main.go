// package slurpd defines a daemon that will scrape news servers and collect articles
// from subscribed groups
package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net"

	"google.golang.org/grpc"

	pb "github.com/dominichamon/slurp/internal/api/slurp"
)

var (
	port   = flag.Int("port", 3232, "The port on which to listen for RPC requests")
	config = flag.String("config", "slurpd.json", "The JSON configuration file")
)

func main() {
	log.Println("loading config")
	data, err := ioutil.ReadFile(*config)
	if err != nil {
		log.Fatalf("Failed to read file: %s", err)
	}

	c, err := LoadConfig(string(data))
	if err != nil {
		log.Fatalf("failed to parse config: %s", err)
	}

	log.Println("starting nntp")
	nntp, err := NewNNTP(c)
	if err != nil {
		log.Fatal(err)
	}

	defer func() {
		if err := nntp.Close(); err != nil {
			log.Fatal(err)
		}
	}()

	l, err := net.Listen("tcp", fmt.Sprintf(":%d", *port))
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("starting slurpd on port %d", *port)
	s := grpc.NewServer()
	pb.RegisterSlurpServer(s, &slurpdServer{n: nntp})
	log.Printf("listening on port %d", *port)
	s.Serve(l)
}
