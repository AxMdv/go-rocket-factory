package main

import (
	"context"
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
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	orderV1API "github.com/AxMdv/go-rocket-factory/order/internal/api/order/v1"
	inventoryClientV1 "github.com/AxMdv/go-rocket-factory/order/internal/client/grpc/inventory/v1"
	paymentClientV1 "github.com/AxMdv/go-rocket-factory/order/internal/client/grpc/payment/v1"
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
	// Таймауты для HTTP-сервера
	readHeaderTimeout = 5 * time.Second
	shutdownTimeout   = 10 * time.Second
)

func main() {
	inventoryConn, err := grpc.NewClient(
		inventoryServiceAddr,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		log.Printf("failed to connect to inventory: %v\n", err)
		return
	}

	paymentConn, err := grpc.NewClient(
		paymentServiceAddr,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		log.Printf("failed to connect to payment: %v\n", err)
		return
	}

	// Создаем gRPC клиенты
	ic := inventoryV1.NewInventoryServiceClient(inventoryConn)
	pc := paymentV1.NewPaymentServiceClient(paymentConn)

	// Создаем реализацию gRPC клиентов (сущности для сервисного слоя)
	invClient := inventoryClientV1.NewInventoryClient(ic)
	payClient := paymentClientV1.NewPaymentClient(pc)

	repo := orderRepository.NewRepository()
	orderSvc := orderService.NewOrderService(repo, invClient, payClient)

	api := orderV1API.NewAPI(orderSvc)
	orderServer, err := orderV1.NewServer(api)
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		if cerr := paymentConn.Close(); cerr != nil {
			log.Printf("failed to close connect to payment: %v", cerr)
		}
	}()
	defer func() {
		if cerr := inventoryConn.Close(); cerr != nil {
			log.Printf("failed to close connect to inventory: %v", cerr)
		}
	}()
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
		log.Printf("🚀 HTTP-сервер запущен на порту %s\n", httpPort)
		err = server.ListenAndServe()
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Printf("❌ Ошибка запуска сервера: %v\n", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("🛑 Завершение работы сервера...")

	ctx, cancel := context.WithTimeout(context.Background(), shutdownTimeout)
	defer cancel()

	err = server.Shutdown(ctx)
	if err != nil {
		log.Printf("❌ Ошибка при остановке сервера: %v\n", err)
	}

	log.Println("✅ Сервер остановлен")
}
