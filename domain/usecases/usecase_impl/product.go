package usecase_impl

import (
	"camera/domain/entities"
	"camera/domain/entities/rules"
	"camera/domain/usecases"
	"camera/domain/usecases/perm_impl"
	"camera/infrastructure/modules/impl/http_error"
	"camera/infrastructure/repositories"
	"camera/infrastructure/storage"
	"camera/settings_loader"
	"context"
	"log"
)

func NewProductUseCase(
	repo repositories.ProductRepository,
	settings settings_loader.SettingsLoader,
	fileStorage storage.FileStorageRepositoryNew,
) usecases.ProductUseCase {
	return perm_impl.NewPermProductUseCase(
		&productUseCase{
			repo:        repo,
			settings:    settings,
			fileStorage: fileStorage,
		})
}

type productUseCase struct {
	repo        repositories.ProductRepository
	settings    settings_loader.SettingsLoader
	fileStorage storage.FileStorageRepositoryNew
}

func (c productUseCase) CreateProductUseCase(ctx context.Context, user entities.User, product entities.Product) (int64, error) {
	err := rules.ProductRules(&product)
	if err != nil {
		log.Println("[CreateProductUseCase] Error productRules", err)
		return 0, err
	}

	products := c.repo.CheckLocalExist(ctx, product.LocalID)

	if product.LocalID == 0 || !products {
		log.Println("[CreateProductUseCase] Error Local não existe", err)
		return 0, http_error.NewUnauthorizedError("Local não existe")
	}

	id, err := c.repo.CreateProductRepository(ctx, product, user)
	if err != nil {
		log.Println("[CreateProductUseCase] Error CreateProductRepository", err)
		return 0, err
	}

	err = c.repo.SetProductStatusCode(ctx, product.Id, entities.StatusExist)
	if err != nil {
		log.Println("[CreateProductUseCase] Error SetProductStatusCode", err)
		return 0, err
	}

	return id, nil
}

func (c productUseCase) ListProductUseCase(
	ctx context.Context,
	user entities.User,
	filter entities.GeneralFilter,
) (*entities.PaginatedListUpdated[entities.Product], error) {
	return c.repo.ListProductRepository(ctx, filter, user)
}

func (c productUseCase) GetProductByIdUseCase(ctx context.Context, user entities.User, productID int64) (*entities.Product, error) {
	return c.repo.GetProductByIdRepository(ctx, productID, user)
}

func (c productUseCase) EditProductUseCase(ctx context.Context, user entities.User, product entities.Product) error {
	err := rules.ProductRulesEdite(&product)
	if err != nil {
		log.Println("[EditProductUseCase] Error ProductRules", err)
		return err
	}

	oldProduct, err := c.repo.GetProductByIdRepository(ctx, product.Id, user)
	if err != nil {
		log.Println("[EditProductUseCase] Error GetProductByIdRepository", err)
		return err
	}

	log.Println(oldProduct)

	err = c.repo.SetProductStatusCode(ctx, product.Id, entities.StatusExist)
	if err != nil {
		log.Println("[EditProductUseCase] Error SetProductStatusCode", err)
		return err
	}

	products := c.repo.CheckLocalExist(ctx, product.LocalID)

	if product.LocalID == 0 || !products {
		log.Println("[EditProductUseCase] Error Local não existe", err)
		return http_error.NewUnauthorizedError("Local não existe")
	}

	return c.repo.EditProductRepository(ctx, product, user)
}

func (c productUseCase) DeleteProductUseCase(ctx context.Context, user entities.User, productID int64) error {
	oldProduct, err := c.repo.GetProductByIdRepository(ctx, productID, user)
	if err != nil {
		log.Println("[DeleteProductUseCase] Error GetProductByIdRepository")
		return err
	}

	log.Println(oldProduct)

	return c.repo.DeleteProduct(ctx, productID)
}

func (c productUseCase) CreateLocalUseCase(ctx context.Context, user entities.User, product entities.Local) (int64, error) {
	err := rules.LocalRules(&product)
	if err != nil {
		log.Println("[CreateProductUseCase] Error productRules", err)
		return 0, err
	}

	id, err := c.repo.CreateLocalRepository(ctx, product, user)
	if err != nil {
		log.Println("[CreateProductUseCase] Error CreateProductRepository", err)
		return 0, err
	}

	err = c.repo.SetLocalStatusCode(ctx, product.Id, entities.StatusExist)
	if err != nil {
		log.Println("[CreateProductUseCase] Error SetProductStatusCode", err)
		return 0, err
	}

	return id, nil
}

func (c productUseCase) ListLocalUseCase(
	ctx context.Context,
	user entities.User,
	filter entities.GeneralFilter,
) (*entities.PaginatedListUpdated[entities.Local], error) {
	return c.repo.ListLocalRepository(ctx, filter, user)
}

func (c productUseCase) GetLocalByIdUseCase(ctx context.Context, user entities.User, productID int64) (*entities.Local, error) {
	return c.repo.GetLocalByIdRepository(ctx, productID, user)
}

func (c productUseCase) EditLocalUseCase(ctx context.Context, user entities.User, product entities.Local) error {
	err := rules.LocalRulesEdite(&product)
	if err != nil {
		log.Println("[EditProductUseCase] Error ProductRules", err)
		return err
	}

	oldProduct, err := c.repo.GetLocalByIdRepository(ctx, product.Id, user)
	if err != nil {
		log.Println("[EditProductUseCase] Error GetProductByIdRepository", err)
		return err
	}

	log.Println(oldProduct)

	err = c.repo.SetProductStatusCode(ctx, product.Id, entities.StatusExist)
	if err != nil {
		log.Println("[EditProductUseCase] Error SetProductStatusCode", err)
		return err
	}

	return c.repo.EditLocalRepository(ctx, product, user)
}

func (c productUseCase) DeleteLocalUseCase(ctx context.Context, user entities.User, productID int64) error {
	oldProduct, err := c.repo.GetLocalByIdRepository(ctx, productID, user)
	if err != nil {
		log.Println("[DeleteLocalUseCase] Error GetLocalByIdRepository")
		return err
	}

	log.Println(oldProduct)

	return c.repo.DeleteLocal(ctx, productID)
}

// TODO:
// TODO:
// TODO:
// TODO:
func (c productUseCase) SetParamiter(ctx context.Context, _ entities.User, productID int64) error {
	return c.repo.SetParamiter(ctx, productID)
}

func (c productUseCase) DeleteReadProduct(ctx context.Context, _ entities.User, productID int64) error {
	return c.repo.DeleteReadProduct(ctx, productID)
}

func (c productUseCase) ListReadProduct(
	ctx context.Context,
	user entities.User,
	filter entities.GeneralFilter,
) (*entities.PaginatedListUpdated[entities.Product], error) {
	return c.repo.ListReadProduct(ctx, filter, user)
}
