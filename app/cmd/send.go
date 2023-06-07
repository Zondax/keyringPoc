package cmd

import (
	"github.com/zondax/keyringPoc/app/client/tx"

	"github.com/spf13/cobra"
)

// sendCmd represents the send command
var sendCmd = &cobra.Command{
	Use:   "send",
	Short: "send tx",
	Long:  `Send ATOM`,
	Run: func(cmd *cobra.Command, args []string) {
		uid, _ := cmd.Flags().GetString("uid")
		plugin, _ := cmd.Flags().GetString("plugin")
		to, _ := cmd.Flags().GetString("to")
		amount, _ := cmd.Flags().GetString("amount")
		node, _ := cmd.Flags().GetString("node")
		tx.Send(uid, plugin, to, amount, node)
	},
}

func init() {
	rootCmd.AddCommand(sendCmd)
	sendCmd.Flags().StringP("to", "t", "", "address to send")
	sendCmd.MarkFlagRequired("to")
	sendCmd.Flags().StringP("amount", "a", "", "amount to send")
	sendCmd.MarkFlagRequired("amount")
	sendCmd.PersistentFlags().StringVarP(&uid, "uid", "u", "", "uid to store the key by")
	sendCmd.MarkPersistentFlagRequired("uid")
}
