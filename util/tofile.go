package util

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
)

func SaveDataToFile(data interface{}, filePath string) error {
	fileContent, err := json.Marshal(data)
	if err != nil {
		return fmt.Errorf("failed to marshal data: %w", err)
	}

	err = ioutil.WriteFile(filePath, fileContent, 0644)
	if err != nil {
		return fmt.Errorf("failed to write file: %w", err)
	}

	fmt.Println("Data saved to file:", filePath)
	return nil
}
