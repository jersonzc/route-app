package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"google.golang.org/grpc"
	"log"
	"net"
	"os"
	"os/signal"
	"route-app/internal/route"
	"time"
)

type application struct {
	route.UnimplementedRouteServer
	infoLog  *log.Logger
	errorLog *log.Logger
}

func main() {
	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	errorLog := log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)

	app := &application{
		infoLog:  infoLog,
		errorLog: errorLog,
	}

	if err := run(app); err != nil {
		errorLog.Fatalln(err)
	}
}

func run(app *application) (err error) {
	// Handle SIGINT (CTRL+C) gracefully.
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()

	// Read vars.
	port := flag.Int("port", 50051, "The server port")
	flag.Parse()

	// Start listener.
	listener, err := net.Listen("tcp", fmt.Sprintf("localhost:%d", *port))
	if err != nil {
		return
	}

	// Start server.
	app.infoLog.Println("Starting server...")
	
	srv := grpc.NewServer()
	route.RegisterRouteServer(srv, app)
	srvErr := make(chan error, 1)
	go func() {
		if myErr := srv.Serve(listener); myErr != nil {
			if !errors.Is(myErr, grpc.ErrServerStopped) {
				srvErr <- myErr
			}
		}
		close(srvErr)
	}()

	// Wait for interruption.
	select {
	case err = <-srvErr:
		return
	case <-ctx.Done():
		app.infoLog.Println("Shutting down...")

		shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		done := make(chan struct{})
		go func() {
			srv.GracefulStop()
			close(done)
		}()

		select {
		case <-done:
			app.infoLog.Println("Server stopped gracefully.")
		case <-shutdownCtx.Done():
			app.infoLog.Println("Forcing stop.")
			srv.Stop()
		}
	}

	return
}
