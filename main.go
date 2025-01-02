package main

import (
	"camera/domain/entities"
	setup "camera/infrastructure"
	"camera/infrastructure/repositories"
	"camera/settings_loader"
	"context"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	"github.com/rs/cors"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"time"
)

const networkSharePath = `\\192.168.1.100\gravacoes` // Caminho do compartilhamento de rede

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

	// Inicia a gravação das câmeras em uma goroutine
	go startCameraRecording(setupConfig.ProductRepository)

	// Inicia o servidor usando o handler com CORS
	log.Println("Servidor iniciado na porta 8080")
	log.Fatal(http.ListenAndServe(":8080", handler)) // Use handler aqui
}

// Inicia o processo de gravação para todas as câmeras
func startCameraRecording(repo repositories.ProductRepository) {
	for {
		// Obtemos todas as câmeras ativas
		ctx := context.Background()
		filter := entities.GeneralFilter{
			Limit:  0,  // Sem limite
			Page:   0,  // Primeira página
			Search: "", // Sem filtro de busca
		}
		user := entities.User{ID: 1} // Exemplo: substitua pelo ID real

		cameras, err := repo.ListProductRepository(ctx, filter, user)
		if err != nil {
			log.Printf("Erro ao listar câmeras: %v", err)
			time.Sleep(1 * time.Second) // Tenta novamente após 1 minuto
			continue
		}

		// Grava cada câmera em uma goroutine
		for _, camera := range cameras.Items {
			go recordCamera(camera)
		}

		// Aguarda até o próximo dia para reiniciar o processo
		now := time.Now()
		nextDay := now.Add(24 * time.Hour).Truncate(24 * time.Hour)
		time.Sleep(time.Until(nextDay))
	}
}

// Grava o stream de uma câmera
func recordCamera(camera entities.Product) {
	// Cria o diretório no compartilhamento de rede
	dir := filepath.Join(networkSharePath, camera.Name)
	if err := os.MkdirAll(dir, os.ModePerm); err != nil {
		log.Printf("Erro ao criar diretório na rede %s: %v", dir, err)
		return
	}

	// Define o nome do arquivo
	now := time.Now()
	fileName := fmt.Sprintf("%s.mp4", now.Format("2006-01-02"))
	outputPath := filepath.Join(dir, fileName)

	// Calcula a duração restante do dia
	tomorrow := now.Add(24 * time.Hour).Truncate(24 * time.Hour)
	duration := time.Until(tomorrow).Seconds()

	// Comando FFmpeg para gravar
	cmd := exec.Command("ffmpeg",
		"-i", fmt.Sprintf("rtsp://%s:%s@%s:%d%s",
			camera.Username, camera.Password, camera.IPAddress, camera.Port, camera.StreamPath),
		"-t", fmt.Sprintf("%.0f", duration),
		"-c:v", "copy",
		"-c:a", "aac",
		outputPath)

	log.Printf("Iniciando gravação da câmera: %s", camera.Name)

	if err := cmd.Run(); err != nil {
		log.Printf("Erro ao gravar câmera %s: %v", camera.Name, err)
		return
	}

	log.Printf("Gravação concluída: %s", outputPath)

	// Limpa gravações antigas
	cleanupOldRecordings(dir)
}

// Remove gravações mais antigas que 30 dias
func cleanupOldRecordings(dir string) {
	files, err := os.ReadDir(dir)
	if err != nil {
		log.Printf("Erro ao listar arquivos em %s: %v", dir, err)
		return
	}

	limitDate := time.Now().AddDate(0, 0, -30)
	for _, file := range files {
		if file.IsDir() {
			continue
		}

		filePath := filepath.Join(dir, file.Name())
		info, err := os.Stat(filePath)
		if err != nil {
			log.Printf("Erro ao obter informações do arquivo %s: %v", filePath, err)
			continue
		}

		if info.ModTime().Before(limitDate) {
			if err := os.Remove(filePath); err != nil {
				log.Printf("Erro ao remover arquivo %s: %v", filePath, err)
			} else {
				log.Printf("Arquivo removido: %s", filePath)
			}
		}
	}
}
