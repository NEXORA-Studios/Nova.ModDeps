package meta

import (
	"encoding/json"
	"os"

	"github.com/NEXORA-Studios/Nova.ModDeps/core"
)

func (m *MetaFunctions) Write(modPackageJson IModPackageJson) error {
	jsonData, err := json.Marshal(modPackageJson)
	if err != nil {
		return err
	}

	return os.WriteFile(core.GetModPackageJsonPath(), jsonData, 0644)
}
