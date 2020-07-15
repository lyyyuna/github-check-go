package cmd

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"

	"github.com/lyyyuna/github-check-go/pkg/check"
	"github.com/lyyyuna/github-check-go/pkg/notify"
	"github.com/spf13/cobra"
)

var notifyCmd = &cobra.Command{
	Use:   "notify",
	Short: "notify from stdin to wechat/telegram",
	Run: func(cmd *cobra.Command, args []string) {
		data, err := ioutil.ReadAll(os.Stdin)
		if err != nil {
			log.Fatalln(err)
		}

		var ci check.CIResult
		err = json.Unmarshal(data, &ci)
		if err != nil {
			log.Fatalln(err)
		}

		if wechat != "" {
			w := notify.NewWechatClient(wechat, &ci)
			w.Notify()
		}
	},
}

var (
	wechat string
)

func init() {
	notifyCmd.Flags().StringVarP(&wechat, "wechat", "w", "", "specify the wechat token")
	rootCmd.AddCommand(notifyCmd)
}
