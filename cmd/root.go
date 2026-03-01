package cmd

import (
	"fmt"
	"os"

	"github.com/m87/rad/radio"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var cfgFile string
var url string

var rootCmd = &cobra.Command{
	Use:   "rad",
	Short: "Online radio player for the terminal",
	Run: func(cmd *cobra.Command, args []string) {
		url = args[0]
		radio := radio.NewRadio(url)
		err := radio.Play()
		if err != nil {
			fmt.Fprintln(os.Stderr, "Error playing radio:", err)
			os.Exit(1)
		}
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
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

func initConfig() {
	if cfgFile != "" {
		viper.SetConfigFile(cfgFile)
	} else {
		home, err := os.UserHomeDir()
		cobra.CheckErr(err)

		viper.AddConfigPath(home)
		viper.SetConfigType("yaml")
		viper.SetConfigName(".rad")
	}

	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err == nil {
		fmt.Fprintln(os.Stderr, "Using config file:", viper.ConfigFileUsed())
	}
}
