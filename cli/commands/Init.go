package commands

import (
	"github.com/NEXORA-Studios/Nova.ModDeps/core/instance"
	"github.com/NEXORA-Studios/Nova.ModDeps/core/meta"
	"github.com/spf13/cobra"
)

var InitCmd = &cobra.Command{
	Use:   "init",
	Short: "初始化 Nova.ModDeps 包文件，需要在 Minecraft 实例根目录下运行",
	Run: func(cmd *cobra.Command, args []string) {
		minecraftVersion, modLoader, err := instance.LoadInstance()
		if err != nil {
			panic(err)
		}
		if minecraftVersion == "" {
			panic("未找到 Minecraft 版本，版本 JSON 存在吗？")
		}
		
		var metaFunc = meta.MetaFunctions{}
		metaFunc.Write(meta.IModPackageJson{
			InternalVersion:  "1.0.0",
			InternalPlatform: "modrinth",
			MinecraftVersion: minecraftVersion,
			ModLoader:        modLoader,
			Mods:             []meta.IModItem{},
		})
	},
}
