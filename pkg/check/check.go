package check

import (
	"context"
	"io/ioutil"
	"log"
	"strings"
	"time"

	"github.com/google/go-github/v32/github"
	"github.com/spf13/viper"
	"golang.org/x/oauth2"
)

type CheckClient struct {
	C         *github.Client
	CommitSHA string
	Repo      string
	Owner     string
	PrNum     string
	Title     string
	Timeout   int
}

type CIResult struct {
	Repo      string `json:"repo"`
	Owner     string `json:"owner"`
	Total     int    `json:"total"`
	Success   int    `json:"success"`
	Complete  int    `json:"complete"`
	PrNum     string `json:"pr"`
	Title     string `json:"title"`
	CommitSHA string `json:"sha"`
}

func NewClient(path string) (*CheckClient, error) {
	b, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}
	token := strings.TrimSpace(string(b))
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: token},
	)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	tc := oauth2.NewClient(ctx, ts)

	return &CheckClient{
		C:         github.NewClient(tc),
		CommitSHA: viper.GetString("commit"),
		Owner:     viper.GetString("owner"),
		Repo:      viper.GetString("repo"),
		Timeout:   viper.GetInt("timeout"),
		PrNum:     viper.GetString("pr"),
		Title:     viper.GetString("title"),
	}, nil
}

func (c *CheckClient) check() (*CIResult, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	listOption := &github.ListOptions{PerPage: 20}
	var statuses []*github.RepoStatus
	for {
		pagedStatus, r, err := c.C.Repositories.ListStatuses(ctx, c.Owner, c.Repo, c.CommitSHA, listOption)
		if err != nil {
			return nil, err
		}
		log.Printf("github api result, the next page is %v", r.NextPage)

		statuses = append(statuses, pagedStatus...)
		if r.NextPage == 0 {
			break
		}
		listOption.Page = r.NextPage
	}

	total := 0
	waitingCount := 0
	successCount := 0
	uniqSatuses := make(map[string]*github.RepoStatus)
	for _, status := range statuses {
		if _, ok := uniqSatuses[status.GetContext()]; !ok {
			uniqSatuses[status.GetContext()] = status
		}
	}
	for _, v := range uniqSatuses {
		total += 1
		if v.GetState() == "success" {
			successCount += 1
		}
		if v.GetState() == "pending" {
			waitingCount += 1
		}
	}

	return &CIResult{
		Repo:      c.Repo,
		Owner:     c.Owner,
		Total:     total,
		Success:   successCount,
		Complete:  total - waitingCount,
		Title:     c.Title,
		PrNum:     c.PrNum,
		CommitSHA: c.CommitSHA}, nil
}

func (c *CheckClient) QueryLoop() *CIResult {
	start := time.Now()
	for {
		now := time.Now()
		if now.After(start.Add(time.Second * time.Duration(c.Timeout))) {
			break
		}

		r, err := c.check()
		if err != nil {
			log.Fatal(err)
			return nil
		}
		log.Printf("check result: %v", r)

		if r.Complete == r.Total {
			return r
		}
		time.Sleep(time.Second * 15)

	}

	return nil
}
