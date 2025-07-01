package modrinth

type ModrinthApiRequest struct{}

type IMrSearchResponse struct {
	Hits      []IMrModProject `json:"hits"`
	Offset    int             `json:"offset"`
	Limit     int             `json:"limit"`
	TotalHits int             `json:"total_hits"`
}

type IMrModProject struct {
	ProjectID         string   `json:"project_id"`
	ProjectType       string   `json:"project_type"`
	Slug              string   `json:"slug"`
	Author            string   `json:"author"`
	Title             string   `json:"title"`
	Description       string   `json:"description"`
	Categories        []string `json:"categories"`
	DisplayCategories []string `json:"display_categories"`
	Versions          []string `json:"versions"`
	Downloads         int64    `json:"downloads"`
	Follows           int      `json:"follows"`
	IconURL           string   `json:"icon_url"`
	DateCreated       string   `json:"date_created"`
	DateModified      string   `json:"date_modified"`
	LatestVersion     string   `json:"latest_version"`
	License           string   `json:"license"`
	ClientSide        string   `json:"client_side"`
	ServerSide        string   `json:"server_side"`
	Gallery           []string `json:"gallery"`
	FeaturedGallery   *string  `json:"featured_gallery"`
	Color             int      `json:"color"`
}

type IMrModVersion struct {
	GameVersions    []string        `json:"game_versions"`
	Loaders         []string        `json:"loaders"`
	ID              string          `json:"id"`
	ProjectID       string          `json:"project_id"`
	AuthorID        string          `json:"author_id"`
	Featured        bool            `json:"featured"`
	Name            string          `json:"name"`
	VersionNumber   string          `json:"version_number"`
	Changelog       string          `json:"changelog"`
	ChangelogURL    *string         `json:"changelog_url"`
	DatePublished   string          `json:"date_published"`
	Downloads       int             `json:"downloads"`
	VersionType     string          `json:"version_type"`
	Status          string          `json:"status"`
	RequestedStatus *string         `json:"requested_status"`
	Files           []IMrModFile    `json:"files"`
	Dependencies    []IMrDependency `json:"dependencies"`
}

type IMrModFile struct {
	Hashes   IMrHashes `json:"hashes"`
	URL      string    `json:"url"`
	Filename string    `json:"filename"`
	Primary  bool      `json:"primary"`
	Size     int       `json:"size"`
	FileType *string   `json:"file_type"`
}

type IMrHashes struct {
	SHA512 string `json:"sha512"`
	SHA1   string `json:"sha1"`
}

type IMrDependency struct {
	VersionID      string  `json:"version_id"`
	ProjectID      string  `json:"project_id"`
	FileName       *string `json:"file_name"`
	DependencyType string  `json:"dependency_type"`
}
