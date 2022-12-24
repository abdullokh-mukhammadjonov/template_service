package unit_test

import (
	"encoding/json"
	"fmt"
	"testing"
)

func TestUnmarshalAuctionResponse(t *testing.T) {
	var resp interface{}
	responseBody := `{
		"result_code": 137,
		"result_msg": "Ushbu buyurtma bo'yicha so'ralgan fayl topilmadi"
	}`

	bytes, err := json.Marshal(responseBody)

	if err != nil {
		t.Errorf("%s", err.Error())
	}

	fmt.Println(string(bytes))

	err = json.Unmarshal(bytes, &resp)

	if err != nil {
		t.Errorf("%s", err.Error())
	}

	fmt.Println(resp)
}
