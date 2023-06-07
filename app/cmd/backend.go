package cmd

import (
	"github.com/spf13/cobra"
	"github.com/zondax/keyringPoc/app/client/backend"
)

// backendCmd represents the backend command
var backendCmd = &cobra.Command{
	Use:   "backend",
	Short: "get keyring backend",
	Long:  `get which backend the keyring is using.`,
	Run: func(cmd *cobra.Command, args []string) {
		plugin, _ := cmd.Flags().GetString("plugin")
		backend.Backend(plugin)
	},
}

func init() {
	rootCmd.AddCommand(backendCmd)
}
