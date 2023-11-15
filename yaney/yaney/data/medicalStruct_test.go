package data

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"testing"
)

func TestMedicalStruct(t *testing.T) {
	var (
		err      error
		content  []byte
		data     MedicalData
		filename string
	)

	filename = "/home/l1nkkk/project/mime/yaney/data/data.json"

	// 1, 把配置文件读进来
	if content, err = ioutil.ReadFile(filename); err != nil {
		return
	}

	// 2, 做JSON反序列化
	if err = json.Unmarshal(content, &data); err != nil {
		return
	}
	fmt.Printf("%v",data)
}