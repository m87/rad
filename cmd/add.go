package cmd

import (
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
			err := viper.WriteConfig()
			if err != nil {
				panic("Failed to write config file: " + err.Error())
			}
		},
	}
	cmd.Flags().StringVarP(&url, "url", "u", "", "URL of the radio station")
	cmd.Flags().StringVarP(&alias, "alias", "a", "", "Alias for the radio station")

	return cmd
}

func init() {
	rootCmd.AddCommand(NewAddCmd())
}
