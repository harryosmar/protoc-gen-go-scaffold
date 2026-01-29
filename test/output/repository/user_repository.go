package repository

import (
	userpb "github.com/harryosmar/protobuf-go/gen/user"
	"gorm.io/gorm"
)

// userRepositoryMySQL implements UserRepository interface
type userRepositoryMySQL struct {
	*BaseGorm[userpb.UserEntityORM, uint32]
}

// NewUserRepositoryMySQL creates a new user repository instance
func NewUserRepositoryMySQL(db *gorm.DB) ServiceRepository[userpb.UserEntityORM, uint32] {
	return &userRepositoryMySQL{
		NewBaseGorm[userpb.UserEntityORM, uint32](db),
	}
}
