// nolint:gosec
package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/joho/godotenv"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	inventoryV1API "github.com/AxMdv/go-rocket-factory/inventory/internal/api/inventory/v1"
	partRepository "github.com/AxMdv/go-rocket-factory/inventory/internal/repository/part"
	partService "github.com/AxMdv/go-rocket-factory/inventory/internal/service/part"
	inventoryV1 "github.com/AxMdv/go-rocket-factory/shared/pkg/proto/inventory/v1"
)

const (
	grpcPort        = 50051
	startupTimeout  = 10 * time.Second
	shutdownTimeout = 10 * time.Second
)

type config struct {
	mongoURI             string
	mongoDatabase        string
	mongoPartsCollection string
}

func main() {
	cfg, err := loadConfig()
	if err != nil {
		log.Printf("failed to load config: %v\n", err)
		return
	}

	startupCtx, cancelStartup := context.WithTimeout(context.Background(), startupTimeout)
	repo, err := partRepository.NewRepository(
		startupCtx,
		cfg.mongoURI,
		cfg.mongoDatabase,
		cfg.mongoPartsCollection,
	)
	cancelStartup()
	if err != nil {
		log.Printf("failed to create part repository: %v\n", err)
		return
	}
	defer func() {
		closeCtx, cancelClose := context.WithTimeout(context.Background(), shutdownTimeout)
		defer cancelClose()

		if closeErr := repo.Close(closeCtx); closeErr != nil {
			log.Printf("failed to close part repository: %v\n", closeErr)
		}
	}()

	lis, err := net.Listen("tcp", fmt.Sprintf("localhost:%d", grpcPort))
	if err != nil {
		log.Printf("failed to listen: %v\n", err)
		return
	}

	server := grpc.NewServer()
	service := partService.NewService(repo)
	api := inventoryV1API.NewAPI(service)

	inventoryV1.RegisterInventoryServiceServer(server, api)
	reflection.Register(server)

	go func() {
		log.Printf("inventory gRPC server listening on %d\n", grpcPort)
		if serveErr := server.Serve(lis); serveErr != nil {
			log.Printf("failed to serve: %v\n", serveErr)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	signal.Stop(quit)

	log.Println("shutting down inventory gRPC server...")
	server.GracefulStop()
	log.Println("inventory gRPC server stopped")
}

func loadConfig() (config, error) {
	if err := godotenv.Load(".env"); err != nil {
		log.Printf(".env file was not loaded: %v; using process environment\n", err)
	}

	mongoURI, err := requiredEnv("MONGO_URI")
	if err != nil {
		return config{}, err
	}
	mongoDatabase, err := requiredEnv("MONGO_DATABASE")
	if err != nil {
		return config{}, err
	}
	mongoPartsCollection, err := requiredEnv("MONGO_PARTS_COLLECTION")
	if err != nil {
		return config{}, err
	}

	return config{
		mongoURI:             mongoURI,
		mongoDatabase:        mongoDatabase,
		mongoPartsCollection: mongoPartsCollection,
	}, nil
}

func requiredEnv(name string) (string, error) {
	value := strings.TrimSpace(os.Getenv(name))
	if value == "" {
		return "", fmt.Errorf("environment variable %s is empty", name)
	}

	return value, nil
}
