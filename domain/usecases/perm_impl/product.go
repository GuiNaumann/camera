package perm_impl

import (
	"bear/domain/entities"
	"bear/domain/usecases"
	"bear/infrastructure/modules/impl/http_error"
	"context"
)

type productPermUseCase struct {
	perm usecases.ProductUseCase
}

func NewPermProductUseCase(ProductUseCase usecases.ProductUseCase) usecases.ProductUseCase {
	return &productPermUseCase{
		perm: ProductUseCase,
	}
}

func (c productPermUseCase) CreateProductUseCase(ctx context.Context, user entities.User, product entities.Product) (int64, error) {
	if !user.IsMaster() && !user.IsFlat3() && !user.IsFlat2() && !user.IsFlat1() {
		return 0, http_error.NewUnauthorizedError(http_error.Unauthorized)
	}

	return c.perm.CreateProductUseCase(ctx, user, product)
}

func (c productPermUseCase) ListProductUseCase(
	ctx context.Context,
	user entities.User,
	filter entities.GeneralFilter,
) (*entities.PaginatedListUpdated[entities.Product], error) {
	if user.IsMaster() || user.IsFlat3() || user.IsFlat2() || user.IsFlat1() {
		return c.perm.ListProductUseCase(ctx, user, filter)
	}

	return nil, http_error.NewUnauthorizedError(http_error.Unauthorized)
}

func (c productPermUseCase) GetProductByIdUseCase(ctx context.Context, user entities.User, productID int64) (*entities.Product, error) {
	if user.IsMaster() || user.IsFlat3() || user.IsFlat2() || user.IsFlat1() {
		return c.perm.GetProductByIdUseCase(ctx, user, productID)
	}

	return nil, http_error.NewUnauthorizedError(http_error.Unauthorized)
}

func (c productPermUseCase) EditProductUseCase(ctx context.Context, user entities.User, product entities.Product) error {
	if !user.IsMaster() && !user.IsFlat3() && !user.IsFlat2() && !user.IsFlat1() {
		return http_error.NewUnauthorizedError(http_error.Unauthorized)
	}

	return c.perm.EditProductUseCase(ctx, user, product)
}

func (c productPermUseCase) DeleteProductUseCase(ctx context.Context, user entities.User, productID int64) error {
	if !user.IsMaster() && !user.IsFlat3() && !user.IsFlat2() && !user.IsFlat1() {
		return http_error.NewUnauthorizedError(http_error.Unauthorized)
	}

	return c.perm.DeleteProductUseCase(ctx, user, productID)
}

func (c productPermUseCase) SetParamiter(ctx context.Context, user entities.User, productID int64) error {
	if !user.IsMaster() && !user.IsFlat3() && !user.IsFlat2() && !user.IsFlat1() {
		return http_error.NewUnauthorizedError(http_error.Unauthorized)
	}

	return c.perm.SetParamiter(ctx, user, productID)
}

func (c productPermUseCase) DeleteReadProduct(ctx context.Context, user entities.User, productID int64) error {
	if !user.IsMaster() && !user.IsFlat3() && !user.IsFlat2() && !user.IsFlat1() {
		return http_error.NewUnauthorizedError(http_error.Unauthorized)
	}

	return c.perm.DeleteReadProduct(ctx, user, productID)
}

func (c productPermUseCase) ListReadProduct(
	ctx context.Context,
	user entities.User,
	filter entities.GeneralFilter,
) (*entities.PaginatedListUpdated[entities.Product], error) {
	if user.IsMaster() || user.IsFlat3() || user.IsFlat2() || user.IsFlat1() {
		return c.perm.ListReadProduct(ctx, user, filter)
	}

	return nil, http_error.NewUnauthorizedError(http_error.Unauthorized)
}
