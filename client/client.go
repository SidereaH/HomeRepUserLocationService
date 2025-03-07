package main

import (
	"context"
	"log"
	"time"

	pb "HomeRepUserLocationService/proto"
	"google.golang.org/grpc"
)

func main() {
	conn, err := grpc.Dial("localhost:50051", grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()

	client := pb.NewLocationServiceClient(conn)

	// Обновление геолокации
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	updateResponse, err := client.UpdateLocation(ctx, &pb.UpdateLocationRequest{
		UserId: 1,
		Lat:    55.7558,
		Lng:    37.6176,
	})
	if err != nil {
		log.Fatalf("could not update location: %v", err)
	}
	log.Printf("Update Location Response: %v", updateResponse)

	// Получение геолокации
	getResponse, err := client.GetLocation(ctx, &pb.GetLocationRequest{
		UserId: 1,
	})
	if err != nil {
		log.Fatalf("could not get location: %v", err)
	}
	log.Printf("Get Location Response: %v", getResponse)
}
