package main

import (
	"github.com/sirupsen/logrus"
	"github.com/zenorachi/image-box/internal/app"
)

const (
	ConfigDir  = "configs"
	ConfigFile = "main"
	ENVFile    = ".env"
)

func main() {
	if err := app.Run(ConfigDir, ConfigFile, ENVFile); err != nil {
		logrus.Errorln(err)
	}
}
