package api_client

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type ModrinthClient struct {
	baseURL    string
	httpClient *http.Client
}

func NewClient() *ModrinthClient {
	return &ModrinthClient{
		baseURL: "https://api.modrinth.com/v2",
		httpClient: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

func (c *ModrinthClient) GetVersion(ctx context.Context, versionID string) (*Version, error) {
	url := fmt.Sprintf("%s/version/%s", c.baseURL, versionID)

	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, err
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	var version Version
	if err := json.NewDecoder(resp.Body).Decode(&version); err != nil {
		return nil, err
	}

	return &version, nil
}

type Version struct {
	ID            string       `json:"id"`
	ProjectID     string       `json:"project_id"`
	Name          string       `json:"name"`
	VersionNumber string       `json:"version_number"`
	Changelog     string       `json:"changelog"`
	Dependencies  []Dependency `json:"dependencies"`
	GameVersions  []string     `json:"game_versions"`
	VersionType   string       `json:"version_type"`
	Loaders       []string     `json:"loaders"`
	Downloads     int          `json:"downloads"`
	DatePublished time.Time    `json:"date_published"`
}

type Dependency struct {
	VersionID      string  `json:"version_id"`
	ProjectID      string  `json:"project_id"`
	FileName       *string `json:"file_name"`
	DependencyType string  `json:"dependency_type"`
}
