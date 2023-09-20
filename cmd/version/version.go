package version

import (
	"fmt"

	"github.com/spf13/cobra"
)

var (
	version = "0.0.2"
)

var VersionCommand = &cobra.Command{
	Use:   "version",
	Short: "Show the using version of templatify",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println(version)
	},
}
