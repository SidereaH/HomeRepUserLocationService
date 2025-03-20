package server

import (
	"context"
	"testing"

	pb "HomeRepUserLocationService/proto"
	"github.com/stretchr/testify/assert"
)

type mockLocationService struct {
	pb.UnimplementedLocationServiceServer
	locations map[int64]*pb.GetLocationResponse
}

func (m *mockLocationService) UpdateLocation(ctx context.Context, req *pb.UpdateLocationRequest) (*pb.UpdateLocationResponse, error) {
	m.locations[req.UserId] = &pb.GetLocationResponse{Location: &pb.GeoPair{Lat: req.GetLocation().Lat, Lng: req.GetLocation().Lng}}
	return &pb.UpdateLocationResponse{Success: true}, nil
}

func (m *mockLocationService) GetLocation(ctx context.Context, req *pb.GetLocationRequest) (*pb.GetLocationResponse, error) {
	loc, exists := m.locations[req.UserId]
	if !exists {
		return nil, nil
	}
	return loc, nil
}

func setupMockService() *mockLocationService {
	return &mockLocationService{locations: make(map[int64]*pb.GetLocationResponse)}
}

func TestUpdateLocation(t *testing.T) {
	service := setupMockService()
	resp, err := service.UpdateLocation(context.Background(), &pb.UpdateLocationRequest{UserId: 1, Location: &pb.GeoPair{Lat: 52.52, Lng: 13.405}})
	assert.NoError(t, err)
	assert.True(t, resp.Success)
}

func TestGetLocation(t *testing.T) {
	service := setupMockService()
	service.locations[1] = &pb.GetLocationResponse{Location: &pb.GeoPair{Lat: 52.52, Lng: 13.405}}
	resp, err := service.GetLocation(context.Background(), &pb.GetLocationRequest{UserId: 1})
	assert.NoError(t, err)
	assert.Equal(t, 52.52, resp.GetLocation().GetLat())
	assert.Equal(t, 13.405, resp.GetLocation().GetLng())
}

func TestGetLocationNotFound(t *testing.T) {
	service := setupMockService()
	resp, err := service.GetLocation(context.Background(), &pb.GetLocationRequest{UserId: 999})
	assert.NoError(t, err)
	assert.Nil(t, resp)
}

func TestUpdateMultipleUsers(t *testing.T) {
	service := setupMockService()
	service.UpdateLocation(context.Background(), &pb.UpdateLocationRequest{UserId: 1, Location: &pb.GeoPair{Lat: 40.7128, Lng: -74.0060}})
	service.UpdateLocation(context.Background(), &pb.UpdateLocationRequest{UserId: 2, Location: &pb.GeoPair{Lat: 34.0522, Lng: -118.2437}})

	resp1, _ := service.GetLocation(context.Background(), &pb.GetLocationRequest{UserId: 1})
	resp2, _ := service.GetLocation(context.Background(), &pb.GetLocationRequest{UserId: 2})

	assert.Equal(t, 40.7128, resp1.GetLocation().GetLat())
	assert.Equal(t, -74.0060, resp1.GetLocation().GetLng())
	assert.Equal(t, 34.0522, resp2.GetLocation().GetLat())
	assert.Equal(t, -118.2437, resp2.GetLocation().GetLng())
}

//func TestUpdateLocationOverride(t *testing.T) {
//	service := setupMockService()
//	service.UpdateLocation(context.Background(), &pb.UpdateLocationRequest{UserId: 1, Lat: 51.5074, Lng: -0.1278})
//	service.UpdateLocation(context.Background(), &pb.UpdateLocationRequest{UserId: 1, Lat: 48.8566, Lng: 2.3522})
//
//	resp, _ := service.GetLocation(context.Background(), &pb.GetLocationRequest{UserId: 1})
//	assert.Equal(t, 48.8566, resp.Lat)
//	assert.Equal(t, 2.3522, resp.Lng)
//}

//func TestInvalidCoordinates(t *testing.T) {
//	service := setupMockService()
//	resp, err := service.UpdateLocation(context.Background(), &pb.UpdateLocationRequest{UserId: 1, Lat: 999, Lng: 999})
//	assert.NoError(t, err)
//	assert.True(t, resp.Success)
//}
//
//func TestNegativeCoordinates(t *testing.T) {
//	service := setupMockService()
//	resp, err := service.UpdateLocation(context.Background(), &pb.UpdateLocationRequest{UserId: 1, Lat: -33.8651, Lng: 151.2099})
//	assert.NoError(t, err)
//	assert.True(t, resp.Success)
//}
//
//func TestZeroCoordinates(t *testing.T) {
//	service := setupMockService()
//	resp, err := service.UpdateLocation(context.Background(), &pb.UpdateLocationRequest{UserId: 1, Lat: 0, Lng: 0})
//	assert.NoError(t, err)
//	assert.True(t, resp.Success)
//}
//
//func TestConcurrentUpdates(t *testing.T) {
//	service := setupMockService()
//	go service.UpdateLocation(context.Background(), &pb.UpdateLocationRequest{UserId: 1, Lat: 10, Lng: 10})
//	go service.UpdateLocation(context.Background(), &pb.UpdateLocationRequest{UserId: 1, Lat: 20, Lng: 20})
//}

//func TestLargeUserId(t *testing.T) {
//	service := setupMockService()
//	resp, err := service.UpdateLocation(context.Background(), &pb.UpdateLocationRequest{UserId: 9223372036854775807, Lat: 37.7749, Lng: -122.4194})
//	assert.NoError(t, err)
//	assert.True(t, resp.Success)
//}
