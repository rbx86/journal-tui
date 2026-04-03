package storage

import (
	"encoding/json"
	"os"
	"path/filepath"

	"journaltui/app"
)

func dataDir() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(home, ".journaltui"), nil
}

func EntriesDir() (string, error) {
	base, err := dataDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(base, "entries"), nil
}

func MetadataPath() (string, error) {
	base, err := dataDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(base, "metadata.json"), nil
}

func Init() error {
	entriesDir, err := EntriesDir()
	if err != nil {
		return err
	}
	return os.MkdirAll(entriesDir, 0755)
}

func LoadMetadata() (map[string]app.EntryMeta, error) {
	path, err := MetadataPath()
	if err != nil {
		return nil, err
	}

	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return make(map[string]app.EntryMeta), nil
		}
		return nil, err
	}

	var entries map[string]app.EntryMeta
	if err := json.Unmarshal(data, &entries); err != nil {
		return nil, err
	}
	return entries, nil
}

func SaveMetadata(entries map[string]app.EntryMeta) error {
	path, err := MetadataPath()
	if err != nil {
		return err
	}

	data, err := json.MarshalIndent(entries, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(path, data, 0644)
}

func SaveEntry(entryID, content string) error {
	dir, err := EntriesDir()
	if err != nil {
		return err
	}

	path := filepath.Join(dir, entryID+".md")
	return os.WriteFile(path, []byte(content), 0644)
}

func LoadEntry(entryID string) (string, error) {
	dir, err := EntriesDir()
	if err != nil {
		return "", err
	}

	path := filepath.Join(dir, entryID+".md")
	data, err := os.ReadFile(path)
	if err != nil {
		return "", err
	}
	return string(data), nil
}

func DeleteEntry(entryID string) error {
	dir, err := EntriesDir()
	if err != nil {
		return err
	}

	path := filepath.Join(dir, entryID+".md")
	return os.Remove(path)
}
