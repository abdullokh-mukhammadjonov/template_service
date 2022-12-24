package helper

import (
	"archive/zip"
	"context"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"strings"

	"github.com/google/uuid"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"gitlab.udevs.io/ekadastr/ek_integration_service/config"
	"gitlab.udevs.io/ekadastr/ek_integration_service/pkg/util"
)

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
		id = id[1:]
		url = "https://api-eqaror.gov.uz/v1/doc/show?id=" + id
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

func OnDevelopment(cfg config.Config) bool {
	return cfg.Environment == "DEV"
}
