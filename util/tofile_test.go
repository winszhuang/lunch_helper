package util

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"testing"
)

func TestSaveDataToFile(t *testing.T) {
	// 建立測試用的資料
	data := struct {
		Name string `json:"name"`
		Age  int    `json:"age"`
	}{
		Name: "John",
		Age:  30,
	}

	// 指定測試用的檔案路徑
	filePath := "test_data.json"

	// 確保測試完成後刪除測試用的檔案
	defer func() {
		err := os.Remove(filePath)
		if err != nil {
			t.Errorf("failed to remove test file: %v", err)
		}
	}()

	// 呼叫函數進行測試
	err := SaveDataToFile(data, filePath)
	if err != nil {
		t.Errorf("SaveDataToFile failed: %v", err)
	}

	// 讀取儲存的檔案內容
	fileContent, err := ioutil.ReadFile(filePath)
	if err != nil {
		t.Errorf("failed to read test file: %v", err)
	}

	// 解析檔案內容為 JSON
	var savedData map[string]interface{}
	err = json.Unmarshal(fileContent, &savedData)
	if err != nil {
		t.Errorf("failed to unmarshal test file content: %v", err)
	}

	// 驗證儲存的資料與原始資料是否一致
	expectedData := map[string]interface{}{
		"name": "John",
		"age":  30,
	}
	expectedJSON, _ := json.Marshal(expectedData)
	savedJSON, _ := json.Marshal(savedData)
	if string(savedJSON) != string(expectedJSON) {
		t.Errorf("saved data does not match expected data. Got %s, want %s", savedJSON, expectedJSON)
	}
}
