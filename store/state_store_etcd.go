package store

import (
	"context"
	"errors"
	"fmt"
	"path"

	etcd "github.com/Ferlab-Ste-Justine/ferleasy/etcd"

	"github.com/Ferlab-Ste-Justine/etcd-sdk/client"
	yaml "gopkg.in/yaml.v2"
)

type StoreEtcd[T any] struct {
	Key string
	LockKey string
	Config etcd.EtcdConfig
	client *client.EtcdClient
}

func (store *StoreEtcd[T]) Initialize() error {
	passErr := store.Config.Auth.ResolvePassword()
	if passErr != nil {
		return errors.New(fmt.Sprintf("Error retrieving etcd store password credentials: %s", passErr.Error()))
	}

	cli, cliErr := client.Connect(context.Background(), client.EtcdClientOptions{
		ClientCertPath:    store.Config.Auth.ClientCert,
		ClientKeyPath:     store.Config.Auth.ClientKey,
		CaCertPath:        store.Config.Auth.CaCert,
		Username:          store.Config.Auth.Username,
		Password:          store.Config.Auth.Password,
		EtcdEndpoints:     store.Config.Endpoints,
		ConnectionTimeout: store.Config.ConnectionTimeout,
		RequestTimeout:    store.Config.RequestTimeout,
		RetryInterval:     store.Config.RetryInterval,
		Retries:           store.Config.Retries,
	})
	if cliErr != nil {
		return errors.New(fmt.Sprintf("Error connecting to etcd store: %s", cliErr.Error()))
	}

	store.client = cli
	return nil
}

func (store *StoreEtcd[T]) ReadContent() (T, error) {
	var content T

	keyInfo, err := store.client.GetKey(path.Join(store.Config.Prefix, store.Key), client.GetKeyOptions{})
	if err != nil {
		return content, errors.New(fmt.Sprintf("Error reading etcd store content: %s", err.Error()))
	}

	if !keyInfo.Found() {
		return content, nil
	}

	err = yaml.Unmarshal([]byte(keyInfo.Value), &content)
	if err != nil {
		return content, errors.New(fmt.Sprintf("Error deserializing etcd store content: %s", err.Error()))
	}

	return content, nil
}

func (store *StoreEtcd[T]) WriteContent(content T) error {
	output, err := yaml.Marshal(&content)
	if err != nil {
		return errors.New(fmt.Sprintf("Error serializing etcd store content: %s", err.Error()))
	}

	_, err = store.client.PutKey(path.Join(store.Config.Prefix, store.Key), string(output))
	if err != nil {
		return errors.New(fmt.Sprintf("Error writing etcd store content: %s", err.Error()))
	}

	return nil
}

func (store *StoreEtcd[T]) Cleanup() error {
	closeErr := store.client.Close()
	if closeErr != nil {
		return errors.New(fmt.Sprintf("Error disconnecting from etcd store: %s", closeErr.Error()))
	}

	return nil
}

func (store *StoreEtcd[T]) Lock() error {
	_, _, lockErr := store.client.AcquireLock(client.AcquireLockOptions{
		Key: path.Join(store.Config.Prefix, store.LockKey),
		Ttl: store.Config.LockTtl,
		Timeout: store.Config.RequestTimeout,
		RetryInterval: store.Config.RetryInterval,
	})
	if lockErr != nil {
		return errors.New(fmt.Sprintf("Error acquiring lock on etcd store content: %s", lockErr.Error()))
	}

	return nil
}

func (store *StoreEtcd[T]) ReleaseLock() error {
	relErr := store.client.ReleaseLock(path.Join(store.Config.Prefix, store.LockKey))
	if relErr != nil {
		return errors.New(fmt.Sprintf("Error releasing lock on etcd store content: %s", relErr.Error()))
	}

	return nil
}