package fs

import (
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/NEXORA-Studios/Nova.ModDeps/core"
)

var logger = core.Logger{}

func DownloadFile(url string, path string) {
	resp, err := http.Get(url)
	if err != nil {
		logger.Error(fmt.Sprintf("下载文件 %s (从 %s) 失败: %v", path, url, err))
		return
	}
	defer resp.Body.Close()

	out, err := os.Create(path)
	if err != nil {
		logger.Error(fmt.Sprintf("创建文件 %s 失败: %v", path, err))
		return
	}
	defer out.Close()

	_, err = io.Copy(out, resp.Body)
	if err != nil {
		logger.Error(fmt.Sprintf("写入文件 %s 失败: %v", path, err))
		return
	}
}
