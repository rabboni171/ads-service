package transport

import (
	authv1 "auth-service/api/gen/proto"
	"auth-service/internal/config"
	"auth-service/internal/lib/types"
	"auth-service/internal/service"
	"context"
	"errors"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// Регистрируем нашу реализацию обработчиков контракта proto в grpc сервер
func AuthRegister(gRPCServer *grpc.Server, auth service.AuthServiceInterface) {
	authv1.RegisterAuthServer(gRPCServer, &serverAPI{authService: auth})
}

type serverAPI struct {
	// позволяет обеспечить обратную совместимость при изменении auth.proto
	// файла и позволит избежать ошибки, если забудем реализовать метод
	authv1.UnimplementedAuthServer
	authService service.AuthServiceInterface
}

// указываем, что serverAPI реализует интерфейс AuthServer
var _ authv1.AuthServer = (*serverAPI)(nil)

func (s *serverAPI) Login(
	ctx context.Context,
	in *authv1.LoginRequest,
) (*authv1.LoginResponse, error) {
	if in.GetEmail() == "" {
		return nil, status.Error(codes.InvalidArgument, types.ErrEmailRequired.Error())
	}
	if in.GetPassword() == "" {
		return nil, status.Error(codes.InvalidArgument, types.ErrPassRequired.Error())
	}
	if in.GetAppId() == 0 {
		return nil, status.Error(codes.InvalidArgument, types.ErrAppIdRequired.Error())
	}

	token, err := s.authService.Login(
		ctx,
		in.GetEmail(),
		in.GetPassword(),
		in.GetAppId(),
		config.Cfg.TokenTTL,
	)
	if err != nil {
		if errors.Is(err, types.ErrInvalidCredentials) {
			return nil, status.Error(codes.InvalidArgument, err.Error())
		}
		if errors.Is(err, types.ErrUserNotFound) {
			return nil, status.Error(codes.NotFound, err.Error())
		}
	
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &authv1.LoginResponse{Token: token}, nil
}

func (s *serverAPI) Register(
	ctx context.Context,
	in *authv1.RegisterRequest,
) (*authv1.RegisterResponse, error) {
	if in.GetEmail() == "" {
		return nil, status.Error(codes.InvalidArgument, types.ErrEmailRequired.Error())
	}
	if in.GetPassword() == "" {
		return nil, status.Error(codes.InvalidArgument, types.ErrPassRequired.Error())
	}

	uid, err := s.authService.RegisterNewUser(ctx, in.GetEmail(), in.GetPassword())
	if err != nil {
		if errors.Is(err, types.ErrUserExists) {
			return nil, status.Error(codes.AlreadyExists, err.Error())
		}

		return nil, status.Error(codes.Internal, "failed to register user: "+err.Error())
	}

	return &authv1.RegisterResponse{UserId: uid}, nil
}
