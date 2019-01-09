package config

import (
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"reflect"

	dotaccess "github.com/go-bongo/go-dotaccess"
	log "github.com/sirupsen/logrus"
	yaml "gopkg.in/yaml.v2"
)

var (
	conf *Config
)

type ConfigMap interface{}

type Config struct {
	M    ConfigMap // user struct
	file *string   // config file path
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
	conf = new(Config)
	conf.M = st
	// Read the config file
	readConfig()
}

func Init() {
	conf = new(Config)
	readConfig()
}

func Unmarshal(key string, st interface{}) error {
	var val interface{}
	var err error

	val, err = dotaccess.Get(conf.M, key)
	if err == nil {
		var b []byte
		b, err = json.Marshal(val)
		if err == nil {
			err = json.Unmarshal(b, st)
			if err == nil {
				fmt.Println(st)
			}
		}
	}
	return err
}

func String(key string) (string, error) {
	val, err := dotaccess.Get(conf.M, key)
	if err != nil {
		return "", err
	}

	v, ok := val.(string)
	if ok {
		return v, nil
	}

	return "", errors.New("Invalid")
}

func Int(key string) (int, error) {
	val, err := dotaccess.Get(conf.M, key)
	if err != nil {
		return 0, err
	}

	v, ok := val.(int)
	if ok {
		return v, nil
	}

	return 0, errors.New("Invalid")
}
