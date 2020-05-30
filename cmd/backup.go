package cmd

import (
	"fmt"
	"tabletop-saves/tts"

	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(backupCmd)
}

var backupCmd = &cobra.Command{
	Use:   "backup [filename]",
	Short: "Tabletop Simulator backups your game to a zip file",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Backup command")
		err := tts.Backup(args[0])
		if err != nil {
			fmt.Println(err)
		}
	},
}
