package cmd

import (
	"fmt"
	"os"

	"diary/pkg/config"
	"diary/pkg/crypto"
	"github.com/spf13/cobra"
)

var setupCmd = &cobra.Command{
	Use:   "setup",
	Short: "Initialize the diary and set the global password",
	Run: func(cmd *cobra.Command, args []string) {
		if _, err := os.Stat(config.ConfigFile); err == nil {
			fmt.Println("Diary is already setup! If you want to reset, delete the .config.json in your diary folder.")
			return
		}

		p1, err := getPassword("Enter new global password: ")
		if err != nil {
			fmt.Fprintln(os.Stderr, "Error reading password:", err)
			return
		}

		p2, err := getPassword("Confirm password: ")
		if err != nil {
			fmt.Fprintln(os.Stderr, "Error reading password:", err)
			return
		}

		if string(p1) != string(p2) {
			fmt.Fprintln(os.Stderr, "Passwords do not match. Aborting setup.")
			return
		}

		hash, err := crypto.HashPassword(p1)
		if err != nil {
			fmt.Fprintln(os.Stderr, "Error hashing password:", err)
			return
		}

		cfg := &config.Config{PasswordHash: hash}
		if err := config.SaveConfig(cfg); err != nil {
			fmt.Fprintln(os.Stderr, "Error saving config:", err)
			return
		}

		fmt.Println("Diary successfully initialized! You can now use 'diary load' to write your first entry.")
	},
}

func init() {
	rootCmd.AddCommand(setupCmd)
}
