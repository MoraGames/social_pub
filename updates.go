package main

import (
	"github.com/MoraGames/social_pub/types"
	"github.com/sirupsen/logrus"
)

// Run the core of the bot
func run(utils types.Utils, data types.Data) {
	// Loop over the updates
	for update := range data.Updates {
		//Log Update
		utils.Logger.WithFields(logrus.Fields{
			"updID":   update.UpdateID,
			"updChat": update.FromChat().Title,
			"updUser": update.SentFrom().UserName,
		}).Debug("Update received")

		// Check if the update is a callback query
		if update.CallbackQuery != nil {
			utils.Logger.WithFields(logrus.Fields{
				"callback": update.CallbackQuery.Data,
			}).Debug("CallbackQuery received")

			manageCallbackQueries(utils, data, update)
			continue
		}

		// Check if the update is a message
		if update.Message != nil {
			// Log Message
			utils.Logger.WithFields(logrus.Fields{
				"msgFrom": update.Message.From.UserName,
				"msgText": update.Message.Text,
				"msgTime": update.Message.Time().Format(utils.TimeFormat),
			}).Info("Message received")

			// Check if the message is a command (and ignore other actions)
			if update.Message.IsCommand() {
				manageCommands(utils, data, update)
				continue
			}

			manageMessages(utils, data, update)
		}
	}
}
