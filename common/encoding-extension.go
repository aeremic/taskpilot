package common

import (
	"encoding/json"
	"os"
)

func GetAndDecodeJsonFile[T any](path string) (*T, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	reader := file
	decoder := json.NewDecoder(reader)

	var t T
	decoderErr := decoder.Decode(&t)
	if decoderErr != nil {
		return nil, decoderErr
	}

	return &t, nil
}
