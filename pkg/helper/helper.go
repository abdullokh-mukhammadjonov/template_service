package helper

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"io/ioutil"
	"mime/multipart"
	"os"
	"path"
	"strings"

	"github.com/abdullokh-mukhammadjonov/template_service/pkg/logger"
	"go.mongodb.org/mongo-driver/mongo"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var (
	MaxLengthPassword                       = 30
	MinLengthPassword                       = 8
	MaxLengthLogin                          = 15
	MinLengthLogin                          = 5
	HasNumber                               = "[0-9]"
	HasAlphabeticCharacter                  = "[A-Za-z]"
	HasNumberOrAlphabeticOrSpecialCharacter = "^[A-Za-z0-9$_@.#]+$"
)

func GetDocumentFileExtension(org_file_name string, new_name string) string {
	slices := strings.Split(org_file_name, ".")
	if len(slices) < 2 {
		return org_file_name
	}
	return new_name + "." + slices[len(slices)-1]
}

func IsFile(fileName string) (string, bool) {
	fileTypes := []string{".xlsx", ".xls", ".doc", ".docx", ".jpeg", ".jpg", ".svg", ".png", ".pdf", ".zip"}
	for _, t := range fileTypes {
		if strings.HasSuffix(strings.ToLower(fileName), t) {
			return t, true
		}
	}
	return "", false
}

func HandleError(log logger.Logger, err error, message string, req interface{}) error {

	if err == mongo.ErrNoDocuments {
		log.Error(message+", Not Found", logger.Error(err), logger.Any("req", req))
		return status.Error(codes.NotFound, "Not Found")
	} else if err != nil {
		log.Error(message, logger.Error(err), logger.Any("req", req))
		return status.Error(codes.Internal, message+err.Error())
	}
	return nil
}
func MarshalUnmarshal(reqStructure, resStructure interface{}) error {
	byte, err := json.Marshal(reqStructure)
	if err != nil {
		return err
	}
	err = json.Unmarshal(byte, resStructure)
	return err
}

func CreateReqBody(bodyMap map[string]interface{}, filePath string, nameInAuction string) (string, io.Reader, error) {
	var err error

	buf := new(bytes.Buffer)
	bw := multipart.NewWriter(buf) // body writer

	f, err := os.Open(filePath)
	if err != nil {
		return "", nil, err
	}
	defer f.Close()

	for fieldName, value := range bodyMap {
		requestField, _ := bw.CreateFormField(fieldName)
		requestField.Write([]byte(value.(string)))
	}

	if nameInAuction == "" {
		nameInAuction = filePath
	}
	fw1, _ := bw.CreateFormFile("file", nameInAuction)
	io.Copy(fw1, f)

	bw.Close() //write the tail boundry
	return bw.FormDataContentType(), buf, nil
}

func CyrillToLatin(cyrill string) (latin string) {

	replacer := strings.NewReplacer(
		"????", "Be",
		"????", "be",
		"????", "Ve",
		"????", "ve",
		"????", "Ge",
		"????", "ge",
		"????", "Ge'",
		"????", "ge'",
		"????", "He",
		"????", "he",
		"????", "De",
		"????", "de",
		"????", "Je",
		"????", "je",
		"????", "Ze",
		"????", "ze",
		"????", "Ye",
		"????", "ye",
		"????", "Ke",
		"????", "ke",
		"????", "Qe",
		"????", "qe",
		"????", "Le",
		"????", "le",
		"????", "Me",
		"????", "me",
		"????", "Ne",
		"????", "ne",
		"????", "Pe",
		"????", "pe",
		"????", "Re",
		"????", "re",
		"????", "Se",
		"????", "se",
		"????", "Te",
		"????", "te",
		"????", "Fe",
		"????", "fe",
		"????", "Xe",
		"????", "xe",
		"????", "TSe",
		"????", "tse",
		"????", "Che",
		"????", "che",
		"????", "She",
		"????", "she",
		"????", "She",
		"????", "she",

		"??", "A",
		"??", "a",
		"??", "B",
		"??", "b",
		"??", "V",
		"??", "v",
		"??", "G",
		"??", "g",
		"??", "G'",
		"??", "g'",
		"??", "H",
		"??", "h",
		"??", "D",
		"??", "d",
		"??", "Ye",
		"??", "ye",
		"??", "J",
		"??", "j",
		"??", "Z",
		"??", "z",
		"??", "I",
		"??", "i",
		"??", "Y",
		"??", "y",
		"??", "K",
		"??", "k",
		"??", "Q",
		"??", "q",
		"??", "L",
		"??", "l",
		"??", "M",
		"??", "m",
		"??", "N",
		"??", "n",
		"??", "O",
		"??", "o",
		"??", "P",
		"??", "p",
		"??", "R",
		"??", "r",
		"??", "S",
		"??", "s",
		"??", "T",
		"??", "t",
		"??", "U",
		"??", "u",
		"??", "F",
		"??", "f",
		"??", "O'",
		"??", "o'",
		"??", "X",
		"??", "x",
		"??", "TS",
		"??", "ts",
		"??", "Ch",
		"??", "ch",
		"??", "Sh",
		"??", "sh",
		"??", "Sh",
		"??", "sh",
		"??", "'",
		"??", "'",
		"??", "I",
		"??", "i",
		"??", "E",
		"??", "e",
		"??", "Yu",
		"??", "yu",
		"??", "Ya",
		"??", "ya",
		"??", "Yo",
		"??", "yo",
		"????", "Ye",
		"????", "ye",
	)
	latin = replacer.Replace(cyrill)
	return latin
}

func ErrorTranslator(function_name string, point string, err error) error {
	switch function_name {
	case "CAO":
		switch point {
		case "9":
			return errors.New("auksionga avval yuklangan hujjatlarni o'chirishda xatolik (kod 9)")
		case "11":
			return errors.New("auksionga avval yuklangan rasmlarni o'chirishda xatolik (kod 11)")
		case "12":
			return errors.New("fayl yuklashda xatolik (kod 12)")
		case "19":
			return errors.New("topografik xarita topilmadi (kod 19)")
		case "22":
			return errors.New("situatsion sxema fayli yuklanmagan (kod 22)")
		case "38":
			if err.Error() == "bad status: 404 Not Found" {
				return errors.New("e-qarorda qaror fayli toplimadi (kod 38)")
			} else {
				return errors.New("qaror faylini olishda xatolik (kod 38)")
			}
		case "44":
			return errors.New("xulosa fayli yuklanmagan (kod 44)")
		case "53":
			return errors.New("asosiy rasm tanlanmagan. Kerakli rasmni yuklang (kod 53)")
		case "78":
			return errors.New("qaror topilmadi (kod 78)")
		case "79", "80":
			return errors.New("auksionga yuborishda xatolik (kod 79)")
		case "81":
			return errors.New("huquq turini aniqlashda muammo (kod 81)")
		}
	}
	return err
}

func ClearFolder(folder string) error {
	dir, err := ioutil.ReadDir("./" + folder)
	if err != nil {
		return err
	}

	for _, d := range dir {
		err := os.RemoveAll(path.Join([]string{folder, d.Name()}...))
		if err != nil {
			return err
		}
	}
	return nil
}
