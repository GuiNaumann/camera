package usecase_impl

import (
	"bear/domain/entities"
	"bear/domain/entities/rules"
	"bear/domain/usecases"
	"bear/domain/usecases/perm_impl"
	"bear/infrastructure/modules/impl/http_error"
	"bear/infrastructure/repositories"
	"bear/infrastructure/storage"
	"bear/settings_loader"
	"context"
	"fmt"
	"github.com/google/uuid"
	"log"
	"path/filepath"
	"strings"
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

	if product.ImageBase64 != "" && !storage.IsURL(product.ImageBase64) {
		generated, err := uuid.NewUUID()
		if err != nil {
			log.Println("[CreateProductUseCase] Error NewUUID", err)
			return 0, http_error.NewUnexpectedError(http_error.Unexpected)
		}

		filePath := fmt.Sprintf("/products/%s", generated)

		product.ImageURL, err = c.fileStorage.SaveBase64(product.ImageBase64, filePath)
		if err != nil {
			log.Println("[CreateProductUseCase] Error SaveBase64", err)
			return 0, http_error.NewUnexpectedError(http_error.Unexpected)
		}
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

	if product.ImageBase64 != "" && !storage.IsURL(product.ImageBase64) {
		if oldProduct.ImageURL != "" {
			_, fileName := filepath.Split(strings.Split(oldProduct.ImageURL, "?")[0])

			err = c.fileStorage.DeletePath(filepath.Join("images", "products", fileName))
			if err != nil {
				log.Println("[EditProductUseCase] Error DeletePath", err)
				return http_error.NewUnexpectedError(http_error.Unexpected)
			}
		}

		generated, err := uuid.NewUUID()
		if err != nil {
			log.Println("[EditProductUseCase] Error NewUUID", err)
			return http_error.NewUnexpectedError(http_error.Unexpected)
		}

		filePath := fmt.Sprintf("/products/%s", generated)

		product.ImageURL, err = c.fileStorage.SaveBase64(product.ImageBase64, filePath)
		if err != nil {
			log.Println("[EditProductUseCase] Error SaveBase64", err)
			return http_error.NewUnexpectedError(http_error.Unexpected)
		}
	}

	err = c.repo.SetProductStatusCode(ctx, product.Id, entities.StatusExist)
	if err != nil {
		log.Println("[EditProductUseCase] Error SetProductStatusCode", err)
		return err
	}

	return c.repo.EditProductRepository(ctx, product, user)
}

func (c productUseCase) DeleteProductUseCase(ctx context.Context, user entities.User, productID int64) error {
	oldProduct, err := c.repo.GetProductByIdRepository(ctx, productID, user)
	if err != nil {
		log.Println("[DeleteProductUseCase] Error GetProductByIdRepository")
		return err
	}

	if oldProduct.ImageBase64 != "" {
		_, fileName := filepath.Split(strings.Split(oldProduct.ImageURL, "?")[0])

		err = c.fileStorage.DeletePath(filepath.Join("images", "productImages", fileName))
		if err != nil {
			log.Println("[DeleteProductUseCase] Error DeletePath")
			return http_error.NewUnexpectedError(http_error.Unexpected)
		}
	}

	return c.repo.DeleteProduct(ctx, productID)
}

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
