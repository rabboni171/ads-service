package auth

import (
	authv1 "ads-service/api/auth-service/gen/proto"
	"context"
	"log/slog"
	"time"

	grpclog "github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/logging"
	grpcretry "github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/retry"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
)

type Client struct {
	Api authv1.AuthClient
}

func New(
	ctx context.Context,
	addr string, // Адрес auth-сервера
	timeout time.Duration, // Таймаут на выполнение каждой попытки
	retriesCount int, // Количетсво попыток
) (*Client, error) {

	// Опции для перехватчика grpcretry, то есть повторые попытки отправления запросов
	retryOpts := []grpcretry.CallOption{
		grpcretry.WithCodes(codes.NotFound, codes.Aborted, codes.DeadlineExceeded),
		grpcretry.WithMax(uint(retriesCount)),
		grpcretry.WithPerRetryTimeout(timeout),
	}

	// Опции для перехватчика grpclog, логируем запросы и ответы
	logOpts := []grpclog.Option{
		grpclog.WithLogOnEvents(grpclog.PayloadReceived, grpclog.PayloadSent),
	}

	// Создаём экземпляр соединения с gRPC-сервером auth с созданными перехватчиками
	cc, err := grpc.DialContext(ctx, addr,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithChainUnaryInterceptor(
			grpclog.UnaryClientInterceptor(InterceptorLogger(), logOpts...),
			grpcretry.UnaryClientInterceptor(retryOpts...),
		))
	if err != nil {
		return nil, err
	}

	// Создаём gRPC-клиент auth
	grpcClient := authv1.NewAuthClient(cc)

	return &Client{
		Api: grpcClient,
	}, nil
}

// Перехватывает логи из grpc в наш логгер
func InterceptorLogger() grpclog.Logger {
	return grpclog.LoggerFunc(func(ctx context.Context, lvl grpclog.Level, msg string, fields ...any) {
		slog.Log(ctx, slog.Level(lvl), msg, fields...)
	})
}
