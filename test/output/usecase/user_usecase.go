package usecase

import (
	"context"
	appError "github.com/example/myapp/error"
	"github.com/example/myapp/repository"
	userpb "github.com/harryosmar/protobuf-go/gen/user"
)

// UserUsecase defines the interface for user business logic
type UserUsecase interface {
	// Methods matching the proto service definition
	CreateUser(ctx context.Context, req *CreateUserRequestDTO) (*CreateUserResponseDTO, error)
	GetUser(ctx context.Context, req *GetUserRequestDTO) (*GetUserResponse, error)
	DeleteUser(ctx context.Context, req *DeleteUserRequestDTO) (*DeleteUserResponseDTO, error)
	UpdateUser(ctx context.Context, req *UpdateUserRequestDTO) (*UpdateUserResponseDTO, error)
	ListUsers(ctx context.Context, req *ListUsersRequestDTO) (*ListUsersResponseDTO, error)
}

// userUsecase implements UserUsecase interface
type userUsecase struct {
	userRepo repository.ServiceRepository[userpb.UserEntityORM, uint32]
}

// NewUserUsecase creates a new user usecase instance
func NewUserUsecase(repo repository.ServiceRepository[userpb.UserEntityORM, uint32]) UserUsecase {
	return &userUsecase{
		userRepo: repo,
	}
}

// ormToDTO converts ORM model to DTO
func (u *userUsecase) ormToDTO(orm *userpb.UserEntityORM) *userpb.UserDTO {
	return &userpb.UserDTO{
		Name:  orm.Name,
		Email: orm.Email,
		Id:    orm.Id,
	}
}

// ormToDTOList converts a slice of ORM models to DTO list
func (u *userUsecase) ormToDTOList(ormRecords []userpb.UserEntityORM) []*userpb.UserDTO {
	var dtoRecords []*userpb.UserDTO
	for _, record := range ormRecords {
		dtoRecords = append(dtoRecords, u.ormToDTO(&record))
	}
	return dtoRecords
}

// CreateUser implements the CreateUser RPC method from the proto service
func (u *userUsecase) CreateUser(ctx context.Context, req *CreateUserRequestDTO) (*CreateUserResponseDTO, error) {
	// Create user entity from DTO
	dto := req.User
	entity := &userpb.UserEntity{
		Id:        dto.Id,
		Name:      dto.Name,
		Email:     dto.Email,
		CreatedAt: dto.CreatedAt,
		UpdatedAt: dto.UpdatedAt,
	}

	// Convert to ORM model for database operations
	userORM, err := entity.ToORM(ctx)
	if err != nil {
		return nil, err
	}

	// Save to database using repository
	newUserORM, err := u.userRepo.Create(ctx, &userORM)
	if err != nil {
		return nil, err
	}

	return &userpb.CreateUserResponseDTO{
		User: &userpb.UserDTO{
			Name:  newUserORM.Name,
			Email: newUserORM.Email,
			Id:    newUserORM.Id,
		},
	}, nil
}

// GetUser implements the GetUser RPC method from the proto service
func (u *userUsecase) GetUser(ctx context.Context, req *GetUserRequestDTO) (*GetUserResponse, error) {
	// Query database for user using repository
	orm, err := u.userRepo.GetById(ctx, req.Id)
	if err != nil {
		return nil, err
	}
	if orm == nil {
		return nil, appError.ErrUserNotFound
	}

	return &userpb.GetUserResponse{
		User: u.ormToDTO(orm),
	}, nil
}

// DeleteUser implements the DeleteUser RPC method from the proto service
func (u *userUsecase) DeleteUser(ctx context.Context, req *DeleteUserRequestDTO) (*DeleteUserResponseDTO, error) {
	err := u.userRepo.Delete(ctx, req.Id)
	if err != nil {
		return nil, err
	}

	return &userpb.DeleteUserResponseDTO{}, nil
}

// UpdateUser implements the UpdateUser RPC method from the proto service
func (u *userUsecase) UpdateUser(ctx context.Context, req *UpdateUserRequestDTO) (*UpdateUserResponseDTO, error) {
	// Get existing user first
	orm, err := u.userRepo.GetById(ctx, req.User.Id)
	if err != nil {
		return nil, err
	}
	if orm == nil {
		return nil, appError.ErrUserNotFound
	}

	// Update fields
	orm.Name = req.User.Name
	orm.Email = req.User.Email

	// Update in database using repository
	if _, err := u.userRepo.Update(ctx, orm); err != nil {
		return nil, err
	}

	return &userpb.UpdateUserResponseDTO{
		User: req.User,
	}, nil
}

// ListUsers implements the ListUsers RPC method from the proto service
func (u *userUsecase) ListUsers(ctx context.Context, req *ListUsersRequestDTO) (*ListUsersResponseDTO, error) {
	ormRecords, paginator, err := u.userRepo.GetPerPage(
		ctx,
		req.Pagination.Page,
		req.Pagination.Limit,
		[]repository.OrderBy{},
		[]repository.Where{},
	)
	if err != nil {
		return nil, err
	}

	return &userpb.ListUsersResponseDTO{
		Users: u.ormToDTOList(ormRecords),
		Pagination: &userpb.PaginationResponse{
			Total: paginator.Total,
			Page:  paginator.Page,
			Limit: paginator.PerPage,
		},
	}, nil
}
