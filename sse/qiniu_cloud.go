package sse

import (
	"context"
	"errors"
	"io"
	"wrap-midjourney/initialization"

	"github.com/qiniu/go-sdk/v7/auth/qbox"
	"github.com/qiniu/go-sdk/v7/storage"
)

type UploadResult struct {
	Key             string
	Hash            string
	PublicAccessURL string
}

func qiniu_cloud(key string, data io.Reader, size int64) (*UploadResult, error) {
	conf := initialization.GetConfig()
	if conf == nil || conf.QiNiuConfig == nil {
		return nil, errors.New("qiniu config is nil")
	}
	bucket := conf.QiNiuConfig.Bucket
	accessKey := conf.QiNiuConfig.AccessKey
	secretKey := conf.QiNiuConfig.SecretKey
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
