package fs

import (
	"fmt"
	"os"
)

func RemoveFile(path string) {
	err := os.Remove(path)
	if err != nil {
		logger.Error(fmt.Sprintf("删除文件 %s 失败: %v", path, err))
		return
	}
}
