package lock

func GetStatus(id string, version string) string {
	lockFile, err := Read()
	if err != nil {
		return "unknown"
	}
	for _, item := range lockFile.Pending {
		if item.ID == id && item.Version == version {
			return "pending"
		}
	}
	for _, item := range lockFile.Installed {
		if item.ID == id && item.Version == version {
			return "installed"
		}
	}
	for _, item := range lockFile.NeedRemove {
		if item.ID == id && item.Version == version {
			return "needRemove"
		}
	}
	return "unknown"
}
