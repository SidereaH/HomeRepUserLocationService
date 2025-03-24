package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"
	"sync"
	"time"

	pb "HomeRepUserLocationService/proto"
	"github.com/jackc/pgx/v5/pgxpool"
	"google.golang.org/grpc"
)

type locationServer struct {
	pb.UnimplementedLocationServiceServer
	mu sync.Mutex
	db *pgxpool.Pool
}

func NewLocationServer(db *pgxpool.Pool) *locationServer {
	return &locationServer{db: db}
}

func (s *locationServer) UpdateLocation(ctx context.Context, req *pb.UpdateLocationRequest) (*pb.UpdateLocationResponse, error) {
	_, err := s.db.Exec(ctx,
		"INSERT INTO user_locations (user_id, latitude, longitude) VALUES ($1, $2, $3)",
		req.UserId, req.Location.Lat, req.Location.Lng,
	)
	if err != nil {
		log.Printf("Failed to update location: %v", err)
		return &pb.UpdateLocationResponse{Success: false}, err
	}
	log.Printf("success adding : %v, %v", req.Location.Lat, req.Location.Lng)
	return &pb.UpdateLocationResponse{Success: true}, nil
}

func (s *locationServer) GetLocation(ctx context.Context, req *pb.GetLocationRequest) (*pb.GetLocationResponse, error) {
	var lat, lng float64
	err := s.db.QueryRow(ctx,
		"SELECT latitude, longitude FROM user_locations WHERE user_id = $1 ORDER BY time DESC LIMIT 1",
		req.UserId,
	).Scan(&lat, &lng)
	if err != nil {
		log.Printf("Failed to get location: %v", err)
		return nil, err
	}
	return &pb.GetLocationResponse{Location: &pb.GeoPair{Lat: lat, Lng: lng}}, nil
}

func (s *locationServer) GetLocationHistory(ctx context.Context, req *pb.GetLocationHistoryRequest) (*pb.GetLocationHistoryResponse, error) {
	rows, err := s.db.Query(ctx,
		"SELECT latitude, longitude, time FROM user_locations WHERE user_id = $1 AND time BETWEEN $2 AND $3 ORDER BY time",
		req.UserId, req.StartTime, req.EndTime,
	)
	if err != nil {
		log.Printf("Failed to get location history: %v", err)
		return nil, err
	}
	defer rows.Close()

	var locations []*pb.GeoPair
	var timestamps []string
	for rows.Next() {
		var lat, lng float64
		var timestamp time.Time
		if err := rows.Scan(&lat, &lng, &timestamp); err != nil {
			log.Printf("Failed to scan row: %v", err)
			continue
		}
		locations = append(locations, &pb.GeoPair{Lat: lat, Lng: lng})
		timestamps = append(timestamps, timestamp.Format(time.RFC3339))
	}

	return &pb.GetLocationHistoryResponse{Locations: locations, Timestamps: timestamps}, nil
}

func main() {
	// Подключение к TimescaleDB
	dbHost := os.Getenv("DB_HOST")
	dbUser := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbName := os.Getenv("DB_NAME")
	log.Printf("postgres://%s:%s@%s:5432/%s", dbUser, dbPassword, dbHost, dbName)
	//connStr := fmt.Sprintf("postgres://%s:%s@%s:5432/%s", dbUser, dbPassword, dbHost, dbName)
	connStr := fmt.Sprint("postgres://user:password@timescaledb:5432/location_service")
	db, err := pgxpool.New(context.Background(), connStr)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	_, err = db.Exec(context.Background(), `
        CREATE TABLE IF NOT EXISTS user_locations (
            user_id BIGINT NOT NULL,
            latitude DOUBLE PRECISION NOT NULL,
            longitude DOUBLE PRECISION NOT NULL,
            time TIMESTAMPTZ DEFAULT NOW()
        );
        
        SELECT create_hypertable('user_locations', 'time', if_not_exists => TRUE);
    `)
	if err != nil {
		log.Fatalf("Failed to create table: %v", err)
	}
	defer db.Close()

	// Запуск gRPC-сервера
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	grpcServer := grpc.NewServer()
	pb.RegisterLocationServiceServer(grpcServer, NewLocationServer(db))

	log.Println("Location Service is running on port 50051...")
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}
