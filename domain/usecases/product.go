package usecases

import (
	"camera/domain/entities"
	"context"
)

type ProductUseCase interface {
	//CreateProductUseCase - Create Product and return id of Product
	CreateProductUseCase(ctx context.Context, user entities.User, product entities.Product) (int64, error)

	//ListProductsUseCase Return a list of all Product with status code 0
	ListProductUseCase(
		ctx context.Context,
		user entities.User,
		filter entities.GeneralFilter,
	) (*entities.PaginatedListUpdated[entities.Product], error)

	//GetProductByIdUseCase Get a Product by id and return the Product
	GetProductByIdUseCase(ctx context.Context, user entities.User, productID int64) (*entities.Product, error)

	//Editertificate Edit information about Product
	EditProductUseCase(ctx context.Context, user entities.User, product entities.Product) error

	//DeleteProduct - delete an Clertificate
	DeleteProductUseCase(ctx context.Context, user entities.User, productID int64) error

	CreateLocalUseCase(ctx context.Context, user entities.User, lcoal entities.Local) (int64, error)

	ListLocalUseCase(
		ctx context.Context,
		user entities.User,
		filter entities.GeneralFilter,
	) (*entities.PaginatedListUpdated[entities.Local], error)

	GetLocalByIdUseCase(ctx context.Context, user entities.User, localID int64) (*entities.Local, error)

	EditLocalUseCase(ctx context.Context, user entities.User, local entities.Local) error

	DeleteLocalUseCase(ctx context.Context, user entities.User, localID int64) error

	//TODO:

	SetParamiter(ctx context.Context, user entities.User, productID int64) error

	DeleteReadProduct(ctx context.Context, user entities.User, productID int64) error

	ListReadProduct(
		ctx context.Context,
		user entities.User,
		filter entities.GeneralFilter,
	) (*entities.PaginatedListUpdated[entities.Product], error)
}
