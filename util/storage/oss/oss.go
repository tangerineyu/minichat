package oss

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net/url"
	"path"
	"strings"

	"github.com/aliyun/aliyun-oss-go-sdk/oss"
)

type AliYunOSS struct {
	cfg    *OSSConfig
	client *oss.Client
}

func (a AliYunOSS) UploadFile(ctx context.Context, objectKey string, contentType string, r io.Reader) (string, error) {
	if a.client == nil || a.cfg == nil {
		return "", errors.New("oss client not initialized")
	}
	objectKey = strings.TrimSpace(objectKey)
	if objectKey == "" {
		return "", errors.New("objectKey is empty")
	}
	if r == nil {
		return "", errors.New("reader is nil")
	}

	bucket, err := a.client.Bucket(a.cfg.BucketName)
	if err != nil {
		return "", fmt.Errorf("get oss bucket failed: %w", err)
	}

	var opts []oss.Option
	opts = append(opts, oss.WithContext(ctx))
	if strings.TrimSpace(contentType) != "" {
		opts = append(opts, oss.ContentType(contentType))
	}

	if err := bucket.PutObject(objectKey, r, opts...); err != nil {
		return "", fmt.Errorf("put oss object failed: %w", err)
	}

	return a.objectURL(objectKey)
}

func (a AliYunOSS) DeleteFile(ctx context.Context, objectKey string) error {
	if a.client == nil || a.cfg == nil {
		return errors.New("oss client not initialized")
	}
	objectKey = strings.TrimSpace(objectKey)
	if objectKey == "" {
		return errors.New("objectKey is empty")
	}

	bucket, err := a.client.Bucket(a.cfg.BucketName)
	if err != nil {
		return fmt.Errorf("get oss bucket failed: %w", err)
	}
	if err := bucket.DeleteObject(objectKey, oss.WithContext(ctx)); err != nil {
		return fmt.Errorf("delete oss object failed: %w", err)
	}
	return nil
}

func (a AliYunOSS) objectURL(objectKey string) (string, error) {
	endpoint := strings.TrimSpace(a.cfg.Endpoint)
	if endpoint == "" {
		return "", errors.New("oss endpoint is empty")
	}

	if !strings.HasPrefix(endpoint, "http://") && !strings.HasPrefix(endpoint, "https://") {
		endpoint = "https://" + endpoint
	}

	u, err := url.Parse(endpoint)
	if err != nil {
		return "", fmt.Errorf("invalid oss endpoint: %w", err)
	}

	host := u.Hostname()
	if host == "" {
		host = u.Host
	}
	u.Host = a.cfg.BucketName + "." + host
	u.Path = path.Join("/", objectKey)
	return u.String(), nil
}

var _ OSSInterface = (*AliYunOSS)(nil)

func NewAliYunOSS(cfg *OSSConfig) (*AliYunOSS, error) {
	if cfg == nil {
		return nil, errors.New("oss config is nil")
	}
	if strings.TrimSpace(cfg.Endpoint) == "" {
		return nil, errors.New("oss endpoint is empty")
	}
	if strings.TrimSpace(cfg.AccessKeyID) == "" || strings.TrimSpace(cfg.AccessKeySecret) == "" {
		return nil, errors.New("oss access key is empty")
	}
	if strings.TrimSpace(cfg.BucketName) == "" {
		return nil, errors.New("oss bucket name is empty")
	}

	client, err := oss.New(cfg.Endpoint, cfg.AccessKeyID, cfg.AccessKeySecret)
	if err != nil {
		return nil, err
	}
	return &AliYunOSS{
		cfg:    cfg,
		client: client,
	}, nil
}
