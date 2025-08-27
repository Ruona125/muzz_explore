package main

import (
	"context"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"

	"google.golang.org/grpc"

	"explore_service/internal/server"
	"explore_service/internal/storage"
	explorepb "explore_service/proto"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
)

// getEnv fetches an environment variable or returns the fallback if unset.
func getEnv(key, fallback string) string {
	if v, ok := os.LookupEnv(key); ok {
		return v
	}
	return fallback
}

func main() {
	ctx := context.Background()
	// Load environment variables from .env if present (local development convenience)
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found; continuing with existing environment variables")
	}
	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		log.Fatal("environment variable DATABASE_URL must be set to a PostgreSQL DSN")
	}
	// Connect to Postgres using pgxpool.  Use background context for
	// connection creation but create a derived context for migration
	// calls where necessary.
	pool, err := pgxpool.New(ctx, dsn)
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}
	defer pool.Close()
	// Initialise the store and run migrations.
	store, err := storage.NewStore(ctx, pool)
	if err != nil {
		log.Fatalf("database migration failed: %v", err)
	}
	// Create the gRPC server and register our ExploreService.
	grpcServer := grpc.NewServer()
	svc := server.NewExploreServer(store, 50)
	explorepb.RegisterExploreServiceServer(grpcServer, svc)
	// Listen on the port specified by the PORT environment variable or
	// default to 50051.  In Docker environments this
	// variable can be set via configuration.
	addr := ":" + getEnv("PORT", "50051")
	lis, err := net.Listen("tcp", addr)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	log.Printf("ExploreService listening on %s", addr)
	// Run the server in a goroutine so that we can handle graceful
	// shutdown via OS signals.
	go func() {
		if err := grpcServer.Serve(lis); err != nil {
			log.Fatalf("gRPC server exited with error: %v", err)
		}
	}()
	// Block until we receive an interrupt or termination signal.
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)
	<-stop
	log.Println("Shutting down ExploreService...")
	grpcServer.GracefulStop()
}
