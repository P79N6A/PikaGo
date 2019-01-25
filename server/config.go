package server

import (
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
)

var ConfigDir string

const ServiceConfigFile = "service_info.yml"

var ServiceConf ServiceConfig

type ServiceConfig struct {
	ServiceName string `yaml:"ServiceName"`
	ServicePort string `yaml:"ServicePort"`
}

func InitConfig() {
	ConfigDir = os.Getenv("CONF_DIR")
	LoadServiceConfig()
}

func LoadServiceConfig() {
	configFile, err := ioutil.ReadFile(filepath.Join(ConfigDir, ServiceConfigFile))
	if err != nil {
		log.Fatalf("ReadFile, err= %v", err)
		return
	}
	yaml.Unmarshal(configFile, &ServiceConf)
	log.Printf("Conf=%v", ServiceConf)
}
