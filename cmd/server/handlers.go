package main

import (
	"context"
	"route-app/internal/route"
)

func (app *application) GetFeature(ctx context.Context, point *route.Point) (*route.Feature, error) {
	app.infoLog.Printf("received: %v", point.String())
	return &route.Feature{Name: "test"}, nil
}

func (app *application) ListFeatures(rect *route.Rectangle, stream route.Route_ListFeaturesServer) error {
	return nil
}

func (app *application) RecordRoute(stream route.Route_RecordRouteServer) error {
	return nil
}

func (app *application) RouteChat(stream route.Route_RouteChatServer) error {
	return nil
}
