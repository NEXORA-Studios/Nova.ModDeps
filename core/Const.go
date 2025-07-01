package core

import (
	"fmt"
	"os"
	"path/filepath"
)

var logger = Logger{}

func GetModPackageJsonPath() string {
	cwd, err := os.Getwd()
	if err != nil {
		logger.Fatal(fmt.Sprintf("尝试使用当前工作目录获取 mod.package.json 路径失败: %v", err))
	}
	return filepath.Join(cwd, "mod.package.json")
}
