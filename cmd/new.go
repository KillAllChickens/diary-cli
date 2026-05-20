package cmd

import (
	"strings"
	"time"

	"diary/pkg/diary"
	"github.com/spf13/cobra"
)

var newCmd = &cobra.Command{
	Use:   "new [date] [title...]",
	Short: "Quickly create a new diary entry",
	Long:  `Instantly creates a brand new diary entry for today (or the specified date). Optional title can be provided. If no title is provided, it will prompt for one.`,
	Run: func(cmd *cobra.Command, args []string) {
		targetDate := time.Now().Format(diary.DateFormat)
		title := ""

		if len(args) > 0 {

			if _, err := time.Parse(diary.DateFormat, args[0]); err == nil {
				targetDate = args[0]
				if len(args) > 1 {
					title = strings.Join(args[1:], " ")
				}
			} else {

				title = strings.Join(args, " ")
			}
		}

		createNewEntry(targetDate, title)
	},
}

func init() {
	rootCmd.AddCommand(newCmd)
}
