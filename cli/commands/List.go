package commands

import (
	"fmt"
	"os"
	"strings"

	"github.com/NEXORA-Studios/Nova.ModDeps/cli/utils"
	"github.com/NEXORA-Studios/Nova.ModDeps/core/api/modrinth"
	"github.com/NEXORA-Studios/Nova.ModDeps/core/lock"
	"github.com/NEXORA-Studios/Nova.ModDeps/core/meta"
	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/spf13/cobra"
)

var ListCmd = &cobra.Command{
	Use:   "list",
	Short: "列出当前已添加的所有 Mod",
	Run: func(cmd *cobra.Command, args []string) {
		metaFunc := meta.MetaFunctions{}
		modPackage, err := metaFunc.Read()
		if err != nil {
			logger.Fatal(fmt.Sprintf("读取 mod.package.json 失败: %v\n", err))
		}
		if len(modPackage.Mods) > 0 {
			t := table.NewWriter()
			t.SetOutputMirror(os.Stdout)
			t.AppendHeader(table.Row{"项目 ID", "版本 ID", "被依赖", "名称", "版本", "状态"})
			for _, mod := range modPackage.Mods {
				var versionName string
				projectMetaData, err := modrinth.GetProjectMetadata(mod.ID)
				versionMetaData, err := modrinth.GetVersionMetadata(mod.Version)
				if err != nil {
					logger.Warn(fmt.Sprintf("获取 %s 版本元数据失败: %v\n", mod.Version, err))
					versionName = "未知"
				} else {
					versionName = projectMetaData.Title
				}
				status := lock.GetStatus(mod.ID, mod.Version)
				switch status {
				case "installed":
					status = fmt.Sprintf("%s[已安装]%s", utils.ColorGreen, utils.ColorReset)
				case "pending":
					status = fmt.Sprintf("%s[待安装]%s", utils.ColorYellow, utils.ColorReset)
				case "needRemove":
					status = fmt.Sprintf("%s[待卸载]%s", utils.ColorRed, utils.ColorReset)
				}
				t.AppendRow(table.Row{
					mod.ID,
					mod.Version,
					strings.Join(mod.RequiredBy, ", "),
					versionName,
					versionMetaData.Name,
					status,
				})
			}
			t.SetStyle(table.StyleLight)
			fmt.Printf("客户端版本: %s | ModLoader: [%s] | Mod 总数: %d\n", modPackage.MinecraftVersion, strings.Join(modPackage.ModLoader, ", "), len(modPackage.Mods))
			t.Render()
		}
	},
}
