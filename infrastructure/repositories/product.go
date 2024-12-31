package repositories

import (
	"bear/domain/entities"
	"context"
)

type ProductRepository interface {
	//CreateProductRepository - Create Product and return id of product
	CreateProductRepository(ctx context.Context, product entities.Product, user entities.User) (int64, error)

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

	SetParamiter(ctx context.Context, productID int64) error

	DeleteReadProduct(ctx context.Context, productID int64) error

	ListReadProduct(
		ctx context.Context,
		filter entities.GeneralFilter,
		user entities.User,
	) (*entities.PaginatedListUpdated[entities.Product], error)
}
