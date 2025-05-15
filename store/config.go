package store

import (
	"errors"

	etcd "github.com/Ferlab-Ste-Justine/ferleasy/etcd"
	git "github.com/Ferlab-Ste-Justine/ferleasy/git"
	s3 "github.com/Ferlab-Ste-Justine/ferleasy/s3"
)

type StoreConfig struct {
	Etcd etcd.EtcdConfig
	Git git.GitConfig
	S3 s3.S3ClientConfig
}

func GetStore[T any](conf *StoreConfig, key string, lockKey string) (Store[T], error) {
	if conf.Etcd.IsDefined() {
		return &StoreEtcd[T]{
			Key: key,
			LockKey: lockKey,
			Config: conf.Etcd,
		}, nil
	}

	if conf.Git.IsDefined() {
		return &StoreGit[T]{
			Key: key,
			Config: conf.Git,
		}, nil
	}

	if conf.S3.IsDefined() {
		return &StoreS3[T]{
			Key: key,
			Config: conf.S3,
		}, nil
	}

	return &StoreEtcd[T]{}, errors.New("No store was found in the configuration")
}