package setup

import (
	"bear/domain/usecases/usecase_impl"
	"bear/infrastructure/modules/impl/auth"
	"bear/infrastructure/modules/impl/images"
	product "bear/infrastructure/modules/impl/product"
	"bear/infrastructure/repositories"
	"bear/infrastructure/repositories/impl"
	"bear/infrastructure/storage/sto"
	"bear/settings_loader"
	"context"
	"database/sql"
	"github.com/gorilla/mux"
	"github.com/gorilla/securecookie"
	_ "github.com/lib/pq" // Driver do PostgreSQL
	"log"
	"net/http"
	"os"
	"strconv"
	"time"
)

// DB é a variável global para a conexão ao banco
type SetupConfig struct {
	DB *sql.DB
}

// NewDatabaseConnection cria uma nova conexão com o banco de dados
func NewDatabaseConnection(settings *settings_loader.SettingsLoader) (*sql.DB, error) {
	dbConfig := settings.GetDatabaseConfig()
	if dbConfig.DatabaseURL == "" {
		log.Fatal("DATABASE_URL não está configurado")
	}

	db, err := sql.Open("postgres", dbConfig.DatabaseURL)
	if err != nil {
		return nil, err
	}

	// Configurações do pool de conexões
	db.SetMaxOpenConns(10)
	db.SetMaxIdleConns(5)
	db.SetConnMaxLifetime(10 * time.Minute)

	// Testa a conexão
	err = db.Ping()
	if err != nil {
		return nil, err
	}

	log.Println("Conectado ao banco de dados PostgreSQL com sucesso")
	return db, nil
}

// Setup inicializa o servidor e os módulos da aplicação
func Setup(router *mux.Router, settings *settings_loader.SettingsLoader) (*SetupConfig, error) {
	// Configurar a conexão ao banco de dados
	db, err := NewDatabaseConnection(settings)
	if err != nil {
		return nil, err
	}

	// Inicializar o módulo de autenticação
	SetupAuthModule(router, db, settings)

	// Configurar rotas privadas
	SetupPrivateRoutes(router, db, settings)

	SetupFileModule(router, settings)

	return &SetupConfig{DB: db}, nil
}

// SetupAuthModule configura as rotas de autenticação
func SetupAuthModule(router *mux.Router, db *sql.DB, settings *settings_loader.SettingsLoader) {
	fileStorage := sto.NewSTOManagerNew(*settings)
	// Inicialize o repositório de autenticação
	authRepo := impl.NewAuthenticationRepository(db, *settings)

	// Inicialize o caso de uso de autenticação
	authUseCase := usecase_impl.NewAuthenticationUseCase(authRepo, *settings, fileStorage)

	// Inicialize o módulo de autenticação
	authModule := &auth.AuthModule{
		Db:          db,
		Cookie:      securecookie.New([]byte(settings.SecurityConfig.CookieEncryptionKey), nil),
		AuthUseCase: authUseCase,
	}

	// Configure as rotas no roteador
	authModule.Setup(router)
}

// SetupPrivateRoutes configura as rotas privadas protegidas por middleware
func SetupPrivateRoutes(router *mux.Router, db *sql.DB, settings *settings_loader.SettingsLoader) {
	fileStorage := sto.NewSTOManagerNew(*settings)
	privateRouter := router.PathPrefix("/private").Subrouter()

	userRepo := impl.NewAuthenticationRepository(db, *settings) // Inicialize o repositório de usuário
	privateRouter.Use(AuthorizationMiddleware(userRepo, settings))

	productRepo := impl.NewProductRepository(db, *settings)
	productUseCase := usecase_impl.NewProductUseCase(productRepo, *settings, fileStorage)
	productModule := &product.ProductModule{
		Db:             db,
		Cookie:         securecookie.New([]byte(os.Getenv("COOKIE_ENCRYPTION_KEY")), nil),
		ProductUseCase: productUseCase,
	}

	productModule.Setup(privateRouter)

	privateRouter.HandleFunc("/profile", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Perfil do usuário - acesso autorizado"))
	}).Methods("GET")
}

func SetupFileModule(router *mux.Router, settings *settings_loader.SettingsLoader) {
	// Criar uma instância do módulo de arquivos
	fileModule := images.NewHTTPFileModule(settings) // Removido o operador &

	// Configurar as rotas relacionadas a arquivos
	fileModule.Setup(nil, router) // Passamos o roteador principal
}

const CtxUserKey = "auth-ctx-user-data"

// AuthorizationMiddleware é um middleware para proteger rotas privadas
func AuthorizationMiddleware(userRepo repositories.AuthenticationRepository, settings *settings_loader.SettingsLoader) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			cookie, err := r.Cookie("auth_token")
			if err != nil {
				log.Println("[AuthorizationMiddleware] Cookie não encontrado:", err)
				http.Error(w, "Não autorizado", http.StatusUnauthorized)
				return
			}

			var tokenData map[string]string
			secureCookie := securecookie.New([]byte(settings.SecurityConfig.CookieEncryptionKey), nil)
			if err := secureCookie.Decode("auth_token", cookie.Value, &tokenData); err != nil {
				log.Println("[AuthorizationMiddleware] Erro ao decodificar cookie:", err)
				http.Error(w, "Token inválido", http.StatusUnauthorized)
				return
			}

			userID := tokenData["user_id"]
			if userID == "" {
				log.Println("[AuthorizationMiddleware] user_id vazio ou não encontrado")
				http.Error(w, "Não autorizado", http.StatusUnauthorized)
				return
			}

			userIDInt, err := strconv.ParseInt(userID, 10, 64)
			if err != nil {
				log.Println("[AuthorizationMiddleware] Erro ao converter userID para int64:", err)
				http.Error(w, "Não autorizado", http.StatusUnauthorized)
				return
			}
			user, err := userRepo.GetUserByID(r.Context(), userIDInt)
			if err != nil || user == nil {
				log.Println("[AuthorizationMiddleware] Usuário não encontrado ou erro:", err)
				http.Error(w, "Não autorizado", http.StatusUnauthorized)
				return
			}

			// Adicione o usuário completo ao contexto
			ctx := context.WithValue(r.Context(), CtxUserKey, user)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

// CloseDB fecha a conexão com o banco ao encerrar a aplicação
func (c *SetupConfig) CloseDB() {
	if c.DB != nil {
		c.DB.Close()
	}
}
