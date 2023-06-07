package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

var (
	plugin string
	uid    string
	node   string
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "app",
	Short: "poc of plugins keyring",
	Long:  `Proof of concept for a keyring that uses hashicorp plugins over grpc.`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	// Run: func(cmd *cobra.Command, args []string) { },
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.PersistentFlags().StringVarP(&plugin, "plugin", "p", "goFile", "Provide a plugin name. Currently available plugins [goFile, pyFile]")
	rootCmd.PersistentFlags().StringVarP(&uid, "uid", "u", "", "uid to store the key by")
	rootCmd.MarkPersistentFlagRequired("uid")
	rootCmd.PersistentFlags().StringVarP(&node, "node", "n", "localhost:26657", "node to connect to. Must have the format url:port")
}
