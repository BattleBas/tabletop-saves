package cmd

import (
	"fmt"
	"os"

	"tabletop-saves/tts"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "tabletop-saves [filename]",
	Short: "Tabletop Simulator backups your game to a zip file",
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		// Do Stuff Here
		fmt.Println(args[0])
		err := tts.Backup(args[0])
		if err != nil {
			fmt.Println(err)
		}
	},
}

// Execute executes the root command
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
