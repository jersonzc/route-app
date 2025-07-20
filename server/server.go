package main

import (
	"context"
	pb "example.com/route"
)

type routeServer struct {
	pb.UnimplementedRouteServer
}

func (s *routeServer) GetFeature(ctx context.Context, point *pb.Point) (*pb.Feature, error) {
	return nil, nil
}

func (s *routeServer) ListFeatures(rect *pb.Rectangle, stream pb.Route_ListFeaturesServer) error {
	return nil
}

func (s *routeServer) RecordRoute(stream pb.Route_RecordRouteServer) error {
	return nil
}

func (s *routeServer) RouteChat(stream pb.Route_RouteChatServer) error {
	return nil
}
