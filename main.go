package main

import (
    "com.example.shopping_server/proto/shopping"
    "context"
    "github.com/golang/protobuf/ptypes/empty"
    "google.golang.org/grpc"
    "log"
    "net"
)

const (
    port = ":50051"
)

type server struct {}

func main() {
    lis, err := net.Listen("tcp", port)
    if err != nil {
        log.Fatalf("failed to listen: %v", err)
    }
    s := grpc.NewServer()
    com_example.RegisterShoppingServiceServer(s, &server{})
    if err := s.Serve(lis); err != nil {
        log.Fatalf("failed to serve: %v", err)
    }
}

func (s *server) Auth(ctx context.Context, request *com_example.AuthRequest) (*com_example.AuthResponse, error) {
    return &com_example.AuthResponse{
        Token: "token",
    }, nil
}

func (s *server) AddToCart(ctx context.Context, request *com_example.AddToCartRequest) (*com_example.AddToCartResponse, error) {
    return &com_example.AddToCartResponse{
        ItemIds: map[string]int32{"11111": 1},
    }, nil
}

func (s *server) CreateOrder(ctx context.Context, request *empty.Empty) (*com_example.OrderResponse, error) {
    return &com_example.OrderResponse{
        ItemIds: map[string]int32{"11111": 1},
        OrderId: "12345",
    }, nil
}

func (s *server) GetOrder(ctx context.Context, request *com_example.GetOrderRequest) (*com_example.OrderResponse, error) {
    return &com_example.OrderResponse{
        OrderId: "12345",
        ItemIds: map[string]int32{"11111": 1},
    }, nil
}
