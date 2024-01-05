package main

import (
	"github.com/MoraGames/social_pub/types"
	"github.com/sirupsen/logrus"
)

func run(utils types.Utils, data types.Data) {
	utils.Logger.Info("Running bot")

	for _, update := range data.Updates {
		utils.Logger.WithFields(logrus.Fields{
			"update": update,
		}).Debug("Update received")

		if update.IsMessage() {
			if update.IsCommand() {
				utils.Logger.WithFields(logrus.Fields{
					"command": update.Command(),
				}).Debug("Command received")

				switch update.Command() {
				case "start":
					utils.Logger.WithFields(logrus.Fields{
						"chat": update.Message.Chat,
					}).Debug("Start command received")
				}
			}
		}
	}
}
