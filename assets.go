package main

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"fmt"
	"mime"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/bootdotdev/learn-file-storage-s3-golang-starter/internal/database"
)

func (cfg apiConfig) ensureAssetsDir() error {
	if _, err := os.Stat(cfg.assetsRoot); os.IsNotExist(err) {
		return os.Mkdir(cfg.assetsRoot, 0755)
	}
	return nil
}

func (cfg apiConfig) getAssetDiskPath(assetPath string) string {
	return filepath.Join(cfg.assetsRoot, assetPath)
}

func (cfg apiConfig) getAssetURL(assetPath string) string {
	return fmt.Sprintf("http://localhost:%s/assets/%s", cfg.port, assetPath)
}

func (cfg apiConfig) getObjectURL(key string) string {
	return fmt.Sprintf("https://%s.s3.%s.amazonaws.com/%s", cfg.s3Bucket, cfg.s3Region, key)
}

func (cfg apiConfig) getObjectLocator(key string) string {
	return fmt.Sprintf("%s,%s", cfg.s3Bucket, key)
}

func (cfg apiConfig) dbVideoToSignedVideo(video database.Video) (database.Video, error) {
	if video.VideoURL == nil {
		return video, nil
	}

	parts := strings.Split(*video.VideoURL, ",")
	if len(parts) != 2 {
		return database.Video{}, errors.New("url does not have the right format")
	}

	bucket, key := parts[0], parts[1]
	expire := time.Hour * 24
	presignedUrl, err := generatePresignedURL(cfg.s3Client, bucket, key, expire)
	if err != nil {
		return database.Video{}, err
	}

	video.VideoURL = &presignedUrl
	return video, nil
}

func generatePresignedURL(s3Client *s3.Client, bucket, key string, expireTime time.Duration) (string, error) {
	presignClient := s3.NewPresignClient(s3Client)
	req, err := presignClient.PresignGetObject(context.Background(), &s3.GetObjectInput{
		Bucket: &bucket,
		Key:    &key,
	}, s3.WithPresignExpires(expireTime))
	if err != nil {
		return "", err
	}
	return req.URL, nil
}

func getMediaType(contentType string) (mediaType string, err error) {
	mediaType, _, err = mime.ParseMediaType(contentType)
	return
}

func getAssetPath(mediaType string) string {
	r := make([]byte, 32)
	if _, err := rand.Read(r); err != nil {
		panic("failed to generate random bytes")
	}
	encoded := base64.URLEncoding.EncodeToString(r)
	ext := mediaTypeToExt(mediaType)
	return fmt.Sprintf("%s%s", encoded, ext)
}

func mediaTypeToExt(mediaType string) string {

	parts := strings.Split(mediaType, "/")
	if len(parts) != 2 {
		return ".bin"
	}
	return "." + parts[1]
}
