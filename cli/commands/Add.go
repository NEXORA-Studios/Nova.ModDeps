package commands

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/NEXORA-Studios/Nova.ModDeps/core/meta"
	"github.com/spf13/cobra"
)

type ModrinthVersion struct {
	ID            string `json:"id"`
	ProjectID     string `json:"project_id"`
	VersionNumber string `json:"version_number"`
	Name          string `json:"name"`
	Dependencies  []struct {
		VersionID      string `json:"version_id"`
		ProjectID      string `json:"project_id"`
		DependencyType string `json:"dependency_type"`
	} `json:"dependencies"`
}

func fetchModrinthVersion(versionID string) (*ModrinthVersion, error) {
	url := "https://api.modrinth.com/v2/version/" + versionID
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("Modrinth API 返回错误: %d", resp.StatusCode)
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	var v ModrinthVersion
	if err := json.Unmarshal(body, &v); err != nil {
		return nil, err
	}
	return &v, nil
}

// 递归添加 Mod 及其依赖
func AddModByVersion(projectID, versionID, requiredBy string) error {
	metaFunc := meta.MetaFunctions{}
	v, err := fetchModrinthVersion(versionID)
	if err != nil {
		return err
	}
	// 检查 requiredBy 是否能被满足，避免有作者没有在 Modrinth 上指定依赖的版本
	for _, dep := range v.Dependencies {
		if dep.DependencyType != "required" || dep.ProjectID == "" || dep.VersionID == "" {
			if _, err := metaFunc.GetModById(dep.ProjectID); err != nil {
				return fmt.Errorf("获取依赖版本失败：项目 %s 的作者没有在 Modrinth 上指定依赖的版本，且目前列表中没有已经配置版本的依赖\n          请先安装项目 %s 的对应版本，然后再试一次", v.ProjectID, dep.ProjectID)
			}
		}
	}
	// 添加本体
	err = metaFunc.UpsertMod(projectID, versionID, v.Name, requiredByList(requiredBy))
	if err != nil {
		return err
	}
	// 递归添加依赖
	for _, dep := range v.Dependencies {
		if dep.DependencyType != "required" || dep.ProjectID == "" {
			continue
		}
		// 递归添加依赖（仅当有具体版本）
		if dep.VersionID != "" {
			err = AddModByVersion(dep.ProjectID, dep.VersionID, versionID)
			if err != nil {
				return err
			}
		} else {
			mod, err := metaFunc.GetModById(dep.ProjectID)
			if err != nil {
				return err
			}
			version := mod.Version
			err = AddModByVersion(dep.ProjectID, version, versionID)
			if err != nil {
				return err
			}
		}
		// 依赖关系写入（无论 version_id 是否为空都写）
		err = metaFunc.UpsertDependency(projectID, dep.ProjectID, dep.VersionID, versionID)
		if err != nil {
			return err
		}
	}
	return nil
}

func requiredByList(requiredBy string) []string {
	if requiredBy == "" {
		return []string{}
	}
	return []string{requiredBy}
}

var AddCmd = &cobra.Command{
	Use:   "add <project_id> <version_id>",
	Short: "添加指定 Mod 及其所有依赖（递归）",
	Args:  cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		projectID := args[0]
		versionID := args[1]
		err := AddModByVersion(projectID, versionID, "")
		if err != nil {
			fmt.Fprintf(os.Stderr, "添加失败: %v\n", err)
			os.Exit(1)
		}
		fmt.Printf("已成功添加 Mod %s 及其依赖\n", projectID)
	},
}
