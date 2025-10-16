package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"

	"github.com/google/uuid"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/reflection"
	"google.golang.org/grpc/status"

	paymentV1 "github.com/AxMdv/go-rocket-factory/shared/pkg/proto/payment/v1"
)

const grpcPort = 50052

// paymentService реализует gRPC сервис оплаты заказов
type paymentService struct {
	paymentV1.UnimplementedPaymentServiceServer
}

// PayOrder обрабатывает команду на оплату, генерирует UUID транзакции, выводит в лог и возвращает вызвавшему
func (s *paymentService) PayOrder(ctx context.Context, req *paymentV1.PayOrderRequest) (*paymentV1.PayOrderResponse, error) {
	orderUUID := req.GetOrderUuid()
	userUUID := req.GetUserUuid()
	if orderUUID == "" || userUUID == "" {
		return nil, status.Error(codes.InvalidArgument, "order_uuid or user_uuid not specified")
	}

	// Проверка  payment_method
	pm := req.GetPaymentMethod()
	if pm == paymentV1.PaymentMethod_PAYMENT_METHOD_UNKNOWN {
		return nil, status.Error(codes.InvalidArgument, "payment_method must be specified and not UNKNOWN")
	}
	switch pm {
	case paymentV1.PaymentMethod_PAYMENT_METHOD_CARD,
		paymentV1.PaymentMethod_PAYMENT_METHOD_SBP,
		paymentV1.PaymentMethod_PAYMENT_METHOD_CREDIT_CARD,
		paymentV1.PaymentMethod_PAYMENT_METHOD_INVESTOR_MONEY:
	default:
		return nil, status.Errorf(codes.InvalidArgument, "unsupported payment_method: %v", pm)
	}

	transactionUUID := uuid.New().String()
	log.Printf("Оплата прошла успешно, transaction_uuid: %s", transactionUUID)

	return &paymentV1.PayOrderResponse{
		TransactionUuid: transactionUUID,
	}, nil
}

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

	grpcServer := grpc.NewServer()
	paymentV1.RegisterPaymentServiceServer(grpcServer, &paymentService{})

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
