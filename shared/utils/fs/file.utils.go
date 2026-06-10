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

func CreateFile(fileName string, path string) (string, error) {
	ext := filepath.Ext(fileName)

	nameWithoutExt := fileName[:len(fileName)-len(ext)]

	fileDir := filepath.Join(path, nameWithoutExt)
	if err := os.MkdirAll(fileDir, os.ModePerm); err != nil {
		return "", err
	}

	filePath := filepath.Join(fileDir, fileName)

	if _, err := os.Stat(filePath); err == nil {
		return filePath, nil
	}

	file, err := os.Create(filePath)
	if err != nil {
		return "", err
	}
	defer file.Close()
	return filePath, nil
}
