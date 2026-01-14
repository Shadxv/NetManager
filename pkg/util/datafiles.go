package util

import (
	"encoding/gob"
	"os"
	"strings"
)

func SaveData[T any](serviceName string, data T) error {
	f, err := os.Create("data/" + strings.ToLower(serviceName) + ".dat")
	if err != nil {
		return err
	}
	defer f.Close()

	enc := gob.NewEncoder(f)
	return enc.Encode(data)
}

func LoadData[T any](serviceName string) (T, error) {
	var result T

	f, err := os.Open("data/" + strings.ToLower(serviceName) + ".dat")
	if err != nil {
		return result, err
	}
	defer f.Close()

	dec := gob.NewDecoder(f)
	err = dec.Decode(&result)
	return result, err
}
