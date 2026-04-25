package cmd

import (
	"bufio"
	"fmt"
	"net"

	"github.com/m87/rad/radio"
	"github.com/spf13/cobra"
)

var statusCmd = &cobra.Command{
	Use:   "status",
	Short: "Show current track metadata",
	RunE: func(cmd *cobra.Command, args []string) error {
		c, err := net.Dial("unix", radio.SockPath())
		if err != nil {
			return fmt.Errorf("could not connect to rad (is it running?)")
		}
		defer c.Close()

		_, _ = fmt.Fprintln(c, "METADATA")

		br := bufio.NewReader(c)
		line, err := br.ReadString('\n')
		if err != nil {
			return err
		}
		fmt.Print(line)
		return nil
	},
}

func init() {
	rootCmd.AddCommand(statusCmd)
}
