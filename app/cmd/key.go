package cmd

import (
	"github.com/spf13/cobra"
)

var mnemonic string

// keyCmd represents the key command
var keyCmd = &cobra.Command{
	Use:   "key",
	Short: "key functionalities",
	Long: `different key functionalities as:
		add
		get
		sign`,
	Run: func(cmd *cobra.Command, args []string) {},
}

func init() {
	rootCmd.AddCommand(keyCmd)
	keyCmd.PersistentFlags().StringVar(&mnemonic, "mnemonic", "", "")
}
