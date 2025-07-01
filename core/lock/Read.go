package lock

import (
	"encoding/json"
	"fmt"
	"os"
)

func Read() (ILockFile, error) {
	lockFile, err := os.ReadFile("mod.lock.json")
	if err != nil {
		logger.Error(fmt.Sprintf("读取 mod.lock.json 失败: %v", err))
		return ILockFile{}, err
	}
	var uLockFile ILockFile
	err = json.Unmarshal(lockFile, &uLockFile)
	if err != nil {
		logger.Error(fmt.Sprintf("解析 mod.lock.json 失败: %v", err))
		return ILockFile{}, err
	}
	return uLockFile, nil
}

func GetBaseDir() string {
	lockFile, err := Read()
	if err != nil {
		return "mods"
	}
	return lockFile.BaseDir
}
