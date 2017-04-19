package main

import (
	"github.com/reechou/robot-auth/config"
	"github.com/reechou/robot-auth/controller"
)

func main() {
	controller.NewLogic(config.NewConfig()).Run()
}
