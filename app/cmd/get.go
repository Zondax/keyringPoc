package cmd

import (
	"github.com/zondax/keyringPoc/app/client/key"

	"github.com/spf13/cobra"
)

// getCmd represents the get command
var getCmd = &cobra.Command{
	Use:   "get",
	Short: "get a key",
	Long:  `Retrieves a key previously created with the same keyring`,
	Run: func(cmd *cobra.Command, args []string) {
		uid, _ := cmd.Flags().GetString("uid")
		plugin, _ := cmd.Flags().GetString("plugin")
		key.Get(uid, plugin)
	},
}

func init() {
	keyCmd.AddCommand(getCmd)
}
