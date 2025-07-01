package commands

import (
	"fmt"

	"github.com/NEXORA-Studios/Nova.ModDeps/core"
	"github.com/NEXORA-Studios/Nova.ModDeps/core/instance"
	"github.com/NEXORA-Studios/Nova.ModDeps/core/lock"
	"github.com/NEXORA-Studios/Nova.ModDeps/core/meta"
	"github.com/spf13/cobra"
)

var initForce bool

var InitCmd = &cobra.Command{
	Use:   "init [package|lock]",
	Short: "初始化 Nova.ModDeps 包文件，需要在 Minecraft 实例根目录下运行",
	Args:  cobra.MaximumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		var initType string
		if len(args) > 0 {
			initType = args[0]
			if initType != "package" && initType != "lock" {
				logger.Error("参数必须是 'package' 或 'lock'")
				return
			}
		}

		// 如果没有指定类型，则初始化两个文件
		if initType == "" {
			initType = "both"
		}

		// 检查文件是否存在并验证是否需要强制覆盖
		if !initForce {
			if initType == "package" || initType == "both" {
				if fileExists(core.GetModPackageJsonPath()) {
					logger.Error("mod.package.json 已存在，使用 --force 标志强制覆盖")
					return
				}
			}
			if initType == "lock" || initType == "both" {
				if fileExists("mod.lock.json") {
					logger.Error("mod.lock.json 已存在，使用 --force 标志强制覆盖")
					return
				}
			}
		}

		// 获取 Minecraft 版本信息（仅在需要初始化 package.json 时）
		var minecraftVersion string
		var modLoader []string
		var err error

		if initType == "package" || initType == "both" {
			minecraftVersion, modLoader, err = instance.LoadInstance()
			if err != nil {
				logger.Fatal(fmt.Sprintf("初始化失败: %v", err))
			}
			if minecraftVersion == "" {
				logger.Fatal("未找到 Minecraft 版本，版本 JSON 存在吗？")
			}
		}

		// 初始化 package.json
		if initType == "package" || initType == "both" {
			var metaFunc = meta.MetaFunctions{}
			err := metaFunc.Write(meta.IModPackageJson{
				InternalVersion:  1,
				InternalPlatform: "modrinth",
				MinecraftVersion: minecraftVersion,
				ModLoader:        modLoader,
				Mods:             []meta.IModItem{},
			})
			if err != nil {
				logger.Error(fmt.Sprintf("初始化 mod.package.json 失败: %v", err))
			} else {
				logger.Info("已初始化 mod.package.json")
			}
		}

		// 初始化 lock.json
		if initType == "lock" || initType == "both" {
			lock.Write(lock.ILockFile{
				LockFileVersion: 1,
				BaseDir:         "mods",
				Pending:         []lock.ILockModItem{},
				Installed:       []lock.ILockModItem{},
				NeedRemove:      []lock.ILockModItem{},
			})
			logger.Info("已初始化 mod.lock.json")
		}
	},
}

func init() {
	InitCmd.Flags().BoolVarP(&initForce, "force", "f", false, "强制覆盖已存在的文件")
}
