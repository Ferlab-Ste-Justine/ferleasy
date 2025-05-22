package git

import (
	"time"
)

type AuthorConfig struct {
	Name  string
	Email string
}

type CommitSignatureConfig struct {
	Key        string
	Passphrase string
}

type GitAuth struct {
	SshKey   string `yaml:"ssh_key"`
	KnownKey string `yaml:"known_key"`
	User     string
}

type GitConfig struct {
	Url                string
	Ref                string
	Path               string
	Auth               GitAuth
	AcceptedSignatures string                `yaml:"accepted_signatures"`
	Author             AuthorConfig
	CommitSignature    CommitSignatureConfig `yaml:"commit_signature"`
	CommitMessage      string                `yaml:"commit_message"`
	PushRetries        int64                 `yaml:"push_retries"`
	PushRetryInterval  time.Duration         `yaml:"push_retry_interval"`
}

func (conf *GitConfig) IsDefined() bool {
	return conf.Url != ""
}