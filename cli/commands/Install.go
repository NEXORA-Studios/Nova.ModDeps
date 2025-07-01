package commands

import (
	"fmt"
	"os"
	"path/filepath"

	"crypto/sha512"
	"io"

	"github.com/NEXORA-Studios/Nova.ModDeps/core"
	"github.com/NEXORA-Studios/Nova.ModDeps/core/fs"
	"github.com/NEXORA-Studios/Nova.ModDeps/core/lock"
	"github.com/NEXORA-Studios/Nova.ModDeps/core/meta"
	"github.com/spf13/cobra"
)

var InstallCmd = &cobra.Command{
	Use:   "install",
	Short: "同步 mods 文件夹与 lockfile",
	Run: func(cmd *cobra.Command, args []string) {
		logger := core.Logger{}
		// 1. 确认 basedir 文件夹存在
		baseDir := filepath.Join(filepath.Dir(core.GetModPackageJsonPath()), lock.GetBaseDir())
		if _, err := os.Stat(baseDir); os.IsNotExist(err) {
			os.MkdirAll(baseDir, 0755)
			logger.Info("已创建 mods 目录: " + baseDir)
		}

		uLockFile, err := lock.Read()
		if err != nil {
			logger.Error("读取 lockfile 失败: " + err.Error())
			return
		}

		// 2. 校验 installed
		for _, item := range uLockFile.Installed {
			fileExists := fileExists(item.Path)
			valid := false
			if fileExists {
				valid = checkSha512(item.Path, item.Sha512)
			}
			if !fileExists || !valid {
				fs.RemoveFile(item.Path)
				lock.RePending(item.ID, item.Version)
				logger.Warn("已回退损坏或缺失的 mod: " + item.Path)
			}
		}

		// 3. needRemove
		metaFunc := meta.MetaFunctions{}

		for _, item := range uLockFile.NeedRemove {
			fs.RemoveFile(item.Path)
			lock.Remove(item.ID, item.Version)
			err = metaFunc.RemoveMod(item.ID)
			if err != nil {
				logger.Error("从 mod.package.json 移除 mod 失败: " + err.Error())
			}
			logger.Info("已删除并移除 needRemove 项: " + item.Path)
		}

		// 4. pending
		for _, item := range uLockFile.Pending {
			fs.DownloadFile(item.Uri, item.Path)
			if checkSha512(item.Path, item.Sha512) {
				lock.AddInstalled(item.ID, item.Version)
				logger.Info("已安装 mod: " + item.Path)
			} else {
				fs.RemoveFile(item.Path)
				logger.Error("下载的 mod 校验失败: " + item.Path)
			}
		}
	},
}

// 工具函数
func fileExists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}

func checkSha512(path, expect string) bool {
	f, err := os.Open(path)
	if err != nil {
		return false
	}
	defer f.Close()
	h := sha512.New()
	if _, err := io.Copy(h, f); err != nil {
		return false
	}
	sum := h.Sum(nil)
	return fmt.Sprintf("%x", sum) == expect
}
