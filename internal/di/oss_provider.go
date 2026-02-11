package di

import (
	"minichat/internal/config"
	ossutil "minichat/util/storage/oss"
)

// ProvideOSS builds an OSSInterface from app config.
// If OSS is not configured (endpoint/bucket empty), it returns (nil, nil),
// allowing the app to run without OSS features.
func ProvideOSS(cfg config.AppConfig) (ossutil.OSSInterface, error) {
	c := cfg.OSS
	if c.Endpoint == "" || c.BucketName == "" || c.AccessKeyID == "" || c.AccessKeySecret == "" {
		return nil, nil
	}
	return ossutil.NewAliYunOSS(&ossutil.OSSConfig{
		Endpoint:        c.Endpoint,
		AccessKeyID:     c.AccessKeyID,
		AccessKeySecret: c.AccessKeySecret,
		BucketName:      c.BucketName,
	})
}
