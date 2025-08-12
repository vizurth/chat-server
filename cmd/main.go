package main

import (
	"context"
	"fmt"
	"github.com/vizurth/chat-server/internal/config"
	"github.com/vizurth/chat-server/internal/postgres"
	"github.com/vizurth/chat-server/internal/server"
	desc "github.com/vizurth/chat-server/pkg/chat"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"log"
	"net"
)

const grpcPort string = ":50052"

func main() {
	ctx := context.Background()

	cfg, _ := config.NewConfig()

	db, err := postgres.New(ctx, cfg.Postgres)

	chatServer := server.NewServer(db)

	lis, err := net.Listen("tcp", fmt.Sprintf(":%s", cfg.Port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	s := grpc.NewServer()
	reflection.Register(s)

	desc.RegisterChatServer(s, chatServer)

	log.Printf("server listening at %v", lis.Addr())
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
