package instance

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
)

func LoadInstance() (string, []string, error) {
	// 版本信息集
	var minecraftVersion string = ""
	var modLoader []string = []string{}

	cwd, err := os.Getwd()
	if err != nil {
		return "", []string{}, err
	}

	// 检查目前目录下是否有 json 文件
	jsonFiles, err := filepath.Glob(filepath.Join(cwd, "*.json"))
	if err != nil {
		return "", []string{}, err
	}
	// 遍历 json 文件，尝试获取 Minecraft 版本
GetMetaLoop:
	for _, jsonFile := range jsonFiles {
		jsonContent, err := os.ReadFile(jsonFile)
		if err != nil {
			return "", []string{}, err
		}
		var versionInfo map[string]any
		json.Unmarshal(jsonContent, &versionInfo)
		// 检查是否有 clientVersion 字段，有则作为 Minecraft 版本，没有则跳过 for 循环
		if _, ok := versionInfo["clientVersion"]; ok {
			minecraftVersion = versionInfo["clientVersion"].(string)
			// 全量匹配 forge neoforge fabric quilt 的 id 字符串，有则作为 ModLoader 类型，可能有多个
			if strings.Contains(string(jsonContent), "net.minecraftforge") {
				modLoader = append(modLoader, "forge")
			}
			if strings.Contains(string(jsonContent), "net.neoforged") {
				modLoader = append(modLoader, "neoforge")
				parts := strings.Split(minecraftVersion, ".")
				if len(parts) > 1 && parts[1] < "21" {
					modLoader = append(modLoader, "forge") // 1.20 以下 Forge == NeoForge，至于快照版先不管了 :/
				}
			}
			if strings.Contains(string(jsonContent), "net.fabricmc:fabric-loader") {
				modLoader = append(modLoader, "fabric")
			}
			if strings.Contains(string(jsonContent), "org.quiltmc") {
				modLoader = append(modLoader, "quilt")
				modLoader = append(modLoader, "fabric") // Quilt == Fabric 嗯，兼容性很好 :)
			}
			break GetMetaLoop
		} else {
			continue
		}
	}

	return minecraftVersion, modLoader, nil
}
