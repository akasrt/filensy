package cmd

import (
	"fmt"
	"os"

	"github.com/akasrt/filensy/internal/config/userconfig"
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
	initDependencies()
	err := rootCmd.Execute()
	if err != nil {
		code, msg := errorx.CLIErrorHandler(err)
		fmt.Fprintln(os.Stderr, msg)
		os.Exit(code)
	}
}

func initDependencies() {
	err := userconfig.Load()
	if err != nil {
		fmt.Printf("Warning: Failed to load user config! Err: %s", err.Error())
	}
}
