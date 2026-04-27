package cmd

import (
	"bufio"
	"encoding/json"
	"fmt"
	"net"

	"github.com/m87/rad/radio"
	"github.com/spf13/cobra"
)

func NewSaveCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "save",
		Short: "Save current track to library",
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

			status := &radio.Status{}
			if err := json.Unmarshal([]byte(line), status); err != nil {
				return fmt.Errorf("failed to parse metadata: %w", err)
			}
			if err := radio.SaveToLibrary(status.Metadata); err != nil {
				return fmt.Errorf("failed to save track: %w", err)
			}

			return nil
		},
	}
}

func init() {
	rootCmd.AddCommand(NewSaveCmd())
}
