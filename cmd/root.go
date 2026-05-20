package cmd

import (
	"fmt"
	"os"

	"diary/pkg/config"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "diary",
	Short: "An encrypted CLI diary",
	Long:  `A feature-full encrypted diary that stores your entries securely.`,
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		if err := config.InitPaths(); err != nil {
			fmt.Fprintf(os.Stderr, "Failed to initialize paths: %v\n", err)
			os.Exit(1)
		}

		if cmd.Name() != "setup" && cmd.Name() != "help" {
			if _, err := os.Stat(config.ConfigFile); os.IsNotExist(err) {
				fmt.Fprintln(os.Stderr, "Diary is not initialized. Please run 'diary setup' to set a global password.")
				os.Exit(1)
			}
		}
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
