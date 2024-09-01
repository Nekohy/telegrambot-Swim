package main

import (
	"goBot/goUnits/logger"
	"strconv"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

func main() {
	// 替换为你的Telegram Bot的API Token
	botToken := "6890025685:AAEeuxYDRNftW5RHfQfvOYir5gle0ZRyq8g"

	bot, err := tgbotapi.NewBotAPI(botToken)
	if err != nil {
		logger.Error("%s", err)
	}
	logger.SetLogLevel(1)
	bot.Debug = true

	logger.Info("Authorized on account %s", bot.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 600

	updates, _ := bot.GetUpdatesChan(u)

	for update := range updates {
		if update.Message == nil { // 忽略任何非消息更新
			continue
		}

		// 检查是否为私聊，并且命令是 /ban
		if update.Message.Chat.IsPrivate() && strings.HasPrefix(update.Message.Text, "/ban") {
			args := strings.Split(update.Message.Text, " ")
			if len(args) != 3 {
				msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Usage: /ban <group_id> <user_id>")
				bot.Send(msg)
				continue
			}

			groupIDStr := args[1]
			userIDStr := args[2]

			groupID, err := strconv.ParseInt(groupIDStr, 10, 64)
			if err != nil {
				msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Invalid group_id format.")
				bot.Send(msg)
				continue
			}

			userID, err := strconv.ParseInt(userIDStr, 10, 64)
			if err != nil {
				msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Invalid user_id format.")
				bot.Send(msg)
				continue
			}

			// 执行封禁操作
			kickConfig := tgbotapi.KickChatMemberConfig{
				ChatMemberConfig: tgbotapi.ChatMemberConfig{
					ChatID: groupID,
					UserID: int(userID), // 进行类型转换
				},
			}

			_, err = bot.KickChatMember(kickConfig)
			if err != nil {
				msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Failed to ban user: "+err.Error())
				bot.Send(msg)
				logger.Error("%s", err.Error())
			} else {
				msg := tgbotapi.NewMessage(update.Message.Chat.ID, "User banned successfully!")
				bot.Send(msg)
			}
		}
		if update.Message.Chat.IsPrivate() && strings.HasPrefix(update.Message.Text, "/help") {
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Usage: /ban <group_id> <user_id>")
			bot.Send(msg)
		}

	}
}
