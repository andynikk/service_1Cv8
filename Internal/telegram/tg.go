package telegram

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

type TgEmoji struct {
	ID      string
	Сaption string
}

type TgMsg struct {
	ID      int64
	API     string
	Message string
	Emoji   TgEmoji
}

type Client struct {
	bot *tgbotapi.BotAPI
}

func New(apiKey string) *Client {
	bot, err := tgbotapi.NewBotAPI(apiKey)
	if err != nil {
		return nil
	}

	amp := &Client{
		bot: bot,
	}

	return amp
}

func (c *Client) SendMessage(text string, chatId int64) error {
	msg := tgbotapi.NewMessage(chatId, text)
	msg.ParseMode = "Markdown"
	_, err := c.bot.Send(msg)

	return err
}
