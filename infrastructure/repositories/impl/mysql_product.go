package impl

import (
	"bear/domain/entities"
	"bear/infrastructure/modules/impl/http_error"
	"bear/infrastructure/repositories"
	"bear/settings_loader"
	"bear/utils"
	"context"
	"database/sql"
	"fmt"
	"github.com/pkg/errors"
	"log"
	"math"
	"strings"
)

func NewProductRepository(
	db *sql.DB,
	settings settings_loader.SettingsLoader,
) repositories.ProductRepository {
	return &productRepository{
		conn:     db,
		settings: settings,
	}
}

type productRepository struct {
	settings settings_loader.SettingsLoader
	conn     *sql.DB
}

func (c productRepository) CreateProductRepository(ctx context.Context, product entities.Product, user entities.User) (int64, error) {
	//language=sql
	query := `
	INSERT INTO product (
		id_user,
		name,
		description,
		quantidade,
		preco,
		tamanho,
		image_url,
		is_active
	)
	VALUES ($1,$2,$3,$4,$5,$6,$7,true)
	RETURNING id
	`

	tx, err := c.conn.Begin()
	if err != nil {
		return 0, err
	}

	stmt, err := tx.PrepareContext(ctx, query)
	if err != nil {
		_ = tx.Rollback()
		return 0, err
	}
	defer stmt.Close()

	var productID int64
	err = stmt.QueryRowContext(
		ctx,
		user.ID,
		product.Name,
		product.Description,
		product.Quantidade,
		product.Preco,
		product.Tamanho,
		product.ImageURL,
	).Scan(&productID)

	if err != nil {
		log.Println("[CreateProductRepository] Error QueryRowContext", err)
		_ = tx.Rollback() // Rollback em caso de erro
		return 0, err
	}

	// Commit da transação
	if err := tx.Commit(); err != nil {
		log.Println("[CreateProductRepository] Error Commit", err)
		return 0, err
	}

	log.Printf("Produto inserido com sucesso: ID %d", productID)

	return productID, nil
}

func (c productRepository) SetProductStatusCode(ctx context.Context, productID int64, statusCode entities.StatusCode) error {
	//language=sql
	query := `
	UPDATE product 
	SET status_code = $1
	WHERE id = $2
	`

	_, err := c.conn.ExecContext(ctx, query, statusCode, productID)
	if err != nil {
		log.Println("[SetproductStatusCode] Error ExecContext", err)
		return http_error.NewUnexpectedError(http_error.Unexpected)
	}

	return nil
}

func (c productRepository) ListProductRepository(
	ctx context.Context,
	filter entities.GeneralFilter,
	user entities.User,
) (*entities.PaginatedListUpdated[entities.Product], error) {
	////language=sql
	queryCount := `
	    SELECT COUNT(*)
	  FROM product
	 WHERE status_code != 2
	 AND id_user = $1`

	//language=sql
	query := `
	SELECT DISTINCT c.id,
	                c.name,
	                c.description,
					c.quantidade,
					c.preco,
					c.tamanho,
					c.image_url,
					c.is_active,
					c.status_code,
					c.created_at,
					c.modified_at
	FROM product c
	WHERE c.status_code != 2
	AND c.id_user = $1 `

	trimSearch := strings.TrimSpace(filter.Search)
	var formattedSearch string
	var searchStartingWith string
	var searchContaining string
	var searchEndingWith string

	if trimSearch != "" {
		cleanedString := utils.CleanMySQLRegexp(trimSearch)
		formattedSearch = cleanedString

		searchStartingWith = cleanedString + "%"
		searchContaining = "%" + cleanedString + "%"
		searchEndingWith = "%" + cleanedString

		query += ` AND LOWER(c.name) LIKE LOWER($2)
		ORDER BY
		CASE
        WHEN c.name LIKE $3 THEN 1 
        WHEN c.name LIKE $4 THEN 2 
        WHEN c.name LIKE $5 THEN 3 
        ELSE 4
        END,
		`
		queryCount += "\n AND name REGEXP $2 "
	}
	var ordinationAsc string
	if filter.OrdinationAsc {
		ordinationAsc = "ASC"
	} else {
		ordinationAsc = "DESC"
	}

	if trimSearch == "" {
		query += ` ORDER BY `
	}
	switch filter.Column {
	case "name":
		query += ` c.name ` + ordinationAsc
		break
	default:
		column := " c.name %s"
		if trimSearch == "" {
			column = " c.modified_at %s"
		}
		query += fmt.Sprintf(column, ordinationAsc)
	}

	if filter.Limit != 0 {
		if filter.Page > 0 {
			filter.Page-- // Ajusta para índice baseado em 0
		}
		firstItem := filter.Page * filter.Limit
		query += fmt.Sprintf(" LIMIT %d OFFSET %d", filter.Limit, firstItem)
	} else {
		filter.Limit = math.MaxInt
	}

	stmt, err := c.conn.PrepareContext(ctx, query)
	if err != nil {
		log.Println("[ListProductRepository] Error PrepareContext", err)
		return nil, http_error.NewUnexpectedError(http_error.Unexpected)
	}
	defer stmt.Close()

	var rows *sql.Rows
	if formattedSearch != "" {
		rows, err = stmt.QueryContext(ctx, user.ID, searchContaining, searchStartingWith, searchContaining, searchEndingWith)
	} else {
		rows, err = stmt.QueryContext(ctx, user.ID)
	}
	if err != nil {
		log.Println("[ListProductRepository] Error QueryContext", err)
		return nil, http_error.NewUnexpectedError(http_error.Unexpected)
	}
	defer rows.Close()

	var products = make([]entities.Product, 0)
	for rows.Next() {
		var product entities.Product
		err = rows.Scan(
			&product.Id,
			&product.Name,
			&product.Description,
			&product.Quantidade,
			&product.Preco,
			&product.Tamanho,
			&product.ImageURL,
			&product.IsActive,
			&product.StatusCode,
			&product.CreatedAt,
			&product.ModifiedAt,
		)
		if err != nil {
			log.Println("[ListProductRepository] Error Scan", err)
			return nil, http_error.NewUnexpectedError(http_error.Unexpected)
		}
		products = append(products, product)
	}

	stmtCount, err := c.conn.PrepareContext(ctx, queryCount)
	if err != nil {
		log.Println("[ListProductRepository] Error stmtCount PrepareContext", err)
		return nil, http_error.NewUnexpectedError(http_error.Unexpected)
	}
	defer stmtCount.Close()

	var totalCount int64
	if formattedSearch != "" {
		err = stmtCount.QueryRowContext(ctx, user.ID, formattedSearch).Scan(&totalCount)
	} else {
		err = stmtCount.QueryRowContext(ctx, user.ID).Scan(&totalCount)
	}
	if err != nil && !errors.Is(sql.ErrNoRows, err) {
		log.Println("[ListProductRepository] Error stmtCount Scan", err)
		return nil, http_error.NewUnexpectedError(http_error.Unexpected)
	}

	mathPage := float64(totalCount) / float64(filter.Limit)
	page := int64(math.Ceil(mathPage))

	return &entities.PaginatedListUpdated[entities.Product]{
		Items:          products,
		RequestedCount: int64(len(products)),
		TotalCount:     totalCount,
		Page:           page,
	}, nil
}

func (c productRepository) GetProductByIdRepository(ctx context.Context, productID int64, user entities.User) (*entities.Product, error) {
	//language=sql
	query := `
	SELECT c.id,
		   c.name,
		   c.description,
		   c.quantidade,
		   c.preco,
		   c.tamanho,
		   c.image_url,
		   c.is_active,
		   c.status_code,
		   c.created_at,
		   c.modified_at
	FROM product c
	WHERE c.id = $1
	  AND c.status_code != $2
	  AND c.id_user = $3`

	var product entities.Product
	err := c.conn.QueryRowContext(ctx, query, productID, entities.StatusDeleted, user.ID).Scan(
		&product.Id,
		&product.Name,
		&product.Description,
		&product.Quantidade,
		&product.Preco,
		&product.Tamanho,
		&product.ImageURL,
		&product.IsActive,
		&product.StatusCode,
		&product.ModifiedAt,
		&product.CreatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			log.Println("[GetProductByIdRepository] Error sql.ErrNoRows", err)
			return nil, http_error.NewBadRequestError(http_error.ProductNotFound)
		}
		log.Println("[GetProductByIdRepository] Error Scan", err)
		return nil, http_error.NewUnexpectedError(http_error.Unexpected)
	}

	return &product, nil
}

func (c productRepository) EditProductRepository(ctx context.Context, product entities.Product, user entities.User) error {
	//language=sql
	command := `
	UPDATE product
	SET name = $1,
	    description = $2,
		quantidade = $3,
		preco = $4,
		tamanho = $5,
		image_url = $6,
		is_active = $7,
		id_user = $8
	WHERE id = $9`

	_, err := c.conn.ExecContext(ctx,
		command,
		product.Name,
		product.Description,
		product.Quantidade,
		product.Preco,
		product.Tamanho,
		product.ImageURL,
		product.IsActive,
		user.ID,
		product.Id,
	)
	if err != nil {
		log.Println("[EditProductRepository] Error ExecContext", err)
		return http_error.NewUnexpectedError(http_error.Unexpected)
	}

	return nil
}

func (c productRepository) DeleteProduct(ctx context.Context, productID int64) error {
	//language=sql
	command := `
	UPDATE product
	SET status_code = 2
	WHERE id = $1`

	_, err := c.conn.ExecContext(
		ctx,
		command,
		productID,
	)
	if err != nil {
		log.Println("[DeleteProduct] Error ExecContext", err)
		return http_error.NewUnexpectedError(http_error.Unexpected)
	}

	return nil
}

func (c productRepository) SetParamiter(ctx context.Context, productID int64) error {
	//language=sql
	command := `
	UPDATE product
	SET parameter = true
	WHERE id = $1`

	_, err := c.conn.ExecContext(
		ctx,
		command,
		productID,
	)
	if err != nil {
		log.Println("[SetParamiter] Error ExecContext", err)
		return http_error.NewUnexpectedError(http_error.Unexpected)
	}

	return nil
}

func (c productRepository) DeleteReadProduct(ctx context.Context, productID int64) error {
	//language=sql
	command := `
	UPDATE product
	SET parameter = false
	WHERE id = $1`

	_, err := c.conn.ExecContext(
		ctx,
		command,
		productID,
	)
	if err != nil {
		log.Println("[DeleteReadProduct] Error ExecContext", err)
		return http_error.NewUnexpectedError(http_error.Unexpected)
	}

	return nil
}

func (c productRepository) ListReadProduct(
	ctx context.Context,
	filter entities.GeneralFilter,
	user entities.User,
) (*entities.PaginatedListUpdated[entities.Product], error) {
	////language=sql
	queryCount := `
	    SELECT COUNT(*)
	  FROM product
	 WHERE status_code != 2
	 AND parameter = true
	 AND id_user = $1`

	//language=sql
	query := `
	SELECT DISTINCT c.id,
	                c.name,
	                c.description,
					c.quantidade,
					c.preco,
					c.tamanho,
					c.image_url,
					c.is_active,
					c.status_code,
					c.created_at,
					c.modified_at
	FROM product c
	WHERE c.status_code != 2
	AND c.parameter = true
	AND c.id_user = $1 `

	trimSearch := strings.TrimSpace(filter.Search)
	var formattedSearch string
	var searchStartingWith string
	var searchContaining string
	var searchEndingWith string

	if trimSearch != "" {
		cleanedString := utils.CleanMySQLRegexp(trimSearch)
		formattedSearch = cleanedString

		searchStartingWith = cleanedString + "%"
		searchContaining = "%" + cleanedString + "%"
		searchEndingWith = "%" + cleanedString

		query += ` AND LOWER(c.name) LIKE LOWER($2)
		ORDER BY
		CASE
        WHEN c.name LIKE $3 THEN 1 
        WHEN c.name LIKE $4 THEN 2 
        WHEN c.name LIKE $5 THEN 3 
        ELSE 4
        END,
		`
		queryCount += "\n AND name REGEXP $2 "
	}
	var ordinationAsc string
	if filter.OrdinationAsc {
		ordinationAsc = "ASC"
	} else {
		ordinationAsc = "DESC"
	}

	if trimSearch == "" {
		query += ` ORDER BY `
	}
	switch filter.Column {
	case "name":
		query += ` c.name ` + ordinationAsc
		break
	default:
		column := " c.name %s"
		if trimSearch == "" {
			column = " c.modified_at %s"
		}
		query += fmt.Sprintf(column, ordinationAsc)
	}

	if filter.Limit != 0 {
		if filter.Page > 0 {
			filter.Page-- // Ajusta para índice baseado em 0
		}
		firstItem := filter.Page * filter.Limit
		query += fmt.Sprintf(" LIMIT %d OFFSET %d", filter.Limit, firstItem)
	} else {
		filter.Limit = math.MaxInt
	}

	stmt, err := c.conn.PrepareContext(ctx, query)
	if err != nil {
		log.Println("[ListReadProduct] Error PrepareContext", err)
		return nil, http_error.NewUnexpectedError(http_error.Unexpected)
	}
	defer stmt.Close()

	var rows *sql.Rows
	if formattedSearch != "" {
		rows, err = stmt.QueryContext(ctx, user.ID, searchContaining, searchStartingWith, searchContaining, searchEndingWith)
	} else {
		rows, err = stmt.QueryContext(ctx, user.ID)
	}
	if err != nil {
		log.Println("[ListReadProduct] Error QueryContext", err)
		return nil, http_error.NewUnexpectedError(http_error.Unexpected)
	}
	defer rows.Close()

	var products = make([]entities.Product, 0)
	for rows.Next() {
		var product entities.Product
		err = rows.Scan(
			&product.Id,
			&product.Name,
			&product.Description,
			&product.Quantidade,
			&product.Preco,
			&product.Tamanho,
			&product.ImageURL,
			&product.IsActive,
			&product.StatusCode,
			&product.CreatedAt,
			&product.ModifiedAt,
		)
		if err != nil {
			log.Println("[ListReadProduct] Error Scan", err)
			return nil, http_error.NewUnexpectedError(http_error.Unexpected)
		}
		products = append(products, product)
	}

	stmtCount, err := c.conn.PrepareContext(ctx, queryCount)
	if err != nil {
		log.Println("[ListReadProduct] Error stmtCount PrepareContext", err)
		return nil, http_error.NewUnexpectedError(http_error.Unexpected)
	}
	defer stmtCount.Close()

	var totalCount int64
	if formattedSearch != "" {
		err = stmtCount.QueryRowContext(ctx, user.ID, formattedSearch).Scan(&totalCount)
	} else {
		err = stmtCount.QueryRowContext(ctx, user.ID).Scan(&totalCount)
	}
	if err != nil && !errors.Is(sql.ErrNoRows, err) {
		log.Println("[ListReadProduct] Error stmtCount Scan", err)
		return nil, http_error.NewUnexpectedError(http_error.Unexpected)
	}

	mathPage := float64(totalCount) / float64(filter.Limit)
	page := int64(math.Ceil(mathPage))

	return &entities.PaginatedListUpdated[entities.Product]{
		Items:          products,
		RequestedCount: int64(len(products)),
		TotalCount:     totalCount,
		Page:           page,
	}, nil
}
