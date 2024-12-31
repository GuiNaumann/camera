package repositories

import (
	"camera/domain/entities"
	"context"
)

// AuthenticationRepository provides access to the data related to repositories
type AuthenticationRepository interface {
	//UserExists - Check if user exists and have password hash stored on local database
	UserExists(ctx context.Context, credential entities.LoginCredentials) (exists bool, err error)

	//ComparePasswordHash - Check if user password hash matches
	ComparePasswordHash(ctx context.Context, login, password string) (bool, error)

	// GetUserByLogin returns the user associated with the provided login
	GetUserByLogin(ctx context.Context, login string) (*entities.User, error)

	GetUserByID(ctx context.Context, ID int64) (*entities.User, error)

	//EmailExists
	EmailExists(ctx context.Context, user entities.User) (bool, error)

	//DocumentExists
	DocumentExists(ctx context.Context, user entities.User) (bool, error)

	//RegisterUser
	RegisterUser(ctx context.Context, user entities.User) error
}
