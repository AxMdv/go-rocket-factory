package main

import (
	"context"
	"database/sql"
	"errors"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/jackc/pgx/v5/pgxpool"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/joho/godotenv"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	orderV1API "github.com/AxMdv/go-rocket-factory/order/internal/api/order/v1"
	inventoryClientV1 "github.com/AxMdv/go-rocket-factory/order/internal/client/grpc/inventory/v1"
	paymentClientV1 "github.com/AxMdv/go-rocket-factory/order/internal/client/grpc/payment/v1"
	"github.com/AxMdv/go-rocket-factory/order/internal/migrator"
	orderRepository "github.com/AxMdv/go-rocket-factory/order/internal/repository/order"
	orderService "github.com/AxMdv/go-rocket-factory/order/internal/service/order"
	orderV1 "github.com/AxMdv/go-rocket-factory/shared/pkg/openapi/order/v1"
	inventoryV1 "github.com/AxMdv/go-rocket-factory/shared/pkg/proto/inventory/v1"
	paymentV1 "github.com/AxMdv/go-rocket-factory/shared/pkg/proto/payment/v1"
)

const (
	httpPort             = "8080"
	inventoryServiceAddr = "localhost:50051"
	paymentServiceAddr   = "localhost:50052"
	readHeaderTimeout    = 5 * time.Second
	shutdownTimeout      = 10 * time.Second
)

func main() {
	ctx := context.Background()

	if err := godotenv.Load(".env"); err != nil {
		log.Printf("failed to load .env file: %v\n", err)
		return
	}

	dbURI := os.Getenv("DB_URI")
	if dbURI == "" {
		log.Println("DB_URI is empty")
		return
	}

	migrationsDir := os.Getenv("MIGRATIONS_DIR")
	if migrationsDir == "" {
		migrationsDir = "migrations"
	}

	migrationDB, err := sql.Open("pgx", dbURI)
	if err != nil {
		log.Printf("failed to open migration db: %v\n", err)
		return
	}
	defer func() {
		if cerr := migrationDB.Close(); cerr != nil {
			log.Printf("failed to close migration db: %v", cerr)
		}
	}()

	if err = migrationDB.PingContext(ctx); err != nil {
		log.Printf("failed to ping migration db: %v\n", err)
		return
	}

	if err = migrator.NewMigrator(migrationDB, migrationsDir).Up(); err != nil {
		log.Printf("failed to apply migrations: %v\n", err)
		return
	}

	dbPool, err := pgxpool.New(ctx, dbURI)
	if err != nil {
		log.Printf("failed to create db pool: %v\n", err)
		return
	}
	defer dbPool.Close()

	if err = dbPool.Ping(ctx); err != nil {
		log.Printf("failed to ping db pool: %v\n", err)
		return
	}

	inventoryConn, err := grpc.NewClient(
		inventoryServiceAddr,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		log.Printf("failed to connect to inventory: %v\n", err)
		return
	}
	defer func() {
		if cerr := inventoryConn.Close(); cerr != nil {
			log.Printf("failed to close connect to inventory: %v", cerr)
		}
	}()

	paymentConn, err := grpc.NewClient(
		paymentServiceAddr,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		log.Printf("failed to connect to payment: %v\n", err)
		return
	}
	defer func() {
		if cerr := paymentConn.Close(); cerr != nil {
			log.Printf("failed to close connect to payment: %v", cerr)
		}
	}()

	ic := inventoryV1.NewInventoryServiceClient(inventoryConn)
	pc := paymentV1.NewPaymentServiceClient(paymentConn)

	invClient := inventoryClientV1.NewInventoryClient(ic)
	payClient := paymentClientV1.NewPaymentClient(pc)

	repo := orderRepository.NewRepository(dbPool)
	orderSvc := orderService.NewOrderService(repo, invClient, payClient)

	api := orderV1API.NewAPI(orderSvc)
	orderServer, err := orderV1.NewServer(api)
	if err != nil {
		log.Fatal(err)
	}

	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Timeout(10 * time.Second))
	r.Mount("/", orderServer)

	server := &http.Server{
		Addr:              net.JoinHostPort("localhost", httpPort),
		Handler:           r,
		ReadHeaderTimeout: readHeaderTimeout,
	}

	go func() {
		log.Printf("HTTP server started on port %s\n", httpPort)
		serveErr := server.ListenAndServe()
		if serveErr != nil && !errors.Is(serveErr, http.ErrServerClosed) {
			log.Printf("failed to start server: %v\n", serveErr)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), shutdownTimeout)
	defer cancel()

	err = server.Shutdown(ctx)
	if err != nil {
		log.Printf("failed to shutdown server: %v\n", err)
	}

	log.Println("server stopped")
}
