package modrinth

import (
	"encoding/json"
	"fmt"

	"github.com/NEXORA-Studios/Nova.ModDeps/core/api"
)

func GetProjectMetadata(id string) (IMrModProject, error) {
	requester := api.Requester{}
	path := fmt.Sprintf("/project/%s", id)
	response, err := requester.Get(path, map[string]string{})
	if err != nil {
		return IMrModProject{}, err
	}
	var project IMrModProject
	err = json.Unmarshal([]byte(response), &project)
	if err != nil {
		return IMrModProject{}, err
	}
	return project, nil
}

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

func GetVersionMetadata(id string) (IMrModVersion, error) {
	requester := api.Requester{}
	path := fmt.Sprintf("/version/%s", id)
	response, err := requester.Get(path, map[string]string{})
	if err != nil {
		return IMrModVersion{}, err
	}
	var version IMrModVersion
	err = json.Unmarshal([]byte(response), &version)
	if err != nil {
		return IMrModVersion{}, err
	}
	return version, nil
}
