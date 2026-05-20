package cmd

import (
	"fmt"
	"os"
	"time"

	"diary/pkg/diary"
	"github.com/manifoldco/promptui"
	"github.com/spf13/cobra"
)

var loadCmd = &cobra.Command{
	Use:   "load [date]",
	Short: "Load and edit a diary entry",
	Long:  `Load an entry for editing. Opens a TUI if no date is provided, or if multiple entries exist for the date. Date format: MM-DD-YYYY`,
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
				editExistingEntry(filtered[0])
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
					editExistingEntry(filtered[index])
				}
				return
			} else {

				createNewEntry(targetDate, "")
				return
			}
		}

		const createNewOpt = "+ Create new entry for today"
		var items []string
		items = append(items, createNewOpt)
		for _, e := range entries {
			items = append(items, e.Display())
		}

		prompt := promptui.Select{
			Label: "Select a diary entry",
			Items: items,
			Size:  10,
		}

		index, _, err := prompt.Run()
		if err != nil {
			return
		}

		if index == 0 {
			createNewEntry(time.Now().Format(diary.DateFormat), "")
		} else {
			editExistingEntry(entries[index-1])
		}
	},
}

func init() {
	rootCmd.AddCommand(loadCmd)
}
