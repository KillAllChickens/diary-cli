package cmd

import (
	"fmt"
	"os"
	"sort"
	"time"

	"diary/pkg/diary"
	"github.com/spf13/cobra"
)

var statsCmd = &cobra.Command{
	Use:   "stats",
	Short: "Show diary statistics and streaks",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Decrypting entries for stats...")
		entries, err := decryptAllEntries()
		if err != nil {
			fmt.Fprintln(os.Stderr, "Failed to decrypt entries:", err)
			return
		}

		if len(entries) == 0 {
			fmt.Println("No entries found.")
			return
		}

		totalEntries := len(entries)
		totalWords := 0
		dateMap := make(map[string]bool)

		for _, e := range entries {
			totalWords += e.WordCount
			dateMap[e.Entry.Date] = true
		}

		var dates []time.Time
		for dateStr := range dateMap {
			t, err := time.Parse(diary.DateFormat, dateStr)
			if err == nil {
				dates = append(dates, t)
			}
		}

		sort.Slice(dates, func(i, j int) bool {
			return dates[i].After(dates[j])
		})

		currentStreak := 0
		if len(dates) > 0 {
			today := time.Now().Truncate(24 * time.Hour)
			lastEntry := dates[0].Truncate(24 * time.Hour)

			if today.Sub(lastEntry) <= 24*time.Hour {
				currentDate := lastEntry
				currentStreak = 1
				for i := 1; i < len(dates); i++ {
					prevDate := dates[i].Truncate(24 * time.Hour)
					if currentDate.Sub(prevDate) == 24*time.Hour {
						currentStreak++
						currentDate = prevDate
					} else if currentDate.Equal(prevDate) {
						continue
					} else {
						break
					}
				}
			}
		}

		fmt.Println("\n--- Diary Statistics ---")
		fmt.Printf("Total Entries : %d\n", totalEntries)
		fmt.Printf("Total Words   : %d\n", totalWords)
		fmt.Printf("Current Streak: %d days\n", currentStreak)
		if totalEntries > 0 {
			fmt.Printf("Avg Words/Entry: %d\n", totalWords/totalEntries)
		}
		fmt.Println("------------------------")
	},
}

func init() {
	rootCmd.AddCommand(statsCmd)
}
