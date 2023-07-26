package sse

import (
	"context"
	"io"

	"github.com/qiniu/go-sdk/v7/auth/qbox"
	"github.com/qiniu/go-sdk/v7/storage"
)

type UploadResult struct {
	Key             string
	Hash            string
	PublicAccessURL string
}

func qiniu_cloud(key string, data io.Reader, size int64) (*UploadResult, error) {
	bucket := "topschat-mj"
	accessKey := "hM2NyESPz0rIRL6YbyADKMA-dzOO8oWqK9CYHmOy"
	secretKey := "RX8sysDxa7EC4h0cKGs3Syp9AR3lmyWjFj5OTdXO"
	putPolicy := storage.PutPolicy{
		Scope: bucket,
	}
	mac := qbox.NewMac(accessKey, secretKey)
	upToken := putPolicy.UploadToken(mac)

	// 构建表单上传的对象
	formUploader := storage.NewFormUploader(&storage.Config{
		Zone:          &storage.ZoneHuadong,
		UseHTTPS:      true,
		UseCdnDomains: false,
	})

	ret := storage.PutRet{}
	// 可选配置
	putExtra := storage.PutExtra{
		Params: map[string]string{
			"x:name": "qiniuyun logo",
		},
	}
	err := formUploader.Put(context.Background(), &ret, upToken, key, data, size, &putExtra)
	if err != nil {
		return nil, err
	}
	publicAccessURL := storage.MakePublicURL("https://oss.topschat.com/", key)
	return &UploadResult{
		Key:             ret.Key,
		Hash:            ret.Hash,
		PublicAccessURL: publicAccessURL,
	}, nil
}
