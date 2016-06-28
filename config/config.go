package config

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
	"path"
	"path/filepath"
)

const (
	configFileName = "config.json"
)

var ConfigParams = &struct {
	ListenPort int `json:"listenPort"`
}{}

func init() {
	if err := loadConfig(); err != nil {
		log.Fatal(err)
	}
}

var loadConfig = func() error {
	f, err := loadFile()
	if err != nil {
		return err
	}

	err = json.Unmarshal(f, ConfigParams)
	if err != nil {
		return err
	}
	return nil
}

var loadFile = func() (b []byte, err error) {
	// Get excutable directory
	exeDir, e := filepath.Abs(filepath.Dir(os.Args[0]))
	if e != nil {
		log.Fatal(e)
	}

	cfgFilePath := path.Join(exeDir, configFileName)
	if _, err = os.Stat(cfgFilePath); err != nil {
		if os.IsNotExist(err) {
			// create a default config file if not exist
			ConfigParams.ListenPort = 8080 // default port 8080
			b, err = json.Marshal(ConfigParams)
			if err != nil {
				return
			}
			err = ioutil.WriteFile(cfgFilePath, b, 0644)
			return b, err
		}
	}

	b, err = ioutil.ReadFile(cfgFilePath)
	return
}
