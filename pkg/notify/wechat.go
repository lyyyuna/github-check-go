package notify

import (
	"fmt"
	"log"
	"net/http"

	"github.com/go-resty/resty/v2"
	"github.com/lyyyuna/github-check-go/pkg/check"
)

type WechatClient struct {
	token   string
	w       *resty.Request
	ci      *check.CIResult
	content string
}

func NewWechatClient(token string, ci *check.CIResult) *WechatClient {
	prlink := fmt.Sprintf("https://github.com/%v/%v/pull/%v", ci.Owner, ci.Repo, ci.PrNum)
	content := fmt.Sprintf(`## CI 检查结束
#### PR 信息
标题: %v
PR 编号: %v
GitHub PR链接: [%v](%v)
#### CI 信息
一共检查 %v 项，通过 %v 项
`, ci.Title, ci.PrNum, prlink, prlink, ci.Total, ci.Success)

	if ci.PrNum == "" {
		content = fmt.Sprintf(`## 合并 CI 检查结束
#### CI 信息
SHA: %v
一共检查 %v 项，通过 %v 项
		`, ci.CommitSHA, ci.Total, ci.Success)
	}

	return &WechatClient{
		token:   token,
		w:       resty.New().SetDisableWarn(true).R(),
		ci:      ci,
		content: content,
	}
}

func (w *WechatClient) Notify() {
	targetUrl := fmt.Sprintf("https://qyapi.weixin.qq.com/cgi-bin/webhook/send?key=%v", w.token)

	type markdown struct {
		Content string `json:"content"`
	}
	req := struct {
		MsgType  string   `json:"msgtype"`
		Markdown markdown `json:"markdown"`
	}{
		MsgType: "markdown",
		Markdown: markdown{
			Content: w.content,
		},
	}

	resp, err := w.w.
		SetBody(req).
		Post(targetUrl)
	if err != nil {
		log.Fatalf("Cannot connect to Wechat, the error is %v", err)
	}
	if resp.StatusCode() != http.StatusOK {
		log.Fatalf("Wechat server response not 200, the response code is: %v, response body is %v", resp.StatusCode(), resp.Body())
	}
}
