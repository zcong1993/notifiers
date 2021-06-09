package notifiers

import (
	"context"
	"strconv"

	"github.com/k3a/html2text"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/pkg/errors"
)

// Telegram impl notifier, notify msg by telegram bot.
type Telegram struct {
	token       string
	tgClient    *tgbotapi.BotAPI
	defaultToId int64
	NoopWaiter
}

// NewTelegram create a instance.
func NewTelegram(token string, defaultToId int64) (*Telegram, error) {
	tgClient, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		return nil, errors.Wrap(err, "init tg bot")
	}
	return &Telegram{
		token:       token,
		tgClient:    tgClient,
		defaultToId: defaultToId,
	}, nil
}

// GetName impl Notifier.GetName.
func (tg *Telegram) GetName() string {
	return "telegram"
}

// Close impl Notifier.Close.
func (tg *Telegram) Close() error {
	tg.tgClient.StopReceivingUpdates()
	return nil
}

// Notify impl Notifier.Notify.
// If to is not set, will send msg to defaultToId.
// The telegram message parse mode is html mode.
func (tg *Telegram) Notify(ctx context.Context, to string, msg Message) error {
	var toId int64
	if to == "" {
		toId = tg.defaultToId
	} else {
		to, err := strconv.ParseInt(to, 10, 0)
		if err != nil {
			return errors.Wrap(err, "covert to to int64")
		}
		toId = to
	}

	tgMsg := tgbotapi.NewMessage(toId, html2text.HTML2Text(msg.Content))
	//tgMsg.ParseMode = tgbotapi.ModeHTML

	_, err := tg.tgClient.Send(tgMsg)
	if err != nil {
		return errors.Wrap(ErrNotify, err.Error())
	}
	return nil
}

var _ Notifier = (*Telegram)(nil)
