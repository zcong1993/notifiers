package tg

import (
	"bytes"
	"fmt"
	"text/template"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/zcong1993/notifiers/types"
)

// Client is ding talk notifier client
type Client struct {
	token    string
	tgClient *tgbotapi.BotAPI
	toID     int64
}

// NewClient construct a ding talk notifier client
func NewClient(token string, toID int64) *Client {
	tgClient, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		panic(err)
	}
	client := &Client{
		token:    token,
		tgClient: tgClient,
		toID:     toID,
	}
	fmt.Printf("tg api username is %s\n", tgClient.Self.UserName)
	return client
}

var msgTpl = template.Must(template.New("telegram").Parse(`
*{{ .Title }}*

{{ if gt (len .URL) 0 }}
URL: [{{ .URL }}]({{ .URL }})
{{ end }}
`))

// Notify impl notifier's notify method
func (tc *Client) Notify(msg *types.Message) error {
	tgMsg, err := buildMsg(msg, tc.toID)
	if err != nil {
		return err
	}
	_, err = tc.tgClient.Send(tgMsg)
	return err
}

func (tc *Client) GetName() string {
	return "telegram"
}

func buildMsg(msg *types.Message, toID int64) (tgbotapi.MessageConfig, error) {
	var res bytes.Buffer
	err := msgTpl.Execute(&res, msg)
	if err != nil {
		return tgbotapi.MessageConfig{}, err
	}

	tgMsg := tgbotapi.NewMessage(toID, res.String())
	tgMsg.ParseMode = "Markdown"

	return tgMsg, nil
}
