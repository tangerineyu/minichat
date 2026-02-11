package oss

import (
	"context"
	"io"
)

type OSSInterface interface {
	// UploadFile 上传文件到OSS，返回文件的访问URL。
	// contentType 建议传入：image/jpeg、image/png 等，便于浏览器与 CDN 正确识别。
	UploadFile(ctx context.Context, objectKey string, contentType string, r io.Reader) (string, error)
	// DeleteFile 从OSS删除文件
	DeleteFile(ctx context.Context, objectKey string) error
}

type OSSConfig struct {
	Endpoint        string
	AccessKeyID     string
	AccessKeySecret string
	BucketName      string
}
