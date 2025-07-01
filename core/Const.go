package core

import (
	"fmt"
	"os"
	"path/filepath"
)

const useCwd = false

func GetModPackageJsonPath() string {
	if useCwd {
		cwd, err := os.Getwd()
		if err != nil {
			panic(fmt.Sprintf("尝试使用当前工作目录获取 mod.package.json 路径失败: %v", err))
		}
		return filepath.Join(cwd, "mod.package.json")
	} else {
		exe, err := os.Executable()
		if err != nil {
			panic(fmt.Sprintf("尝试使用当前可执行文件路径获取 mod.package.json 路径失败: %v", err))
		}
		exeDir := filepath.Dir(exe)
		return filepath.Join(exeDir, "mod.package.json")
	}
}
