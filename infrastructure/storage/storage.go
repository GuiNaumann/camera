package storage

import (
	"io"
	"os"
	"strings"
)

type FileStorageRepository interface {
	SaveImage(imageFile string, pathDomain string, folderPath string, fileName string) (string, error)

	SetFullPath(fullPath string)

	IsURL(fileB64 string) bool

	TypeOfBase64(base64 string, prefix string) string

	CreateFileFromBase64(base64, domainPath string, fullFolderPath, folderPath, filename string) (string, error)
}

type FileStorageRepositoryNew interface {
	SaveBase64(fileBase64 string, filePath string) (string, error)

	CreateFile(
		path string,
		reader io.Reader,
		fileName string,
	) error

	DeletePath(fileName string) error

	CreateTempFile(body io.ReadCloser, extension string) (string, error)

	CreateEmptyFolder(path string) (string, error)

	CreateFolderIfNotExists(folder string) error

	CheckPathPrefix(dir string) (string, error)

	Open(filePath string) (*os.File, error)

	OpenXLSXImport(filePath string) (*os.File, error)

	GetBytes(filePath string) ([]byte, error)

	Store(filePath string, b []byte) (string, error)
}

func IsURL(fileString string) bool {
	indexByte := strings.IndexByte(fileString, ':')
	if indexByte == -1 {
		return false
	}

	var indexBase64 string
	indexBase64 = fileString[0:indexByte]

	return indexBase64 == "http" || indexBase64 == "https"
}
