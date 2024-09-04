package main

import (
	"encoding/json"
	"goBot/goUnits/logger/logger"
	"os"
	"strconv"
	"strings"

	tgbotapi "github.com/ijnkawakaze/telegram-bot-api"
	//tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

type Config struct {
	Token    string `json:"token"`
	Loglevel int    `json:"loglevel"`
}

var Bot, Err = tgbotapi.NewBotAPI(getToken("config.json"))

func verifiedUser(uid, gid int64, gname string) bool {
	USerconfig := tgbotapi.ChatConfigWithUser{
		ChatID: gid,
		UserID: uid,
	}
	chatMenberConfig := tgbotapi.GetChatMemberConfig{
		USerconfig,
	}

	getChatMenber, err := Bot.GetChatMember(chatMenberConfig)

	if getChatMenber.Status == "creator" || getChatMenber.Status == "administrator" {
		return true
	}

	if err != nil {

		logger.Error("Get chat error: %s \n ChatId: %d \n UserId : %d \n gname : %s \n", err, gid, uid, gname)
		//fmt.Printf("chatmenber :%s", getChatMenber)
	}
	return false
}
func getToken(file string) (token string) {

	configFile, err := os.Open(file)
	defer configFile.Close()
	var config Config
	decoder := json.NewDecoder(configFile)
	if err := decoder.Decode(&config); err != nil {
		logger.Error("Error decoding config file: %s", err)
		return ""
	}

	if err != nil {
		logger.Error("%s", err)
	}
	token = config.Token
	logger.Debug("%s", token)
	return token
}
func main() {

	if Err != nil {
		logger.Error("%s", Err)
	}
	logger.SetLogLevel(1)
	Bot.Debug = true

	logger.Info("Authorized on account %s", Bot.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 600

	updates := Bot.GetUpdatesChan(u)
	for update := range updates {
		if update.Message == nil { // 忽略任何非消息更新
			continue
		} else {
			logger.Error("Error in getting update!")
		}

		// 检查是否为私聊，并且命令是 /ban
		if update.Message.Chat.IsPrivate() && strings.HasPrefix(update.Message.Text, "/ban") {
			args := strings.Split(update.Message.Text, " ")
			if len(args) != 3 {
				msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Usage: /ban <group_id> <user_id>")
				Bot.Send(msg)

				continue
			}

			groupIDStr := args[1]
			userIDStr := args[2]
			groupID, err := strconv.ParseInt(groupIDStr, 10, 64)
			if err != nil {
				msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Invalid group_id format.")
				Bot.Send(msg)
				continue
			}

			userID, err := strconv.ParseInt(userIDStr, 10, 64)
			if err != nil {
				msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Invalid user_id format.")
				Bot.Send(msg)
				continue
			}
			chatConfig := tgbotapi.ChatInfoConfig{
				ChatConfig: tgbotapi.ChatConfig{ChatID: groupID},
			}
			groupNameSrt, _ := Bot.GetChat(chatConfig)
			checkID := update.Message.Chat.ID
			logger.Info("CheckID:%d", checkID)
			if verifiedUser(checkID, groupID, groupNameSrt.Title) == true {

				// 执行封禁操作
				/* kickConfig := tgbotapi.KickChatMemberConfig{
					ChatMemberConfig: tgbotapi.ChatMemberConfig{
						ChatID: groupID,
						UserID: userID, // 进行类型转换
					},
				} */

				_, err = Bot.BanChatMember(groupID, userID)
				if err != nil {
					msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Failed to ban user: "+err.Error())
					Bot.Send(msg)
					logger.Error("%s", err.Error())
				} else {
					msg := tgbotapi.NewMessage(update.Message.Chat.ID, "User banned successfully!")
					Bot.Send(msg)
				}
			} else {
				msg := tgbotapi.NewMessage((update.Message.Chat.ID), "You have no access to action")
				Bot.Send(msg)
				logger.Error("Found access dinded ID:%d", update.Message.From.ID)
			}
		}
		if update.Message.Chat.IsPrivate() && strings.HasPrefix(update.Message.Text, "/help") {
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Usage: /ban <group_id> <user_id>")
			Bot.Send(msg)
		}

	}
}
