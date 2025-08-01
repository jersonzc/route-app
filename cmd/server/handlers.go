package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"google.golang.org/protobuf/proto"
	"io"
	"math"
	"os"
	"route-app/internal/route"
	"time"
)

func (app *application) GetFeature(ctx context.Context, point *route.Point) (*route.Feature, error) {
	for _, feature := range app.savedFeatures {
		if proto.Equal(feature.Location, point) {
			return feature, nil
		}
	}
	return &route.Feature{Location: point}, nil
}

func (app *application) ListFeatures(rect *route.Rectangle, stream route.Route_ListFeaturesServer) error {
	for _, feature := range app.savedFeatures {
		if inRange(feature.Location, rect) {
			if err := stream.Send(feature); err != nil {
				return err
			}
		}
	}
	return nil
}

func (app *application) RecordRoute(stream route.Route_RecordRouteServer) error {
	var pointCount, featureCount, distance int32
	var lastPoint *route.Point
	startTime := time.Now()
	for {
		point, err := stream.Recv()
		if err == io.EOF {
			endTime := time.Now()
			return stream.SendAndClose(&route.RouteSummary{
				PointCount:   pointCount,
				FeatureCount: featureCount,
				Distance:     distance,
				ElapsedTime:  int32(endTime.Sub(startTime).Seconds()),
			})
		}
		if err != nil {
			return err
		}
		pointCount++
		for _, feature := range app.savedFeatures {
			if proto.Equal(feature.Location, point) {
				featureCount++
			}
		}
		if lastPoint != nil {
			distance += calcDistance(lastPoint, point)
		}
		lastPoint = point
	}
}

func (app *application) RouteChat(stream route.Route_RouteChatServer) error {
	for {
		in, err := stream.Recv()
		if err == io.EOF {
			return nil
		}
		if err != nil {
			return err
		}

		key := serialize(in.Location)

		app.mu.Lock()
		app.routeNotes[key] = append(app.routeNotes[key], in)
		rn := make([]*route.RouteNote, len(app.routeNotes[key]))
		copy(rn, app.routeNotes[key])
		app.mu.Unlock()

		for _, note := range rn {
			if err := stream.Send(note); err != nil {
				return err
			}
		}
	}
}

func (app *application) initialize(filePath string) error {
	if filePath == "" {
		return errors.New("file path is empty")
	}

	var data []byte
	data, err := os.ReadFile(filePath)
	if err != nil {
		return fmt.Errorf("could not read file: %s", err.Error())
	}

	if err := json.Unmarshal(data, &app.savedFeatures); err != nil {
		return fmt.Errorf("could not unmarshal route features: %s", err.Error())
	}

	return nil
}

func inRange(point *route.Point, rect *route.Rectangle) bool {
	left := math.Min(float64(rect.Lo.Longitude), float64(rect.Hi.Longitude))
	right := math.Max(float64(rect.Lo.Longitude), float64(rect.Hi.Longitude))
	top := math.Max(float64(rect.Lo.Latitude), float64(rect.Hi.Latitude))
	bottom := math.Min(float64(rect.Lo.Latitude), float64(rect.Hi.Latitude))

	if float64(point.Longitude) >= left &&
		float64(point.Longitude) <= right &&
		float64(point.Latitude) >= bottom &&
		float64(point.Latitude) <= top {
		return true
	}
	return false
}

func toRadians(num float64) float64 {
	return num * math.Pi / float64(180)
}

// calcDistance calculates the distance between two points using the "haversine" formula.
// The formula is based on http://mathforum.org/library/drmath/view/51879.html.
func calcDistance(p1 *route.Point, p2 *route.Point) int32 {
	const CordFactor float64 = 1e7
	const R = float64(6371000) // earth radius in metres
	lat1 := toRadians(float64(p1.Latitude) / CordFactor)
	lat2 := toRadians(float64(p2.Latitude) / CordFactor)
	lng1 := toRadians(float64(p1.Longitude) / CordFactor)
	lng2 := toRadians(float64(p2.Longitude) / CordFactor)
	dlat := lat2 - lat1
	dlng := lng2 - lng1

	a := math.Sin(dlat/2)*math.Sin(dlat/2) +
		math.Cos(lat1)*math.Cos(lat2)*
			math.Sin(dlng/2)*math.Sin(dlng/2)
	c := 2 * math.Atan2(math.Sqrt(a), math.Sqrt(1-a))

	distance := R * c
	return int32(distance)
}

func serialize(point *route.Point) string {
	return fmt.Sprintf("%d %d", point.Latitude, point.Longitude)
}
