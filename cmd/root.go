package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/m87/rad/radio"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var cfgFile string
var url string

var rootCmd = &cobra.Command{
	Use:   "rad",
	Short: "Online radio player for the terminal",
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		err := play(args[0])
		if err != nil {
			fmt.Fprintln(os.Stderr, "Error playing radio:", err)
			os.Exit(1)
		}
	},
}

func Execute() {
	if len(os.Args) > 1 && shouldPlayDirectly(os.Args[1]) {
		initConfig()
		err := play(os.Args[1])
		if err != nil {
			fmt.Fprintln(os.Stderr, "Error playing radio:", err)
			os.Exit(1)
		}
		return
	}

	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func play(input string) error {
	url := input

	stations := viper.GetStringMap("stations")
	if strings.HasPrefix(input, "@") {
		if stations == nil {
			return fmt.Errorf("no stations found in config")
		}

		alias := strings.TrimPrefix(input, "@")
		urlValue, ok := stations[alias]
		if !ok {
			return fmt.Errorf("station alias not found: %s", alias)
		}

		resolvedURL, ok := urlValue.(string)
		if !ok {
			return fmt.Errorf("station alias %s has invalid value type", alias)
		}
		url = resolvedURL
	} else if stations != nil {
		if urlValue, ok := stations[input]; ok {
			resolvedURL, ok := urlValue.(string)
			if !ok {
				return fmt.Errorf("station alias %s has invalid value type", input)
			}
			url = resolvedURL
		}
	}

	r := radio.NewRadio(url)
	return r.Play()
}

func shouldPlayDirectly(firstArg string) bool {
	if firstArg == "" || strings.HasPrefix(firstArg, "-") {
		return false
	}

	for _, command := range rootCmd.Commands() {
		if command.Name() == firstArg {
			return false
		}
		for _, alias := range command.Aliases {
			if alias == firstArg {
				return false
			}
		}
	}

	return true
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

		defaultConfigPath := filepath.Join(home, ".rad.yaml")
		if _, err := os.Stat(defaultConfigPath); os.IsNotExist(err) {
			err = os.WriteFile(defaultConfigPath, []byte("{}\n"), 0o644)
			cobra.CheckErr(err)
		}

		viper.SetConfigFile(defaultConfigPath)

	}

	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err == nil {
		fmt.Fprintln(os.Stderr, "Using config file:", viper.ConfigFileUsed())
	}
}
