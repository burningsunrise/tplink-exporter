package parser

import (
	"io/ioutil"

	"github.com/burningsunrise/tplink-exporter/config"

	log "github.com/sirupsen/logrus"
	"gopkg.in/yaml.v2"
)

type YamlConfig struct {
	User     string   `yaml:"user"`
	Password string   `yaml:"password"`
	Devices  []string `yaml:",flow"`
}

func (y *YamlConfig) GetConfig() *YamlConfig {
	yamlFile, err := ioutil.ReadFile("config.yaml")
	if err != nil {
		log.Printf("yamlFile.Get err   #%v ", err)
	}
	err = yaml.Unmarshal(yamlFile, &y)
	if err != nil {
		log.Fatalf("Unmarshal: %v", err)
	}
	if len(y.Devices) <= 0 {
		log.WithFields(log.Fields{
			"devices": "missing",
		}).Fatal("you must add devices in a list in a yaml file")
	}
	if y.Password == "" || y.User == "" {
		log.WithFields(log.Fields{
			"yaml": "incomplete",
		}).Info("password or username missing from yaml file, trying .env")
		y.Password = config.Config("PASSWORD")
		y.User = config.Config("USER")
		if y.Password == "" || y.User == "" {
			log.WithFields(log.Fields{
				"credentials": "missing",
			}).Fatal("you must either add a user / password key in a yaml file or add USER and PASSWORD env viarables")
		}
	}
	return y
}
