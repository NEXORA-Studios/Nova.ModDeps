package lock

func AddInstalled(id string, version string) {
	uLockFile, err := Read()
	if err != nil {
		return
	}
	// 从 Pending 中找到对应的版本，移动到 Installed 中
	for i, pending := range uLockFile.Pending {
		if pending.ID == id && pending.Version == version {
			uLockFile.Installed = append(uLockFile.Installed, uLockFile.Pending[i])
			uLockFile.Pending = append(uLockFile.Pending[:i], uLockFile.Pending[i+1:]...)
			break
		}
	}
	Write(uLockFile)
}
