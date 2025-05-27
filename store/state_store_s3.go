package store

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"path"

	"github.com/Ferlab-Ste-Justine/ferleasy/s3"

	minio "github.com/minio/minio-go/v7"
	yaml "gopkg.in/yaml.v2"
)

type StoreS3[T any] struct {
	Key string
	Config s3.S3ClientConfig
}

func (store *StoreS3[T]) Initialize() error {
	return store.Config.Auth.GetKeyAuth()
}

func (store *StoreS3[T]) ReadContent() (T, error) {
	var content T

	conn, connErr := s3.Connect(store.Config)
	if connErr != nil {
		return content, errors.New(fmt.Sprintf("Error creating s3 store client: %s", connErr.Error()))
	}

	exists, existsErr := s3.KeyExists(store.Config.Bucket, path.Join(store.Config.Path, "state.yml"), conn)
	if existsErr != nil {
		return content, errors.New(fmt.Sprintf("Error checking s3 store content existence: %s", existsErr.Error()))
	}

	if !exists {
		return content, nil
	}

	objRead, readErr := conn.GetObject(context.Background(), store.Config.Bucket, path.Join(store.Config.Path, store.Key), minio.GetObjectOptions{})
	if readErr != nil {
		return content, errors.New(fmt.Sprintf("Error getting s3 store content: %s", readErr.Error()))
	}

    data, transfErr := io.ReadAll(objRead)
    if transfErr != nil {
        return content, errors.New(fmt.Sprintf("Error downloading s3 store content: %s", transfErr.Error()))
    }

	unmarErr := yaml.Unmarshal(data, &content)
	if unmarErr != nil {
		return content, errors.New(fmt.Sprintf("Error deserializing s3 store content: %s", unmarErr.Error()))
	}

	return content, nil
}

func (store *StoreS3[T]) WriteContent(content T) error {
	output, err := yaml.Marshal(&content)
	if err != nil {
		return errors.New(fmt.Sprintf("Error serializing s3 store content: %s", err.Error()))
	}

	conn, connErr := s3.Connect(store.Config)
	if connErr != nil {
		return errors.New(fmt.Sprintf("Error creating s3 store client: %s", connErr.Error()))
	}

	_, putErr := conn.PutObject(
		context.Background(),
		store.Config.Bucket,
		path.Join(store.Config.Path, store.Key),
		bytes.NewReader(output),
		int64(len(output)),
		minio.PutObjectOptions{},
	)
	if putErr != nil {
		return errors.New(fmt.Sprintf("Error writing s3 store content: %s", putErr.Error()))
	}

	return nil
}

func (st *StoreS3[T]) Lock() error {
	return nil
}

func (st *StoreS3[T]) ReleaseLock() error {
	return nil
}

func (st *StoreS3[T]) Cleanup() error {
	return nil
}