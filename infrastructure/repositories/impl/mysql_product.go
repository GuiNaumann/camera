package impl

import (
	"camera/domain/entities"
	"camera/infrastructure/modules/impl/http_error"
	"camera/infrastructure/repositories"
	"camera/settings_loader"
	"camera/utils"
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

func (c productRepository) CreateProductRepository(ctx context.Context, camera entities.Product, user entities.User) (int64, error) {
	//language=sql
	query := `
	INSERT INTO camera (
		id_user,
		id_local,
		name,
		description,
		ip_address,
		port,
		username,
		password,
		stream_path,
		camera_type,
		is_active
	)
	VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, true)
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

	var cameraID int64
	err = stmt.QueryRowContext(
		ctx,
		user.ID,
		camera.Name,
		camera.Description,
		camera.IPAddress,
		camera.Port,
		camera.Username,
		camera.Password,
		camera.StreamPath,
		camera.CameraType,
	).Scan(&cameraID)

	if err != nil {
		log.Println("[CreateCameraRepository] Error QueryRowContext", err)
		_ = tx.Rollback() // Rollback em caso de erro
		return 0, err
	}

	// Commit da transação
	if err := tx.Commit(); err != nil {
		log.Println("[CreateCameraRepository] Error Commit", err)
		return 0, err
	}

	log.Printf("Câmera inserida com sucesso: ID %d", cameraID)

	return cameraID, nil
}

func (c productRepository) SetProductStatusCode(ctx context.Context, productID int64, statusCode entities.StatusCode) error {
	//language=sql
	query := `
	UPDATE camera 
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
	//language=sql
	queryCount := `
	    SELECT COUNT(*)
	    FROM camera
	    WHERE is_active = true
	      AND id_user = $1`

	//language=sql
	query := `
	SELECT DISTINCT c.id,
	                c.name,
	                c.description,
	                c.ip_address,
	                c.port,
	                c.username,
	                c.password,
	                c.stream_path,
	                c.camera_type,
	                c.is_active,
	                c.status_code,
	                c.created_at,
	                c.modified_at
	FROM camera c
	WHERE c.is_active = true
	  AND c.id_user = $1`

	trimSearch := strings.TrimSpace(filter.Search)
	var searchContaining string

	if trimSearch != "" {
		searchContaining = "%" + strings.ToLower(trimSearch) + "%"
		query += ` AND LOWER(c.name) LIKE $2`
		queryCount += " AND LOWER(name) LIKE $2"
	}

	query += ` ORDER BY c.modified_at DESC`

	if filter.Limit > 0 {
		firstItem := filter.Page * filter.Limit
		query += fmt.Sprintf(" LIMIT %d OFFSET %d", filter.Limit, firstItem)
	}

	stmt, err := c.conn.PrepareContext(ctx, query)
	if err != nil {
		log.Println("[ListCameraRepository] Error PrepareContext", err)
		return nil, http_error.NewUnexpectedError(http_error.Unexpected)
	}
	defer stmt.Close()

	var rows *sql.Rows
	if trimSearch != "" {
		rows, err = stmt.QueryContext(ctx, user.ID, searchContaining)
	} else {
		rows, err = stmt.QueryContext(ctx, user.ID)
	}
	if err != nil {
		log.Println("[ListCameraRepository] Error QueryContext", err)
		return nil, http_error.NewUnexpectedError(http_error.Unexpected)
	}
	defer rows.Close()

	var cameras = make([]entities.Product, 0)
	for rows.Next() {
		var camera entities.Product
		err = rows.Scan(
			&camera.Id,
			&camera.Name,
			&camera.Description,
			&camera.IPAddress,
			&camera.Port,
			&camera.Username,
			&camera.Password,
			&camera.StreamPath,
			&camera.CameraType,
			&camera.IsActive,
			&camera.StatusCode,
			&camera.CreatedAt,
			&camera.ModifiedAt,
		)
		if err != nil {
			log.Println("[ListCameraRepository] Error Scan", err)
			return nil, http_error.NewUnexpectedError(http_error.Unexpected)
		}
		cameras = append(cameras, camera)
	}

	stmtCount, err := c.conn.PrepareContext(ctx, queryCount)
	if err != nil {
		log.Println("[ListCameraRepository] Error stmtCount PrepareContext", err)
		return nil, http_error.NewUnexpectedError(http_error.Unexpected)
	}
	defer stmtCount.Close()

	var totalCount int64
	if trimSearch != "" {
		err = stmtCount.QueryRowContext(ctx, user.ID, searchContaining).Scan(&totalCount)
	} else {
		err = stmtCount.QueryRowContext(ctx, user.ID).Scan(&totalCount)
	}
	if err != nil {
		log.Println("[ListCameraRepository] Error stmtCount Scan", err)
		return nil, http_error.NewUnexpectedError(http_error.Unexpected)
	}

	return &entities.PaginatedListUpdated[entities.Product]{
		Items:          cameras,
		RequestedCount: int64(len(cameras)),
		TotalCount:     totalCount,
		Page:           filter.Page + 1,
	}, nil
}

func (c productRepository) GetProductByIdRepository(ctx context.Context, productID int64, user entities.User) (*entities.Product, error) {
	//language=sql
	query := `
	SELECT c.id,
	       c.name,
	       c.description,
	       c.ip_address,
	       c.port,
	       c.username,
	       c.password,
	       c.stream_path,
	       c.camera_type,
	       c.is_active,
	       c.status_code,
	       c.created_at,
	       c.modified_at
	FROM camera c
	WHERE c.id = $1
	  AND c.status_code != $2
	  AND c.id_user = $3`

	var camera entities.Product
	err := c.conn.QueryRowContext(ctx, query, productID, entities.StatusDeleted, user.ID).Scan(
		&camera.Id,
		&camera.Name,
		&camera.Description,
		&camera.IPAddress,
		&camera.Port,
		&camera.Username,
		&camera.Password,
		&camera.StreamPath,
		&camera.CameraType,
		&camera.IsActive,
		&camera.StatusCode,
		&camera.CreatedAt,
		&camera.ModifiedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			log.Println("[GetCameraByIdRepository] Error sql.ErrNoRows", err)
			return nil, http_error.NewBadRequestError(http_error.CameraNotFound)
		}
		log.Println("[GetCameraByIdRepository] Error Scan", err)
		return nil, http_error.NewUnexpectedError(http_error.Unexpected)
	}

	return &camera, nil
}

func (c productRepository) EditProductRepository(ctx context.Context, camera entities.Product, user entities.User) error {
	//language=sql
	command := `
	UPDATE camera
	SET name = $1,
	    description = $2,
	    ip_address = $3,
	    port = $4,
	    username = $5,
	    password = $6,
	    stream_path = $7,
	    camera_type = $8,
	    is_active = $9,
	    id_user = $10
	    id_local = $11
	WHERE id = $12`

	_, err := c.conn.ExecContext(ctx,
		command,
		camera.Name,
		camera.Description,
		camera.IPAddress,
		camera.Port,
		camera.Username,
		camera.Password,
		camera.StreamPath,
		camera.CameraType,
		camera.IsActive,
		user.ID,
		camera.LocalID,
		camera.Id,
	)
	if err != nil {
		log.Println("[EditCameraRepository] Error ExecContext", err)
		return http_error.NewUnexpectedError(http_error.Unexpected)
	}

	return nil
}

func (c productRepository) DeleteProduct(ctx context.Context, productID int64) error {
	//language=sql
	command := `
	UPDATE camera
	SET status_code = 2
	WHERE id = $1`

	_, err := c.conn.ExecContext(
		ctx,
		command,
		productID,
	)
	if err != nil {
		log.Println("[DeleteCamera] Error ExecContext", err)
		return http_error.NewUnexpectedError(http_error.Unexpected)
	}

	return nil
}

//TODO: ================================================================================================================
//TODO: LOCAL ==========================================================================================================
//TODO: ================================================================================================================

func (c productRepository) CreateLocalRepository(ctx context.Context, local entities.Local, user entities.User) (int64, error) {
	//language=sql
	query := `
	INSERT INTO local (
		id_user,
		name,
		description,
		state,
		city,
		street,
		number,
		zip_code,
		complement,
		is_active
	)
	VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, true)
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

	var localID int64
	err = stmt.QueryRowContext(
		ctx,
		user.ID,
		local.Name,
		local.Description,
		local.State,
		local.City,
		local.Street,
		local.Number,
		local.ZipCode,
		local.Complement,
	).Scan(&localID)

	if err != nil {
		log.Println("[CreateLocalRepository] Error QueryRowContext", err)
		_ = tx.Rollback()
		return 0, err
	}

	if err := tx.Commit(); err != nil {
		log.Println("[CreateLocalRepository] Error Commit", err)
		return 0, err
	}

	log.Printf("Local inserido com sucesso: ID %d", localID)

	return localID, nil
}

func (c productRepository) SetLocalStatusCode(ctx context.Context, localID int64, statusCode entities.StatusCode) error {
	//language=sql
	query := `
	UPDATE local
	SET status_code = $1
	WHERE id = $2
	`

	_, err := c.conn.ExecContext(ctx, query, statusCode, localID)
	if err != nil {
		log.Println("[SetLocalStatusCode] Error ExecContext", err)
		return http_error.NewUnexpectedError(http_error.Unexpected)
	}

	return nil
}

func (c productRepository) ListLocalRepository(
	ctx context.Context,
	filter entities.GeneralFilter,
	user entities.User,
) (*entities.PaginatedListUpdated[entities.Local], error) {
	//language=sql
	queryCount := `
	    SELECT COUNT(*)
	    FROM local
	    WHERE is_active = true
	      AND id_user = $1`

	//language=sql
	query := `
	SELECT DISTINCT l.id,
	                l.name,
	                l.description,
	                l.state,
	                l.city,
	                l.street,
	                l.number,
	                l.zip_code,
	                l.complement,
	                l.is_active,
	                l.status_code,
	                l.created_at,
	                l.modified_at
	FROM local l
	WHERE l.is_active = true
	  AND l.id_user = $1`

	trimSearch := strings.TrimSpace(filter.Search)
	var searchContaining string

	if trimSearch != "" {
		searchContaining = "%" + strings.ToLower(trimSearch) + "%"
		query += ` AND LOWER(l.name) LIKE $2`
		queryCount += " AND LOWER(name) LIKE $2"
	}

	query += ` ORDER BY l.modified_at DESC`

	if filter.Limit > 0 {
		firstItem := filter.Page * filter.Limit
		query += fmt.Sprintf(" LIMIT %d OFFSET %d", filter.Limit, firstItem)
	}

	stmt, err := c.conn.PrepareContext(ctx, query)
	if err != nil {
		log.Println("[ListLocalRepository] Error PrepareContext", err)
		return nil, http_error.NewUnexpectedError(http_error.Unexpected)
	}
	defer stmt.Close()

	var rows *sql.Rows
	if trimSearch != "" {
		rows, err = stmt.QueryContext(ctx, user.ID, searchContaining)
	} else {
		rows, err = stmt.QueryContext(ctx, user.ID)
	}
	if err != nil {
		log.Println("[ListLocalRepository] Error QueryContext", err)
		return nil, http_error.NewUnexpectedError(http_error.Unexpected)
	}
	defer rows.Close()

	var locals = make([]entities.Local, 0)
	for rows.Next() {
		var local entities.Local
		err = rows.Scan(
			&local.Id,
			&local.Name,
			&local.Description,
			&local.State,
			&local.City,
			&local.Street,
			&local.Number,
			&local.ZipCode,
			&local.Complement,
			&local.IsActive,
			&local.StatusCode,
			&local.CreatedAt,
			&local.ModifiedAt,
		)
		if err != nil {
			log.Println("[ListLocalRepository] Error Scan", err)
			return nil, http_error.NewUnexpectedError(http_error.Unexpected)
		}
		locals = append(locals, local)
	}

	stmtCount, err := c.conn.PrepareContext(ctx, queryCount)
	if err != nil {
		log.Println("[ListLocalRepository] Error stmtCount PrepareContext", err)
		return nil, http_error.NewUnexpectedError(http_error.Unexpected)
	}
	defer stmtCount.Close()

	var totalCount int64
	if trimSearch != "" {
		err = stmtCount.QueryRowContext(ctx, user.ID, searchContaining).Scan(&totalCount)
	} else {
		err = stmtCount.QueryRowContext(ctx, user.ID).Scan(&totalCount)
	}
	if err != nil && !errors.Is(sql.ErrNoRows, err) {
		log.Println("[ListLocalRepository] Error stmtCount Scan", err)
		return nil, http_error.NewUnexpectedError(http_error.Unexpected)
	}

	return &entities.PaginatedListUpdated[entities.Local]{
		Items:          locals,
		RequestedCount: int64(len(locals)),
		TotalCount:     totalCount,
		Page:           filter.Page + 1,
	}, nil
}

func (c productRepository) GetLocalByIdRepository(ctx context.Context, localID int64, user entities.User) (*entities.Local, error) {
	//language=sql
	query := `
	SELECT l.id,
	       l.name,
	       l.description,
	       l.state,
	       l.city,
	       l.street,
	       l.number,
	       l.zip_code,
	       l.complement,
	       l.is_active,
	       l.status_code,
	       l.created_at,
	       l.modified_at
	FROM local l
	WHERE l.id = $1
	  AND l.status_code != $2
	  AND l.id_user = $3`

	var local entities.Local
	err := c.conn.QueryRowContext(ctx, query, localID, entities.StatusDeleted, user.ID).Scan(
		&local.Id,
		&local.Name,
		&local.Description,
		&local.State,
		&local.City,
		&local.Street,
		&local.Number,
		&local.ZipCode,
		&local.Complement,
		&local.IsActive,
		&local.StatusCode,
		&local.CreatedAt,
		&local.ModifiedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			log.Println("[GetLocalByIdRepository] Error sql.ErrNoRows", err)
			return nil, http_error.NewBadRequestError(http_error.LocalNotFound)
		}
		log.Println("[GetLocalByIdRepository] Error Scan", err)
		return nil, http_error.NewUnexpectedError(http_error.Unexpected)
	}

	return &local, nil
}

func (c productRepository) EditLocalRepository(ctx context.Context, local entities.Local, user entities.User) error {
	//language=sql
	command := `
	UPDATE local
	SET name = $1,
	    description = $2,
	    state = $3,
	    city = $4,
	    street = $5,
	    number = $6,
	    zip_code = $7,
	    complement = $8,
	    is_active = $9,
	    id_user = $10
	WHERE id = $11`

	_, err := c.conn.ExecContext(ctx,
		command,
		local.Name,
		local.Description,
		local.State,
		local.City,
		local.Street,
		local.Number,
		local.ZipCode,
		local.Complement,
		local.IsActive,
		user.ID,
		local.Id,
	)
	if err != nil {
		log.Println("[EditLocalRepository] Error ExecContext", err)
		return http_error.NewUnexpectedError(http_error.Unexpected)
	}

	return nil
}

func (c productRepository) DeleteLocal(ctx context.Context, localID int64) error {
	//language=sql
	command := `
	UPDATE local
	SET status_code = 2
	WHERE id = $1`

	_, err := c.conn.ExecContext(ctx, command, localID)
	if err != nil {
		log.Println("[DeleteLocal] Error ExecContext", err)
		return http_error.NewUnexpectedError(http_error.Unexpected)
	}

	return nil
}

// TODO: ================================================================================================================
// TODO: ================================================================================================================
// TODO: ================================================================================================================

func (c productRepository) ListCameraStreams(
	ctx context.Context,
	filter entities.GeneralFilter,
	user entities.User,
	screenCount int, // Número de telas (1 a 64)
) ([]entities.CameraStream, error) {
	// Limitar o número de telas para evitar excesso
	if screenCount < 1 {
		screenCount = 1
	} else if screenCount > 64 {
		screenCount = 64
	}

	//language=sql
	query := `
    SELECT id,
           name,
           description,
           ip_address,
           port,
           username,
           password,
           stream_path
    FROM camera
    WHERE is_active = true
      AND id_user = $1
    ORDER BY name
    LIMIT $2`

	stmt, err := c.conn.PrepareContext(ctx, query)
	if err != nil {
		log.Println("[ListCameraStreams] Error PrepareContext", err)
		return nil, http_error.NewUnexpectedError(http_error.Unexpected)
	}
	defer stmt.Close()

	rows, err := stmt.QueryContext(ctx, user.ID, screenCount)
	if err != nil {
		log.Println("[ListCameraStreams] Error QueryContext", err)
		return nil, http_error.NewUnexpectedError(http_error.Unexpected)
	}
	defer rows.Close()

	var streams []entities.CameraStream
	for rows.Next() {
		var id int64
		var name, description, ip, streamPath, username, password string
		var port int

		err := rows.Scan(&id, &name, &description, &ip, &port, &username, &password, &streamPath)
		if err != nil {
			log.Println("[ListCameraStreams] Error Scan", err)
			return nil, http_error.NewUnexpectedError(http_error.Unexpected)
		}

		// Gerar a URL RTSP
		streamURL := fmt.Sprintf("rtsp://%s:%s@%s:%d%s", username, password, ip, port, streamPath)

		streams = append(streams, entities.CameraStream{
			CameraID:    id,
			Name:        name,
			Description: description,
			StreamURL:   streamURL,
		})
	}

	return streams, nil
}

// TODO: ================================================================================================================
// TODO: ================================================================================================================
// TODO: ================================================================================================================

// TODO: ================================================================================================================
// TODO: NAO SERÁ USADO =================================================================================================
// TODO: ================================================================================================================
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
			//&product.Quantidade,
			//&product.Preco,
			//&product.Tamanho,
			//&product.ImageURL,
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
