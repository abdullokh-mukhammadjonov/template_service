package unit_test

import (
	"fmt"
	"os"
	"testing"

	"gitlab.udevs.io/ekadastr/ek_integration_service/config"
	"gitlab.udevs.io/ekadastr/ek_integration_service/pkg/helper"
	// "github.com/bxcodec/faker/v3"
)

type orderFiles struct {
	Urls    []map[string]string
	OrderID int64
	Want    string
}

func TestZipify(t *testing.T) {
	cfg := config.Load()
	cases := []orderFiles{
		// {
		// 	Urls: []map[string]string{
		// 		{
		// 			"url": "https://files.e-auksion.uz/order-contract?hash=b599c1cf8cad25c5cc8b7a5535c0e2ea4b6b718c",
		// 			"as":  "shartnoma",
		// 		},
		// 		{
		// 			"url": "https://files.e-auksion.uz/protocol-file?hash=21844288b477a2d6e25b34ee1a27e9b94608a03b",
		// 			"as":  "bayonnoma",
		// 		},
		// 	},
		// 	OrderID: 571062,
		// 	Want:    "protocol_bayonnoma/5710621/",
		// },
		// {
		// 	Urls: []map[string]string{
		// 		{
		// 			"url": "https://files.e-auksion.uz/order-contract?hash=ea7431cc996101ebb6b9727f41eb1bcfcb285262",
		// 			"as":  "shartnoma",
		// 		},
		// 		{
		// 			"url": "https://files.e-auksion.uz/protocol-file?hash=f4516fa2e9eeaa22b5cb2f0dd8a6bf74cba3e49e",
		// 			"as":  "bayonnoma",
		// 		},
		// 	},
		// 	OrderID: 572410,
		// 	Want:    "protocol_bayonnoma/572410/",
		// },
		{
			Urls: []map[string]string{
				{
					"url": "https://files.e-auksion.uz/order-contract?hash=ce6df0e5a5359e9db24874b457ff67d49a2bc228",
					"as":  "shartnoma",
				},
				{
					"url": "https://files.e-auksion.uz/protocol-file?hash=86b4f0dd829e51f2576ddc92d2904141e6364694",
					"as":  "bayonnoma",
				},
			},
			OrderID: 583317,
			Want:    "protocol_bayonnoma/583317/",
		},
	}

	for _, cs := range cases {
		got, err := helper.FetchFilesAndZip(cs.Urls, fmt.Sprintf("protocol_bayonnoma/%d/", cs.OrderID))

		if err != nil {
			t.Errorf("Zipify(%d) == %q, want %q", cs.OrderID, got, cs.Want)
			fmt.Println("error :", err)
		}

		loc, err := helper.UploadToMinio(cfg, got, fmt.Sprintf("winner_docs_%d", cs.OrderID))
		if err != nil {
			t.Errorf("TestZipify => UploadToMinio \n%s", err.Error())
		}

		fmt.Print("  -- -- saved_in :", cfg.MinioDomain+"/"+cfg.BucketName+"/"+loc, "\n")

		err = os.RemoveAll(fmt.Sprintf("protocol_bayonnoma/%d", cs.OrderID))
		if err != nil {
			t.Errorf("ERROR :" + err.Error())
		}
		fmt.Print("  -- -- -- SUCCESS \n\n")
	}
	// RUN: go test -v ./... | sed ''/PASS/s//$(printf "\033[32mPASS\033[0m")/'' | sed ''/FAIL/s//$(printf "\033[31mFAIL\033[0m")/'' | sed ''/Downloading/s//$(printf "\033[30mDownloading\033[0m")/'' | sed ''/saved_in/s//$(printf "\033[36msaved_in\033[0m")/'' | sed ''/SUCCESS/s//$(printf "\033[35mSUCCESS\033[0m")/''
}
