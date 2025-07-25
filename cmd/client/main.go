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
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	res, err := srv.GetFeature(ctx, &route.Point{Latitude: 1, Longitude: 1})
	if err != nil {
		log.Fatalf("could not get feature: %v", err)
	}
	log.Printf("feature name: %s", res.Name)
}
