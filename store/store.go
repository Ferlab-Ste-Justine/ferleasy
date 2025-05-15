package store

import (
	"fmt"
)

type Store[T any] interface {
	Initialize() error
	ReadContent() (T, error)
	WriteContent(T) error
	Lock() error
	ReleaseLock() error
	Cleanup() error
}

type StoreScopedFn[T any] func(T) (T, error)

func ProcessStoreContent[T any] (fn StoreScopedFn[T], store Store[T]) error {
	initErr := store.Initialize()
	if initErr != nil {
		return initErr
	}
	defer func() {
		cleanupErr := store.Cleanup()
		if cleanupErr != nil {
			fmt.Printf("Error cleaning up store connection: %s\n", cleanupErr.Error())
		}
	}()

	lockErr := store.Lock()
	if lockErr != nil {
		return lockErr
	}
	defer func() {
		releaseErr := store.ReleaseLock()
		if releaseErr != nil {
			fmt.Printf("Error releasing lock: %s\n", releaseErr.Error())
		}
	}()

	content, readErr := store.ReadContent()
	if readErr != nil {
		return readErr
	}

	var processErr error
	content, processErr = fn(content)
	if processErr != nil {
		return processErr
	}

	writeErr := store.WriteContent(content)
	if writeErr != nil {
		return writeErr
	}

	return nil
}