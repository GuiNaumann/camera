package images

import (
	"bear/infrastructure/modules/impl/http_error"
	"bear/settings_loader"
	"fmt"
	"github.com/gorilla/mux"
	"io"
	"log"
	"mime"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

type httpFileModule struct {
	settings *settings_loader.SettingsLoader
}

func NewHTTPFileModule(settings *settings_loader.SettingsLoader) *httpFileModule {
	return &httpFileModule{settings: settings}
}

const (
	ProductImagesFolder = "images/products"
)

const (
	RootFolder                 StorageFolder = ""
	ImagesFolder               StorageFolder = "images"
	PDFFolder                  StorageFolder = "pdf"
	ExcelFolder                StorageFolder = "xlsx"
	SCORMFolder                StorageFolder = "scorm"
	ComplementaryContentFolder StorageFolder = "complementary_contents"
	MP3Folder                  StorageFolder = "mp3"
)

type StorageFolder string

func (h httpFileModule) Setup(_ *mux.Router, router *mux.Router) {
	// Outras rotas já existentes
	router.HandleFunc("/products/{filename}", h.handleFile(ProductImagesFolder))
}

func (h httpFileModule) handleFile(folderName StorageFolder) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r) // Captura variáveis da URL, como {filename}
		filename := vars["filename"]
		filePath := filepath.Join(string(folderName), filename)

		b, err := h.GetFileBytes(filePath)
		if err != nil {
			log.Println("[handleFile] Error GetFileBytes", err)
			http_error.HandleError(w, err)
			return
		}

		mimeType := mime.TypeByExtension(filepath.Ext(filename))
		if mimeType == "" {
			mimeType = "application/octet-stream"
		}

		w.Header().Set("Content-Type", mimeType)
		w.WriteHeader(http.StatusOK)
		_, err = w.Write(b)
		if err != nil {
			log.Println("[handleFile] Error Write", err)
		}
	}
}

func (h httpFileModule) GetFileBytes(relativePath string) ([]byte, error) {
	useRelativePath := strings.Replace(filepath.Clean(relativePath), "\\", "/", -1)
	split := strings.Split(useRelativePath, "/")
	folderName := split[0]

	log.Println("")
	log.Println("", relativePath)
	log.Println("")
	log.Println("")
	log.Println("")
	log.Println("", useRelativePath)
	log.Println("")
	log.Println("", split)
	log.Println("")
	log.Println("")
	log.Println("AQUIIIIIII", folderName)
	log.Println("")
	log.Println("")
	log.Println("")

	switch StorageFolder(folderName) {
	case ImagesFolder,
		PDFFolder,
		SCORMFolder,
		ExcelFolder,
		MP3Folder,
		ComplementaryContentFolder:
		return h.GetBytes(useRelativePath)
	default:
		return nil, http_error.NewBadRequestError(fmt.Sprintf("Caminho %s não entrado (%s)", folderName, useRelativePath))
	}
}

func (h httpFileModule) GetBytes(relativePath string) ([]byte, error) {
	pathConfig := h.settings.GetPathConfig()

	filePath := relativePath

	filePath = filepath.Clean(filePath)
	pathConfigs := filepath.Clean(pathConfig.FileServerRootPath)

	hasRootPrefix := strings.HasPrefix(filePath, pathConfigs)
	if !hasRootPrefix {
		filePath = filepath.Join(pathConfig.FileServerRootPath, relativePath)
	}

	log.Println("")
	log.Println("")
	log.Println("")
	log.Println("")
	log.Println("filePath", filePath)
	log.Println("")
	log.Println("")
	log.Println("")
	log.Println("")

	file, err := os.Open(filePath)
	if err != nil {
		log.Println("[GetBytes] Error Open", err)
		return nil, err
	}
	defer file.Close()

	b, err := io.ReadAll(file)
	if err != nil {
		log.Println("[GetBytes] Error ReadAll(file)", err)
		return nil, err
	}

	return b, nil
}
