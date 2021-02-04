package notifiers

import (
	"context"
	"strconv"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/pkg/errors"
)

type Telegram struct {
	token    string
	tgClient *tgbotapi.BotAPI
}

func NewTelegram(token string) (*Telegram, error) {
	tgClient, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		return nil, errors.Wrap(err, "init tg bot")
	}
	return &Telegram{
		token:    token,
		tgClient: tgClient,
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
	toId, err := strconv.ParseInt(to, 10, 0)
	if err != nil {
		return errors.Wrap(err, "covert to to int64")
	}

	tgMsg := tgbotapi.NewMessage(toId, msg.Content)

	_, err = tg.tgClient.Send(tgMsg)
	if err != nil {
		return errors.Wrap(ErrNotify, err.Error())
	}
	return nil
}
