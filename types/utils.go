package types

import (
	"github.com/MoraGames/social_pub/config"
	"github.com/sirupsen/logrus"
)

type Utils struct {
	Config     *config.Config
	Logger     *logrus.Logger
	TimeFormat string
}
