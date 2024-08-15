package main

import (
	"auth-service/internal/config"
	"auth-service/internal/logger"
	"auth-service/internal/service"
	"auth-service/internal/storage"
	"auth-service/internal/transport"
	"context"
	"fmt"
	"log/slog"
	"net"
)

func main() {
	config.MustLoad()
	logger.MustLoad(config.Cfg.Env)

	ctx := context.Background()
	storage := storage.NewStorage(ctx)
	service := service.NewService(storage.Auth)

	// Создаём слушальшика, который будет слушить TCP-сообщения, адресованные
	// Нашему gRPC-серверу
	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", config.Cfg.GRPC.Port))
	if err != nil {
		panic(err.Error())
	}

	gRPCServer := transport.NewServer(ctx, service)

	slog.Info("grpc server started", slog.String("addr", listener.Addr().String()))
	// запускаем на нашем слушальщике обработчик grpc
	if err := gRPCServer.Serve(listener); err != nil {
		panic(err.Error())
	}
}
