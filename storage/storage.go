package storage

import (
	"context"
	"database/sql"

	"github.com/NEXORA-Studios/Nova.ModDeps/api_client"
	_ "github.com/mattn/go-sqlite3"
)

type Storage struct {
	db *sql.DB
}

func NewStorage(dbPath string) (*Storage, error) {
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		return nil, err
	}

	if err := createTables(db); err != nil {
		return nil, err
	}

	return &Storage{db: db}, nil
}

func createTables(db *sql.DB) error {
	_, err := db.Exec(`
        CREATE TABLE IF NOT EXISTS versions (
            id TEXT PRIMARY KEY,
            project_id TEXT NOT NULL,
            name TEXT NOT NULL,
            version_number TEXT NOT NULL,
            changelog TEXT,
            version_type TEXT NOT NULL,
            downloads INTEGER NOT NULL,
            date_published DATETIME NOT NULL
        );

        CREATE TABLE IF NOT EXISTS dependencies (
            version_id TEXT NOT NULL,
            dependency_id TEXT NOT NULL,
            dependency_type TEXT NOT NULL,
            FOREIGN KEY(version_id) REFERENCES versions(id)
        );

        CREATE TABLE IF NOT EXISTS game_versions (
            version_id TEXT NOT NULL,
            game_version TEXT NOT NULL,
            FOREIGN KEY(version_id) REFERENCES versions(id)
        );

        CREATE TABLE IF NOT EXISTS loaders (
            version_id TEXT NOT NULL,
            loader TEXT NOT NULL,
            FOREIGN KEY(version_id) REFERENCES versions(id)
        );
    `)
	return err
}

func (s *Storage) SaveVersion(ctx context.Context, version *api_client.Version) error {
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	_, err = tx.ExecContext(ctx,
		`INSERT OR REPLACE INTO versions 
         (id, project_id, name, version_number, changelog, version_type, downloads, date_published)
         VALUES (?, ?, ?, ?, ?, ?, ?, ?)`,
		version.ID, version.ProjectID, version.Name, version.VersionNumber,
		version.Changelog, version.VersionType, version.Downloads, version.DatePublished)
	if err != nil {
		return err
	}

	for _, dep := range version.Dependencies {
		_, err = tx.ExecContext(ctx,
			`INSERT OR REPLACE INTO dependencies 
             (version_id, dependency_id, dependency_type)
             VALUES (?, ?, ?)`,
			version.ID, dep.VersionID, dep.DependencyType)
		if err != nil {
			return err
		}
	}

	for _, gameVersion := range version.GameVersions {
		_, err = tx.ExecContext(ctx,
			`INSERT OR REPLACE INTO game_versions 
             (version_id, game_version)
             VALUES (?, ?)`,
			version.ID, gameVersion)
		if err != nil {
			return err
		}
	}

	for _, loader := range version.Loaders {
		_, err = tx.ExecContext(ctx,
			`INSERT OR REPLACE INTO loaders 
             (version_id, loader)
             VALUES (?, ?)`,
			version.ID, loader)
		if err != nil {
			return err
		}
	}

	return tx.Commit()
}

func (s *Storage) GetVersion(ctx context.Context, id string) (*api_client.Version, error) {
	var version api_client.Version
	err := s.db.QueryRowContext(ctx,
		`SELECT id, project_id, name, version_number, changelog, version_type, downloads, date_published
         FROM versions WHERE id = ?`, id).Scan(
		&version.ID, &version.ProjectID, &version.Name, &version.VersionNumber,
		&version.Changelog, &version.VersionType, &version.Downloads, &version.DatePublished)
	if err != nil {
		return nil, err
	}

	rows, err := s.db.QueryContext(ctx,
		`SELECT dependency_id, dependency_type FROM dependencies WHERE version_id = ?`, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var dep api_client.Dependency
		if err := rows.Scan(&dep.VersionID, &dep.DependencyType); err != nil {
			return nil, err
		}
		version.Dependencies = append(version.Dependencies, dep)
	}

	rows, err = s.db.QueryContext(ctx,
		`SELECT game_version FROM game_versions WHERE version_id = ?`, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var gameVersion string
		if err := rows.Scan(&gameVersion); err != nil {
			return nil, err
		}
		version.GameVersions = append(version.GameVersions, gameVersion)
	}

	rows, err = s.db.QueryContext(ctx,
		`SELECT loader FROM loaders WHERE version_id = ?`, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var loader string
		if err := rows.Scan(&loader); err != nil {
			return nil, err
		}
		version.Loaders = append(version.Loaders, loader)
	}

	return &version, nil
}
