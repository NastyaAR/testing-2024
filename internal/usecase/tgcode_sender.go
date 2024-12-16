package usecase

import (
	"avito-test-task/internal/domain"
	"context"
	"strconv"

	"github.com/go-telegram/bot"
	"go.uber.org/zap"
)

type TelegramCodeSender struct {
	bt *bot.Bot
}

func NewTelegramCodeSender(token string) *TelegramCodeSender {
	bt, _ := bot.New(token)
	return &TelegramCodeSender{bt}
}

func (t *TelegramCodeSender) SendCode(ctx context.Context, user *domain.User, code int, lg *zap.Logger) error {
	params := bot.SendMessageParams{
		ChatID: "1186604465",
		Text:   strconv.Itoa(code),
	}

	msg, err := t.bt.SendMessage(context.Background(), &params)
	if err != nil {
		lg.Warn("cant send code to telegram")
		return err
	}

	pinParams := bot.PinChatMessageParams{
		ChatID:    "1186604465",
		MessageID: msg.ID,
	}

	ok, err := t.bt.PinChatMessage(context.Background(), &pinParams)
	if !ok {
		lg.Warn("cant pin code to telegram")
		return err
	}

	return nil
}
