package config

import (
	"encoding/json"
	"flag"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"reflect"

	log "github.com/sirupsen/logrus"
	yaml "gopkg.in/yaml.v2"
)

var (
	conf *Config
)

type Config struct {
	M    interface{} // user struct
	file *string     // config file path
}

func readConfig() {
	pwd := path.Dir(os.Args[0])

	fle := os.Getenv("CONFIG_FILE")
	conf.file = &fle

	if *conf.file == "" {
		conf.file = flag.String("config", pwd+"/conf.json", "--config <configFile.json>")
		flag.Parse()
	}

	if b, err := ioutil.ReadFile(*conf.file); err != nil {
		log.WithFields(log.Fields{
			"configFile": *conf.file,
			"err":        err,
		}).Fatal("Unable to read config file")
	} else {
		// if conf.M == nil {
		// 	// not binded yet
		// 	log.Debug("Unbinded struct creating new interface")
		// 	conf.M = new(interface{})
		// }
		if filepath.Ext(*conf.file) == ".json" {
			if err := json.Unmarshal(b, &conf.M); err != nil {
				log.WithFields(log.Fields{
					"configFile": *conf.file,
					"err":        err,
				}).Fatal("Unable to Parse json file")
			}
		} else if filepath.Ext(*conf.file) == ".yml" || filepath.Ext(*conf.file) == ".yaml" {
			if reflect.ValueOf(conf.M).Kind() == reflect.Ptr {
				if err := yaml.Unmarshal(b, conf.M); err != nil {
					log.WithFields(log.Fields{
						"configFile": *conf.file,
						"err":        err,
					}).Fatal("Unable to Parse yaml file")
				}
			} else {
				if err := yaml.Unmarshal(b, &conf.M); err != nil {
					log.WithFields(log.Fields{
						"configFile": *conf.file,
						"err":        err,
					}).Fatal("Unable to Parse yaml file")
				}
			}
		} else {
			log.WithFields(log.Fields{
				"configFile": *conf.file,
			}).Fatal("Unknown file extension, supported .json, .yml || .yaml")
		}
		// log.WithFields(log.Fields{"conf.M": conf.M}).Debug()
	}
}

func Bind(st interface{}) {
	conf.M = st
	// Read the config file
	readConfig()
}

func Init() {
	readConfig()
}

func init() {
	conf = new(Config)
}
