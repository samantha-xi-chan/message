package util_struct

import (
	"encoding/json"
	"fmt"
	"log"
)

func ConvertJsonStr2Interface(input string) (x interface{}, err error) {
	var output map[string]interface{}
	if err = json.Unmarshal([]byte(input), &output); err != nil {
		fmt.Println("解码JSON失败:", err)
		return
	}

	return output, nil
}

func MultiConvertJsonStr2Interface(inputs []string) (x []interface{}, err error) {
	size := len(inputs)
	log.Println("MultiConvertJsonStr2Interface size: ", size)

	var outputs []interface{}

	for i := 0; i < size; i++ {
		input := inputs[i]
		var output map[string]interface{}
		if err = json.Unmarshal([]byte(input), &output); err != nil {
			log.Println("解码JSON失败:", err)
			return nil, err
		}
		outputs = append(outputs, output)
	}

	return outputs, nil
}
