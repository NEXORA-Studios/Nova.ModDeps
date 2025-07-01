package lock

import (
	"fmt"
	"path/filepath"

	"github.com/NEXORA-Studios/Nova.ModDeps/core/api/modrinth"
)

func UpsertPending(id string, version string) {
	uLockFile, err := Read()
	if err != nil {
		return
	}
	versionMetadata, err := modrinth.GetVersionMetadata(version)
	if err != nil {
		logger.Error(fmt.Sprintf("获取版本元数据失败: %v", err))
		return
	}
	var iPath string
	var iUri string
	var iSha512 string
	uris := versionMetadata.Files
	for _, uri := range uris {
		if uri.Primary {
			iUri = uri.URL
			iSha512 = uri.Hashes.SHA512
			iPath = filepath.Join(uLockFile.BaseDir, uri.Filename)
			break
		}
	}

	// 检查是否已存在相同 id 和 version，存在则更新，否则添加
	found := false
	for idx, item := range uLockFile.Pending {
		if item.ID == id && item.Version == version {
			uLockFile.Pending[idx] = ILockModItem{
				ID:      id,
				Version: version,
				Path:    iPath,
				Uri:     iUri,
				Sha512:  iSha512,
			}
			found = true
			break
		}
	}
	if !found {
		uLockFile.Pending = append(uLockFile.Pending, ILockModItem{
			ID:      id,
			Version: version,
			Path:    iPath,
			Uri:     iUri,
			Sha512:  iSha512,
		})
	}
	Write(uLockFile)
}

func RePending(id string, version string) {
	// 由于本地文件丢失而从 Installed 中找到对应的版本，移动到 Pending 中
	uLockFile, err := Read()
	if err != nil {
		return
	}
	for i, item := range uLockFile.Installed {
		if item.ID == id && item.Version == version {
			uLockFile.Pending = append(uLockFile.Pending, item)
			uLockFile.Installed = append(uLockFile.Installed[:i], uLockFile.Installed[i+1:]...)
		}
	}
}
