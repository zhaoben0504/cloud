package server

import (
	"bytes"
	"cloud/server/model"
	"cloud/tool"
	"context"
	"encoding/json"
	uuid2 "github.com/hashicorp/go-uuid"
	"github.com/tencentyun/cos-go-sdk-v5"
	"io"
	"math/rand"
	"mime/multipart"
	"net/http"
	"net/url"
	"os"
	"path"
	"strconv"
	"strings"
	"time"
)

func GenerateToken(user *model.UserBasic) (string, error) {
	tokenInfo, err := json.Marshal(user)
	if err != nil {
		return "", err
	}

	token := tool.MD5(string(tokenInfo))
	if err != nil {
		return "", err
	}

	// 塞入Redis
	err = SetInfoInRedis(token, string(tokenInfo), TokenExpire)
	if err != nil {
		return "", err
	}

	return token, nil
}

func GetUserInfoFromToken(token string) (userInfo *model.UserBasic, err error) {
	err = json.Unmarshal([]byte(token), &userInfo)
	if err != nil {
		tool.Logger.Error(err)
		return nil, err
	}

	b, err := GetEngine().Where("id=? AND deleted_at=0", userInfo.Id).Get(&userInfo)
	if err != nil {
		tool.Logger.Error(err)
		return nil, err
	}
	if !b {
		tool.Logger.Error("not found")
		return nil, nil
	}
	return userInfo, nil
}

func SendEmailCode() string {
	return "123456"
}

func GenerateEmailCode() string {
	str := "1234567890"
	code := ""
	rand.Seed(time.Now().UnixNano())
	for i := 0; i < EmailCodeLen; i++ {
		code += string(str[rand.Intn(len(str))])
	}
	return code
}

func GenerateUUID() string {
	str, err := uuid2.GenerateUUID()
	if err != nil {
		return ""
	}
	return str[0:15]
}

// UploadCos upload file to COS
func UploadCos(fileHeader multipart.FileHeader) (string, error) {
	file, err := fileHeader.Open()
	if err != nil {
		tool.Logger.Error(err.Error())
		return "", err
	}
	defer func(file multipart.File) {
		err = file.Close()
		if err != nil {
			tool.Logger.Error(err.Error())
			return
		}
	}(file)

	u, _ := url.Parse(COSADDR)
	b := &cos.BaseURL{BucketURL: u}
	client := cos.NewClient(b, &http.Client{
		Transport: &cos.AuthorizationTransport{
			SecretID:  os.Getenv(CloudId),
			SecretKey: os.Getenv(CloudKey),
		},
	})

	key := "mystorage/" + GenerateUUID() + path.Ext(fileHeader.Filename)

	_, err = client.Object.Put(
		context.Background(), key, file, nil,
	)
	if err != nil {
		tool.Logger.Error(err.Error())
		return "", err
	}
	return COSADDR + "/" + key, nil
}

func CosInitPart(ext string) (string, string, error) {
	u, _ := url.Parse(COSADDR)
	b := &cos.BaseURL{BucketURL: u}
	client := cos.NewClient(b, &http.Client{
		Transport: &cos.AuthorizationTransport{
			SecretID:  os.Getenv(CloudId),
			SecretKey: os.Getenv(CloudKey),
		},
	})
	key := "mystorage/" + GenerateUUID() + "." + ext
	v, _, err := client.Object.InitiateMultipartUpload(context.Background(), key, nil)
	if err != nil {
		return "", "", err
	}
	return key, v.UploadID, nil
}

func CosPartUpload(r *http.Request) (string, error) {
	u, _ := url.Parse(COSADDR)
	b := &cos.BaseURL{BucketURL: u}
	client := cos.NewClient(b, &http.Client{
		Transport: &cos.AuthorizationTransport{
			SecretID:  os.Getenv(CloudId),
			SecretKey: os.Getenv(CloudKey),
		},
	})
	key := r.PostForm.Get("key")
	UploadID := r.PostForm.Get("uploadId")
	partNumber, err := strconv.Atoi(r.PostForm.Get("partNumber"))
	f, _, err := r.FormFile("file")
	if err != nil {
		return "", err
	}
	buf := bytes.NewBuffer(nil)
	io.Copy(buf, f)
	resp, err := client.Object.UploadPart(
		context.Background(), key, UploadID, partNumber, bytes.NewReader(buf.Bytes()), nil,
	)
	if err != nil {
		return "", nil
	}
	return strings.Trim(resp.Header.Get("ETag"), "\""), nil
}

// 分片上传的结束
func CosChunkUploadFinish(key, uploadId string, co []cos.Object) error {
	u, _ := url.Parse(COSADDR)
	b := &cos.BaseURL{BucketURL: u}
	client := cos.NewClient(b, &http.Client{
		Transport: &cos.AuthorizationTransport{
			SecretID:  os.Getenv(CloudId),
			SecretKey: os.Getenv(CloudKey),
		},
	})

	opt := &cos.CompleteMultipartUploadOptions{}
	opt.Parts = append(opt.Parts, co...)
	_, _, err := client.Object.CompleteMultipartUpload(
		context.Background(), key, uploadId, opt,
	)
	return err
}

func SetInfoInRedis(key, value string, duration int) error {
	_, err := GetRedisClient().Set(context.Background(), key, value, time.Second*time.Duration(duration)).Result()
	if err != nil {
		tool.Logger.Error(err.Error())
		return err
	}
	return nil
}
