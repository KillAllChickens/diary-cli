package cmd

import (
	"fmt"
	"os"

	"diary/pkg/diary"
	"github.com/manifoldco/promptui"
	"github.com/spf13/cobra"
)

var readCmd = &cobra.Command{
	Use:   "read [date]",
	Short: "Read a diary entry to stdout",
	Run: func(cmd *cobra.Command, args []string) {
		entries, err := diary.ListEntries()
		if err != nil {
			fmt.Fprintln(os.Stderr, "Failed to list entries:", err)
			return
		}

		if len(args) > 0 {
			targetDate := args[0]
			var filtered []diary.Entry
			for _, e := range entries {
				if e.Date == targetDate {
					filtered = append(filtered, e)
				}
			}

			if len(filtered) == 1 {
				readExistingEntry(filtered[0])
				return
			} else if len(filtered) > 1 {
				var items []string
				for _, e := range filtered {
					items = append(items, e.Display())
				}

				prompt := promptui.Select{
					Label: fmt.Sprintf("Multiple entries found for %s. Select one:", targetDate),
					Items: items,
					Size:  10,
				}

				index, _, err := prompt.Run()
				if err == nil {
					readExistingEntry(filtered[index])
				}
				return
			} else {
				fmt.Println("Entry not found for that date.")
				return
			}
		}

		if len(entries) == 0 {
			fmt.Println("No entries found.")
			return
		}

		var items []string
		for _, e := range entries {
			items = append(items, e.Display())
		}

		prompt := promptui.Select{
			Label: "Select a diary entry to read",
			Items: items,
			Size:  10,
		}

		index, _, err := prompt.Run()
		if err != nil {
			return
		}

		readExistingEntry(entries[index])
	},
}

func init() {
	rootCmd.AddCommand(readCmd)
}
