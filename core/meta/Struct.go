package meta

type MetaFunctions struct{}

type IModPackageJson struct {
	InternalVersion  int        `json:"__version__"`
	InternalPlatform string     `json:"__platform__"`
	MinecraftVersion string     `json:"minecraftVersion"`
	ModLoader        []string   `json:"modLoader"`
	Mods             []IModItem `json:"mods"`
}

type IModItem struct {
	ID           string               `json:"id"`
	Name         string               `json:"name"`
	Version      string               `json:"version"`
	Dependencies []IModDependencyItem `json:"dependencies,omitempty"`
	RequiredBy   []string             `json:"requiredBy"`
}

type IModDependencyItem struct {
	ID      string `json:"id"`
	Version string `json:"version"`
}
