package cmd

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
)

func stateDir() string {
	home, _ := os.UserHomeDir()
	return filepath.Join(home, ".local", "state", "rad")
}
func sockPath() string { return filepath.Join(stateDir(), "rad.sock") }

var statusCmd = &cobra.Command{
	Use: "status",
	Run: func(cmd *cobra.Command, args []string) {
		c, err := net.Dial("unix", sockPath())
		if err != nil {
			return
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
