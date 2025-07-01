package modrinth

import (
	"encoding/json"
	"fmt"

	"github.com/NEXORA-Studios/Nova.ModDeps/core/api"
)

// GetProjectVersion 根据项目 id 获取所有版本信息
func GetProjectVersion(id string) ([]IMrModVersion, error) {
	requester := api.Requester{}
	path := fmt.Sprintf("/project/%s/version", id)
	response, err := requester.Get(path, map[string]string{})
	if err != nil {
		return nil, err
	}
	var versions []IMrModVersion
	err = json.Unmarshal([]byte(response), &versions)
	if err != nil {
		return nil, err
	}
	return versions, nil
}
