package http_error

import (
	"encoding/json"
	"log"
	"net/http"
)

func HandleError(w http.ResponseWriter, err error) {
	var statusCode int
	switch err.(type) {
	case BadRequestError:
		statusCode = http.StatusBadRequest
	case NotFoundError:
		statusCode = http.StatusNotFound
	case UnexpectedError:
		statusCode = http.StatusInternalServerError
	case UnauthorizedError:
		statusCode = http.StatusUnauthorized
	default:
		statusCode = 500
	}

	w.WriteHeader(statusCode)
	b, err := json.Marshal(err)
	if err != nil {
		log.Println("[HandleError] Error Marshal", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	_, err = w.Write(b)
	if err != nil {
		log.Println("[HandleError] Error Write", err)
		return
	}
}
