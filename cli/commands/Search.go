package commands

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/NEXORA-Studios/Nova.ModDeps/core/api/modrinth"
	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/spf13/cobra"
)

func filterMatch(slice []string, targets []string) []string {
	var result []string
	targetSet := make(map[string]struct{}, len(targets))

	for _, t := range targets {
		targetSet[t] = struct{}{}
	}

	for _, item := range slice {
		if _, ok := targetSet[item]; ok {
			result = append(result, item)
		}
	}
	return result
}

var SearchCmd = &cobra.Command{
	Use:   "search <name> [page]",
	Short: "搜索 Modrinth 上的模组",
	Args:  cobra.RangeArgs(1, 2),
	Run: func(cmd *cobra.Command, args []string) {
		name := args[0]
		page := 0
		if len(args) == 2 {
			p, err := strconv.Atoi(args[1])
			if err != nil {
				panic("页码必须是数字")
			}
			page = p - 1
		}
		offset := page * 10
		result, err := modrinth.Search(name, offset)
		if err != nil {
			panic(err)
		}
		if len(result.Hits) == 0 {
			fmt.Println("未找到相关模组")
			return
		}

		fmt.Printf("共找到 %d 个结果，当前为 %d/%d 页\n", result.TotalHits, offset/10+1, (result.TotalHits+9)/10)

		t := table.NewWriter()
		t.SetOutputMirror(os.Stdout)
		t.AppendHeader(table.Row{"ID", "名称", "类型", "版本支持", "加载器支持"})

		for _, mod := range result.Hits {
			var modType []string
			if mod.ClientSide == "required" {
				modType = append(modType, "客户端")
			}
			if mod.ServerSide == "required" {
				modType = append(modType, "服务端")
			}

			var modVersionRange []string
			modVersionRange = append(modVersionRange, mod.Versions[0])
			if len(mod.Versions) > 1 {
				modVersionRange = append(modVersionRange, mod.Versions[len(mod.Versions)-1])
			}

			var modLoadersType []string
			modLoadersType = append(modLoadersType, filterMatch(mod.Categories, []string{"forge", "fabric", "neoforge", "quilt"})...)

			t.AppendRow(table.Row{
				mod.ProjectID,
				mod.Title,
				strings.Join(modType, " & "),
				strings.Join(modVersionRange, " ~ "),
				strings.Join(modLoadersType, " / "),
			})
		}

		t.SetStyle(table.StyleLight)
		t.Render()
	},
}
