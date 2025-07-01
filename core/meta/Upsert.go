package meta

func (m *MetaFunctions) UpsertMod(id string, version string, name string, requiredBy []string) error {
	modPackageJson, err := m.Read()
	if err != nil {
		return err
	}

	for i, mod := range modPackageJson.Mods {
		if mod.ID == id {
			modPackageJson.Mods[i].Version = version
			modPackageJson.Mods[i].Name = name
			mergedRequiredBy := mergeStringSlices(mod.RequiredBy, requiredBy)
			modPackageJson.Mods[i].RequiredBy = mergedRequiredBy
			return m.Write(modPackageJson)
		}
	}

	modPackageJson.Mods = append(modPackageJson.Mods, IModItem{
		ID:           id,
		Name:         name,
		Version:      version,
		Dependencies: []IModDependencyItem{},
		RequiredBy:   requiredBy,
	})

	return m.Write(modPackageJson)
}

func (m *MetaFunctions) UpsertDependency(modId string, dependencyId string, dependencyVersion string, requiredBy string) error {
	modPackageJson, err := m.Read()
	if err != nil {
		return err
	}

	mod, err := m.GetModById(modId)
	if err != nil {
		return err
	}

	for i, dependency := range mod.Dependencies {
		if dependency.ID == dependencyId {
			mod.Dependencies[i].Version = dependencyVersion
			// 更新 requiredBy
			for j, m := range modPackageJson.Mods {
				if m.ID == dependencyId {
					if !contains(m.RequiredBy, requiredBy) {
						modPackageJson.Mods[j].RequiredBy = append(modPackageJson.Mods[j].RequiredBy, requiredBy)
					}
				}
			}
			return m.Write(modPackageJson)
		}
	}

	mod.Dependencies = append(mod.Dependencies, IModDependencyItem{
		ID:      dependencyId,
		Version: dependencyVersion,
	})

	for j, m := range modPackageJson.Mods {
		if m.ID == dependencyId {
			if !contains(m.RequiredBy, requiredBy) {
				modPackageJson.Mods[j].RequiredBy = append(modPackageJson.Mods[j].RequiredBy, requiredBy)
			}
		}
	}

	return m.Write(modPackageJson)
}

func contains(arr []string, str string) bool {
	for _, a := range arr {
		if a == str {
			return true
		}
	}
	return false
}

func mergeStringSlices(a, b []string) []string {
	m := make(map[string]struct{})
	for _, v := range a {
		m[v] = struct{}{}
	}
	for _, v := range b {
		m[v] = struct{}{}
	}
	result := make([]string, 0, len(m))
	for k := range m {
		result = append(result, k)
	}
	return result
}
