package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
)

var searchCmd = &cobra.Command{
	Use:   "search <query>",
	Short: "Search your entire diary history",
	Long:  `Temporarily decrypts all your entries in memory and performs a full-text search for the given query.`,
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		query := strings.ToLower(strings.Join(args, " "))

		fmt.Println("Decrypting entries for search...")
		entries, err := decryptAllEntries()
		if err != nil {
			fmt.Fprintln(os.Stderr, "Failed to decrypt entries:", err)
			return
		}

		fmt.Printf("\nSearch results for '%s':\n", query)
		found := false

		for _, e := range entries {

			lowerBody := strings.ToLower(e.Content)
			if strings.Contains(lowerBody, query) || strings.Contains(strings.ToLower(e.Entry.Title), query) {
				found = true
				fmt.Printf("\n--- Match in: %s ---\n", e.Entry.Display())

				lines := strings.Split(e.Content, "\n")
				for _, line := range lines {
					if strings.Contains(strings.ToLower(line), query) {

						fmt.Printf("  > %s\n", strings.TrimSpace(line))
						break
					}
				}
			}
		}

		if !found {
			fmt.Println("No matches found.")
		}
		fmt.Println()
	},
}

func init() {
	rootCmd.AddCommand(searchCmd)
}
