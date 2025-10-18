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

// Константы статусов заказа
const (
	StatusPendingPayment = "PENDING_PAYMENT"
	StatusPaid           = "PAID"
	StatusCancelled      = "CANCELLED"
)

func main() {
	ctx := context.Background()

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

	// Создаем gRPC клиенты
	inventoryClient := inventoryV1.NewInventoryServiceClient(inventoryConn)
	paymentClient := paymentV1.NewPaymentServiceClient(paymentConn)

	orderService := NewOrderService(inventoryClient, paymentClient)

	handler := NewOrderHandler(orderService)
	orderServer, err := orderV1.NewServer(handler)
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
