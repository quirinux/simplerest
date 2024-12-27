package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"simplerest/libs/server"
	"simplerest/libs/settings"
)

func main() {
	var config_path string
	flag.StringVar(&config_path, "config", "simplerest.toml", "Config file path")
	flag.Parse()

	var full_file_path string
	var err error
	if full_file_path, err = filepath.Abs(config_path); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	full_path := filepath.Dir(full_file_path)
	if err := os.Chdir(full_path); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	_settings, err := settings.Parse(full_file_path)
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
