package commands

import (
	"fmt"
	"path/filepath"

	"github.com/NEXORA-Studios/Nova.ModDeps/core/api/modrinth"
	"github.com/NEXORA-Studios/Nova.ModDeps/core/lock"
	"github.com/NEXORA-Studios/Nova.ModDeps/core/meta"
	"github.com/spf13/cobra"
)

var LoadCmd = &cobra.Command{
	Use:   "load <from>",
	Short: "在 mod.package.json 和 mod.lock.json 之间互相补充加载",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		from := args[0]

		if from != "package" && from != "lock" {
			logger.Error("参数 'from' 必须是 'package' 或 'lock'")
			return
		}

		metaFunc := meta.MetaFunctions{}

		if from == "package" {
			// 从 package.json 加载到 lock.json
			err := loadFromPackageToLock(metaFunc)
			if err != nil {
				logger.Error(fmt.Sprintf("从 package.json 加载到 lock.json 失败: %v", err))
				return
			}
			logger.Info("已从 mod.package.json 加载到 mod.lock.json")
		} else {
			// 从 lock.json 加载到 package.json
			err := loadFromLockToPackage(metaFunc)
			if err != nil {
				logger.Error(fmt.Sprintf("从 lock.json 加载到 package.json 失败: %v", err))
				return
			}
			logger.Info("已从 mod.lock.json 加载到 mod.package.json")
		}
	},
}

// 从 package.json 加载到 lock.json
func loadFromPackageToLock(metaFunc meta.MetaFunctions) error {
	// 读取 package.json
	modPackage, err := metaFunc.Read()
	if err != nil {
		return fmt.Errorf("读取 mod.package.json 失败: %v", err)
	}

	// 读取当前的 lock.json
	lockFile, err := lock.Read()
	if err != nil {
		return fmt.Errorf("读取 mod.lock.json 失败: %v", err)
	}

	// 清空 pending 列表，重新从 package.json 构建
	lockFile.Pending = []lock.ILockModItem{}

	// 遍历 package.json 中的所有 mod
	for _, mod := range modPackage.Mods {
		// 检查是否已经在 installed 中
		found := false
		for _, installed := range lockFile.Installed {
			if installed.ID == mod.ID && installed.Version == mod.Version {
				found = true
				break
			}
		}

		if !found {
			// 获取 mod 的元数据以确定文件路径
			versionMetadata, err := modrinth.GetVersionMetadata(mod.Version)
			if err != nil {
				logger.Warn(fmt.Sprintf("获取 mod %s 版本元数据失败: %v，添加到 pending", mod.ID, err))
				// 无法获取元数据时，使用占位符信息添加到 pending
				lockItem := lock.ILockModItem{
					ID:      mod.ID,
					Version: mod.Version,
					Path:    filepath.Join(lockFile.BaseDir, fmt.Sprintf("%s-%s.jar", mod.ID, mod.Version)),
					Uri:     "",
					Sha512:  "",
				}
				lockFile.Pending = append(lockFile.Pending, lockItem)
				continue
			}

			// 确定文件路径和元数据
			var filePath string
			var fileURL string
			var fileSha512 string
			for _, file := range versionMetadata.Files {
				if file.Primary {
					filePath = filepath.Join(lockFile.BaseDir, file.Filename)
					fileURL = file.URL
					fileSha512 = file.Hashes.SHA512
					break
				}
			}

			// 检查文件是否存在
			if fileExists(filePath) {
				// 文件存在，添加到 installed
				lockItem := lock.ILockModItem{
					ID:      mod.ID,
					Version: mod.Version,
					Path:    filePath,
					Uri:     fileURL,    // 保留下载链接以备需要
					Sha512:  fileSha512, // 保留校验和
				}
				lockFile.Installed = append(lockFile.Installed, lockItem)
				logger.Info(fmt.Sprintf("发现已存在的文件，将 mod %s 添加到 installed: %s", mod.ID, filePath))
			} else {
				// 文件不存在，添加到 pending
				lockItem := lock.ILockModItem{
					ID:      mod.ID,
					Version: mod.Version,
					Path:    filePath,
					Uri:     fileURL,
					Sha512:  fileSha512,
				}
				lockFile.Pending = append(lockFile.Pending, lockItem)
				logger.Info(fmt.Sprintf("文件不存在，将 mod %s 添加到 pending: %s", mod.ID, filePath))
			}
		}
	}

	// 保存更新后的 lock.json
	lock.Write(lockFile)
	return nil
}

// 从 lock.json 加载到 package.json
func loadFromLockToPackage(metaFunc meta.MetaFunctions) error {
	// 读取 lock.json
	lockFile, err := lock.Read()
	if err != nil {
		return fmt.Errorf("读取 mod.lock.json 失败: %v", err)
	}

	// 读取当前的 package.json
	modPackage, err := metaFunc.Read()
	if err != nil {
		return fmt.Errorf("读取 mod.package.json 失败: %v", err)
	}

	// 清空 mods 列表，重新从 lock.json 构建
	modPackage.Mods = []meta.IModItem{}

	// 从 installed 和 pending 中获取所有 mod
	allMods := make(map[string]lock.ILockModItem)

	// 添加 installed 中的 mod
	for _, mod := range lockFile.Installed {
		allMods[mod.ID] = mod
	}

	// 添加 pending 中的 mod（如果不在 installed 中）
	for _, mod := range lockFile.Pending {
		if _, exists := allMods[mod.ID]; !exists {
			allMods[mod.ID] = mod
		}
	}

	// 为每个 mod 获取元数据并添加到 package.json
	for _, lockMod := range allMods {
		// 从 Modrinth API 获取 mod 的名称
		projectMeta, err := modrinth.GetProjectMetadata(lockMod.ID)
		if err != nil {
			logger.Warn(fmt.Sprintf("获取 mod %s 的元数据失败: %v，使用占位符名称", lockMod.ID, err))
			projectMeta.Title = fmt.Sprintf("Mod-%s", lockMod.ID)
		}

		// 添加到 package.json
		err = metaFunc.UpsertMod(lockMod.ID, lockMod.Version, projectMeta.Title, []string{})
		if err != nil {
			return fmt.Errorf("添加 mod %s 到 package.json 失败: %v", lockMod.ID, err)
		}
	}

	return nil
}
