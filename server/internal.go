package server

import (
	"bytes"
	"cloud/server/model"
	"cloud/tool"
	"context"
	"encoding/json"
	"errors"
	uuid2 "github.com/hashicorp/go-uuid"
	"github.com/tencentyun/cos-go-sdk-v5"
	"io"
	"math/rand"
	"mime/multipart"
	"net/http"
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

// UploadFile upload file to COS
func UploadFile(fileHeader multipart.FileHeader) (string, error) {
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

	filePath := "mystorage/" + GenerateUUID() + path.Ext(fileHeader.Filename)
	err = uploadCos(filePath, file)
	if err != nil {
		tool.Logger.Error(err.Error())
		return "", err
	}

	return filePath, nil
}

// DownloadFile download file from COS
func DownloadFile(fileId string) ([]byte, error) {
	if fileId == "" {
		tool.Logger.Error(errors.New(GetMsgByCode("zh", ParamErrCode)))
		return nil, errors.New(GetMsgByCode("zh", ParamErrCode))
	}

	filePath := "mystorage/" + fileId
	file, err := downloadCos(filePath)
	if err != nil {
		tool.Logger.Error(err.Error())
		return nil, err
	}
	return file, nil
}

func CosInitPart(ext string) (string, string, error) {
	key := "mystorage/" + GenerateUUID() + "." + ext
	v, _, err := GetCosClient().Object.InitiateMultipartUpload(context.Background(), key, nil)
	if err != nil {
		return "", "", err
	}
	return key, v.UploadID, nil
}

func CosPartUpload(r *http.Request) (string, error) {
	key := r.PostForm.Get("key")
	UploadID := r.PostForm.Get("uploadId")
	partNumber, err := strconv.Atoi(r.PostForm.Get("partNumber"))
	f, _, err := r.FormFile("file")
	if err != nil {
		return "", err
	}
	buf := bytes.NewBuffer(nil)
	io.Copy(buf, f)
	resp, err := GetCosClient().Object.UploadPart(
		context.Background(), key, UploadID, partNumber, bytes.NewReader(buf.Bytes()), nil,
	)
	if err != nil {
		return "", nil
	}
	return strings.Trim(resp.Header.Get("ETag"), "\""), nil
}

// CosChunkUploadFinish 分片上传的结束
func CosChunkUploadFinish(key, uploadId string, co []cos.Object) error {
	opt := &cos.CompleteMultipartUploadOptions{}
	opt.Parts = append(opt.Parts, co...)
	_, _, err := GetCosClient().Object.CompleteMultipartUpload(
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

func uploadCos(filePath string, file multipart.File) error {
	if filePath == "" || file == nil {
		tool.Logger.Error(errors.New(GetMsgByCode("zh", ParamErrCode)))
		return errors.New(GetMsgByCode("zh", ParamErrCode))
	}
	_, err := server.CosClient.Object.Put(
		context.Background(), filePath, file, nil,
	)
	if err != nil {
		tool.Logger.Error(err.Error())
		return err
	}
	return nil
}

func downloadCos(filePath string) ([]byte, error) {
	if filePath == "" {
		tool.Logger.Error(errors.New(GetMsgByCode("zh", ParamErrCode)))
		return nil, errors.New(GetMsgByCode("zh", ParamErrCode))
	}
	object, err := GetCosClient().Object.Get(
		context.Background(), filePath, &cos.ObjectGetOptions{},
	)
	if err != nil {
		tool.Logger.Error(err.Error())
		return nil, err
	}

	objectBytes, err := io.ReadAll(object.Body)
	if err != nil {
		tool.Logger.Error(err.Error())
		return nil, err
	}
	return objectBytes, nil
}
