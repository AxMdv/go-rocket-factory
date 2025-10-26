// nolint:gosec
package main

import (
	"context"
	"fmt"
	"log"
	"math/rand"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/brianvoe/gofakeit/v6"
	"github.com/google/uuid"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/reflection"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"

	inventoryV1 "github.com/AxMdv/go-rocket-factory/shared/pkg/proto/inventory/v1"
)

const grpcPort = 50051

// inventoryService реализует gRPC сервис для работы со складом запчастей
type inventoryService struct {
	inventoryV1.UnimplementedInventoryServiceServer
	parts map[string]*inventoryV1.Part
}

func (s *inventoryService) generateParts(count int) {
	// Используем общий генератор создающий и возвращающий карту деталей,
	// затем просто присваиваем её в s.parts.
	s.parts = createParts(count)
}

// createParts генерирует и возвращает map[string]*inventoryV1.Part с помощью gofakeit.
func createParts(count int) map[string]*inventoryV1.Part {
	// Инициализация генератора
	gofakeit.Seed(time.Now().UnixNano())

	now := timestamppb.New(time.Now())
	parts := make(map[string]*inventoryV1.Part, count)

	categories := []inventoryV1.Category{
		inventoryV1.Category_CATEGORY_ENGINE,
		inventoryV1.Category_CATEGORY_FUEL,
		inventoryV1.Category_CATEGORY_PORTHOLE,
		inventoryV1.Category_CATEGORY_WING,
	}

	// Вспомогательные генераторы
	genDimensions := func() *inventoryV1.Dimensions {
		return &inventoryV1.Dimensions{
			Length: gofakeit.Float64Range(0.1, 50.0), // метры
			Width:  gofakeit.Float64Range(0.1, 10.0),
			Height: gofakeit.Float64Range(0.1, 10.0),
			Weight: gofakeit.Float64Range(0.5, 5000.0), // кг
		}
	}
	genManufacturer := func() *inventoryV1.Manufacturer {
		return &inventoryV1.Manufacturer{
			Name:    gofakeit.Company(),
			Country: gofakeit.Country(),
			Website: gofakeit.URL(),
		}
	}
	genTags := func() []string {
		base := []string{"space", "rocket", "module", "spare", "core", "nano", "mkII", "mkIII"}
		n := gofakeit.Number(1, 4)
		out := make([]string, 0, n)
		for i := 0; i < n; i++ {
			out = append(out, base[rand.Intn(len(base))])
		}
		return out
	}
	genMetadata := func() map[string]*inventoryV1.Value {
		// Пример смешанного metadata: string, number, bool
		return map[string]*inventoryV1.Value{
			"batch": {
				Kind: &inventoryV1.Value_StringValue{
					StringValue: gofakeit.UUID(),
				},
			},
			"lifetime_hours": {
				Kind: &inventoryV1.Value_DoubleValue{
					DoubleValue: gofakeit.Float64Range(100, 100000),
				},
			},
			"refurbished": {
				Kind: &inventoryV1.Value_BoolValue{
					BoolValue: gofakeit.Bool(),
				},
			},
		}
	}

	for i := 0; i < count; i++ {
		id := uuid.New().String()
		name := gofakeit.Company() + " Part"
		desc := gofakeit.Sentence(8)

		cat := categories[gofakeit.Number(0, len(categories)-1)]
		price := gofakeit.Price(100.0, 500000.0)
		stock := int64(gofakeit.Number(0, 500))

		part := &inventoryV1.Part{
			Uuid:          id,
			Name:          name,
			Description:   desc,
			Price:         price,
			StockQuantity: stock,
			Category:      cat,
			Dimensions:    genDimensions(),
			Manufacturer:  genManufacturer(),
			Tags:          genTags(),
			Metadata:      genMetadata(),
			CreatedAt:     now,
			UpdatedAt:     now,
		}
		parts[id] = part
	}

	return parts
}

// GetPart возвращает информацию о детали по её идентификатору.
func (s *inventoryService) GetPart(ctx context.Context, req *inventoryV1.GetPartRequest) (*inventoryV1.GetPartResponse, error) {
	part, ok := s.parts[req.GetUuid()]
	if !ok {
		return nil, status.Errorf(codes.NotFound, "part with UUID %s not found", req.GetUuid())
	}
	return &inventoryV1.GetPartResponse{
		Part: part,
	}, nil
}

func (s *inventoryService) ListParts(ctx context.Context, req *inventoryV1.ListPartsRequest) (*inventoryV1.ListPartsResponse, error) {
	filter := req.GetFilter()
	if filter == nil {
		parts := make([]*inventoryV1.Part, 0, len(s.parts))
		for _, part := range s.parts {
			parts = append(parts, part)
		}
		return &inventoryV1.ListPartsResponse{Parts: parts}, nil
	}
	// storedParts := make([]*inventoryV1.Part, 0, len(s.parts))
	// for _, part := range s.parts {
	// 	storedParts = append(storedParts, part)
	// }
	filteredParts := filterParts(s.parts, filter)
	return &inventoryV1.ListPartsResponse{
		Parts: filteredParts,
	}, nil
}

func main() {
	lis, err := net.Listen("tcp", fmt.Sprintf("localhost:%d", grpcPort))
	if err != nil {
		log.Printf("failed to listen: %v\n", err)
		return
	}

	// Создаем gRPC сервер
	s := grpc.NewServer()

	// Регистрируем наш сервис
	service := &inventoryService{
		parts: make(map[string]*inventoryV1.Part),
	}
	service.generateParts(15)

	inventoryV1.RegisterInventoryServiceServer(s, service)

	// Включаем рефлексию для отладки
	reflection.Register(s)

	go func() {
		log.Printf("🚀 inventory gRPC server listening on %d\n", grpcPort)
		err = s.Serve(lis)
		if err != nil {
			log.Printf("failed to serve: %v\n", err)
			return
		}
	}()
	defer func() {
		if cerr := lis.Close(); cerr != nil {
			log.Printf("failed to close listener: %v\n", cerr)
		}
	}()
	// Graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("🛑 Shutting down gRPC inventory server...")
	s.GracefulStop()
	log.Println("✅ Server stopped")
}
