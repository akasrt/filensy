package cmd

import (
	"fmt"
	"os"

	"github.com/akasrt/filensy/internal/util/errorx"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:           "filensy",
	Short:         "filensy: securely save and share files",
	SilenceUsage:  true,
	SilenceErrors: true,
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		code, msg := errorx.CLIErrorHandler(err)
		fmt.Fprintln(os.Stderr, msg)
		os.Exit(code)
	}
}
