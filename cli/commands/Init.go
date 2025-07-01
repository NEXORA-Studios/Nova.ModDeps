package commands

import (
	"fmt"

	"github.com/NEXORA-Studios/Nova.ModDeps/core/instance"
	"github.com/NEXORA-Studios/Nova.ModDeps/core/lock"
	"github.com/NEXORA-Studios/Nova.ModDeps/core/meta"
	"github.com/spf13/cobra"
)

var InitCmd = &cobra.Command{
	Use:   "init",
	Short: "初始化 Nova.ModDeps 包文件，需要在 Minecraft 实例根目录下运行",
	Run: func(cmd *cobra.Command, args []string) {
		minecraftVersion, modLoader, err := instance.LoadInstance()
		if err != nil {
			logger.Fatal(fmt.Sprintf("初始化失败: %v", err))
		}
		if minecraftVersion == "" {
			logger.Fatal("未找到 Minecraft 版本，版本 JSON 存在吗？")
		}

		var metaFunc = meta.MetaFunctions{}
		metaFunc.Write(meta.IModPackageJson{
			InternalVersion:  1,
			InternalPlatform: "modrinth",
			MinecraftVersion: minecraftVersion,
			ModLoader:        modLoader,
			Mods:             []meta.IModItem{},
		})

		lock.Write(lock.ILockFile{
			LockFileVersion: 1,
			BaseDir:         "mods",
			Pending:         []lock.ILockModItem{},
			Installed:       []lock.ILockModItem{},
			NeedRemove:      []lock.ILockModItem{},
		})
	},
}
