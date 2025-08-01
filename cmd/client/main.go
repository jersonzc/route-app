package main

import (
	"context"
	"flag"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"log"
	"route-app/internal/route"
	"time"
)

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
}
