package settings_loader

import (
	"fmt"
	"github.com/pelletier/go-toml"
	"log"
)

// SecurityConfig armazena configurações relacionadas à segurança (chave de criptografia do cookie e segredo do JWT)
type SecurityConfig struct {
	CookieEncryptionKey string
	JWTSecret           string
}

// DatabaseConfig armazena a URL do banco de dados
type DatabaseConfig struct {
	DatabaseURL string
}

// PathConfig contém os caminhos utilizados pela aplicação
type PathConfig struct {
	LogoPath            string
	FavIconPath         string
	EmailImagesRootPath string
	FileServerRootPath  string
	RedirectUrl         string
	HTMLRootPath        string
}

// TLSConfig define as configurações de TLS
type TLSConfig struct {
	IsTLS bool
	Cert  string
	Key   string
}

// ServerDomainConfig define configurações de domínio
type ServerDomainConfig struct {
	ServerDomain   string
	UseCloudSql    string
	EnableRedirect bool
}

// SettingsLoader carrega todas as configurações do ambiente
type SettingsLoader struct {
	SecurityConfig     SecurityConfig
	DatabaseConfig     DatabaseConfig
	PathConfig         PathConfig
	TLSConfig          TLSConfig
	ServerDomainConfig ServerDomainConfig
}

// NewSettingsLoader cria uma nova instância do SettingsLoader e carrega as configurações do ambiente
func NewSettingsLoader() *SettingsLoader {
	// Carregar o arquivo TOML
	config, err := toml.LoadFile("settings.toml")
	if err != nil {
		log.Fatalf("Erro ao carregar settings.toml: %v", err)
	}
	log.Printf("Arquivo settings.toml carregado com sucesso.")

	// Preencher os valores
	pathConfig := PathConfig{
		LogoPath:            config.Get("PathConfig.LogoPath").(string),
		FavIconPath:         config.Get("PathConfig.FavIconPath").(string),
		EmailImagesRootPath: config.Get("PathConfig.EmailImagesRootPath").(string),
		FileServerRootPath:  config.Get("PathConfig.FileServerRootPath").(string),
		RedirectUrl:         config.Get("PathConfig.RedirectUrl").(string),
		HTMLRootPath:        config.Get("PathConfig.HTMLRootPath").(string),
	}

	return &SettingsLoader{
		SecurityConfig: SecurityConfig{
			CookieEncryptionKey: config.Get("SecurityConfig.COOKIE_ENCRYPTION_KEY").(string),
			JWTSecret:           config.Get("SecurityConfig.JWT_SECRET").(string),
		},
		DatabaseConfig: DatabaseConfig{
			DatabaseURL: config.Get("DatabaseConfig.DATABASE_URL").(string),
		},
		PathConfig: pathConfig,
		TLSConfig: TLSConfig{
			IsTLS: config.Get("TLSConfig.IsTLS").(bool),
			Cert:  config.Get("TLSConfig.Cert").(string),
			Key:   config.Get("TLSConfig.Key").(string),
		},
		ServerDomainConfig: ServerDomainConfig{
			ServerDomain:   config.Get("ServerDomainConfig.ServerDomain").(string),
			UseCloudSql:    config.Get("ServerDomainConfig.UseCloudSql").(string),
			EnableRedirect: config.Get("ServerDomainConfig.EnableRedirect").(bool),
		},
	}
}

// GetSecurityConfig retorna as configurações de segurança
func (s *SettingsLoader) GetSecurityConfig() SecurityConfig {
	return s.SecurityConfig
}

// GetDatabaseConfig retorna as configurações do banco de dados
func (s *SettingsLoader) GetDatabaseConfig() DatabaseConfig {
	return s.DatabaseConfig
}

// GetPathConfig retorna as configurações de caminhos
func (s *SettingsLoader) GetPathConfig() PathConfig {
	return s.PathConfig
}

// GetFullDomain retorna o domínio completo (com protocolo)
func (s *SettingsLoader) GetFullDomain() (string, error) {
	if s.TLSConfig.IsTLS {
		return fmt.Sprintf("https://%s", s.ServerDomainConfig.ServerDomain), nil
	} else {
		return fmt.Sprintf("http://%s", s.ServerDomainConfig.ServerDomain), nil
	}
}
