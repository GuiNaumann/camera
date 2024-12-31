package impl

import (
	"bear/domain/entities"
	"bear/infrastructure/modules/impl/http_error"
	"bear/infrastructure/repositories"
	"bear/settings_loader"
	"context"
	"database/sql"
	"golang.org/x/crypto/bcrypt"
	"log"
)

func NewAuthenticationRepository(
	db *sql.DB,
	settings settings_loader.SettingsLoader,
) repositories.AuthenticationRepository {
	return &authRepository{
		conn:     db,
		settings: settings,
	}
}

type authRepository struct {
	settings settings_loader.SettingsLoader
	conn     *sql.DB
}

func (r *authRepository) UserExists(
	ctx context.Context,
	credential entities.LoginCredentials,
) (exists bool, err error) {
	// language=sql
	query := `
		SELECT id
		FROM public.user au
		WHERE (au.user_document = $1 AND au.user_document != '') 
			OR (au.user_email = $2 AND au.user_email != '') 
		LIMIT 1
	`

	row, err := r.conn.QueryContext(
		ctx,
		query,
		credential.Login,
		credential.Login,
	)
	if err != nil {
		log.Println("[UserExists] Error QueryContext", err)
		return false, http_error.NewBadRequestError(http_error.Unexpected)
	}
	defer row.Close()

	return row.Next(), nil
}

func (r *authRepository) ComparePasswordHash(
	ctx context.Context,
	login string,
	password string,
) (bool, error) {
	// language=sql
	query := `
	SELECT password_hash 
	FROM public.user
	WHERE status_code = 0 
	  AND ((user_document = $1 AND user_document != '') 
	   OR (user_email = $2 AND user_email != ''))
	LIMIT 1`

	stmt, err := r.conn.PrepareContext(ctx, query)
	if err != nil {
		log.Println("[ComparePasswordHash] Error PrepareContext", err)
		return false, http_error.NewBadRequestError(http_error.Unexpected)
	}
	defer stmt.Close()

	row, err := stmt.QueryContext(ctx, login, login)
	if err != nil {
		log.Println("[ComparePasswordHash] Error QueryContext", err)
		return false, http_error.NewBadRequestError(http_error.Unexpected)
	}
	defer row.Close()

	if row.Next() {
		var it entities.LoginCredentials
		err = row.Scan(
			&it.Password,
		)
		if err != nil {
			return false, err
		}

		//if password and hash match, the user will log in
		err := bcrypt.CompareHashAndPassword([]byte(it.Password), []byte(password))
		if err != nil {
			log.Println("")
			log.Println("it.Password", it.Password)
			log.Println("password", password)
			log.Println("")
			log.Println("[ComparePasswordHash] Error CompareHashAndPassword", err)
			return false, nil
		}
	}

	return true, nil
}

func (r *authRepository) GetUserByLogin(
	ctx context.Context,
	login string,
) (*entities.User, error) {
	// language=sql
	query := `
	SELECT au.id, 
	       au.user_name,      
	       au.user_document,
	       au.user_email,
	       ut.numeric_code,
	       au.is_active,
	       au.status_code,
	       au.password_modified_at,
	       COALESCE(au.badge_code, 0)
	  FROM public.user au
		INNER JOIN public.user_type aut on au.id = aut.id_user
		INNER JOIN public.type_user ut on aut.id_user_type = ut.id
	 WHERE (au.user_document = $1 AND au.user_document != '') 
	    OR (au.user_email = $2 AND au.user_email != '')`

	rows, err := r.conn.QueryContext(
		ctx,
		query,
		login,
		login,
	)
	if err != nil {
		return nil, http_error.NewBadRequestError(http_error.Unexpected)
	}
	defer rows.Close()

	for rows.Next() {
		var it entities.User
		err = rows.Scan(
			&it.ID,
			&it.Name,
			&it.Document,
			&it.Email,
			&it.UserType,
			&it.IsActive,
			&it.StatusCode,
			&it.PasswordModifiedAt,
			&it.BadgeCode,
		)
		if err != nil {
			return nil, err
		}

		if it.StatusCode == 2 {
			return nil, http_error.NewUnexpectedError(http_error.UserDeleted)
		}

		if !it.IsActive {
			return nil, http_error.NewUnexpectedError(http_error.UserInactivated)
		}

		return &it, nil
	}
	return nil, http_error.NewUnexpectedError(http_error.UserNotFound)
}

func (r *authRepository) GetUserByID(
	ctx context.Context,
	ID int64,
) (*entities.User, error) {
	// language=sql
	query := `
	SELECT au.id, 
	       au.user_name,      
	       au.user_picture,      
	       au.user_document,
	       au.user_email,
	       ut.numeric_code,
	       au.is_active,
	       au.status_code,
	       au.password_modified_at,
	       COALESCE(au.badge_code, 0)
	  FROM public.user au
		INNER JOIN public.user_type aut on au.id = aut.id_user
		INNER JOIN public.type_user ut on aut.id_user_type = ut.id
	 WHERE au.id = $1`

	rows, err := r.conn.QueryContext(
		ctx,
		query,
		ID,
	)
	if err != nil {
		return nil, http_error.NewBadRequestError(http_error.Unexpected)
	}
	defer rows.Close()

	for rows.Next() {
		var it entities.User
		err = rows.Scan(
			&it.ID,
			&it.Name,
			&it.Image,
			&it.Document,
			&it.Email,
			&it.UserType,
			&it.IsActive,
			&it.StatusCode,
			&it.PasswordModifiedAt,
			&it.BadgeCode,
		)
		if err != nil {
			return nil, err
		}

		if it.StatusCode == 2 {
			return nil, http_error.NewUnexpectedError(http_error.UserDeleted)
		}

		if !it.IsActive {
			return nil, http_error.NewUnexpectedError(http_error.UserInactivated)
		}

		return &it, nil
	}
	return nil, http_error.NewUnexpectedError(http_error.UserNotFound)
}

func (r *authRepository) EmailExists(
	ctx context.Context,
	user entities.User,
) (bool, error) {
	// language=sql
	query := `
	SELECT user_email 
    FROM public.user
    WHERE user_email = $1 and status_code = 0;`

	stmt, err := r.conn.PrepareContext(ctx, query)
	if err != nil {
		return false, http_error.NewBadRequestError(http_error.Unexpected)
	}
	defer stmt.Close()

	row, err := stmt.QueryContext(ctx, user.Email)
	if err != nil {
		return false, http_error.NewBadRequestError(http_error.Unexpected)
	}
	defer row.Close()

	return row.Next(), nil
}

func (r *authRepository) DocumentExists(
	ctx context.Context,
	user entities.User,
) (bool, error) {
	// language=sql
	query := `
	SELECT user_document 
	FROM public.user
	WHERE user_document = $1 
	  AND status_code = 0;`

	stmt, err := r.conn.PrepareContext(ctx, query)
	if err != nil {
		log.Println("[DocumentExists] Error")
		err := http_error.NewBadRequestError(http_error.Unexpected)
		return false, err
	}
	defer stmt.Close()

	row, err := stmt.QueryContext(ctx, user.Document)
	if err != nil {
		log.Println("[DocumentExists] Error")
		return false, http_error.NewBadRequestError(http_error.Unexpected)
	}
	defer row.Close()

	return row.Next(), nil
}

func (r *authRepository) RegisterUser(ctx context.Context, user entities.User) error {
	// language=sql
	query := `
	INSERT INTO public.user(
		user_name,
	    user_picture,
		user_document,
		user_email, 
		zip_code,
		federal_unit,
		city,
		street,
		district,
		address_number,
		password_hash,
		is_active,
		terms_accepted
	) 
	VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, 1, (SELECT id FROM public.terms ORDER BY created_at DESC LIMIT 1))
	RETURNING id`

	tx, err := r.conn.Begin()
	if err != nil {
		return err
	}

	stmt, err := tx.PrepareContext(ctx, query)
	if err != nil {
		_ = tx.Rollback()
		return err
	}
	defer stmt.Close()

	// Gerar o hash da senha
	sb, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		_ = tx.Rollback()
		return err
	}
	user.Password = string(sb)

	log.Println("[RegisterUser] Password hash generated:", user.Password)

	// Executar a inserção e obter o ID do usuário criado
	var userID int64
	err = stmt.QueryRowContext(
		ctx,
		user.Name,
		user.ImageB64,
		user.Document,
		user.Email,
		user.Address.ZipCode,
		user.Address.StateID,
		user.Address.City,
		user.Address.Street,
		user.Address.District,
		user.Address.AddressNumber,
		user.Password,
	).Scan(&userID)
	if err != nil {
		_ = tx.Rollback()
		return err
	}

	// Chamar a função UserType passando o ID do usuário criado
	err = r.UserType(ctx, userID, tx)
	if err != nil {
		_ = tx.Rollback()
		return err
	}

	// Commit da transação
	return tx.Commit()
}
func (r *authRepository) UserType(ctx context.Context, id int64, tx *sql.Tx) error {
	//language=sql
	query := `
	INSERT INTO public.user_type (id_user, id_user_type) 
	VALUES ($1, $2)`

	stmt, err := tx.PrepareContext(ctx, query)
	if err != nil {
		_ = tx.Rollback()
		return err
	}
	defer stmt.Close()

	_, err = stmt.ExecContext(ctx, id, 1)
	if err != nil {
		_ = tx.Rollback()
		return err
	}

	return nil
}
