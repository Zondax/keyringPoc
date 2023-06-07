/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"github.com/spf13/cobra"
	"github.com/zondax/keyringPoc/app/client/key"
)

// signCmd represents the sign command
var signCmd = &cobra.Command{
	Use:   "sign",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		uid, _ := cmd.Flags().GetString("uid")
		plugin, _ := cmd.Flags().GetString("plugin")
		msg, _ := cmd.Flags().GetString("message")
		key.Sign(uid, plugin, msg)
	},
}

func init() {
	keyCmd.AddCommand(signCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// signCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// signCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	signCmd.Flags().String("message", "", "Help message for toggle")
	signCmd.MarkFlagRequired("message")
}
