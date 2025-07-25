package main

import (
	"context"
	"route-app/internal/route"
)

func (s *application) GetFeature(ctx context.Context, point *route.Point) (*route.Feature, error) {
	s.infoLog.Printf("received: %v", point.String())
	return &route.Feature{Name: "test"}, nil
}

func (s *application) ListFeatures(rect *route.Rectangle, stream route.Route_ListFeaturesServer) error {
	return nil
}

func (s *application) RecordRoute(stream route.Route_RecordRouteServer) error {
	return nil
}

func (s *application) RouteChat(stream route.Route_RouteChatServer) error {
	return nil
}
