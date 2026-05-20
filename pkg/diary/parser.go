package diary

import (
	"bytes"
	"strings"

	"gopkg.in/yaml.v3"
)

type Frontmatter struct {
	Title   string   `yaml:"title,omitempty"`
	Tags    []string `yaml:"tags,omitempty"`
	Mood    string   `yaml:"mood,omitempty"`
	Weather string   `yaml:"weather,omitempty"`
}

type DecryptedEntry struct {
	Entry       Entry
	Content     string
	Frontmatter Frontmatter
	WordCount   int
}

func ParseContent(raw string) DecryptedEntry {
	var fm Frontmatter
	var body string

	if strings.HasPrefix(raw, "---\n") || strings.HasPrefix(raw, "---\r\n") {
		parts := strings.SplitN(raw, "---", 3)
		if len(parts) >= 3 {

			yamlData := parts[1]
			body = strings.TrimSpace(parts[2])
			_ = yaml.Unmarshal([]byte(yamlData), &fm)
		} else {
			body = strings.TrimSpace(raw)
		}
	} else {
		body = strings.TrimSpace(raw)
	}

	words := len(bytes.Fields([]byte(body)))

	return DecryptedEntry{
		Content:     body,
		Frontmatter: fm,
		WordCount:   words,
	}
}
