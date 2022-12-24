package util

import (
	"encoding/json"
	"fmt"
)

func JSONStringify(obj interface{}) string {
	beautifulJsonByte, err := json.MarshalIndent(obj, "", "  ")
	body := ""
	if err != nil {
		body = fmt.Sprintf("%v", obj)
	} else {
		body = string(beautifulJsonByte)
	}
	return body
}

// this function prints any data type as json
func PrintJson(v interface{}) {
	beautifulJsonByte, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		fmt.Println(err)
		panic(err)
	}
	fmt.Println(string(beautifulJsonByte))
}
