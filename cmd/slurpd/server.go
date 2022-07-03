package main

import (
	"context"
	"log"

	pb "github.com/dominichamon/slurp/internal/api/slurp"
)

type slurpdServer struct {
	pb.UnimplementedSlurpServer

	n *NNTP
}

func (s *slurpdServer) Addresses(_ context.Context, _ *pb.AddressesRequest) (*pb.AddressesResponse, error) {
	log.Printf("getting addresses")
	addrs, err := s.n.Addresses()
	if err != nil {
		return nil, err
	}

	resp := &pb.AddressesResponse{
		Address: make([]string, len(addrs)),
	}

	for i, a := range addrs {
		resp.Address[i] = a
	}
	return resp, nil
}

func (s *slurpdServer) Groups(_ context.Context, r *pb.GroupsRequest) (*pb.GroupsResponse, error) {
	log.Printf("getting groups for server %s", r.Server)
	groups, err := s.n.Groups(r.Server)
	if err != nil {
		return nil, err
	}

	resp := &pb.GroupsResponse{
		Group: make([]string, len(groups)),
	}
	for i, g := range groups {
		resp.Group[i] = g.Name
	}
	return resp, nil
}
