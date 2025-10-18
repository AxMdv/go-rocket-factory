package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"

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

// Структура для десериализации JSON (примерно совпадает с protobuf, но с enum в строковом виде)
type partJSON struct {
	Uuid          string                        `json:"uuid"`
	Name          string                        `json:"name"`
	Description   string                        `json:"description"`
	Price         float64                       `json:"price"`
	StockQuantity int64                         `json:"stock_quantity"`
	Category      string                        `json:"category"`
	Dimensions    *inventoryV1.Dimensions       `json:"dimensions"`
	Manufacturer  *inventoryV1.Manufacturer     `json:"manufacturer"`
	Tags          []string                      `json:"tags"`
	Metadata      map[string]*inventoryV1.Value `json:"metadata"`
}

// Функция загрузки деталей в map[string]*Part
func (s *inventoryService) loadPartsFromJSON(filepath string) error {
	data, err := os.ReadFile(filepath)
	if err != nil {
		return err
	}

	var partsFromJSON []partJSON
	if err := json.Unmarshal(data, &partsFromJSON); err != nil {
		return err
	}

	now := timestamppb.New(time.Now())
	s.parts = make(map[string]*inventoryV1.Part, len(partsFromJSON))

	for _, p := range partsFromJSON {
		// Конвертируем строку category в enum
		catValue, ok := inventoryV1.Category_value[p.Category]
		if !ok {
			catValue = int32(inventoryV1.Category_CATEGORY_UNKNOWN)
		}

		s.parts[p.Uuid] = &inventoryV1.Part{
			Uuid:          p.Uuid,
			Name:          p.Name,
			Description:   p.Description,
			Price:         p.Price,
			StockQuantity: p.StockQuantity,
			Category:      inventoryV1.Category(catValue),
			Dimensions:    p.Dimensions,
			Manufacturer:  p.Manufacturer,
			Tags:          p.Tags,
			Metadata:      p.Metadata,
			CreatedAt:     now,
			UpdatedAt:     now,
		}
	}

	return nil
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
	defer func() {
		if cerr := lis.Close(); cerr != nil {
			log.Printf("failed to close listener: %v\n", cerr)
		}
	}()

	// Создаем gRPC сервер
	s := grpc.NewServer()

	// Регистрируем наш сервис
	service := &inventoryService{
		parts: make(map[string]*inventoryV1.Part),
	}
	service.loadPartsFromJSON("./parts.json")
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

	// Graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("🛑 Shutting down gRPC inventory server...")
	s.GracefulStop()
	log.Println("✅ Server stopped")
}
