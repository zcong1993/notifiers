package ding

import (
	"bytes"
	"errors"
	"html/template"
	"net/http"

	"github.com/imroc/req"
	"github.com/zcong1993/notifiers/types"
)

const endPoint = "https://oapi.dingtalk.com/robot/send?access_token="

// Client is ding talk notifier client
type Client struct {
	token      string
	httpClient *req.Req
}

// ActionCard is ding talk msg sub struct
type ActionCard struct {
	Title          string `json:"title"`
	Text           string `json:"text"`
	HideAvatar     string `json:"hideAvatar"`
	BtnOrientation string `json:"btnOrientation"`
	SingleTitle    string `json:"singleTitle"`
	SingleURL      string `json:"singleURL"`
}

// DingTalkMsg is ding talk msg struct
type RequestMsg struct {
	Msgtype    string `json:"msgtype"`
	ActionCard `json:"actionCard"`
}

// DingTalkResp is ding talk resp struct
type Resp struct {
	ErrMsg  string `json:"errmsg"`
	ErrCode int    `json:"errcode"`
}

var msgTpl = template.Must(template.New("ding").Parse(`
### {{ .Title }}

{{ .Content }}

{{ if gt (len .Tags) 0 }}
{{ range $tag := .Tags }}
- {{ $tag }}
{{ end }}
{{ end }}
`))

// NewClient construct a ding talk notifier client
func NewClient(token string) *Client {
	client := &Client{
		token:      token,
		httpClient: req.New(),
	}
	return client
}

// SetHTTPClient can replace http client to yourself
func (dc *Client) SetHTTPClient(hc *http.Client) {
	dc.httpClient.SetClient(hc)
}

// Notify impl notifier's notify method
func (dc *Client) Notify(msg *types.Message) error {
	dingMsg, err := buildMsg(msg)
	if err != nil {
		return err
	}
	res, err := dc.httpClient.Post(endPoint+dc.token, req.BodyJSON(dingMsg))
	if err != nil {
		return err
	}
	var resp Resp
	err = res.ToJSON(&resp)
	if err != nil {
		return err
	}
	if resp.ErrCode != 0 {
		return errors.New(resp.ErrMsg)
	}
	return nil
}

func buildMsg(msg *types.Message) (*RequestMsg, error) {
	var res bytes.Buffer
	err := msgTpl.Execute(&res, msg)
	if err != nil {
		return nil, err
	}

	return &RequestMsg{
		Msgtype: "actionCard",
		ActionCard: ActionCard{
			Title:          msg.Title,
			Text:           res.String(),
			HideAvatar:     "0",
			BtnOrientation: "0",
			SingleTitle:    "阅读全文",
			SingleURL:      msg.URL,
		},
	}, nil
}
