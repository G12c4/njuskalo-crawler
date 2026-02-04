package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

var bot *tgbotapi.BotAPI

func InitBot() error {
	token := os.Getenv("TELEGRAM_BOT_TOKEN")
	if token == "" {
		return fmt.Errorf("TELEGRAM_BOT_TOKEN environment variable not set")
	}

	var err error
	bot, err = tgbotapi.NewBotAPI(token)
	if err != nil {
		return fmt.Errorf("failed to initialize telegram bot: %w", err)
	}

	fmt.Printf("Telegram bot authorized as @%s\n", bot.Self.UserName)
	return nil
}

func SendCarMessage(car CarDetail) error {
	if bot == nil {
		return fmt.Errorf("bot not initialized")
	}

	chatIDStr := os.Getenv("TELEGRAM_CHAT_ID")
	if chatIDStr == "" {
		return fmt.Errorf("TELEGRAM_CHAT_ID environment variable not set")
	}

	chatID, err := strconv.ParseInt(chatIDStr, 10, 64)
	if err != nil {
		return fmt.Errorf("invalid TELEGRAM_CHAT_ID: %w", err)
	}

	message := FormatCarMessage(car)

	msg := tgbotapi.NewMessage(chatID, message)
	msg.ParseMode = "HTML"
	msg.DisableWebPagePreview = false

	_, err = bot.Send(msg)
	return err
}

func FormatCarMessage(car CarDetail) string {
	var builder strings.Builder

	builder.WriteString(fmt.Sprintf("<b>%s</b>\n\n", car.Title))
	builder.WriteString(fmt.Sprintf("<b>Cijena:</b> %s\n", car.Price))
	builder.WriteString(fmt.Sprintf("<b>Lokacija:</b> %s\n\n", car.Location))

	builder.WriteString("<b>Detalji:</b>\n")
	if car.Year != "" {
		builder.WriteString(fmt.Sprintf("• Godina: %s\n", car.Year))
	}
	if car.Mileage != "" {
		builder.WriteString(fmt.Sprintf("• Kilometraža: %s\n", car.Mileage))
	}
	if car.Gearbox != "" {
		builder.WriteString(fmt.Sprintf("• Mjenjač: %s\n", car.Gearbox))
	}
	if car.Power != "" {
		builder.WriteString(fmt.Sprintf("• Snaga: %s\n", car.Power))
	}
	if car.Engine != "" {
		builder.WriteString(fmt.Sprintf("• Motor: %s\n", car.Engine))
	}
	if car.Type != "" {
		builder.WriteString(fmt.Sprintf("• Tip: %s\n", car.Type))
	}
	if car.ServiceBook != "" {
		builder.WriteString(fmt.Sprintf("• Servisna knjiga: %s\n", car.ServiceBook))
	}

	builder.WriteString(fmt.Sprintf("\n<a href=\"%s\">Pogledaj oglas</a>", car.URL))

	return builder.String()
}

func SendNotification(message string) error {
	if bot == nil {
		return fmt.Errorf("bot not initialized")
	}

	chatIDStr := os.Getenv("TELEGRAM_CHAT_ID")
	if chatIDStr == "" {
		return fmt.Errorf("TELEGRAM_CHAT_ID environment variable not set")
	}

	chatID, err := strconv.ParseInt(chatIDStr, 10, 64)
	if err != nil {
		return fmt.Errorf("invalid TELEGRAM_CHAT_ID: %w", err)
	}

	msg := tgbotapi.NewMessage(chatID, message)
	msg.ParseMode = "HTML"

	_, err = bot.Send(msg)
	return err
}
