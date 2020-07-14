package cmd

import "github.com/spf13/cobra"

var notifyCmd = &cobra.Command{
	Use:   "notify",
	Short: "notify from stdin to wechat/telegram",
	Run: func(cmd *cobra.Command, args []string) {

	},
}

func init() {
	rootCmd.AddCommand(notifyCmd)
}
