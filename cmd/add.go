package cmd

import (
	"github.com/spf13/cobra"
)

func NewAddCmd() *cobra.Command {
	var (
		url   string
		alias string
	)

	cmd := &cobra.Command{
		Use:     "add",
		Aliases: []string{"a"},
		Short:   "Add a new radio station",
		RunE: func(cmd *cobra.Command, args []string) error {
			return saveStation(Station{Alias: alias, URL: url})
		},
	}
	cmd.Flags().StringVarP(&url, "url", "u", "", "URL of the radio station")
	cmd.Flags().StringVarP(&alias, "alias", "a", "", "Alias for the radio station")
	cmd.MarkFlagRequired("url")
	cmd.MarkFlagRequired("alias")

	return cmd
}

func init() {
	rootCmd.AddCommand(NewAddCmd())
}
