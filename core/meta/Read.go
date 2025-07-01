package meta

import (
	"encoding/json"
	"errors"
	"io"
	"os"

	"github.com/NEXORA-Studios/Nova.ModDeps/core"
)

func (m *MetaFunctions) Read() (IModPackageJson, error) {
	path := core.GetModPackageJsonPath()
	jsonFile, err := os.Open(path)
	if err != nil {
		return IModPackageJson{}, err
	}
	defer jsonFile.Close()

	byteValue, err := io.ReadAll(jsonFile)
	if err != nil {
		return IModPackageJson{}, err
	}

	var modPackageJson IModPackageJson
	err = json.Unmarshal(byteValue, &modPackageJson)
	if err != nil {
		return IModPackageJson{}, err
	}

	return modPackageJson, nil
}

func (m *MetaFunctions) GetMetaVersion() (string, error) {
	modPackageJson, err := m.Read()
	if err != nil {
		return "", err
	}

	return modPackageJson.InternalVersion, nil
}

func (m *MetaFunctions) GetMetaPlatform() (string, error) {
	modPackageJson, err := m.Read()
	if err != nil {
		return "", err
	}

	return modPackageJson.InternalPlatform, nil
}

func (m *MetaFunctions) GetMinecraftVersion() (string, error) {
	modPackageJson, err := m.Read()
	if err != nil {
		return "", err
	}
	return modPackageJson.MinecraftVersion, nil
}

func (m *MetaFunctions) GetModLoader() ([]string, error) {
	modPackageJson, err := m.Read()
	if err != nil {
		return nil, err
	}
	return modPackageJson.ModLoader, nil
}

func (m *MetaFunctions) GetMods() ([]IModItem, error) {
	modPackageJson, err := m.Read()
	if err != nil {
		return nil, err
	}

	return modPackageJson.Mods, nil
}

func (m *MetaFunctions) GetModById(id string) (IModItem, error) {
	mods, err := m.GetMods()
	if err != nil {
		return IModItem{}, err
	}

	for _, mod := range mods {
		if mod.ID == id {
			return mod, nil
		}
	}

	return IModItem{}, errors.New("mod not found")
}

func (m *MetaFunctions) GetDependenciesById(id string) ([]IModDependencyItem, error) {
	mods, err := m.GetModById(id)
	if err != nil {
		return nil, err
	}

	return mods.Dependencies, nil
}

// IModItem.RequiredBy 字段用于记录哪些 mod 依赖此 mod，便于依赖安全移除
