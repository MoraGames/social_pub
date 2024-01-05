package main

import (
	"encoding/json"
	"io"
	"log"
	"os"
	"time"

	"github.com/MoraGames/social_pub/config"
	"github.com/MoraGames/social_pub/logger"
	"github.com/MoraGames/social_pub/types"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/sirupsen/logrus"
)

func main() {
	//get the configurations
	conf, err := config.NewConfig()
	if err != nil {
		log.Fatalln(err)
	}

	//setup the logger
	l := logger.NewLogger(conf.Log.Type, conf.Log.Format, conf.Log.Level)
	l.WithFields(logrus.Fields{
		"lvl": conf.Log.Level,
		"typ": conf.Log.Type,
		"frm": conf.Log.Format,
	}).Debug("Logger initialized")

	logFile, err := os.OpenFile("files/log.json", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Println(err)
	}

	mw := io.MultiWriter(os.Stdout, logFile)
	l.SetOutput(mw)

	//link Telegram API
	apiToken := os.Getenv("TELEGRAM_API_TOKEN")
	if apiToken == "" {
		l.WithFields(logrus.Fields{
			"env": "TELEGRAM_API_TOKEN",
		}).Panic("Env not set")
	}

	//get the bot API
	bot, err := tgbotapi.NewBotAPI(apiToken)
	if err != nil {
		l.WithFields(logrus.Fields{
			"err": err,
		}).Panic("Error while getting bot API")
	}
	l.WithFields(logrus.Fields{
		"id":       bot.Self.ID,
		"username": bot.Self.UserName,
	}).Info("Account authorized")

	bot.Debug = false
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 180

	//get current time location
	timeLocation, err := time.LoadLocation("Local")
	if err != nil {
		l.WithFields(logrus.Fields{
			"err": err,
		}).Warn("Time location not get (using UTC)")
	}

	updates := bot.GetUpdatesChan(u)
	l.WithFields(logrus.Fields{
		"debugMode": bot.Debug,
		"timeout":   u.Timeout,
	}).Debug("Update channel retreived")

	// Intialize Utils and Data struct
	utils := types.Utils{
		Config:     conf,
		Logger:     l,
		TimeFormat: "15:04:05.000000 MST",
	}
	data := types.Data{
		Bot:     bot,
		Updates: updates,
	}

	// Read from specified files and reload the data into the structs
	ReloadStatus(
		[]types.Reload{
			{FileName: "files/users.json", DataStruct: &Users, IfOkay: nil, IfFail: nil},
		},
		utils,
	)

	run(utils, data)
}

func ReloadStatus(reloads []types.Reload, utils types.Utils) {
	utils.Logger.WithFields(logrus.Fields{
		"reloads": reloads,
	}).Info("Reloading data")

	numOfFail, numOfFailFunc, numOfOkay, numOfOkayFunc := 0, 0, 0, 0
	for _, reload := range reloads {
		hasFailed := false

		utils.Logger.WithFields(logrus.Fields{
			"IfFail() exist": reload.IfFail != nil,
			"IfOkay() exist": reload.IfOkay != nil,
		}).Debug("Reloading " + reload.FileName)

		file, err := os.ReadFile(reload.FileName)
		if err != nil {
			hasFailed = true
			utils.Logger.WithFields(logrus.Fields{
				"file": reload.FileName,
				"err":  err,
			}).Error("Error while reading file")
		}

		if len(file) != 0 {
			err = json.Unmarshal(file, reload.DataStruct)
			if err != nil {
				hasFailed = true
				utils.Logger.WithFields(logrus.Fields{
					"data": reload.DataStruct,
					"err":  err,
				}).Error("Error while unmarshalling data")
			}
		} else {
			hasFailed = true
			utils.Logger.WithFields(logrus.Fields{
				"file": reload.FileName,
			}).Error("File is empty")
		}

		if hasFailed {
			numOfFail++

			utils.Logger.WithFields(logrus.Fields{
				"file": reload.FileName,
			}).Info("Reloading has failed")

			if reload.IfFail != nil {
				numOfFailFunc++
				reload.IfFail(utils)
				utils.Logger.WithFields(logrus.Fields{
					"file": reload.FileName,
				}).Debug("Reload.IfFail() executed")
			}
		} else {
			numOfOkay++
			utils.Logger.WithFields(logrus.Fields{
				"file": reload.FileName,
			}).Info("Reloading has succeed")

			if reload.IfOkay != nil {
				numOfOkayFunc++
				reload.IfOkay(utils)
				utils.Logger.WithFields(logrus.Fields{
					"file": reload.FileName,
				}).Debug("Reload.IfOkay() executed")
			}
		}
	}

	utils.Logger.WithFields(logrus.Fields{
		"fails":     numOfFail,
		"failsFunc": numOfFailFunc,
		"okays":     numOfOkay,
		"okaysFunc": numOfOkayFunc,
		"totat":     len(reloads),
	}).Info("Reloading data completed")
}

func WriteMessage(bot *tgbotapi.BotAPI, chatID int64, replyMessageID int, text string) {
	msg := tgbotapi.NewMessage(chatID, text)
	if replyMessageID != -1 {
		msg.ReplyToMessageID = replyMessageID
	}
	bot.Send(msg)
}
