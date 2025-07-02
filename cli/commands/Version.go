package commands

import (
	"encoding/json"
	"fmt"
	"os"
	"slices"
	"strconv"

	"github.com/NEXORA-Studios/Nova.ModDeps/cli/utils"
	"github.com/NEXORA-Studios/Nova.ModDeps/core/api/modrinth"
	"github.com/NEXORA-Studios/Nova.ModDeps/core/meta"
	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/spf13/cobra"
)

var versionJson bool

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

// PrintProjectVersionsWithPagination 打印指定项目的所有版本，支持本地分页
func PrintProjectVersionsWithPagination(projectID string, page int, size int) {
	versions, err := modrinth.GetProjectVersion(projectID)
	if err != nil {
		logger.Fatal(fmt.Sprintf("获取版本失败: %v\n", err))
	}
	var localMeta *meta.IModPackageJson
	metaFunc := meta.MetaFunctions{}
	m, err := metaFunc.Read()
	if err != nil {
		logger.Fatal(fmt.Sprintf("读取 mod.package.json 失败: %v\n请确认是否在 Minecraft 实例根目录下运行，和 mod.package.json 文件是否存在\n若不存在，请先使用 \"novadm init\" 初始化", err))
	}
	localMeta = &m
	if size <= 0 {
		size = 10 // 默认每页10条
	}
	total := len(versions)
	start := (page - 1) * size
	if start >= total {
		logger.Warn("没有更多数据了。")
		return
	}
	end := min(start+size, total)

	if versionJson {
		// JSON 输出
		outputVersionJSON(projectID, versions[start:end], total, page, (total+size-1)/size, localMeta)
	} else {
		// 表格输出
		outputVersionTable(projectID, versions[start:end], total, page, (total+size-1)/size, localMeta)
	}
}

var VersionCmd = &cobra.Command{
	Use:   "version <project_id> [page]",
	Short: "查看指定项目的所有版本",
	Args:  cobra.RangeArgs(1, 2),
	Run: func(cmd *cobra.Command, args []string) {
		projectID := args[0]
		page := 1
		if len(args) >= 2 {
			p, err := strconv.Atoi(args[1])
			if err == nil && p > 0 {
				page = p
			}
		}
		PrintProjectVersionsWithPagination(projectID, page, 10)
	},
}

func init() {
	VersionCmd.Flags().BoolVarP(&versionJson, "json", "j", false, "以 JSON 格式输出")
}

// JSON 输出函数
func outputVersionJSON(projectID string, versions []modrinth.IMrModVersion, total, page, totalPages int, localMeta *meta.IModPackageJson) {
	type VersionInfo struct {
		ID            string   `json:"id"`
		Name          string   `json:"name"`
		DatePublished string   `json:"datePublished"`
		GameVersions  []string `json:"gameVersions"`
		Loaders       []string `json:"loaders"`
		Compatibility struct {
			GameMatch   bool `json:"gameMatch"`
			LoaderMatch bool `json:"loaderMatch"`
		} `json:"compatibility"`
	}

	type VersionOutput struct {
		ProjectID   string        `json:"projectId"`
		Total       int           `json:"total"`
		CurrentPage int           `json:"currentPage"`
		TotalPages  int           `json:"totalPages"`
		Versions    []VersionInfo `json:"versions"`
	}

	output := VersionOutput{
		ProjectID:   projectID,
		Total:       total,
		CurrentPage: page,
		TotalPages:  totalPages,
		Versions:    []VersionInfo{},
	}

	for _, v := range versions {
		matchGame := slices.Contains(v.GameVersions, localMeta.MinecraftVersion)
		matchLoader := false
		for _, loader := range v.Loaders {
			if slices.Contains(localMeta.ModLoader, loader) {
				matchLoader = true
				break
			}
		}

		versionInfo := VersionInfo{
			ID:            v.ID,
			Name:          v.Name,
			DatePublished: v.DatePublished,
			GameVersions:  v.GameVersions,
			Loaders:       v.Loaders,
			Compatibility: struct {
				GameMatch   bool `json:"gameMatch"`
				LoaderMatch bool `json:"loaderMatch"`
			}{
				GameMatch:   matchGame,
				LoaderMatch: matchLoader,
			},
		}

		output.Versions = append(output.Versions, versionInfo)
	}

	jsonData, err := json.MarshalIndent(output, "", "  ")
	if err != nil {
		logger.Error(fmt.Sprintf("JSON 序列化失败: %v", err))
		return
	}

	fmt.Println(string(jsonData))
}

// 表格输出函数
func outputVersionTable(projectID string, versions []modrinth.IMrModVersion, total, page, totalPages int, localMeta *meta.IModPackageJson) {
	fmt.Printf("项目 %s 共找到 %d 个版本，当前为 %d/%d 页\n", projectID, total, page, totalPages)

	t := table.NewWriter()
	t.SetOutputMirror(os.Stdout)
	t.AppendHeader(table.Row{"版本 ID", "名称", "发布时间", "支持游戏版本", "加载器"})

	for _, v := range versions {
		gameVersions := ""
		if len(v.GameVersions) > 0 {
			gameVersions = v.GameVersions[0]
			if len(v.GameVersions) > 1 {
				gameVersions += " ~ " + v.GameVersions[len(v.GameVersions)-1]
			}
		}
		loaders := ""
		if len(v.Loaders) > 0 {
			loaders = v.Loaders[0]
			if len(v.Loaders) > 1 {
				loaders += " / " + v.Loaders[len(v.Loaders)-1]
			}
		}

		matchGame := slices.Contains(v.GameVersions, localMeta.MinecraftVersion)
		matchLoader := false
		for _, loader := range v.Loaders {
			if slices.Contains(localMeta.ModLoader, loader) {
				matchLoader = true
				break
			}
		}

		switch {
		case matchGame && matchLoader:
			gameVersions = utils.ColorGreen + gameVersions + utils.ColorReset
			loaders = utils.ColorGreen + loaders + utils.ColorReset

		case matchGame && !matchLoader:
			gameVersions = utils.ColorYellow + gameVersions + utils.ColorReset
			loaders = utils.ColorRed + loaders + utils.ColorReset

		case !matchGame && matchLoader:
			gameVersions = utils.ColorRed + gameVersions + utils.ColorReset
			loaders = utils.ColorYellow + loaders + utils.ColorReset

		case !matchGame && !matchLoader:
			gameVersions = utils.ColorRed + gameVersions + utils.ColorReset
			loaders = utils.ColorRed + loaders + utils.ColorReset
		}

		row := table.Row{
			v.ID,
			v.Name,
			v.DatePublished,
			gameVersions,
			loaders,
		}

		t.AppendRow(row)
	}
	t.SetStyle(table.StyleLight)
	t.Render()
}
