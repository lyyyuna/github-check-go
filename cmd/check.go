package cmd

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/lyyyuna/github-check-go/pkg/check"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var checkCmd = &cobra.Command{
	Use:   "check",
	Short: "check if all CI finish",
	Run: func(cmd *cobra.Command, args []string) {
		c, err := check.NewClient(token)
		if err != nil {
			log.Fatalf("Fail to open the github oauth token file: %v", err)
		}
		r := c.QueryLoop()
		b, err := json.Marshal(&r)
		if err != nil {
			log.Fatalf("Fail to encode json: %v", err)
		}
		fmt.Println(string(b))
	},
}

var (
	token     string
	commitSHA string
	owner     string
	repo      string
	timeout   int
)

func init() {
	checkCmd.Flags().StringVarP(&token, "token", "t", "", "specify github token file path")
	checkCmd.Flags().StringVarP(&commitSHA, "commit", "c", "", "specify the commit sha")
	checkCmd.Flags().StringVarP(&owner, "owner", "o", "", "specify the owner name")
	checkCmd.Flags().StringVarP(&repo, "repo", "r", "", "specify the repo name")
	checkCmd.Flags().IntVarP(&timeout, "timeout", "", 60, "specify the timeout")
	viper.BindPFlags(checkCmd.Flags())
	rootCmd.AddCommand(checkCmd)
}
