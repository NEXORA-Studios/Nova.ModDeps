package api

import (
	"crypto/sha256"
	"encoding/hex"
	"io"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
)

type Requester struct{}

func cacheDir() string {
	exePath, _ := os.Executable()
	dir := filepath.Join(filepath.Dir(exePath), "temp", "ndm")
	os.MkdirAll(dir, 0755)
	return dir
}

func cacheFileName(key string) string {
	h := sha256.Sum256([]byte(key))
	return filepath.Join(cacheDir(), hex.EncodeToString(h[:]))
}

func (r *Requester) Get(path string, query map[string]string) (string, error) {
	baseURL := "https://api.modrinth.com/v2"
	u, err := url.Parse(baseURL + path)
	if err != nil {
		return "", err
	}
	q := u.Query()
	for k, v := range query {
		q.Set(k, v)
	}
	u.RawQuery = q.Encode()
	cacheKey := u.String()
	cacheFile := cacheFileName(cacheKey)

	// 优先查磁盘缓存
	if data, err := os.ReadFile(cacheFile); err == nil {
		return string(data), nil
	}

	// 未命中缓存，请求
	response, err := http.Get(cacheKey)
	if err != nil {
		return "", err
	}
	defer response.Body.Close()

	body, err := io.ReadAll(response.Body)
	if err != nil {
		return "", err
	}

	// 写入磁盘缓存
	os.WriteFile(cacheFile, body, 0644)

	return string(body), nil
}
