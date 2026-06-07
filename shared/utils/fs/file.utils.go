package fs

import (
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
)

func ReadFile[T any](path string) (T, error) {
	var response T
	data, err := os.ReadFile(path)

	if err != nil {
		return response, nil
	}
	err = json.Unmarshal(data, &response)

	if err != nil {
		return response, err
	}

	return response, nil
}

func WriteJSONFile(path string, fileName string, data any, permCode int32) error {

	_, err := os.Stat(path)

	if err != nil {
		return errors.New("Unable to find the directory " + path)
	}
	fileData, err := json.Marshal(data)

	if err != nil {
		return errors.New("Can't Marshal data  :: " + err.Error())
	}
	writeFilePath := filepath.Join(path, fileName)

	err = os.WriteFile(writeFilePath, fileData, os.FileMode(permCode))

	if err != nil {
		return errors.New("Unable to write JSON file :: " + err.Error())
	}

	return nil

}
