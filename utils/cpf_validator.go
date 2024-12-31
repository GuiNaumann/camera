package utils

import (
	"log"
	"strconv"
)

var _personWeights = []int{11, 10, 9, 8, 7, 6, 5, 4, 3, 2}

// CheckPersonDocument - Check if it's a valid personal document
func CheckPersonDocument(document string) bool {
	if document == "" || len(document) != 11 {
		// Making sure the length is valid
		log.Println("falls here")
		return false
	}

	// Cutting the main part of the document
	baseDocument := document[:len(document)-2]
	// Calculating expected verification numbers
	ver1 := _calcVerificationNumber(baseDocument, _personWeights)
	ver2 := _calcVerificationNumber(baseDocument+strconv.Itoa(ver1), _personWeights)
	// Assembling the generated valid document
	generated := baseDocument + strconv.Itoa(ver1) + strconv.Itoa(ver2)
	// Checking if the document received is the same as the generated one
	return document == generated
}

// Calculate the verification number of a document
func _calcVerificationNumber(document string, weights []int) int {
	var sum = 0
	for i := len(document) - 1; i >= 0; i-- {
		number, _ := strconv.Atoi(document[i : i+1])
		sum += number * weights[len(weights)-len(document)+i]
	}

	check := 11 - sum%11
	if check > 9 {
		return 0
	} else {
		return check
	}
}
