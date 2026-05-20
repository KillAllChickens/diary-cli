package cmd

import (
	"fmt"
	"os"

	"diary/pkg/diary"
	"github.com/manifoldco/promptui"
	"github.com/spf13/cobra"
)

func deleteExistingEntry(entry diary.Entry) {

	_, err := getAndVerifyPassword()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return
	}

	prompt := promptui.Prompt{
		Label:     fmt.Sprintf("Are you sure you want to delete '%s'", entry.Display()),
		IsConfirm: true,
	}

	result, err := prompt.Run()
	if err != nil || (result != "y" && result != "Y" && result != "") {
		fmt.Println("Deletion cancelled.")
		return
	}

	if err := os.Remove(entry.Path); err != nil {
		fmt.Fprintf(os.Stderr, "Failed to delete entry: %v\n", err)
		return
	}

	fmt.Println("Entry deleted successfully.")
}

var deleteCmd = &cobra.Command{
	Use:     "delete [date]",
	Aliases: []string{"rm", "remove"},
	Short:   "Delete a diary entry",
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
				deleteExistingEntry(filtered[0])
				return
			} else if len(filtered) > 1 {
				var items []string
				for _, e := range filtered {
					items = append(items, e.Display())
				}

				prompt := promptui.Select{
					Label: fmt.Sprintf("Multiple entries found for %s. Select one to delete:", targetDate),
					Items: items,
					Size:  10,
				}

				index, _, err := prompt.Run()
				if err == nil {
					deleteExistingEntry(filtered[index])
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
			Label: "Select a diary entry to delete",
			Items: items,
			Size:  10,
		}

		index, _, err := prompt.Run()
		if err != nil {
			return
		}

		deleteExistingEntry(entries[index])
	},
}

func init() {
	rootCmd.AddCommand(deleteCmd)
}
