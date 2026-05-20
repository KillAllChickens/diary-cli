package diary

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"diary/pkg/config"
)

const DateFormat = "01-02-2006"

type Entry struct {
	Date  string
	Title string
	Path  string
}

func (e Entry) Display() string {
	if e.Title != "" {
		return fmt.Sprintf("%s (%s)", e.Title, e.Date)
	}
	return e.Date
}

func ListEntries() ([]Entry, error) {
	files, err := os.ReadDir(config.DiaryDir)
	if err != nil {
		return nil, err
	}

	var entries []Entry
	for _, f := range files {
		if !f.IsDir() && strings.HasSuffix(f.Name(), ".md") {
			name := strings.TrimSuffix(f.Name(), ".md")
			parts := strings.SplitN(name, "_", 2)
			date := parts[0]
			title := ""
			if len(parts) > 1 {
				title = strings.ReplaceAll(parts[1], "_", " ")
			}
			entries = append(entries, Entry{
				Date:  date,
				Title: title,
				Path:  filepath.Join(config.DiaryDir, f.Name()),
			})
		}
	}

	sort.Slice(entries, func(i, j int) bool {
		ti, _ := time.Parse(DateFormat, entries[i].Date)
		tj, _ := time.Parse(DateFormat, entries[j].Date)
		return ti.After(tj)
	})

	return entries, nil
}

func EditTempFile(initialData []byte) ([]byte, error) {
	tempFile, err := os.CreateTemp("", "virt-diary-*.md")
	if err != nil {
		return nil, fmt.Errorf("failed to create temp file: %w", err)
	}
	tempName := tempFile.Name()
	defer os.Remove(tempName)

	if len(initialData) > 0 {
		if _, err = tempFile.Write(initialData); err != nil {
			tempFile.Close()
			return nil, fmt.Errorf("failed to write to temp file: %w", err)
		}
	}
	tempFile.Close()

	editor := os.Getenv("EDITOR")
	if editor == "" {
		editor = "nvim"
	}

	execCmd := exec.Command(editor, tempName)
	execCmd.Stdin = os.Stdin
	execCmd.Stdout = os.Stdout
	execCmd.Stderr = os.Stderr

	if err = execCmd.Run(); err != nil {
		return nil, fmt.Errorf("editor returned an error: %w", err)
	}

	editedData, err := os.ReadFile(tempName)
	if err != nil {
		return nil, fmt.Errorf("failed to read edited temp file: %w", err)
	}

	wipeData := make([]byte, len(editedData))
	os.WriteFile(tempName, wipeData, 0600)

	return editedData, nil
}
