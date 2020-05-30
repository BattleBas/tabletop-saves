package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "tabletop-saves [filename]",
	Short: "Tabletop Simulator backups your game to a zip file",
	Args:  cobra.MinimumNArgs(1),
}

// Execute executes the root command
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
