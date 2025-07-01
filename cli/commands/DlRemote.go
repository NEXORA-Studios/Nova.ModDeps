package commands

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"

	"github.com/NEXORA-Studios/Nova.ModDeps/core"
	"github.com/NEXORA-Studios/Nova.ModDeps/core/lock"
	"github.com/NEXORA-Studios/Nova.ModDeps/core/meta"
	"github.com/spf13/cobra"
)

var dlremoteForce bool
var dlremoteBackup bool

var DlremoteCmd = &cobra.Command{
	Use:   "dlremote <url>",
	Short: "从远程获取 mod.package.json 和 mod.lock.json 的内容",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		baseURL := args[0]

		// 确保 URL 以 / 结尾
		if baseURL[len(baseURL)-1] != '/' {
			baseURL += "/"
		}

		// 构建完整的 URL
		packageURL := baseURL + "mod.package.json"
		lockURL := baseURL + "mod.lock.json"

		logger.Info(fmt.Sprintf("正在从远程获取配置文件..."))
		logger.Info(fmt.Sprintf("Package URL: %s", packageURL))
		logger.Info(fmt.Sprintf("Lock URL: %s", lockURL))

		// 检查本地文件是否存在
		packageExists := fileExists(core.GetModPackageJsonPath())
		lockExists := fileExists("mod.lock.json")

		// 如果需要备份且文件存在
		if dlremoteBackup {
			if packageExists {
				err := backupFile(core.GetModPackageJsonPath(), "mod.package.json.backup")
				if err != nil {
					logger.Error(fmt.Sprintf("备份 mod.package.json 失败: %v", err))
				} else {
					logger.Info("已备份 mod.package.json")
				}
			}
			if lockExists {
				err := backupFile("mod.lock.json", "mod.lock.json.backup")
				if err != nil {
					logger.Error(fmt.Sprintf("备份 mod.lock.json 失败: %v", err))
				} else {
					logger.Info("已备份 mod.lock.json")
				}
			}
		}

		// 检查是否需要强制覆盖
		if !dlremoteForce {
			if packageExists {
				logger.Error("mod.package.json 已存在，使用 --force 标志强制覆盖")
				return
			}
			if lockExists {
				logger.Error("mod.lock.json 已存在，使用 --force 标志强制覆盖")
				return
			}
		}

		// 下载 mod.package.json
		logger.Info("正在下载 mod.package.json...")
		err := downloadAndValidatePackage(packageURL)
		if err != nil {
			logger.Error(fmt.Sprintf("下载 mod.package.json 失败: %v", err))
			return
		}
		logger.Info("已成功下载 mod.package.json")

		// 下载 mod.lock.json
		logger.Info("正在下载 mod.lock.json...")
		err = downloadAndValidateLock(lockURL)
		if err != nil {
			logger.Error(fmt.Sprintf("下载 mod.lock.json 失败: %v", err))
			return
		}
		logger.Info("已成功下载 mod.lock.json")

		logger.Info("远程配置文件下载完成！")
		logger.Info("提示：使用 'novadm load package' 来同步本地文件状态")
	},
}

func init() {
	DlremoteCmd.Flags().BoolVarP(&dlremoteForce, "force", "f", false, "强制覆盖已存在的文件")
	DlremoteCmd.Flags().BoolVarP(&dlremoteBackup, "backup", "b", false, "在覆盖前备份现有文件")
}

// 下载并验证 package.json
func downloadAndValidatePackage(url string) error {
	// 下载文件
	resp, err := http.Get(url)
	if err != nil {
		return fmt.Errorf("HTTP 请求失败: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return fmt.Errorf("HTTP 状态码错误: %d", resp.StatusCode)
	}

	// 读取响应内容
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("读取响应内容失败: %v", err)
	}

	// 验证 JSON 格式
	var packageData meta.IModPackageJson
	err = json.Unmarshal(body, &packageData)
	if err != nil {
		return fmt.Errorf("JSON 格式验证失败: %v", err)
	}

	// 验证必要字段
	if packageData.InternalVersion == 0 {
		return fmt.Errorf("缺少 __version__ 字段")
	}
	if packageData.InternalPlatform == "" {
		return fmt.Errorf("缺少 __platform__ 字段")
	}
	if packageData.MinecraftVersion == "" {
		return fmt.Errorf("缺少 minecraftVersion 字段")
	}

	// 写入文件
	err = os.WriteFile(core.GetModPackageJsonPath(), body, 0644)
	if err != nil {
		return fmt.Errorf("写入文件失败: %v", err)
	}

	return nil
}

// 下载并验证 lock.json
func downloadAndValidateLock(url string) error {
	// 下载文件
	resp, err := http.Get(url)
	if err != nil {
		return fmt.Errorf("HTTP 请求失败: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return fmt.Errorf("HTTP 状态码错误: %d", resp.StatusCode)
	}

	// 读取响应内容
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("读取响应内容失败: %v", err)
	}

	// 验证 JSON 格式
	var lockData lock.ILockFile
	err = json.Unmarshal(body, &lockData)
	if err != nil {
		return fmt.Errorf("JSON 格式验证失败: %v", err)
	}

	// 验证必要字段
	if lockData.LockFileVersion == 0 {
		return fmt.Errorf("缺少 lockfileVersion 字段")
	}
	if lockData.BaseDir == "" {
		return fmt.Errorf("缺少 basedir 字段")
	}

	// 写入文件
	err = os.WriteFile("mod.lock.json", body, 0644)
	if err != nil {
		return fmt.Errorf("写入文件失败: %v", err)
	}

	return nil
}

// 备份文件
func backupFile(sourcePath, backupName string) error {
	// 读取源文件
	sourceData, err := os.ReadFile(sourcePath)
	if err != nil {
		return err
	}

	// 写入备份文件
	backupPath := filepath.Join(filepath.Dir(sourcePath), backupName)
	err = os.WriteFile(backupPath, sourceData, 0644)
	if err != nil {
		return err
	}

	return nil
}
