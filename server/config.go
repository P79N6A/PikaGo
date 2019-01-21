package server

import (
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"os"
	"path/filepath"
)

var ConfigDir string

const ServiceConfigFile = "service_info.yml"

var Config config

type config struct {
	ServiceConfig ServiceConfig
}

type ServiceConfig struct {
	ServiceName string `yaml:"ServiceName"`
	ServicePort int `yaml:"ServicePort"`
}

func InitConfig() {
	ConfigDir = os.Getenv("CONF_DIR")
	LoadServiceConfig()
}

func LoadServiceConfig() {
	configFile, err := ioutil.ReadFile(filepath.Join(ConfigDir, ServiceConfigFile))
	if err != nil {
		return
	}
	yaml.Unmarshal(configFile, &Config.ServiceConfig)
}
