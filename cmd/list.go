package cmd

import (
	"fmt"
	"os"

	"diary/pkg/diary"
	"github.com/spf13/cobra"
)

var tagFilter string

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List all diary entries",
	Run: func(cmd *cobra.Command, args []string) {
		if tagFilter != "" {
			fmt.Printf("Decrypting entries to filter by tag '%s'...\n", tagFilter)
			entries, err := decryptAllEntries()
			if err != nil {
				fmt.Fprintf(os.Stderr, "Failed to decrypt entries: %v\n", err)
				return
			}

			fmt.Printf("Diary Entries (Tag: %s):\n", tagFilter)
			found := false
			for _, e := range entries {
				for _, t := range e.Frontmatter.Tags {
					if t == tagFilter {
						fmt.Printf("  - %s\n", e.Entry.Display())
						found = true
						break
					}
				}
			}
			if !found {
				fmt.Println("  No entries found.")
			}
			return
		}

		entries, err := diary.ListEntries()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Failed to list entries: %v\n", err)
			return
		}

		fmt.Println("Diary Entries:")
		if len(entries) == 0 {
			fmt.Println("  No entries found.")
			return
		}

		for _, e := range entries {
			fmt.Printf("  - %s\n", e.Display())
		}
	},
}

func init() {
	listCmd.Flags().StringVarP(&tagFilter, "tag", "t", "", "Filter entries by a specific tag")
	rootCmd.AddCommand(listCmd)
}
