package storage

import (
	"archive/zip"
	"encoding/json"
	"os"
	"path/filepath"
	"time"

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

func ExportEntries(progressCallback func(percent int)) error {
	home, err := os.UserHomeDir()
	if err != nil {
		return err
	}

	entriesDir, err := EntriesDir()
	if err != nil {
		return err
	}

	outputPath := filepath.Join(home, "Downloads", "journaltui-export.zip")

	files, err := filepath.Glob(filepath.Join(entriesDir, "*.md"))
	if err != nil {
		return err
	}

	zipFile, err := os.Create(outputPath)
	if err != nil {
		return err
	}
	defer zipFile.Close()

	writer := zip.NewWriter(zipFile)
	defer writer.Close()

	total := len(files)
	if total == 0 {
		progressCallback(100)
		time.Sleep(1 * time.Second)
		return nil
	}

	for i, file := range files {
		data, err := os.ReadFile(file)
		if err != nil {
			return err
		}

		f, err := writer.Create(filepath.Base(file))
		if err != nil {
			return err
		}

		_, err = f.Write(data)
		if err != nil {
			return err
		}

		percent := int(float64(i+1) / float64(total) * 100)
		progressCallback(percent)
		time.Sleep(50 * time.Millisecond) // give tview time to render each frame
	}

	// Hold at 100% briefly so user sees the completed bar
	time.Sleep(500 * time.Millisecond)

	return nil
}
