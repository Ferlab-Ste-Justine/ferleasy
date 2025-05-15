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
	return nil
}

func (store *StoreS3[T]) ReadContent() (T, error) {
	var content T

	conn, connErr := s3.Connect(store.Config)
	if connErr != nil {
		return content, connErr
	}

	exists, existsErr := s3.KeyExists(store.Config.Bucket, path.Join(store.Config.Path, "state.yml"), conn)
	if existsErr != nil {
		return content, existsErr
	}

	if !exists {
		return content, nil
	}

	objRead, readErr := conn.GetObject(context.Background(), store.Config.Bucket, path.Join(store.Config.Path, store.Key), minio.GetObjectOptions{})
	if readErr != nil {
		return content, readErr
	}

    data, transfErr := io.ReadAll(objRead)
    if transfErr != nil {
        return content, transfErr
    }

	unmarErr := yaml.Unmarshal(data, &content)
	if unmarErr != nil {
		return content, errors.New(fmt.Sprintf("Error deserializing the content info: %s", unmarErr.Error()))
	}

	return content, nil
}

func (store *StoreS3[T]) WriteContent(content T) error {
	output, err := yaml.Marshal(&content)
	if err != nil {
		return errors.New(fmt.Sprintf("Error serializing the content info: %s", err.Error()))
	}

	conn, connErr := s3.Connect(store.Config)
	if connErr != nil {
		return connErr
	}

	_, putErr := conn.PutObject(
		context.Background(),
		store.Config.Bucket,
		path.Join(store.Config.Path, store.Key),
		bytes.NewReader(output),
		int64(len(output)),
		minio.PutObjectOptions{},
	)

	return putErr
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