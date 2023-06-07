package cmd

import (
	"github.com/spf13/cobra"
	"github.com/zondax/keyringPoc/app/client/key"
)

// signCmd represents the sign command
var signCmd = &cobra.Command{
	Use:   "sign",
	Short: "sign a message",
	Long:  `Sign the provided message using a key.`,
	Run: func(cmd *cobra.Command, args []string) {
		uid, _ := cmd.Flags().GetString("uid")
		plugin, _ := cmd.Flags().GetString("plugin")
		msg, _ := cmd.Flags().GetString("message")
		key.Sign(uid, plugin, msg)
	},
}

func init() {
	keyCmd.AddCommand(signCmd)
	signCmd.Flags().String("message", "", "Help message for toggle")
	signCmd.MarkFlagRequired("message")
}
