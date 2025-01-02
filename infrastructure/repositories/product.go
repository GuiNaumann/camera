package repositories

import (
	"camera/domain/entities"
	"context"
)

type ProductRepository interface {
	//CreateProductRepository - Create Product and return id of product
	CreateProductRepository(ctx context.Context, product entities.Product, user entities.User) (int64, error)

	CheckLocalExist(ctx context.Context, localID int64) bool

	SetProductStatusCode(ctx context.Context, productID int64, statusCode entities.StatusCode) error

	//ListProductRepository Return a list of all Product with status code 0
	ListProductRepository(
		ctx context.Context,
		filter entities.GeneralFilter,
		user entities.User,
	) (*entities.PaginatedListUpdated[entities.Product], error)

	//GetProductByIdRepository Get a Product by id
	GetProductByIdRepository(ctx context.Context, productID int64, user entities.User) (*entities.Product, error)

	//EditProductRepository - Edit the instructor
	EditProductRepository(ctx context.Context, product entities.Product, user entities.User) error

	//DeleteProduct - Set status code of Product to StatusDeleted
	DeleteProduct(ctx context.Context, productID int64) error

	CreateLocalRepository(ctx context.Context, product entities.Local, user entities.User) (int64, error)

	SetLocalStatusCode(ctx context.Context, productID int64, statusCode entities.StatusCode) error

	//ListProductRepository Return a list of all Product with status code 0
	ListLocalRepository(
		ctx context.Context,
		filter entities.GeneralFilter,
		user entities.User,
	) (*entities.PaginatedListUpdated[entities.Local], error)

	//GetProductByIdRepository Get a Product by id
	GetLocalByIdRepository(ctx context.Context, localID int64, user entities.User) (*entities.Local, error)

	//EditProductRepository - Edit the instructor
	EditLocalRepository(ctx context.Context, local entities.Local, user entities.User) error

	//DeleteProduct - Set status code of Product to StatusDeleted
	DeleteLocal(ctx context.Context, localID int64) error

	//TODO:
	//TODO:
	//TODO:
	//TODO:
	SetParamiter(ctx context.Context, productID int64) error

	DeleteReadProduct(ctx context.Context, productID int64) error

	ListReadProduct(
		ctx context.Context,
		filter entities.GeneralFilter,
		user entities.User,
	) (*entities.PaginatedListUpdated[entities.Product], error)
}
