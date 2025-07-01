package meta

import (
	"fmt"
)

func (m *MetaFunctions) RemoveMod(id string) error {
	modPackageJson, err := m.Read()
	if err != nil {
		return err
	}

	found := false
	mods := modPackageJson.Mods[:0]

	for _, mod := range modPackageJson.Mods {
		if mod.ID == id {
			found = true
			continue
		}

		newRequiredBy := []string{}
		for _, r := range mod.RequiredBy {
			if r != id {
				newRequiredBy = append(newRequiredBy, r)
			}
		}
		mod.RequiredBy = newRequiredBy
		mods = append(mods, mod)
	}

	if !found {
		return fmt.Errorf("模组 %s 不存在", id)
	}

	modPackageJson.Mods = mods
	return m.Write(modPackageJson)
}
