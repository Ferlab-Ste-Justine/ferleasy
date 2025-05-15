package config

import (
	"errors"
	"fmt"
	"io/ioutil"
	"time"

	"github.com/Ferlab-Ste-Justine/ferleasy/state"
	"github.com/Ferlab-Ste-Justine/ferleasy/store"

	ferleaseconf "github.com/Ferlab-Ste-Justine/ferlease/config"
	yaml "gopkg.in/yaml.v2"
)

type ConfigSync struct {
	Author            ferleaseconf.AuthorConfig
	CommitMessage     string                       `yaml:"commit_message"`
	PushRetries       int64                        `yaml:"push_retries"`
	PushRetryInterval time.Duration                `yaml:"push_retry_interval"`
	Orchestrations    []ferleaseconf.Orchestration
	StateStore        store.StoreConfig            `yaml:"state_store"`
}

type Config struct {
	ReleasesStore store.StoreConfig `yaml:"releases_store"`
	EntryPolicy   state.EntryPolicy `yaml:"entry_policy"`
	Sync          ConfigSync
}


func GetConfig(path string) (*Config, error) {
	var c Config

	b, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("Error reading the configuration file: %s", err.Error()))
	}
	err = yaml.Unmarshal(b, &c)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("Error parsing the configuration file: %s", err.Error()))
	}

	return &c, nil
}