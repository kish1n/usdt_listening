package helpers

import (
	"github.com/go-chi/chi"
	"github.com/google/uuid"
	"github.com/pkg/errors"
	"io/ioutil"
	"log"
	"net/http"
)

func GenerateUUID() string {
	id, err := uuid.NewUUID()
	if err != nil {
		log.Fatalf("Failed to generate UUID: %v", err)
	}
	return id.String()
}

func ReadABIFile(filePath string) (string, error) {
	data, err := ioutil.ReadFile(filePath)
	if err != nil {
		return "", err
	}
	return string(data), nil
}

func GetAddress(r *http.Request, parameter string) (string, error) {
	shortened := chi.URLParam(r, parameter)
	if shortened == "" {
		return "", errors.New("from_address parameter is required")
	}
	return shortened, nil
}
