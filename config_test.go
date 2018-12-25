package config

import (
	"os"
	"testing"

	log "github.com/sirupsen/logrus"
)

func init() {
	log.SetLevel(log.DebugLevel)
	os.Setenv("CONFIG_FILE", "config.json")
}

func TestReadConfig(t *testing.T) {
	readConfig()
}

func TestBind(t *testing.T) {
	s := new(struct {
		A string `json:"a"`
		C struct {
			Foo string `json:"foo"`
		} `json:"c"`
	})

	Bind(s)

	log.Debug(s.C.Foo)
	if s.C.Foo != "bar" {
		t.Error("Foo is not bar ")
	}
}
