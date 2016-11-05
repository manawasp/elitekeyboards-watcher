package main

import (
	"github.com/BurntSushi/toml"
	log "github.com/Sirupsen/logrus"

	"github.com/manawasp/elitekeyboards-watcher/email"
	kbs "github.com/manawasp/elitekeyboards-watcher/keyboards"
)

const SENDGRID_KEY string = "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"

type AppConfig struct {
	DB          string `toml:"DB"`
	HTML        string `toml:"HTML"`
	URL         string `toml:"URL"`
	SendgridKey string `toml:"SENDGRID_API_KEY"`
}

func main() {
	// Retrieve AppConfig
	var conf AppConfig
	if _, err := toml.DecodeFile("config.toml", &conf); err != nil {
		// handle error
		log.Errorf("Error: Unable to decode config file, %v", err)
	}

	// Get new stats from the website
	keyboards := kbs.WebParse(conf.URL)

	// Load previous stats and compare them
	arr := kbs.Diff(keyboards, kbs.PreviousState(conf.DB))
	if len(arr) > 0 {
		email.Send(SENDGRID_KEY, conf.HTML, arr)
		kbs.Save(conf.DB, keyboards)
	}
}
