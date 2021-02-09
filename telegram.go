package notifiers

import (
	"context"
	"strconv"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/pkg/errors"
)

type Telegram struct {
	token       string
	tgClient    *tgbotapi.BotAPI
	defaultToId int64
}

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

func (tg *Telegram) GetName() string {
	return "telegram"
}

func (tg *Telegram) Close() error {
	tg.tgClient.StopReceivingUpdates()
	return nil
}

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

	tgMsg := tgbotapi.NewMessage(toId, msg.Content)
	tgMsg.ParseMode = tgbotapi.ModeHTML

	_, err := tg.tgClient.Send(tgMsg)
	if err != nil {
		return errors.Wrap(ErrNotify, err.Error())
	}
	return nil
}
