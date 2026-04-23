package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
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
		Run: func(cmd *cobra.Command, args []string) {
			stations := viper.GetStringMap("stations")
			if stations == nil {
				stations = make(map[string]interface{})
			}
			stations[alias] = url
			viper.Set("stations", stations)
			if err := viper.WriteConfig(); err != nil {
				fmt.Fprintln(os.Stderr, "Failed to write config file:", err)
				os.Exit(1)
			}
		},
	}
	cmd.Flags().StringVarP(&url, "url", "u", "", "URL of the radio station")
	cmd.Flags().StringVarP(&alias, "alias", "a", "", "Alias for the radio station")
	_ = cmd.MarkFlagRequired("url")
	_ = cmd.MarkFlagRequired("alias")

	return cmd
}

func init() {
	rootCmd.AddCommand(NewAddCmd())
}
