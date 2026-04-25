package cmd

import (
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var cfgFile string
var player string
var Version = "dev"

var rootCmd = &cobra.Command{
	Use:   "rad",
	Short: "Online radio player for the terminal",
	Args:  cobra.MinimumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		return play(args[0], player)
	},
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}
func init() {
	cobra.OnInitialize(initConfig)
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.rad.yaml)")
	rootCmd.Flags().StringVarP(&player, "player", "p", "native", "Audio player to use (mpv or native)")
}

func initConfig() {
	if cfgFile != "" {
		viper.SetConfigFile(cfgFile)
	} else {
		home, err := os.UserHomeDir()
		cobra.CheckErr(err)

		defaultConfigPath := filepath.Join(home, ".rad.yaml")
		if _, err := os.Stat(defaultConfigPath); os.IsNotExist(err) {
			err = os.WriteFile(defaultConfigPath, []byte("{}\n"), 0o644)
			cobra.CheckErr(err)
		}

		viper.SetConfigFile(defaultConfigPath)
	}

	viper.AutomaticEnv()
	_ = viper.ReadInConfig()
}
