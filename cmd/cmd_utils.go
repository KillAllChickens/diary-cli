package cmd

import (
	"os"

	"diary/pkg/crypto"
	"diary/pkg/diary"
	"github.com/charmbracelet/glamour"
)

func decryptAllEntries() ([]diary.DecryptedEntry, error) {
	entries, err := diary.ListEntries()
	if err != nil {
		return nil, err
	}
	if len(entries) == 0 {
		return nil, nil
	}

	password, err := getAndVerifyPassword()
	if err != nil {
		return nil, err
	}

	var decrypted []diary.DecryptedEntry
	for _, e := range entries {
		fileData, err := os.ReadFile(e.Path)
		if err != nil {
			continue
		}

		plaintext, err := crypto.Decrypt(fileData, password)
		if err != nil {
			continue
		}

		parsed := diary.ParseContent(string(plaintext))
		parsed.Entry = e
		decrypted = append(decrypted, parsed)
	}

	return decrypted, nil
}

func renderMarkdown(content string) string {
	r, err := glamour.NewTermRenderer(
		glamour.WithAutoStyle(),
		glamour.WithWordWrap(80),
	)
	if err != nil {
		return content
	}

	out, err := r.Render(content)
	if err != nil {
		return content
	}
	return out
}
