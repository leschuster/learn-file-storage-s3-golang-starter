package main

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"mime"
	"os"
	"path/filepath"
	"strings"
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
