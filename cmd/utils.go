package cmd

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"diary/pkg/config"
	"diary/pkg/crypto"
	"diary/pkg/diary"

	"golang.org/x/term"
)

func getPassword(prompt string) ([]byte, error) {
	fmt.Print(prompt)
	password, err := term.ReadPassword(int(os.Stdin.Fd()))
	fmt.Println()
	if err != nil {
		return nil, err
	}
	return password, nil
}

type SessionData struct {
	ExpiresAt int64  `json:"expires_at"`
	Password  string `json:"password"`
}

func getSessionFilePath() string {
	return filepath.Join(os.TempDir(), fmt.Sprintf("virt-diary-session-%d.json", os.Getuid()))
}

func getCachedPassword() (string, bool) {
	data, err := os.ReadFile(getSessionFilePath())
	if err != nil {
		return "", false
	}
	var session SessionData
	if err := json.Unmarshal(data, &session); err != nil {
		return "", false
	}
	if time.Now().Unix() > session.ExpiresAt {
		os.Remove(getSessionFilePath())
		return "", false
	}
	return session.Password, true
}

func cachePassword(password string) {
	session := SessionData{
		ExpiresAt: time.Now().Add(15 * time.Minute).Unix(),
		Password:  password,
	}
	data, err := json.Marshal(session)
	if err == nil {
		os.WriteFile(getSessionFilePath(), data, 0600)
	}
}

func ClearPasswordCache() {
	os.Remove(getSessionFilePath())
}

func getAndVerifyPassword() ([]byte, error) {
	cfg, err := config.LoadConfig()
	if err != nil {
		if os.IsNotExist(err) {
			return nil, errors.New("diary not setup. Please run 'diary setup' first")
		}
		return nil, fmt.Errorf("error loading config: %w", err)
	}

	if cachedPwd, ok := getCachedPassword(); ok {
		if crypto.VerifyPassword([]byte(cachedPwd), cfg.PasswordHash) {

			cachePassword(cachedPwd)
			return []byte(cachedPwd), nil
		}
	}

	password, err := getPassword("Enter global diary password: ")
	if err != nil {
		return nil, err
	}

	if !crypto.VerifyPassword(password, cfg.PasswordHash) {
		return nil, errors.New("incorrect password")
	}

	cachePassword(string(password))
	return password, nil
}

func generateUniqueFilename(targetDate, title string) string {
	baseName := targetDate
	safeTitle := strings.ReplaceAll(title, " ", "_")
	if safeTitle != "" {
		baseName = fmt.Sprintf("%s_%s", targetDate, safeTitle)
	}

	filename := filepath.Join(config.DiaryDir, baseName+".md")
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		return filename
	}

	for i := 2; ; i++ {
		filename = filepath.Join(config.DiaryDir, fmt.Sprintf("%s_%d.md", baseName, i))
		if _, err := os.Stat(filename); os.IsNotExist(err) {
			return filename
		}
	}
}

func createNewEntry(targetDate, presetTitle string) {
	password, err := getAndVerifyPassword()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return
	}

	template := fmt.Sprintf(`---
title: "%s"
tags: []
mood: ""
weather: ""
---

`, presetTitle)

	fmt.Println("New entry. Proceeding to editor...")

	editedData, err := diary.EditTempFile([]byte(template))
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return
	}

	parsed := diary.ParseContent(string(editedData))
	finalTitle := strings.TrimSpace(parsed.Frontmatter.Title)
	if finalTitle == "" {
		finalTitle = strings.TrimSpace(presetTitle)
	}

	filename := generateUniqueFilename(targetDate, finalTitle)

	ciphertext, err := crypto.Encrypt(editedData, password)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Encryption failed: %v\n", err)
		return
	}

	if err := os.WriteFile(filename, ciphertext, 0600); err != nil {
		fmt.Fprintf(os.Stderr, "Failed to save entry: %v\n", err)
		return
	}

	fmt.Printf("Entry for %s saved securely.\n", targetDate)
}

func editExistingEntry(entry diary.Entry) {
	password, err := getAndVerifyPassword()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return
	}

	fileData, err := os.ReadFile(entry.Path)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error reading file: %v\n", err)
		return
	}

	plaintext, err := crypto.Decrypt(fileData, password)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to decrypt entry: %v\n", err)
		return
	}

	editedData, err := diary.EditTempFile(plaintext)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return
	}

	ciphertext, err := crypto.Encrypt(editedData, password)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Encryption failed: %v\n", err)
		return
	}

	if err := os.WriteFile(entry.Path, ciphertext, 0600); err != nil {
		fmt.Fprintf(os.Stderr, "Failed to save entry: %v\n", err)
		return
	}

	fmt.Printf("Entry %s saved securely.\n", entry.Display())
}

func readExistingEntry(entry diary.Entry) {
	fileData, err := os.ReadFile(entry.Path)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error reading file: %v\n", err)
		return
	}

	password, err := getAndVerifyPassword()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return
	}

	plaintext, err := crypto.Decrypt(fileData, password)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to decrypt entry: %v\n", err)
		return
	}

	parsed := diary.ParseContent(string(plaintext))

	fmt.Printf("\n--- Entry: %s ---\n", entry.Display())
	if len(parsed.Frontmatter.Tags) > 0 || parsed.Frontmatter.Mood != "" || parsed.Frontmatter.Weather != "" {
		fmt.Print("Metadata: ")
		if parsed.Frontmatter.Mood != "" {
			fmt.Printf("Mood: %s | ", parsed.Frontmatter.Mood)
		}
		if parsed.Frontmatter.Weather != "" {
			fmt.Printf("Weather: %s | ", parsed.Frontmatter.Weather)
		}
		if len(parsed.Frontmatter.Tags) > 0 {
			fmt.Printf("Tags: %v", parsed.Frontmatter.Tags)
		}
		fmt.Print("\n\n")
	}

	fmt.Println(renderMarkdown(parsed.Content))
	fmt.Println("-------------------")
}
