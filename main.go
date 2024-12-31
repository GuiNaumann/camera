package main

import (
	"github.com/joho/godotenv"
	"log"
	"net/http"
	"os"

	"bear/infrastructure"
	"bear/settings_loader"
	"github.com/gorilla/mux"
	"github.com/rs/cors"
)

func init() {
	dir, _ := os.Getwd()
	log.Printf("Diretório atual: %s", dir)

	// Carrega as variáveis do .env
	err := godotenv.Load()
	if err != nil {
		log.Println("Erro ao carregar .env")
	}
}

func initLogger() {
	// Crie a pasta de logs, se ela não existir
	logDir := "./logs"
	if _, err := os.Stat(logDir); os.IsNotExist(err) {
		err := os.Mkdir(logDir, 0755)
		if err != nil {
			log.Fatalf("Erro ao criar a pasta de logs: %v", err)
		}
	}

	// Abra o arquivo de log para escrita (modo append)
	logFile, err := os.OpenFile(logDir+"/server.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatalf("Erro ao abrir/criar arquivo de log: %v", err)
	}

	// Configurar o log para gravar no arquivo e no console
	log.SetOutput(logFile)
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)
}

func main() {
	router := mux.NewRouter()
	initLogger()
	log.Println("Servidor iniciado com sucesso!")

	// Configurar o middleware CORS
	handler := cors.New(cors.Options{
		AllowedOrigins:   []string{"http://localhost:8081"}, // Frontend URL
		AllowedMethods:   []string{"POST", "GET", "OPTIONS"},
		AllowedHeaders:   []string{"Content-Type", "Authorization"},
		AllowCredentials: true, // Habilita envio de cookies
	}).Handler(router)

	// Carregar as configurações
	settings := settings_loader.NewSettingsLoader()

	// Configura o projeto chamando o Setup da infraestrutura
	setupConfig, err := setup.Setup(router, settings)
	if err != nil {
		log.Fatalf("Erro ao configurar a infraestrutura: %v", err)
	}

	defer setupConfig.CloseDB()

	// Inicia o servidor usando o handler com CORS
	log.Println("Servidor iniciado na porta 8080")
	log.Fatal(http.ListenAndServe(":8080", handler)) // Use handler aqui
}
