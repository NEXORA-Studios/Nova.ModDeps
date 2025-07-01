package lock

type ILockFile struct {
	LockFileVersion int            `json:"lockfileVersion"`
	BaseDir         string         `json:"basedir"`
	Pending         []ILockModItem `json:"pending"`
	Installed       []ILockModItem `json:"installed"`
	NeedRemove      []ILockModItem `json:"needRemove"`
}

type ILockModItem struct {
	ID      string `json:"id"`
	Version string `json:"version"`
	Path    string `json:"path"`
	Uri     string `json:"uri"`
	Sha512  string `json:"sha512"`
}
