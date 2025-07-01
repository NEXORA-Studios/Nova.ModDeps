package lock

import (
	"encoding/json"
	"fmt"
	"os"
)

func Write(uLockFile ILockFile) {
	lockFile, err := json.Marshal(uLockFile)
	if err != nil {
		logger.Error(fmt.Sprintf("写入 mod.lock.json 失败: %v", err))
		return
	}
	os.WriteFile("mod.lock.json", lockFile, 0644)
}
