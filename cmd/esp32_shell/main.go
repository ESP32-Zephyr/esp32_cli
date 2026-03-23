package main

import (

	"github.com/ESP32-Zephyr/esp32_cli/internal/command"
)

func main() {
	command.Init()
	command.Execute()
}
