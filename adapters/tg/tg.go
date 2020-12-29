package tg

import (
	"bytes"
	"context"
	"fmt"
	"text/template"

	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	"github.com/pkg/errors"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/zcong1993/notifiers/types"
)

// Client is ding talk notifier client
type Client struct {
	token    string
	tgClient *tgbotapi.BotAPI
	toID     int64
	logger   log.Logger
}

// NewClient construct a ding talk notifier client
func NewClient(token string, toID int64, logger log.Logger) (*Client, error) {
	if logger == nil {
		logger = log.NewNopLogger()
	}
	tgClient, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		return nil, errors.Wrap(err, "init tg bot")
	}
	client := &Client{
		token:    token,
		tgClient: tgClient,
		toID:     toID,
		logger:   logger,
	}
	level.Info(client.logger).Log("msg", fmt.Sprintf("tg api username is %s\n", tgClient.Self.UserName))
	return client, nil
}

var msgTpl = template.Must(template.New("telegram").Parse(`
*{{ .Title }}*

{{ if gt (len .URL) 0 }}
URL: [{{ .URL }}]({{ .URL }})
{{ end }}
`))

// Notify impl notifier's notify method
func (tc *Client) Notify(_ context.Context, msg *types.Message) error {
	tgMsg, err := buildMsg(msg, tc.toID)
	if err != nil {
		return errors.Wrap(err, "build msg")
	}
	_, err = tc.tgClient.Send(tgMsg)
	return errors.Wrap(err, "send tg msg")
}

func (tc *Client) GetName() string {
	return "telegram"
}

func (tc *Client) Close() error {
	tc.tgClient.StopReceivingUpdates()
	return nil
}

func buildMsg(msg *types.Message, toID int64) (*tgbotapi.MessageConfig, error) {
	var res bytes.Buffer
	err := msgTpl.Execute(&res, msg)
	if err != nil {
		return nil, err
	}

	tgMsg := tgbotapi.NewMessage(toID, res.String())
	tgMsg.ParseMode = "Markdown"

	return &tgMsg, nil
}

var _ types.Notifier = (*Client)(nil)
