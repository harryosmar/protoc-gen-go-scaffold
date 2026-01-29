package server

import (
	"context"
	"fmt"

	appError "github.com/example/myapp/error"
	"github.com/example/myapp/logger"
	"github.com/example/myapp/usecase"
	userpb "github.com/harryosmar/protobuf-go/gen/user"
	"go.uber.org/zap"
)

// UserServer implements the User with appUsecase pattern
type UserServer struct {
	userpb.UnimplementedUserServiceServer
	userUsecase usecase.UserUsecase
}

// NewUserServer creates a new UserServer instance
func NewUserServer(userUsecase usecase.UserUsecase) *UserServer {
	return &UserServer{
		userUsecase: userUsecase,
	}
}

// CreateUser implements the CreateUser RPC method
func (s *UserServer) CreateUser(ctx context.Context, req *CreateUserRequestDTO) (*CreateUserResponseDTO, error) {
	var (
		log = logger.FromContext(ctx)
		err error
	)
	defer func() {
		if err != nil {
			log.Error("UserServer.CreateUser err", zap.Error(err))
		}
	}()
	log.Info("User.CreateUser called", zap.String("req", fmt.Sprintf("%+v", req)))

	if err = req.Validate(); err != nil {
		return nil, appError.ErrInvalidArgument.WithMessage("validation failed: %v", err)
	}

	return s.userUsecase.CreateUser(ctx, req)
}

// GetUser implements the GetUser RPC method
func (s *UserServer) GetUser(ctx context.Context, req *GetUserRequestDTO) (*GetUserResponse, error) {
	var (
		log = logger.FromContext(ctx)
		err error
	)
	defer func() {
		if err != nil {
			log.Error("UserServer.GetUser err", zap.Error(err))
		}
	}()
	log.Info("User.GetUser called", zap.String("req", fmt.Sprintf("%+v", req)))

	if err = req.Validate(); err != nil {
		return nil, appError.ErrInvalidArgument.WithMessage("validation failed: %v", err)
	}

	return s.userUsecase.GetUser(ctx, req)
}

// DeleteUser implements the DeleteUser RPC method
func (s *UserServer) DeleteUser(ctx context.Context, req *DeleteUserRequestDTO) (*DeleteUserResponseDTO, error) {
	var (
		log = logger.FromContext(ctx)
		err error
	)
	defer func() {
		if err != nil {
			log.Error("UserServer.DeleteUser err", zap.Error(err))
		}
	}()
	log.Info("User.DeleteUser called", zap.String("req", fmt.Sprintf("%+v", req)))

	if err = req.Validate(); err != nil {
		return nil, appError.ErrInvalidArgument.WithMessage("validation failed: %v", err)
	}

	return s.userUsecase.DeleteUser(ctx, req)
}

// UpdateUser implements the UpdateUser RPC method
func (s *UserServer) UpdateUser(ctx context.Context, req *UpdateUserRequestDTO) (*UpdateUserResponseDTO, error) {
	var (
		log = logger.FromContext(ctx)
		err error
	)
	defer func() {
		if err != nil {
			log.Error("UserServer.UpdateUser err", zap.Error(err))
		}
	}()
	log.Info("User.UpdateUser called", zap.String("req", fmt.Sprintf("%+v", req)))

	if err = req.Validate(); err != nil {
		return nil, appError.ErrInvalidArgument.WithMessage("validation failed: %v", err)
	}

	return s.userUsecase.UpdateUser(ctx, req)
}

// ListUsers implements the ListUsers RPC method
func (s *UserServer) ListUsers(ctx context.Context, req *ListUsersRequestDTO) (*ListUsersResponseDTO, error) {
	var (
		log = logger.FromContext(ctx)
		err error
	)
	defer func() {
		if err != nil {
			log.Error("UserServer.ListUsers err", zap.Error(err))
		}
	}()
	log.Info("User.ListUsers called", zap.String("req", fmt.Sprintf("%+v", req)))

	if err = req.Validate(); err != nil {
		return nil, appError.ErrInvalidArgument.WithMessage("validation failed: %v", err)
	}

	return s.userUsecase.ListUsers(ctx, req)
}
