package commands

import (
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
)

var CacheCleanCmd = &cobra.Command{
	Use:   "cache-clean",
	Short: "清理 API 请求缓存",
	Run: func(cmd *cobra.Command, args []string) {
		exePath, _ := os.Executable()
		dir := filepath.Join(filepath.Dir(exePath), "temp", "ndm")
		err := os.RemoveAll(dir)
		if err != nil {
			logger.Error("清理缓存失败: " + err.Error())
			return
		}
		logger.Info("缓存已清理: " + dir)
	},
}
