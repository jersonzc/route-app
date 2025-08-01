package main

import (
	"context"
	"flag"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"io"
	"log"
	"route-app/internal/route"
	"time"
)

// printFeature gets the feature for the given point.
func printFeature(client route.RouteClient, point *route.Point) {
	log.Printf("Getting feature for point (%d, %d)", point.Latitude, point.Longitude)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	feature, err := client.GetFeature(ctx, point)
	if err != nil {
		log.Fatalf("client.GetFeature failed: %v", err)
	}

	log.Println(feature)
}

// printFeatures lists all the features within the given bounding Rectangle.
func printFeatures(client route.RouteClient, rect *route.Rectangle) {
	log.Printf("Looking for features within %v", rect)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	stream, err := client.ListFeatures(ctx, rect)
	if err != nil {
		log.Fatalf("client.ListFeatures failed: %v", err)
	}

	for {
		feature, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatalf("client.ListFeatures failed: %v", err)
		}

		log.Printf("Feature: name: %q, point:(%v, %v)", feature.GetName(), feature.GetLocation().GetLatitude(), feature.GetLocation().GetLongitude())
	}
}

func main() {
	addr := flag.String("addr", "localhost:50051", "The server address in the format of host:port")
	flag.Parse()

	// Connect to server.
	conn, err := grpc.NewClient(*addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()

	// Contact the server.
	srv := route.NewRouteClient(conn)

	// Looking for a valid feature
	printFeature(srv, &route.Point{Latitude: 409146138, Longitude: -746188906})

	// Feature missing.
	printFeature(srv, &route.Point{Latitude: 0, Longitude: 0})

	// Looking for features between 40, -75 and 42, -73.
	printFeatures(srv, &route.Rectangle{
		Lo: &route.Point{Latitude: 400000000, Longitude: -750000000},
		Hi: &route.Point{Latitude: 420000000, Longitude: -730000000},
	})
}
