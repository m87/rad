package cmd

import (
	"bufio"
	"fmt"
	"net"
	"os"

	"github.com/m87/rad/radio"
	"github.com/spf13/cobra"
)

var statusCmd = &cobra.Command{
	Use:   "status",
	Short: "Show current track metadata",
	Run: func(cmd *cobra.Command, args []string) {
		c, err := net.Dial("unix", radio.SockPath())
		if err != nil {
			fmt.Fprintln(os.Stderr, "Could not connect to rad (is it running?)")
			os.Exit(1)
		}
		defer c.Close()

		_, _ = fmt.Fprintln(c, "METADATA")

		br := bufio.NewReader(c)
		line, err := br.ReadString('\n')
		if err != nil {
			return
		}
		fmt.Print(line)
	},
}

func init() {
	rootCmd.AddCommand(statusCmd)
}
