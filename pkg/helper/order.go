package helper

import (
	"archive/zip"
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/google/uuid"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"gitlab.udevs.io/ekadastr/ek_integration_service/config"
	"gitlab.udevs.io/ekadastr/ek_integration_service/genproto/entity_service"
	"gitlab.udevs.io/ekadastr/ek_integration_service/genproto/integration_service"
	"gitlab.udevs.io/ekadastr/ek_integration_service/modules/ek_variables/ek_integration_service"
	"gitlab.udevs.io/ekadastr/ek_integration_service/pkg/util"
)

func RequestToChange(request map[string]interface{}, fileName, folderName, auctionUrl string) (response ek_integration_service.BaseOrderResponse, err error) {
	localFileUrl := folderName + fileName
	file, err := os.Open(localFileUrl)
	if err != nil {
		_ = os.Remove(localFileUrl)
		fmt.Println("Error upload open file: --> ", err, localFileUrl)
		return response, err
	}
	defer file.Close()

	contentType, requestBodyBuffer, err := CreateReqBody(request, localFileUrl, fileName)

	r, _ := http.NewRequest("POST", auctionUrl+"/ws/api/services/common/request-to-change", requestBodyBuffer)
	r.Header.Add("Content-Type", contentType)
	client := &http.Client{}
	resp, err := client.Do(r)
	if err != nil {
		_ = os.Remove(localFileUrl)
		fmt.Println("Error upload request: --> ", err)
		return response, err
	}

	err = json.NewDecoder(resp.Body).Decode(&response)
	if err != nil {
		_ = os.Remove(localFileUrl)
		fmt.Println("Error upload decode: --> ", err)
		return response, err
	}

	_ = os.Remove(localFileUrl)

	fmt.Println(response)

	return response, err
}

func CancelOrderFromAuction(document map[string]interface{}, fileName, auctionUrl string) (response ek_integration_service.CreateDocumentResponse, err error) {
	file, err := os.Open(fileName)
	if err != nil {
		_ = os.Remove(fileName)
		fmt.Println("Error upload open file: --> ", err, fileName)
	}
	defer file.Close()

	contentType, requestBodyBuffer, err := CreateReqBody(document, fileName, "")

	fmt.Println(document)
	fmt.Println(fileName)

	r, _ := http.NewRequest("POST", auctionUrl+"/ws/api/services/common/cancel-order", requestBodyBuffer)
	r.Header.Add("Content-Type", contentType)
	client := &http.Client{}
	resp, err := client.Do(r)
	if err != nil {
		_ = os.Remove(fileName)
		fmt.Println("Error upload request: --> ", err)
	}
	err = json.NewDecoder(resp.Body).Decode(&response)
	if err != nil {
		_ = os.Remove(fileName)
		fmt.Println("Error upload decode: --> ", err)
	}

	_ = os.Remove(fileName)

	fmt.Println(response)

	return response, err
}

func UploadDocumentToAuction(cfg config.Config, document map[string]interface{}, fileName, folderName string, nameInAuction string, auctionUrl string) (response ek_integration_service.CreateDocumentResponse, err error) {
	localFileUrl := folderName + fileName
	auctionFileName := GetDocumentFileExtension(fileName, nameInAuction)

	file, err := os.Open(localFileUrl)
	if err != nil {
		_ = os.Remove(localFileUrl)
		return response, err
	}
	defer file.Close()

	contentType, requestBodyBuffer, err := CreateReqBody(document, localFileUrl, auctionFileName)

	// fmt.Println(document)
	fmt.Println(localFileUrl)
	var r *http.Request
	if OnDevelopment(cfg) {
		// LOCAL
		viaRequestBody := ek_integration_service.ViaServerRequest{
			AuctionUrl:  auctionUrl + "/ws/api/services/documents",
			RequestType: http.MethodPost,
			Body: map[string]interface{}{
				"file_upload": requestBodyBuffer,
			},
			ContentType: contentType,
		}
		viaRequestBodyByte, _ := json.Marshal(viaRequestBody)
		viaRequestBodyBuffer := bytes.NewBuffer(viaRequestBodyByte)
		r, _ = http.NewRequest(http.MethodPost, "https://api-admin.yerelektron.uz/v1/query-auction-via-server", viaRequestBodyBuffer)
	} else {
		// PROD
		r, _ = http.NewRequest("POST", auctionUrl+"/ws/api/services/documents", requestBodyBuffer)
		r.Header.Add("Content-Type", contentType)
	}

	client := &http.Client{}
	resp, err := client.Do(r)
	if err != nil {
		_ = os.Remove(localFileUrl)
		return response, err
	}
	err = json.NewDecoder(resp.Body).Decode(&response)

	_ = os.Remove(localFileUrl)

	return response, err
}

func UploadImageToAuction(cfg config.Config, image map[string]interface{}, fileName, folderName, nameInAuction, auctionUrl string) (response ek_integration_service.CreateImageResponse, err error) {
	localFileUrl := folderName + fileName
	auctionFileName := GetDocumentFileExtension(fileName, nameInAuction)

	file, err := os.Open(localFileUrl)
	if err != nil {
		fmt.Println("Error: --> ", err)
		fmt.Println(err)
	}
	defer file.Close()

	contentType, requestBodyBuffer, err := CreateReqBody(image, localFileUrl, auctionFileName)

	fmt.Println(image)
	fmt.Println(localFileUrl)
	var r *http.Request
	if OnDevelopment(cfg) {
		// LOCAL
		viaRequestBody := ek_integration_service.ViaServerRequest{
			AuctionUrl:  auctionUrl + "/ws/api/services/images",
			RequestType: http.MethodPost,
			Body: map[string]interface{}{
				"file_upload": requestBodyBuffer,
			},
			ContentType: contentType,
		}
		viaRequestBodyByte, _ := json.Marshal(viaRequestBody)
		viaRequestBodyBuffer := bytes.NewBuffer(viaRequestBodyByte)
		r, _ = http.NewRequest(http.MethodPost, "https://api-admin.yerelektron.uz/v1/query-auction-via-server", viaRequestBodyBuffer)
	} else {
		// PROD
		r, _ = http.NewRequest("POST", auctionUrl+"/ws/api/services/images", requestBodyBuffer)
		r.Header.Add("Content-Type", contentType)
	}

	client := &http.Client{}
	resp, err := client.Do(r)
	if err != nil {
		fmt.Println("Error: --> ", err)
		fmt.Println(err)
		_ = os.Remove(localFileUrl)
	}
	err = json.NewDecoder(resp.Body).Decode(&response)
	if err != nil {
		fmt.Println("Error: --> ", err)
		fmt.Println(err)
	}

	_ = os.Remove(localFileUrl)

	fmt.Println(response)

	return response, err
}

func GetFileFromMinio(cfg config.Config, fileName string, folderName string) error {
	minioClient, err := minio.New(cfg.MinioDomain, &minio.Options{
		Creds:  credentials.NewStaticV4(cfg.MinioAccessKeyID, cfg.MinioSecretAccesKey, ""),
		Secure: true,
	})
	if err != nil {
		fmt.Println("Error minio client create: --> ", err)
		return err
	}

	exists, _ := minioClient.BucketExists(context.Background(), cfg.BucketName)
	if !exists {
		err = minioClient.MakeBucket(context.Background(), cfg.BucketName, minio.MakeBucketOptions{Region: cfg.MinioDomain})
		fmt.Println(err)
		if err != nil {
			return err
		}
	}

	err = minioClient.FGetObject(context.Background(), cfg.BucketName, fileName, folderName+fileName, minio.GetObjectOptions{})
	if err != nil {
		fmt.Println("Error: --> ", err, "bucketname: ", cfg.BucketName, "filename: ", fileName)
		return err
	}
	return nil
}

func UploadToMinio(cfg config.Config, filepath string, nameAddon string) (string, error) {
	minioClient, err := minio.New(cfg.MinioDomain, &minio.Options{
		Creds:  credentials.NewStaticV4(cfg.MinioAccessKeyID, cfg.MinioSecretAccesKey, ""),
		Secure: true,
	})
	if err != nil {
		fmt.Println("1", err)
		return "", err
	}

	exists, err := minioClient.BucketExists(context.Background(), cfg.BucketName)
	if err != nil {
		fmt.Println("3", err)
		return "", err
	}

	if !exists {
		err = minioClient.MakeBucket(context.Background(), cfg.BucketName, minio.MakeBucketOptions{Region: ""})
		if err != nil {
			fmt.Println("4", err)
			return "", err
		}
	}

	file, err := os.Open(filepath)
	if err != nil {
		fmt.Println(err)
		return "", err
	}
	defer file.Close()

	fileStat, err := file.Stat()
	if err != nil {
		fmt.Println(err)
		return "", err
	}

	newName := nameAddon + "_" + uuid.New().String() + "." + strings.Split(filepath, ".")[len(strings.Split(filepath, "."))-1]
	contentType := "application/octet-stream"
	extention, bol := IsFile(newName)
	
	if bol && extention == ".zip" {
		contentType = "application/zip"
	}
	
	uploadInfo, err := minioClient.PutObject(context.Background(), cfg.BucketName, newName, file, fileStat.Size(), minio.PutObjectOptions{ContentType: contentType})
	if err != nil {
		return "", err
	} else {
		fmt.Println("Succesfully uploaded:", uploadInfo)
	}

	return newName, err
}

func MakeGetOrderRequest(id, typeCode int64, cfg config.Config) (resp ek_integration_service.GetAuctionOrderResponse, req_body string, res_body string, err error) {
	client := http.Client{}
	var req *http.Request
	reqBody := map[string]interface{}{
		"order":    id,
		"language": "uz",
	}
	if typeCode == 1 {
		reqBody["username"] = cfg.AuctionUsername
		reqBody["password"] = cfg.AuctionPassword
	} else {
		reqBody["username"] = cfg.AuctionUsernameTypeSix
		reqBody["password"] = cfg.AuctionPasswordTypeSix
	}
	// add request body to get order details from auction
	reqBodyByte, err := json.Marshal(reqBody)
	if err != nil {
		return resp, "", "", err
	}

	if OnDevelopment(cfg) {
		/// IN LOCAL
		URL := fmt.Sprintf("https://api-admin.yerelektron.uz/v1/get-order-in-auction?order_id=%d&category_id=33", id)
		req, err = http.NewRequest("GET", URL, nil)
		if err != nil {
			return resp, "", "", err
		}
	} else {
		// /// IN PROD
		// create request credentials for client
		reqBodyBuff := bytes.NewBuffer(reqBodyByte)
		req, err = http.NewRequest("GET", cfg.AuctionOrderGetURL, reqBodyBuff)
		if err != nil {
			return resp, "", "", err
		}
	}

	req.Header.Set("Content-Type", "application/json")

	// send prepared request to auction
	response, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return resp, "", "", err
	}

	data, err := ioutil.ReadAll(response.Body)
	if err != nil {
		fmt.Println(err)
		return resp, "", "", err
	}

	// if reponse status code is error code return new error
	if response.StatusCode > 300 {
		return resp, string(reqBodyByte), string(data), fmt.Errorf("error while sending order get request to auction. Status code: %d", response.StatusCode)
	}

	// parse response body from auction to function's response object

	err = json.Unmarshal(data, &resp)
	if err != nil {
		fmt.Println(err)
		return resp, "", "", err
	}

	return resp, string(reqBodyByte), string(data), nil
}

func MakeCreateOrderRequest(req *integration_service.CreateUpdateOrderRequest, orderId, typeCode int64, url string, cfg config.Config) (resp *ek_integration_service.CreateOrderResponse, orderBody map[string]interface{}, req_body string, res_body string, err error) {
	var response *http.Response
	var postBody []byte
	auctionOrder := map[string]interface{}{
		"language":        "uz",
		"soato":           req.GetSoato(),
		"name":            req.GetName() + " yer uchastkasi",
		"price":           req.GetPrice(),
		"address":         req.GetAddress(),
		"category":        req.Category,
		"closed":          req.Closed,
		"prop_set":        req.PropSet,
		"additional_info": "",
		"lat":             req.Latitude,
		"lng":             req.Longitude,
	}

	if orderId > 0 {
		auctionOrder["order"] = orderId
	}
	if typeCode == 1 {
		auctionOrder["username"] = cfg.AuctionUsername
		auctionOrder["password"] = cfg.AuctionPassword
	} else {
		auctionOrder["username"] = cfg.AuctionUsernameTypeSix
		auctionOrder["password"] = cfg.AuctionPasswordTypeSix
	}

	if OnDevelopment(cfg) {
		// LOCAL
		client := http.Client{}
		viaRequestBody := ek_integration_service.ViaServerRequest{
			AuctionUrl:  url,
			RequestType: http.MethodPost,
			Body:        auctionOrder,
		}
		viaRequestBodyByte, _ := json.Marshal(viaRequestBody)
		viaRequestBodyBuffer := bytes.NewBuffer(viaRequestBodyByte)
		viaReq, err := http.NewRequest(http.MethodPost, "https://api-admin.yerelektron.uz/v1/query-auction-via-server", viaRequestBodyBuffer)
		if err != nil {
			return resp, auctionOrder, "", "", err
		}
		viaReq.Header.Set("Content-Type", "application/json")
		response, err = client.Do(viaReq)
		if err != nil {
			return resp, auctionOrder, "", "", err
		}
	} else {
		// PROD
		postBody, err = json.Marshal(auctionOrder)
		if err != nil {
			fmt.Println("Error create 1: --> ", err)
			return resp, auctionOrder, "", "", err
		}

		responseBody := bytes.NewBuffer(postBody)
		response, err = http.Post(url, "application/json", responseBody)
		if err != nil {
			fmt.Println("Error create 2 : --> ", err)
			return resp, auctionOrder, "", "", err
		}
	}

	defer response.Body.Close()

	err = json.NewDecoder(response.Body).Decode(&resp)
	if err != nil {
		fmt.Println("Error create 3: --> ", err)
		return resp, auctionOrder, "", "", err
	}
	if resp.OrderId == 0 {
		resp.OrderId = orderId
	}
	auctionOrder["entity_id"] = req.EntityId
	auctionOrder["order_id"] = resp.OrderId
	resBytes, err := io.ReadAll(response.Body)
	if err != nil {
		return resp, auctionOrder, "", "", err
	}
	return resp, auctionOrder, string(postBody), string(resBytes), err
}

func SaveOrderDetails(cfg config.Config, request map[string]interface{}, url string) (response ek_integration_service.BaseOrderResponse, err error) {

	postBody, err := json.Marshal(request)
	if err != nil {
		fmt.Println("Error create 1: --> ", err)
		return response, err
	}

	requestBody := bytes.NewBuffer(postBody)
	var res *http.Response

	if OnDevelopment(cfg) {
		// LOCAL
		viaRequestBody := ek_integration_service.ViaServerRequest{
			AuctionUrl:  url,
			RequestType: http.MethodPost,
			Body:        request,
			ContentType: "application/json",
		}
		viaRequestBodyByte, _ := json.Marshal(viaRequestBody)
		viaRequestBodyBuffer := bytes.NewBuffer(viaRequestBodyByte)
		res, err = http.Post("https://api-admin.yerelektron.uz/v1/query-auction-via-server", "application/json", viaRequestBodyBuffer)
	} else {
		// PROD
		res, err = http.Post(url, "application/json", requestBody)
	}

	if err != nil {
		fmt.Println("Error create 2 : --> ", err)
		return response, err
	}
	defer res.Body.Close()

	err = json.NewDecoder(res.Body).Decode(&response)
	if err != nil {
		fmt.Println("Error create 3: --> ", err)
		return response, err
	}
	fmt.Println(response, err)

	return response, err
}

func AddUserToClosedAuction(request map[string]interface{}, url string) (response ek_integration_service.BaseOrderResponse, err error) {

	postBody, err := json.Marshal(request)
	if err != nil {
		fmt.Println("AddUserToClosedAuction.Json.Marshal: --> ", err)
		return response, err
	}

	responseBody := bytes.NewBuffer(postBody)
	// fmt.Println("POST body:", request)
	// fmt.Println("URL", url)

	res, err := http.Post(url, "application/json", responseBody)
	if err != nil {
		fmt.Println("AddUserToClosedAuction.Http.POST: --> ", err)
		return response, err
	}
	defer res.Body.Close()

	err = json.NewDecoder(res.Body).Decode(&response)
	if err != nil {
		fmt.Println("AddUserToClosedAuction.DecodeResponse: --> ", err)
		return response, err
	}
	fmt.Println(response, err)

	return response, err
}

func MakeDeleteImageDocumentRequest(request map[string]interface{}, url string, cfg config.Config) (resp ek_integration_service.BaseOrderResponse, err error) {
	client := http.Client{}
	requestBodyByte, _ := json.Marshal(request)
	var req *http.Request

	if OnDevelopment(cfg) {
		// LOCAL
		viaRequestBody := ek_integration_service.ViaServerRequest{
			AuctionUrl:  url,
			RequestType: http.MethodDelete,
			Body:        request,
		}
		viaRequestBodyByte, _ := json.Marshal(viaRequestBody)
		viaRequestBodyBuffer := bytes.NewBuffer(viaRequestBodyByte)

		req, err = http.NewRequest(http.MethodPost, "https://api-admin.yerelektron.uz/v1/query-auction-via-server", viaRequestBodyBuffer)
	} else {
		// PROD
		requestBody := bytes.NewBuffer(requestBodyByte)
		strings.NewReader(string(requestBodyByte))

		req, err = http.NewRequest(http.MethodDelete, url, requestBody)
	}

	if err != nil {
		return resp, errors.New("1 ==> " + err.Error())
	}
	req.Header.Set("Content-Type", "application/json")
	fmt.Println("URL: ", url)
	response, err := client.Do(req)
	if err != nil {
		return resp, errors.New("2 ==> " + err.Error())
	}

	if response.StatusCode == 404 {
		resp.ResultCode = 0
		return resp, nil
	}

	data, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return resp, errors.New("3 ==> " + err.Error())
	}

	err = json.Unmarshal(data, &resp)
	if err != nil {
		log.Printf("error decoding response: %v", err)
		return resp, errors.New("4 ==> req :" + util.JSONStringify(request) + err.Error() + string(data))
	}

	// if no files found, do not stop process
	if resp.ResultCode == 137 {
		resp.ResultCode = 0
		return resp, nil
	}

	resp.ReqBody = string(requestBodyByte)
	resp.ResBody = string(data)

	return resp, err
}

func MakeSendOrderRequest(cfg config.Config, body map[string]interface{}, url string) error {
	var (
		resp ek_integration_service.BaseOrderResponse
	)

	postBody, err := json.Marshal(body)
	if err != nil {
		fmt.Println("Error send 1: --> ", err)
		return err
	}

	responseBody := bytes.NewBuffer(postBody)
	fmt.Println(body)
	fmt.Println(url)

	var response *http.Response

	if OnDevelopment(cfg) {
		// LOCAL
		viaRequestBody := ek_integration_service.ViaServerRequest{
			AuctionUrl:  url,
			RequestType: http.MethodPost,
			Body:        body,
			ContentType: "application/json",
		}
		viaRequestBodyByte, _ := json.Marshal(viaRequestBody)
		viaRequestBodyBuffer := bytes.NewBuffer(viaRequestBodyByte)
		response, err = http.Post("https://api-admin.yerelektron.uz/v1/query-auction-via-server", "application/json", viaRequestBodyBuffer)
	} else {
		// PROD
		response, err = http.Post(url, "application/json", responseBody)
	}

	if err != nil {
		fmt.Println("Error send 2 : --> ", err)
		return err
	}
	defer response.Body.Close()

	err = json.NewDecoder(response.Body).Decode(&resp)
	if err != nil {
		fmt.Println("Error send 3: --> ", err)
		return err
	}

	if resp.ResultCode != 0 {
		return errors.New("error from auction: " + resp.ResultMsg)
	}

	fmt.Println(response)

	return nil
}

func DownloadFile(filepath string, url string, isEqaror bool) (err error) {
	// Create the file
	out, err := util.CreateDolders(filepath)
	if err != nil {
		return errors.New("ochishda xatolik 1 ===> " + err.Error())
	}
	defer out.Close()
	if isEqaror {
		arr := strings.Split(url, "")
		bo := true
		id := ""
		for _, v := range arr {
			if v != "=" && bo {
				continue
			}
			bo = false
			if v == "&" {
				break
			}
			id = id + v
		}
		id =id[1:]
		url = "https://api-eqaror.gov.uz/v1/doc/show?id="+id
	}
	// Get the data
	resp, err := http.Get(url)
	fmt.Println("  --  Downloading as ", filepath, "...")
	if err != nil {
		return errors.New("qarorni olish xatolik ===> " + err.Error())
	}
	defer resp.Body.Close()

	// Check server response
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("qaror faylini olishda xatolik yuz berdi (kod 404)")
	}

	// Writer the body to file
	_, err = io.Copy(out, resp.Body)
	if err != nil {
		return errors.New("faylni yozishda xatolik ===> " + err.Error())
	}

	return nil
}

func MakeGetProtocolRequest(id, typeCode int64, cfg config.Config) (resp ek_integration_service.GetAuctionProtoolResponse, req_body string, res_body string, err error) {
	client := http.Client{}
	var req *http.Request
	reqBody := map[string]interface{}{
		"order":    id,
		"language": "uz",
	}
	if typeCode == 1 {
		reqBody["username"] = cfg.AuctionUsername
		reqBody["password"] = cfg.AuctionPassword
	} else {
		reqBody["username"] = cfg.AuctionUsernameTypeSix
		reqBody["password"] = cfg.AuctionPasswordTypeSix
	}
	// add request body to get order details from auction
	reqBodyByte, err := json.Marshal(reqBody)
	if err != nil {
		return resp, "", "", err
	}

	if OnDevelopment(cfg) {
		// /// IN LOCAL
		URL := fmt.Sprintf("https://api-admin.yerelektron.uz/v1/get-order-in-auction?order_id=%d&category_id=111", id)
		req, err = http.NewRequest("GET", URL, nil)
		if err != nil {
			return resp, "", "", err
		}
	} else {
		// /// IN PROD
		// create request credentials for client
		reqBodyBuff := bytes.NewBuffer(reqBodyByte)
		req, err = http.NewRequest("GET", cfg.AuctionGetProtocolURL, reqBodyBuff)
		if err != nil {
			return resp, "", "", err
		}
	}

	req.Header.Set("Content-Type", "application/json")

	// send prepared request to auction
	response, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return resp, "", "", err
	}

	data, err := ioutil.ReadAll(response.Body)
	if err != nil {
		fmt.Println(err)
		return resp, "", "", err
	}

	// if reponse status code is error code return new error
	if response.StatusCode > 300 {
		return resp, string(reqBodyByte), string(data), fmt.Errorf("error while sending order get request to auction. Status code: %d", response.StatusCode)
	}

	// parse response body from auction to function's response object
	err = json.Unmarshal(data, &resp)
	if err != nil {
		fmt.Println(err)
		return resp, "", "", err
	}

	return resp, string(reqBodyByte), string(data), nil
}

func Zipify(baseFolder string, files []string) (string, error) {
	outFile, err := os.Create(baseFolder + "docs.zip")
	if err != nil {
		fmt.Println(err)
		return "", err
	}
	defer outFile.Close()

	w := zip.NewWriter(outFile)

	// Add some files to the archive.
	for _, file := range files {
		dat, err := ioutil.ReadFile(file)
		if err != nil {
			fmt.Println(err)
		}

		filePathInZip := file
		fileNameParts := strings.Split(file, "/")
		if len(fileNameParts) > 1 {
			filePathInZip = fileNameParts[len(fileNameParts)-1]
		}

		f, err := w.Create(filePathInZip)
		if err != nil {
			fmt.Println(err)
		}
		_, err = f.Write(dat)
		if err != nil {
			fmt.Println(err)
		}
	}

	err = w.Close()
	if err != nil {
		fmt.Println(err)
		return "", err
	}

	return baseFolder + "docs.zip", nil
}

func FetchFilesAndZip(fileUrls []map[string]string, folderToSave string) (string, error) {
	files := []string{}

	for i, url_key := range fileUrls {

		fullPath := fmt.Sprintf("%s%s%d.pdf", folderToSave, url_key["as"], i+1)
		fmt.Println(" - ", fullPath)
		err := DownloadFile(fullPath, url_key["url"], false)
		if err != nil {
			return "", err
		}
		files = append(files, fullPath)
	}

	return Zipify(folderToSave, files)
}

func GotDifferentAuctionPushDetails(newDetails map[string]string, props []*entity_service.GetEntityProperty) bool {
	oldDetails := map[string]string{}
	for _, prop := range props {
		oldDetails[prop.Property.Id] = prop.Value
	}

	for key, val := range newDetails {
		if val != "" {
			if oldDetails[key] != val {
				return true
			}
		}
	}

	return false
}

func OnDevelopment(cfg config.Config) bool {
	return cfg.Environment == "DEV"
}
