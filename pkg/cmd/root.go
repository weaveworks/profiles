package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var (
	rootCmd = &cobra.Command{
		Use:   "PCTL [command]",
		Short: "PCTL is a CLI for interacting with profiles",
	}
)

func init() {
	rootCmd.PersistentFlags().BoolP("help", "h", false, "help for this command")
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
