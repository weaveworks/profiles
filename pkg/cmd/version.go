package cmd

import (
	"fmt"

	"github.com/weaveworks/profiles/pkg/version"

	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(versionCmd)
}

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Output the version of PCTL",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println(version.Get())
	},
}
