package lock

func AddNeedRemove(id string, version string) {
	uLockFile, err := Read()
	if err != nil {
		return
	}
	// 从 Installed 中找到对应的版本，移动到 NeedRemove 中
	for i, installed := range uLockFile.Installed {
		if installed.ID == id && installed.Version == version {
			uLockFile.NeedRemove = append(uLockFile.NeedRemove, uLockFile.Installed[i])
			uLockFile.Installed = append(uLockFile.Installed[:i], uLockFile.Installed[i+1:]...)
			break
		}
	}
	Write(uLockFile)
}

func Remove(id string, version string) {
	uLockFile, err := Read()
	if err != nil {
		return
	}
	// 从 NeedRemove 中找到对应的版本，删除
	for i, needRemove := range uLockFile.NeedRemove {
		if needRemove.ID == id && needRemove.Version == version {
			uLockFile.NeedRemove = append(uLockFile.NeedRemove[:i], uLockFile.NeedRemove[i+1:]...)
			break
		}
	}
	Write(uLockFile)
}
