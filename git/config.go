package git

import (
	"time"
)

type CommitHash struct {
	Url  string
	Ref  string
	Path string
	Hash string 
}

type AuthorConfig struct {
	Name  string
	Email string
}

type CommitSignatureConfig struct {
	Key        string
	Passphrase string
}

type GitAuth struct {
	SshKeyPath     string `yaml:"ssh_key_path"`
	KnownHostsPath string `yaml:"known_hosts_path"`
	User           string
}

type GitConfig struct {
	Url                string
	Ref                string
	Path               string
	Auth               GitAuth
	GpgPublicKeysPaths []string              `yaml:"gpg_public_keys_paths"`
	Author             AuthorConfig
	CommitSignature    CommitSignatureConfig `yaml:"commit_signature"`
	CommitMessage      string                `yaml:"commit_message"`
	PushRetries        int64                 `yaml:"push_retries"`
	PushRetryInterval  time.Duration         `yaml:"push_retry_interval"`
}

func (conf *GitConfig) IsDefined() bool {
	return conf.Url != ""
}