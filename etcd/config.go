package etcd

import (
	"time"
)

type EtcdConfig struct {
	Prefix            string
	Endpoints         []string
	ConnectionTimeout time.Duration	`yaml:"connection_timeout"`
	RequestTimeout    time.Duration `yaml:"request_timeout"`
	RetryInterval     time.Duration `yaml:"retry_interval"`
	Retries           uint64
	Auth              Auth
	LockTtl           int64
}

func (conf *EtcdConfig) IsDefined() bool {
	return len(conf.Endpoints) > 0
}