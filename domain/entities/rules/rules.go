package rules

import (
	"camera/domain/entities"
	"camera/infrastructure/modules/impl/http_error"
	"camera/utils"
	"log"
	"strings"
)

const (
	lowerChars   = `abcdefghijklmnopqrstuvwxyz`
	upperChars   = `ABCDEFGHIJKLMNOPQRSTUVWXYZ`
	sepChars     = `-_.@`
	numberChars  = `0123456789`
	specialChars = `!@#$%^&*`
)

func ValidateUserRegister(user *entities.User) error {
	user.Name = strings.TrimSpace(user.Name)
	user.Email = strings.TrimSpace(user.Email)

	if user.Name == "" {
		return http_error.NewBadRequestError(http_error.EmptyNameFieldError)
	}

	if user.Email == "" {
		return http_error.NewBadRequestError(http_error.EmptyEmailFieldError)
	}

	if user.IsForeigner == true {
		if *user.Document != "" {
			cpfIsValid := utils.CheckPersonDocument(*user.Document)
			if !cpfIsValid {
				return http_error.NewBadRequestError(http_error.DocumentNotValid)
			}
		}
	} else {
		if *user.Document == "" {
			return http_error.NewBadRequestError(http_error.EmptyDocumentFieldError)
		}

		cpfIsValid := utils.CheckPersonDocument(*user.Document)
		if !cpfIsValid {
			return http_error.NewBadRequestError(http_error.DocumentNotValid)
		}
	}

	if !checkValidEmail(user.Email) {
		return http_error.NewBadRequestError(http_error.InvalidEmailFieldError)
	}

	if user.Password == "" {
		return http_error.NewBadRequestError(http_error.EmptyPasswordField)
	}

	if strings.Contains(user.Password, " ") {
		return http_error.NewBadRequestError(http_error.PasswordIsntValid)
	}

	if user.Password != user.PasswordConfirmation {
		return http_error.NewBadRequestError(http_error.PasswordDoesntMatch)
	}

	isValid := CheckValidPassword(user.Password)
	if !isValid {
		return http_error.NewBadRequestError(http_error.PasswordIsntValid)
	}

	return nil
}

// checkValidEmail - Return true when email is valid.
func checkValidEmail(email string) bool {
	emailFormatted := strings.ToLower(strings.TrimSpace(email))

	for _, it := range emailFormatted {
		if !strings.ContainsRune(lowerChars, it) && !strings.ContainsRune(sepChars, it) && !strings.ContainsRune(numberChars, it) {
			log.Println("[checkValidEmail] Error invalid character special or any character special on start of string", email)
			return false
		}
	}

	sliceEmailAt := strings.Split(emailFormatted, "@")
	if len(sliceEmailAt) > 2 || len(sliceEmailAt) == 1 {
		log.Println("[checkValidEmail] Error more than one @ or without @", email)
		return false
	}

	if !strings.ContainsRune(sliceEmailAt[1], '.') {
		log.Println("[checkValidEmail] Error doesn't have . after @", email)
		return false
	}

	sliceEmailPoint := strings.Split(sliceEmailAt[1], ".")
	for _, it := range sliceEmailPoint {
		if it == "" {
			log.Println("[checkValidEmail] Error doesn't have string after @ or .", email)
			return false
		}
	}

	return true
}

func CheckValidPassword(password string) bool {
	var hasNumberChar bool
	var hasLowerChar bool
	var hasSpecialChar bool
	var hasUpperChar bool

	if len(password) < 8 {
		return false
	}

	if strings.Contains(password, " ") {
		return false
	}

	for _, it := range password {
		if strings.ContainsRune(numberChars, it) && !hasNumberChar {
			hasNumberChar = true
		}

		if strings.ContainsRune(specialChars, it) && !hasSpecialChar {
			hasSpecialChar = true
		}

		if strings.ContainsRune(upperChars, it) && !hasUpperChar {
			hasUpperChar = true
		}

		if strings.ContainsRune(lowerChars, it) && !hasLowerChar {
			hasLowerChar = true
		}
	}

	isValid := hasLowerChar && hasNumberChar
	isValid = hasSpecialChar && hasUpperChar

	return isValid
}

func ProductRules(product *entities.Product) error {
	product.ImageBase64 = strings.TrimSpace(product.ImageBase64)
	product.Name = strings.TrimSpace(product.Name)

	if product.ImageBase64 == "" {
		return http_error.NewBadRequestError(http_error.EmptyImageError)
	}

	if product.Name == "" {
		return http_error.NewBadRequestError(http_error.EmptyProductFieldError)
	}

	return nil
}

func ProductRulesEdite(product *entities.Product) error {
	product.Name = strings.TrimSpace(product.Name)

	if product.Name == "" {
		return http_error.NewBadRequestError(http_error.EmptyProductFieldError)
	}

	return nil
}
