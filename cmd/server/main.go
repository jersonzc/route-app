package main

import (
	"context"
	"route-app/internal/route"
)

type routeServer struct {
	route.UnimplementedRouteServer
}

func (s *routeServer) GetFeature(ctx context.Context, point *route.Point) (*route.Feature, error) {
	return nil, nil
}

func (s *routeServer) ListFeatures(rect *route.Rectangle, stream route.Route_ListFeaturesServer) error {
	return nil
}

func (s *routeServer) RecordRoute(stream route.Route_RecordRouteServer) error {
	return nil
}

func (s *routeServer) RouteChat(stream route.Route_RouteChatServer) error {
	return nil
}
