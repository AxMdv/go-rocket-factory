package main

import (
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"

	paymentV1API "github.com/AxMdv/go-rocket-factory/payment/internal/api/payment/v1"
	paymentService "github.com/AxMdv/go-rocket-factory/payment/internal/service/payment"
	paymentV1 "github.com/AxMdv/go-rocket-factory/shared/pkg/proto/payment/v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

const grpcPort = 50052

func main() {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", grpcPort))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	defer func() {
		if cerr := lis.Close(); cerr != nil {
			log.Printf("failed to close listener: %v", cerr)
		}
	}()
	// Создаем gRPC сервер
	grpcServer := grpc.NewServer()

	// Регистрируем наш сервис
	service := paymentService.NewService()
	api := paymentV1API.NewAPI(service)

	paymentV1.RegisterPaymentServiceServer(grpcServer, api)

	// Включаем рефлексию для отладки
	reflection.Register(grpcServer)

	go func() {
		log.Printf("🚀 gRPC PaymentService listening on %d\n", grpcPort)
		if err := grpcServer.Serve(lis); err != nil {
			log.Fatalf("failed to serve: %v", err)
		}
	}()

	// Graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("🛑 Shutting down gRPC payment server...")
	grpcServer.GracefulStop()
	log.Println("✅ Server stopped")
}
