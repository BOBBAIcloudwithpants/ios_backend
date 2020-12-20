package service

import (
	"context"
	"fmt"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"github.com/spf13/viper"
	"crypto/md5"
	"log"
	"mime/multipart"
	"time"
	b64 "encoding/base64"

)

var MinioClient *minio.Client

type File struct{
	F multipart.File
	H *multipart.FileHeader
}

func init() {
	viper.SetConfigName("configure")
	viper.SetConfigType("json")
	viper.AddConfigPath("$GOPATH/src/github.com/service-computing-2020/bbs_backend/config/")
	viper.AddConfigPath("config/")
	if err := viper.ReadInConfig(); err != nil {
		panic(fmt.Errorf("Fatal error config file: %s \n", err))
	}

	var err error
	endPoint := viper.GetString(`minio.endPoint`)
	accessKeyID := viper.GetString(`minio.accessKeyID`)
	secretAccessKey := viper.GetString(`minio.secretAccessKey`)
	secure := viper.GetBool(`minio.secure`)
	MinioClient, err = minio.New(endPoint, &minio.Options{
		Creds: credentials.NewStaticV4(accessKeyID, secretAccessKey, ""),
		Secure: secure,
	})
	if err != nil {
		log.Fatalln(err)
	} else {
		log.Println("%#v\n", MinioClient)
	}
}

// 传入api路径和文件扩展名，如 getUploadName('/api/users/1/avatar', '.png')
func GetUploadName(ext string) string {
	now := time.Now()
	sec := now.String()
	prefix := md5.Sum([]byte(sec))
	input := prefix[:]
	sEnc := b64.StdEncoding.EncodeToString(input)
	name := fmt.Sprintf("%s%s",sEnc, ext)
	fmt.Println(name)
	return name
}

// 同上


// 上传多个文件，如果有文件上传出错则回滚之前的文件, 返回成功上传的文件名
func MultipleFilesUpload(files []File, bucketName string, path string ,ext string)([]string, error) {

	var names []string
	for idx, f := range files {
		new_path := fmt.Sprintf("%s%d", path, idx)

		filename, err := FileUpload(f.F, f.H, bucketName, new_path, ext)
		if err != nil {
			for _, del_f := range names {
				err := FileDelete(del_f, bucketName)
				if err != nil {
					panic(err.Error())
				}
			}
			return nil, err
		}
		names = append(names, filename)
	}
	return names, nil
}


func FileDelete(filename string, bucketName string) error{
	ctx := context.Background()
	return MinioClient.RemoveObject(ctx, bucketName, filename, minio.RemoveObjectOptions{})
}

func FileUpload(file multipart.File,header *multipart.FileHeader, bucketName string, path string, ext string)(filename string, err error) {
	ctx := context.Background()
	filename = GetUploadName(ext)
	_, err = MinioClient.PutObject(ctx, bucketName, filename, file, header.Size, minio.PutObjectOptions{ContentType: "application/octet-stream"})
	if err != nil {
		return "", err
	} else {
		return filename, err
	}
}

//func FileDownload(filename string, bucketName string, ext string) (*minio.Object, error) {
//	ctx := context.Background()
//	return MinioClient.GetObject(ctx, bucketName, GetDownloadName(filename, ext), minio.GetObjectOptions{})
//}

func FileDownloadByName(filename string, bucketName string) (*minio.Object, error) {
	ctx := context.Background()
	return MinioClient.GetObject(ctx, bucketName, filename, minio.GetObjectOptions{})
}


