package cmd

import (
	"fmt"

	"github.com/akasrt/filensy/internal/config/userconfig"
	"github.com/akasrt/filensy/internal/util/errorx"
	"github.com/spf13/cobra"
)

var configCmd = &cobra.Command{
	Use:   "config [key] [value]",
	Short: "configure settings",
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) != 2 {
			return fmt.Errorf("requires exactly 2 arguments: [key] and [value]")
		}
		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		conf := userconfig.GetConfig()
		key := args[0]
		value := args[1]
		switch key {
		case "auth":
			conf.AuthKey = value
		case "dir":
			conf.Directory = value
		default:
			return errorx.ErrInvalidConfigKey
		}

		err := userconfig.SetConfig(conf)
		if err != nil {
			return err
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(configCmd)
}
