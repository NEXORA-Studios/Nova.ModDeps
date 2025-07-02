package commands

import (
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/NEXORA-Studios/Nova.ModDeps/core/api/modrinth"
	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/spf13/cobra"
)

var searchJson bool

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
				logger.Error("页码必须是数字")
				return
			}
			page = p - 1
		}
		offset := page * 10
		result, err := modrinth.Search(name, offset)
		if err != nil {
			logger.Error(fmt.Sprintf("搜索失败: %v", err))
			return
		}
		if len(result.Hits) == 0 {
			logger.Warn("未找到相关模组")
			return
		}

		if searchJson {
			// JSON 输出
			outputSearchJSON(result, offset)
		} else {
			// 表格输出
			outputSearchTable(result, offset)
		}
	},
}

func init() {
	SearchCmd.Flags().BoolVarP(&searchJson, "json", "j", false, "以 JSON 格式输出")
}

// JSON 输出函数
func outputSearchJSON(result modrinth.IMrSearchResponse, offset int) {
	type SearchModInfo struct {
		ID            string   `json:"id"`
		Name          string   `json:"name"`
		Type          []string `json:"type"`
		VersionRange  []string `json:"versionRange"`
		LoaderSupport []string `json:"loaderSupport"`
		Description   string   `json:"description"`
		Author        string   `json:"author"`
		Categories    []string `json:"categories"`
	}

	type SearchOutput struct {
		TotalHits    int             `json:"totalHits"`
		CurrentPage  int             `json:"currentPage"`
		TotalPages   int             `json:"totalPages"`
		Results      []SearchModInfo `json:"results"`
	}

	output := SearchOutput{
		TotalHits:   result.TotalHits,
		CurrentPage: offset/10 + 1,
		TotalPages:  (result.TotalHits + 9) / 10,
		Results:     []SearchModInfo{},
	}

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

		modInfo := SearchModInfo{
			ID:            mod.ProjectID,
			Name:          mod.Title,
			Type:          modType,
			VersionRange:  modVersionRange,
			LoaderSupport: modLoadersType,
			Description:   mod.Description,
			Author:        mod.Author,
			Categories:    mod.Categories,
		}

		output.Results = append(output.Results, modInfo)
	}

	jsonData, err := json.MarshalIndent(output, "", "  ")
	if err != nil {
		logger.Error(fmt.Sprintf("JSON 序列化失败: %v", err))
		return
	}

	fmt.Println(string(jsonData))
}

// 表格输出函数
func outputSearchTable(result modrinth.IMrSearchResponse, offset int) {
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
}
