package main

import (
	"flag"
	"fmt"
	"os"
	"simplerest/libs/server"
	"simplerest/libs/settings"
)

func main() {
	configPath := flag.String("config", "simplerest.toml", "Config file path")
	flag.Parse()
	_settings, err := settings.Parse(*configPath)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	_settings.Display()

	_server := server.New(_settings)

	if err := _server.Initialize(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	if err := _server.Run(); err != nil {
		fmt.Println(err.Error())
	}
	fmt.Println("exiting")
}
