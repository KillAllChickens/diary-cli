package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var lockCmd = &cobra.Command{
	Use:   "lock",
	Short: "Clear the password cache",
	Long:  `Immediately drops the cached session password so the next diary command will require manual authentication.`,
	Run: func(cmd *cobra.Command, args []string) {
		ClearPasswordCache()
		fmt.Println("Diary locked. Cached password has been cleared.")
	},
}

func init() {
	rootCmd.AddCommand(lockCmd)
}
