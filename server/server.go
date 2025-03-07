package main

import (
	"context"
	"log"
	"net"
	"sync"

	pb "HomeRepUserLocationService/proto"
	"google.golang.org/grpc"
)

type locationServer struct {
	pb.UnimplementedLocationServiceServer
	mu            sync.Mutex
	userLocations map[int64]*pb.GetLocationResponse // Используем GetLocationResponse вместо GeoPair
}

func (s *locationServer) UpdateLocation(ctx context.Context, req *pb.UpdateLocationRequest) (*pb.UpdateLocationResponse, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.userLocations[req.UserId] = &pb.GetLocationResponse{
		Lat: req.Lat,
		Lng: req.Lng,
	}

	return &pb.UpdateLocationResponse{Success: true}, nil
}

func (s *locationServer) GetLocation(ctx context.Context, req *pb.GetLocationRequest) (*pb.GetLocationResponse, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	location, exists := s.userLocations[req.UserId]
	if !exists {
		return nil, grpc.Errorf(grpc.Code(nil), "location not found")
	}

	return &pb.GetLocationResponse{
		Lat: location.Lat,
		Lng: location.Lng,
	}, nil
}

func main() {
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	grpcServer := grpc.NewServer()
	pb.RegisterLocationServiceServer(grpcServer, &locationServer{
		userLocations: make(map[int64]*pb.GetLocationResponse),
	})

	log.Println("Location Service is running on port 50051...")
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
