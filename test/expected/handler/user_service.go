package handler

import (
	"context"
	userpb "github.com/harryosmar/protobuf-go/gen/user"
	"github.com/harryosmar/protoc-gen-go-scaffold/repository"
	proto "github.com/harryosmar/protoc-gen-go-scaffold/testgen"
)

type UserServiceHandler interface {
	CreateUser(context.Context, *proto.CreateUserRequest) (*proto.CreateUserResponse, error)
	GetUser(context.Context, *proto.GetUserRequest) (*proto.GetUserResponse, error)
}

type userServiceHandler struct {
	userpb.UnimplementedUserServiceServer
	repo repository.UserServiceRepository
}

func NewUserServiceHandler(repo repository.UserServiceRepository) UserServiceHandler {
	return &userServiceHandler{repo: repo}
}

func (u *userServiceHandler) CreateUser(ctx context.Context, req *proto.CreateUserRequest) (*proto.CreateUserResponse, error) {
	return u.repo.CreateUser(ctx, req)
}

func (u *userServiceHandler) GetUser(ctx context.Context, req *proto.GetUserRequest) (*proto.GetUserResponse, error) {
	return u.repo.GetUser(ctx, req)
}
