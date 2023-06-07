package cmd

import (
	"github.com/spf13/cobra"
	"github.com/zondax/keyringPoc/app/client/key"
)

// addCmd represents the add command
var addCmd = &cobra.Command{
	Use:   "add",
	Short: "add a key",
	Long:  `Creates and storages a new key`,
	Run: func(cmd *cobra.Command, args []string) {
		uid, _ := cmd.Flags().GetString("uid")
		plugin, _ := cmd.Flags().GetString("plugin")
		mnemonic, _ := cmd.Flags().GetString("mnemonic")
		key.Add(uid, plugin, mnemonic)
	},
}

func init() {
	keyCmd.AddCommand(addCmd)
	addCmd.Flags().StringP("mnemonic", "m", "", "mnemonic of a key you'd like to import")
}
