package cmd

import (
	"github.com/spf13/cobra"
)

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Show rad version",
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Println(Version)
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}
