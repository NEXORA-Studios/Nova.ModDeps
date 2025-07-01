package commands

import (
	"fmt"
	"os"
	"strings"

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
			fmt.Fprintf(os.Stderr, "读取 mod.package.json 失败: %v\n", err)
			os.Exit(1)
		}
		t := table.NewWriter()
		t.SetOutputMirror(os.Stdout)
		t.AppendHeader(table.Row{"ID", "版本", "名称", "被依赖"})
		for _, mod := range modPackage.Mods {
			t.AppendRow(table.Row{
				mod.ID,
				mod.Version,
				mod.Name,
				strings.Join(mod.RequiredBy, ", "),
			})
		}
		t.SetStyle(table.StyleLight)
		t.Render()
	},
}
