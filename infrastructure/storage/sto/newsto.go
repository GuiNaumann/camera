package sto

import (
	"camera/infrastructure/modules/impl/http_error"
	"camera/infrastructure/storage"
	"camera/settings_loader"
	"camera/utils"
	"encoding/base64"
	"errors"
	"fmt"
	"image"
	"image/jpeg"
	"image/png"
	"io"
	"io/fs"
	"io/ioutil"
	"log"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"strings"
)

const (
	ImagesRelativePath = "/images"
)

func NewSTOManagerNew(settings settings_loader.SettingsLoader) storage.FileStorageRepositoryNew {
	return &newStorageManager{
		settings: settings,
	}
}

type newStorageManager struct {
	settings settings_loader.SettingsLoader
}

func (g *newStorageManager) Open(filePath string) (*os.File, error) {
	usePath, err := g.CheckPathPrefix(filePath)
	if err != nil {
		return nil, err
	}

	log.Printf("[Open] %s", usePath)

	return os.Open(usePath)
}

func (g *newStorageManager) OpenXLSXImport(filePath string) (*os.File, error) {
	usePath, err := g.CheckPathPrefixXLSXImport(filePath)
	if err != nil {
		return nil, err
	}

	log.Printf("[OpenXLSXImport] %s", usePath)

	return os.Open(usePath)
}

func (g *newStorageManager) Store(filePath string, b []byte) (string, error) {

	pathConfig := g.settings.GetPathConfig()

	absolutePath := filePath

	absolutePath = filepath.Clean(absolutePath)
	pathConfigs := filepath.Clean(pathConfig.FileServerRootPath)

	// Check if the filePath contains the storage absolute path and add if not contains.
	hasRootPrefix := strings.HasPrefix(absolutePath, pathConfigs)
	if !hasRootPrefix {
		absolutePath = filepath.Join(pathConfig.FileServerRootPath, filePath)
	}

	log.Println("pathConfig.FileServerRootPath", pathConfig.FileServerRootPath)
	log.Println("absolutePath", absolutePath)

	err := os.WriteFile(absolutePath, b, os.ModePerm)
	if err != nil {
		log.Println("[Store] Error WriteFile", err)
		return "", err
	}

	return absolutePath, nil
}

func (g *newStorageManager) GetBytes(relativePath string) ([]byte, error) {
	pathConfig := g.settings.GetPathConfig()

	filePath := relativePath

	filePath = filepath.Clean(filePath)
	pathConfigs := filepath.Clean(pathConfig.FileServerRootPath)

	hasRootPrefix := strings.HasPrefix(filePath, pathConfigs)
	if !hasRootPrefix {
		filePath = filepath.Join(pathConfig.FileServerRootPath, relativePath)
	}

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

func (g *newStorageManager) CheckPathPrefix(dir string) (string, error) {
	pathConfig := g.settings.GetPathConfig()
	dir = filepath.Clean(dir)
	pathConfigs := filepath.Clean(pathConfig.FileServerRootPath)

	hasRootPrefix := strings.HasPrefix(dir, pathConfigs)
	if hasRootPrefix {
		return dir, nil
	}

	result := filepath.Join(pathConfig.FileServerRootPath, dir)

	return result, nil
}

func (g *newStorageManager) CheckPathPrefixXLSXImport(dir string) (string, error) {
	pathConfig := g.settings.GetPathConfig()
	dir = filepath.Clean(dir)
	pathConfigs := filepath.Clean(pathConfig.FileServerRootPath)

	hasRootPrefix := strings.HasPrefix(dir, pathConfigs)
	if hasRootPrefix {
		return dir, nil
	}

	return dir, nil
}

// CreateFolderIfNotExists creates a folder if it does not exist.
func (g *newStorageManager) CreateFolderIfNotExists(folder string) error {
	useFolder, err := g.CheckPathPrefix(folder)
	if err != nil {
		log.Println("[CreateFolderIfNotExists] Error GetPathConfig", err)
		return err
	}

	if _, err := os.Stat(useFolder); os.IsNotExist(err) {
		mkdirErr := os.MkdirAll(useFolder, 0750)
		if mkdirErr != nil {
			log.Println("[CreateFolderIfNotExists] Error os.MkdirAll(folder, os.ModePerm)", err)
			return err
		}
	} else if err != nil {
		log.Println("[CreateFolderIfNotExists] Error os.Stat(folder)", err)
		return err
	}

	return nil
}

func (g *newStorageManager) CreateEmptyFolder(path string) (string, error) {
	pathConfig := g.settings.GetPathConfig()

	pathToCreate := filepath.Clean(path)

	// Clean file and storage path.
	storageAbsolutePath := filepath.Clean(pathConfig.FileServerRootPath)

	// Check if the filePath contains the storage absolute path and add if not contains.
	hasRootPrefix := strings.HasPrefix(pathToCreate, storageAbsolutePath)

	if !hasRootPrefix {
		pathToCreate = filepath.Join(storageAbsolutePath, pathToCreate)
	} else {
		log.Printf(
			"[CreateEmptyFolder] The relativePath (%s) should not contains the root path. Refactor it.",
			storageAbsolutePath,
		)
	}

	err := os.MkdirAll(pathToCreate, 0777)
	if err != nil {
		log.Println("[CreateEmptyFolder] Errors Mkdir")
		return "", err
	}

	return pathToCreate, nil
}

func (g *newStorageManager) DeletePath(relativePath string) error {
	// Check the filePath contains alphanumeric character.
	pathContainsAlphaNum, err := regexp.Match("([a-zA-Z]|[0-9])", []byte(relativePath))
	if err != nil {
		log.Println("[DeletePath] Error Match", err)
		return err
	}
	if !pathContainsAlphaNum {
		err = errors.New("file path is not valid.")
		log.Println("[DeletePath] Error !pathContainsAlphaNum", err)
		return err
	}

	pathConfig := g.settings.GetPathConfig()
	// Check the storage path is defined.
	storageAbsolutePath := strings.TrimSpace(pathConfig.FileServerRootPath)
	if storageAbsolutePath == "" {
		err = errors.New("storage path is not defined check the settings file.")
		log.Println("[DeletePath] Error empty storageAbsolutePath", err)
		return err
	}

	// Clean file and storage path.
	storageAbsolutePath = filepath.Clean(pathConfig.FileServerRootPath)
	relativePath = filepath.Clean(relativePath)

	// Check if the filePath contains the storage absolute path and add if not contains.
	var pathToRemove string
	filePathHasRootPath := strings.HasPrefix(relativePath, storageAbsolutePath)
	if !filePathHasRootPath {
		pathToRemove = filepath.Join(storageAbsolutePath, relativePath)
	} else {
		log.Printf("[DeletePath] The relativePath (%s) should not contains the root path. Refactor it.", relativePath)
	}

	// Remove the file.
	err = os.RemoveAll(pathToRemove)
	if err != nil {
		log.Println("[DeletePath] Error RemoveAll", err)
		return err
	}
	log.Printf("[DeletePath] PATH (%S), REMOVED SUCCESSFULLY", pathToRemove)

	return nil
}

func (g *newStorageManager) encodeImage(reader io.Reader, writer io.Writer, imgType string, img image.Image) error {
	var err error
	switch imgType {
	case "jpg", "jpeg":
		err = jpeg.Encode(writer, img, &jpeg.Options{Quality: 75})
		if err != nil {
			log.Println("[encodeImage] Error jpeg.Encode")
			return http_error.NewUnexpectedError(http_error.Unexpected)
		}

		break
	case "png":
		err = png.Encode(writer, img)
		if err != nil {
			log.Println("[encodeImage] Error png.Encode")
			return http_error.NewUnexpectedError(http_error.Unexpected)
		}
		break
	case "webp":
		buf, err := ioutil.ReadAll(reader)
		if err != nil {
			log.Println("[encodeImage] Error webp")
			return http_error.NewUnexpectedError(http_error.Unexpected)
		}

		_, err = writer.Write(buf)
		if err != nil {
			log.Println("[encodeImage] Error Write")
			return http_error.NewUnexpectedError(http_error.Unexpected)
		}
		break

	default:
		log.Println("imgType -> ", imgType)
		return http_error.NewUnexpectedError(http_error.InvalidImageExtension)
	}

	return nil
}

func (g *newStorageManager) extractDataFromB64(base64 string) string {
	return base64[strings.IndexByte(base64, ',')+1:]
}

func (g *newStorageManager) Base64HasMimeType(fileBase64 string) bool {
	return strings.Contains(fileBase64, ":") ||
		strings.Contains(fileBase64, ";") ||
		strings.Contains(fileBase64, ",")
}

func (g *newStorageManager) CreateFileFromBase64(
	fileBase64 string,
	domainPath string,
	fullFolderPath string,
	folderPath string,
	filename string,
) (string, error) {
	var err error

	if !g.Base64HasMimeType(fileBase64) {
		return "", http_error.NewUnexpectedError(http_error.InvalidMetadata)
	}

	base64WithoutHeader := g.extractDataFromB64(fileBase64)
	decodedBase64, err := base64.StdEncoding.DecodeString(base64WithoutHeader)
	if err != nil {
		log.Println("[CreateFileFromBase64] Error DecodeString")
		return "", http_error.NewUnexpectedError(http_error.InvalidMetadata)
	}

	fileExtension, err := g.TypeOfBase64(fileBase64)
	if err != nil {
		log.Println("[CreateFileFromBase64] Error TypeOfBase64", err)
		return "", err
	}

	filePath := fullFolderPath + "/" + filename + "." + fileExtension
	file, err := os.Create(filePath)
	if err != nil {
		log.Println("[CreateFileFromBase64] Error OpenFile")
		return "", http_error.NewUnexpectedError(http_error.Unexpected)
	}
	defer file.Close()

	_, err = file.Write(decodedBase64)
	if err != nil {
		log.Println("[CreateFileFromBase64] Error Write")
		return "", http_error.NewUnexpectedError(http_error.Unexpected)
	}

	err = file.Sync()
	if err != nil {
		log.Println("[CreateFileFromBase64] Error Sync")
		return "", http_error.NewUnexpectedError(http_error.Unexpected)
	}

	return domainPath + folderPath + "/" + filename + fileExtension, nil
}

func (g *newStorageManager) TypeOfBase64(base64 string) (string, error) {
	if strings.Contains(base64, "/") && strings.Contains(base64, ";") {
		return base64[strings.Index(base64, "/")+1 : strings.Index(base64, ";")], nil
	}

	return "", errors.New("Invalid base64")
}

func (g *newStorageManager) SaveBase64(fileBase64 string, filePath string) (string, error) {
	fileType, err := g.TypeOfBase64(fileBase64)
	if err != nil {
		log.Println("[SaveBase64] Error TypeOfBase64", err)
		return "", err
	}

	switch fileType {
	case "webp", "png", "jpg", "jpeg":
		return g.saveImage(fileBase64, filePath)
	case "pdf", "mp4":
		return g.saveFile(fileBase64, fileType, filePath)
	}

	return "", nil
}

func (g *newStorageManager) saveImage(imageBase64 string, filePath string) (string, error) {
	//In case that not is base64 just return, because can be a URL
	if storage.IsURL(imageBase64) {
		return imageBase64, nil
	}

	if !g.Base64HasMimeType(imageBase64) {
		return "", http_error.NewUnexpectedError(http_error.InvalidMetadata)
	}

	imgType, err := g.TypeOfBase64(imageBase64)
	if err != nil {
		log.Println("[saveImage] Error TypeOfBase64", err)
		return "", err
	}

	filePathType := filePath + "." + imgType

	pathConfig := g.settings.GetPathConfig()
	log.Printf("_________----------------________________-----_____--FileServerRootPath: %s", pathConfig.FileServerRootPath)

	fullFilePath := path.Join(pathConfig.FileServerRootPath, "images", filePathType)

	log.Println("||||||||||||||||||||||||||||||| : ", fullFilePath)
	// C:\dev\camera_run\storage/images/products/20193929-befc-11ef-86b0-00ff3d6dafd2.png

	err = g.DeletePath(filePath)
	if err != nil {
		log.Println("[saveImage] DeletePath", err)
		return "", http_error.NewUnexpectedError(http_error.Unexpected)
	}

	dataB64 := g.extractDataFromB64(imageBase64)
	reader := base64.NewDecoder(base64.StdEncoding, strings.NewReader(dataB64))

	var img image.Image

	if imgType != "webp" {
		img, _, err = image.Decode(reader)
		if err != nil {
			log.Println("[saveImage] Error Decode reader", err)
			return "", http_error.NewUnexpectedError(http_error.Unexpected)
		}
	}

	parentFolder, _ := filepath.Split(fullFilePath)

	err = os.MkdirAll(parentFolder, 0777)
	if err != nil {
		log.Println("[saveImage] error creating parent folder", err)
		return "", http_error.NewUnexpectedError(http_error.Unexpected)
	}

	log.Printf("File path: %s", parentFolder)
	stats, err := os.Stat(parentFolder)
	if err != nil {
		log.Println("[saveImage] Error getting dir stats", err)
		return "", http_error.NewUnexpectedError(http_error.Unexpected)
	}

	log.Printf("IsDir: %v, Mode: %v", stats.IsDir(), stats.Mode())

	writer, err := os.OpenFile(fullFilePath, os.O_WRONLY|os.O_CREATE|os.O_APPEND|os.O_TRUNC, 0777)
	if err != nil {
		log.Println("[saveImage] Error OpenFile", err)
		return "", http_error.NewUnexpectedError(http_error.Unexpected)
	}

	err = g.encodeImage(reader, writer, imgType, img)
	if err != nil {
		log.Println("[saveImage] Error writer encodeImage")
		return "", err
	}

	writer.Close()

	readerExif := base64.NewDecoder(base64.StdEncoding, strings.NewReader(dataB64))

	ori, err := utils.ReadOrientation(readerExif)
	if err != nil {
		log.Println("[saveImage] Error ReadOrientation")
		return "", http_error.NewUnexpectedError(http_error.Unexpected)
	}

	if ori != 0 {
		data, err := ioutil.ReadFile(fullFilePath)
		if err != nil {
			log.Println("[SaveImage] Error GetBytes")
			return "", http_error.NewUnexpectedError(http_error.Unexpected)
		}

		var reorientedImg image.Image

		switch ori {
		case 3:
			reorientedImg, err = utils.RotateImage(data, 180)
			if err != nil {
				log.Println("[SaveImage] Error 3 RotateImage")
				return "", http_error.NewUnexpectedError(http_error.Unexpected)
			}
			break
		case 6:
			reorientedImg, err = utils.RotateImage(data, 90)
			if err != nil {
				log.Println("[SaveImage] Error 6 RotateImage")
				return "", http_error.NewUnexpectedError(http_error.Unexpected)
			}
			break
		case 8:
			reorientedImg, err = utils.RotateImage(data, 270)
			if err != nil {
				log.Println("[SaveImage] Error 8 RotateImage")
				return "", http_error.NewUnexpectedError(http_error.Unexpected)
			}
			break
		default:
			reorientedImg = img
			break
		}

		auxWriter, err := os.OpenFile(fullFilePath, os.O_WRONLY|os.O_CREATE|os.O_APPEND|os.O_TRUNC, 0777)
		if err != nil {
			log.Println("[SaveImage] Error auxWriter OpenFile")
			return "", http_error.NewUnexpectedError(http_error.Unexpected)
		}
		defer auxWriter.Close()

		err = g.encodeImage(readerExif, auxWriter, imgType, reorientedImg)
		if err != nil {
			log.Println("[SaveImage] Error auxWriter encodeImage")
			return "", err
		}
	}

	getFullDomain, err := g.settings.GetFullDomain()
	if err != nil {
		log.Println("[SaveImage] Error GetFullDomain")
		return "", err
	}
	log.Println("ZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZ")
	log.Println("", getFullDomain)
	log.Println("ZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZ")
	log.Println(fmt.Sprintf("%s/%s", getFullDomain, filePathType))
	log.Println("ZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZ")
	return fmt.Sprintf("%s/%s", getFullDomain, filePathType), nil

}

func (g *newStorageManager) saveFile(
	fileBase64 string,
	fileType string,
	filePath string,
) (string, error) {
	decodedBase64, err := base64.StdEncoding.DecodeString(g.extractDataFromB64(fileBase64))
	if err != nil {
		log.Println("[CreateFileFromBase64] Error DecodeString")
		return "", http_error.NewUnexpectedError(http_error.InvalidMetadata)
	}

	pathConfig := g.settings.GetPathConfig()

	fullFilePath := path.Join(pathConfig.FileServerRootPath, filePath)
	err = g.DeletePath(fullFilePath)
	if err != nil {
		log.Println("[SaveFile] Error DeletePath: ", err.Error())
		return "", err
	}

	rootPath, fileName := path.Split(fullFilePath)
	fileName += fmt.Sprintf(".%s", fileType)
	filePath += fmt.Sprintf(".%s", fileType)

	// Verify if path exist, case not exist create
	_, err = os.Stat(rootPath)
	if err != nil {
		if os.IsNotExist(err) {
			if err != nil && !errors.Is(err, fs.ErrNotExist) {
				log.Println("[GeneratePDF] Error Stat", err)
				return "", err
			}

			errDir := os.Mkdir(rootPath, 0777)
			if errDir != nil {
				log.Println("[CreateFile] Error Mkdir", errDir)
				return "", errDir
			}
		}
		return "", err
	}

	err = os.WriteFile(filepath.Join(rootPath, fileName), decodedBase64, 0777)
	if err != nil {
		log.Println("[SaveFile] Error WriteFile: ", err.Error())
		return "", err
	}

	GetFullDomain, err := g.settings.GetFullDomain()
	if err != nil {
		log.Println("[SaveFile] Error GetFullDomain")
		return "", err
	}
	return GetFullDomain + filePath, nil
}

func (g *newStorageManager) CreateFile(
	path string,
	reader io.Reader,
	fileName string,
) error {
	file, err := io.ReadAll(reader)
	if err != nil {
		log.Println("[CreateFile] Error RemoveAll: ", err.Error())
		return err
	}

	// Verify if path exist, case not exist create
	if _, err = os.Stat(path); os.IsNotExist(err) {
		if err != nil && !errors.Is(err, fs.ErrNotExist) {
			log.Println("[CreateFile] Error Stat", err)
			return err
		}

		errDir := os.MkdirAll(path, 0777)
		if errDir != nil {
			log.Println("[CreateFile] Error Mkdir", errDir)
			return errDir
		}
	}

	err = os.WriteFile(filepath.Join(path, fileName), file, 0777)
	if err != nil {
		log.Println("[CreateFile] Error WriteFile", err.Error())
		return err
	}

	return nil
}

func (g *newStorageManager) CreateTempFile(body io.ReadCloser, extension string) (string, error) {
	// Create a temporary file to statement the required data
	tempFile, err := ioutil.TempFile("", fmt.Sprintf("*.%s", filepath.Base(extension)))
	if err != nil {
		return "", err
	}
	defer tempFile.Close()

	// read the contents of our uploaded file and dumps in the temp file
	_, err = io.Copy(tempFile, body)
	if err != nil {
		return "", err
	}

	tempFilePath := tempFile.Name()

	// Close the temporary file
	err = tempFile.Close()
	if err != nil {
		return "", err
	}

	return tempFilePath, nil
}
