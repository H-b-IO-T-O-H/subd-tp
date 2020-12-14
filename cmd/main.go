package main

import (
	yamlConfig "github.com/rowdyroad/go-yaml-config"
)

func main() {
	var config api.Config
	_ = yamlConfig.LoadConfig(&config, "configs/config.yaml", nil)
	app := api.NewApp(config)
	defer app.Close()
	app.Run()
}
