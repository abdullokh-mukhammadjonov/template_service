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

	"gitlab.udevs.io/ekadastr/ek_integration_service/pkg/logger"
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
		"Бе", "Be",
		"бе", "be",
		"Ве", "Ve",
		"ве", "ve",
		"Ге", "Ge",
		"ге", "ge",
		"Ғе", "Ge'",
		"ғе", "ge'",
		"Ҳе", "He",
		"ҳе", "he",
		"Де", "De",
		"де", "de",
		"Же", "Je",
		"же", "je",
		"Зе", "Ze",
		"зе", "ze",
		"Йе", "Ye",
		"йе", "ye",
		"Ке", "Ke",
		"ке", "ke",
		"Қе", "Qe",
		"қе", "qe",
		"Ле", "Le",
		"ле", "le",
		"Ме", "Me",
		"ме", "me",
		"Не", "Ne",
		"не", "ne",
		"Пе", "Pe",
		"пе", "pe",
		"Ре", "Re",
		"ре", "re",
		"Се", "Se",
		"се", "se",
		"Те", "Te",
		"те", "te",
		"Фе", "Fe",
		"фе", "fe",
		"Хе", "Xe",
		"хе", "xe",
		"Це", "TSe",
		"це", "tse",
		"Че", "Che",
		"че", "che",
		"Ше", "She",
		"ше", "she",
		"Ще", "She",
		"ще", "she",

		"А", "A",
		"а", "a",
		"Б", "B",
		"б", "b",
		"В", "V",
		"в", "v",
		"Г", "G",
		"г", "g",
		"Ғ", "G'",
		"ғ", "g'",
		"Ҳ", "H",
		"ҳ", "h",
		"Д", "D",
		"д", "d",
		"Е", "Ye",
		"е", "ye",
		"Ж", "J",
		"ж", "j",
		"З", "Z",
		"з", "z",
		"И", "I",
		"и", "i",
		"Й", "Y",
		"й", "y",
		"К", "K",
		"к", "k",
		"Қ", "Q",
		"қ", "q",
		"Л", "L",
		"л", "l",
		"М", "M",
		"м", "m",
		"Н", "N",
		"н", "n",
		"О", "O",
		"о", "o",
		"П", "P",
		"п", "p",
		"Р", "R",
		"р", "r",
		"С", "S",
		"с", "s",
		"Т", "T",
		"т", "t",
		"У", "U",
		"у", "u",
		"Ф", "F",
		"ф", "f",
		"Ў", "O'",
		"ў", "o'",
		"Х", "X",
		"х", "x",
		"Ц", "TS",
		"ц", "ts",
		"Ч", "Ch",
		"ч", "ch",
		"Ш", "Sh",
		"ш", "sh",
		"Щ", "Sh",
		"щ", "sh",
		"Ъ", "'",
		"ъ", "'",
		"Ы", "I",
		"ы", "i",
		"Э", "E",
		"э", "e",
		"Ю", "Yu",
		"ю", "yu",
		"Я", "Ya",
		"я", "ya",
		"Ё", "Yo",
		"ё", "yo",
		"Ъе", "Ye",
		"ъе", "ye",
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
